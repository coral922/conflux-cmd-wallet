package moonswap

import (
	conflux "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func NewPair(address types.Address) (*Pair, error) {
	c, err := cfx().GetContract([]byte(MoonPairABI), &address)
	if err != nil {
		return nil, err
	}
	return (*Pair)(c), nil
}

func (p *Pair) Token0() (*Token, error) {
	return p.token(true)
}

func (p *Pair) Token1() (*Token, error) {
	return p.token(false)
}

func (p *Pair) Tokens() ([2]*Token, error) {
	t0, err := p.Token0()
	if err != nil {
		return [2]*Token{}, err
	}
	t1, err := p.Token1()
	if err != nil {
		return [2]*Token{}, nil
	}
	return [2]*Token{t0, t1}, nil
}

func (p *Pair) Reserves() ([2]*big.Int, error) {
	type reserves struct {
		R0 *big.Int `abi:"_reserve0"`
		R1 *big.Int `abi:"_reserve1"`
		T  uint32   `abi:"_blockTimestampLast"`
	}
	var R reserves
	err := (*conflux.Contract)(p).Call(nil, &R, "getReserves")
	if err != nil {
		return [2]*big.Int{}, err
	}
	return [2]*big.Int{R.R0, R.R1}, nil
}

func (p *Pair) token(first bool) (*Token, error) {
	method := "token0"
	if !first {
		method = "token1"
	}
	var a common.Address
	err := (*conflux.Contract)(p).Call(nil, &a, method)
	if err != nil {
		return nil, err
	}
	commonAddr := cfxaddress.MustNewFromCommon(a, cfxaddress.NetowrkTypeMainnetID)
	contract, err := cfx().GetContract([]byte(MoonTokenABI), &commonAddr)
	if err != nil {
		return nil, err
	}
	return (*Token)(contract), nil
}
