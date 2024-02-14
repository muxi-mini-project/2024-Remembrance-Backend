package routers

import (
	"remembrance/app/controller"
	"remembrance/app/severs"

	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine {
	a := gin.Default()
	//登录，注册
	LoginGroup := a.Group("/api/login")
	{
		LoginGroup.GET("/get_code", controller.Get_code)      //获取验证码
		LoginGroup.POST("/check_code", controller.Check_Code) //检查验证码
		LoginGroup.PUT("/register", controller.Register)      //注册
		LoginGroup.POST("/login", controller.Login)           //登录
	}
	//用户相关路由
	UserGroup := a.Group("/api/user")
	{
		UserGroup.GET("/getinfo", controller.GetUserInfo)            //获取个人信息
		UserGroup.POST("/changepassword", controller.ChangePassword) //更改密码
		UserGroup.POST("/changename", controller.Changename)         //更改用户名
		UserGroup.PUT("/group/creat", controller.CreateGroup)        //创建群
		UserGroup.POST("/group/join", controller.JoinGroup)          //加入群
		UserGroup.POST("/group/out", controller.OutGroup)            //退出或踢出群
		UserGroup.GET("/group/line", severs.HandleConnections)       //websocket
		//UserGroup.POST("/group/out", controller.OutGroup)            //踢出群

	}

	//记忆相关
	PhotoPost := a.Group("/api/photo")
	{
		//PhotoPost.POST("/test", controller.Test)         //测试上传图片
		PhotoPost.GET("/gettoken", controller.Get_token) //获取token

		PhotoPost.PUT("/personal/createalbum", controller.CreatePersonalAlbum) //创建个人相册
		PhotoPost.GET("/personal/get", controller.GetPersonalPhoto)            //获得个人记忆
		PhotoPost.PUT("/personal/post", controller.PostPersonalPhoto)          //发布个人记忆
		PhotoPost.PUT("/common/photo/post", controller.PostCommonPhoto)        //发布共同记忆
		PhotoPost.GET("/common/photo/get", controller.GetCommonPhoto)          //获取共同记忆
		PhotoPost.PUT("/common/comment/post", controller.PostComment)          //发布共同评论
		PhotoPost.GET("/common/comment/get", controller.GetCommonComment)      //获取共同评论
		PhotoPost.POST("/group/post", controller.PostGroupPhoto)               //发布多人记忆
		PhotoPost.GET("/group/get", controller.GetGroupPhoto)                  //获取多人记忆
	}

	// //用于测试的路由
	// Test := a.Group("/api/test")
	// {
	// 	Test.POST("", controller.Test)
	// }

	return a
}
