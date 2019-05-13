# Module game

# Game模块Extension规范

Extension至少实现两个接口

- IGameExtension(必须)
- IOnNoneRequestExtension、IOnBinaryRequestExtension、IOnObjectRequestExtension(选一)
- IGoroutineExtension、IBatchExtension、IBeforeRequestExtension、IAfterRequestExtension(可选)