package common_commands

import "ledis-server/redis"

type snapshotCmd struct {
	rds redis.Redis
}

func NewSnapshotCmd(rds redis.Redis) redis.ICommandHandler {
	return &snapshotCmd{rds: rds}
}

func (cmd *snapshotCmd) CommandName() string {
	return "SNAPSHOT"
}

func (cmd *snapshotCmd) Execute(args ...string) (any, error) {
	return nil, cmd.rds.MakeSnapshot()
}
