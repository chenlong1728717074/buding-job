package log

import (
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"strings"
	"time"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	coloredFormatter := NewColoredFormatter()
	log.SetFormatter(coloredFormatter)
	// 设置 logrus 输出到文件
	log.AddHook(NewCustomLfsHook("logs", logrus.InfoLevel, coloredFormatter))
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

type ColoredFormatter struct {
	logrus.TextFormatter
}

func NewColoredFormatter() *ColoredFormatter {
	return &ColoredFormatter{
		TextFormatter: logrus.TextFormatter{
			ForceColors:            true,
			DisableTimestamp:       true, // 禁用 TextFormatter 的时间戳
			DisableLevelTruncation: true, // 禁用 TextFormatter 的日志级别
		},
	}
}

func (f *ColoredFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := f.levelColor(entry.Level)
	level := fmt.Sprintf("%s%s\x1b[0m", levelColor, strings.ToUpper(entry.Level.String()))
	message := fmt.Sprintf("[%s] %s - %s", entry.Time.Format(time.RFC3339), level, entry.Message)
	return []byte(message + "\n"), nil
}

// levelColor 返回与日志级别对应的 ANSI 颜色码
func (f *ColoredFormatter) levelColor(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "\x1b[34m" // 蓝色
	case logrus.InfoLevel:
		return "\x1b[32m" // 绿色
	case logrus.WarnLevel:
		return "\x1b[33m" // 黄色
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return "\x1b[31m" // 红色
	default:
		return "\x1b[0m" // 默认颜色
	}
}
