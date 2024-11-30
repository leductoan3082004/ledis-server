package httphandler

import (
	"context"
	"github.com/gin-gonic/gin"
	"ledis-server/logging"
	"ledis-server/middleware"
	"ledis-server/redis"
	"ledis-server/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	server := &http.Server{
		Addr:    ":6379",
		Handler: r,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logging.GetLogger().Info("Starting HTTP server on :6379")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.GetLogger().Errorf("HTTP server ListenAndServe failed: %v", err)
		}
	}()

	sigReceived := <-stopChan
	logging.GetLogger().Infof("Received signal %v, shutting down gracefully...", sigReceived)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logging.GetLogger().Errorf("HTTP server shutdown failed: %v", err)
		return err
	}

	logging.GetLogger().Info("HTTP server gracefully stopped")
	return nil
}
