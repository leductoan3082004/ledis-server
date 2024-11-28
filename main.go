package main

import (
	httphandler "ledis-server/http_handler"
	"ledis-server/logging"
	"ledis-server/redis"
	"ledis-server/redis/commands"
)

func main() {
	rds := redis.NewRedis()
	commandManager := commands.NewCommandManager()
	commandManager.Register(commands.NewSetCmd(rds))

	if err := httphandler.StartHTTPHandler(); err != nil {
		logging.GetLogger().Fatalln(err)
	}
}
