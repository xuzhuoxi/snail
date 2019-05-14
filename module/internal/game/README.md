# Module game

## Extension

### How to write a extension

If your are using the snail built-in game module. then:

1. New a struct implement IGameExtension.
2. Implement one of the interfaces: IOnNoneRequestExtension、IOnBinaryRequestExtension or IOnObjectRequestExtension, and finish your game logic.
   And finish your game logic at the function under the interface.
3. Implement other interfaces if your need: IGoroutineExtension、IBatchExtension、IBeforeRequestExtension、IAfterRequestExtension and so on.
4. [Here](./extension/demo) is a demo.

### How to add extension to the game module

Call the function "RegisterExtension" in the package "github.com/xuzhuoxi/snail/module/internal/game/ifc" to register your extension go game module. Maybe like that:
```go
ifc.RegisterExtension(func() ifc.IGameExtension {
	return NewNoneDemoExtension("NoneDemo")
})
```

### Note

- Extension must implement two interfaces: IGameExtension and one of them(IOnNoneRequestExtension、IOnBinaryRequestExtension、IOnObjectRequestExtension).
- Other extension function, you can see:

  [/engine/extension](/engine/extension)<br>
  [/module/internal/game/ifc/iextension.go](/module/internal/game/ifc/iextension.go)<br>
  [github.com/xuzhuoxi/infra-go/infra-go/extendx](https://github.com/xuzhuoxi/infra-go/tree/master/extendx)<br>
