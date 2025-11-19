package kafka

type ProjectSvcConf1 struct {
	Kafka KafkaConf1
}

type KafkaConf1 struct {
	NameServers string
	TopicMap    map[string]string
	TopicName1  string
	GroupName1  string
}

func (k *KafkaConf1) GetNameServers() string {
	return k.NameServers
}

func (k *KafkaConf1) GetTopicByKey(key string) string {
	return k.TopicName1
}

func (k *KafkaConf1) GetGroupByKey(key string) string {
	return k.GroupName1
}

func (k *KafkaConf1) GetReconsumeTimes() int {
	return 1
}

func (k *KafkaConf1) GetRePushTimes() int {
	return 1
}
