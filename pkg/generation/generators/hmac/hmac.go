package hmac

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// Hmac generates a hmac value for webhooks
func Hmac() (string, error) {
	value, err := randStringBytes(41)
	if err != nil {
		return value, fmt.Errorf("generating hmac: %w", err)
	}
	return value, nil
}

func randStringBytes(n int) (string, error) {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, (n+1)/2) // can be simplified to n/2 if n is always even

	if _, err := src.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b)[:n], nil
}
