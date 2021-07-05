package wetgear

import (
	"github.com/bwmarrin/discordgo"
)

const (
	EmojiFastBackward  = "\u23EA"
	EmojiBackwardArrow = "\u2B05\uFE0F"
	EmojiForwardArrow  = "\u27A1\uFE0F"
	EmojiFastForward   = "\u23E9"
)

type PaginationPage func() *discordgo.MessageEmbed

type Pagination struct {
	Embeds    []PaginationPage
	RepliesTo *discordgo.MessageReference
	Router    *Router
	Done      bool
	Current   uint
	ChannelID string
	MessageID string
	AuthorID  string
}

func (r *Router) CreatePagination(channelID, authorID string) *Pagination {
	return &Pagination{ChannelID: channelID, Router: r, Done: false, AuthorID: authorID}
}

func (p *Pagination) AddEmbeds(embeds ...PaginationPage) *Pagination {
	p.Embeds = embeds
	return p
}

func (p *Pagination) SetRepliesTo(repliesTo *discordgo.MessageReference) *Pagination {
	p.RepliesTo = repliesTo
	return p
}

func (p *Pagination) addReactions() error {
	reactions := []string{EmojiFastBackward, EmojiBackwardArrow, EmojiForwardArrow, EmojiFastForward}

	for _, reaction := range reactions {
		err := p.Router.Session.MessageReactionAdd(p.ChannelID, p.MessageID, reaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pagination) Update(index uint) error {
	p.Current = index
	_, err := p.Router.Session.ChannelMessageEditEmbed(p.ChannelID, p.MessageID, p.Embeds[index]())
	return err
}

func (p *Pagination) Spawn() error {
	defer func() {
		p.Done = true
	}()
	if len(p.Embeds) > 0 && p.Embeds[0] != nil {
		var err error
		var message *discordgo.Message
		if p.RepliesTo != nil {
			message, err = ChannelMessageSendEmbedReply(p.Router.Session, p.ChannelID, p.Embeds[0](), p.RepliesTo)
		} else {
			message, err = p.Router.Session.ChannelMessageSendEmbed(p.ChannelID, p.Embeds[0]())
		}

		if err != nil {
			return err
		}

		p.MessageID = message.ID

		p.Router.Paginations.Lock()
		defer p.Router.Paginations.Unlock()
		p.Router.Paginations.Values[message.ID] = p

		err = p.addReactions()
		if err != nil {
			return err
		}
	}

	return nil
}
