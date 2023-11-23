package log

import "github.com/sirupsen/logrus"

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true, // 强制使用颜色（可选）
	})
}

func GetLog() *logrus.Logger {
	return log
}
