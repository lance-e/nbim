package logichttp

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	engine := gin.Default()

	engine.POST("/register/device", RegisterDevice) //注册设备
	engine.POST("/signin", SignIn)                  //登陆
	v1 := engine.Group("/v1")
	v1.Use(JWT()) //jwt鉴权
	{
		user := v1.Group("/user")
		user.POST("/get", GetUser)
		user.POST("/update", UpdateUser)
		user.POST("/search", SearchUser)

		friend := v1.Group("/friend")
		friend.POST("/add", AddFriend)
		friend.POST("/agree", AgreeFriend)
		friend.POST("/set", SetFriend)
		friend.POST("/all", GetAllFriends)

		group := v1.Group("/group")
		group.POST("/create", CreateGroup)
		group.POST("/update", UpdateGroup)
		group.POST("/get", GetGroup)
		friend.POST("/all", GetAllGroup)

		groupMember := v1.Group("/group-member")
		groupMember.POST("/add", AddGroupMember)
		groupMember.POST("/update", UpdateGroupMember)
		groupMember.POST("/delete", DeleteGroupMember)
		groupMember.POST("/get", GetGroupMember)
	}

	return engine
}
