package auth

import (
	"fmt"
	"testing"
	"time"
)

func TestAuth_Sign(t *testing.T) {
	token := SignJwtToken("api", "token", &JwtOptions{
		ExpiresAt: time.Now().Add(20 * time.Minute).Unix(),
	})

	fmt.Println(token)

	_, _ = VerifyJwtToken(token, "token")
}
