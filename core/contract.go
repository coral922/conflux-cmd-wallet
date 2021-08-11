package core

import (
	conflux "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

type Contract struct {
	*conflux.Contract
}

func NewContract(c *conflux.Contract) *Contract {
	return &Contract{
		Contract : c,
	}
}

func (c *Contract) CreateUnsignedTransaction(from types.Address,
	amount *hexutil.Big,
	method string,
	args ...interface{}) (types.UnsignedTransaction, error) {
	data, err := c.GetData(method, args...)
	if err != nil {
		return types.UnsignedTransaction{}, errors.Wrap(err, "failed to encode call data")
	}
	tx, err := c.Client.CreateUnsignedTransaction(from, *c.Address, amount, data)
	if err != nil {
		return types.UnsignedTransaction{}, errors.Wrap(err, "failed to create tx")
	}
	return tx, err
}