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
// @Description	获取 UserId  AlbumName
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			userid				body		Message					true	"userid"
// @Param			personalalbumname	body		Message					true	"personalalbumname"
// @Success		200					{object}	response.OkMesData		`{"message":"成功"}`
// @Failure		400					{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/personal/createalbum [put]
func CreatePersonalAlbum(c *gin.Context) {
	var album models.PersonalAlbum
	c.ShouldBindJSON(&album)
	//album.Photo_num = 0
	common.DB.Create(&album)
	response.Ok(c)
}

// @Summary		获取个人相册
// @Description	根据 UserId 返回 相册名albumname  相册id albumid
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			userid				body		uint					true	"userid"
// @Success		200					{object}	response.OkMesData		`{"message":"成功"}`
// @Failure		400					{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/personal/getpersonalalbum [get]
func GetPersonalAlbum(c *gin.Context) {
	var mes Message
	c.ShouldBindJSON(&mes)
	var album []models.PersonalAlbum
	common.DB.Limit(10).Table("personalalbums").Where("User_id = ?", mes.UserId).Find(&album)
	response.OkData(c, album)
}

// @Summary		获取个人相册指定照片
// @Description	根据根据个人相册的id 返回对应的照片
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			userid				body		Message					true	"userid"
// @Success		200					{object}	response.OkMesData		`{"message":"成功"}`
// @Failure		400					{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/personal/getfromalbum [get]
func GetPersonalPhotoFromAlbum(c *gin.Context) {
	var mes Message
	c.ShouldBindJSON(&mes)
	album := mes.GetPersonalAlbum()
	var photo []models.PersonalAlbum_Photo
	common.DB.Limit(10).Table("personalalbums").Where("PersonalAlbum_id = ?", album.ID).Find(&photo)
	i := 0
	personalphotos := make([]models.PersonalPhoto, 0)
	for _, p := range photo {
		common.DB.Table("PersonalPhotos").Where("id = ?", p.Photo_id).First(&personalphotos[i])
		i++
	}
	response.OkData(c, personalphotos)
}

// @Summary		发布个人记忆
// @Description	需要 UserId 图片url text
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			cloudurl	body		Message					true	"cloudurl"
// @Param			text		body		Message					true	"text"
// @Param			userid		body		Message					true	"userid"
// @Success		200			{object}	response.OkMesData		`{"message":"成功"}`
// @Failure		400			{object}	response.FailMesData	`{"message":"Failure"}`
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

// @Summary		获取个人记忆
// @Description	根据 UserId 返回个人记忆
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			userid	body		Message					true	"userid"
// @Success		200		{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/personal/get [get]
func GetPersonalPhoto(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	var photos []models.PersonalPhoto
	common.DB.Limit(7).Table("personalphotos").Where("User_id = ?", mes.UserId).Find(&photos)
	response.OkData(c, photos)
}

// @Summary		发布共同记忆
// @Description	需要 UserId 图片url text location
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			cloudurl	body		Message					true	"cloudurl"
// @Param			text		body		Message				true	"text"
// @Param			userid		body		Message					true	"userid"
// @Param			location	body		Message					true	"location"
// @Success		200			{object}	response.OkMesData		`{"message":"成功"}`
// @Failure		400			{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/common/photo/post [put]
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

// @Summary		发布多人记忆
// @Description	需要 UserId groupid 图片url text
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			cloudurl	body		Message					true	"cloudurl"
// @Param			text		body		Message					true	"text"
// @Param			userid		body		Message					true	"userid"
// @Param			groupid		body		Message					true	"groupid"
// @Success		200			{object}	response.OkMesData		`{"message":"成功"}`
// @Failure		400			{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/group/post [put]
func PostGroupPhoto(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	//获取图片信息
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
// func GetGroupPhoto(c *gin.Context) {
// 	var mes Message
// 	c.BindJSON(&mes)
// 	var photos models.GroupPhoto
// 	common.DB.Limit(7).Table("groupphotos").Where("Group_id = ?", mes.GroupId).Find(&photos)
// 	response.OkData(c, photos)

// }

// @Summary		获取共同记忆
// @Description	需要 userid(用于记录搜索历史) location 传回的信息中包含url photoid text
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			location	body		string					true	"location"
// @Success		200			{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400			{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/common/comment/get [get]
func GetCommonPhoto(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	//查找图片
	var photos []models.CommonPhoto
	common.DB.Limit(7).Table("commonphotos").Where("location = ?", mes.Location).Find(&photos)
	//记录
	search := mes.GetSearch()
	common.DB.Create(&search)
	response.OkData(c, photos)

}

// @Summary		获取搜素历史
// @Description	需要 userid 返回5条搜素历史
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			location	body		string					true	"location"
// @Success		200			{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400			{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/common/comment/getsearch [get]
func GetSearch(c *gin.Context) {
	var mes Message
	var search []models.Search
	common.DB.Limit(5).Table("searchs").Where("User_id = ?", mes.UserId).Find(&search)
	response.OkData(c, search)
}

// @Summary		发布共同评论
// @Description	需要 UserId photoid text
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			photoid	body		Message					true	"photoid"
// @Param			text	body		Message					true	"text"
// @Param			userid	body		Message					true	"userid"
// @Param			groupid	body		Message					true	"groupid"
// @Success		200		{object}	response.OkMesData		`{"message":"成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/common/comment/post [put]
func PostComment(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	CreateCommonComment(mes.UserId, mes.PhotoId, mes.Text)
	response.Ok(c)
}

// @Summary		获取共同评论
// @Description	需要 photoid
// @Tags			controller
// @Accept			json
// @Produce		json
// @Param			photoid	body		Message					true	"photoid"
// @Success		200		{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400		{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/common/comment/get [get]
func GetCommonComment(c *gin.Context) {
	var mes Message
	c.BindJSON(&mes)
	var comments []models.CommonComment
	common.DB.Limit(7).Table("comments").Where("Photo_id = ?", mes.PhotoId).Find(&comments)
	response.OkData(c, comments)
}

// @Summary		获取token
// @Description	发送请求即可
// @Tags			controller
// @Accept			json
// @Produce		json
// @Success		200	{object}	response.OkMesData		`{"message":"获取成功"}`
// @Failure		400	{object}	response.FailMesData	`{"message":"Failure"}`
// @Router			/api/photo/gettoken [get]
func Get_token(c *gin.Context) {
	token := tube.GetQNToken()
	response.OkMsg(c, token)
}
