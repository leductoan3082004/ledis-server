package common_commands

import (
	"errors"
	"ledis-server/redis"
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
		return nil, errors.New("expecting 1 arguments")
	}
	return cmd.rds.TTL(args[0])
}
