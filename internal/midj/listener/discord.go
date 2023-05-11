package listener

import (
	"aiStudio/internal/conf"
	"fmt"
	discord "github.com/bwmarrin/discordgo"
)

var discordClient *discord.Session

var (
	cg = make(map[string]struct{})
)

func Init(c []conf.Midj) (err error) {
	for _, cc := range c {
		discordClient, err = discord.New("Bot " + cc.BotToken)
		if err != nil {
			return fmt.Errorf("error creating Discord session, %s", err)
		}

		err = discordClient.Open()
		if err != nil {
			return fmt.Errorf("error opening connection, %s", err)
		}
		// defer dg.Close()

		for _, id := range cc.ChannelID {
			cg[cc.GuildID+"-"+id] = struct{}{}
		}
		// 注册事件
		discordClient.AddHandler(DiscordMsgCreate)
		discordClient.AddHandler(DiscordMsgUpdate)
	}
	return nil
}
