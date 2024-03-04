package controller

import (
	"remembrance/app/common"
	"remembrance/app/models"
)

// 个人相册
func Creat_album_photo(albumid uint, photoid uint) {
	re := models.PersonalAlbum_Photo{
		PersonalAlbum_id: albumid,
		Photo_id:         photoid,
	}
	common.DB.Create(&re)
}

// func Create_photo_group(photoid uint, groupid uint) {
// 	re := models.Photo_Group{
// 		Photo_id:         photoid,
// 		GroupAndAlbum_id: groupid,
// 	}
// 	common.DB.Create(&re)

// }

// 共同评论
func CreateCommonComment(userid uint, photoid uint, text string) {
	re := models.CommonComment{
		User_id:  userid,
		Photo_id: photoid,
		Text:     text,
	}
	common.DB.Create(&re)
}

// 用户与群
func CreatUser_Group(userid uint, groupid uint, purview string) {
	re := models.User_Group{
		User_id:  userid,
		Group_id: groupid,
		Purview:  purview,
	}
	common.DB.Create(&re)
}
