# redisx 的使用
使用前，请先调用一下 Init() 方法，进行初始化

Init 的入参是类型 `RedisConf`，它是一个接口类型，也就是参数必须实现 `RedisConf` 接口，主要作用是初始化 redis client 时，获取
参数进行初始化。`RedisConf` 的定义如下：

```
type RedisConf interface {
	GetAddr() string
	GetPwd() string
	GetDb() int
}
```
