# Module game

## Extension

### How to write a extension

If your are using the snail built-in game module. then:

1. Define a extension(struct), embedding GameExtensionSupport and implement IGameExtension.
2. You can call the function SetRequestHandler, SetRequestHandlerBinary, SetRequestHandlerJson or SetRequestHandlerObject to associate a protocol response logic.
3. Finish the protocol response logic said above.
4. Implement other interfaces if your need: IGoroutineExtension、IBeforeRequestExtension、IAfterRequestExtension and so on.
5. [Here](./extension/demo) is a demo.

### How to add extension to the game module

Call the function "RegisterExtension" in the package "github.com/xuzhuoxi/snail/module/internal/game/ifc" to register your extension go game module. Maybe like that:
```go
ifc.RegisterExtension(func() ifc.IGameExtension {
	return NewDemoExtension("Demo")
})
```

### Note

- Extension must implement two interfaces: IGameExtension.
- Other extension function, you can see:

  [/module/internal/game/ifc/iextension.go](/module/internal/game/ifc/iextension.go)<br>
  [github.com/xuzhuoxi/infra-go/infra-go/extendx](https://github.com/xuzhuoxi/infra-go/tree/master/extendx)<br>
