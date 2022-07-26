package monitoring

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// Envelope contains data that is required for all the chain.
type Envelope struct {
	// latest transmission details
	ConfigDigest    types.ConfigDigest
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp time.Time

	// latest contract config
	ContractConfig types.ContractConfig

	// extra
	BlockNumber             uint64
	Transmitter             types.Account
	LinkBalance             *big.Int
	LinkAvailableForPayment *big.Int

	// The "fee coin" is different for each chain.
	JuelsPerFeeCoin   *big.Int
	AggregatorRoundID uint32
}

// TxResultsBySender contains the numbers of transactions succeeded/failed published by a single sender.
type TxResultsBySender struct {
	NumSucceeded uint64
	NumFailed    uint64
}

type TxResults struct {
	BySender map[types.Account]TxResultsBySender
}
