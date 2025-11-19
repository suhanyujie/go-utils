package kafka

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/suhanyujie/go_utils/jsonx"
	"github.com/suhanyujie/go_utils/logx"
	"github.com/suhanyujie/go_utils/pkg/mq/mq_model"
)

type exampleConsumerGroupHandler struct {
	fu          func(message *mq_model.MqMessageExt) error
	errCallback func(message *mq_model.MqMessageExt)
	proxy       *Proxy
}

func (exampleConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	return nil
}
func (exampleConsumerGroupHandler) Cleanup(c sarama.ConsumerGroupSession) error {
	return nil
}
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg != nil {
			//获取重试次数
			ReconsumeTimes := 0
			if msg.Headers != nil {
				for _, header := range msg.Headers {
					if string(header.Key) == RecordHeaderReconsumeTimes {
						v, _ := strconv.ParseInt(string(header.Value), 10, 32)
						ReconsumeTimes = int(v)
					}
				}
			}
			afterConsumerTimes := ReconsumeTimes - 1

			msgExt := &mq_model.MqMessageExt{
				MqMessage: mq_model.MqMessage{
					Topic:          msg.Topic,
					Body:           string(msg.Value),
					Keys:           string(msg.Key),
					Partition:      msg.Partition,
					Offset:         msg.Offset,
					ReconsumeTimes: &afterConsumerTimes,
				},
			}
			err1 := h.fu(msgExt)
			if err1 != nil {
				logx.GetLogger().Errorf("Kafka 业务消费异常 %v", err1)
				logx.GetLogger().Errorf("Topic: %s, 处理失败的消息：%s ", msg.Topic, jsonx.ToJsonIgnoreErr(msgExt.MqMessage))
				sess.MarkMessage(msg, "consumer err")

				logx.GetLogger().Infof("剩余消费次数%d", afterConsumerTimes)
				if afterConsumerTimes > -1 {
					_, pushErr := h.proxy.PushMessage(&msgExt.MqMessage)
					if pushErr != nil {
						logx.GetLogger().Errorf("重试推送失败, 消息内容:%s", jsonx.ToJsonIgnoreErr(msgExt.MqMessage))
					}
				} else {
					logx.GetLogger().Errorf("无重试次数, 消息最终消费失败, 消息内容%s", jsonx.ToJsonIgnoreErr(msgExt.MqMessage))
					h.errCallback(msgExt)
				}
			} else {
				sess.MarkMessage(msg, "")
			}
		}
	}
	return nil
}

func (proxy *Proxy) ConsumeMessage(topic string, groupId string, fu func(message *mq_model.MqMessageExt) error, errCallback func(message *mq_model.MqMessageExt)) error {
	kafkaConfig := getKafkaConfig()
	logx.GetLogger().Infof("Kafka config %s", jsonx.ToJsonIgnoreErr(kafkaConfig))
	logx.GetLogger().Infof("Starting a new Sarama consumer, topic %s, groupId %s", topic, groupId)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	//config.ChannelBufferSize = 2560
	//config.Consumer.Fetch.Min = 1024 * 1024
	//config.Consumer.Fetch.Default = 1024 * 1024 * 2
	//config.Consumer.Fetch.Max = 1024 * 1024 * 10
	//config.Consumer.MaxWaitTime = 2 * time.Second
	//config.Consumer.MaxProcessingTime = 1 * time.Second
	//config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Version = version

	topics := []string{topic}

	for {
		logx.GetLogger().Infof("kafka consumer 开始连接...")
		ctx, _ := context.WithCancel(context.Background())
		client, err := sarama.NewConsumerGroup(strings.Split(kafkaConfig.GetNameServers(), ","), groupId, config)
		if err != nil {
			logx.GetLogger().Errorf("Error creating consumer group client: %v", err)
			return err
		}

		handler := exampleConsumerGroupHandler{
			fu:          fu,
			proxy:       proxy,
			errCallback: errCallback,
		}

		for {
			// logx.GetLogger().Infof("准备消费, topic %s, groupId %s", topic, groupId)
			if err := client.Consume(ctx, topics, &handler); err != nil {
				logx.GetLogger().Errorf("[ConsumeMessage] Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				logx.GetLogger().Errorf("异常退出 %v", ctx.Err())
				break
			}
		}

		err = client.Close()
		if err != nil {
			logx.GetLogger().Errorf("[ConsumeMessage] 关闭连接失败 %v", err)
		}
		time.Sleep(2 * time.Second)
		logx.GetLogger().Info("[ConsumeMessage] 准备重连...")
	}

	logx.GetLogger().Info("[ConsumeMessage] kafka 消费结束")

	return nil
}
