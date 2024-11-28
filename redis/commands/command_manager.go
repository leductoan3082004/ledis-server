package commands

import (
	"ledis-server/logging"
	"ledis-server/utils"
)

type commandManager struct {
	commandMapper map[string]ICommandHandler
}

func NewCommandManager() ICommandManager {
	return &commandManager{
		commandMapper: make(map[string]ICommandHandler),
	}
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
