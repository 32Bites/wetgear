package wetgear

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func AddCommandToMap(command *Command, commandMap map[string]*Command) {
	for _, alias := range command.Aliases {
		if _, exists := commandMap[alias]; exists {
			continue
		}
		commandMap[alias] = command
	}
}

func ChannelMessageSendEmbedReply(session *discordgo.Session, channelID string, embed *discordgo.MessageEmbed, reference *discordgo.MessageReference) (*discordgo.Message, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	return session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Reference: reference,
		Embed:     embed,
	})
}
