package provider

import (
	"fmt"

	"github.com/nsqio/go-nsq"
	"go-chat/config"
)

// NewNsqProducer 初始化生产者
func NewNsqProducer(conf *config.Config) *nsq.Producer {

	nsqConfig := nsq.NewConfig()
	producer, err := nsq.NewProducer(conf.Nsq.Addr, nsqConfig)
	if err != nil {
		panic(fmt.Sprintf("create producer failed, err:%s\n", err.Error()))
	}

	return producer
}
