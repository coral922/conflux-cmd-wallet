package tx

import (
	"cfxWorld/core"
	"go.uber.org/dig"
)

const DefaultSlippageTolerance = 1

type Service struct {
	c   *core.Client
	txm *core.TxMgr
	am  *core.AccountMgr
}

type Dep struct {
	dig.In
	C   *core.Client `name:"rw"`
	TXM *core.TxMgr
	AM  *core.AccountMgr
}

func NewService(dep Dep) *Service {
	return &Service{
		c:   dep.C,
		txm: dep.TXM,
		am:  dep.AM,
	}
}
