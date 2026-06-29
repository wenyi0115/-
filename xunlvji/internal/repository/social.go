package repository

import (
	"context"
	"errors"
	"time"

	"github.com/xunlvji/xunlvji/internal/model"
	"gorm.io/gorm"
)

// 社交模块错误定义
var (
	ErrFollowNotFound     = errors.New("follow relation not found")
	ErrFavoriteNotFound   = errors.New("favorite not found")
	ErrCommentNotFound    = errors.New("comment not found")
	ErrCommentLikeExists  = errors.New("comment like already exists")
	ErrTipVoteNotFound    = errors.New("tip vote not found")
)

// 目标类型常量
const (
	TargetTypeRoute  = "route"
	TargetTypePlace  = "place"
	TargetTypeGuide  = "guide"
	TargetTypeTrip   = "trip"
	TargetTypeCheckin = "checkin"
)

// =============== Follow ===============

// FollowRepository 关注关系数据访问接口
type FollowRepository interface {
	// Create 创建关注关系（幂等：已存在返回 ErrAlreadyFollowed 由上层判断）
	Create(ctx context.Context, f *model.Follow) error
	// Delete 取消关注
	Delete(ctx context.Context, followerID, followingID uint64) error
	// Exists 是否存在关注关系
	Exists(ctx context.Context, followerID, followingID uint64) (bool, error)
	// ListFollowing 获取关注列表（按关注时间倒序）
	ListFollowing(ctx context.Context, followerID uint64, offset, limit int) ([]model.Follow, int64, error)
	// ListFollowers 获取粉丝列表
	ListFollowers(ctx context.Context, followingID uint64, offset, limit int) ([]model.Follow, int64, error)
	// CountFollowers 粉丝数
	CountFollowers(ctx context.Context, userID uint64) (int64, error)
	// CountFollowing 关注数
	CountFollowing(ctx context.Context, userID uint64) (int64, error)
	// IsFollowingBatch 批量查询当前用户是否关注了一组用户
	IsFollowingBatch(ctx context.Context, followerID uint64, targetIDs []uint64) (map[uint64]bool, error)
}

type followRepo struct {
	db *gorm.DB
}

// NewFollowRepository 创建关注仓库
func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepo{db: db}
}

func (r *followRepo) Create(ctx context.Context, f *model.Follow) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *followRepo) Delete(ctx context.Context, followerID, followingID uint64) error {
	res := r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&model.Follow{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrFollowNotFound
	}
	return nil
}

func (r *followRepo) Exists(ctx context.Context, followerID, followingID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r *followRepo) ListFollowing(ctx context.Context, followerID uint64, offset, limit int) ([]model.Follow, int64, error) {
	var (
		list  []model.Follow
		total int64
	)
	db := r.db.WithContext(ctx).Model(&model.Follow{}).Where("follower_id = ?", followerID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *followRepo) ListFollowers(ctx context.Context, followingID uint64, offset, limit int) ([]model.Follow, int64, error) {
	var (
		list  []model.Follow
		total int64
	)
	db := r.db.WithContext(ctx).Model(&model.Follow{}).Where("following_id = ?", followingID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *followRepo) CountFollowers(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Follow{}).
		Where("following_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *followRepo) CountFollowing(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Follow{}).
		Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *followRepo) IsFollowingBatch(ctx context.Context, followerID uint64, targetIDs []uint64) (map[uint64]bool, error) {
	result := make(map[uint64]bool, len(targetIDs))
	if len(targetIDs) == 0 {
		return result, nil
	}
	var rows []model.Follow
	err := r.db.WithContext(ctx).
		Select("following_id").
		Where("follower_id = ? AND following_id IN ?", followerID, targetIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, f := range rows {
		result[f.FollowingID] = true
	}
	return result, nil
}

// =============== Favorite ===============

// FavoriteRepository 收藏数据访问接口
type FavoriteRepository interface {
	Create(ctx context.Context, f *model.Favorite) error
	Delete(ctx context.Context, userID uint64, targetType string, targetID uint64) error
	Exists(ctx context.Context, userID uint64, targetType string, targetID uint64) (bool, error)
	List(ctx context.Context, userID uint64, targetType string, offset, limit int) ([]model.Favorite, int64, error)
	CountByTarget(ctx context.Context, targetType string, targetID uint64) (int64, error)
}

type favoriteRepo struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏仓库
func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepo{db: db}
}

func (r *favoriteRepo) Create(ctx context.Context, f *model.Favorite) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *favoriteRepo) Delete(ctx context.Context, userID uint64, targetType string, targetID uint64) error {
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Delete(&model.Favorite{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrFavoriteNotFound
	}
	return nil
}

func (r *favoriteRepo) Exists(ctx context.Context, userID uint64, targetType string, targetID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Favorite{}).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Count(&count).Error
	return count > 0, err
}

func (r *favoriteRepo) List(ctx context.Context, userID uint64, targetType string, offset, limit int) ([]model.Favorite, int64, error) {
	var (
		list  []model.Favorite
		total int64
	)
	db := r.db.WithContext(ctx).Model(&model.Favorite{}).Where("user_id = ?", userID)
	if targetType != "" {
		db = db.Where("target_type = ?", targetType)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *favoriteRepo) CountByTarget(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Favorite{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error
	return count, err
}

// =============== Comment ===============

// CommentQuery 评论列表查询参数
type CommentQuery struct {
	TargetType string
	TargetID   uint64
	Sort       string // newest / hottest
}

// CommentRepository 评论数据访问接口
type CommentRepository interface {
	Create(ctx context.Context, c *model.Comment) error
	GetByID(ctx context.Context, id uint64) (*model.Comment, error)
	Delete(ctx context.Context, id uint64, userID uint64) error // 仅本人可删，RowsAffected校验
	ListTopLevel(ctx context.Context, q CommentQuery, offset, limit int) ([]model.Comment, int64, error)
	ListReplies(ctx context.Context, parentID uint64, offset, limit int) ([]model.Comment, int64, error)
	CountReplies(ctx context.Context, parentID uint64) (int64, error)
	IncrLikeCount(ctx context.Context, id uint64, delta int) error
	CountCommentsByTarget(ctx context.Context, targetType string, targetID uint64) (int64, error)
	// 限频：统计用户近 N 分钟内发布的评论数
	CountByUserSince(ctx context.Context, userID uint64, since time.Time) (int64, error)
}

type commentRepo struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓库
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepo{db: db}
}

func (r *commentRepo) Create(ctx context.Context, c *model.Comment) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *commentRepo) GetByID(ctx context.Context, id uint64) (*model.Comment, error) {
	var c model.Comment
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *commentRepo) Delete(ctx context.Context, id uint64, userID uint64) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Comment{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCommentNotFound
	}
	return nil
}

