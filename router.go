package wetgear

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var Logger *log.Logger = log.New(os.Stdout, "[WETGEAR] ", log.LstdFlags)

// PrefixSettings defines the settings for a router's prefixes
type PrefixSettings struct {
	Prefixes   []string
	IgnoreCase bool
	HandlePing bool // whether or not to treat a ping as a command prefix. Example: "@BotName help" does the same thing as "!help" if "!" is a prefix.
}

// HelpSettings defines the settings for the router's built-in help command.
type HelpSettings struct {
	Enabled          bool
	NotFoundExecutor CommandNotFoundExecutor
	FoundExecutor    CommandFoundExecutor
	Aliases          []string
}

// Router is a command router
type Router struct {
	Session          *discordgo.Session
	Commands         map[string]*Command
	PrefixSettings   PrefixSettings
	BotsAllowed      bool
	GlobalMiddlwares []CommandMiddleware
	RemoveHandlers   []func()
	HelpSettings     HelpSettings
	Paginations      struct {
		Values map[string]*Pagination
		sync.RWMutex
	}
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
	baseRouter.Commands = map[string]*Command{}
	baseRouter.Paginations = struct {
		Values map[string]*Pagination
		sync.RWMutex
	}{
		Values: map[string]*Pagination{},
	}

	baseRouter.RemoveHandlers = []func(){session.AddHandler(baseRouter.messageCreateHandler)}
	baseRouter.RemoveHandlers = append(baseRouter.RemoveHandlers, session.AddHandler(func(session *discordgo.Session, event *discordgo.MessageReactionAdd) {
		baseRouter.reactionHandler(event.Emoji.Name, event.MessageID, event.UserID)
	}), session.AddHandler(func(session *discordgo.Session, event *discordgo.MessageReactionRemove) {
		baseRouter.reactionHandler(event.Emoji.Name, event.MessageID, event.UserID)
	}))

	if baseRouter.HelpSettings.Enabled {
		if len(baseRouter.HelpSettings.Aliases) == 0 {
			baseRouter.HelpSettings.Aliases = []string{"help", "h"}
		}
		helpCommand := NewCommand(baseRouter).AddAliases(baseRouter.HelpSettings.Aliases...).
			SetName("Help").
			SetDescription("Provides Help Information for a command.").
			SetExecutor(helpExecutor)
		baseRouter.AddCommand(helpCommand)
	}

	return baseRouter, nil
}

func (r *Router) GetMentions() []string {
	return []string{r.Session.State.User.Mention(), fmt.Sprintf("<@!%s>", r.Session.State.User.ID)}
}

func (r *Router) AddCommand(command *Command) *Router {
	addCommandToMap(command, r.Commands)
	return r
}

func (r *Router) reactionHandler(emoji, messageID, userID string) {
	r.Paginations.Lock()
	defer r.Paginations.Unlock()
	if pagination, exists := r.Paginations.Values[messageID]; exists {
		var err error

		if !pagination.Done || pagination.AuthorID != userID {
			return
		}

		switch emoji {
		case EmojiFastBackward:
			err = pagination.Update(0)
		case EmojiFastForward:
			err = pagination.Update(uint(len(pagination.Embeds) - 1))
		case EmojiBackwardArrow:
			newIndex := int(pagination.Current) - 1
			if newIndex <= len(pagination.Embeds)-1 && newIndex >= 0 {
				pagination.Update(uint(newIndex))
			}
		case EmojiForwardArrow:
			newIndex := int(pagination.Current) + 1
			if newIndex <= len(pagination.Embeds)-1 {
				pagination.Update(uint(newIndex))
			}
		}

		if err != nil {
			Logger.Println("An error occurred in reactionHandler:", err.Error())
		}
	}
}

func (r *Router) messageCreateHandler(session *discordgo.Session, event *discordgo.MessageCreate) {
	// Check if bot and if bots are allowed, and make sure that the message is not empty
	if (event.Author.Bot && !r.BotsAllowed) || event.Content == "" {
		return
	}
	contentRunes := []rune(event.Content)
	command := ""

	// Check if the message starts with a mention, and that the message is longer than the mention
	if r.PrefixSettings.HandlePing {
		for _, mention := range r.GetMentions() {
			if strings.HasPrefix(event.Content, mention+" ") && len(contentRunes) > len([]rune(mention+" ")) {
				command = strings.TrimSpace(string(contentRunes[len([]rune(mention)):]))
				break
			}
		}
	}

	// If it is not a mention command
	if command == "" {
		for _, prefix := range r.PrefixSettings.Prefixes {
			content := event.Content
			if r.PrefixSettings.IgnoreCase {
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
		if cmd, exists := r.Commands[args[0].Raw()]; exists {
			if cmd != nil {
				cmd.execute(event, args...)
			}
		}
	}
}
