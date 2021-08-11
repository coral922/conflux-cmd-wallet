package crawler

import (
	"cfxWorld/lib/standard"
	conflux "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

type CoveredBySponsorInfo struct {
	Gas        bool
	Collateral bool
}

func SponsorOf(contractAddr types.Address) (CoveredBySponsorInfo, error) {
	var info CoveredBySponsorInfo
	b, err := cfx().GetSponsorInfo(contractAddr)
	if err != nil {
		return info, err
	}
	if b.SponsorBalanceForGas.ToInt().Cmp(big.NewInt(0)) == 1 {
		info.Gas = true
	}
	if b.SponsorBalanceForCollateral.ToInt().Cmp(big.NewInt(0)) == 1 {
		info.Collateral = true
	}
	return info, nil
}

func Contract(abi []byte, contractAddr types.Address) (*conflux.Contract, error) {
	return cfx().GetContract(abi, &contractAddr)
}

func CRC20TotalSupply(contractAddr types.Address) (*hexutil.Big, error) {
	t, err := cfx().GetContract([]byte(standard.CRC20BaseABI), &contractAddr)
	if err != nil {
		return nil, err
	}
	var b *big.Int
	err = t.Call(nil, &b, "totalSupply")
	return (*hexutil.Big)(b), err
}
