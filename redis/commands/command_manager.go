package commands

import (
	"ledis-server/logging"
	"ledis-server/redis"
	"ledis-server/redis/commands/list_commands"
	"ledis-server/redis/commands/string_commands"
	"ledis-server/utils"
)

type commandManager struct {
	commandMapper map[string]ICommandHandler
	rds           redis.Redis
}

func NewCommandManager(rds redis.Redis) ICommandManager {
	commandManager := &commandManager{
		commandMapper: make(map[string]ICommandHandler),
	}
	commandManager.
		Register(string_commands.NewGetCmd(rds)).
		Register(string_commands.NewSetCmd(rds))

	commandManager.Register(list_commands.NewLLenCmd(rds))
	return commandManager
}

func (cm *commandManager) Register(handler ICommandHandler) ICommandManager {
	if _, ok := cm.commandMapper[handler.CommandName()]; ok {
		logging.GetLogger().Fatalln(utils.ErrCommandRegisteredDuplicate(handler.CommandName()))
	}
	cm.commandMapper[handler.CommandName()] = handler
	return cm
}

func (cm *commandManager) Execute(command string, args ...string) (any, error) {
	handler, ok := cm.commandMapper[command]
	if !ok {
		return nil, utils.ErrCommandDoesNotExist
	}
	return handler.Execute(args...)
}
