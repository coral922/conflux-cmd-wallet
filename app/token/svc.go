package token

import (
	"cfxWorld/core"
	"go.uber.org/dig"
)

type Service struct {
	tm *core.TokenMgr
}

type Dep struct {
	dig.In
	TM *core.TokenMgr
}

func NewService(dep Dep) *Service {
	return &Service{
		tm: dep.TM,
	}
}
