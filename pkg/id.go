package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/Rehtt/Kit/random"
	"time"
)

func GenId(v any) string {
	s := sha256.New()
	json.NewEncoder(s).Encode(v)
	s.Write([]byte(random.RandName()))
	s.Write([]byte(time.Now().String()))
	return hex.EncodeToString(s.Sum(nil))
}
