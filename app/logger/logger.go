package logger

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Initlogger() {
	// 创建一个日志记录器
	logger := zap.NewExample().Sugar()
	defer logger.Sync() // 释放资源
	// 将GIN框架日志记录到文件中
	file, err := os.OpenFile("gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatalf("无法打开或创建日志文件：%s", err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout) // 将日志同时输出到文件和控制台
}
