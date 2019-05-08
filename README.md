# snail

# Usage

# How to get
```
go get -u github.com/xuzhuoxi/snail
```

# Related Library

- infra-go [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)<br>
基础库，整个snail框架中的大部分简单复用的逻辑都抽象到这个基础库中。

- goxc [https://github.com/laher/goxc](https://github.com/laher/goxc)<br>
打包依赖库，主要用于交叉编译

- json-iterator [https://github.com/json-iterator/go](https://github.com/json-iterator/go)<br>
带对应结构体的Json解释库

## Contact
xuzhuoxi<br>
<xuzhuoxi@gmail.com> or <mailxuzhuoxi@163.com>

## License
IconGen source code is available under the MIT [License](/LICENSE).

# Package Description
包分类及文件功能说明
<details>
<summary>Expand view</summary>
<pre><code>.
├── conf: 配置解释
│   ├── conf.go: 与配置相关结构体定义、解释，读取行为
├── engine: 引擎库
│   ├── extension: 扩展支持
│   │   ├── container.go: 扩展容器接口及基础结构体定义
│   │   ├── extension.go: 扩展接口及基础结构体定义
│   ├── mmo: MMO世界支持
│   │   ├── basis: 接口声明及公共结构体
│   │   │   ├── channel.go: Channel接口及常量定义，以及相关处理函数
│   │   │   ├── child.go: MMO子实体接口及常量定义，以及相关处理函数
│   │   │   ├── container.go: MMO实体容器接口及常量定义，以及相关处理函数
│   │   │   ├── entity.go:　MMO实体接口及常量定义，以及相关处理函数
│   │   │   ├── events.go: MMO事件接口及常量定义，以及相关处理函数
│   │   │   ├── group.go: MMO实体分组接口及常量定义，以及相关处理函数
│   │   │   ├── index.go: MMO实体索引接口及常量定义，以及相关处理函数
│   │   │   ├── manager.go: MMO管理器接口及常量定义，以及相关处理函数
│   │   │   ├── position.go: MMO坐标定义及行为
│   │   │   ├── proto.go: MMO协议号分组及定义
│   │   │   ├── team.go: MMO队伍及团队接口及常量定义，以及相关处理函数
│   │   │   ├── user.go: MMO玩家接口及常量定义，以及相关处理函数
│   │   │   ├── variable.go: MMO实体变量接口及常量定义，以及相关处理函数
│   │   ├── entity:basis包中实体接口对应的实现
│   │   │   ├── channel.go: Channel实现
│   │   │   ├── child.go: MMO子实体支持，并发安全
│   │   │   ├── container.go: MMO实体容器支持，并发安全
│   │   │   ├── group.go: MMO实体分组支持，并发安全
│   │   │   ├── room.go: MMO房间支持，并发安全
│   │   │   ├── team.go: MMO队伍支持，并发安全
│   │   │   ├── teamcorps.go: MMO团队支持，并发安全
│   │   │   ├── user.go: MMO玩家支持，并发安全
│   │   │   ├── userblackwhite.go: MMO玩家黑白名单支持，并发安全
│   │   │   ├── variable.go: MMO实体变量(包括用户变量、环境变量)支持，并发安全
│   │   │   ├── world.go: MMO世界支持，并发安全
│   │   │   ├── zone.go: MMO分区支持，并发安全
│   │   ├── index:basis包中index文件接口对应的实现
│   │   │   ├── channelidx.go: Channel索引管理，依赖entityidx中的EntityIndex
│   │   │   ├── entityidx.go:　实体索引管理EntityIndex，并发安全
│   │   │   ├── roomidx.go:　房间索引管理，依赖entityidx中的EntityIndex
│   │   │   ├── teamcorpsidx.go: 团队索引管理，依赖entityidx中的EntityIndex
│   │   │   ├── teamidx.go: 队伍索引管理，依赖entityidx中的EntityIndex
│   │   │   ├── useridx.go: 玩家索引管理，依赖entityidx中的EntityIndex
│   │   │   ├── zoneidx.go: 分区索引管理，依赖entityidx中的EntityIndex
│   │   ├── manager:
│   │   │   ├── broadcast.go: 消息广播管理
│   │   │   ├── entity.go: 实体管理，包括实体的创建，查找功能以及MMO世界的创建
│   │   │   ├── user.go: 玩家的在环境实体间转移管理，包括进入世界、分区，房间以及在房间间转移等操作
│   │   │   ├── variable.go: 变量监听管理，监听变量(环境变量、用户变量)更新，进行更新消息广播
│   │   ├── proto:
│   │   │   ├── define.go: 基础通信协议定义
│   │   ├── mmo.go: MMO管理入口
├── module:
│   ├── imodule:
│   │   ├── module.go: 模块基础接口与实现。模块注册，模块实现化等相关功能
│   │   ├── rpc.go: 模块RPC通信实现，支持自定义扩展
│   │   ├── state.go: Socket Server的状态支持，包括响应时间，连接数等信息记录
│   ├── internal:
│   │   ├── admin: 游戏管理模块，目前为空
│   │   ├── game: 游戏逻辑模块，已经实现包括向route模块进行集群登记、状态更新、集群注销，业务扩展管理等功能。其它具体游戏逻辑业务可通过扩展进行增加
│   │   ├── route: 游戏路由模块，包括登录分配，服务器集群登记等功能
│   │   ├── cmds.go: 服务器命令行管理，目前实现包括模块启动、模块关闭、模块信息状态查询功能
│   │   ├── internal.go: 模块内部管理入口
│   ├── module.go: 模块对外管理入口
</code></pre>
</details>