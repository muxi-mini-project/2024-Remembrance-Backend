package gorm

import (
	"remembrance/app/common"
	"remembrance/app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Newgorm() *gorm.DB {
	db, err := gorm.Open(mysql.Open(common.Dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}
	//迁移
	db.AutoMigrate(&models.User{})                //用户
	db.AutoMigrate(&models.EmailCode{})           //邮箱验证码
	db.AutoMigrate(&models.PersonalPhoto{})       //个人照片
	db.AutoMigrate(&models.PersonalAlbum{})       //个人相册
	db.AutoMigrate(&models.PersonalAlbum_Photo{}) //个人照片与相册
	db.AutoMigrate(&models.CommonPhoto{})         //共同照片
	db.AutoMigrate(&models.CommonAlbum{})         //共同相册
	db.AutoMigrate(&models.CommonComment{})       //共同评论
	db.AutoMigrate(&models.Group{})               //群
	db.AutoMigrate(&models.User_Group{})          //用户与群
	db.AutoMigrate(&models.GroupPhoto{})          //多人照片
	db.AutoMigrate(&models.GroupComment{})        //多人评论
	db.AutoMigrate(&models.GroupCode{})           //群验证码
	return db
}
