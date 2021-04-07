package wetgear

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Context struct {
	MessageCreate *discordgo.MessageCreate
	Command       *Command
	Alias         string
	Arguments     []Argument
}

func (c *Context) GetSession() *discordgo.Session {
	if c.Command == nil {
		return nil
	}

	return c.Command.Session
}

func (c *Context) ArgumentsString() string {
	args := make([]string, 0)
	for _, arg := range c.Arguments {
		if arg.Quoted() {
			args = append(args, arg.SurroundQuotes())
		} else {
			args = append(args, arg.Raw())
		}
	}

	return strings.Join(args, " ")
}
