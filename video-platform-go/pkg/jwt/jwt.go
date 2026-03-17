package jwt

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	apperrors "github.com/video-platform/go/pkg/errors"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type JWTManager struct {
	secretKey          []byte
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
	issuer             string
	refreshTokenStore  map[string]*Claims
}

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration, issuer string) (*JWTManager, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("secret key must be at least 32 characters")
	}
	
	return &JWTManager{
		secretKey:         []byte(secret),
		accessTokenTTL:    accessTTL,
		refreshTokenTTL:   refreshTTL,
		issuer:            issuer,
		refreshTokenStore: make(map[string]*Claims),
	}, nil
}

func (m *JWTManager) GenerateTokenPair(ctx context.Context, userID int64, username, role string) (*TokenPair, error) {
	now := time.Now()
	
	accessClaims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}
	
	refreshClaims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}
	
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(m.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}
	
	refreshTokenString, err := m.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(m.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}
	
	m.refreshTokenStore[refreshTokenString] = refreshClaims
	
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(m.accessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (m *JWTManager) GenerateAccessToken(ctx context.Context, userID int64, username, role string) (string, error) {
	now := time.Now()
	
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}
	
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	
	return token, nil
}

func (m *JWTManager) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})
	
	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, apperrors.ErrTokenExpired
		}
		return nil, apperrors.ErrInvalidToken
	}
	
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, apperrors.ErrInvalidToken
	}
	
	return claims, nil
}

func (m *JWTManager) ValidateRefreshToken(ctx context.Context, refreshToken string) (*Claims, error) {
	claims, err := m.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	
	return claims, nil
}

func (m *JWTManager) RefreshTokenPair(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := m.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	
	return m.GenerateTokenPair(ctx, claims.UserID, claims.Username, claims.Role)
}

func (m *JWTManager) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	delete(m.refreshTokenStore, refreshToken)
	return nil
}

func (m *JWTManager) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	claims, ok := ctx.Value("userClaims").(*Claims)
	if !ok {
		return 0, false
	}
	return claims.UserID, true
}

func GetUsernameFromContext(ctx context.Context) (string, bool) {
	claims, ok := ctx.Value("userClaims").(*Claims)
	if !ok {
		return "", false
	}
	return claims.Username, true
}

func GetRoleFromContext(ctx context.Context) (string, bool) {
	claims, ok := ctx.Value("userClaims").(*Claims)
	if !ok {
		return "", false
	}
	return claims.Role, true
}

func ExtractTokenFromHeader(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", apperrors.ErrInvalidToken
	}
	return parts[1], nil
}
