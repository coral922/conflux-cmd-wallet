package nft

import (
	"cfxWorld/core"
	"go.uber.org/dig"
)

type Service struct {
	nm *core.NftMgr
}

type Dep struct {
	dig.In
	NM *core.NftMgr
}

func NewService(dep Dep) *Service {
	return &Service{
		nm: dep.NM,
	}
}
