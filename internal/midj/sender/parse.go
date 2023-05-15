package sender

import (
	discord "github.com/bwmarrin/discordgo"
	"path/filepath"
	"strings"
)

func ParseMsgIdAndHash(m *discord.Message) (msgId, msgHash string, imageUrl *string) {
	msgId = m.ID
	for _, att := range m.Attachments {
		name := att.Filename
		name = strings.ReplaceAll(name, filepath.Ext(name), "")
		s := strings.Split(name, "_")
		if len(s) == 0 {
			continue
		}
		msgHash = s[len(s)-1]
		imageUrl = &att.URL
		break
	}
	return
}
