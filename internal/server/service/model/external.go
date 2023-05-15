package model

import "time"

type ExternalInfo struct {
	Token              string        `json:"token"`
	RedisKey           string        `json:"-"`
	LockKey            string        `json:"-"`
	Number             int           `json:"number"`
	ExpirationDuration time.Duration `json:"expiration_duration"`
}

type GenerateType string

const (
	GEN = GenerateType("generate")
	UP  = GenerateType("upscale")
	VA  = GenerateType("variate")
)

type GenerateImageRequest struct {
	Type   GenerateType `json:"type"`   // 类型
	Prompt string       `json:"prompt"` // generate
	GenId  string       `json:"gen_id"` // 消息id
	V      int          `json:"v"`      // Variate
	U      int          `json:"u"`      // upscale
}
type GenerateImageResponse struct {
	Id string `json:"id"`
}
type ProgressResponse struct {
	Progress int    `json:"progress"`
	Prompt   string `json:"prompt"`
	Status   string `json:"status"`
	ImageUrl string `json:"image_url"`
}
