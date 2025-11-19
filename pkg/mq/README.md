## kafka 库使用方式
### 实现 kafka 配置接口
接口定义见 `pkg/mq/kafka/kafka_conf_if.go`，实现对应方法，就能让库内部获取配置值。可以参考[农场示例](http://192.168.3.3:9980/game-server/farm_go/-/blob/a8759bad6e3ef08b8dff4ffbe5868cf8d69f8f71/global/conf/kafka_conf.go)。

### 在服务启动时，调用初始化
参考[农场使用](http://192.168.3.3:9980/game-server/farm_go/-/blob/a8759bad6e3ef08b8dff4ffbe5868cf8d69f8f71/internal/api_server/api_route/server.go#L91) 

在初始化前，确保调用了 go_utils 包的[日志初始化](http://192.168.3.3:9980/game-server/farm_go/-/blob/a8759bad6e3ef08b8dff4ffbe5868cf8d69f8f71/internal/api_server/api_route/server.go#L72)。

### 项目内适用
参考[农场使用](http://192.168.3.3:9980/game-server/farm_go/-/blob/a8759bad6e3ef08b8dff4ffbe5868cf8d69f8f71/internal/domain/farm_domain/farm_domain.go#L1065-1080)。
