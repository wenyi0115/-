package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xunlvji/xunlvji/internal/middleware"
	"github.com/xunlvji/xunlvji/internal/repository"
	"github.com/xunlvji/xunlvji/internal/service"
	"github.com/xunlvji/xunlvji/pkg/errcode"
	"github.com/xunlvji/xunlvji/pkg/response"
	"github.com/xunlvji/xunlvji/pkg/util"
	"go.uber.org/zap"
)

// SocialController 社交互动控制器
type SocialController struct {
	socialSvc *service.SocialService
}

// NewSocialController 创建社交控制器
func NewSocialController(socialSvc *service.SocialService) *SocialController {
	return &SocialController{socialSvc: socialSvc}
}

// === 6.1 关注用户 ===
// POST /api/v1/social/follow
func (ctrl *SocialController) Follow(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrMissingParam, "user_id 不能为空")
		return
	}

	userID := uint64(middleware.GetUserID(c))
	if userID == req.UserID {
		response.Error(c, errcode.ErrFollowSelf, errcode.GetMsg(errcode.ErrFollowSelf))
		return
	}
	result, err := ctrl.socialSvc.Follow(c.Request.Context(), userID, req.UserID)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.2 取消关注 ===
// POST /api/v1/social/unfollow
func (ctrl *SocialController) Unfollow(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrMissingParam, "user_id 不能为空")
		return
	}

	userID := uint64(middleware.GetUserID(c))
	result, err := ctrl.socialSvc.Unfollow(c.Request.Context(), userID, req.UserID)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.3 关注列表 ===
