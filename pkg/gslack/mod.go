package gslack

import (
	"github.com/slack-go/slack"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

type Client struct {
	*slack.Client
	ChannelID string
}

func NewClient(slackToken string, channelID string) Client {

	api := slack.New(slackToken)

	//channelID := "C1234567890"
	//_, _, err := api.PostMessage(
	//	channelID,
	//	slack.MsgOptionText("Hello from Go 👋", false),
	//)
	//if err != nil {
	//	Log.Fatal(err)
	//}

	return Client{
		Client:    api,
		ChannelID: channelID,
	}
}

func (c *Client) Text(text string) error {
	_, _, err := c.PostMessage(
		c.ChannelID,
		slack.MsgOptionText(text, false),
	)

	if err != nil {
		Log.Warnf("PostMessage ChannelID=%s, Error=%s", c.ChannelID, err.Error())
	}

	return err
}

//📩 消息
//api.PostMessage()
//api.UpdateMessage()
//api.DeleteMessage()

//👥 用户 & Channel
//api.GetUsers()
//api.GetConversationInfo()
//api.GetConversations()

//🔘 交互组件
//slack.NewButtonBlockElement()
//slack.NewSectionBlock()
//slack.NewActionBlock()

//⚡ Slash Command
//slack.SlashCommandParse(r)
