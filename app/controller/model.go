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
	Text              string `json:"text"`
	GroupName         string `json:"groupname"`
	GroupCode         string `json:"groupcode"`
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
