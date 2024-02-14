package controller

import (
	"fmt"
	"remembrance/app/common"
	"remembrance/app/models"
	"remembrance/app/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	common.DB.Table("users").Where("ID = ?", user.ID).First(&user)
	response.OkData(c, user)
}

// 更改密码
func ChangePassword(c *gin.Context) {
	//获取信息
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user.Email != mes.Email {
		response.FailMsg(c, "前后邮箱不一致")
		return
	}
	//查询原密码
	var preuser models.User
	common.DB.Table("users").First(&preuser, "Email = ?", mes.Email)
	if user.Password == preuser.Password {
		response.FailMsg(c, "更改后的密码不应与之前相同")
		return
	}
	//改密码
	common.DB.Model(&user).Where("email = ?", user.Email).Update("password", user.Password)
	response.Ok(c)
}

// 更改昵称
func Changename(c *gin.Context) {
	var user models.User
	//获取信息
	c.BindJSON(&user)
	//更改
	common.DB.Table("users").Where("email = ", user.Email).Update("name", user.Name)
	response.Ok(c)
}

// 创建群
// 需要创建者id 群名 code
func CreateGroup(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	//检查验证码是否与已生效的重复
	group := mes.GetGroup()
	//创建群
	common.DB.Create(&group)
	//建立关系
	CreatUser_Group(mes.UserId, group.ID, "creater")
	//记录验证码
	err := common.RDB.Set("verification"+string(rune(group.ID)), group.Code, 10*time.Minute).Err()
	if err != nil {
		panic(err)
	}
	response.Ok(c)
}

// 加入群
func JoinGroup(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	group := mes.GetGroup()
	//验证
	val, err := common.RDB.Get("verification" + string(rune(group.ID))).Result()
	if err == redis.Nil {
		fmt.Println("验证码不存在或已过期")
		response.FailMsg(c, "验证码不存在或已过期")
	} else if err != nil {
		panic(err)
	} else if val == mes.GroupCode {
		fmt.Println("验证码正确")
		response.OkMsg(c, "验证码正确")
	} else {
		fmt.Println("验证码错误")
		response.FailMsg(c, "验证码错误")
	}

}

func OutGroup(c *gin.Context) {
	//主动退出则传退出者的userid
	//被踢则传被踢的人的userid
	var mes Message
	c.BindJSON(&mes)
	var group models.Group

	usergroup := mes.GetUser_Group()
	groupphoto := mes.GetGroupPhoto()
	//删除关系
	common.DB.Delete(&usergroup)
	//该群人数减1
	common.DB.Table("Groups").Where("Id = ?", mes.GroupId).First(&group).Update("PeopleNum", group.PeopleNum-1)
	//判断是否保留
	if mes.IfKeepGroupPhoto {
		//保留
		response.Ok(c)
	} else {
		//不保留,删掉
		common.DB.Delete(&groupphoto)
		response.Ok(c)
	}
}
