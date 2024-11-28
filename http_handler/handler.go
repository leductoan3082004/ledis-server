package httphandler

import (
	"github.com/gin-gonic/gin"
	"ledis-server/middleware"
	"ledis-server/redis"
	"ledis-server/utils"
	"net/http"
)

func StartHTTPHandler(commandManager redis.ICommandManager) error {
	r := gin.Default()
	r.Use(middleware.Recover())

	r.GET(
		"/ping", func(c *gin.Context) {
			c.JSON(
				200, gin.H{
					"message": "pong",
				},
			)
		},
	)

	r.POST(
		"/", func(c *gin.Context) {
			var redisCommandRequest utils.RedisCommandRequest
			if err := c.ShouldBind(&redisCommandRequest); err != nil {
				panic(utils.ErrInvalidRequest(err))
			}

			res, err := commandManager.Execute(redisCommandRequest.Command, redisCommandRequest.Args...)

			if err != nil {
				panic(err)
			}

			c.JSON(http.StatusOK, utils.SimpleSuccessResponse(res))
		},
	)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

	return nil
}
