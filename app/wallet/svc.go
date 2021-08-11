package wallet

import (
	"cfxWorld/core"
	"go.uber.org/dig"
)

type Service struct {
	c  *core.Client
	am *core.AccountMgr
	tm *core.TokenMgr
	nm *core.NftMgr
}

type Dep struct {
	dig.In
	C  *core.Client `name:"rw"`
	AM *core.AccountMgr
	TM *core.TokenMgr
	NM *core.NftMgr
}

func NewService(dep Dep) *Service {
	return &Service{
		c:  dep.C,
		am: dep.AM,
		tm: dep.TM,
		nm: dep.NM,
	}
}
