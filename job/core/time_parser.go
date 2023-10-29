package core

import (
	"github.com/gorhill/cronexpr"
	"time"
)

const (
	Cron     = 1
	Interval = 2
)

type TimeParser interface {
	Next(time.Time) time.Time
	NextTime() time.Time
}

type CronTimeParser struct {
	expression *cronexpr.Expression
}

func NewCronTimeParser(cron string) *CronTimeParser {
	parse, _ := cronexpr.Parse(cron)
	return &CronTimeParser{
		expression: parse,
	}
}

func (parser *CronTimeParser) Next(now time.Time) time.Time {
	return parser.expression.Next(now)
}

func (parser *CronTimeParser) NextTime() time.Time {
	return parser.Next(time.Now())
}

type FixTimeParser struct {
	intervalTime int64
}

func NewFixTimeParser(intervalTime int64) *FixTimeParser {
	return &FixTimeParser{
		intervalTime: intervalTime,
	}
}

func (parser *FixTimeParser) Next(now time.Time) time.Time {
	return now.Add(time.Second * time.Duration(parser.intervalTime))
}

func (parser *FixTimeParser) NextTime() time.Time {
	return parser.Next(time.Now())
}
