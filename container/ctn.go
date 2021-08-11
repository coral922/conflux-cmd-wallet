package container

import (
	"cfxWorld/app/nft"
	"cfxWorld/app/token"
	"cfxWorld/app/tx"
	"cfxWorld/app/wallet"
	"cfxWorld/core"
	"cfxWorld/lib/storage"
	"go.uber.org/dig"
	"log"
	"os"
)

var Ctn = dig.New()

type provider struct {
	Fn interface{}
	As string
}

var providers = []provider{
	{core.NewClient, "rw"},
	{core.NewAccountMgr, ""},
	{core.NewTokenMgr, ""},
	{core.NewNftMgr, ""},
	{core.NewTxMgr, ""},
	{func() (*storage.Storage, error) {
		err := os.MkdirAll(core.DBFolder, 0666)
		if err != nil {
			return nil, err
		}
		return storage.NewStorage(core.DBBaseFile)
	}, ""},
	{wallet.NewService, ""},
	{tx.NewService, ""},
	{token.NewService, ""},
	{nft.NewService, ""},
}

func init() {
	for _, provider := range providers {
		if err := Ctn.Provide(provider.Fn, dig.Name(provider.As)); err != nil {
			log.Fatalf("Init failed : %v", err)
		}
	}
}
