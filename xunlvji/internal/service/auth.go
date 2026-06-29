// Package service 业务逻辑层
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/xunlvji/xunlvji/internal/model"
	"github.com/xunlvji/xunlvji/internal/repository"
	"github.com/xunlvji/xunlvji/pkg/jwt"
	"github.com/xunlvji/xunlvji/pkg/wechat"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 认证服务错误
var (
	ErrWechatLoginFailed    = errors.New("wechat login failed")
	ErrAccountBanned        = errors.New("account banned")
	ErrAccountCanceled      = errors.New("account canceled")
	ErrRefreshTokenInvalid  = errors.New("refresh token invalid")
)

// LoginResult 登录结果
type LoginResult struct {
	AccessToken      string     `json:"token"`
	RefreshToken     string     `json:"refresh_token"`
	ExpireAt         time.Time  `json:"expire_at"`
	RefreshExpireAt  time.Time  `json:"refresh_expire_at"`
	User             *UserInfo  `json:"user"`
}

// UserInfo 用户简要信息（登录返回）
type UserInfo struct {
	ID           uint64 `json:"id"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	CreditScore  int    `json:"credit_score"`
	IsNewUser    bool   `json:"is_new_user"`
}

// AuthService 认证服务
type AuthService struct {
	db         *gorm.DB
	userRepo   repository.UserRepository
	oauthRepo  repository.OAuthRepository
	wechatCli  *wechat.Client
	jwtMgr     *jwt.Manager
}

// NewAuthService 创建认证服务
func NewAuthService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	oauthRepo repository.OAuthRepository,
	wechatCli *wechat.Client,
	jwtMgr *jwt.Manager,
) *AuthService {
	return &AuthService{
		db:        db,
		userRepo:  userRepo,
		oauthRepo: oauthRepo,
		wechatCli: wechatCli,
		jwtMgr:    jwtMgr,
	}
}

// WechatLogin 微信登录
// 流程：
//  1. code → openid/session_key（微信 code2session）
//  2. 按 openid 查 user_oauth，存在则登录老用户，不存在则注册新用户
//  3. 签发 JWT 双 Token
func (s *AuthService) WechatLogin(ctx context.Context, code string) (*LoginResult, error) {
	// 1. 调用微信 code2session
	wxResp, err := s.wechatCli.Code2Session(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWechatLoginFailed, err)
	}

	// 2. 查找或创建用户
	user, isNew, err := s.findOrCreateUser(ctx, wxResp.OpenID, wxResp.UnionID)
	if err != nil {
		return nil, err
	}

	// 3. 校验账号状态
	if err := s.checkUserStatus(user); err != nil {
		return nil, err
	}

	// 4. 签发 JWT
	tokenPair, err := s.jwtMgr.GenerateTokenPair(int64(user.ID), wxResp.OpenID, user.Nickname, ptrToStr(user.Avatar))
	if err != nil {
		return nil, fmt.Errorf("签发 token 失败: %w", err)
	}

	// 5. 组装返回
	return &LoginResult{
		AccessToken:     tokenPair.AccessToken,
		RefreshToken:    tokenPair.RefreshToken,
		ExpireAt:        time.Now().Add(time.Duration(tokenPair.ExpiresIn) * time.Second),
		RefreshExpireAt: time.Now().Add(30 * 24 * time.Hour), // 30天
		User: &UserInfo{
			ID:          user.ID,
			Nickname:    user.Nickname,
			Avatar:      ptrToStr(user.Avatar),
			CreditScore: user.CreditScore,
			IsNewUser:   isNew,
		},
	}, nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error) {
	tokenPair, err := s.jwtMgr.Refresh(refreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrRefreshTokenInvalid
		}
		return nil, ErrRefreshTokenInvalid
	}

	// 解析出用户ID，查询用户最新信息
	claims, err := s.jwtMgr.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}
	user, err := s.userRepo.GetByID(ctx, uint64(claims.UserID))
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if err := s.checkUserStatus(user); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:     tokenPair.AccessToken,
		RefreshToken:    tokenPair.RefreshToken,
		ExpireAt:        time.Now().Add(time.Duration(tokenPair.ExpiresIn) * time.Second),
		RefreshExpireAt: time.Now().Add(30 * 24 * time.Hour),
		User: &UserInfo{
			ID:          user.ID,
			Nickname:    user.Nickname,
			Avatar:      ptrToStr(user.Avatar),
			CreditScore: user.CreditScore,
			IsNewUser:   false,
		},
	}, nil
}

// findOrCreateUser 查找或创建用户（事务）
func (s *AuthService) findOrCreateUser(ctx context.Context, openID, unionID string) (*model.User, bool, error) {
	// 先查 openid
	oauth, err := s.oauthRepo.GetByProvider(ctx, repository.ProviderWechat, openID)
	if err != nil && !errors.Is(err, repository.ErrOAuthNotFound) {
		return nil, false, fmt.Errorf("查询 oauth 失败: %w", err)
	}

	// 老用户：直接返回
	if oauth != nil {
		user, err := s.userRepo.GetByID(ctx, oauth.UserID)
		if err != nil {
			return nil, false, fmt.Errorf("查询用户失败: %w", err)
		}
		return user, false, nil
	}

	// 新用户：事务创建 user + user_oauth
	var newUser *model.User
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 生成默认昵称（微信用户 + openid 后6位）
		nickname := generateDefaultNickname(openID)
		user := &model.User{
			Nickname:    nickname,
			Role:        "user",
			CreditScore: 100,
			Language:    "zh-CN",
			Status:      repository.UserStatusNormal,
		}
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}

		// 创建 oauth 关联
		oauthRecord := &model.UserOAuth{
			UserID:   user.ID,
			Provider: repository.ProviderWechat,
			OauthID:  openID,
		}
		if unionID != "" {
			oauthRecord.UnionID = &unionID
		}
		if err := tx.Create(oauthRecord).Error; err != nil {
			return fmt.Errorf("创建 oauth 关联失败: %w", err)
		}

		// 创建默认用户设置（可选，按 user_setting 表结构）
		// 这里简化处理，不强制创建；查询时若不存在用默认值
		newUser = user
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	zap.L().Info("新用户注册",
		zap.Uint64("user_id", newUser.ID),
		zap.String("provider", repository.ProviderWechat),
	)
	return newUser, true, nil
}

// checkUserStatus 校验用户状态
func (s *AuthService) checkUserStatus(user *model.User) error {
	switch user.Status {
	case repository.UserStatusDisabled:
		return ErrAccountBanned
	case repository.UserStatusCanceled:
		return ErrAccountCanceled
	}
	return nil
}

// generateDefaultNickname 生成默认昵称
func generateDefaultNickname(openID string) string {
	suffix := openID
	if len(suffix) > 6 {
		suffix = suffix[len(suffix)-6:]
	}
	return "微信用户" + suffix
}

// ptrToStr 指针字符串转字符串
func ptrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ptrToTime 时间指针转时间
func ptrToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// maskPhone 脱敏手机号 138****8000
func maskPhone(phone string) string {
	if len(phone) < 11 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// maskEmail 脱敏邮箱 u****@example.com
func maskEmail(email string) string {
	at := strings.Index(email, "@")
	if at <= 0 {
		return email
	}
	prefix := email[:at]
	domain := email[at:]
	if len(prefix) <= 1 {
		return prefix + "****" + domain
	}
	return prefix[:1] + "****" + domain
}
