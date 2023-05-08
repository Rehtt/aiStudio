package model

import "time"

type ExternalInfo struct {
	Key                string
	Number             int
	ExpirationDuration time.Duration
}
