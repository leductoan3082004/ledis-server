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
	cmd.rds.Lock()
	defer cmd.rds.Unlock()

	if err := cmd.rds.LoadSnapshot(); err != nil {
		return "Failed to restore snapshot, please try again later !!!", err
	}
	return "Restore snapshot successfully !!! 😁", nil
}
