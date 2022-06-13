package main

import (
	"github.com/Shopify/sarama"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"logagent/kafka"
	"logagent/tailfile"
	"time"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
	ChanSize int64 `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

func run() (err error)  {
	for {
		line, ok := <-tailfile.TailTask.Lines
		if !ok {
			logrus.Warn("tail file close reopen filename:%s", tailfile.TailTask.Filename)
			time.Sleep(time.Second)
			continue
		}
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder(line.Text)
		kafka.MsgChan <- msg
	}
	return
}

func main() {
	var configObj = new(Config)
	//cfg, err := ini.Load("./conf/config.ini")
	//if err != nil {
	//	logrus.Error("load", err)
	//	return
	//}
	//kafkaAddress := cfg.Section("kafka").Key("address").String()
	//logrus.Println(kafkaAddress)

	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Error("load config failed, err:%v", err)
		return
	}
	logrus.Printf("%#v\n", configObj)

	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Error("kafka Init failed, err:%v", err)
		return
	}
	logrus.Info("init tail file success")
	err = tailfile.Init(configObj.CollectConfig.LogFilePath)
	if err != nil {
		logrus.Error("kafka Init failed, err:%v", err)
		return
	}

	logrus.Info("kafka Init success")

	err = run()
	if err != nil {
		logrus.Error("run failed, err:%v", err)
		return
	}
}
