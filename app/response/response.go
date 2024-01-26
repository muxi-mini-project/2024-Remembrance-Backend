package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS = 200
	FAIL    = 400
)

// 自定义通用消息
func Message(ctx *gin.Context, status int, message string, data ...any) {
	var obj gin.H
	if len(data) == 0 {
		obj = gin.H{
			"code":    status,
			"message": message,
		}
	} else {
		obj = gin.H{
			"code":    status,
			"message": message,
			"data":    data[0],
		}
	}
	ctx.JSON(http.StatusOK, obj)
}

// 默认的成功响应
func Ok(ctx *gin.Context) {
	obj := gin.H{
		"code":    SUCCESS,
		"message": "操作成功",
	}
	ctx.JSON(http.StatusOK, obj)
}

// 携带消息的成功响应
func OkMsg(ctx *gin.Context, message string) {
	obj := gin.H{
		"code":    SUCCESS,
		"message": message,
	}
	ctx.JSON(http.StatusOK, obj)
}

// 携带数据的成功响应
func OkData(ctx *gin.Context, data any) {
	obj := gin.H{
		"code":    SUCCESS,
		"message": "操作成功",
		"data":    data,
	}
	ctx.JSON(http.StatusOK, obj)
}

// 携带消息和数据的成功响应
func OkMsgData(ctx *gin.Context, message string, data any) {
	obj := gin.H{
		"code":    SUCCESS,
		"message": message,
		"data":    data,
	}
	ctx.JSON(http.StatusOK, obj)
}

// 默认的失败响应
func Fail(ctx *gin.Context) {
	obj := gin.H{
		"code":    FAIL,
		"message": "操作失败",
	}
	ctx.JSON(http.StatusOK, obj)
}

// 携带消息的失败响应
func FailMsg(ctx *gin.Context, message string) {
	obj := gin.H{
		"code":    FAIL,
		"message": message,
	}
	ctx.JSON(http.StatusOK, obj)
}

// 携带数据的失败响应
func FailData(ctx *gin.Context, data any) {
	obj := gin.H{
		"code":    FAIL,
		"message": "操作失败",
		"data":    data,
	}
	ctx.JSON(http.StatusOK, obj)
}
