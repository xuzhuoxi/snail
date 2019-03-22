//
//Created by xuzhuoxi
//on 2019-03-03.
//@author xuzhuoxi
//
package user

import (
	"github.com/xuzhuoxi/infra-go/extendx"
	"github.com/xuzhuoxi/snail/module/internal/game/extension"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

const (
	LoginId   = "LI"
	ReLoginId = "RLI"
)

func NewLoginExtension(Name string, SingleCase ifc.IGameSingleCase) *LoginExtension {
	return &LoginExtension{GameExtensionSupport: extension.NewGameExtensionSupport(Name, SingleCase)}
}

type LoginExtension struct {
	extension.GameExtensionSupport
}

func (e *LoginExtension) InitProtocolId() {
	e.ProtoIdToValue[LoginId] = struct{}{}
	e.ProtoIdToValue[ReLoginId] = struct{}{}
}

func (e *LoginExtension) OnRequest(resp extendx.IExtensionResponse, protoId string, uid string, data []byte, data2 ...[]byte) {
	password := string(data)
	if e.check(uid, password) {
		e.SingleCase.AddressProxy().MapIdAddress(uid, resp.SenderAddress())
		switch protoId {
		case LoginId:
			break
		case ReLoginId:
			break
		}
		e.GetLogger().Debugln("LoginExtension.OnRequest:", "Check Succ!", protoId, uid, password)
	} else {
		e.GetLogger().Debugln("LoginExtension.OnRequest:", "Check Fail!", protoId, uid, password)
	}
}

func (e *LoginExtension) check(uid string, password string) bool {
	return uid == password
}
