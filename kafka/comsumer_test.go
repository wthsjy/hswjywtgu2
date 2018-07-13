package kafka

import (
	"testing"
	"strings"
	"fmt"
)

func TestKafComsumer_Comsumer(t *testing.T) {
	dataChan := make(chan interface{})
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
		for v := range dataChan {
			data, ok := v.(string)
			if ok {
				fmt.Printf("%s\n", data)
			}
		}
	}()
}
