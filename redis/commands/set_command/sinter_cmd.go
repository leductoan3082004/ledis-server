package set_command

import (
	"github.com/scylladb/go-set/strset"
	"ledis-server/redis"
	"ledis-server/redis/types"
	"ledis-server/utils"
)

type sInterCmd struct {
	rds redis.Redis
}

func NewSInterCmd(rds redis.Redis) redis.ICommandHandler {
	return &sInterCmd{rds: rds}
}

func (cmd *sInterCmd) CommandName() string {
	return "SINTER"
}

func (cmd *sInterCmd) Execute(args ...string) (any, error) {
	if len(args) < 1 {
		return nil, utils.ErrArgsLengthNotMatch
	}

	cmd.rds.Lock()
	defer cmd.rds.Unlock()

	items := cmd.rds.Gets(args...)
	for _, item := range items {
		if item.Type() != utils.SetType {
			return nil, utils.ErrTypeMismatch(utils.SetType, item.Type())
		}
	}

	setItems := make([]*strset.Set, 0, len(items))
	for _, item := range items {
		v := item.(*types.SetType)
		setItems = append(setItems, v.GetSet())
	}

	return types.SInter(setItems...), nil
}
