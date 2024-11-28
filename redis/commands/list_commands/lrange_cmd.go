package list_commands

import (
	"fmt"
	"ledis-server/redis"
	"ledis-server/redis/types"
	"ledis-server/utils"
	"strconv"
)

type lrangeCmd struct {
	rds redis.Redis
}

func NewLRangeCmd(rds redis.Redis) redis.ICommandHandler {
	return &lrangeCmd{rds: rds}
}

func (cmd *lrangeCmd) CommandName() string {
	return "LRANGE"
}

// Execute LLEN key
func (cmd *lrangeCmd) Execute(args ...string) (any, error) {
	if len(args) != 3 {
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

	cmd.rds.RLock()
	defer cmd.rds.RUnlock()

	start, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, err
	}

	l := v.(*types.ListType)
	return l.LRange(start, end), nil
}
