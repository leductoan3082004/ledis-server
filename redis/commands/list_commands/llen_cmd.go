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
	cmd.rds.Lock()
	defer cmd.rds.Unlock()

	key := args[0]
	v, exist := cmd.rds.GetOrExpired(key)

	if !exist {
		return nil, utils.ErrKeyDoesNotExist(key)
	}
	if v.Type() != utils.ListType {
		return nil, utils.ErrTypeMismatch(utils.ListType, v.Type())
	}

	l := v.(*types.ListType)
	return l.LLen(), nil
}
