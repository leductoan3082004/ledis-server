package common_commands

import "ledis-server/redis"

type restoreCmd struct {
	rds redis.Redis
}

func NewRestoreCmd(rds redis.Redis) redis.ICommandHandler {
	return &restoreCmd{rds: rds}
}

func (cmd *restoreCmd) CommandName() string {
	return "RESTORE"
}

func (cmd *restoreCmd) Execute(args ...string) (any, error) {
	return nil, cmd.rds.LoadSnapshot()
}
