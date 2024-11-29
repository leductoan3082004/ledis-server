package httphandler

import (
	"github.com/gin-gonic/gin"
	"ledis-server/logging"
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
				logging.GetLogger().Error(err)
				panic(err)
			}
			logging.GetLogger().Debugf(
				"response for %s command and args %v is: %v", redisCommandRequest.Command, redisCommandRequest.Args,
				res,
			)
			c.JSON(http.StatusOK, utils.SimpleSuccessResponse(res))
		},
	)

	if err := r.Run(":6379"); err != nil {
		panic(err)
	}

	return nil
}
