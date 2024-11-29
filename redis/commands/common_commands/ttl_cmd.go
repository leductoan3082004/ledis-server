package common_commands

import (
	"ledis-server/redis"
	"ledis-server/utils"
)

type ttlCmd struct {
	rds redis.Redis
}

func NewTTLCmd(rds redis.Redis) redis.ICommandHandler {
	return &ttlCmd{rds: rds}
}

func (cmd *ttlCmd) CommandName() string {
	return "TTL"
}

func (cmd *ttlCmd) Execute(args ...string) (any, error) {
	if len(args) != 1 {
		return nil, utils.ErrArgsLengthNotMatch
	}
	cmd.rds.Lock()
	defer cmd.rds.Unlock()
	return cmd.rds.TTL(args[0])
}
