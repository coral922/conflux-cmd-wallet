package crawler

import (
	conflux "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"log"
	"math/big"
	"sync"
)

var client *conflux.Client

var l = sync.Mutex{}

var cfxClientOption = struct {
	nodeUrl string
	option conflux.ClientOption
}{}

func LazyLoad(nodeURL string, option ...conflux.ClientOption) {
	cfxClientOption.nodeUrl = nodeURL
	if len(option) > 0 {
		cfxClientOption.option = option[0]
	}
}

func cfx() *conflux.Client {
	l.Lock()
	defer l.Unlock()
	if client != nil {
		return client
	} else {
		c, err := conflux.NewClient(cfxClientOption.nodeUrl, cfxClientOption.option)
		if err != nil {
			log.Fatal(err)
		}
		client = c
		return client
	}
}

func IsCfxBalanceEnough(tx types.UnsignedTransaction) (enough bool, least *hexutil.Big) {
	balance, err := cfx().GetBalance(*tx.From)
	if err != nil {
		log.Println(err)
		return enough, (*hexutil.Big)(big.NewInt(0))
	}
	s, err := SponsorOf(*tx.To)
	if err != nil {
		log.Println(err)
	}
	total := new(big.Int)
	total.Add(total, tx.Value.ToInt())
	if !s.Gas {
		gasTotal := new(big.Int)
		gasTotal.Mul(tx.Gas.ToInt(), tx.GasPrice.ToInt())
		total.Add(total, gasTotal)
	}
	if !s.Collateral {
		storage := new(big.Int).SetUint64(uint64(*tx.StorageLimit))
		total.Add(total, storage)
	}
	b := balance.ToInt()
	return b.Cmp(total) >= 0, (*hexutil.Big)(total)
}