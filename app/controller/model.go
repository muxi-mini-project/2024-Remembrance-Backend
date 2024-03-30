package controller

import (
	"remembrance/app/common"
	"remembrance/app/models"
)

type Message struct {
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

func (mes Message) GetUser() models.User {
	var user models.User
	user.ID = mes.UserId
	common.DB.First(&user)
	return user
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
	Search := models.Search{
		User_id: mes.UserId,
		Text:    mes.Location,
	}
	return Search
}

// func (mes Message)GetGroupCode() models.GroupCode{
// 	code := models.GroupCode{
// 		Group_id:  mes.GroupId,
// 		Code:      mes.GroupCode,
// 	}
// 	return code
// }
