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
		sendMsg(m.Message, FirstTrigger)
		return
	}

	for _, attachment := range m.Attachments {
		if attachment.Width > 0 && attachment.Height > 0 {
			sendMsg(m.Message, GenerateEnd)
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
		sendMsg(m.Message, GenerateEditError)
		return
	}

	if len(m.Embeds) > 0 {
		send(m.Embeds)
	} else {
		sendMsg(m.Message, UpdateMsg)
	}
}

type ReqCb struct {
	Embeds     []*discord.MessageEmbed `json:"embeds,omitempty"`
	DiscordMsg *discord.Message        `json:"discord_msg,omitempty"`
	Type       Scene                   `json:"type"`
}

func send(embeds []*discord.MessageEmbed) {
	body := ReqCb{
		Embeds: embeds,
		Type:   RichText,
	}
	request(body)
}
func sendMsg(m *discord.Message, t Scene) {
	request(
		ReqCb{
			DiscordMsg: m,
			Type:       t,
		},
	)
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
