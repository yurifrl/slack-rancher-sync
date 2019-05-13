package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/nlopes/slack"
	cli "gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

var (
	version = "0.0.0"
)

func main() {
	app := cli.NewApp()
	app.Action = run
	app.Version = version
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
			Value: "secrets/config.yaml",
		},
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:   "slack.token",
				Usage:  "Slack Token",
				EnvVar: "SLACK_TOKEN",
			},
		),
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:   "slack.channel.prefix",
				Usage:  "Channel Prefix",
				EnvVar: "SLACK_CHANNEL_PREFIX",
				Value:  "test-",
			},
		),
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:   "rancher.endpoint",
				Usage:  "Rancher Endpoint",
				EnvVar: "RANCHER_ENDPOINT",
			},
		),
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:   "rancher.user",
				Usage:  "Rancher User",
				EnvVar: "RANCHER_User",
			},
		),
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:   "rancher.token",
				Usage:  "Rancher Token",
				EnvVar: "RANCHER_Token",
			},
		),
	}
	app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))
	app.Flags = flags

	pp.Println("Starting")
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err, "cli boot failed")
	}

}
func run(c *cli.Context) {
	x := Sync{
		SlackToken:      c.String("slack.token"),
		ChannelPrefix:   c.String("slack.channel.prefix"),
		RancherEndpoint: c.String("rancher.endpoint"),
		RancherUser:     c.String("rancher.user"),
		RancherToken:    c.String("rancher.token"),
	}
	pp.Println(x)

	x.Exec()
}

type Sync struct {
	SlackToken      string
	ChannelPrefix   string
	RancherEndpoint string
	RancherUser     string
	RancherToken    string
}

type SlackChannelWithUser struct {
	Channel *slack.Channel
	Users   []*slack.User
}

func (c *Sync) Exec() {
	slackState := GetSlackState(c)

	for _, root := range slackState {
		for _, user := range root.Users {
			id := url.QueryEscape("adfs_user://" + user.Profile.Email)
			pp.Println(id)
		}
	}
}

func GetSlackState(c *Sync) (buffer map[string]*SlackChannelWithUser) {
	buffer = make(map[string]*SlackChannelWithUser)
	api := slack.New(c.SlackToken)
	conversations, _, _ := api.GetConversations(&slack.GetConversationsParameters{
		ExcludeArchived: "true",
		Types:           []string{"private_channel"},
	})
	for _, conversation := range conversations {
		if strings.HasPrefix(conversation.Name, c.ChannelPrefix) {
			members, _, _ := api.GetUsersInConversation(&slack.GetUsersInConversationParameters{
				ChannelID: conversation.ID,
			})

			var foundUsers []*slack.User
			for _, member := range members {
				user, _ := api.GetUserInfo(member)
				foundUsers = append(foundUsers, user)
			}

			buffer[conversation.ID] = &SlackChannelWithUser{
				Channel: &conversation,
				Users:   foundUsers,
			}
		}
	}

	return buffer
}
