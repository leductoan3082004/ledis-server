package set_command

import (
	"ledis-server/redis"
	"ledis-server/redis/types"
	"ledis-server/utils"
)

type sMembersCmd struct {
	rds redis.Redis
}

func NewSMemberCmd(rds redis.Redis) redis.ICommandHandler {
	return &sMembersCmd{rds: rds}
}

func (cmd *sMembersCmd) CommandName() string {
	return "SMEMBERS"
}

// Execute RPUSH key value1 value2, ...
func (cmd *sMembersCmd) Execute(args ...string) (any, error) {
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
	if v.Type() != utils.SetType {
		return nil, utils.ErrTypeMismatch(utils.SetType, v.Type())
	}

	l := v.(*types.SetType)
	return l.SMembers(), nil
}
