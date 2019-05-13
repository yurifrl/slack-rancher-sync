package slack

import (
	"strings"

	"github.com/nlopes/slack"
)

type Config struct {
	Token         string
	ChannelPrefix string
}

type ChannelWithUser struct {
	Channel *slack.Channel
	Users   []*slack.User
}

func (c *Config) GetState() (buffer map[string]ChannelWithUser) {
	buffer = make(map[string]ChannelWithUser)
	api := slack.New(c.Token)
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

			buffer[conversation.ID] = ChannelWithUser{
				Channel: &conversation,
				Users:   foundUsers,
			}
		}
	}

	return buffer
}
