package controller

import (
	"errors"
	"fmt"
	"remembrance/app/common"
	"remembrance/app/common/tube"
	"remembrance/app/models"
	"remembrance/app/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary		创建个人相册
// @Description	获取 UserId AlbumName
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			userid				body		uint	true	"userid"
// @Param			personalalbumname	body		string	true	"personalalbumname"
// @Success		200					{object}	response.Okres
// @Failure		400					{object}	response.Failres
// @Router			/api/photo/personal/createalbum [put]
func CreatePersonalAlbum(c *gin.Context) {
	var album models.PersonalAlbum
	c.ShouldBindJSON(&album)
	//album.Photo_num = 0
	common.DB.Create(&album)
	response.Ok(c)
}

// @Summary		发布个人相册
// @Description	获取 UserId AlbumName
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			cloudurl	body		string	true	"cloudurl"
// @Param			text		body		string	true	"text"
// @Param			userid		body		uint	true	"userid"
// @Success		200			{object}	response.Okres
// @Failure		400			{object}	response.Failres
// @Router			/api/photo/personal/post [put]
func PostPersonalPhoto(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	//获取指定信息
	var album models.PersonalAlbum
	var user models.User
	//获取图片信息
	photo := mes.GetPersonalPhoto()
	//找到该用户的个人相册
	common.DB.Table("PersonalAlbum").Where("User_id = ", mes.PersonalAlbumName).First(&album)
	//印记数加一
	common.DB.Model(&user).First(&user, "ID = ?", mes.UserId).Update("StampNum", user.StampNum+1).Update("PostNum", user.PostPersonalNum+1)
	//将图片的url入库
	common.DB.Create(&photo)
	common.DB.Table("photos").Where("Cloudurl = ?", photo.Cloudurl).First(&photo)
	//与相册关联
	Creat_album_photo(album.ID, photo.ID)
	common.DB.Model(&album).First(&album, "id = ?", album.ID).Update("Photo_num", album.Photo_num+1)
	response.Ok(c)
}

// 获取个人记忆
func GetPersonalPhoto(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	var photos []models.PersonalPhoto
	common.DB.Limit(7).Table("personalphotos").Where("User_id = ?", mes.UserId).Find(&photos)
	response.OkData(c, photos)
}

// 发布共同记忆
func PostCommonPhoto(c *gin.Context) {
	//var mes message
	var photo models.CommonPhoto
	c.BindJSON(&photo)
	//入库
	common.DB.Create(&photo)
	//查询相册
	var album models.CommonAlbum
	err := common.DB.Table("commonalbums").Where("location = ?", photo.Location).First(&album).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有找到匹配的记录，创建一个相册
			album.Location = photo.Location
			common.DB.Create(&album)
			response.Ok(c)
		} else {
			// 其他查询错误
			fmt.Printf("查询错误: %s\n", err.Error())
		}
	} else {
		// 找到匹配的记录，可以使用 user 变量
		response.Ok(c)
	}
}

// 发布多人记忆
func PostGroupPhoto(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	//获取指定信息
	// //获取图片信息
	photo := mes.GetGroupPhoto()
	//入库
	common.DB.Create(&photo)
	//获取用户信息
	user := mes.GetUser()
	//印记数加一
	common.DB.Model(&user).First(&user, "ID = ?", user.ID).Update("StampNum", user.StampNum+1).Update("PostNum", user.PostGroupNum+1)
	response.Ok(c)
}

// 获取多人记忆(websocket)
func GetGroupPhoto(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	var photos models.GroupPhoto
	common.DB.Limit(7).Table("groupphotos").Where("Group_id = ?", mes.GroupId).Find(&photos)
	response.OkData(c, photos)

}

// 获取共同记忆
func GetCommonPhoto(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	//查找图片
	var photos []models.CommonPhoto
	common.DB.Limit(7).Table("commonphotos").Where("location = ?", mes.Location).Find(&photos)
	response.OkData(c, photos)

}

// 共同发布评论
func PostComment(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	CreateCommonComment(mes.UserId, mes.PhotoId, mes.Text)
	response.Ok(c)
}

// 获取共同评论
func GetCommonComment(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	var comments []models.CommonComment
	common.DB.Limit(7).Table("comments").Where("Photo_id = ?", mes.PhotoId).Find(&comments)
	response.OkData(c, comments)
}

// 传输token
func Get_token(c *gin.Context) {
	token := tube.GetQNToken()
	response.OkMsg(c, token)
}
