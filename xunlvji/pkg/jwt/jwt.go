package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xunlvji/xunlvji/internal/config"
)

// Claims 自定义 JWT 声明
type Claims struct {
	UserID   int64  `json:"user_id"`
	OpenID   string `json:"open_id,omitempty"` // 微信小程序 OpenID
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair Access Token + Refresh Token
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // access_token 剩余秒数
}

// TokenType Token 类型
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

var (
	ErrTokenInvalid  = errors.New("token 无效")
	ErrTokenExpired  = errors.New("token 已过期")
	ErrTokenType     = errors.New("token 类型不匹配")
)

// Manager JWT 管理器
type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

// NewManager 创建 JWT 管理器
func NewManager(cfg config.JWTConfig) *Manager {
	return &Manager{
		secret:     []byte(cfg.Secret),
		accessTTL:  time.Duration(cfg.AccessTTL) * time.Second,
		refreshTTL: time.Duration(cfg.RefreshTTL) * time.Second,
		issuer:     cfg.Issuer,
	}
}

// GenerateTokenPair 生成 Access Token + Refresh Token
func (m *Manager) GenerateTokenPair(userID int64, openID, nickname, avatar string) (*TokenPair, error) {
	access, err := m.generateToken(userID, openID, nickname, avatar, TokenTypeAccess, m.accessTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := m.generateToken(userID, openID, nickname, avatar, TokenTypeRefresh, m.refreshTTL)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(m.accessTTL.Seconds()),
	}, nil
}

// generateToken 生成单个 Token
func (m *Manager) generateToken(userID int64, openID, nickname, avatar string, tokenType TokenType, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		OpenID:   openID,
		Nickname: nickname,
		Avatar:   avatar,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Audience:  jwt.ClaimStrings{string(tokenType)},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken 解析并验证 Token
func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名方法: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

// ParseAccessToken 解析 Access Token（校验类型）
func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if !m.isType(claims, TokenTypeAccess) {
		return nil, ErrTokenType
	}
	return claims, nil
}

// ParseRefreshToken 解析 Refresh Token（校验类型）
func (m *Manager) ParseRefreshToken(tokenString string) (*Claims, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if !m.isType(claims, TokenTypeRefresh) {
		return nil, ErrTokenType
	}
	return claims, nil
}

// Refresh 刷新 Token（用 refresh_token 换取新的 access_token）
func (m *Manager) Refresh(refreshToken string) (*TokenPair, error) {
	claims, err := m.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	return m.GenerateTokenPair(claims.UserID, claims.OpenID, claims.Nickname, claims.Avatar)
}

// isType 校验 Token 类型
func (m *Manager) isType(claims *Claims, want TokenType) bool {
	if len(claims.Audience) == 0 {
		return false
	}
	for _, aud := range claims.Audience {
		if aud == string(want) {
			return true
		}
	}
	return false
}
