package generate

import (
	"crypto/rand"
	"fmt"
)

func GenerateCookieToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
