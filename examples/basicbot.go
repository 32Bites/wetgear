package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/32bites/wetgear"
	"github.com/bwmarrin/discordgo"
)

func main() {
	sess, err := discordgo.New("Bot " + os.Getenv("WETGEAR_TEST_TOKEN"))
	if err != nil {
		panic(err)
	}
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	router := &wetgear.Router{
		PrefixSettings: wetgear.PrefixSettings{
			Prefixes:   []string{"!"},
			IgnoreCase: true,
			HandlePing: true,
		},
		HelpSettings: wetgear.HelpSettings{
			Enabled: true,
		}, // Will use Default aliases.
		BotsAllowed: true,
	}

	router, err = wetgear.NewRouter(sess, router)
	if err != nil {
		panic(err)
	}

	ping := wetgear.
		NewCommand(router).
		AddAliases("ping").
		SetName("Ping").
		SetDescription("Creates some embeds").
		SetExecutor(func(ctx wetgear.Context) {
			fmt.Println(ctx.MessageCreate.Author.ID)
			ctx.Command.Router.CreatePagination(ctx.MessageCreate.ChannelID, ctx.MessageCreate.Author.ID).AddEmbeds(
				wetgear.EmbedToPage(&discordgo.MessageEmbed{Description: "Page 1"}),
				wetgear.EmbedToPage(&discordgo.MessageEmbed{Description: "Page 2"}),
				wetgear.EmbedToPage(&discordgo.MessageEmbed{Description: "Page 3"}),
			).Spawn()
		})
	pong := wetgear.NewCommand(router).AddAliases("pong").SetExecutor(func(ctx wetgear.Context) {
		ctx.GetSession().ChannelMessageSend(ctx.MessageCreate.ChannelID, ":ping_pong:")
	})
	ping.AddSubCommand(pong)
	router.AddCommand(ping)

	// now !ping and !ping pong are valid commands

	err = sess.Open()
	if err != nil {
		panic(err)
	}
	fmt.Println("Running Bot as", sess.State.User.String())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	sess.Close()
}
