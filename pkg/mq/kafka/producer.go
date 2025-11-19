package kafka

import (
	"context"
	"github.com/IBM/sarama"
	errgroup "github.com/suhanyujie/go_utils/errorgroup"
	"github.com/suhanyujie/go_utils/logx"
	"os"
	"os/signal"
	"sync"
	"time"
)

type ProducerMessage struct {
	Topic   string
	Key     string
	Payload []byte
}

type Producer struct {
	ctx       context.Context
	cancel    context.CancelFunc
	producers sync.Map
}

type Options struct {
	Address    string `xml:"addr"`
	GameId     string `xml:"group"`
	Group      string `xml:"gameId"`
	Topic      string `xml:"topic"`
	MaxRetry   int    `xml:"retry"`
	MaxTimeout int    `xml:"timeout"`
}

type Conn struct {
	ctx     context.Context
	cancel  context.CancelFunc
	p       sarama.AsyncProducer
	options *Options
}

func NewProducer(options ...*Options) (p *Producer, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	p = &Producer{
		ctx:       ctx,
		cancel:    cancel,
		producers: sync.Map{},
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for _, option := range options {
		conn, err0 := newProducerConn(option)
		if err0 != nil {
			cancel()
			err = err0
			return
		}
		p.producers.Store(option.Topic, conn)

		go func(topic string) {
			for {
				select {
				case suc := <-conn.p.Successes():
					msg, _ := suc.Value.Encode()
					logx.GetLogger().Infof("[ProducerConn] send msg success offset:%d,topic:%s,key:%s,data:%s", suc.Offset, topic, suc.Key, string(msg))
				case fail := <-conn.p.Errors():
					logx.GetLogger().Errorf("[ProducerConn] send msg failed, err: %v", fail.Err.Error())
					connErr := fail.Err
					if connErr == sarama.ErrNotConnected || connErr == sarama.ErrClosedClient || connErr == sarama.ErrOutOfBrokers {
						logx.GetLogger().Errorf("[ProducerConn] client connect close, err: %v", connErr)
						p.producers.Delete(topic)
						return
					}
				case <-signals:
					logx.GetLogger().Infof("[ProducerConn] os interrupt")
					p.Close()
					return
				}
			}
		}(option.Topic)
	}

	return
}

func newProducerConn(options *Options) (c *Conn, err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Timeout = time.Second * time.Duration(options.MaxTimeout)
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = options.MaxRetry

	client, err := sarama.NewClient([]string{options.Address}, config)
	if err != nil {
		logx.GetLogger().Errorf("[ProducerConn] new client err: %v", err)
		return
	}

	p, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		logx.GetLogger().Errorf("[ProducerConn] new async producer err: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(options.MaxTimeout))
	c = &Conn{
		ctx:     ctx,
		cancel:  cancel,
		p:       p,
		options: options,
	}

	return
}

func (p *Producer) Send(message *ProducerMessage) {
	p.SendBatch([]*ProducerMessage{message})
}

func (p *Producer) SendBatch(messages []*ProducerMessage) {
	eg := errgroup.WithContext(p.ctx, logx.GetLogger())
	for _, message := range messages {
		func(msg *ProducerMessage) {
			eg.Go(func(ctx context.Context) error {
				p.sendAsync(msg)
				return nil
			})
		}(message)
	}
	_ = eg.Wait()
}

func (p *Producer) sendAsync(message *ProducerMessage) {
	value, ok := p.producers.Load(message.Topic)
	if ok {
		conn := value.(*Conn)
		msg := &sarama.ProducerMessage{
			Topic: message.Topic,
			Key:   sarama.StringEncoder(message.Key),
			Value: sarama.ByteEncoder(message.Payload),
		}

		conn.p.Input() <- msg
	}

	return
}

func (p *Producer) Close() {
	p.producers.Range(func(key, value any) bool {
		conn := value.(*Conn)
		if err := conn.p.Close(); err != nil {
			logx.GetLogger().Errorf("[ProducerConn] close err: %v", err)
		}
		logx.GetLogger().Infof("[ProducerConn] close finish:%v", key)
		return true
	})
}
