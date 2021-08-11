package moonswap

import (
	conflux "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
)

func NewToken(address types.Address) (*Token, error) {
	c, err := cfx().GetContract([]byte(MoonTokenABI), &address)
	if err != nil {
		return nil, err
	}
	return (*Token)(c), nil
}

func (p *Token) Name() (string, error) {
	var n string
	err := (*conflux.Contract)(p).Call(nil, &n, "name")
	return n, err
}

func (p *Token) Symbol() (string, error) {
	var n string
	err := (*conflux.Contract)(p).Call(nil, &n, "symbol")
	return n, err
}

func (p *Token) Decimals() (uint8, error) {
	var n uint8
	err := (*conflux.Contract)(p).Call(nil, &n, "decimals")
	return n, err
}
