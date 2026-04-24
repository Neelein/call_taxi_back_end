package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func CreateSignedID() string {
	sessionKey := []byte(os.Getenv("SESSIONKEY"))
	uuid := uuid.New().String()
	h := hmac.New(sha256.New, sessionKey)
	h.Write([]byte(uuid))
	signature := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s:%s", uuid, signature)
}
