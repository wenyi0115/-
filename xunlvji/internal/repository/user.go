// Package repository 数据访问层
package repository

import (
	"context"
	"errors"
	"time"

	"github.com/xunlvji/xunlvji/internal/model"
	"gorm.io/gorm"
)

// 用户状态常量
const (
	UserStatusDisabled  int8 = 0 // 禁用
	UserStatusNormal    int8 = 1 // 正常
	UserStatusCanceled  int8 = 2 // 注销
)

// 第三方登录提供方
const (
	ProviderWechat = "wechat"
	ProviderApple  = "apple"
	ProviderGoogle = "google"
)

// 错误定义
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrOAuthNotFound     = errors.New("oauth not found")
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *model.User) error
	// GetByID 根据ID查询用户
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	// GetByPhone 根据手机号查询用户
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	// Update 更新用户信息（全字段）
	Update(ctx context.Context, user *model.User) error
	// UpdateFields 按字段更新用户
	UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error
}

// OAuthRepository 第三方账号关联数据访问接口
type OAuthRepository interface {
	// GetByProvider 根据 provider+oauthID 查询关联
	GetByProvider(ctx context.Context, provider, oauthID string) (*model.UserOAuth, error)
	// GetByUnionID 根据 unionID 查询关联（微信开放平台）
	GetByUnionID(ctx context.Context, unionID string) (*model.UserOAuth, error)
	// Create 创建第三方账号关联
	Create(ctx context.Context, oauth *model.UserOAuth) error
}

// userRepo 用户仓库实现
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepo) UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

// oauthRepo 第三方账号仓库实现
type oauthRepo struct {
	db *gorm.DB
}

// NewOAuthRepository 创建第三方账号仓库
func NewOAuthRepository(db *gorm.DB) OAuthRepository {
	return &oauthRepo{db: db}
}

func (r *oauthRepo) GetByProvider(ctx context.Context, provider, oauthID string) (*model.UserOAuth, error) {
	var o model.UserOAuth
	err := r.db.WithContext(ctx).
		Where("provider = ? AND oauth_id = ?", provider, oauthID).
		First(&o).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOAuthNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (r *oauthRepo) GetByUnionID(ctx context.Context, unionID string) (*model.UserOAuth, error) {
	var o model.UserOAuth
	err := r.db.WithContext(ctx).
		Where("provider = ? AND union_id = ?", ProviderWechat, unionID).
		First(&o).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOAuthNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (r *oauthRepo) Create(ctx context.Context, oauth *model.UserOAuth) error {
	return r.db.WithContext(ctx).Create(oauth).Error
}

// 辅助函数：取当前时间指针
func nowPtr() *time.Time {
	t := time.Now()
	return &t
}
