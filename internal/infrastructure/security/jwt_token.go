package security

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"consumo-real-server/internal/shared/auth"
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: ttl}
}

type jwtClaims struct {
	jwt.RegisteredClaims
	EmpresaID int64  `json:"empresa_id"`
	Papel     string `json:"papel"`
}

func (m *JWTManager) Gerar(claims auth.Claims) (string, error) {
	now := time.Now()
	parsed := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(claims.UsuarioID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
		EmpresaID: claims.EmpresaID,
		Papel:     claims.Papel,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, parsed).SignedString(m.secret)
}

func (m *JWTManager) Validar(tokenStr string) (auth.Claims, error) {
	parsed := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, parsed, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return auth.Claims{}, auth.ErrTokenInvalido
	}

	usuarioID, err := strconv.ParseInt(parsed.Subject, 10, 64)
	if err != nil || usuarioID <= 0 {
		return auth.Claims{}, errors.New("subject inválido no token")
	}

	return auth.Claims{
		UsuarioID: usuarioID,
		EmpresaID: parsed.EmpresaID,
		Papel:     parsed.Papel,
	}, nil
}
