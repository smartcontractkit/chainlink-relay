package ops

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/smartcontractkit/chainlink-relay/ops/adapter"
	"github.com/smartcontractkit/chainlink-relay/ops/chainlink"
	"github.com/smartcontractkit/chainlink-relay/ops/database"
	"github.com/smartcontractkit/chainlink-relay/ops/relay"
	"github.com/smartcontractkit/chainlink-relay/ops/utils"
	"github.com/smartcontractkit/integrations-framework/client"
)

// Deployer interface for deploying contracts
type Deployer interface {
	Load() error                            // upload contracts (may not be necessary)
	DeployLINK() error                      // deploy LINK contract
	DeployOCR() error                       // deploy OCR contract
	TransferLINK() error                    // transfer LINK to OCR contract
	InitOCR(keys []map[string]string) error // initialize OCR contract with provided keys
	OCR2Address() string                    // fetch deployed OCR contract address
}

func New(ctx *pulumi.Context, deployer Deployer) error {
	img := map[string]*utils.Image{}

	// fetch postgres
	img["psql"] = &utils.Image{
		Name: "postgres-image",
		Tag:  "postgres:latest", // always use latest postgres
	}

	// fetch chainlink image
	img["chainlink"] = &utils.Image{
		Name: "chainlink-image",
		Tag:  "public.ecr.aws/chainlink/chainlink:" + config.Require(ctx, "CL-NODE_VERSION"),
	}

	// fetch list of EAs
	eas := []string{}
	if err := config.GetObject(ctx, "EA-NAMES", &eas); err != nil {
		return err
	}
	for _, n := range eas {
		img[n] = &utils.Image{
			Name: n + "-adapter-image",
			Tag:  fmt.Sprintf("public.ecr.aws/chainlink/adapters/%s-adapter:develop-latest", n),
		}
	}

	// pull remote images
	for i := range img {
		if err := img[i].Pull(ctx); err != nil {
			return err
		}
	}

	// validate number of relays
	relayNum := config.GetInt(ctx, "R-COUNT")
	if relayNum < 4 {
		return fmt.Errorf("Minimum number of relays (4) not met (%d)", relayNum)
	}
	// build local image for relay
	img["relay"] = &utils.Image{
		Name: "relay-image",
		Tag:  "relay:" + config.Require(ctx, "R-VERSION"), // always use latest postgres
	}
	if err := img["relay"].Build(ctx); err != nil {
		return err
	}

	// start pg + create DBs
	db, err := database.New(ctx, img["psql"].Img)
	if err != nil {
		return err
	}
	if !ctx.DryRun() {
		// wait for readiness check
		if err := db.Ready(); err != nil {
			return err
		}

		// create DB names
		dbNames := []string{"chainlink", "relay_bootstrap"}
		for i := 0; i < relayNum; i++ {
			dbNames = append(dbNames, fmt.Sprintf("relay_%d", i))
		}

		// create DBs
		for _, n := range dbNames {
			if err := db.Create(n); err != nil {
				return err
			}
		}
	}

	// start cl
	cl, err := chainlink.New(ctx, img["chainlink"].Img, db.Port)
	if err != nil {
		return err
	}
	if !ctx.DryRun() {
		// wait for readiness check
		if err := cl.Ready(); err != nil {
			return err
		}

		// delete all jobs if they exist
		// any related jobs to EIs must be removed for the EI to be replaced later
		if err := cl.DeleteAllJobs(); err != nil {
			return err
		}
	}

	// start adapters + add to CL node
	for i, ea := range eas {
		a, err := adapter.New(ctx, img[ea], i)
		if err != nil {
			return err
		}
		// add to chainlink node
		if !ctx.DryRun() {
			if err := cl.AddBridge(a.Name, a.URL); err != nil {
				return err
			}
		}
	}

	// add webhooks and EA to CL + start relays
	relays := map[string]*relay.Relay{}
	for i := 0; i <= relayNum; i++ {
		indexStr := ""
		if i == 0 {
			indexStr = "bootstrap"
		} else {
			indexStr = strconv.Itoa(i - 1)
		}

		eiSecrets := map[string]string{}
		if !ctx.DryRun() {
			// create EI secrets
			eiSecrets, err = cl.AddEI("relay_"+indexStr, fmt.Sprintf("http://localhost:%d/jobs", config.RequireInt(ctx, "R-PORT-START")+i))
			if err != nil {
				return err
			}

			// create EA endpoints
			if err := cl.AddBridge("relay_"+indexStr, fmt.Sprintf("http://localhost:%d/runs", config.RequireInt(ctx, "R-PORT-START")+i)); err != nil {
				return err
			}

		}

		// start container
		r, err := relay.New(ctx, img["relay"].Local, db.Port, i, eiSecrets)
		if err != nil {
			return err
		}

		relays[indexStr] = &r
	}

	// fetch keys from relays
	if !ctx.DryRun() {
		for k := range relays {
			if err := relays[k].GetKeys(); err != nil {
				return err
			}
		}

		// deploy contracts
		// upload contracts
		if err = deployer.Load(); err != nil {
			return err
		}
		// deploy LINK
		if err = deployer.DeployLINK(); err != nil {
			return err
		}

		// deploy OCR2 contract (w/ dummy access controller addresses)
		if err = deployer.DeployOCR(); err != nil {
			return err
		}

		// transfer tokens to OCR2 contract
		if err = deployer.TransferLINK(); err != nil {
			return err
		}

		// set OCR2 config
		var keys []map[string]string
		for k, v := range relays {
			// skip if bootstrap node
			if k == "bootstrap" {
				continue
			}

			parsedKeys := map[string]string{}
			// remove prefixes if present
			for k, val := range v.Keys {
				parsedKeys[k] = val
				// replace value with val without prefix if prefix exists
				sArr := strings.Split(val, "_")
				if len(sArr) == 2 {
					parsedKeys[k] = sArr[1]
				}
			}

			keys = append(keys, parsedKeys)
		}
		if err = deployer.InitOCR(keys); err != nil {
			return err
		}

		// create job specs
		p2pBootstrap := relays["bootstrap"].Keys["P2PID"] + "@" + relays["bootstrap"].P2P
		i := 0
		for k := range relays {
			name := "relay_" + k

			// if bootstrap, change the other parameters
			bootstrap := "false"
			if k == "bootstrap" {
				bootstrap = "true"
			}

			// create specs + add to CL node
			ea := eas[i%len(eas)]
			msg := utils.LogStatus(fmt.Sprintf("Adding job spec to '%s' with '%s' EA", name, ea))
			spec := &client.WebhookJobSpec{
				Name:      name + " job",
				Initiator: name,
				InitiatorSpec: fmt.Sprintf("{\\\"contractAddress\\\": \\\"%s\\\",\\\"p2pBootstrapPeers\\\": [\\\"%s\\\"],\\\"isBootstrapPeer\\\": %s,\\\"keyBundleID\\\": \\\"%s\\\",\\\"observationTimeout\\\": \\\"10s\\\",\\\"blockchainTimeout\\\": \\\"20s\\\",\\\"contractConfigTrackerSubscribeInterval\\\": \\\"2m\\\",\\\"contractConfigConfirmations\\\": 3}",
					deployer.OCR2Address(),     // contractAddress
					p2pBootstrap,               //p2pBootstrapPeers
					bootstrap,                  //isBootstrapPeer
					relays[k].Keys["OCRKeyID"], // keyBundleID
				),
				ObservationSource: fmt.Sprintf(`
         ea  [type=bridge name=%s requestData=<{"data":{"from":"LINK", "to":"USD"}}>]
         parse [type="jsonparse" path="result"]
         multiply [type="multiply" times=100000000]
         return [type=bridge name="%s" requestData=<{"jobID":$(jobSpec.externalJobID), "result":$(multiply)}>]

         ea -> parse -> multiply -> return`, ea, name),
			}
			_, err = cl.Call.CreateJob(spec)
			if msg.Check(err) != nil {
				return err
			}
			i++
		}
	}

	return nil
}
