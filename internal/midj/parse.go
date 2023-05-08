package midj

import (
	discord "github.com/bwmarrin/discordgo"
	"path/filepath"
	"strings"
)

func ParseMsgIdAndHash(m *discord.Message) (msgId string, msgHash string) {
	msgId = m.ID
	for _, att := range m.Attachments {
		name := att.Filename
		name = strings.ReplaceAll(name, filepath.Ext(name), "")
		s := strings.Split(name, "_")
		if len(s) == 0 {
			continue
		}
		msgHash = s[len(s)-1]
		break
	}
	return
}
