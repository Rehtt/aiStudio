package model

import "time"

type ExternalInfo struct {
	Token              string
	RedisKey           string
	LockKey            string
	Number             int
	ExpirationDuration time.Duration
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
