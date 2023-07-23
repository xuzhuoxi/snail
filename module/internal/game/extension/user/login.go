//
//Created by xuzhuoxi
//on 2019-03-03.
//@author xuzhuoxi
//
package user

import (
	"github.com/xuzhuoxi/infra-go/extendx/protox"
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

func (e *LoginExtension) InitExtension() error {
	e.GetLogger().Debugln("LoginExtension.InitExtension", e.Name)
	e.SetRequestHandlerString(LoginId, e.onRequestLogin)
	e.SetRequestHandlerString(ReLoginId, e.onRequestReLogin)
	return nil
}

func (e *LoginExtension) DestroyExtension() error {
	e.ClearRequestHandler(ReLoginId)
	e.ClearRequestHandler(LoginId)
	e.GetLogger().Debugln("LoginExtension.DestroyExtension", e.Name)
	return nil
}

func (e *LoginExtension) onRequestLogin(resp protox.IExtensionStringResponse, req protox.IExtensionStringRequest) {
	password := req.StringData()[0]
	if e.check(req.ClientId(), password) {
		ifc.AddressProxy.MapIdAddress(req.ClientId(), req.ClientAddress())
		time.Sleep(time.Millisecond * 20)
		resp.SendStringResponse("ok")
		e.GetLogger().Traceln("LoginExtension.onRequestLogin:", "Check Succ!", req.ProtoId(), req.ClientId(), password)
	} else {
		e.GetLogger().Warnln("LoginExtension.onRequestLogin:", "Check Fail!", req.ProtoId(), req.ClientId(), password)
	}
}

func (e *LoginExtension) onRequestReLogin(resp protox.IExtensionStringResponse, req protox.IExtensionStringRequest) {
	password := req.StringData()[0]
	if e.check(req.ClientId(), password) {
		ifc.AddressProxy.MapIdAddress(req.ClientId(), req.ClientAddress())
		time.Sleep(time.Millisecond * 20)
		resp.SendStringResponse("ok")
		e.GetLogger().Traceln("LoginExtension.onRequestReLogin:", "Check Succ!", req.ProtoId(), req.ClientId(), password)
	} else {
		e.GetLogger().Warnln("LoginExtension.onRequestReLogin:", "Check Fail!", req.ProtoId(), req.ClientId(), password)
	}
}

func (e *LoginExtension) check(uid string, password string) bool {
	return uid == password
}
