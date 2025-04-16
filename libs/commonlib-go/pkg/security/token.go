package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewTokenHelper() TokenHelper {
	return TokenHelper{}
}

type tokenConfig struct {
	accessDuration  time.Duration
	refreshDuration time.Duration
	accessSecret    string
	refreshSecret   string
}

var tokenCfg = tokenConfig{
	accessDuration:  time.Minute * 5,
	refreshDuration: time.Minute * 30,
	accessSecret:    "access_secret",
	refreshSecret:   "refresh_secret",
}

type TokenHelper struct{}

func (h *TokenHelper) GenerateAccessToken(req AccessTokenClaimRequest) (string, error) {
	jwtId := uuid.New().String()
	now := time.Now()
	expirationTime := now.Add(tokenCfg.accessDuration)

	claims := AccessTokenClaims{
		AccessTokenClaimRequest: req,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jwtId,
			Issuer:    "super admin",
			Subject:   "access_token",
			Audience:  jwt.ClaimStrings{"admin"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenCfg.accessSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *TokenHelper) GenerateRefreshToken(req RefreshTokenClaimRequest) (string, error) {
	jwtId := uuid.New().String()
	now := time.Now()
	expirationTime := now.Add(tokenCfg.refreshDuration)

	claims := RefreshTokenClaims{
		RefreshTokenClaimRequest: req,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jwtId,
			Issuer:    "super_admin",
			Subject:   "refresh_token",
			Audience:  jwt.ClaimStrings{"admin"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenCfg.refreshSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *TokenHelper) ParseAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenCfg.accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func (h *TokenHelper) ParseRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenCfg.refreshSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

type (
	AccessTokenClaims struct {
		AccessTokenClaimRequest
		jwt.RegisteredClaims
	}

	AccessTokenClaimRequest struct {
		RefreshTokenId int64
		AdminId        int64
		Email          string
		FirstName      *string
		LastName       *string
		Roles          []string
	}

	RefreshTokenClaims struct {
		RefreshTokenClaimRequest
		jwt.RegisteredClaims
	}

	RefreshTokenClaimRequest struct {
		AdminId int64
	}
)
