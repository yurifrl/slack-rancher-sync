// This package formats things to be printed out, it does not print just formarts
package apis

import (
	"net/url"

	"gopkg.in/urfave/cli.v2"

	"github.com/cloud104/slack-rancher-sync/pkg/apis/rancher"
	"github.com/cloud104/slack-rancher-sync/pkg/apis/slack"
)

type Apis struct {
	Slack   slack.Config
	Rancher rancher.Config
}

// ListResponse returns the users to be returnerd
type ListResponse struct {
	Email string `json:"email"`
}

// Creates Reconciler based on cli inputs
func NewCliApiRenconciler(c *cli.Context) Apis {
	return Apis{
		Slack: slack.Config{
			Token:         c.String("slack.token"),
			ChannelPrefix: c.String("slack.channel.prefix"),
		},
		Rancher: rancher.Config{
			Endpoint: c.String("rancher.endpoint"),
			User:     c.String("rancher.user"),
			Token:    c.String("rancher.token"),
		},
	}
}

func (c *Apis) ListSlackUsers() (resp []ListResponse, err error) {
	// Call Apis
	state := c.Slack.GetState()
	// loop To get Values
	for _, root := range state {
		for _, user := range root.Users {
			resp = append(resp, ListResponse{Email: user.Profile.Email})
		}
	}
	return resp, err
}

func (c *Apis) Reconcile() (resp []string, err error) {
	state := c.Slack.GetState()
	for _, root := range state {
		for _, user := range root.Users {
			resp = append(resp, url.QueryEscape("adfs_user://"+user.Profile.Email))
		}
	}
	return resp, err
}
