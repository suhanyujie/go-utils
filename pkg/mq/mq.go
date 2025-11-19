package mq

import (
	"github.com/suhanyujie/go_utils/pkg/mq/kafka"
	"github.com/suhanyujie/go_utils/pkg/mq/mq_model"
)

type MQClient interface {
	PushMessage(messages ...*mq_model.MqMessage) (*[]mq_model.MqMessageExt, error)
	ConsumeMessage(topic string, groupId string, fu func(message *mq_model.MqMessageExt) error, errCallback func(message *mq_model.MqMessageExt)) error
}

var (
	kafkaClient MQClient = &kafka.Proxy{}
)

func GetMQClient() MQClient {
	return kafkaClient
}
