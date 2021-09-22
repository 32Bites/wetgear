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
			name := "Commands"
			description := "Here are the commands."
			aliases := ""

			subCommandsOrig := ctx.Command.Router.Commands

			if command != nil {
				name = command.Name
				if name == "" {
					name = command.Aliases[0]
				}
				if command.Description == "" {
					description = ""
				} else {
					description = "```\n" + command.Description + "\n```"
				}
				aliases = strings.Join(command.Aliases, " ")

				subCommandsOrig = command.SubCommands
			}

			subCommands := ""
			for name := range subCommandsOrig {
				subCommands += name + " "
			}

			message := fmt.Sprintf("Name: %s\nDescription: %s\nAliases: %s\nSub Commands: %s\n",
				name,
				description,
				aliases,
				subCommands)
			session.ChannelMessageSendReply(ctx.MessageCreate.ChannelID, message, ctx.MessageCreate.Reference())
		}
	}
}

func helpExecutor(ctx Context) {
	if session := ctx.GetSession(); session != nil {
		if len(ctx.Arguments) == 0 {
			ctx.Command.Router.HelpSettings.FoundExecutor(nil)(ctx)
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
				ctx.Command.Router.HelpSettings.NotFoundExecutor(currentCommand, arg.Raw())(ctx)
				return
			}
		}

		ctx.Command.Router.HelpSettings.FoundExecutor(currentCommand)(ctx)
	}
}
