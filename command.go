package wetgear

import (
	"github.com/bwmarrin/discordgo"
)

type CommandExecutor func(ctx Context)
type CommandMiddleware func(exe CommandExecutor) CommandExecutor

// Command represents a discord command.
type Command struct {
	Description         string
	SubCommands         map[string]*Command
	Aliases             []string
	Name                string
	Router              *Router
	Session             *discordgo.Session
	CommandExecutor     CommandExecutor
	Middlwares          []CommandMiddleware
	RequiredPermissions int64 // Use CombinePermissions to use multiple
}

func NewCommand(router *Router) *Command {
	if router == nil {
		return nil
	}
	if router.Session == nil {
		return nil
	}

	return &Command{
		Router:      router,
		Session:     router.Session,
		SubCommands: map[string]*Command{},
	}
}

func (c *Command) execute(msg *discordgo.MessageCreate, args ...Argument) {
	if msg == nil || len(args) < 1 {
		return
	}
	// Check if subcommand, if subcommand execute and return
	if len(args) > 1 {
		if subCmd, exists := c.SubCommands[args[1].Raw()]; exists {
			subCmd.execute(msg, args[1:]...)
			return
		}
	}

	// Check Permissions
	if c.RequiredPermissions != 0 {
		if msg.Member != nil {
			if perm, err := c.Session.State.MessagePermissions(msg.Message); err != nil {
				return
			} else {
				if (perm & c.RequiredPermissions) < 0 {
					return
				}
			}
		}
	}

	// To avoid crashes
	if c.CommandExecutor == nil {
		c.CommandExecutor = func(ctx Context) {
			Logger.Printf("Command being executed \"%s\" has no CommandExecutor. Please use the SetExecutor method to override this message.\n", ctx.MessageCreate.Content)
		}
	}

	// Apply middlewares
	exec := c.CommandExecutor
	// Global
	for _, middleware := range c.Router.GlobalMiddlwares {
		exec = middleware(exec)
	}
	// Local
	for _, middlware := range c.Middlwares {
		exec = middlware(exec)
	}

	exec(Context{
		MessageCreate: msg,
		Command:       c,
		Arguments:     args[1:],
		Alias:         args[0].Raw(),
	})
}

func (c *Command) AddSubCommand(command *Command) *Command {
	AddCommandToMap(command, c.SubCommands)
	return c
}

func (c *Command) SetExecutor(executor CommandExecutor) *Command {
	c.CommandExecutor = executor
	return c
}

func (c *Command) SetName(name string) *Command {
	c.Name = name
	return c
}

func (c *Command) SetDescription(description string) *Command {
	c.Description = description
	return c
}

func (c *Command) AddAliases(aliases ...string) *Command {
	c.Aliases = append(c.Aliases, aliases...)
	return c
}

func (c *Command) AddMiddlewares(middlewares ...CommandMiddleware) *Command {
	for _, middleware := range middlewares {
		if middleware == nil {
			continue
		}
		c.Middlwares = append(c.Middlwares, middleware)
	}
	return c
}
