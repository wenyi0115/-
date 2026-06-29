package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/xunlvji/xunlvji/internal/model"
	"github.com/xunlvji/xunlvji/internal/repository"
	"github.com/xunlvji/xunlvji/pkg/util"
	"gorm.io/gorm"
)

// === 服务级错误 ===

var (
	ErrSocialTargetInvalid   = errors.New("invalid social target type")
	ErrSocialContentEmpty    = errors.New("comment content is empty")
	ErrSocialContentTooLong  = errors.New("comment content too long")
	ErrSocialReplyTooDeep    = errors.New("only two-level replies are allowed")
	ErrSocialNoPermission    = errors.New("no permission to operate")
	ErrSocialCommentTooFast  = errors.New("commenting too fast")
	ErrSocialTipIndexInvalid = errors.New("tip index out of range")
	ErrSocialCannotVoteSelf  = errors.New("cannot vote your own tip")
)

// 评论限频配置
const (
	commentRateLimitWindow = time.Hour // 1 小时窗口
	commentRateLimitMax    = 30        // 30 条/小时
	commentMaxLen          = 500
	replyPreviewCount      = 3 // 一级评论附带最多 3 条二级回复预览
)

// SocialService 社交互动服务
type SocialService struct {
	db             *gorm.DB
	userRepo       repository.UserRepository
	followRepo     repository.FollowRepository
	favoriteRepo   repository.FavoriteRepository
	commentRepo    repository.CommentRepository
	commentLikeRepo repository.CommentLikeRepository
	tipVoteRepo    repository.TipVoteRepository
	checkinRepo    repository.CheckinRepository
	placeRepo      repository.PlaceRepository
	routeRepo      repository.RouteRepository
	guideRepo      repository.GuideRepository
}

// NewSocialService 创建社交服务
func NewSocialService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
	favoriteRepo repository.FavoriteRepository,
	commentRepo repository.CommentRepository,
	commentLikeRepo repository.CommentLikeRepository,
	tipVoteRepo repository.TipVoteRepository,
	checkinRepo repository.CheckinRepository,
	placeRepo repository.PlaceRepository,
	routeRepo repository.RouteRepository,
	guideRepo repository.GuideRepository,
) *SocialService {
	return &SocialService{
		db:              db,
		userRepo:        userRepo,
		followRepo:      followRepo,
		favoriteRepo:    favoriteRepo,
		commentRepo:     commentRepo,
		commentLikeRepo: commentLikeRepo,
		tipVoteRepo:     tipVoteRepo,
		checkinRepo:     checkinRepo,
		placeRepo:       placeRepo,
		routeRepo:       routeRepo,
		guideRepo:       guideRepo,
	}
}

// === DTO 定义 ===

// FollowOpResult 关注/取关操作结果
type FollowOpResult struct {
	IsFollowing   bool `json:"is_following"`
	FollowerCount int  `json:"follower_count"`
}

// FollowUserItem 关注/粉丝列表项
type FollowUserItem struct {
	ID            uint64     `json:"id"`
	Nickname      string     `json:"nickname"`
	Avatar        string     `json:"avatar"`
	FollowerCount int        `json:"follower_count"`
	IsFollowing   bool       `json:"is_following"`
	FollowedAt    *time.Time `json:"followed_at"`
}

// FavoriteOpResult 收藏操作结果
type FavoriteOpResult struct {
	IsFavorited    bool `json:"is_favorited"`
	FavoriteCount int  `json:"favorite_count"`
}

// FavoriteListItem 收藏列表项
type FavoriteListItem struct {
	TargetType  string    `json:"target_type"`
	TargetID    uint64    `json:"target_id"`
	Title       string    `json:"title"`
	CoverImage  string    `json:"cover_image"`
	Category    string    `json:"category"`
	CategoryName string   `json:"category_name"`
	Rating      float64   `json:"rating"`
	FavoritedAt time.Time `json:"favorited_at"`
}

