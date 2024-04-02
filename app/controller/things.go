package controller

import (
	"remembrance/app/common/tool"
	"remembrance/app/response"

	"github.com/gin-gonic/gin"
)

func GetRandString(c *gin.Context) {
	var mes Message
	c.BindJSON(mes)
	str := tool.Randnum(mes.Number)
	response.OkData(c, str)
}
