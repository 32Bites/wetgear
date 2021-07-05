package wetgear

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

const (
	ColorBlurple         = 0x7289DA
	ColorFullWhite       = 0xFFFFFF
	ColorDarkButNotBlack = 0x2C2F33
	ColorNotQuiteBlack   = 0x23272A
	ColorCrimson         = 0xDC143C
	ColorDarkSalmon      = 0xE9967A
	ColorLightCoral      = 0xF08080
	ColorPink            = 0xFFC0CB
	ColorHotPink         = 0xFF69B4
	ColorDeepPink        = 0xFF1493
	ColorTomato          = 0xFF6347
	ColorGold            = 0xFFD700
	ColorLemonChiffon    = 0xFFFACD
	ColorPapayawhip      = 0xFFEFD5
	ColorMoccasin        = 0xFFE4B5
	ColorDarkKhaki       = 0xBDB76B
	ColorViolet          = 0xEE82EE
	ColorMediumOrchid    = 0xBA55D3
	ColorIndigo          = 0x4B0082
	ColorSlateBlue       = 0x6A5ACD
	ColorChartreuse      = 0x7FFF00
	ColorSeaGreen        = 0x2E8B57
	ColorAqua            = 0x00FFFF
	ColorSteelBlue       = 0x4682B4
	ColorMaroon          = 0x800000
	ColorTeal            = 0x008080
	ColorMidnightBlue    = 0x191970
	ColorChocolate       = 0xD2691E
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

func EmbedToPage(embed *discordgo.MessageEmbed) PaginationPage {
	return func() *discordgo.MessageEmbed {
		return embed
	}
}

func CombinePermissions(perms ...int) int {
	perm := 0
	for _, prm := range perms {
		perm = perm | prm
	}
	return perm
}
