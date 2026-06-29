// Package controller 控制器层（HTTP 入口）
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
	"go.uber.org/zap"
)

// UserController 用户控制器
type UserController struct {
	userSvc *service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userSvc *service.UserService) *UserController {
	return &UserController{userSvc: userSvc}
}

// GetProfile 获取当前登录用户信息
// GET /api/v1/user/profile  @auth
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrNotLoggedIn, errcode.GetMsg(errcode.ErrNotLoggedIn))
		return
	}

	profile, err := ctrl.userSvc.GetProfile(c.Request.Context(), uint64(userID))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.Error(c, errcode.ErrUserNotFound, errcode.GetMsg(errcode.ErrUserNotFound))
			return
		}
		zap.L().Error("获取用户信息失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Int64("user_id", userID),
			zap.Error(err),
		)
		response.Error(c, errcode.ErrInternal, errcode.GetMsg(errcode.ErrInternal))
		return
	}

	response.Success(c, profile)
}

// GetOtherProfile 获取其他用户公开信息
// GET /api/v1/user/{user_id}/profile  @optionalAuth
func (ctrl *UserController) GetOtherProfile(c *gin.Context) {
	userIDStr := c.Param("user_id")
	targetID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || targetID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的用户ID")
		return
	}

	// 当前登录用户ID（未登录为0）
	currentUserID := uint64(middleware.GetUserID(c))

	profile, err := ctrl.userSvc.GetOtherProfile(c.Request.Context(), targetID, currentUserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.Error(c, errcode.ErrUserNotFound, errcode.GetMsg(errcode.ErrUserNotFound))
			return
		}
		zap.L().Error("获取其他用户信息失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint64("target_id", targetID),
			zap.Error(err),
		)
		response.Error(c, errcode.ErrInternal, errcode.GetMsg(errcode.ErrInternal))
		return
	}

	response.Success(c, profile)
}

// updateProfileRequest 更新用户资料请求
type updateProfileRequest struct {
	Nickname   *string  `json:"nickname,omitempty"`
	Avatar     *string  `json:"avatar,omitempty"`
	Bio        *string  `json:"bio,omitempty"`
	Birthday   *string  `json:"birthday,omitempty"`   // YYYY-MM-DD
	PreferTags []string `json:"prefer_tags,omitempty"`
}

// UpdateProfile 更新用户资料
// PUT /api/v1/user/profile  @auth
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrNotLoggedIn, errcode.GetMsg(errcode.ErrNotLoggedIn))
		return
	}

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrInvalidParam, "请求参数格式错误")
		return
	}

	// 构建更新字段（仅更新非空字段）
	fields := make(map[string]interface{})
	if req.Nickname != nil {
		fields["nickname"] = *req.Nickname
	}
	if req.Avatar != nil {
		fields["avatar"] = *req.Avatar
	}
	if req.Bio != nil {
		fields["bio"] = *req.Bio
	}
	if req.Birthday != nil {
		fields["birthday"] = *req.Birthday
	}
	if req.PreferTags != nil {
		tagsJSON, _ := service.MarshalTags(req.PreferTags)
		fields["prefer_tags"] = tagsJSON
	}

	if len(fields) == 0 {
		response.Error(c, errcode.ErrMissingParam, "未提供任何更新字段")
		return
	}

	if err := ctrl.userSvc.UpdateProfile(c.Request.Context(), uint64(userID), fields); err != nil {
		if errors.Is(err, service.ErrUserBanned) {
			response.Error(c, errcode.ErrAccountBanned, errcode.GetMsg(errcode.ErrAccountBanned))
			return
		}
		// 参数校验类错误
		response.Error(c, errcode.ErrInvalidParam, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", nil)
}
