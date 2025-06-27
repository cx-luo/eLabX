```shell
api开发
├── conf                    #项目配置文件目录
│   └── config.toml         #大家可以选择自己熟悉的配置文件管理工具包例如：toml、xml等等
├── controllers             #控制器目录，按模块存放控制器（或者叫控制器函数），必要的时候可以继续划分子目录。
│   └── user.go
├── models                  #模型目录，负责项目的数据存储部分，例如各个模块的Mysql表的读写模型。
│   ├── food.go
│   ├── user.go
│	└── init.go				#模型初始化
├── logs                    #日志文件目录，主要保存项目运行过程中产生的日志。
├── main.go                 #项目入口，这里负责Gin框架的初始化，注册路由信息，关联控制器函数等。
```
以上两种目录结构基本上大同小异，一般的项目构建都可以使用这样的项目结构，清晰明了。

3.项目结构扩展
对于一般的项目都是使用上面的目录结构，但是有时候项目复杂了，会进一步拆分，把数据校验，模型定义，数据操作，服务，响应进行解耦。
```shell
数据校验->控制器->调用服务->数据操作->数据模型->数据响应
```

```shell
├── conf                    #项目配置文件目录
│   └── config.toml         #大家可以选择自己熟悉的配置文件管理工具包例如：toml、xml、ini等等
├── requests                #定义入参即入参校验规则
│   └── user_request.go
│   └── food_request.go
├── responses                #定义响应的数据
│   └── user_response.go
│   └── food_response.go
├── services                #服务定义目录
|	└── v1					#服务v1版本
│   |	└── user_service.go
│   |	└── food_service.go
|	└── indigov2					#服务v2版本
│   |	└── user_service.go
│   |	└── food_service.go
├── api             		#api目录，按模块存放控制器（或者叫控制器函数），必要的时候可以继续划分子目录。
│   └── v1					#apiv1版本
│   |	└── user.go
│   |	└── food.go
│   └── indigov2					#apiv2版本
│   |	└── user.go
│   |	└── food.go
├── router					#路由目录
│   └── v1					#路由v1版本
│   |	└── user.go
│   |	└── food.go
│   └── indigov2					#路由v2版本
│   |	└── user.go
│   |	└── food.go
├── init.go					#路由初始化
├── pkg						#自定义的工具类等
│   └── e					#项目统一的响应定义，如错误码，通用的错误信息，响应的结构体
│   └── util				#工具类目录
├── models                  #模型目录，负责项目的数据存储部分，例如各个模块的Mysql表的读写模型。
│   ├── food.go
│   ├── user.go
│	└── init.go				#模型初始化
├── repositories            #数据操作层，定义各种数据操作。
│   └── user_repository.go
│   └── food_repository.go
├── logs                    #日志文件目录，主要保存项目运行过程中产生的日志。
├── main.go                 #项目入口，这里负责Gin框架的初始化，注册路由信息，关联控制器函数等。

```