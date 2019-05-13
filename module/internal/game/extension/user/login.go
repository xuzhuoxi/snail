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
	"time"
)

const (
	LoginId   = "LI"
	ReLoginId = "RLI"
)

func NewLoginExtension(Name string) *LoginExtension {
	return &LoginExtension{GameExtensionSupport: extension.NewGameExtensionSupport(Name)}
}

type LoginExtension struct {
	extension.GameExtensionSupport
}

func (e *LoginExtension) InitProtocolId() {
	e.ProtoIdToValue[LoginId] = struct{}{}
	e.ProtoIdToValue[ReLoginId] = struct{}{}
}

func (e *LoginExtension) InitExtension() error {
	e.GetLogger().Debugln("LoginExtension.InitExtension", e.Name)
	return nil
}

func (e *LoginExtension) OnRequest(resp extendx.IExtensionBinaryResponse, protoId string, uid string, data []byte, data2 ...[]byte) {
	password := string(data)
	if e.check(uid, password) {
		ifc.AddressProxy.MapIdAddress(uid, resp.SenderAddress())
		time.Sleep(time.Millisecond * 20)
		switch protoId {
		case LoginId:
			break
		case ReLoginId:
			break
		}
		data := append([]byte(protoId), []byte(uid)...)
		resp.SendBinaryResponse(data)
		e.GetLogger().Traceln("LoginExtension.OnRequest:", "Check Succ!", protoId, uid, password, data)
	} else {
		e.GetLogger().Warnln("LoginExtension.OnRequest:", "Check Fail!", protoId, uid, password)
	}
}

func (e *LoginExtension) check(uid string, password string) bool {
	return uid == password
}
