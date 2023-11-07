package tokenManager

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

var (
	ErrInvalidToken = fmt.Errorf("token is invalid")
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

func (m *Manager) NewJWT(uid int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		Subject:   strconv.Itoa(int(uid)),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) NewRT(uid int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Subject:   strconv.Itoa(int(uid)),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(token, secret string) (*jwt.StandardClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
