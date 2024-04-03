package controller

import (
	"errors"
	"fmt"
	"remembrance/app/common"
	"remembrance/app/models"

	"gorm.io/gorm"
)

type Message struct {
	Number            int    `json:"num"`
	UserId            uint   `json:"userid"`
	GroupId           uint   `json:"groupid"`
	PhotoId           uint   `json:"photoid"`
	Email             string `json:"email"`
	Cloudurl          string `json:"cloudurl"`
	Location          string `json:"location"`
	PersonalAlbumName string `json:"personalalbumname"`
	PersonalAlbumId   uint   `json:"personalAlbumId"`
	Text              string `json:"text"`
	GroupName         string `json:"groupname"`
	GroupCode         string `json:"groupcode"`
	IfKeepGroupPhoto  bool   `json:"ifkeepgroupphoto"`
}

func (mes Message) GetPersonalAlbum() models.PersonalAlbum {
	m := models.PersonalAlbum{
		PersonalAlbumName: mes.PersonalAlbumName,
		User_id:           mes.UserId,
	}
	return m
}

func (mes Message) GetPersonalPhoto() models.PersonalPhoto {
	m := models.PersonalPhoto{
		User_id:  mes.UserId,
		Cloudurl: mes.Cloudurl,
		Text:     mes.Text,
	}
	return m
}

func (mes Message) GetGroupPhoto() models.GroupPhoto {
	m := models.GroupPhoto{
		Group_id: mes.GroupId,
		User_id:  mes.UserId,
		Cloudurl: mes.Cloudurl,
		Text:     mes.Text,
	}
	return m
}

func (mes Message) GetUser() (error, models.User) {
	var user models.User
	if err := common.DB.Table("users").Where("id = ?", mes.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有找到匹配的记录
			return err, user
		} else {
			// 其他查询错误
			fmt.Printf("查询错误: %s\n", err.Error())
			return err, user
		}
	} else {
		// 找到匹配的记录，可以使用 user 变量
		return nil, user
	}
}

func (mes Message) GetGroup() models.Group {
	group := models.Group{
		Name: mes.GroupName,
		Code: mes.GroupCode,
	}
	return group

}

func (mes Message) GetUser_Group() models.User_Group {
	group_user := models.User_Group{
		User_id:  mes.UserId,
		Group_id: mes.GroupId,
	}
	return group_user
}

func (mes Message) GetSearch() models.Search {
	search := models.Search{
		User_id: mes.UserId,
		Text:    mes.Location,
	}
	return search
}

// func (mes Message)GetGroupCode() models.GroupCode{
// 	code := models.GroupCode{
// 		Group_id:  mes.GroupId,
// 		Code:      mes.GroupCode,
// 	}
// 	return code
// }
