package kafka

import (
	"fmt"
	"strings"
	"testing"
)

func TestKafComsumer_Comsumer(t *testing.T) {
	comsumer := KafComsumer{
		Addr:     strings.Split("127.0.0.1:9092", ","),
		Topics:   strings.Split("topic_test", ","),
		GroupId:  "group_id",
		dataChan: make(chan []byte),
	}
	err := comsumer.Comsumer()
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}

	go func() {
		for bs := range comsumer.dataChan {
			fmt.Printf("%s\n", bs)
		}
	}()
}
