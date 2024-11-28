package string_commands

import (
	"ledis-server/redis"
	"ledis-server/redis/commands"
	"ledis-server/utils"
)

type getCmd struct {
	rds redis.Redis
}

func NewGetCmd(rds redis.Redis) commands.ICommandHandler {
	return &getCmd{rds: rds}
}

func (cmd *getCmd) CommandName() string {
	return "GET"
}

// Execute GET KEY
func (cmd *getCmd) Execute(args ...string) (any, error) {
	if len(args) != 1 {
		return nil, utils.ErrArgsLengthNotMatch
	}
	item, exist := cmd.rds.Get(args[0])

	if !exist {
		return nil, utils.ErrKeyDoesNotExist(args[0])
	}

	cmd.rds.RLock()
	defer cmd.rds.RUnlock()
	v := *item.Value().(*string)
	return v, nil
}