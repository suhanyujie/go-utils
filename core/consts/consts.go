package consts

const (
	PlatformHttpApiOk = 200
	GrpcSvcCodeOk     = 200

	TimeDateFormat = "2006-01-02"

	MetaKeyForGrpcPort       = "grpc.server.port"
	MetaKeyPriorityForDevEnv = "priority"
)

var (
	NacosApiServerPort = 8848

	// 游戏服务注册到 nacos 时，需要向实例中写入 meta 数据，这些 meta 数据的 key 在这里定义
	NacosMetaKeyOfGameId = "gameId"
	NacosMetaKeyOfWsPort = "wsPort"
)

const (
	GrpcErrMsg = "grpc server error"
)
