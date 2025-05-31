package auth

import (
	"LostAndFound/internal/domain/repository"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	CONFIG_AUTH_PATH = "CONFIG_AUTH_PATH"
)

type TokenManager struct {
	SecretKey string        `yaml:"secret_key"`
	TokenTTL  time.Duration `yaml:"token_ttl"`
	CacheRepo repository.CacheRepo
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenManager(cacheRepo repository.CacheRepo) (*TokenManager, error) {
	configPath := os.Getenv(CONFIG_AUTH_PATH)
	if configPath == "" {
		return nil, fmt.Errorf("%s environment variable not set", CONFIG_AUTH_PATH)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist %s", CONFIG_AUTH_PATH, configPath)
	}

	tokenManager := TokenManager{
		CacheRepo: cacheRepo,
	}

	if err := cleanenv.ReadConfig(configPath, &tokenManager); err != nil {
		return nil, fmt.Errorf("cannot load config file: %s", err)
	}
	return &tokenManager, nil
}

func (tm *TokenManager) Generate(userID, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.TokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.SecretKey))
}

func (tm *TokenManager) GetToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing token")
	}

	fields := strings.Split(authHeader, " ")
	if len(fields) != 2 || fields[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}
	return fields[1], nil
}

func (tm *TokenManager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tm.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
