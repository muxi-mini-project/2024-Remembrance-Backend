package codecheck

import (
	"remembrance/app/models"
	"time"
)

func IsCodeValid(code models.GroupCode, typ string) bool {
	durationLimit := 1 * time.Minute
	switch typ {
	case "email":
		durationLimit = 5 * time.Minute
	case "group":
		durationLimit = 10 * time.Minute
	}
	// 设置限时（例如5分钟）
	// 检查时间差异
	currentTime := time.Now()
	if currentTime.Sub(code.TimeStamp) > durationLimit {
		return false //超时
	} else {
		return true //未超时
	}

}
