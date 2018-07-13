package kafka

import (
	"testing"
	"strings"
	"fmt"
)

func TestKafComsumer_Comsumer(t *testing.T) {
	dataChan := make(chan []byte)
	comsumer := KafComsumer{
		Addr:     strings.Split("127.0.0.1:9092", ","),
		Topics:   strings.Split("topic_test", ","),
		GroupId:  "group_id",
		dataChan: dataChan,
	}
	err := comsumer.Comsumer()
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}

	go func() {
		for bs := range dataChan {
			fmt.Printf("%s\n", bs)
		}
	}()
}
