package commands

import (
	"ledis-server/redis"
	"ledis-server/utils"
)

type getCmd struct {
	rds redis.Redis
}

func NewGetCmd(rds redis.Redis) *getCmd {
	return &getCmd{rds: rds}
}

func (cmd *getCmd) CommandName() string {
	return "GET"
}

func (cmd *getCmd) Execute(args ...string) (any, error) {
	if len(args) != 1 {
		return nil, utils.ErrArgsLengthNotMatch
	}
	cmd.rds.RLock()
	defer cmd.rds.RUnlock()

	item, exist := cmd.rds.Get(args[0])

	if !exist {
		return nil, utils.ErrKeyDoesNotExist(args[0])
	}

	v := *item.Value().(*string)
	return v, nil
}
