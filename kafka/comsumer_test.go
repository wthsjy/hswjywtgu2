package kafka

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

type C struct {
}

func (c *C) Process(bs []byte) {
	fmt.Println(string(bs[:]))
}
func TestKafComsumer_Comsumer(t *testing.T) {
	comsumer := KafComsumer{
		Addr:    strings.Split("127.0.0.1:9092", ","),
		Topics:  strings.Split("topic_test", ","),
		GroupId: "group_id",
		Process: &C{},
	}
	for {
		err := comsumer.Comsumer()
		if err != nil {
			fmt.Printf("error:%s\n", err.Error())
		}
		time.Sleep(time.Second)
	}

}
