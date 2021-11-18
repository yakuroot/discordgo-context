package framework

type (
	RunCommand func(ctx Context)

	CommandStruct struct {
		Name    string
		Aliases []string
		Run     RunCommand
	}

	CommandHandler struct {
		cmds []CommandStruct
	}
)

func CreateCommandHandler() *CommandHandler {
	return &CommandHandler{make([]CommandStruct, 0)}
}

func (commandHandler *CommandHandler) RegisterCommand(
	name string,
	aliases []string,
	run RunCommand,
) {
	commandStruct := CommandStruct{
		Name:    name,
		Aliases: aliases,
		Run:     run,
	}

	commandHandler.cmds = append(commandHandler.cmds, commandStruct)
}

func (commandHandler *CommandHandler) GetCommands(commandName string) (RunCommand, bool) {
	for _, command := range commandHandler.cmds {
		if commandName == command.Name {
			return command.Run, true
		}

		for _, aliase := range command.Aliases {
			if aliase == commandName {
				return command.Run, true
			}
		}
	}

	return nil, false
}
