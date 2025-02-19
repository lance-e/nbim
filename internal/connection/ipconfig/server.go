package ipconfig

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunMain() {
	ctx := context.Background()

	//先启动etcd服务发现
	go DataHandler(&ctx)

	//mock测试
	/*  testctx := context.Background() */
	/* TestEtcd(&testctx, "127.0.0.1", "7896", "node1") */
	/* TestEtcd(&testctx, "192.168.0.1", "7897", "node2") */
	/* TestEtcd(&testctx, "186.0.0.2", "7898", "node3") */

	//http服务
	engine := gin.Default()
	engine.GET("/ip/list", func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "internal server error",
				})
			}
		}()

		eds := Dispatch(c)
		//取前五个
		if len(eds) > 5 {
			eds = eds[:5]
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": eds,
		})

	})
	engine.Run(":6789")
}
