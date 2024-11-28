package set_command

import (
	"fmt"
	"ledis-server/redis"
	"ledis-server/redis/types"
	"ledis-server/utils"
)

type saddCmd struct {
	rds redis.Redis
}

func NewSAddCmd(rds redis.Redis) redis.ICommandHandler {
	return &saddCmd{rds: rds}
}

func (cmd *saddCmd) CommandName() string {
	return "SADD"
}

// Execute RPUSH key value1 value2, ...
func (cmd *saddCmd) Execute(args ...string) (any, error) {
	if len(args) <= 2 {
		return nil, utils.ErrArgsLengthNotMatch
	}

	key := args[0]
	v, exist := cmd.rds.GetOrExpired(key)

	if !exist {
		cmd.rds.Set(key, types.NewSetType())
	}

	v, _ = cmd.rds.Get(key)

	if v.Type() != utils.SetType {
		return nil, fmt.Errorf("type mismatch, this key is not a set, it is a %s", utils.TypeToString[v.Type()])
	}

	v, _ = cmd.rds.Get(key)

	cmd.rds.Lock()
	defer cmd.rds.Unlock()

	l := v.(*types.SetType)

	return l.SAdd(args[1:]...), nil
}