// GET /api/v1/social/following/list
func (ctrl *SocialController) FollowingList(c *gin.Context) {
	p := util.ParsePagination(c)
	userID := uint64(middleware.GetUserID(c))

	result, err := ctrl.socialSvc.ListFollowing(c.Request.Context(), userID, p.Page, p.PageSize)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.4 粉丝列表 ===
// GET /api/v1/social/follower/list
func (ctrl *SocialController) FollowerList(c *gin.Context) {
	p := util.ParsePagination(c)
	userID := uint64(middleware.GetUserID(c))

	result, err := ctrl.socialSvc.ListFollowers(c.Request.Context(), userID, p.Page, p.PageSize)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.5 收藏 ===
// POST /api/v1/social/favorite
func (ctrl *SocialController) Favorite(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint64 `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrMissingParam, "target_type 和 target_id 不能为空")
		return
	}

	userID := uint64(middleware.GetUserID(c))
	result, err := ctrl.socialSvc.Favorite(c.Request.Context(), userID, req.TargetType, req.TargetID)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.6 取消收藏 ===
// POST /api/v1/social/unfavorite
func (ctrl *SocialController) Unfavorite(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint64 `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrMissingParam, "target_type 和 target_id 不能为空")
		return
	}

	userID := uint64(middleware.GetUserID(c))
	result, err := ctrl.socialSvc.Unfavorite(c.Request.Context(), userID, req.TargetType, req.TargetID)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.7 收藏列表 ===
// GET /api/v1/social/favorite/list
func (ctrl *SocialController) FavoriteList(c *gin.Context) {
	p := util.ParsePagination(c)
	targetType := c.Query("target_type")

	userID := uint64(middleware.GetUserID(c))
	result, err := ctrl.socialSvc.ListFavorites(c.Request.Context(), userID, targetType, p.Page, p.PageSize)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.8 发布评论 ===
// POST /api/v1/social/comment
func (ctrl *SocialController) CreateComment(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint64 `json:"target_id" binding:"required"`
		Content    string `json:"content"`
		ParentID   *uint64 `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrMissingParam, "target_type 和 target_id 不能为空")
		return
	}

	userID := uint64(middleware.GetUserID(c))
	input := service.CommentCreateInput{
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Content:    req.Content,
		ParentID:   req.ParentID,
	}
	result, err := ctrl.socialSvc.CreateComment(c.Request.Context(), userID, input)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.9 评论列表 ===
// GET /api/v1/social/comment/list
func (ctrl *SocialController) CommentList(c *gin.Context) {
	targetType := c.Query("target_type")
	targetIDStr := c.Query("target_id")
	if targetType == "" || targetIDStr == "" {
		response.Error(c, errcode.ErrMissingParam, "target_type 和 target_id 不能为空")
		return
	}
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil || targetID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的 target_id")
		return
	}
	sort := c.DefaultQuery("sort", "newest")
	if sort != "newest" && sort != "hottest" {
		sort = "newest"
	}
	p := util.ParsePagination(c)

	result, err := ctrl.socialSvc.ListComments(c.Request.Context(), targetType, targetID, sort, p.Page, p.PageSize)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.10 二级回复列表 ===
// GET /api/v1/social/comment/replies
func (ctrl *SocialController) CommentReplies(c *gin.Context) {
	commentIDStr := c.Query("comment_id")
	if commentIDStr == "" {
		response.Error(c, errcode.ErrMissingParam, "comment_id 不能为空")
		return
	}
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil || commentID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的 comment_id")
		return
	}
	p := util.ParsePagination(c)

	result, err := ctrl.socialSvc.ListCommentReplies(c.Request.Context(), commentID, p.Page, p.PageSize)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.11 删除评论 ===
// DELETE /api/v1/social/comment/:comment_id
func (ctrl *SocialController) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil || commentID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的 comment_id")
		return
	}

	userID := uint64(middleware.GetUserID(c))
	if err := ctrl.socialSvc.DeleteComment(c.Request.Context(), commentID, userID); err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// === 6.21 评论点赞（Toggle） ===
// POST /api/v1/social/comments/:id/like
func (ctrl *SocialController) ToggleCommentLike(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || commentID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的评论ID")
		return
	}

	userID := uint64(middleware.GetUserID(c))
	result, err := ctrl.socialSvc.ToggleCommentLike(c.Request.Context(), commentID, userID)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.12 避雷"有用"点赞 ===
// POST /api/v1/social/tip-vote
func (ctrl *SocialController) TipVote(c *gin.Context) {
	var req struct {
		CheckinID uint64 `json:"checkin_id" binding:"required"`
		TipIndex  int    `json:"tip_index"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrMissingParam, "checkin_id 不能为空")
		return
	}

	userID := uint64(middleware.GetUserID(c))
	result, err := ctrl.socialSvc.TipVote(c.Request.Context(), userID, req.CheckinID, req.TipIndex)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// === 6.13 取消避雷点赞 ===
// POST /api/v1/social/tip-unvote
func (ctrl *SocialController) UnvoteTip(c *gin.Context) {
	var req struct {
		CheckinID uint64 `json:"checkin_id" binding:"required"`
		TipIndex  int    `json:"tip_index"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrMissingParam, "checkin_id 不能为空")
		return
	}

	userID := uint64(middleware.GetUserID(c))
	result, err := ctrl.socialSvc.UnvoteTip(c.Request.Context(), userID, req.CheckinID, req.TipIndex)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// handleError 统一错误处理
func (ctrl *SocialController) handleError(c *gin.Context, err error) {
	zap.L().Error("social controller error",
		zap.String("path", c.Request.URL.Path),
		zap.Error(err))

	switch {
	case errors.Is(err, service.ErrSocialTargetInvalid):
		response.Error(c, errcode.ErrInvalidParam, "无效的 target_type")
	case errors.Is(err, service.ErrSocialContentEmpty):
		response.Error(c, errcode.ErrCommentEmpty, errcode.GetMsg(errcode.ErrCommentEmpty))
	case errors.Is(err, service.ErrSocialContentTooLong):
		response.Error(c, errcode.ErrCommentTooLong, errcode.GetMsg(errcode.ErrCommentTooLong))
	case errors.Is(err, service.ErrSocialReplyTooDeep):
		response.Error(c, errcode.ErrInvalidParam, "仅支持二级回复")
	case errors.Is(err, service.ErrSocialNoPermission):
		response.Error(c, errcode.ErrNoPermission, "无权操作他人评论")
	case errors.Is(err, service.ErrSocialCommentTooFast):
		response.Error(c, errcode.ErrCommentTooFast, errcode.GetMsg(errcode.ErrCommentTooFast))
	case errors.Is(err, service.ErrSocialTipIndexInvalid):
		response.Error(c, errcode.ErrInvalidParam, "避雷建议索引无效")
	case errors.Is(err, service.ErrSocialCannotVoteSelf):
		// 关注自己 / 给自己避雷点赞都走此分支（关注自己在 controller 顶部已拦截）
		response.Error(c, errcode.ErrOperationSelf, "不能给自己的避雷建议点赞")
	case errors.Is(err, repository.ErrCommentNotFound):
		response.Error(c, errcode.ErrCommentNotFound, errcode.GetMsg(errcode.ErrCommentNotFound))
	case errors.Is(err, repository.ErrTipVoteNotFound):
		response.Error(c, errcode.ErrInvalidParam, "未点赞该避雷建议")
	case errors.Is(err, repository.ErrUserNotFound):
		response.Error(c, errcode.ErrUserNotFound, "用户不存在")
	default:
		response.Error(c, errcode.ErrInternal, "服务异常，请稍后重试")
	}
}
