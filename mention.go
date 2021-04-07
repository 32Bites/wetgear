package wetgear

import (
	"fmt"
	"regexp"
)

const (
	MentionUser                MentionType = 0
	MentionUserNickname        MentionType = 1
	MentionChannel             MentionType = 2
	MentionRole                MentionType = 3
	MentionCustomEmoji         MentionType = 4
	MentionCustomEmojiAnimated MentionType = 5
)

var MentionUserRegex = regexp.MustCompile(`^<@(?P<Identifier>\d+)>$`)
var MentionUserNicknameRegex = regexp.MustCompile(`^<@!(?P<Identifier>\d+})>$`)
var MentionChannelRegex = regexp.MustCompile(`^<#(?P<Identifier>\d+)>$`)
var MentionRoleRegex = regexp.MustCompile(`^<@&(?P<Identifier>\d+)>$`)
var MentionCustomEmojiRegex = regexp.MustCompile(`^<:(?P<EmojiName>\w+):(?P<Identifier>\d+)>$`)
var MentionCustomEmojiAnimatedRegex = regexp.MustCompile(`^<a:(?P<EmojiName>\w+):(?P<Identifier>\d+)>$`)

// MentionType describes mention types
type MentionType int

// Mention represents a mention
type Mention struct {
	MentionType MentionType
	ID          string
	EmojiName   string // The name of the emoji, if the mention is an emoji
}

func (m Mention) Stringify() string {
	switch m.MentionType {
	case MentionUser:
		return fmt.Sprintf("<@%s>", m.ID)
	case MentionUserNickname:
		return fmt.Sprintf("<@!%s>", m.ID)
	case MentionChannel:
		return fmt.Sprintf("<#%s>", m.ID)
	case MentionRole:
		return fmt.Sprintf("<@&%s>", m.ID)
	case MentionCustomEmoji:
		return fmt.Sprintf("<:%s:%s>", m.EmojiName, m.ID)
	case MentionCustomEmojiAnimated:
		return fmt.Sprintf("<a:%s:%s", m.EmojiName, m.ID)
	default:
		return ""
	}
}

// IsMention tests string to determine whether or not it matches any mention regular expressions
func IsMention(content string) bool {
	return GetMentionType(content) != -1
}

// GetMentionType attempts to get the MentionType of a string, upon failure returns -1 as a MentionType
func GetMentionType(content string) MentionType {
	if MentionUserRegex.MatchString(content) {
		return MentionUser
	} else if MentionUserNicknameRegex.MatchString(content) {
		return MentionUserNickname
	} else if MentionChannelRegex.MatchString(content) {
		return MentionChannel
	} else if MentionRoleRegex.MatchString(content) {
		return MentionRole
	} else if MentionCustomEmojiRegex.MatchString(content) {
		return MentionCustomEmoji
	} else if MentionCustomEmojiAnimatedRegex.MatchString(content) {
		return MentionCustomEmojiAnimated
	} else {
		return -1
	}
}

// GetMention tries to make a Mention out of a string. Upon failure, return nil.
func GetMention(content string) *Mention {
	var regex *regexp.Regexp

	if MentionUserRegex.MatchString(content) {
		regex = MentionUserRegex
	} else if MentionUserNicknameRegex.MatchString(content) {
		regex = MentionUserNicknameRegex
	} else if MentionChannelRegex.MatchString(content) {
		regex = MentionChannelRegex
	} else if MentionRoleRegex.MatchString(content) {
		regex = MentionRoleRegex
	} else if MentionCustomEmojiRegex.MatchString(content) {
		regex = MentionCustomEmojiRegex
	} else if MentionCustomEmojiAnimatedRegex.MatchString(content) {
		regex = MentionCustomEmojiAnimatedRegex
	} else {
		return nil
	}
	id := ""
	emojiName := ""

	matches := regex.FindStringSubmatch(content)
	names := regex.SubexpNames()

	for i, match := range matches {
		if i != 0 {
			if names[i] == "Identifier" {
				id = match
			} else if names[i] == "EmojiName" {
				emojiName = match
			}
		}
	}

	return &Mention{
		MentionType: GetMentionType(content),
		ID:          id,
		EmojiName:   emojiName,
	}
}
