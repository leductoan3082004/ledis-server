package list_commands

import (
	"fmt"
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

	key := args[0]
	v, exist := cmd.rds.GetOrExpired(key)

	if !exist {
		return nil, utils.ErrKeyDoesNotExist(key)
	}

	if v.Type() != utils.ListType {
		return nil, fmt.Errorf("type mismatch, this key is not a list, it is a %s", utils.TypeToString[v.Type()])
	}

	cmd.rds.Lock()
	defer cmd.rds.Unlock()

	l := v.(*types.ListType)
	res := l.LPop()
	if res != nil {
		return *res, nil
	}
	return res, nil
}
