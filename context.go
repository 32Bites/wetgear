package wetgear

import "github.com/bwmarrin/discordgo"

type Context struct {
	MessageCreate *discordgo.MessageCreate
	Command       *Command
	Alias string
	Arguments     []Argument
}

func (c *Context) GetSession() *discordgo.Session {
	if c.Command == nil {
		return nil
	}

	return c.Command.Session
}