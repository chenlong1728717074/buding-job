package utils

import "buding-job/common/log"

func Recover(msg string) {
	if err := recover(); err != nil {
		log.GetLog().Error(msg, ":", err)
	}
}
