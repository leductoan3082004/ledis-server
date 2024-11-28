package set_command

import (
	"fmt"
	"ledis-server/redis"
	"ledis-server/redis/types"
	"ledis-server/utils"
)

type sremCmd struct {
	rds redis.Redis
}

func NewSRemCmd(rds redis.Redis) redis.ICommandHandler {
	return &sremCmd{rds: rds}
}

func (cmd *sremCmd) CommandName() string {
	return "SREM"
}

// Execute RPUSH key value1 value2, ...
func (cmd *sremCmd) Execute(args ...string) (any, error) {
	if len(args) < 2 {
		return nil, utils.ErrArgsLengthNotMatch
	}

	key := args[0]
	v, exist := cmd.rds.GetOrExpired(key)

	if !exist {
		return nil, utils.ErrKeyDoesNotExist(key)
	}

	if v.Type() != utils.SetType {
		return nil, fmt.Errorf("type mismatch, this key is not a set, it is a %s", utils.TypeToString[v.Type()])
	}

	cmd.rds.Lock()
	defer cmd.rds.Unlock()

	l := v.(*types.SetType)

	return l.SRem(args[1:]...), nil
}