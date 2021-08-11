package core

import (
	"cfxWorld/lib/crawler"
	"fmt"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
)

type TxOption func(tx *types.UnsignedTransactionBase)

func applyTxOpts(tx *types.UnsignedTransactionBase, opts ...TxOption) {
	for _, opt := range opts {
		opt(tx)
	}
}

func TxURL(txHash types.Hash) string {
	return fmt.Sprintf("%s/transaction/%s", crawler.CfxScanBaseURL, txHash.String())
}


