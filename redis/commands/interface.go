package commands

type ICommandHandler interface {
	CommandName() string
	Execute(args ...string) (any, error)
}

type ICommandManager interface {
	Register(handler ICommandHandler) ICommandManager
	Execute(command string, args ...string) (any, error)
}
