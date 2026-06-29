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

// TripController 行程控制器
type TripController struct {
	tripSvc *service.TripService
}

// NewTripController 创建行程控制器
func NewTripController(tripSvc *service.TripService) *TripController {
	return &TripController{tripSvc: tripSvc}
}

// Create 创建行程
// POST /api/v1/trips/create
func (ctrl *TripController) Create(c *gin.Context) {
	var input service.TripCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errcode.ErrInvalidParam, "请求参数格式错误")
		return
	}
	if input.Name == "" {
		response.Error(c, errcode.ErrMissingParam, "行程名称不能为空")
		return
	}
	if input.TripDate == "" {
		response.Error(c, errcode.ErrMissingParam, "出行日期不能为空")
		return
	}

	userID := middleware.GetUserID(c)
	result, err := ctrl.tripSvc.Create(c.Request.Context(), uint64(userID), input)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// List 行程列表
// GET /api/v1/trips/list
func (ctrl *TripController) List(c *gin.Context) {
	p := util.ParsePagination(c)
	status := c.Query("status")
	if status != "" && status != "pending" && status != "active" && status != "completed" {
		response.Error(c, errcode.ErrInvalidParam, "status 参数无效")
		return
	}

	userID := uint64(middleware.GetUserID(c))
	result, err := ctrl.tripSvc.List(c.Request.Context(), userID, status, p)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// Detail 行程详情
// GET /api/v1/trips/:trip_id
func (ctrl *TripController) Detail(c *gin.Context) {
	tripID, err := strconv.ParseUint(c.Param("trip_id"), 10, 64)
	if err != nil || tripID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的行程ID")
		return
	}

	detail, err := ctrl.tripSvc.GetDetail(c.Request.Context(), tripID)
	if err != nil {
		if errors.Is(err, repository.ErrTripNotFound) {
			response.Error(c, errcode.ErrTripNotFound, errcode.GetMsg(errcode.ErrTripNotFound))
			return
		}
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, detail)
}

// Start 开始行程
// POST /api/v1/trips/:trip_id/start
func (ctrl *TripController) Start(c *gin.Context) {
	tripID, err := strconv.ParseUint(c.Param("trip_id"), 10, 64)
	if err != nil || tripID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的行程ID")
		return
	}
	userID := uint64(middleware.GetUserID(c))
	if err := ctrl.tripSvc.Start(c.Request.Context(), userID, tripID); err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, gin.H{"id": tripID, "status": "active"})
}

// Complete 结束行程
// POST /api/v1/trips/:trip_id/complete
func (ctrl *TripController) Complete(c *gin.Context) {
	tripID, err := strconv.ParseUint(c.Param("trip_id"), 10, 64)
	if err != nil || tripID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的行程ID")
		return
	}
	userID := uint64(middleware.GetUserID(c))
	if err := ctrl.tripSvc.Complete(c.Request.Context(), userID, tripID); err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, gin.H{"id": tripID, "status": "completed"})
}

// Delete 删除行程
// DELETE /api/v1/trips/:trip_id
func (ctrl *TripController) Delete(c *gin.Context) {
	tripID, err := strconv.ParseUint(c.Param("trip_id"), 10, 64)
	if err != nil || tripID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的行程ID")
		return
	}
	userID := uint64(middleware.GetUserID(c))
	if err := ctrl.tripSvc.Delete(c.Request.Context(), userID, tripID); err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, gin.H{"id": tripID})
}

