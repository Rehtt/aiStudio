package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/Rehtt/Kit/random"
	jsoniter "github.com/json-iterator/go"
	"time"
)

func GenId(v ...any) string {
	s := sha256.New()
	for _, vv := range v {
		jsoniter.NewEncoder(s).Encode(vv)
	}
	s.Write([]byte(random.RandName()))
	s.Write([]byte(time.Now().String()))
	return hex.EncodeToString(s.Sum(nil))
}
