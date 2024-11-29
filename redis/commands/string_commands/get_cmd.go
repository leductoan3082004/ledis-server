package string_commands

import (
	"ledis-server/redis"
	"ledis-server/utils"
)

type getCmd struct {
	rds redis.Redis
}

func NewGetCmd(rds redis.Redis) redis.ICommandHandler {
	return &getCmd{rds: rds}
}

func (cmd *getCmd) CommandName() string {
	return "GET"
}

// Execute GET KEY
func (cmd *getCmd) Execute(args ...string) (any, error) {
	if len(args) != 1 {
		return nil, utils.ErrArgsLengthNotMatch
	}

	cmd.rds.Lock()
	defer cmd.rds.Unlock()

	item, exist := cmd.rds.GetOrExpired(args[0])

	if !exist {
		return nil, nil
	}

	if item.Type() != utils.StringType {
		return nil, utils.ErrTypeMismatch(utils.StringType, item.Type())
	}

	v := *item.Value().(*string)
	return v, nil
}
