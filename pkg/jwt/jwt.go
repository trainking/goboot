package jwt

import (
	"encoding/json"
	"errors"
	"time"

	jwtV4 "github.com/golang-jwt/jwt/v4"
)

type (
	jwtCustomClaims struct {
		jwtV4.RegisteredClaims
		Data interface{} `json:"data"`
	}

	JwtToken struct {
		secret string
	}
)

// NewJwtToken is create jwt tonen generate.
func NewJwtToken(secret string) *JwtToken {
	return &JwtToken{secret: secret}
}

// Generate 创建一个jwt
func (t *JwtToken) Generate(data interface{}, expire time.Duration) (string, error) {
	j := &jwtCustomClaims{
		Data: data,
	}
	j.RegisteredClaims = jwtV4.RegisteredClaims{
		ExpiresAt: jwtV4.NewNumericDate(time.Now().Add(expire)),
	}

	jwtT := jwtV4.NewWithClaims(jwtV4.SigningMethodHS256, j)
	tokenString, err := jwtT.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 验证一个jwt
func (t *JwtToken) ParseToken(tokenString string, data interface{}) error {
	token, err := jwtV4.Parse(tokenString, func(token *jwtV4.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtV4.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(t.secret), nil
	})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwtV4.MapClaims); ok && token.Valid {
		if data != nil {
			dataB, _ := json.Marshal(claims["data"])
			if err := json.Unmarshal(dataB, data); err != nil {
				return err
			}
		}

		return nil
	}

	return errors.New("invilid")
}
