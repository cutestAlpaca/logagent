package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client sarama.SyncProducer
	MsgChan chan *sarama.ProducerMessage
)

func Init(address []string, chanSize int64) (err error)  {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //ack
	config.Producer.Partitioner = sarama.NewRandomPartitioner //分区
	config.Producer.Return.Successes = true                   //成功交付

	client, err = sarama.NewSyncProducer([]string{"127.0.0.1:29092"}, config)
	if err != nil {
		logrus.Error("producer close, err:", err)
		return
	}

	MsgChan = make(chan *sarama.ProducerMessage, chanSize)
	go sendMsg()
	return
	//defer client.Close()
}

func sendMsg()  {
	for {
		select {
		case msg := <- MsgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Warn("send msg failed, err:", err)
				return
			}
			logrus.Info(msg.Value)
			logrus.Info("send msg to kafka success. pid:%v offset:%v", pid, offset)
		}
	}
}
