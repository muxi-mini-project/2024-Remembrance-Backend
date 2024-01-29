package controller

import (
	"remembrance/app/common"
	"remembrance/app/common/codecheck"
	"remembrance/app/models"
	"remembrance/app/response"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	common.DB.Table("users").Where("ID = ?", user.ID).First(&user)
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

	group := models.Group{
		Name: mes.GroupName,
	}
	//创建群
	common.DB.Create(&group)
	//建立关系
	CreatUser_Group(mes.UserId, group.ID, "creater")
	//记录验证码
	code := models.GroupCode{
		Group_id:  group.ID,
		Code:      mes.GroupCode,
		TimeStamp: time.Now(),
	}
	common.DB.Create(&code)
	response.Ok(c)
}

// 加入群
func JoinGroup(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	code := models.GroupCode{
		Code: mes.GroupCode,
	}
	//验证时间
	if codecheck.IsCodeValid(code, "group") {
		//找到群id
		common.DB.Table("GroupCodes").Where("code = ?", code.Code).First(&code)
		//var group models.Group
		//common.DB.Table("Groups").Where("ID = ?", code.Group_id).First(&group)
		CreatUser_Group(mes.UserId, code.Group_id, "member")
		response.Ok(c)
	} else {
		response.FailMsg(c, "群码已过期")
	}

}

func ExitGroup(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	var group models.Group

	usergroup := models.User_Group{
		Group_id: mes.GroupId,
		User_id:  mes.UserId,
	}
	//删除关系
	common.DB.Delete(&usergroup)
	//该群人数减1
	common.DB.Table("Groups").Where("Id = ?", mes.GroupId).First(&group).Update("PeopleNum", group.PeopleNum-1)
}
