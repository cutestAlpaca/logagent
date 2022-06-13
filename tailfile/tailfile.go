package tailfile

import (
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

var TailTask *tail.Tail

func Init(filename string) (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: true,
		Poll:      true,
	}

	TailTask, err = tail.TailFile(filename, config)
	if err != nil {
		logrus.Error("tail file %s err:%v", filename, err)
		return
	}
	return
}
