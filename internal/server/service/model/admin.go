package model

import "time"

type AddUser struct {
	Token       string        `json:"token"`
	Number      int           `json:"number"`
	ExpDuration time.Duration `json:"exp_duration"`
}
