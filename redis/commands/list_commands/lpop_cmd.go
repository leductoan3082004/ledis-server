package list_commands

import (
	"ledis-server/redis"
	"ledis-server/redis/types"
	"ledis-server/utils"
)

type lpopCmd struct {
	rds redis.Redis
}

func NewLPopCmd(rds redis.Redis) redis.ICommandHandler {
	return &lpopCmd{rds: rds}
}

func (cmd *lpopCmd) CommandName() string {
	return "LPOP"
}

// Execute LLEN key
func (cmd *lpopCmd) Execute(args ...string) (any, error) {
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
	res := l.LPop()
	if res != nil {
		return *res, nil
	}
	return res, nil
}
