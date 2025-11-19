package redisx

type RedisConf interface {
	GetAddr() string
	GetPwd() string
	GetDb() int
}
