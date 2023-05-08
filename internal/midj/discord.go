package midj

import (
	"aiStudio/internal/conf"
	dis "aiStudio/internal/midj/handle"
	"fmt"
	discord "github.com/bwmarrin/discordgo"
)

var discordClient *discord.Session

var (
	channleId string
	serverId  string
	userToken string
)

func Init(c *conf.Midj) (err error) {
	discordClient, err = discord.New("Bot " + c.BotToken)
	if err != nil {
		return fmt.Errorf("error creating Discord session, %s", err)
	}

	err = discordClient.Open()
	if err != nil {
		return fmt.Errorf("error opening connection, %s", err)
	}
	// defer dg.Close()

	channleId = c.ChannelID
	serverId = c.ServerID
	userToken = c.UserToken
	dis.Init(c.ChannelID)
	// 注册事件
	discordClient.AddHandler(dis.DiscordMsgCreate)
	discordClient.AddHandler(dis.DiscordMsgUpdate)

	return nil
}

func GetDiscordClient() *discord.Session {
	return discordClient
}
