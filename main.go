package main

import (
	"encoding/gob"
	httphandler "ledis-server/http_handler"
	"ledis-server/logging"
	"ledis-server/redis"
	"ledis-server/redis/commands"
	"ledis-server/redis/types"
)

func main() {
	gob.RegisterName("types.ListType", &types.ListType{})
	gob.RegisterName("types.SetType", &types.SetType{})
	gob.RegisterName("types.StringType", &types.StringType{})
	rds := redis.NewRedis()
	commandManager := commands.NewCommandManager(rds)

	if err := httphandler.StartHTTPHandler(commandManager); err != nil {
		logging.GetLogger().Fatalln(err)
	}
}
