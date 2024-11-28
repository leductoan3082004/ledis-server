package common_commands

import "ledis-server/redis"

type flushDBCmd struct {
	rds redis.Redis
}

func NewFlushDBCmd(rds redis.Redis) redis.ICommandHandler {
	return &flushDBCmd{rds: rds}
}

func (cmd *flushDBCmd) CommandName() string {
	return "FLUSHDB"
}

func (cmd *flushDBCmd) Execute(args ...string) (any, error) {
	cmd.rds.FlushDB()
	return nil, nil
}
