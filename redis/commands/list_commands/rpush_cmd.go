package list_commands

import (
	"ledis-server/redis"
	"ledis-server/redis/types"
	"ledis-server/utils"
)

type rpushCmd struct {
	rds redis.Redis
}

func NewRPushCmd(rds redis.Redis) redis.ICommandHandler {
	return &rpushCmd{rds: rds}
}

func (cmd *rpushCmd) CommandName() string {
	return "RPUSH"
}

// Execute RPUSH key value1 value2, ...
func (cmd *rpushCmd) Execute(args ...string) (any, error) {
	if len(args) <= 1 {
		return nil, utils.ErrArgsLengthNotMatch
	}

	cmd.rds.Lock()
	defer cmd.rds.Unlock()

	key := args[0]
	v, exist := cmd.rds.GetOrExpired(key)

	if !exist {
		cmd.rds.Set(key, types.NewListType())
	}

	v, _ = cmd.rds.GetOrExpired(key)

	if v.Type() != utils.ListType {
		return nil, utils.ErrTypeMismatch(utils.ListType, v.Type())
	}

	l := v.(*types.ListType)
	for i := 1; i < len(args); i++ {
		l.RPush(&args[i])
	}

	return l.LLen(), nil
}
