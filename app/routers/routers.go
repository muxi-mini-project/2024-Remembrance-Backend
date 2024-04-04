package routers

import (
	"remembrance/app/controller"
	"remembrance/app/severs"

	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine {
	a := gin.Default()
	//登录，注册
	LoginGroup := a.Group("/api")
	{
		LoginGroup.POST("/get_code", controller.Get_code) //获取验证码
		//LoginGroup.GET("/get_code", controller.Get_code)      //获取验证码
		LoginGroup.POST("/check_code", controller.Check_Code) //检查验证码
		LoginGroup.PUT("/register", controller.Register)      //注册
		LoginGroup.POST("", controller.Login)                 //登录
	}
	//用户相关路由
	UserGroup := a.Group("/api/user")
	{
		UserGroup.POST("/getinfo", controller.GetUserInfo)           //获取个人信息
		UserGroup.POST("/changepassword", controller.ChangePassword) //更改密码
		UserGroup.POST("/changename", controller.Changename)         //更改用户名
		UserGroup.POST("/group/get", controller.GetGroup)            //获取群
		UserGroup.PUT("/group/creat", controller.CreateGroup)        //创建群
		UserGroup.POST("/group/join", controller.JoinGroup)          //加入群
		UserGroup.POST("/group/out", controller.OutGroup)            //退出或踢出群
		UserGroup.POST("/group/delete", controller.OutGroup)         //退出或踢出群
		UserGroup.POST("/group/line", severs.HandleConnections)      //websocket
		//UserGroup.POST("/group/out", controller.OutGroup)            //踢出群
		//UserGroup.GET("/getinfo", controller.GetUserInfo)            //获取个人信息
		//UserGroup.GET("/group/get", controller.GetGroup)             //获取群

	}

	//记忆相关
	PhotoPost := a.Group("/api/photo")
	{
		//PhotoPost.POST("/test", controller.Test)         //测试上传图片
		PhotoPost.GET("/gettoken", controller.Get_QNtoken) //获取qntoken

		PhotoPost.PUT("/personal/createalbum", controller.CreatePersonalAlbum)         //创建个人相册
		PhotoPost.POST("/personal/getpersonalalbum", controller.GetPersonalAlbum)      //获取个人相册
		PhotoPost.POST("/personal/deletealbum", controller.DeletePersonalAlbum)        //删除个人相册
		PhotoPost.POST("/personal/getfromalbum", controller.GetPersonalPhotoFromAlbum) //根据相册获得个人记忆
		PhotoPost.POST("/personal/get", controller.GetPersonalPhoto)                   //获得个人记忆
		PhotoPost.POST("/personal/numget", controller.GetNumPersonalPhoto)             //获得指定数量的最新个人记忆
		PhotoPost.PUT("/personal/post", controller.PostPersonalPhoto)                  //发布个人记忆
		PhotoPost.POST("/personal/delete", controller.DeletePersonalPhoto)             //删除个人记忆
		PhotoPost.PUT("/common/photo/post", controller.PostCommonPhoto)                //发布共同记忆
		PhotoPost.POST("/common/photo/delete", controller.DeleteCommonPhoto)           //删除共同记忆
		PhotoPost.POST("/common/photo/getself", controller.GetSelfCommonPhoto)         //获取自己发布的共同记忆
		PhotoPost.POST("/common/photo/get", controller.GetCommonPhoto)                 //获取指定地点共同记忆
		PhotoPost.GET("/common/photo/randget", controller.GetRandCommonPhoto)          //获取随机共同记忆
		PhotoPost.POST("/common/comment/getsearch", controller.GetSearch)              //获取搜索历史
		PhotoPost.PUT("/common/comment/post", controller.PostComment)                  //发布共同评论
		PhotoPost.POST("/common/comment/get", controller.GetCommonComment)             //获取共同评论
		//PhotoPost.GET("/personal/getpersonalalbum", controller.GetPersonalAlbum)      //获取个人相册
		PhotoPost.POST("/group/get", controller.GetGroupPhoto)  //获取多人记忆
		PhotoPost.PUT("/group/post", controller.PostGroupPhoto) //发布多人记忆
	}

	//工具
	ToolGroup := a.Group("api/tool")
	{
		ToolGroup.GET("/getqntoken", controller.Get_QNtoken)       //获取qntoken
		ToolGroup.POST("/getrandstring", controller.GetRandString) //获取指定位数的随机字符

	}

	//用于测试的路由
	Test := a.Group("/api/test")
	{
		Test.Any("", controller.Test)
	}

	return a
}
