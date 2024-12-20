package list_commands

import (
	"ledis-server/redis"
	"ledis-server/redis/types"
	"ledis-server/utils"
)

type rpopCmd struct {
	rds redis.Redis
}

func NewRPopCmd(rds redis.Redis) redis.ICommandHandler {
	return &rpopCmd{rds: rds}
}

func (cmd *rpopCmd) CommandName() string {
	return "RPOP"
}

// Execute LLEN key
func (cmd *rpopCmd) Execute(args ...string) (any, error) {
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
	res := l.RPop()
	if res != nil {
		return *res, nil
	}
	return res, nil
}
