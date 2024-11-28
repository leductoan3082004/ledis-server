package string_commands

import (
	"ledis-server/redis"
	"ledis-server/redis/commands"
	"ledis-server/redis/types"
	"ledis-server/utils"
)

type setCmd struct {
	rds redis.Redis
}

func NewSetCmd(rds redis.Redis) commands.ICommandHandler {
	return &setCmd{rds: rds}
}

func (c *setCmd) CommandName() string {
	return "SET"
}

// Execute USAGE SET KEY VALUE
func (c *setCmd) Execute(args ...string) (any, error) {
	if len(args) != 2 {
		return nil, utils.ErrArgsLengthNotMatch
	}
	c.rds.Set(args[0], types.NewStringType(args[1]))
	return args[1], nil
}