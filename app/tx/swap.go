package tx

import (
	"cfxWorld/config"
	"cfxWorld/core"
	"cfxWorld/lib/moonswap"
	"cfxWorld/lib/util"
	"context"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	moon "github.com/coral922/moonswap-sdk-go/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

type SwapReceipt struct {
	*types.TransactionReceipt
	SendToken  string
	SendAmount *hexutil.Big
	GetToken   string
	GetAmount  *hexutil.Big
}

func (s *Service) MakeSwap(mode int, input, output, amount string) (*moon.Trade, error) {
	a, err := util.FloatStrToAmount(amount)
	if err != nil {
		return nil, err
	}
	return s.txm.MakeSwap(mode, input, output, (*big.Int)(a))
}

func (s *Service) DoSwap(addr string, trade *moon.Trade, password string, slipTolerance int) (*SwapReceipt, error) {
	fromAddr, err := s.am.ParseStringToAddress(addr, true)
	if err != nil {
		return nil, err
	}
	err = s.am.Unlock(*fromAddr, core.ToRealPW(password))
	if err != nil {
		return nil, err
	}
	if slipTolerance <= 0 {
		slipTolerance = DefaultSlippageTolerance
	}

	tx, err := s.txm.SwapTx(*fromAddr, trade, slipTolerance)
	if err != nil {
		return nil, err
	}
	hash, err := s.txm.Do(tx)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.C.Tx.MaxWait)
	defer cancel()
	receipt, err := s.c.Wait(ctx, hash)
	if err != nil {
		return nil, err
	}
	var sendA, getA *hexutil.Big
	result := moonswap.ParseSwapResult(trade, receipt)
	if result != nil {
		sendA = (*hexutil.Big)(result.Input.Raw())
		getA = (*hexutil.Big)(result.Output.Raw())
	}
	return &SwapReceipt{
		TransactionReceipt: receipt,
		SendToken:          trade.InputAmount().Symbol,
		SendAmount:         sendA,
		GetToken:           trade.OutputAmount().Symbol,
		GetAmount:          getA,
	}, nil
}
