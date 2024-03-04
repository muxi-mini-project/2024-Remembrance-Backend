// 获取用户信息
//	@Summary		获取用户信息
//	@Description	根据用户ID获取用户详细信息
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"用户ID"
//	@Success		200	{object}	models.User			"请求成功"
//	@Failure		400	{object}	response.ErrorMsg	"请求失败"
//	@Router			/user/{id} [get]

// 更改密码
//
//	@Summary		更改密码
//	@Description	用户更改自己的密码
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			id				body		int					true	"用户ID"
//	@Param			old_password	body		string				true	"旧密码"
//	@Param			new_password	body		string				true	"新密码"
//	@Success		200				{object}	response.OkMsg		"密码更改成功"
//	@Failure		400				{object}	response.ErrorMsg	"密码更改失败"
//	@Router			/user/password [post]
package controller

import (
	"errors"
	"fmt"
	"remembrance/app/common"
	"remembrance/app/models"
	"remembrance/app/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// @Summary		获取用户信息
// @Description	根据userid获取用户信息
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			email		body		models.User				true	"email"
// @Param			password	body		models.User				true	"password"
// @Success		200			{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400			{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/getinfo [get]
func GetUserInfo(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	common.DB.Table("users").Where("ID = ?", user.ID).First(&user)
	response.OkData(c, user)
}

// @Summary		更改密码
// @Description	前端需保证前后两次邮箱相同，并将邮箱与更改后的密码上传
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			email		body		models.User				true	"email"
// @Param			password	body		models.User				true	"password"
// @Success		200			{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400			{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/changepassword [post]
func ChangePassword(c *gin.Context) {
	//获取信息
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// //检查前后邮箱是否相同
	// if user.Email != mes.Email {
	// 	response.FailMsg(c, "前后邮箱不一致")
	// 	return
	// }
	//查询原密码
	var preuser models.User
	common.DB.Table("users").First(&preuser, "Email = ?", user.Email)
	if user.Password == preuser.Password {
		response.FailMsg(c, "更改后的密码不应与之前相同")
		return
	}
	//改密码
	common.DB.Model(&user).Where("email = ?", user.Email).Update("password", user.Password)
	response.Ok(c)
}

// @Summary		更改昵称
// @Description	保证前后两次邮箱相同，并将邮箱与更改后的密码上传
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			userid	body		models.User				true	"userid"
// @Param			name	body		models.User				true	"name"
// @Success		200		{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/changename [post]
func Changename(c *gin.Context) {
	var user models.User
	//获取信息
	c.BindJSON(&user)
	//更改
	common.DB.Table("users").Where("email = ", user.Email).Update("name", user.Name)
	response.Ok(c)
}

// @Summary		创建群
// @Description	需要创建者的id 群名 code （目前群名不能重复）
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			userid	body		models.Group			true	"userid"
// @Param			name	body		models.Group			true	"name"
// @Param			code	body		models.Group			true	"code"
// @Success		200		{object}	response.OkMesData		`{"message":"创建成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/group/creat [put]
func CreateGroup(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	//检查验证码是否与已生效的重复
	group := mes.GetGroup()
	var g models.Group
	//检查群名是否已经被使用
	if err := common.DB.Table("groups").First(&g, "Name = ?", mes.GroupName).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有找到匹配的记录
			//创建群
			common.DB.Create(&group)
			response.OkMsg(c, "群已创建")
		} else {
			// 其他查询错误
			fmt.Printf("查询错误: %s\n", err.Error())
		}
	} else {
		// 找到匹配的记录，可以使用 user 变量
		//fmt.Printf("找到用户记录: %+v\n", user)
		response.FailMsg(c, "该群名已使用")
	}

	//建立关系
	CreatUser_Group(mes.UserId, group.ID, "creater")
	//记录验证码
	err := common.RDB.Set("verification"+string(rune(group.ID)), group.Code, 10*time.Minute).Err()
	if err != nil {
		panic(err)
	}
	response.Ok(c)
}

// @Summary		加入群
// @Description	需要加入者的id 加入的群名 对应的code
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			userid	body		models.Group			true	"userid"
// @Param			name	body		models.Group			true	"name"
// @Param			code	body		models.Group			true	"code"
// @Success		200		{object}	response.OkMesData		`{"message":"成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/group/join [post]
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

// @Summary		退出群
// @Description	主动退出则传退出者的userid，被踢则传被踢的人的userid 还需要群名
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			userid	body		models.Group			true	"userid"
// @Param			name	body		models.Group			true	"name"
// @Param			code	body		models.Group			true	"code"
// @Success		200		{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/group/out [post]
func OutGroup(c *gin.Context) {
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
