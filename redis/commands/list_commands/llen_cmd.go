package list_commands

import (
	"ledis-server/redis"
	"ledis-server/redis/types"
	"ledis-server/utils"
)

type llenCmd struct {
	rds redis.Redis
}

func NewLLenCmd(rds redis.Redis) redis.ICommandHandler {
	return &llenCmd{rds: rds}
}

func (cmd *llenCmd) CommandName() string {
	return "LLEN"
}

// Execute LLEN key
func (cmd *llenCmd) Execute(args ...string) (any, error) {
	if len(args) != 1 {
		return nil, utils.ErrArgsLengthNotMatch
	}

	key := args[0]
	v, exist := cmd.rds.Get(key)

	if !exist {
		return nil, utils.ErrKeyDoesNotExist(key)
	}

	cmd.rds.RLock()
	defer cmd.rds.RUnlock()

	l := v.(*types.ListType)
	return l.LLen(), nil
}
