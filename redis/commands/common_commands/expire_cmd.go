package common_commands

import (
	"errors"
	"ledis-server/redis"
	"strconv"
)

type expireCmd struct {
	rds redis.Redis
}

func NewExpireCmd(rds redis.Redis) redis.ICommandHandler {
	return &expireCmd{rds: rds}
}

func (cmd *expireCmd) CommandName() string {
	return "EXPIRE"
}

func (cmd *expireCmd) Execute(args ...string) (any, error) {
	if len(args) != 2 {
		return nil, errors.New("expecting 2 arguments")
	}

	ttlInSeconds, err := strconv.Atoi(args[1])

	if err != nil {
		return nil, err
	}

	cmd.rds.Lock()
	defer cmd.rds.Unlock()

	return nil, cmd.rds.Expire(args[0], ttlInSeconds)
}
