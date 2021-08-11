package moonswap

import (
	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func GetPairCount() (int, error) {
	var num *big.Int
	err := Factory().Call(nil, &num, "allPairsLength")
	return int(num.Int64()), err
}

func GetPairByIndex(index int) (*Pair, error) {
	var pa common.Address
	err := Factory().Call(nil, &pa, "allPairs", big.NewInt(int64(index)))
	if err != nil {
		return nil, err
	}
	commonAddr := cfxaddress.MustNewFromCommon(pa, cfxaddress.NetowrkTypeMainnetID)
	contract, err := cfx().GetContract([]byte(MoonPairABI), &commonAddr)
	if err != nil {
		return nil, err
	}
	return (*Pair)(contract), nil
}
