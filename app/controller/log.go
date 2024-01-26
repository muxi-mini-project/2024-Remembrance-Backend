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

//检查用户是否存在

func Login(c *gin.Context) {
	var loguser, user models.User
	// 使用 ShouldBindJSON 解析请求中的 JSON 数据并将其绑定到 user 结构体
	if err := c.ShouldBindJSON(&loguser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//查询用户信息
	common.DB.Table("users").First(&user, "Email = ?", loguser.Email)
	fmt.Println(user)
	//匹配密码
	if user.Password == loguser.Password {
		//密码正确
		response.OkMsgData(c, user.Email, user.ID)
	} else {
		//密码错误
		response.Fail(c)
	}

}

func Get_code(c *gin.Context) {
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
					mes.Code = email.SendCode(mes.Email)
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
				mes.Code = email.SendCode(mes.Email)
				response.OkMsg(c, "验证码已发送")
			}

		}
	}

}

func Check_Code(c *gin.Context) {
	var usermes email.Message
	//获取邮箱验证码
	if err := c.ShouldBindJSON(&usermes); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if mes.Email == usermes.Email {
		if mes.Code == usermes.Code {
			//验证成功
			//重置验证码
			mes.Code = ""
			//发送信息
			response.OkMsg(c, "验证成功")
		} else {
			//向前端返回错误
			response.FailMsg(c, "验证失败")
		}
	}
}

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
}
