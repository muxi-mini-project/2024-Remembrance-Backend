package controller

import (
	"errors"
	"fmt"
	"remembrance/app/common"
	"remembrance/app/common/email"
	"remembrance/app/models"
	"remembrance/app/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var mes email.Message

// @Summary		登录
// @Description	获取email password 并检查是否正确
// @Tags			login
// @Accept			application/json
// @Produce		application/json
// @Param			email		body		models.User			true	"email"
// @Param			password	body		models.User			true	"password"
// @Success		200			{object}	response.OkMesData		`{"message":"登录成功"}`
// @Failure		400			{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/login [post]
func Login(c *gin.Context) {
	var loguser, user models.User
	// 使用 ShouldBindJSON 解析请求中的 JSON 数据并将其绑定到 user 结构体
	if err := c.ShouldBindJSON(&loguser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//查询用户信息
	err := common.DB.Table("users").First(&user, "Email = ?", loguser.Email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有找到匹配的记录
			response.FailMsg(c, "该用户不存在")
		} else {
			// 其他查询错误
			response.FailMsg(c, "other")
			//fmt.Printf("查询错误: %s\n", err.Error())
			//fmt.Println(user)
		}
	}

	if user.Password == loguser.Password {
		//密码正确
		response.OkMsgData(c, user.Email, user.ID)
	} else {
		//密码错误
		response.Fail(c)
	}

}

// @Summary		获取验证码
// @Description	根据情况获取不同时限的验证码
// @Tags			login
// @Accept			application/json
// @Produce		application/json
// @Param			email	body		email.Message		true	"email"
// @Param			gettype	body		email.Message		true	"请求类型:注册'register',改密码'change'"
// @Success		200		{object}	response.OkMesData		"{"message":"获取成功"}"
// @Failure		400		{object}	response.FailMesData	"{"message":"Failure"}"
// @Router			/api/get_code [get]
func Get_code(c *gin.Context) {
	//需要 用户邮箱 验证码用途（改密码，注册）
	// 使用 ShouldBindJSON 解析请求中的 JSON 数据并将其绑定到 mes 结构体
	if err := c.ShouldBindJSON(&mes); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	switch mes.Type {
	case "register":
		{ //检查账户是否存在
			var user models.User
			if err := common.DB.Table("users").First(&user, "Email = ?", mes.Email).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 没有找到匹配的记录
					//向目标邮箱发送验证码
					email.SendCode(mes.Email, mes.Type)
					response.OkMsg(c, "验证码已发送")
				} else {
					// 其他查询错误
					fmt.Printf("查询错误: %s\n", err.Error())
				}
			} else {
				// 找到匹配的记录，可以使用 user 变量
				//fmt.Printf("找到用户记录: %+v\n", user)
				response.FailMsg(c, "该邮箱已存在")
			}
		}
	case "change":
		{ //检查账户是否存在
			var user models.User
			if err := common.DB.Table("users").First(&user, "Email = ?", mes.Email).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 没有找到匹配的记录
					response.FailMsg(c, "该邮箱未注册")
				} else {
					// 其他查询错误
					fmt.Printf("查询错误: %s\n", err.Error())
				}
			} else {
				// 找到匹配的记录，可以使用 user 变量
				//向目标邮箱发送验证码
				email.SendCode(mes.Email, mes.Type)
				response.OkMsg(c, "验证码已发送")
			}

		}
	}

}

// @Summary		检查验证码
// @Description	根据情况检查验证码 会返回 "验证码不存在或已过期"/"验证码正确"/"验证码错误"
// @Tags			login
// @Accept			application/json
// @Produce		application/json
// @Param			email	body		email.Message		true	"email"
// @Param code body email.Message true "code"
// @Param			gettype	body		email.Message		true	"请求类型:注册'register',改密码'change'"
// @Success		200		{object}	response.OkMesData		"{"message":"成功"}"
// @Failure		400		{object}	response.FailMesData	"{"message":"Failure"}"
// @Router			/api/check_code [post]
func Check_Code(c *gin.Context) {
	//需要 用户邮箱 验证码 验证码用途（改密码，注册）
	var usermes email.Message
	//获取邮箱验证码
	if err := c.ShouldBindJSON(&usermes); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//检查验证码
	status, st := usermes.CheckCode()
	response.Message(c, status, st)

}

// @Summary		注册
// @Description	注册
// @Produce		json
// @Tags			login
// @Accept			application/json
// @Produce		application/json
// @Param			email		body		models.User		true	"email"
// @Param			password	body		models.User		true	"password"
// @Success		200			{object}	response.OkMesData	`{"message":"注册成功"}`
// @Failure		400		{object}	response.FailMesData	"{"message":"Failure"}"
// @Router			/api/register [put]
func Register(c *gin.Context) {
	//获取信息
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//为用户建表
	common.DB.Create(&user)
	//创建一个默认相册
	album := models.PersonalAlbum{
		PersonalAlbumName: "我的记忆",
		User_id:           user.ID,
		Photo_num:         0,
	}
	common.DB.Create(&album)
	response.OkMsg(c, "注册成功")
}
