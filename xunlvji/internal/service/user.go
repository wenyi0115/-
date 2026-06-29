// Package service 业务逻辑层
package service

import (
	"context"
	"errors"

	"github.com/xunlvji/xunlvji/internal/model"
	"github.com/xunlvji/xunlvji/internal/repository"
)

// 用户服务错误
var (
	ErrUserBanned = errors.New("user banned")
)

// UserProfile 用户完整信息（接口返回）
type UserProfile struct {
	ID               uint64         `json:"id"`
	Nickname         string         `json:"nickname"`
	Avatar           string         `json:"avatar"`
	Phone            string         `json:"phone"`
	Email            string         `json:"email"`
	CreditScore      int            `json:"credit_score"`
	PreferTags       []string       `json:"prefer_tags"`
	VisitedCityCount int            `json:"visited_city_count"`
	FollowerCount    int            `json:"follower_count"`
	FollowingCount   int            `json:"following_count"`
	LikeCount        int            `json:"like_count"`
	Language         string         `json:"language"`
	CreatedAt        string         `json:"created_at"`
}

// OtherUserProfile 其他用户公开信息
type OtherUserProfile struct {
	ID               uint64   `json:"id"`
	Nickname         string   `json:"nickname"`
	Avatar           string   `json:"avatar"`
	Bio              string   `json:"bio"`
	VisitedCityCount int      `json:"visited_city_count"`
	FollowerCount    int      `json:"follower_count"`
	FollowingCount   int      `json:"following_count"`
	LikeCount        int      `json:"like_count"`
	PreferTags       []string `json:"prefer_tags"`
	IsFollowing      bool     `json:"is_following"` // 当前登录用户是否关注了该用户
}

// UserService 用户服务
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetProfile 获取当前用户完整信息
func (s *UserService) GetProfile(ctx context.Context, userID uint64) (*UserProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.toProfile(user), nil
}

// GetOtherProfile 获取其他用户公开信息
// currentUserID: 当前登录用户ID（0表示未登录），用于判断 is_following
func (s *UserService) GetOtherProfile(ctx context.Context, userID, currentUserID uint64) (*OtherUserProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	// 被禁用/注销的用户对外不可见
	if user.Status != repository.UserStatusNormal {
		return nil, repository.ErrUserNotFound
	}

	profile := &OtherUserProfile{
		ID:               user.ID,
		Nickname:         user.Nickname,
		Avatar:           ptrToStr(user.Avatar),
		Bio:              ptrToStr(user.Bio),
		VisitedCityCount: user.VisitedCityCount,
		FollowerCount:    user.FollowerCount,
		FollowingCount:   user.FollowingCount,
		LikeCount:        user.LikeCount,
		PreferTags:       parseTags(user.PreferTags),
		IsFollowing:      false, // TODO: 后续接入 FollowRepository 后填充
	}
	_ = currentUserID // 暂未使用，关注关系后续接入
	return profile, nil
}

// UpdateProfile 更新用户资料
// 支持更新的字段：nickname, avatar, bio, birthday, prefer_tags
func (s *UserService) UpdateProfile(ctx context.Context, userID uint64, fields map[string]interface{}) error {
	// 校验昵称长度
	if nickname, ok := fields["nickname"].(string); ok {
		if len(nickname) < 1 || len(nickname) > 64 {
			return errors.New("昵称长度需在1-64之间")
		}
	}
	// 校验简介长度
	if bio, ok := fields["bio"].(string); ok {
		if len(bio) > 200 {
			return errors.New("个人简介过长（限200字）")
		}
	}
	return s.userRepo.UpdateFields(ctx, userID, fields)
}

// toProfile 转换为 UserProfile（含脱敏）
func (s *UserService) toProfile(user *model.User) *UserProfile {
	phone := ""
	if user.Phone != nil {
		phone = maskPhone(*user.Phone)
	}
	email := ""
	if user.Email != nil {
		email = maskEmail(*user.Email)
	}
	return &UserProfile{
		ID:               user.ID,
		Nickname:         user.Nickname,
		Avatar:           ptrToStr(user.Avatar),
		Phone:            phone,
		Email:            email,
		CreditScore:      user.CreditScore,
		PreferTags:       parseTags(user.PreferTags),
		VisitedCityCount: user.VisitedCityCount,
		FollowerCount:    user.FollowerCount,
		FollowingCount:   user.FollowingCount,
		LikeCount:        user.LikeCount,
		Language:         user.Language,
		CreatedAt:        user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// parseTags 解析 JSON 标签为字符串切片
func parseTags(raw []byte) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var tags []string
	if err := unmarshalJSON(raw, &tags); err != nil {
		return []string{}
	}
	return tags
}
