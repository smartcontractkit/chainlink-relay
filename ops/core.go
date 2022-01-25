package ops

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/smartcontractkit/chainlink-relay/ops/adapter"
	"github.com/smartcontractkit/chainlink-relay/ops/chainlink"
	"github.com/smartcontractkit/chainlink-relay/ops/database"
	"github.com/smartcontractkit/chainlink-relay/ops/utils"
	"github.com/smartcontractkit/integrations-framework/client"
)

// Deployer interface for deploying contracts
type Deployer interface {
	Load() error                             // upload contracts (may not be necessary)
	DeployLINK() error                       // deploy LINK contract
	DeployOCR() error                        // deploy OCR contract
	TransferLINK() error                     // transfer LINK to OCR contract
	InitOCR(keys []chainlink.NodeKeys) error // initialize OCR contract with provided keys
	Fund(addresses []string) error           // fund the nodes
	OCR2Address() string                     // fetch deployed OCR contract address
	Addresses() map[int]string               // map of all deployed addresses (ocr2, validators, etc)
}

// ObservationSource creates the observation source for the CL node jobs
type ObservationSource func(priceAdapter string) string

// RelayConfig creates the stringified config for the job spec
type RelayConfig func(ctx *pulumi.Context, addresses map[int]string) (map[string]string, error)

func New(ctx *pulumi.Context, deployer Deployer, obsSource ObservationSource, juelsObsSource ObservationSource, relayConfigFunc RelayConfig) error {
	// check these two parameters at the beginning to prevent getting to the end and erroring if they are not present
	chain := config.Require(ctx, "CL-RELAY_NAME")

	img := map[string]*utils.Image{}

	// fetch postgres
	img["psql"] = &utils.Image{
		Name: "postgres-image",
		Tag:  "postgres:latest", // always use latest postgres
	}

	buildLocal := config.GetBool(ctx, "CL-BUILD_LOCALLY")
	if !buildLocal {
		// fetch chainlink image
		img["chainlink"] = &utils.Image{
			Name: "chainlink-remote-image",
			Tag:  "public.ecr.aws/chainlink/chainlink:" + config.Require(ctx, "CL-NODE_VERSION"),
		}
	}
	// TODO: build local chainlink image

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

	// build local chainlink node
	if buildLocal {
		img["chainlink"] = &utils.Image{
			Name: "chainlink-local-build",
			Tag:  "chainlink:local",
		}
		if err := img["chainlink"].Build(ctx, config.Require(ctx, "CL-BUILD_CONTEXT"), config.Require(ctx, "CL-BUILD_DOCKERFILE")); err != nil {
			return err
		}
	}

	// validate number of relays
	nodeNum := config.GetInt(ctx, "CL-COUNT")
	if nodeNum < 4 {
		return fmt.Errorf("Minimum number of chainlink nodes (4) not met (%d)", nodeNum)
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
		dbNames := []string{"chainlink_bootstrap"}
		for i := 0; i < nodeNum; i++ {
			dbNames = append(dbNames, fmt.Sprintf("chainlink_%d", i))
		}

		// create DBs
		for _, n := range dbNames {
			if err := db.Create(n); err != nil {
				return err
			}
		}
	}

	// start EAs
	adapters := []client.BridgeTypeAttributes{}
	for i, ea := range eas {
		a, err := adapter.New(ctx, img[ea], i)
		if err != nil {
			return err
		}
		adapters = append(adapters, a)
	}

	// start chainlink nodes
	nodes := map[string]*chainlink.Node{}
	for i := 0; i <= nodeNum; i++ {
		// start container
		cl, err := chainlink.New(ctx, img["chainlink"], db.Port, i)
		if err != nil {
			return err
		}
		nodes[cl.Name] = &cl // store in map
	}

	if config.GetBool(ctx, "ENV-ONLY_BOOT_CONTAINERS") {
		fmt.Println("ONLY BOOTING CONTAINERS")
		return nil
	}

	if !ctx.DryRun() {
		for _, cl := range nodes {
			// wait for readiness check
			if err := cl.Ready(); err != nil {
				return err
			}

			// delete all jobs if any exist
			if err := cl.DeleteAllJobs(); err != nil {
				return err
			}

			// add adapters to CL node
			for _, a := range adapters {
				if err := cl.AddBridge(a.Name, a.URL); err != nil {
					return err
				}
			}
		}
	}

	if !ctx.DryRun() {
		// fetch keys from relays
		for k := range nodes {
			if err := nodes[k].GetKeys(chain); err != nil {
				return err
			}
		}

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
		var keys []chainlink.NodeKeys
		for k := range nodes {
			// skip if bootstrap node
			if k == "chainlink-bootstrap" {
				continue
			}
			keys = append(keys, nodes[k].Keys)
		}
		if err = deployer.InitOCR(keys); err != nil {
			return err
		}

		// create relay config
		relayConfig, err := relayConfigFunc(ctx, deployer.Addresses())
		if err != nil {
			return err
		}

		// create job specs
		var addresses []string
		i := 0
		for k := range nodes {
			// adding chain node to CL node
			switch relayConfig["nodeType"] {
			case "terra":
				msg := utils.LogStatus(fmt.Sprintf("Adding terra node to '%s'", k))
				attrs := client.TerraNodeAttributes{
					Name:          "Terra Node Localhost",
					TerraChainID:  relayConfig["chainID"],
					TendermintURL: relayConfig["tendermintURL"],
					FCDURL:        relayConfig["fcdURL"],
				}
				_, err = nodes[k].Call.CreateTerraNode(&attrs)
				if msg.Check(err) != nil {
					return err
				}
			default:
				fmt.Printf("WARN: No chain config to add to '%s'\n", k)
			}

			// create specs + add to CL node
			ea := eas[i%len(eas)]
			msg := utils.LogStatus(fmt.Sprintf("Adding job spec to '%s' with '%s' EA", k, ea))

			spec := &client.OCR2TaskJobSpec{
				Name:        "local testing job",
				ContractID:  deployer.OCR2Address(),
				Relay:       chain,
				RelayConfig: relayConfig,
				P2PPeerID:   nodes[k].Keys.P2PID,
				P2PBootstrapPeers: []client.P2PData{
					nodes["chainlink-bootstrap"].P2P,
				},
				IsBootstrapPeer:       k == "chainlink-bootstrap",
				OCRKeyBundleID:        nodes[k].Keys.OCR2KeyID,
				TransmitterID:         nodes[k].Keys.OCR2TransmitterID,
				ObservationSource:     obsSource(ea),
				JuelsPerFeeCoinSource: juelsObsSource(ea),
			}
			_, err = nodes[k].Call.CreateJob(spec)
			if msg.Check(err) != nil {
				return err
			}
			i++

			// retrieve transmitter address for funding
			addresses = append(addresses, nodes[k].Keys.OCR2Transmitter)
		}

		// fund nodes
		if err = deployer.Fund(addresses); err != nil {
			return err
		}
	}

	return nil
}