// CommentUserBrief 评论用户简略信息
type CommentUserBrief struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// CommentReplyItem 二级回复项
type CommentReplyItem struct {
	ID        uint64           `json:"id"`
	Content   string           `json:"content"`
	User      CommentUserBrief `json:"user"`
	ReplyTo   *CommentUserBrief `json:"reply_to,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}

// CommentListItem 评论列表项
type CommentListItem struct {
	ID         uint64              `json:"id"`
	Content    string              `json:"content"`
	User       CommentUserBrief    `json:"user"`
	ReplyCount int64               `json:"reply_count"`
	Replies    []CommentReplyItem  `json:"replies"`
	LikeCount  int                 `json:"like_count"`
	CreatedAt  time.Time           `json:"created_at"`
}

// CommentCreateResult 发布评论结果
type CommentCreateResult struct {
	ID        uint64           `json:"id"`
	Content   string           `json:"content"`
	User      CommentUserBrief `json:"user"`
	ParentID  *uint64          `json:"parent_id"`
	CreatedAt time.Time        `json:"created_at"`
}

// CommentLikeResult 评论点赞结果
type CommentLikeResult struct {
	Liked     bool `json:"liked"`
	LikeCount int  `json:"like_count"`
}

// TipVoteResult 避雷点赞结果
type TipVoteResult struct {
	IsVoted   bool `json:"is_voted"`
	VoteCount int  `json:"vote_count"`
}

// === 6.1 关注用户 ===

// Follow 关注指定用户
func (s *SocialService) Follow(ctx context.Context, userID, targetUserID uint64) (*FollowOpResult, error) {
	if userID == targetUserID {
		return nil, ErrSocialCannotVoteSelf // 复用为"不能关注自己"
	}

	// 幂等：已关注则直接返回
	exists, err := s.followRepo.Exists(ctx, userID, targetUserID)
	if err != nil {
		return nil, err
	}
	if !exists {
		// 事务：插入关注关系 + 双方计数更新
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			follow := &model.Follow{
				FollowerID:  userID,
				FollowingID: targetUserID,
			}
			if err := tx.Create(follow).Error; err != nil {
				return err
			}
			// 被关注者粉丝数 +1
			if err := tx.Model(&model.User{}).Where("id = ?", targetUserID).
				UpdateColumn("follower_count", gorm.Expr("follower_count + 1")).Error; err != nil {
				return err
			}
			// 关注者关注数 +1
			if err := tx.Model(&model.User{}).Where("id = ?", userID).
				UpdateColumn("following_count", gorm.Expr("following_count + 1")).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// 查询最新粉丝数
	followerCount, err := s.followRepo.CountFollowers(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	return &FollowOpResult{
		IsFollowing:   true,
		FollowerCount: int(followerCount),
	}, nil
}

// === 6.2 取消关注 ===

// Unfollow 取消关注
func (s *SocialService) Unfollow(ctx context.Context, userID, targetUserID uint64) (*FollowOpResult, error) {
	exists, err := s.followRepo.Exists(ctx, userID, targetUserID)
	if err != nil {
		return nil, err
	}
	if exists {
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			res := tx.Where("follower_id = ? AND following_id = ?", userID, targetUserID).
				Delete(&model.Follow{})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return nil // 已被删，幂等
			}
			if err := tx.Model(&model.User{}).Where("id = ?", targetUserID).
				UpdateColumn("follower_count", gorm.Expr("follower_count - 1")).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", userID).
				UpdateColumn("following_count", gorm.Expr("following_count - 1")).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	followerCount, err := s.followRepo.CountFollowers(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	return &FollowOpResult{
		IsFollowing:   false,
		FollowerCount: int(followerCount),
	}, nil
}

// === 6.3 关注列表 ===

// ListFollowing 关注列表
func (s *SocialService) ListFollowing(ctx context.Context, userID uint64, page, pageSize int) (*util.PageResult, error) {
	offset := (page - 1) * pageSize
	rels, total, err := s.followRepo.ListFollowing(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, err
	}
	if len(rels) == 0 {
		return s.emptyPage(page, pageSize, total), nil
	}

	// 收集被关注者 ID 并批量查询用户信息
	targetIDs := make([]uint64, 0, len(rels))
	for _, r := range rels {
		targetIDs = append(targetIDs, r.FollowingID)
	}
	userMap, err := s.batchUsers(ctx, targetIDs)
	if err != nil {
		return nil, err
	}
	// 批量查询当前用户是否回粉了这些被关注者（始终 true，因为就是关注列表本身）
	followedAtMap := make(map[uint64]time.Time, len(rels))
	for _, r := range rels {
		followedAtMap[r.FollowingID] = r.CreatedAt
	}

	list := make([]FollowUserItem, 0, len(rels))
	for _, id := range targetIDs {
		u, ok := userMap[id]
		if !ok {
			continue
		}
		fa := followedAtMap[id]
		list = append(list, buildFollowUserItem(u, true, &fa))
	}
	return &util.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  int64(offset+pageSize) < total,
	}, nil
}

// === 6.4 粉丝列表 ===

// ListFollowers 粉丝列表
func (s *SocialService) ListFollowers(ctx context.Context, userID uint64, page, pageSize int) (*util.PageResult, error) {
	offset := (page - 1) * pageSize
	rels, total, err := s.followRepo.ListFollowers(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, err
	}
	if len(rels) == 0 {
		return s.emptyPage(page, pageSize, total), nil
	}

	followerIDs := make([]uint64, 0, len(rels))
	for _, r := range rels {
		followerIDs = append(followerIDs, r.FollowerID)
	}
	userMap, err := s.batchUsers(ctx, followerIDs)
	if err != nil {
		return nil, err
	}
	// 查询当前用户是否回粉这些粉丝
	followingMap, err := s.followRepo.IsFollowingBatch(ctx, userID, followerIDs)
	if err != nil {
		return nil, err
	}
	followedAtMap := make(map[uint64]time.Time, len(rels))
	for _, r := range rels {
		followedAtMap[r.FollowerID] = r.CreatedAt
	}

	list := make([]FollowUserItem, 0, len(rels))
	for _, id := range followerIDs {
		u, ok := userMap[id]
		if !ok {
			continue
		}
		fa := followedAtMap[id]
		list = append(list, buildFollowUserItem(u, followingMap[id], &fa))
	}
	return &util.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  int64(offset+pageSize) < total,
	}, nil
}

// === 6.5 收藏 ===

// Favorite 收藏
func (s *SocialService) Favorite(ctx context.Context, userID uint64, targetType string, targetID uint64) (*FavoriteOpResult, error) {
	if !isValidFavoriteTarget(targetType) {
		return nil, ErrSocialTargetInvalid
	}
	exists, err := s.favoriteRepo.Exists(ctx, userID, targetType, targetID)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			fav := &model.Favorite{
				UserID:     userID,
				TargetType: targetType,
				TargetID:   targetID,
			}
			if err := tx.Create(fav).Error; err != nil {
				return err
			}
			// 仅 Route 表维护 favorite_count 字段，其他表后续扩展
			if targetType == repository.TargetTypeRoute {
				if err := tx.Model(&model.Route{}).Where("id = ?", targetID).
					UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	count, err := s.favoriteRepo.CountByTarget(ctx, targetType, targetID)
	if err != nil {
		return nil, err
	}
	return &FavoriteOpResult{
		IsFavorited:   true,
		FavoriteCount: int(count),
	}, nil
}

// === 6.6 取消收藏 ===

// Unfavorite 取消收藏
func (s *SocialService) Unfavorite(ctx context.Context, userID uint64, targetType string, targetID uint64) (*FavoriteOpResult, error) {
	if !isValidFavoriteTarget(targetType) {
		return nil, ErrSocialTargetInvalid
	}
	exists, err := s.favoriteRepo.Exists(ctx, userID, targetType, targetID)
	if err != nil {
		return nil, err
	}
	if exists {
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			res := tx.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
				Delete(&model.Favorite{})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return nil
			}
			if targetType == repository.TargetTypeRoute {
				if err := tx.Model(&model.Route{}).Where("id = ?", targetID).
					UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	count, err := s.favoriteRepo.CountByTarget(ctx, targetType, targetID)
	if err != nil {
		return nil, err
	}
	return &FavoriteOpResult{
		IsFavorited:   false,
		FavoriteCount: int(count),
	}, nil
}

// === 6.7 收藏列表 ===

// ListFavorites 收藏列表
func (s *SocialService) ListFavorites(ctx context.Context, userID uint64, targetType string, page, pageSize int) (*util.PageResult, error) {
	if targetType != "" && !isValidFavoriteTarget(targetType) {
		return nil, ErrSocialTargetInvalid
	}
	offset := (page - 1) * pageSize
	favs, total, err := s.favoriteRepo.List(ctx, userID, targetType, offset, pageSize)
	if err != nil {
		return nil, err
	}
	if len(favs) == 0 {
		return s.emptyPage(page, pageSize, total), nil
	}

	list := make([]FavoriteListItem, 0, len(favs))
	for _, f := range favs {
		item, err := s.buildFavoriteItem(ctx, f)
		if err != nil {
			continue // 目标已删除则跳过
		}
		list = append(list, item)
	}
	return &util.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  int64(offset+pageSize) < total,
	}, nil
}

// === 6.8 发布评论 ===

// CommentCreateInput 评论创建入参
type CommentCreateInput struct {
	TargetType string
	TargetID   uint64
	Content    string
	ParentID   *uint64
}

// CreateComment 发布评论
func (s *SocialService) CreateComment(ctx context.Context, userID uint64, input CommentCreateInput) (*CommentCreateResult, error) {
	if !isValidCommentTarget(input.TargetType) {
		return nil, ErrSocialTargetInvalid
	}
	content := strings.TrimSpace(input.Content)
	if content == "" {
		return nil, ErrSocialContentEmpty
	}
	if len([]rune(content)) > commentMaxLen {
		return nil, ErrSocialContentTooLong
	}

	// 限频：每用户每小时 30 条
	since := time.Now().Add(-commentRateLimitWindow)
	recentCount, err := s.commentRepo.CountByUserSince(ctx, userID, since)
	if err != nil {
		return nil, err
	}
	if recentCount >= int64(commentRateLimitMax) {
		return nil, ErrSocialCommentTooFast
	}

	// 二级回复校验：parent_id 必须是一级评论
	var parentID *uint64
	if input.ParentID != nil {
		parent, err := s.commentRepo.GetByID(ctx, *input.ParentID)
		if err != nil {
			return nil, err
		}
		if parent.ParentID != nil {
			return nil, ErrSocialReplyTooDeep
		}
		parentID = input.ParentID
	}

	// 写入评论 + 目标 comment_count +1
	c := &model.Comment{
		TargetType: input.TargetType,
		TargetID:   input.TargetID,
		UserID:     userID,
		ParentID:   parentID,
		Content:    content,
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		return s.incrTargetCommentCount(tx, input.TargetType, input.TargetID, 1)
	})
	if err != nil {
		return nil, err
	}

	// 组装用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &CommentCreateResult{
		ID:        c.ID,
		Content:   c.Content,
		User:      toCommentUserBrief(user),
		ParentID:  parentID,
		CreatedAt: c.CreatedAt,
	}, nil
}

// === 6.9 评论列表 ===

// ListComments 获取评论列表（一级 + 3 条二级预览）
func (s *SocialService) ListComments(ctx context.Context, targetType string, targetID uint64, sort string, page, pageSize int) (*util.PageResult, error) {
	if !isValidCommentTarget(targetType) {
		return nil, ErrSocialTargetInvalid
	}
	offset := (page - 1) * pageSize
	q := repository.CommentQuery{
		TargetType: targetType,
		TargetID:   targetID,
		Sort:       sort,
	}
	topComments, total, err := s.commentRepo.ListTopLevel(ctx, q, offset, pageSize)
	if err != nil {
		return nil, err
	}
	if len(topComments) == 0 {
		return s.emptyPage(page, pageSize, total), nil
	}

	// 收集所有一级评论 ID 和用户 ID
	commentIDs := make([]uint64, 0, len(topComments))
	userIDs := make([]uint64, 0, len(topComments))
	for _, c := range topComments {
		commentIDs = append(commentIDs, c.ID)
		userIDs = append(userIDs, c.UserID)
	}
	// 批量加载每个一级评论的前 3 条回复
	replyMap, replyUserIDs, err := s.batchLoadReplies(ctx, commentIDs, replyPreviewCount)
	if err != nil {
		return nil, err
	}
	// 加载回复涉及的所有用户（包括 reply_to 和一级评论作者）
	allUserIDs := append(userIDs, replyUserIDs...)
	allUserMap, err := s.batchUsers(ctx, uniqueUint64(allUserIDs))
	if err != nil {
		return nil, err
	}

	list := make([]CommentListItem, 0, len(topComments))
	for _, c := range topComments {
		replies := replyMap[c.ID]
		replyItems := make([]CommentReplyItem, 0, len(replies))
		// 一级评论 user 作为 reply_to
		parentUserBrief := buildUserBrief(allUserMap, c.UserID)
		for _, r := range replies {
			replyItems = append(replyItems, CommentReplyItem{
				ID:        r.ID,
				Content:   r.Content,
				User:      buildUserBrief(allUserMap, r.UserID),
				ReplyTo:   &parentUserBrief,
				CreatedAt: r.CreatedAt,
			})
		}
		replyCount, _ := s.commentRepo.CountReplies(ctx, c.ID)

		list = append(list, CommentListItem{
			ID:         c.ID,
			Content:    c.Content,
			User:       buildUserBrief(allUserMap, c.UserID),
			ReplyCount: replyCount,
			Replies:    replyItems,
			LikeCount:  c.LikeCount,
			CreatedAt:  c.CreatedAt,
		})
	}
	return &util.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  int64(offset+pageSize) < total,
	}, nil
}

// === 6.10 二级回复列表 ===

// ListCommentReplies 二级回复列表
func (s *SocialService) ListCommentReplies(ctx context.Context, parentID uint64, page, pageSize int) (*util.PageResult, error) {
	offset := (page - 1) * pageSize
	replies, total, err := s.commentRepo.ListReplies(ctx, parentID, offset, pageSize)
	if err != nil {
		return nil, err
	}
	if len(replies) == 0 {
		return s.emptyPage(page, pageSize, total), nil
	}

	// 获取父评论作者作为 reply_to
	parent, err := s.commentRepo.GetByID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	userIDs := make([]uint64, 0, len(replies)+1)
	userIDs = append(userIDs, parent.UserID)
	for _, r := range replies {
		userIDs = append(userIDs, r.UserID)
	}
	userMap, err := s.batchUsers(ctx, uniqueUint64(userIDs))
	if err != nil {
		return nil, err
	}

	parentBrief := buildUserBrief(userMap, parent.UserID)
	list := make([]CommentReplyItem, 0, len(replies))
	for _, r := range replies {
		list = append(list, CommentReplyItem{
			ID:        r.ID,
			Content:   r.Content,
			User:      buildUserBrief(userMap, r.UserID),
			ReplyTo:   &parentBrief,
			CreatedAt: r.CreatedAt,
		})
	}
	return &util.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  int64(offset+pageSize) < total,
	}, nil
}

// === 6.11 删除评论 ===

// DeleteComment 删除评论（仅本人；一级评论级联删除二级回复）
func (s *SocialService) DeleteComment(ctx context.Context, commentID, userID uint64) error {
	c, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}
	if c.UserID != userID {
		return ErrSocialNoPermission
	}
	// 事务：删除自身 + 删除子回复 + 目标 comment_count -（1 + 子回复数）
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先统计子回复数量（仅一级评论触发，删除前统计）
		delta := -1
		if c.ParentID == nil {
			var subCount int64
			if err := tx.Model(&model.Comment{}).Where("parent_id = ?", commentID).Count(&subCount).Error; err != nil {
				return err
			}
			delta = -1 - int(subCount)
			// 删除所有子回复
			if err := tx.Where("parent_id = ?", commentID).Delete(&model.Comment{}).Error; err != nil {
				return err
			}
		}
		// 删除评论本身（校验 user_id 防越权）
		res := tx.Where("id = ? AND user_id = ?", commentID, userID).Delete(&model.Comment{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return repository.ErrCommentNotFound
		}
		// 同步减目标 comment_count
		_ = s.incrTargetCommentCount(tx, c.TargetType, c.TargetID, delta)
		return nil
	})
}

// === 6.21 评论点赞（Toggle） ===

// ToggleCommentLike 评论点赞 / 取消点赞
func (s *SocialService) ToggleCommentLike(ctx context.Context, commentID, userID uint64) (*CommentLikeResult, error) {
	c, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	exists, err := s.commentLikeRepo.Exists(ctx, commentID, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		// 取消点赞
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).
				Delete(&model.CommentLike{}).Error; err != nil {
				return err
			}
			return tx.Model(&model.Comment{}).Where("id = ?", commentID).
				UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
		})
		if err != nil {
			return nil, err
		}
		return &CommentLikeResult{Liked: false, LikeCount: c.LikeCount - 1}, nil
	}
	// 点赞
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		like := &model.CommentLike{CommentID: commentID, UserID: userID}
		if err := tx.Create(like).Error; err != nil {
			return err
		}
		return tx.Model(&model.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
	})
	if err != nil {
		return nil, err
	}
	return &CommentLikeResult{Liked: true, LikeCount: c.LikeCount + 1}, nil
}

// === 6.12 / 6.13 避雷"有用"点赞 ===

// TipVote 避雷点赞
func (s *SocialService) TipVote(ctx context.Context, userID uint64, checkinID uint64, tipIndex int) (*TipVoteResult, error) {
	checkin, err := s.checkinRepo.GetByID(ctx, checkinID)
	if err != nil {
		return nil, err
	}
	// 不能给自己的避雷点赞
	if checkin.UserID == userID {
		return nil, ErrSocialCannotVoteSelf
	}
	// 校验 tip_index 范围
	var tips []json.RawMessage
	if len(checkin.PitfallTips) > 0 {
		if err := json.Unmarshal(checkin.PitfallTips, &tips); err != nil {
			return nil, ErrSocialTipIndexInvalid
		}
	}
	if tipIndex < 0 || tipIndex >= len(tips) {
		return nil, ErrSocialTipIndexInvalid
	}
	// 幂等
	exists, err := s.tipVoteRepo.Exists(ctx, checkinID, tipIndex, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		vote := &model.TipVote{
			CheckinID: checkinID,
			TipIndex:  tipIndex,
			UserID:    userID,
		}
		if err := s.tipVoteRepo.Create(ctx, vote); err != nil {
			return nil, err
		}
	}
	count, err := s.tipVoteRepo.CountByCheckinTip(ctx, checkinID, tipIndex)
	if err != nil {
		return nil, err
	}
	return &TipVoteResult{IsVoted: true, VoteCount: int(count)}, nil
}

// UnvoteTip 取消避雷点赞
func (s *SocialService) UnvoteTip(ctx context.Context, userID uint64, checkinID uint64, tipIndex int) (*TipVoteResult, error) {
	exists, err := s.tipVoteRepo.Exists(ctx, checkinID, tipIndex, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		if err := s.tipVoteRepo.Delete(ctx, checkinID, tipIndex, userID); err != nil && !errors.Is(err, repository.ErrTipVoteNotFound) {
			return nil, err
		}
	}
	count, err := s.tipVoteRepo.CountByCheckinTip(ctx, checkinID, tipIndex)
	if err != nil {
		return nil, err
	}
	return &TipVoteResult{IsVoted: false, VoteCount: int(count)}, nil
}

// === 内部辅助 ===

// batchUsers 批量查询用户信息
func (s *SocialService) batchUsers(ctx context.Context, ids []uint64) (map[uint64]*model.User, error) {
	result := make(map[uint64]*model.User, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	// 直接走 DB 批量查（UserRepository 接口未提供批量方法）
	var users []*model.User
	err := s.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		result[u.ID] = u
	}
	return result, nil
}

// batchLoadReplies 批量加载多个一级评论的回复预览
func (s *SocialService) batchLoadReplies(ctx context.Context, parentIDs []uint64, limit int) (map[uint64][]model.Comment, []uint64, error) {
	result := make(map[uint64][]model.Comment, len(parentIDs))
	if len(parentIDs) == 0 {
		return result, nil, nil
	}
	var allReplies []model.Comment
	err := s.db.WithContext(ctx).
		Where("parent_id IN ?", parentIDs).
		Order("parent_id ASC, created_at ASC").
		Find(&allReplies).Error
	if err != nil {
		return nil, nil, err
	}
	replyUserIDs := make([]uint64, 0)
	for _, r := range allReplies {
		if r.ParentID == nil {
			continue
		}
		pid := *r.ParentID
		bucket := result[pid]
		if len(bucket) >= limit {
			continue
		}
		bucket = append(bucket, r)
		result[pid] = bucket
		replyUserIDs = append(replyUserIDs, r.UserID)
	}
	return result, replyUserIDs, nil
}

// buildFavoriteItem 根据目标类型查表构建收藏项
func (s *SocialService) buildFavoriteItem(ctx context.Context, f model.Favorite) (FavoriteListItem, error) {
	item := FavoriteListItem{
		TargetType:  f.TargetType,
		TargetID:    f.TargetID,
		FavoritedAt: f.CreatedAt,
	}
	switch f.TargetType {
	case repository.TargetTypeRoute:
		r, err := s.routeRepo.GetByID(ctx, f.TargetID)
		if err != nil {
			return item, err
		}
		item.Title = r.Title
		item.CoverImage = firstJSONImage(r.CoverImage)
		item.Category = r.Category
		item.CategoryName = RouteCategoryName(r.Category)
		item.Rating = r.Rating.InexactFloat64()
	case repository.TargetTypePlace:
		p, err := s.placeRepo.GetByID(ctx, f.TargetID)
		if err != nil {
			return item, err
		}
		item.Title = p.Name
		item.CoverImage = firstJSONImage(p.CoverImage)
		item.Category = p.Category
		item.CategoryName = PlaceCategoryName(p.Category)
		item.Rating = p.Rating.InexactFloat64()
	case repository.TargetTypeGuide:
		g, err := s.guideRepo.GetByID(ctx, f.TargetID)
		if err != nil {
			return item, err
		}
		item.Title = g.Title
		item.CoverImage = firstJSONImage(g.CoverImage)
		item.Category = g.Category
		item.CategoryName = GuideCategoryName(g.Category)
		item.Rating = 0.0 // 攻略无评分字段，预留
	case repository.TargetTypeTrip:
		// P1 阶段扩展
		var t model.Trip
		if err := s.db.WithContext(ctx).First(&t, f.TargetID).Error; err != nil {
			return item, err
		}
		item.Title = t.Name
	}
	return item, nil
}

// incrTargetCommentCount 同步目标的 comment_count
func (s *SocialService) incrTargetCommentCount(tx *gorm.DB, targetType string, targetID uint64, delta int) error {
	switch targetType {
	case repository.TargetTypeRoute:
		return tx.Model(&model.Route{}).Where("id = ?", targetID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", delta)).Error
	case repository.TargetTypePlace:
		return tx.Model(&model.Place{}).Where("id = ?", targetID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", delta)).Error
	case repository.TargetTypeGuide:
		return tx.Model(&model.Guide{}).Where("id = ?", targetID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", delta)).Error
	}
	return nil
}

func (s *SocialService) emptyPage(page, pageSize int, total int64) *util.PageResult {
	return &util.PageResult{
		List:     []interface{}{},
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  false,
	}
}

// === 包级辅助函数 ===

func isValidFavoriteTarget(t string) bool {
	return t == repository.TargetTypeRoute ||
		t == repository.TargetTypePlace ||
		t == repository.TargetTypeGuide ||
		t == repository.TargetTypeTrip
}

func isValidCommentTarget(t string) bool {
	return t == repository.TargetTypeRoute ||
		t == repository.TargetTypePlace ||
		t == repository.TargetTypeGuide
}

func toCommentUserBrief(u *model.User) CommentUserBrief {
	b := CommentUserBrief{ID: u.ID, Nickname: u.Nickname}
	if u.Avatar != nil {
		b.Avatar = *u.Avatar
	}
	return b
}

func buildUserBrief(m map[uint64]*model.User, id uint64) CommentUserBrief {
	if u, ok := m[id]; ok {
		return toCommentUserBrief(u)
	}
	return CommentUserBrief{ID: id}
}

func buildFollowUserItem(u *model.User, isFollowing bool, followedAt *time.Time) FollowUserItem {
	item := FollowUserItem{
		ID:            u.ID,
		Nickname:      u.Nickname,
		FollowerCount: u.FollowerCount,
		IsFollowing:   isFollowing,
		FollowedAt:    followedAt,
	}
	if u.Avatar != nil {
		item.Avatar = *u.Avatar
	}
	return item
}

// firstJSONImage 从 json.RawMessage 中取第一张图 URL
func firstJSONImage(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var arr []string
	if err := json.Unmarshal(raw, &arr); err != nil {
		return ""
	}
	if len(arr) == 0 {
		return ""
	}
	return arr[0]
}

// uniqueUint64 去重 uint64 切片
func uniqueUint64(ids []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
