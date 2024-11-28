package commands

import (
	"ledis-server/logging"
	"ledis-server/redis"
	"ledis-server/redis/commands/common_commands"
	"ledis-server/redis/commands/list_commands"
	"ledis-server/redis/commands/string_commands"
	"ledis-server/utils"
)

type commandManager struct {
	commandHandlerMapper map[string]redis.ICommandHandler
	rds                  redis.Redis
}

func NewCommandManager(rds redis.Redis) redis.ICommandManager {
	commandManager := &commandManager{
		commandHandlerMapper: make(map[string]redis.ICommandHandler),
	}
	createCommand := func(newFunc func(redis.Redis) redis.ICommandHandler) redis.ICommandHandler {
		return newFunc(rds)
	}

	commandManager.
		Register(createCommand(string_commands.NewGetCmd)).
		Register(createCommand(string_commands.NewSetCmd)).
		Register(createCommand(list_commands.NewLLenCmd)).
		Register(createCommand(common_commands.NewDelCmd)).
		Register(createCommand(common_commands.NewKeysCmd))

	return commandManager
}

func (cm *commandManager) Register(handler redis.ICommandHandler) redis.ICommandManager {
	if _, ok := cm.commandHandlerMapper[handler.CommandName()]; ok {
		logging.GetLogger().Fatalln(utils.ErrCommandRegisteredDuplicate(handler.CommandName()))
	}
	cm.commandHandlerMapper[handler.CommandName()] = handler
	return cm
}

func (cm *commandManager) Execute(command string, args ...string) (any, error) {
	handler, ok := cm.commandHandlerMapper[command]
	if !ok {
		return nil, utils.ErrCommandDoesNotExist
	}
	return handler.Execute(args...)
}
