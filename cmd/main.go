package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloud104/slack-rancher-sync/pkg/apis"
	cli "gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"
)

var (
	version = "0.0.0"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
			Value:   "secrets/config.yaml",
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "slack.token",
			Aliases: []string{"slack-token"},
			Usage:   "Slack Token",
			EnvVars: []string{"SLACK_TOKEN"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "slack.channel.prefix",
			Aliases: []string{"slack-channel-prefix"},
			Usage:   "Channel Prefix",
			EnvVars: []string{"SLACK_CHANNEL_PREFIX"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "rancher.endpoint",
			Aliases: []string{"rancher-endpoint"},
			Usage:   "Rancher Endpoint",
			EnvVars: []string{"RANCHER_ENDPOINT"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "rancher.user",
			Aliases: []string{"rancher-user"},
			Usage:   "Rancher User",
			EnvVars: []string{"RANCHER_User"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "rancher.token",
			Aliases: []string{"rancher-token"},
			Usage:   "Rancher Token",
			EnvVars: []string{"RANCHER_Token"},
		}),
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "list-slack-users",
				Action: listSlackUsers,
			},
			{
				Name:   "reconcile",
				Action: reconcile,
			},
		},
		Action: reconcile,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Flags:  flags,
	}

	app.Run(os.Args)
}

// Commands
func listSlackUsers(c *cli.Context) (err error) {
	api := apis.NewCliApiRenconciler(c)
	resp, err := api.ListSlackUsers()
	if err != nil {
		return err
	}
	if j, err := json.Marshal(resp); err == nil {
		fmt.Println(string(j))
	}

	return err
}

func reconcile(c *cli.Context) (err error) {
	api := apis.NewCliApiRenconciler(c)
	resp, err := api.Reconcile()
	if err != nil {
		return err
	}
	if j, err := json.Marshal(resp); err == nil {
		fmt.Println(string(j))
	}

	return err
}