// Update 更新行程基本信息
// PUT /api/v1/trips/:trip_id
func (ctrl *TripController) Update(c *gin.Context) {
	tripID, err := strconv.ParseUint(c.Param("trip_id"), 10, 64)
	if err != nil || tripID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的行程ID")
		return
	}
	var body struct {
		Name     string `json:"name"`
		TripDate string `json:"trip_date"`
		EndDate  string `json:"end_date"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.ErrInvalidParam, "请求参数格式错误")
		return
	}
	userID := uint64(middleware.GetUserID(c))
	if err := ctrl.tripSvc.Update(c.Request.Context(), userID, tripID, body.Name, body.TripDate, body.EndDate); err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, gin.H{"id": tripID})
}

// AddPoint 添加点位到行程
// POST /api/v1/trips/:trip_id/points
func (ctrl *TripController) AddPoint(c *gin.Context) {
	tripID, err := strconv.ParseUint(c.Param("trip_id"), 10, 64)
	if err != nil || tripID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的行程ID")
		return
	}
	var input service.TripPointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errcode.ErrInvalidParam, "请求参数格式错误")
		return
	}
	if input.Name == "" {
		response.Error(c, errcode.ErrMissingParam, "点位名称不能为空")
		return
	}
	userID := uint64(middleware.GetUserID(c))
	if err := ctrl.tripSvc.AddPoint(c.Request.Context(), userID, tripID, input); err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, gin.H{"trip_id": tripID})
}

// DeletePoint 删除行程点位
// DELETE /api/v1/trips/:trip_id/point/:point_id
func (ctrl *TripController) DeletePoint(c *gin.Context) {
	tripID, err := strconv.ParseUint(c.Param("trip_id"), 10, 64)
	if err != nil || tripID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的行程ID")
		return
	}
	pointID, err := strconv.ParseUint(c.Param("point_id"), 10, 64)
	if err != nil || pointID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的点位ID")
		return
	}
	userID := uint64(middleware.GetUserID(c))
	if err := ctrl.tripSvc.DeletePoint(c.Request.Context(), userID, tripID, pointID); err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, gin.H{"point_id": pointID})
}

// SkipPoint 跳过点位
// POST /api/v1/trips/:trip_id/point/:point_id/skip
func (ctrl *TripController) SkipPoint(c *gin.Context) {
	ctrl.updatePointStatus(c, "skipped")
}

// UnskipPoint 取消跳过点位
// POST /api/v1/trips/:trip_id/point/:point_id/unskip
func (ctrl *TripController) UnskipPoint(c *gin.Context) {
	ctrl.updatePointStatus(c, "pending")
}

func (ctrl *TripController) updatePointStatus(c *gin.Context, target string) {
	tripID, err := strconv.ParseUint(c.Param("trip_id"), 10, 64)
	if err != nil || tripID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的行程ID")
		return
	}
	pointID, err := strconv.ParseUint(c.Param("point_id"), 10, 64)
	if err != nil || pointID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的点位ID")
		return
	}
	userID := uint64(middleware.GetUserID(c))
	var err2 error
	if target == "skipped" {
		err2 = ctrl.tripSvc.SkipPoint(c.Request.Context(), userID, tripID, pointID)
	} else {
		err2 = ctrl.tripSvc.UnskipPoint(c.Request.Context(), userID, tripID, pointID)
	}
	if err2 != nil {
		ctrl.handleError(c, err2)
		return
	}
	response.Success(c, gin.H{"point_id": pointID, "checkin_status": target})
}

// handleError 处理行程服务错误
func (ctrl *TripController) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrTripNotFound):
		response.Error(c, errcode.ErrTripNotFound, errcode.GetMsg(errcode.ErrTripNotFound))
	case errors.Is(err, repository.ErrTripPointNotFound):
		response.Error(c, errcode.ErrTripPointNotFound, errcode.GetMsg(errcode.ErrTripPointNotFound))
	case errors.Is(err, service.ErrInvalidTripName):
		response.Error(c, errcode.ErrInvalidParam, "行程名称长度需为1-50个字符")
	case errors.Is(err, service.ErrInvalidTripDate):
		response.Error(c, errcode.ErrInvalidParam, "出行日期格式错误或结束日期早于开始日期")
	case errors.Is(err, service.ErrTripCreateLimit):
		response.Error(c, errcode.ErrTooManyRequests, "每天最多创建5个行程")
	case errors.Is(err, service.ErrTripNoPermission):
		response.Error(c, errcode.ErrNoPermission, errcode.GetMsg(errcode.ErrNoPermission))
	case errors.Is(err, repository.ErrTripStatusInvalid), errors.Is(err, service.ErrPointStatusInvalid):
		response.Error(c, errcode.ErrTripStatusInvalid, errcode.GetMsg(errcode.ErrTripStatusInvalid))
	default:
		zap.L().Error("行程服务异常",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Error(err),
		)
		response.Error(c, errcode.ErrInternal, errcode.GetMsg(errcode.ErrInternal))
	}
}
