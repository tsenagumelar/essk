package authn

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tsenagumelar/essk/services/backend/internal/config"
)

type Claims struct {
	UserID   uuid.UUID  `json:"user_id"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
	Email    string     `json:"email"`
	jwt.RegisteredClaims
}

type TokenService struct {
	cfg config.AuthConfig
}

func NewTokenService(cfg config.AuthConfig) TokenService {
	return TokenService{cfg: cfg}
}

func (s TokenService) CreateAccessToken(userID uuid.UUID, tenantID *uuid.UUID, email string, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(s.cfg.AccessTokenTTL)
	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    s.cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.SigningKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func (s TokenService) ParseAccessToken(rawToken string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (any, error) {
		return []byte(s.cfg.SigningKey), nil
	}, jwt.WithIssuer(s.cfg.Issuer))
	if err != nil {
		return Claims{}, err
	}
	if !token.Valid {
		return Claims{}, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
