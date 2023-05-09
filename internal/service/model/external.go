package model

import "time"

type ExternalInfo struct {
	Key                string
	Number             int
	ExpirationDuration time.Duration
	LockKey            string
}

type GenerateImageRequest struct {
	Prompt string `json:"prompt"`
}
type GenerateImageResponse struct {
	Id string `json:"id"`
}
