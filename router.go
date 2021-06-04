package wetgear

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Logger *log.Logger = log.New(os.Stdout, "[WETGEAR] ", log.LstdFlags)

// PrefixSettings defines the settings for a router's prefixes
type PrefixSettings struct {
	Prefixes   []string
	IgnoreCase bool
	HandlePing bool // whether or not to treat a ping as a command prefix. Example: "@BotName help" does the same thing as "!help" if "!" is a prefix.
}

// Router is a command router
type Router struct {
	Session          *discordgo.Session
	Commands         map[string]*Command
	PrefixSettings   PrefixSettings
	BotsAllowed      bool
	GlobalMiddlwares []CommandMiddleware
	RemoveHandler    func()
}

// NewRouter creates a *Router and configures a *discordgo.Session to work with the command system
func NewRouter(session *discordgo.Session, baseRouter *Router) (*Router, error) {
	if session == nil {
		return nil, errors.New("provided discord session is nil")
	}
	if baseRouter == nil {
		return nil, errors.New("provided base router is nil")
	}
	baseRouter.Session = session

	remove := session.AddHandler(func(session *discordgo.Session, msg *discordgo.MessageCreate) {
		// Check if bot and if bots are allowed, and make sure that the message is not empty
		if (msg.Author.Bot && !baseRouter.BotsAllowed) || msg.Content == "" {
			return
		}
		contentRunes := []rune(msg.Content)
		command := ""

		// Check if the message starts with a mention, and that the message is longer than the mention
		if baseRouter.PrefixSettings.HandlePing {
			for _, mention := range baseRouter.GetMentions() {
				if strings.HasPrefix(msg.Content, mention+" ") && len(contentRunes) > len([]rune(mention+" ")) {
					command = strings.TrimSpace(string(contentRunes[len([]rune(mention)):]))
					break
				}
			}
		}

		// If it is not a mention command
		if command == "" {
			for _, prefix := range baseRouter.PrefixSettings.Prefixes {
				content := msg.Content
				if baseRouter.PrefixSettings.IgnoreCase {
					content = strings.ToUpper(content)
					prefix = strings.ToUpper(prefix)
				}

				if strings.HasPrefix(content, prefix) {
					command = strings.Trim(string(contentRunes[len([]rune(prefix)):]), " ")
					break
				}
			}
		}

		args := ParseArguments(command)
		if len(args) > 0 {
			if cmd, exists := baseRouter.Commands[args[0].Raw()]; exists {
				if cmd != nil {
					cmd.execute(msg, args...)
				}
			}
		}
	})
	baseRouter.RemoveHandler = remove
	return baseRouter, nil
}

func (r *Router) GetMentions() []string {
	return []string{r.Session.State.User.Mention(), fmt.Sprintf("<@!%s>", r.Session.State.User.ID)}
}

func (r *Router) AddCommand(command *Command) *Router {
	AddCommandToMap(command, r.Commands)
	return r
}
