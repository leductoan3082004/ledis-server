package common_commands

import (
	"ledis-server/redis"
	"ledis-server/utils"
)

type delCmd struct {
	rds redis.Redis
}

func NewDelCmd(rds redis.Redis) redis.ICommandHandler {
	return &delCmd{rds: rds}
}

func (cmd *delCmd) CommandName() string {
	return "DEL"
}

func (cmd *delCmd) Execute(args ...string) (any, error) {
	if len(args) != 1 {
		return nil, utils.ErrArgsLengthNotMatch
	}
	cmd.rds.Delete(args[0])
	return nil, nil
}
