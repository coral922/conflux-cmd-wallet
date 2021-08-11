package tx

import (
	"cfxWorld/config"
	"cfxWorld/core"
	"cfxWorld/lib/util"
	"context"
	"fmt"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

type TransferReceipt struct {
	*types.TransactionReceipt
	Token      string
	NftTokenID string
	Value      *hexutil.Big
}

func (s *Service) Send(token, from, to, amount, password string) (*TransferReceipt, error) {
	a, err := util.FloatStrToAmount(amount)
	if err != nil {
		return nil, err
	}
	fromAddr, err := s.am.ParseStringToAddress(from, true)
	if err != nil {
		return nil, err
	}

	toAddr, err := s.am.ParseStringToAddress(to, false)
	if err != nil {
		return nil, err
	}

	if !s.txm.TransferSupportedToken(token) {
		return nil, fmt.Errorf("unsupported token [%s]", token)
	}

	err = s.am.Unlock(*fromAddr, core.ToRealPW(password))
	if err != nil {
		return nil, err
	}

	tx, err := s.txm.SendTx(token, *fromAddr, *toAddr, a)
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
	return &TransferReceipt{
		TransactionReceipt: receipt,
		Token:              token,
		Value:              a,
	}, nil
}

func (s *Service) SendNft(token, tokenID, from, to, password string) (*TransferReceipt, error) {
	id, err := util.IntStrToBig(tokenID)
	if err != nil {
		return nil, err
	}
	fromAddr, err := s.am.ParseStringToAddress(from, true)
	if err != nil {
		return nil, err
	}

	toAddr, err := s.am.ParseStringToAddress(to, false)
	if err != nil {
		return nil, err
	}

	if !s.txm.TransferSupportedNft(token) {
		return nil, fmt.Errorf("unsupported nft token [%s]", token)
	}

	err = s.am.Unlock(*fromAddr, core.ToRealPW(password))
	if err != nil {
		return nil, err
	}

	tx, err := s.txm.SendNftTx(token, id, *fromAddr, *toAddr)
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
	return &TransferReceipt{
		TransactionReceipt: receipt,
		Token:              token,
		NftTokenID:         tokenID,
		Value:              (*hexutil.Big)(big.NewInt(1)),
	}, nil
}
