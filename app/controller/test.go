package controller

import (
	"remembrance/app/common"
	"remembrance/app/models"
	"remembrance/app/response"

	"github.com/gin-gonic/gin"
)

// @Summary		测试
// @Description	测试
// @Tags			controller
// @Accept			json
// @Success		200	{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400	{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/test [get]
func Test(c *gin.Context) {
	var mes Message
	c.ShouldBindJSON(&mes)
	response.OkMsgData(c, "访问成功", mes)
	var user models.User
	common.DB.Table("users").Where("id = ?", 1).First(&user)
	response.OkData(c, user)
}
