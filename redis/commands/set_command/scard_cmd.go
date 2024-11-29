package set_command

import (
	"ledis-server/redis"
	"ledis-server/redis/types"
	"ledis-server/utils"
)

type scardCmd struct {
	rds redis.Redis
}

func NewSCardCmd(rds redis.Redis) redis.ICommandHandler {
	return &scardCmd{rds: rds}
}

func (cmd *scardCmd) CommandName() string {
	return "SCARD"
}

func (cmd *scardCmd) Execute(args ...string) (any, error) {
	if len(args) != 1 {
		return nil, utils.ErrArgsLengthNotMatch
	}

	key := args[0]
	v, exist := cmd.rds.GetOrExpired(key)

	if !exist {
		return nil, utils.ErrKeyDoesNotExist(key)
	}

	if v.Type() != utils.SetType {
		return nil, utils.ErrTypeMismatch(utils.SetType, v.Type())
	}

	cmd.rds.RLock()
	defer cmd.rds.RUnlock()

	l := v.(*types.SetType)

	return l.SCard(), nil
}
