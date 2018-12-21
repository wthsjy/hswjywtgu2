package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"time"
)

type KafComsumer struct {
	Addr     []string
	Topics   []string
	GroupId  string
	dataChan chan []byte
	Conf     *cluster.Config
}

func defaultConfig() *cluster.Config {
	conf := cluster.NewConfig()
	conf.Group.Return.Notifications = true
	conf.Consumer.Offsets.CommitInterval = 1 * time.Second
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	return conf
}

// Comsumer kafka 消费
func (c *KafComsumer) Comsumer() error {
	if c.Conf == nil {
		c.Conf = defaultConfig()
	}
	comsumer, err := cluster.NewConsumer(c.Addr, c.GroupId, c.Topics, c.Conf)
	defer func() {
		close(c.dataChan)
		if comsumer != nil {
			comsumer.Close()
		}
	}()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	for {
		select {
		case err := <-comsumer.Errors():
			fmt.Printf("[ERROR] %s\n", err.Error())
			return nil
		case <-comsumer.Notifications():
		case msg := <-comsumer.Messages():
			c.dataChan <- msg.Value
			comsumer.MarkOffset(msg, "") //MarkOffset 并不是实时写入kafka，有可能在程序crash时丢掉未提交的offset
		}
	}
	return nil
}
