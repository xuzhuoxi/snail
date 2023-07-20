// Package rabbit
// Created by xuzhuoxi
// on 2019-02-19.
// @author xuzhuoxi
//
package root

import (
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

func NewGameSingleCase() ifc.IGameSingleCase {
	return &GameSingleCase{}
}

type GameSingleCase struct {
	isInit bool
	logger logx.ILogger
}

func (s *GameSingleCase) Init() {
	if s.isInit {
		return
	}
	s.isInit = true
}
func (s *GameSingleCase) GetLogger() logx.ILogger {
	return s.logger
}

func (s *GameSingleCase) SetLogger(logger logx.ILogger) {
	s.logger = logger
}
