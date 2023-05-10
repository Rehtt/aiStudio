package model

import "gorm.io/gorm"

type RecordTable struct {
	gorm.Model
	Token  string `json:"token" gorm:"index"`
	Prompt string `json:"prompt"`
	GenID  string `json:"gen_id" gorm:"index"`
}
