package listener

import (
	"strings"

	discord "github.com/bwmarrin/discordgo"
)

type Scene string

const (
	/**
	 * 首次触发生成
	 */
	FirstTrigger Scene = "FirstTrigger"
	/**
	 * 生成图片结束
	 */
	GenerateEnd Scene = "GenerateEnd"
	/**
	 * 发送的指令midjourney生成过程中发现错误
	 */
	GenerateEditError Scene = "GenerateEditError"
	/**
	 * 进度消息
	 */
	UpdateMsg Scene = "update"

	/**
	 * 富文本
	 */
	RichText Scene = "RichText"
	/**
	 * 发送的指令midjourney直接报错或排队阻塞不在该项目中处理 在业务服务中处理
	 * 例如：首次触发生成多少秒后没有回调业务服务判定会指令错误或者排队阻塞
	 */
)

func DiscordMsgCreate(s *discord.Session, m *discord.MessageCreate) {

	if _, ok := cg[m.GuildID+"-"+m.ChannelID]; !ok {
		return
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	// 重新生成不发送
	// TODO 优化 使用 From
	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		createMsg(m, FirstTrigger)
		return
	}

	for _, attachment := range m.Attachments {
		if attachment.Width > 0 && attachment.Height > 0 {
			replay(m)
			return
		}
	}
}

func DiscordMsgUpdate(s *discord.Session, m *discord.MessageUpdate) {
	if _, ok := cg[m.GuildID+"-"+m.ChannelID]; !ok {
		return
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "(Stopped)") {
		updateMsg(m, GenerateEditError)
		return
	}

	if len(m.Embeds) > 0 {
		send(m.Embeds)
	} else {
		updateMsg(m, UpdateMsg)
	}
}

type ReqCb struct {
	Embeds        []*discord.MessageEmbed `json:"embeds,omitempty"`
	DiscordCreate *discord.MessageCreate  `json:"discord_create,omitempty"`
	DiscordUpdate *discord.MessageUpdate  `json:"discord_update,omitempty"`
	Type          Scene                   `json:"type"`
}

func replay(m *discord.MessageCreate) {
	body := ReqCb{
		DiscordCreate: m,
		Type:          GenerateEnd,
	}
	request(body)
}

func send(embeds []*discord.MessageEmbed) {
	body := ReqCb{
		Embeds: embeds,
		Type:   RichText,
	}
	request(body)
}

func updateMsg(m *discord.MessageUpdate, t Scene) {
	body := ReqCb{
		DiscordUpdate: m,
		Type:          t,
	}
	request(body)
}

func createMsg(m *discord.MessageCreate, t Scene) {
	request(ReqCb{
		DiscordCreate: m,
		Type:          t,
	})
}

//func request(r ReqCb) {
//	req(&r)
//	//data, err := json.Marshal(params)
//	//if err != nil {
//	//	fmt.Println("json marshal error: ", err)
//	//	return
//	//}
//	//os.WriteFile("test.json", data, 0644)
//	// req, err := http.NewRequest("POST", initialization.GetConfig().CB_URL, strings.NewReader(string(data)))
//	// if err != nil {
//	// 	fmt.Println("http request error: ", err)
//	// 	return
//	// }
//	// req.Header.Set("Content-Type", "application/json")
//	// client := &http.Client{}
//	// resp, err := client.Do(req)
//	// if err != nil {
//	// 	fmt.Println("http request error: ", err)
//	// 	return
//	// }
//	// defer resp.Body.Close()
//}
//
//var req func(body *ReqCb)
//
//func Request(f func(body *ReqCb)) {
//	req = f
//}
