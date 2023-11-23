package log

import (
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"time"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true, // 强制使用颜色（可选）
	})
	// 设置 logrus 输出到文件
	log.AddHook(NewCustomLfsHook("logs", logrus.InfoLevel, &logrus.TextFormatter{
		ForceColors: true,
	}))
}

func GetLog() *logrus.Logger {
	return log
}

// NewCustomLfsHook 创建一个自定义 LFS Hook，将日志输出到文件
func NewCustomLfsHook(logDir string, level logrus.Level, formatter logrus.Formatter) logrus.Hook {
	// 通过时间生成日志文件名
	fileName := generateLogFileName()

	// 组合日志文件的完整路径
	logPath := filepath.Join(logDir, fileName)
	lfHook := lfshook.NewHook(
		lfshook.WriterMap{
			level: &lumberjack.Logger{
				Filename:   logPath,
				MaxSize:    50, // 每个日志文件的最大尺寸，单位：MB
				MaxBackups: 7,  // 最多保留的旧日志文件数
				MaxAge:     1,  // 保留的旧日志文件的最大天数
				Compress:   true,
			},
		},
		formatter,
	)

	return lfHook
}

// generateLogFileName 生成以时间为名称的日志文件名
func generateLogFileName() string {
	currentTime := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s.log", currentTime)
}
