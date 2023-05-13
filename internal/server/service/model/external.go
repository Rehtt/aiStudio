package model

import "time"

type ExternalInfo struct {
	Token              string        `json:"token"`
	RedisKey           string        `json:"-"`
	LockKey            string        `json:"-"`
	Number             int           `json:"number"`
	ExpirationDuration time.Duration `json:"expiration_duration"`
}

type GenerateImageRequest struct {
	Prompt string `json:"prompt"`
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
