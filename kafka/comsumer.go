package kafka

import (
	"github.com/bsm/sarama-cluster"
	"fmt"
	"time"
	"github.com/Shopify/sarama"
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
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer comsumer.Close()
	go func(c *cluster.Consumer) {
		errors := c.Errors()
		noti := c.Notifications()
		for {
			select {
			case err := <-errors:
				fmt.Printf("[ERROR] %s\n", err.Error())
			case <-noti:
			}
		}
	}(comsumer)
	for msg := range comsumer.Messages() {
		c.dataChan <- msg.Value
		comsumer.MarkOffset(msg, "") //MarkOffset 并不是实时写入kafka，有可能在程序crash时丢掉未提交的offset
	}
	return nil
}
