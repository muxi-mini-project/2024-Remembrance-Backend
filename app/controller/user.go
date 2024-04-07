package controller

import (
	"errors"
	"fmt"
	"net/http"
	"remembrance/app/common"
	"remembrance/app/common/tool"
	"remembrance/app/models"
	"remembrance/app/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary		获取用户信息
// @Description	根据userid获取用户信息
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			email		body		models.S_User			true	"email"
// @Param			password	body		models.S_User			true	"password"
// @Success		200			{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400			{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/getinfo [get]
func GetUserInfo(c *gin.Context) {
	var user, u models.User
	c.BindJSON(&user)
	common.DB.Table("users").Where("ID = ?", user.ID).First(&u)
	response.OkData(c, u)
}

// @Summary		更改密码
// @Description	检查验证码正确后，将 邮箱 与更改后的 密码 上传 （保证检查验证码时的邮箱，与上传的邮箱一致）
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			email		body		models.S_User			true	"email"
// @Param			password	body		models.S_User			true	"password"
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
// @Description	保证前后两次邮箱相同，并将邮箱与更改后的昵称上传
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			userid	body		models.S_User			true	"userid"
// @Param			name	body		models.S_User			true	"name"
// @Success		200		{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/changename [post]
func Changename(c *gin.Context) {
	var user models.User
	//获取信息
	c.BindJSON(&user)
	//更改
	common.DB.Table("users").Where("email = ?", user.Email).Update("name", user.Name)
	response.Ok(c)
}

// @Summary		创建群
// @Description	需要创建者的 id 群名 code （目前群名不能重复）
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			userid	body		models.S_Group			true	"userid"
// @Param			name	body		models.S_Group			true	"name"
// @Param			code	body		models.S_Group			true	"code"
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
	if err := common.DB.Table("groups").Where("name = ?", mes.GroupName).First(&g).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有找到匹配的记录
		} else {
			// 其他查询错误
			fmt.Printf("查询错误: %s\n", err.Error())
			response.FailMsgData(c, "未知错误", err)
			return
		}
	} else {
		// 找到匹配的记录
		response.FailMsg(c, "该群名已使用")
		return
	}

	code := tool.Randnumber()
	group.Code = code
	common.DB.Create(&group)
	common.DB.Table("groups").Where("Name = ?", group.Name).First(&group)
	//建立关系
	CreatUser_Group(mes.UserId, group.ID, "creater")

	response.OkData(c, group)
	// //记录验证码
	// err := common.RDB.Set("verification"+string(rune(group.ID)), group.Code, 10*time.Minute).Err()
	// if err != nil {
	// 	panic(err)
	// }
	// common.DB.Table("")

}

// @Summary		加入群
// @Description	需要 加入者的id 加入的群名 对应的code
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			userid	body		models.S_Group			true	"userid"
// @Param			name	body		models.S_Group			true	"name"
// @Param			code	body		models.S_Group			true	"code"
// @Success		200		{object}	response.OkMesData		`{"message":"成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/group/join [post]
func JoinGroup(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	group := mes.GetGroup()
	var g models.Group
	common.DB.Table("groups").Where("Name = ?", group.Name).First(&g)
	if group.Code != g.Code {
		response.FailMsg(c, "code错误")
		return
	}

	CreatUser_Group(mes.UserId, g.ID, "member")

	response.OkData(c, g)
	// //验证
	// val, err := common.RDB.Get("verification" + string(rune(group.ID))).Result()
	// if err == redis.Nil {
	// 	fmt.Println("验证码不存在或已过期")
	// 	response.FailMsg(c, "验证码不存在或已过期")
	// } else if err != nil {
	// 	panic(err)
	// } else if val == mes.GroupCode {
	// 	fmt.Println("验证码正确")
	// 	response.OkMsg(c, "验证码正确")
	// } else {
	// 	fmt.Println("验证码错误")
	// 	response.FailMsg(c, "验证码错误")
	// }

}

// @Summary		退出群
// @Description	主动退出则传退出者的userid，被踢则传被踢的人的userid  还需要群名
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			userid	body		models.S_Group			true	"userid"
// @Param			name	body		models.S_Group			true	"name"
// @Param			code	body		models.S_Group			true	"code"
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
	common.DB.Table("`Groups`").Where("Id = ?", mes.GroupId).First(&group).Update("People_Num", group.PeopleNum-1)
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

// @Summary		解散群
// @Description	传入被解散的群的id就行，确保只有群主能执行这个操作
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			groupid	body		models.S_Group			true	"groupid"
// @Success		200		{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/group/delete [post]
func DeleteGroup(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)

	group := mes.GetGroup()
	//删除群
	common.DB.Delete(&group)
	//删除关系
	var usergroup models.User_Group
	common.DB.Table("User_Groups").Where("Group_id = ?", mes.GroupId).Delete(&usergroup)
	response.Ok(c)
}

// @Summary		获取群
// @Description	传入userid 获得该用户所加入的群信息
// @Tags			user
// @Accept			application/json
// @Produce		application/json
// @Param			userid	body		models.S_Group			true	"userid"
// @Success		200		{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/user/group/get [get]
func GetGroup(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	var gp []models.User_Group
	//获取对应的群id
	common.DB.Limit(20).Table("user_groups").Where("user_id = ?", mes.UserId).Find(&gp)
	//获取对应群信息
	var groups []models.Group
	for _, group := range gp {
		var info models.Group
		common.DB.Limit(20).Table("groups").Where("id = ?", group.Group_id).Find(&info)
		groups = append(groups, info)
	}
	obj := gin.H{
		"code":         200,
		"message":      "获取成功",
		"group":        groups,
		"relationship": gp,
	}
	c.JSON(http.StatusOK, obj)
}
