package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	username = flag.String("user", "", "email address")
	password = flag.String("pass", "", "password")
	token    = flag.String("t", "", "access token")
	summary  = flag.Bool("s", false, "don't list channels")
)

func list(d *discordgo.Session, event *discordgo.Ready) {
	gid := ""
	for {
		guilds, err := d.UserGuilds(100, "", gid)
		if err != nil {
			log.Fatal("error getting guilds:", err)
		}

		if len(guilds) == 0 {
			break
		}

		for _, i := range guilds {
			gid = i.ID

			if flag.Arg(0) != "" && gid != flag.Arg(0) {
				continue
			}

			fmt.Println(gid, i.Name)

			if *summary {
				continue
			}

			gch, err := d.GuildChannels(gid)
			if err != nil {
				log.Fatalf("error getting channels for %s: %s", gid, err)
			}

			for _, i := range gch {
				var symbol string

				switch i.Type {
				case discordgo.ChannelTypeGuildText:
					symbol = "#"
				case discordgo.ChannelTypeGuildVoice:
					symbol = "🔊 "
				default:
					symbol = "(?) "
				}

				fmt.Println("\t", i.ID, symbol+i.Name)
			}
		}
	}
	if *summary {
		goto exit
	}
	/*{
		fmt.Println("@me")
		c, err := d.UserChannels()
		if err != nil {
			log.Fatal("error getting dm channels:", err)
		}

		for _, i := range c {
			var symbol string

			switch i.Type {
			case discordgo.ChannelTypeDM:
				symbol = ""
			case discordgo.ChannelTypeGroupDM:
				symbol = "(g) "
			default:
				symbol = "(?) "
			}

			if i.Name != "" {
				fmt.Println("\t", i.ID, symbol+i.Name)
			} else {
				fmt.Printf("\t %s %s", i.ID, symbol)
				for idx, r := range i.Recipients {
					fmt.Print(r.Username)
					if idx != len(i.Recipients)-1 {
						fmt.Print(", ")
					} else {
						fmt.Print("\n")
					}
				}
			}
		}

	}*/
exit:
	os.Exit(0)
}

func main() {
	flag.Parse()

	d, err := discordgo.New(*token)
	if err != nil {
		log.Fatal("login failed:", err)
	}

	d.AddHandler(list)
	err = d.Open()
	defer d.Close()
	if err != nil {
		log.Fatal("opening the websocket connection failed:", err)
	}

	if *token == "" {
		log.Println("token:", d.Token)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
