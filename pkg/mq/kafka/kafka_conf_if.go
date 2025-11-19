package kafka

// 在项目中，针对kafka配置字段，实现 KafkaConfIf 接口即可
type KafkaConfIf interface {
	GetNameServers() string
	GetTopicByKey(key string) string
	GetGroupByKey(key string) string
	GetReconsumeTimes() int
	GetRePushTimes() int
}

type KafkaConf struct {
	Pubs []*Options `xml:"pub"`
}
