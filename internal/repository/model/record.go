package model

import "gorm.io/gorm"

type RecordTable struct {
	*gorm.Model
	Token  string `json:"token" gorm:"index"`
	Prompt string `json:"prompt"`
	GenID  string `json:"gen_id" gorm:"index"`

	Progress    int     `json:"progress"`
	Status      string  `json:"status"`
	GuildID     string  `json:"guild_id"`
	ChannelID   string  `json:"channel_id"`
	MsgID       string  `json:"msg_id" gorm:"index"`
	ParentMsgID string  `json:"parent_id" gorm:"index"`
	Option      string  `json:"option"`
	ImageUrl    *string `json:"image_url"`
	MHash       string  `json:"m_hash"`
	Remarks     string  `json:"remarks"`
}
