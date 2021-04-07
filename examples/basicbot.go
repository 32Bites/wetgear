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
		BotsAllowed: true,
		Commands:    map[string]*wetgear.Command{},
	}

	router, err = wetgear.NewRouter(sess, router)
	if err != nil {
		panic(err)
	}

	ping := wetgear.NewCommand(router).AddAliases("ping").SetExecutor(func(ctx wetgear.Context) {
		ctx.GetSession().ChannelMessageSend(ctx.MessageCreate.ChannelID, "Pong!")
	})
	pong := wetgear.NewCommand(router).AddAliases("pong").SetExecutor(func (ctx wetgear.Context) {
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