func (r *commentRepo) ListTopLevel(ctx context.Context, q CommentQuery, offset, limit int) ([]model.Comment, int64, error) {
	var (
		list  []model.Comment
		total int64
	)
	db := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("target_type = ? AND target_id = ? AND parent_id IS NULL", q.TargetType, q.TargetID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	order := "created_at DESC"
	if q.Sort == "hottest" {
		order = "like_count DESC, created_at DESC"
	}
	if err := db.Order(order).Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *commentRepo) ListReplies(ctx context.Context, parentID uint64, offset, limit int) ([]model.Comment, int64, error) {
	var (
		list  []model.Comment
		total int64
	)
	db := r.db.WithContext(ctx).Model(&model.Comment{}).Where("parent_id = ?", parentID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at ASC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *commentRepo) CountReplies(ctx context.Context, parentID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}

func (r *commentRepo) IncrLikeCount(ctx context.Context, id uint64, delta int) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

func (r *commentRepo) CountCommentsByTarget(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error
	return count, err
}

func (r *commentRepo) CountByUserSince(ctx context.Context, userID uint64, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Count(&count).Error
	return count, err
}

// =============== CommentLike ===============

// CommentLikeRepository 评论点赞数据访问接口
type CommentLikeRepository interface {
	Create(ctx context.Context, l *model.CommentLike) error
	Delete(ctx context.Context, commentID, userID uint64) error
	Exists(ctx context.Context, commentID, userID uint64) (bool, error)
}

type commentLikeRepo struct {
	db *gorm.DB
}

// NewCommentLikeRepository 创建评论点赞仓库
func NewCommentLikeRepository(db *gorm.DB) CommentLikeRepository {
	return &commentLikeRepo{db: db}
}

func (r *commentLikeRepo) Create(ctx context.Context, l *model.CommentLike) error {
	return r.db.WithContext(ctx).Create(l).Error
}

func (r *commentLikeRepo) Delete(ctx context.Context, commentID, userID uint64) error {
	res := r.db.WithContext(ctx).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Delete(&model.CommentLike{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCommentLikeExists // 复用为"未点赞"标识，上层忽略 RowsAffected=0
	}
	return nil
}

func (r *commentLikeRepo) Exists(ctx context.Context, commentID, userID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CommentLike{}).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Count(&count).Error
	return count > 0, err
}

// =============== TipVote ===============

// TipVoteRepository 避雷点赞数据访问接口
type TipVoteRepository interface {
	Create(ctx context.Context, v *model.TipVote) error
	Delete(ctx context.Context, checkinID uint64, tipIndex int, userID uint64) error
	Exists(ctx context.Context, checkinID uint64, tipIndex int, userID uint64) (bool, error)
	CountByCheckinTip(ctx context.Context, checkinID uint64, tipIndex int) (int64, error)
}

type tipVoteRepo struct {
	db *gorm.DB
}

// NewTipVoteRepository 创建避雷点赞仓库
func NewTipVoteRepository(db *gorm.DB) TipVoteRepository {
	return &tipVoteRepo{db: db}
}

func (r *tipVoteRepo) Create(ctx context.Context, v *model.TipVote) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *tipVoteRepo) Delete(ctx context.Context, checkinID uint64, tipIndex int, userID uint64) error {
	res := r.db.WithContext(ctx).
		Where("checkin_id = ? AND tip_index = ? AND user_id = ?", checkinID, tipIndex, userID).
		Delete(&model.TipVote{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrTipVoteNotFound
	}
	return nil
}

func (r *tipVoteRepo) Exists(ctx context.Context, checkinID uint64, tipIndex int, userID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.TipVote{}).
		Where("checkin_id = ? AND tip_index = ? AND user_id = ?", checkinID, tipIndex, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *tipVoteRepo) CountByCheckinTip(ctx context.Context, checkinID uint64, tipIndex int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.TipVote{}).
		Where("checkin_id = ? AND tip_index = ?", checkinID, tipIndex).
		Count(&count).Error
	return count, err
}
