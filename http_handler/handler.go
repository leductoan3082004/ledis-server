package httphandler

import (
	"github.com/gin-gonic/gin"
	"ledis-server/middleware"
	"ledis-server/utils"
)

func StartHTTPHandler() error {
	r := gin.Default()
	r.Use(middleware.Recover())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/", func(c *gin.Context) {
		var redisCommandRequest utils.RedisCommandRequest
		if err := c.ShouldBind(&redisCommandRequest); err != nil {
			panic(utils.ErrInvalidRequest(err))
		}
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

	return nil
}
