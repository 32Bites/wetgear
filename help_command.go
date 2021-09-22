package wetgear

import (
	"fmt"
	"strings"
)

type CommandNotFoundExecutor func(parentCommand *Command, commandName string) CommandExecutor
type CommandFoundExecutor func(command *Command) CommandExecutor

func DefaultCommandNotFound(parentCommand *Command, commandName string) CommandExecutor {
	return func(ctx Context) {
		if session := ctx.GetSession(); session != nil {
			parentName := "the command router"
			if parentCommand != nil {
				parentName = `"` + parentCommand.Name + `"`
			}

			message := fmt.Sprintf("Could not find \"%s\" in %s.", commandName, parentName)
			session.ChannelMessageSendReply(ctx.MessageCreate.ChannelID, message, ctx.MessageCreate.Reference())
		}
	}
}

func DefaultCommandFound(command *Command) CommandExecutor {
	return func(ctx Context) {
		if session := ctx.GetSession(); session != nil {
			name := command.Name
			if name == "" {
				name = command.Aliases[0]
			}

			description := "```\n" + command.Description + "\n```"
			if command.Description == "" {
				description = ""
			}

			subCommands := ""
			for name := range command.SubCommands {
				subCommands += name + " "
			}

			message := fmt.Sprintf("Name: %s\nDescription: %s\nAliases: %s\nSub Commands: %s\n",
				name,
				description,
				strings.Join(command.Aliases, " "),
				subCommands)
			session.ChannelMessageSendReply(ctx.MessageCreate.ChannelID, message, ctx.MessageCreate.Reference())
		}
	}
}

func helpExecutor(ctx Context) {
	if session := ctx.GetSession(); session != nil {
		if len(ctx.Arguments) == 0 {
			return
		}
		var currentCommand *Command = nil

		for _, arg := range ctx.Arguments {
			commandMap := ctx.Command.Router.Commands
			if currentCommand != nil {
				commandMap = currentCommand.SubCommands
			}

			if command, exists := commandMap[arg.Raw()]; exists {
				currentCommand = command
			} else {
				if ctx.Command.Router.HelpSettings.NotFoundExecutor == nil {
					ctx.Command.Router.HelpSettings.NotFoundExecutor = DefaultCommandNotFound
				}

				ctx.Command.Router.HelpSettings.NotFoundExecutor(currentCommand, arg.Raw())(ctx)
				return
			}
		}

		if ctx.Command.Router.HelpSettings.FoundExecutor == nil {
			ctx.Command.Router.HelpSettings.FoundExecutor = DefaultCommandFound
		}

		ctx.Command.Router.HelpSettings.FoundExecutor(currentCommand)(ctx)
	}
}
