package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserId int64     `json:"userId"`
	Email  string    `json:"email,omitempty"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret            string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewJWTService(secret string, accessExp, refreshExp time.Duration) *JWTService {
	return &JWTService{
		secret:            secret,
		accessExpiration:  accessExp,
		refreshExpiration: refreshExp,
	}
}

// GenerateAccessToken создает access токен для пользователя
func (j *JWTService) GenerateAccessToken(userId int64, email string) (string, error) {
	claims := Claims{
		UserId: userId,
		Email:  email,
		Type:   AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// GenerateRefreshToken создает refresh токен для пользователя
func (j *JWTService) GenerateRefreshToken(userId int64) (string, error) {
	claims := Claims{
		UserId: userId,
		Type:   RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// ValidateToken проверяет валидность токена и возвращает claims
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется правильный алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ExtractUserId извлекает userId из токена
func (j *JWTService) ExtractUserId(tokenString string) (int64, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserId, nil
}

// ExtractEmail извлекает email из токена
func (j *JWTService) ExtractEmail(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Email, nil
}
