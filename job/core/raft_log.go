package core

import (
	clog "buding-job/common/log"
	"github.com/hashicorp/go-hclog"
	"github.com/sirupsen/logrus"
	"io"
	"log"
)

type RaftLog struct {
	logger *logrus.Logger
	name   string
}

func NewRaftLog() *RaftLog {
	return &RaftLog{
		logger: clog.GetLog(),
	}
}

func (a *RaftLog) Log(level hclog.Level, msg string, args ...interface{}) {
	switch level {
	case hclog.Trace:
		a.logger.Tracef(msg, args...)
	case hclog.Debug:
		a.logger.Debugf(msg, args...)
	case hclog.Info:
		a.logger.Infof(msg, args...)
	case hclog.Warn:
		a.logger.Warnf(msg, args...)
	case hclog.Error:
		a.logger.Errorf(msg, args...)
	}
}
func (a *RaftLog) Trace(msg string, args ...interface{}) {
	a.logger.Tracef(msg, args...)
}

func (a *RaftLog) Debug(msg string, args ...interface{}) {
	a.logger.Debugf(msg, args...)
}

func (a *RaftLog) Info(msg string, args ...interface{}) {
	a.logger.Infof(msg, args...)
}

func (a *RaftLog) Warn(msg string, args ...interface{}) {
	a.logger.Warnf(msg, args...)
}

func (a *RaftLog) Error(msg string, args ...interface{}) {
	a.logger.Errorf(msg, args...)
}

func (a *RaftLog) IsTrace() bool { return a.logger.IsLevelEnabled(logrus.TraceLevel) }
func (a *RaftLog) IsDebug() bool { return a.logger.IsLevelEnabled(logrus.DebugLevel) }
func (a *RaftLog) IsInfo() bool  { return a.logger.IsLevelEnabled(logrus.InfoLevel) }
func (a *RaftLog) IsWarn() bool  { return a.logger.IsLevelEnabled(logrus.WarnLevel) }
func (a *RaftLog) IsError() bool { return a.logger.IsLevelEnabled(logrus.ErrorLevel) }

func (a *RaftLog) ImpliedArgs() []interface{} {
	return nil // 根据需要调整
}

func (a *RaftLog) With(args ...interface{}) hclog.Logger {
	return &RaftLog{
		logger: a.logger.WithFields(logrus.Fields{"args": args}).Logger,
		name:   a.name,
	}
}

func (a *RaftLog) Named(name string) hclog.Logger {
	return &RaftLog{
		logger: a.logger.WithField("name", name).Logger,
		name:   name,
	}
}

func (a *RaftLog) Name() string {
	return a.name
}

func (a *RaftLog) ResetNamed(name string) hclog.Logger {
	return &RaftLog{
		logger: a.logger, // 重置名字，但保持其他配置
		name:   name,
	}
}

func (a *RaftLog) SetLevel(level hclog.Level) {
	a.logger.SetLevel(logrus.Level(level))
}

func (a *RaftLog) GetLevel() hclog.Level {
	return hclog.Level(a.logger.GetLevel())
}

func (a *RaftLog) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(a.StandardWriter(opts), "", 0)
}

func (a *RaftLog) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return a.logger.Out
}
