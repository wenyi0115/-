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

// ContentController 内容控制器
type ContentController struct {
	contentSvc *service.ContentService
}

// NewContentController 创建内容控制器
func NewContentController(contentSvc *service.ContentService) *ContentController {
	return &ContentController{contentSvc: contentSvc}
}

// RecommendFeed 推荐流
// GET /api/v1/content/feed/recommend
func (ctrl *ContentController) RecommendFeed(c *gin.Context) {
	p := util.ParsePagination(c)
	currentUserID := uint64(middleware.GetUserID(c))

	result, err := ctrl.contentSvc.ListRecommendFeed(c.Request.Context(), currentUserID, p)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// GuideFeed 攻略流
// GET /api/v1/content/feed/guide
func (ctrl *ContentController) GuideFeed(c *gin.Context) {
	p := util.ParsePagination(c)
	category := c.Query("category")
	city := c.Query("city")
	sort := c.DefaultQuery("sort", "newest")

	result, err := ctrl.contentSvc.ListGuideFeed(c.Request.Context(), category, city, sort, p)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// PlaceFeed 打卡地流
// GET /api/v1/content/feed/place
func (ctrl *ContentController) PlaceFeed(c *gin.Context) {
	p := util.ParsePagination(c)
	category := c.Query("category")
	city := c.Query("city")
	sort := c.DefaultQuery("sort", "newest")

	result, err := ctrl.contentSvc.ListPlaceFeed(c.Request.Context(), category, city, sort, p)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// RouteFeed 路线流
// GET /api/v1/content/feed/route
func (ctrl *ContentController) RouteFeed(c *gin.Context) {
	p := util.ParsePagination(c)
	// 路线流默认 page_size=10
	if c.Query("page_size") == "" {
		p.PageSize = 10
		p.Offset = (p.Page - 1) * p.PageSize
	}
	category := c.Query("category")
	city := c.Query("city")
	sort := c.DefaultQuery("sort", "hot")

	result, err := ctrl.contentSvc.ListRouteFeed(c.Request.Context(), category, city, sort, p)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// NearbyFeed 附近流
// GET /api/v1/content/feed/nearby
func (ctrl *ContentController) NearbyFeed(c *gin.Context) {
	p := util.ParsePagination(c)
	locationMode := c.DefaultQuery("location_mode", "travel")
	city := c.Query("city")
	lng := util.ParseFloat(c.Query("longitude"), 0)
	lat := util.ParseFloat(c.Query("latitude"), 0)
	radius := util.ParseInt(c, "radius", 10000)

	result, err := ctrl.contentSvc.ListNearbyFeed(c.Request.Context(), locationMode, city, lng, lat, radius, p)
	if err != nil {
		if errors.Is(err, service.ErrLocationDenied) {
			response.Error(c, errcode.ErrLocationDenied, errcode.GetMsg(errcode.ErrLocationDenied))
			return
		}
		if errors.Is(err, service.ErrMissingCity) {
			response.Error(c, errcode.ErrMissingParam, "plan/manual 模式必须提供 city 参数")
			return
		}
		if errors.Is(err, service.ErrInvalidLocationMode) {
			response.Error(c, errcode.ErrInvalidParam, "location_mode 必须为 travel/plan/manual")
			return
		}
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, result)
}

// RouteDetail 路线详情
// GET /api/v1/content/route/:route_id
func (ctrl *ContentController) RouteDetail(c *gin.Context) {
	routeID, err := strconv.ParseUint(c.Param("route_id"), 10, 64)
	if err != nil || routeID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的路线ID")
		return
	}

	detail, err := ctrl.contentSvc.GetRouteDetail(c.Request.Context(), routeID)
	if err != nil {
		if errors.Is(err, repository.ErrRouteNotFound) {
			response.Error(c, errcode.ErrRouteNotFound, errcode.GetMsg(errcode.ErrRouteNotFound))
			return
		}
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, detail)
}

// PlaceDetail 地点详情
// GET /api/v1/content/place/:place_id
func (ctrl *ContentController) PlaceDetail(c *gin.Context) {
	placeID, err := strconv.ParseUint(c.Param("place_id"), 10, 64)
	if err != nil || placeID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的地点ID")
		return
	}

	lng := util.ParseFloat(c.Query("longitude"), 0)
	lat := util.ParseFloat(c.Query("latitude"), 0)

	detail, err := ctrl.contentSvc.GetPlaceDetail(c.Request.Context(), placeID, lng, lat)
	if err != nil {
		if errors.Is(err, repository.ErrPlaceNotFound) {
			response.Error(c, errcode.ErrPlaceNotFound, errcode.GetMsg(errcode.ErrPlaceNotFound))
			return
		}
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, detail)
}

// GuideDetail 攻略详情
// GET /api/v1/content/guide/:guide_id
func (ctrl *ContentController) GuideDetail(c *gin.Context) {
	guideID, err := strconv.ParseUint(c.Param("guide_id"), 10, 64)
	if err != nil || guideID == 0 {
		response.Error(c, errcode.ErrInvalidParam, "无效的攻略ID")
		return
	}

	detail, err := ctrl.contentSvc.GetGuideDetail(c.Request.Context(), guideID)
	if err != nil {
		if errors.Is(err, repository.ErrGuideNotFound) {
			response.Error(c, errcode.ErrGuideNotFound, errcode.GetMsg(errcode.ErrGuideNotFound))
			return
		}
		ctrl.handleError(c, err)
		return
	}
	response.Success(c, detail)
}

// handleError 处理内容服务通用错误
func (ctrl *ContentController) handleError(c *gin.Context, err error) {
	zap.L().Error("内容服务异常",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.Error(err),
	)
	response.Error(c, errcode.ErrInternal, errcode.GetMsg(errcode.ErrInternal))
}
