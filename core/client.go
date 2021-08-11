package core

import (
	"cfxWorld/config"
	"context"
	"fmt"
	conflux "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type Client struct {
	*conflux.Client
}

func NewClient() (*Client, error) {
	c, err := conflux.NewClient(config.C.NodeURL, conflux.ClientOption{
		KeystorePath:   config.C.Wallet.KeyStorePath,
		RetryCount:     config.C.Client.RequestMaxRetry,
		RetryInterval:  config.C.Client.RequestRetryInterval,
		RequestTimeout: config.C.Client.RequestTimeout,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		Client: c,
	}, nil
}

func (c *Client) WaitForPacked(ctx context.Context, hash types.Hash) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var errWaitPacked error
	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Second)
			tx, err := c.GetTransactionByHash(hash)
			if err != nil {
				errWaitPacked = errors.Wrapf(err, "failed to get transaction by hash %v", hash)
			}
			if tx.Status != nil {
				return
			}
			select {
			case <-ctx.Done():
				errWaitPacked = fmt.Errorf("wait tx %s packed time out", hash)
				return
			default:
			}
		}
	}()
	wg.Wait()
	return errWaitPacked
}

func (c *Client) WaitForReceipt(ctx context.Context, hash types.Hash) (*types.TransactionReceipt, error) {
	var errWaitRec error
	var txReceipt *types.TransactionReceipt
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Second)
			txReceipt, errWaitRec = c.GetTransactionReceipt(hash)
			if errWaitRec != nil {
				errWaitRec = errors.Wrapf(errWaitRec, "failed to get transaction receipt by hash %v", hash)
			}
			if txReceipt != nil {
				return
			}
			select {
			case <-ctx.Done():
				errWaitRec = fmt.Errorf("get tx %s receipt time out", hash)
				return
			default:
			}
		}
	}()
	wg.Wait()
	return txReceipt, errWaitRec
}

func (c *Client) Wait(ctx context.Context, hash types.Hash) (*types.TransactionReceipt, error) {
	err := c.WaitForPacked(ctx, hash)
	if err != nil {
		return nil, err
	}
	return c.WaitForReceipt(ctx, hash)
}
