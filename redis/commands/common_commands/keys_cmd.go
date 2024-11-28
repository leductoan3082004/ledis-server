package common_commands

import (
	"ledis-server/redis"
)

type keyCmd struct {
	rds redis.Redis
}

func NewKeysCmd(rds redis.Redis) redis.ICommandHandler {
	return &keyCmd{rds: rds}
}

func (cmd *keyCmd) CommandName() string {
	return "KEYS"
}

func (cmd *keyCmd) Execute(args ...string) (any, error) {
	return cmd.rds.Keys(), nil
}
