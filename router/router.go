package router

import (
	"HiChat/middleware"
	"HiChat/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")

	user := v1.Group("user")
	{
		user.GET("/list", middleware.JWY() ,service.List)
		user.POST("/login", service.LoginByNameAndPassWord)
		user.POST("/register",service.NewUser)
		user.POST("/users", middleware.JWY() ,service.UpdataUser)
		user.DELETE("/delete", middleware.JWY() ,service.DeleteUser)

		//消息相關
		user.GET("/SendUserMsg", service.SendUserMsg)
	}

	relation := v1.Group("relation").Use(middleware.JWY())
	{
		relation.GET("/list", service.FriendList)
		relation.POST("/add", service.AddFriendByName)
		relation.POST("/newGroup", service.NewGroup)
        relation.POST("/groupList", service.GroupList)
        relation.POST("/joinGroup", service.JoinGroup)
	}

	upload := v1.Group("upload").Use(middleware.JWY())
	{
		upload.POST("/image",service.UploadImage)
	}

	return router
}
