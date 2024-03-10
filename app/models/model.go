package models

import (
	"time"

	"gorm.io/gorm"
)

// 用户
type User struct {
	gorm.Model
	Name            string `json:"username"`
	Password        string `json:"password"`
	Email           string `json:"email"`
	PostPersonalNum int    `json:"pospersonaltnum" gorm:"default:0"`
	PostGroupNum    int    `json:"postgroupnum" gorm:"default:0"`
	PostCommonNum   int    `json:"postcommonnum" gorm:"default:0"`
	StampNum        int    `json:"stampnum" gorm:"default:0"`
}

// 个人照片
type PersonalPhoto struct {
	gorm.Model
	User_id  uint   `json:"userid"`
	Cloudurl string `json:"cloudurl"`
	Text     string `json:"text"`
	//Location string `json:"location"`
}

// 个人相册与照片
type PersonalAlbum_Photo struct {
	gorm.Model
	PersonalAlbum_id uint
	Photo_id         uint
}

// 个人相册
type PersonalAlbum struct {
	gorm.Model
	PersonalAlbumName string `json:"personalalbumname" gorm:"default:'默认相册'"`
	User_id           uint   `json:"userid"`
	Photo_num         int    `gorm:"default:0"`
}

// 共同照片
type CommonPhoto struct {
	gorm.Model
	User_id  uint   `json:"userid"`
	Cloudurl string `json:"cloudurl"`
	Text     string `json:"text"`
	Location string `json:"location"`
}

// 共同相册
type CommonAlbum struct {
	gorm.Model
	Location string `json:"location"`
}

// 共同评论
type CommonComment struct {
	gorm.Model
	User_id  uint   `json:"userid"`
	Photo_id uint   `json:"photoid"`
	Text     string `json:"text"`
}

// 群
type Group struct {
	gorm.Model
	Name      string `jsom:"groupname"`
	PeopleNum int    `json:"peoplenum" gorm:"default:1"`
	Code      string `json:"groupcode"`
}

// 用户与群
type User_Group struct {
	gorm.Model
	User_id  uint   `json:"userid"`
	Group_id uint   `json:"groupid"`
	Purview  string `json:"purview"` //权限 member 成员  creater 群主
}

// 多人照片
type GroupPhoto struct {
	gorm.Model
	Group_id uint   `json:"groupid"`
	User_id  uint   `json:"userid"`
	Cloudurl string `json:"cloudurl"`
	Text     string `json:"text"`
}

// 群评论
type GroupComment struct {
	gorm.Model
	User_id  uint   `json:"userid"`
	Photo_id uint   `json:"photoid"`
	Text     string `json:"text"`
}

// 验证码
type EmailCode struct {
	Email     string `json:"email"`
	Code      string `json:"code"`
	TimeStamp time.Time
}

// 搜索历史
type Search struct {
	gorm.Model
	User_id uint   `json:"userid"`
	Text    string `json:"text"`
}

type GroupCode struct {
	Group_id  uint   `json:"groupid"`
	Code      string `json:"code"`
	TimeStamp time.Time
}

// 仅用于注释
type S_User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt DeletedAt `gorm:"index"`
	Name            string `json:"username"`
	Password        string `json:"password"`
	Email           string `json:"email"`
	PostPersonalNum int    `json:"pospersonaltnum" gorm:"default:0"`
	PostGroupNum    int    `json:"postgroupnum" gorm:"default:0"`
	PostCommonNum   int    `json:"postcommonnum" gorm:"default:0"`
	StampNum        int    `json:"stampnum" gorm:"default:0"`
}

type S_Group struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt DeletedAt `gorm:"index"`
	Name      string `jsom:"groupname"`
	PeopleNum int    `json:"peoplenum" gorm:"default:1"`
	Code      string `json:"groupid"`
}

type S_PersonalAlbum struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt DeletedAt `gorm:"index"`
	PersonalAlbumName string `json:"personalalbumname" gorm:"default:'默认相册'"`
	User_id           uint   `json:"userid"`
	Photo_num         int    `gorm:"default:0"`
}

type S_PersonalPhoto struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt DeletedAt `gorm:"index"`
	User_id  uint   `json:"userid"`
	Cloudurl string `json:"cloudurl"`
	Text     string `json:"text"`
	//Location string `json:"location"`
}

type S_CommonPhoto struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt DeletedAt `gorm:"index"`
	User_id  uint   `json:"userid"`
	Cloudurl string `json:"cloudurl"`
	Text     string `json:"text"`
	Location string `json:"location"`
}

type S_GroupPhoto struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt DeletedAt `gorm:"index"`
	Group_id uint   `json:"groupid"`
	User_id  uint   `json:"userid"`
	Cloudurl string `json:"cloudurl"`
	Text     string `json:"text"`
}

type S_CommonComment struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt DeletedAt `gorm:"index"`
	User_id  uint   `json:"userid"`
	Photo_id uint   `json:"photoid"`
	Text     string `json:"text"`
}
