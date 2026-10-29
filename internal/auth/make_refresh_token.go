package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	refToken := make([]byte, 32)
	_, err := rand.Read(refToken)
	if err != nil {
		return ""
	}
	hexStr := hex.EncodeToString(refToken)
	return hexStr
}
