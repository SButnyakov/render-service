package tokenManager

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type Manager struct {
	signingKey string
}

func New(signingKey string) (*Manager, error) {
	fn := "lib.api.tokenManager.New"

	if signingKey == "" {
		return nil, fmt.Errorf("%s: empty signing key", fn)
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(uid int64, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(duration).Unix(),
		Subject:   strconv.Itoa(int(uid)),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) NewRefreshToken() (string, error) {
	return "", nil
}
