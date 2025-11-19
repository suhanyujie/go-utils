package kafka

import (
	"strconv"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/suhanyujie/go_utils/jsonx"
	"github.com/suhanyujie/go_utils/logx"
	"github.com/suhanyujie/go_utils/pkg/mq/mq_model"
	"github.com/suhanyujie/go_utils/redisx"
	"go.uber.org/zap"
)

// 日志相关常量
const (
	LogAppKey       = "appName"
	LogTagKey       = "tag"
	LogMqMessageKey = "mqMessage"
)

type Proxy struct {
	// key: topic + partitioner
	producers map[string]sarama.AsyncProducer
}

var (
	version = sarama.V3_3_1_0

	producerConfig = sarama.NewConfig()
	kafkaConfObj   KafkaConfIf
)

const (
	RecordHeaderReconsumeTimes = "ReconsumeTimes"
	RecordHeaderRePushTimes    = "RePushTimes"
)

func init() {
	//生产者通用配置
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.Return.Errors = true
	producerConfig.Version = version
}

type MsgMetadata struct {
	//消费重试次数
	ReconsumeTimes int
	//推送重试次数
	RePushTimes int
}

func SetKafkaConfIf(confObj KafkaConfIf) {
	kafkaConfObj = confObj
}

func getKafkaConfig() KafkaConfIf {
	return kafkaConfObj
}

func (proxy *Proxy) getProducerAutoConnect(topic string) (*sarama.AsyncProducer, error) {
	//key := topic + "#" + strconv.Itoa(int(partition))
	producer, err := proxy.getProducer(topic)
	if err != nil {
		logx.GetLogger().Errorf("err: %v", err)
		return nil, err
	}
	return producer, nil
}

func (proxy *Proxy) getProducer(topic string) (*sarama.AsyncProducer, error) {
	key := topic
	if proxy.producers == nil {
		proxy.producers = map[string]sarama.AsyncProducer{}
	}

	if v, ok := proxy.producers[key]; ok && v != nil {
		return &v, nil
	}

	suc, err := redisx.SetNxExpire(key, 1, 5*time.Second)
	if err != nil {
		logx.GetLogger().Errorf("err:%v", err)
		return nil, err
	}
	if suc {
		//如果获取到锁，则开始初始化
		defer func() {
			redisx.UnLock(key)
			if _, err := redisx.UnLock(key); err != nil {
				logx.GetLogger().Errorf("err:%v", err)
			}
		}()
	}

	//二次确认
	if v, ok := proxy.producers[key]; ok && v != nil {
		return &v, nil
	}

	//重新构造producer
	producer, err1 := proxy.buildProducer()
	if err1 != nil {
		logx.GetLogger().Errorf("err:%v", err1)
		return nil, err1
	}

	proxy.producers[key] = *producer
	return producer, nil
}

func (proxy *Proxy) CloseConnect(topic string) error {
	proxy.producers[topic] = nil
	return nil
}

func (proxy *Proxy) buildProducer() (*sarama.AsyncProducer, error) {
	kafkaConfig := getKafkaConfig()
	logx.GetLogger().Infof("build producer")

	producer, err := sarama.NewAsyncProducer(strings.Split(kafkaConfig.GetNameServers(), ","), producerConfig)
	if err != nil {
		logx.GetLogger().Infof("producer_test create producer error :%#v", err)
		return nil, err
	}
	return &producer, nil
}

func (proxy *Proxy) PushMessage(messages ...*mq_model.MqMessage) (*[]mq_model.MqMessageExt, error) {
	if messages == nil || len(messages) == 0 {
		return nil, errors.New("[PushMessage] message is empty")
	}

	msgExts := make([]mq_model.MqMessageExt, len(messages))
	for i, message := range messages {
		// 传递 metadata，方便消费端重试
		ReconsumeTimes := getKafkaConfig().GetReconsumeTimes()
		RePushTimes := getKafkaConfig().GetRePushTimes()
		if message.ReconsumeTimes != nil {
			ReconsumeTimes = *message.ReconsumeTimes
		}
		if message.RePushTimes != nil {
			RePushTimes = *message.RePushTimes
		}

		// send message
		msg := &sarama.ProducerMessage{
			Topic: message.Topic,
			//Partition: message.Partition,
			Key:   sarama.StringEncoder(message.Keys),
			Value: sarama.ByteEncoder(message.Body),
			Headers: []sarama.RecordHeader{
				{
					Key:   []byte(RecordHeaderReconsumeTimes),
					Value: []byte(strconv.Itoa(ReconsumeTimes)),
				},
			},
		}

		var pushErr error = nil
		for rePushTime := 0; rePushTime <= RePushTimes; rePushTime++ {
			p, err1 := proxy.getProducerAutoConnect(message.Topic)
			if err1 != nil {
				logx.GetLogger().Errorf("err:%v", err1)
				return nil, err1
			}
			producer := *p

			if rePushTime > 0 {
				logx.GetLogger().Infof("[PushMessage] 重试次数%d，最大次数%d, 上次失败原因%v, 消息内容%s", rePushTime, message.RePushTimes, pushErr, jsonx.ToJsonIgnoreErr(message))
			}
			producer.Input() <- msg
			select {
			case suc := <-producer.Successes():
				logx.GetLogger().Infof("[PushMessage] 推送成功, offset: %d,  timestamp: %s， 消息内容%s", suc.Offset, suc.Timestamp.String(), jsonx.ToJsonIgnoreErr(message))
				pushErr = nil
			case fail := <-producer.Errors():
				logx.GetLogger().Errorf("[PushMessage] err: %v", fail.Err.Error())
				pushErr = fail.Err
				//return nil, errors.BuildSystemErrorInfo(errors.KafkaMqSendMsgError, fail)

				if pushErr == sarama.ErrNotConnected || pushErr == sarama.ErrClosedClient || pushErr == sarama.ErrOutOfBrokers {
					logx.GetLogger().Errorf("断开连接... %v", pushErr)
					//重连
					closeErr := proxy.CloseConnect(message.Topic)
					if closeErr != nil {
						logx.GetLogger().Errorf("err:%v", closeErr)
					}
				}
				time.Sleep(time.Duration(3) * time.Second)
			}
			if pushErr == nil {
				break
			}
		}
		if pushErr != nil {
			// 最终推送失败，记 log
			logx.GetLogger().Error("[PushMessage] 消息推送失败，无重试次数", zap.String(LogMqMessageKey, jsonx.ToJsonIgnoreErr(message)))
			return nil, pushErr
		}
		logx.GetLogger().Info("[PushMessage] 消息发送成功 %v", zap.String(LogMqMessageKey, jsonx.ToJsonIgnoreErr(message)))
		msgExts[i] = mq_model.MqMessageExt{
			MqMessage: mq_model.MqMessage{
				Topic:     msg.Topic,
				Body:      message.Body,
				Keys:      message.Keys,
				Partition: msg.Partition,
				Offset:    msg.Offset,
			},
		}
	}
	return &msgExts, nil
}
