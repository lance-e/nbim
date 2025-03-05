package logichttp

import (
	"context"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// 生成元数据，userid和deviceid
func GetMetaData(c *gin.Context) context.Context {
	me := c.PostForm("caller_id")
	device := c.PostForm("device_id")

	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", me,
		"device_id", device,
	))
	return ctx
}
