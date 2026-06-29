package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/xunlvji/xunlvji/internal/model"
	"github.com/xunlvji/xunlvji/internal/repository"
	"github.com/xunlvji/xunlvji/pkg/util"
	"gorm.io/gorm"
)

// === 行程状态 ===

// 行程状态
const (
	TripStatusPending   = "pending"
	TripStatusActive    = "active"
	TripStatusCompleted = "completed"
)

// 行程状态中文名
var tripStatusNames = map[string]string{
	"pending":   "待出发",
	"active":    "进行中",
	"completed": "已完成",
}

// TripStatusName 获取行程状态中文名
func TripStatusName(code string) string {
	return tripStatusNames[code]
}

// 点位打卡状态
const (
	PointStatusPending = "pending"
	PointStatusChecked = "checked"
	PointStatusSkipped = "skipped"
)

// 点位打卡状态中文名
var pointStatusNames = map[string]string{
	"pending": "待打卡",
	"checked": "已打卡",
	"skipped": "已跳过",
}

// PointStatusName 获取点位状态中文名
func PointStatusName(code string) string {
	return pointStatusNames[code]
}

// === DTO 定义 ===

// TripCreateInput 创建行程输入
type TripCreateInput struct {
	Name     string             `json:"name"`
	TripDate string             `json:"trip_date"` // YYYY-MM-DD
	EndDate  string             `json:"end_date"`  // YYYY-MM-DD，可选
	RouteID  *uint64            `json:"route_id"`
	Points   []TripPointInput   `json:"points"`
}

// TripPointInput 点位输入
type TripPointInput struct {
	PlaceID      *uint64 `json:"place_id"`
	Name         string  `json:"name"`
	PointTime    string  `json:"point_time"`
	StayDuration *int    `json:"stay_duration"`
	Cost         *float64 `json:"cost"`
	Category     string  `json:"category"`
	Transport    string  `json:"transport"`
}

// TripListItemDTO 行程列表项
type TripListItemDTO struct {
	ID             uint64  `json:"id"`
	Name           string  `json:"name"`
	TripDate       string  `json:"trip_date"`
	Status         string  `json:"status"`
	StatusName     string  `json:"status_name"`
	RouteID        *uint64 `json:"route_id,omitempty"`
	RouteTitle     string  `json:"route_title,omitempty"`
	PointCount     int     `json:"point_count"`
	CheckedCount   int     `json:"checked_count"`
	SkippedCount   int     `json:"skipped_count"`
	Progress       float64 `json:"progress,omitempty"`
	CountdownDays  int     `json:"countdown_days,omitempty"`
	CreatedAt      string  `json:"created_at"`
}

// TripPointDTO 行程点位
type TripPointDTO struct {
	ID               uint64      `json:"id"`
	SortOrder        int         `json:"sort_order"`
	PlaceID          *uint64     `json:"place_id,omitempty"`
	Name             string      `json:"name"`
	PointTime        string      `json:"point_time,omitempty"`
	StayDuration     *int        `json:"stay_duration,omitempty"`
	Cost             *float64    `json:"cost,omitempty"`
	Category         string      `json:"category,omitempty"`
	CategoryName     string      `json:"category_name,omitempty"`
	Transport        string      `json:"transport,omitempty"`
	CheckinStatus    string      `json:"checkin_status"`
	CheckinStatusName string     `json:"checkin_status_name"`
	CheckinRecordID  *uint64     `json:"checkin_record_id,omitempty"`
	Place            *PlaceBriefDTO `json:"place,omitempty"`
}

// PlaceBriefDTO 地点简要信息
type PlaceBriefDTO struct {
	ID           uint64  `json:"id"`
	Name         string  `json:"name"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	Address      string  `json:"address,omitempty"`
	OpeningHours string  `json:"opening_hours,omitempty"`
	IsOpen       bool    `json:"is_open"`
}

// TripDetailDTO 行程详情
type TripDetailDTO struct {
	ID               uint64          `json:"id"`
	Name             string          `json:"name"`
	TripDate         string          `json:"trip_date"`
	EndDate          string          `json:"end_date,omitempty"`
	Status           string          `json:"status"`
	StatusName       string          `json:"status_name"`
	RouteID          *uint64         `json:"route_id,omitempty"`
	RouteTitle       string          `json:"route_title,omitempty"`
	PointCount       int             `json:"point_count"`
	CheckedCount     int             `json:"checked_count"`
	SkippedCount     int             `json:"skipped_count"`
	Progress         float64         `json:"progress"`
	ElapsedMinutes   int             `json:"elapsed_minutes,omitempty"`
	EstimatedMinutes int             `json:"estimated_minutes,omitempty"`
	TotalCost        float64         `json:"total_cost,omitempty"`
	CurrentPoint     *TripPointDTO   `json:"current_point,omitempty"`
	Points           []TripPointDTO  `json:"points"`
	CreatedAt        string          `json:"created_at"`
}

// TripCreateResultDTO 创建行程结果
type TripCreateResultDTO struct {
	ID         uint64  `json:"id"`
	Name       string  `json:"name"`
	TripDate   string  `json:"trip_date"`
	EndDate    string  `json:"end_date,omitempty"`
	Status     string  `json:"status"`
	RouteID    *uint64 `json:"route_id,omitempty"`
	PointCount int     `json:"point_count"`
	CreatedAt  string  `json:"created_at"`
}

// === 行程服务 ===

// TripService 行程服务
type TripService struct {
	tripRepo  repository.TripRepository
	routeRepo repository.RouteRepository
	placeRepo repository.PlaceRepository
}

// NewTripService 创建行程服务
func NewTripService(tripRepo repository.TripRepository, routeRepo repository.RouteRepository, placeRepo repository.PlaceRepository) *TripService {
	return &TripService{
		tripRepo:  tripRepo,
		routeRepo: routeRepo,
		placeRepo: placeRepo,
	}
}

// Create 创建行程
func (s *TripService) Create(ctx context.Context, userID uint64, input TripCreateInput) (*TripCreateResultDTO, error) {
	// 参数校验
	if input.Name == "" || len(input.Name) > 50 {
		return nil, ErrInvalidTripName
	}
	tripDate, err := time.Parse("2006-01-02", input.TripDate)
	if err != nil {
		return nil, ErrInvalidTripDate
	}
	var endDate time.Time
	if input.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", input.EndDate)
		if err != nil {
			return nil, ErrInvalidTripDate
		}
		if endDate.Before(tripDate) {
			return nil, ErrInvalidTripDate
		}
	} else {
		endDate = tripDate
	}

	// 频率限制：每天最多5个行程
	todayCount, err := s.tripRepo.CountByUserToday(ctx, userID, time.Now())
	if err != nil {
		return nil, err
	}
	if todayCount >= 5 {
		return nil, ErrTripCreateLimit
	}

	// 构建行程
	trip := &model.Trip{
		UserID:   userID,
		Name:     input.Name,
		TripDate: tripDate,
		EndDate:  &endDate,
		Status:   TripStatusPending,
	}
	if input.RouteID != nil {
		trip.RouteID = input.RouteID
	}

	// 构建点位
	var points []model.TripPoint
	if len(input.Points) > 0 {
		points = make([]model.TripPoint, 0, len(input.Points))
		for i, p := range input.Points {
			tp := model.TripPoint{
				SortOrder:     i + 1,
				PlaceID:       p.PlaceID,
				Name:          p.Name,
				StayDuration:  p.StayDuration,
				CheckinStatus: PointStatusPending,
			}
			if p.PointTime != "" {
				tp.PointTime = &p.PointTime
			}
			if p.Cost != nil {
				c := decimal.NewFromFloat(*p.Cost)
				tp.Cost = &c
			}
			if p.Category != "" {
				tp.Category = &p.Category
			}
			if p.Transport != "" {
				tp.Transport = &p.Transport
			}
			points = append(points, tp)
		}
	} else if input.RouteID != nil {
		// 基于路线创建：从路线复制点位
		routePoints, err := s.routeRepo.GetPoints(ctx, *input.RouteID)
		if err != nil {
			return nil, err
		}
		points = make([]model.TripPoint, 0, len(routePoints))
		for i, rp := range routePoints {
			tp := model.TripPoint{
				SortOrder:     i + 1,
				PlaceID:       rp.PlaceID,
				Name:          rp.Name,
				PointTime:     rp.PointTime,
				StayDuration:  rp.StayDuration,
				Cost:          rp.Cost,
				Category:      rp.Category,
				Transport:     rp.Transport,
				CheckinStatus: PointStatusPending,
			}
			points = append(points, tp)
		}
	}

	// 事务创建行程 + 点位
	err = s.tripRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := s.tripRepo.Create(ctx, tx, trip); err != nil {
			return err
		}
		// 设置点位的 TripID
		for i := range points {
			points[i].TripID = trip.ID
		}
		if len(points) > 0 {
			if err := s.tripRepo.CreatePoints(ctx, tx, points); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := &TripCreateResultDTO{
		ID:         trip.ID,
		Name:       trip.Name,
		TripDate:   trip.TripDate.Format("2006-01-02"),
		Status:     trip.Status,
		RouteID:    trip.RouteID,
		PointCount: len(points),
		CreatedAt:  trip.CreatedAt.Format(time.RFC3339),
	}
	if trip.EndDate != nil {
		result.EndDate = trip.EndDate.Format("2006-01-02")
	}
	return result, nil
}

// List 行程列表
func (s *TripService) List(ctx context.Context, userID uint64, status string, p util.Pagination) (*util.PageResult, error) {
	list, total, err := s.tripRepo.List(ctx, repository.TripQuery{
		UserID: userID,
		Status: status,
	}, p.Offset, p.PageSize)
	if err != nil {
		return nil, err
	}

	items := make([]TripListItemDTO, 0, len(list))
	for _, t := range list {
		item := TripListItemDTO{
			ID:         t.ID,
			Name:       t.Name,
			TripDate:   t.TripDate.Format("2006-01-02"),
			Status:     t.Status,
			StatusName: TripStatusName(t.Status),
			RouteID:    t.RouteID,
			CreatedAt:  t.CreatedAt.Format(time.RFC3339),
		}
		// 补充点位统计
		points, _ := s.tripRepo.ListPoints(ctx, t.ID)
		item.PointCount = len(points)
		for _, pt := range points {
			switch pt.CheckinStatus {
			case PointStatusChecked:
				item.CheckedCount++
			case PointStatusSkipped:
				item.SkippedCount++
			}
		}
		// 进度计算
		validCount := item.PointCount - item.SkippedCount
		if validCount > 0 {
			item.Progress = roundTo1Decimal(float64(item.CheckedCount) / float64(validCount) * 100)
		}
		// 倒计时（仅 pending）
		if t.Status == TripStatusPending {
			days := int(t.TripDate.Sub(time.Now().Truncate(24*time.Hour)).Hours() / 24)
			if days < 0 {
				days = 0
			}
			item.CountdownDays = days
		}
		// 路线标题
		if t.RouteID != nil {
			if route, err := s.routeRepo.GetByID(ctx, *t.RouteID); err == nil {
				item.RouteTitle = route.Title
			}
		}
		items = append(items, item)
	}
	return util.NewPageResult(items, total, p), nil
}

// GetDetail 行程详情
func (s *TripService) GetDetail(ctx context.Context, tripID uint64) (*TripDetailDTO, error) {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return nil, err
	}

	points, err := s.tripRepo.ListPoints(ctx, tripID)
	if err != nil {
		return nil, err
	}

	dto := &TripDetailDTO{
		ID:         trip.ID,
		Name:       trip.Name,
		TripDate:   trip.TripDate.Format("2006-01-02"),
		Status:     trip.Status,
		StatusName: TripStatusName(trip.Status),
		RouteID:    trip.RouteID,
		PointCount: len(points),
		CreatedAt:  trip.CreatedAt.Format(time.RFC3339),
	}
	if trip.EndDate != nil {
		dto.EndDate = trip.EndDate.Format("2006-01-02")
	}
	if trip.RouteID != nil {
		if route, err := s.routeRepo.GetByID(ctx, *trip.RouteID); err == nil {
			dto.RouteTitle = route.Title
		}
	}

	// 转换点位
	totalCost := decimal.Zero
	estimatedMinutes := 0
	pointDTOs := make([]TripPointDTO, 0, len(points))
	for _, p := range points {
		pt := TripPointDTO{
			ID:                p.ID,
			SortOrder:         p.SortOrder,
			PlaceID:           p.PlaceID,
			Name:              p.Name,
			StayDuration:      p.StayDuration,
			CheckinStatus:     p.CheckinStatus,
			CheckinStatusName: PointStatusName(p.CheckinStatus),
		}
		if p.PointTime != nil {
			pt.PointTime = *p.PointTime
		}
		if p.Cost != nil {
			c := p.Cost.InexactFloat64()
			pt.Cost = &c
			totalCost = totalCost.Add(*p.Cost)
		}
		if p.Category != nil {
			pt.Category = *p.Category
			pt.CategoryName = PlaceCategoryName(*p.Category)
		}
		if p.Transport != nil {
			pt.Transport = *p.Transport
		}
		if p.StayDuration != nil {
			estimatedMinutes += *p.StayDuration
		}
		// 补充地点信息
		if p.PlaceID != nil {
			if place, err := s.placeRepo.GetByID(ctx, *p.PlaceID); err == nil {
				pt.Place = &PlaceBriefDTO{
					ID:           place.ID,
					Name:         place.Name,
					Longitude:    place.Longitude.InexactFloat64(),
					Latitude:     place.Latitude.InexactFloat64(),
					Address:      ptrToStr(place.Address),
					OpeningHours: ptrToStr(place.OpeningHours),
					IsOpen:       isOpenNow(ptrToStr(place.OpeningHours)),
				}
			}
		}
		pointDTOs = append(pointDTOs, pt)

		if p.CheckinStatus == PointStatusChecked {
			dto.CheckedCount++
		}
		if p.CheckinStatus == PointStatusSkipped {
			dto.SkippedCount++
		}
	}
	dto.Points = pointDTOs
	dto.TotalCost = totalCost.InexactFloat64()
	dto.EstimatedMinutes = estimatedMinutes

	// 进度
	validCount := dto.PointCount - dto.SkippedCount
	if validCount > 0 {
		dto.Progress = roundTo1Decimal(float64(dto.CheckedCount) / float64(validCount) * 100)
	}

	// 当前点位（第一个 pending）
	for i := range pointDTOs {
		if pointDTOs[i].CheckinStatus == PointStatusPending {
			dto.CurrentPoint = &pointDTOs[i]
			break
		}
	}

	return dto, nil
}

// Start 开始行程
func (s *TripService) Start(ctx context.Context, userID, tripID uint64) error {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip.UserID != userID {
		return ErrTripNoPermission
	}
	if trip.Status != TripStatusPending {
		return repository.ErrTripStatusInvalid
	}
	return s.tripRepo.UpdateStatus(ctx, tripID, TripStatusActive)
}

// Complete 结束行程
func (s *TripService) Complete(ctx context.Context, userID, tripID uint64) error {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip.UserID != userID {
		return ErrTripNoPermission
	}
	if trip.Status != TripStatusActive {
		return repository.ErrTripStatusInvalid
	}
	return s.tripRepo.UpdateStatus(ctx, tripID, TripStatusCompleted)
}

// Delete 删除行程
func (s *TripService) Delete(ctx context.Context, userID, tripID uint64) error {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip.UserID != userID {
		return ErrTripNoPermission
	}
	return s.tripRepo.Delete(ctx, tripID)
}

// Update 更新行程基本信息
func (s *TripService) Update(ctx context.Context, userID, tripID uint64, name, tripDate, endDate string) error {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip.UserID != userID {
		return ErrTripNoPermission
	}
	if name != "" {
		if len(name) > 50 {
			return ErrInvalidTripName
		}
		trip.Name = name
	}
	if tripDate != "" {
		t, err := time.Parse("2006-01-02", tripDate)
		if err != nil {
			return ErrInvalidTripDate
		}
		trip.TripDate = t
	}
	if endDate != "" {
		t, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return ErrInvalidTripDate
		}
		trip.EndDate = &t
	}
	return s.tripRepo.Update(ctx, trip)
}

// AddPoint 添加点位到行程
func (s *TripService) AddPoint(ctx context.Context, userID, tripID uint64, input TripPointInput) error {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip.UserID != userID {
		return ErrTripNoPermission
	}
	// 获取当前最大序号
	points, err := s.tripRepo.ListPoints(ctx, tripID)
	if err != nil {
		return err
	}
	maxOrder := 0
	for _, p := range points {
		if p.SortOrder > maxOrder {
			maxOrder = p.SortOrder
		}
	}
	point := &model.TripPoint{
		TripID:        tripID,
		SortOrder:     maxOrder + 1,
		PlaceID:       input.PlaceID,
		Name:          input.Name,
		StayDuration:  input.StayDuration,
		CheckinStatus: PointStatusPending,
	}
	if input.PointTime != "" {
		point.PointTime = &input.PointTime
	}
	if input.Cost != nil {
		c := decimal.NewFromFloat(*input.Cost)
		point.Cost = &c
	}
	if input.Category != "" {
		point.Category = &input.Category
	}
	if input.Transport != "" {
		point.Transport = &input.Transport
	}
	return s.tripRepo.AddPoint(ctx, point)
}

// DeletePoint 删除行程点位
func (s *TripService) DeletePoint(ctx context.Context, userID, tripID, pointID uint64) error {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip.UserID != userID {
		return ErrTripNoPermission
	}
	return s.tripRepo.DeletePoint(ctx, pointID)
}

// SkipPoint 跳过点位
func (s *TripService) SkipPoint(ctx context.Context, userID, tripID, pointID uint64) error {
	return s.updatePointStatus(ctx, userID, tripID, pointID, PointStatusSkipped, true)
}

// UnskipPoint 取消跳过点位
func (s *TripService) UnskipPoint(ctx context.Context, userID, tripID, pointID uint64) error {
	return s.updatePointStatus(ctx, userID, tripID, pointID, PointStatusPending, false)
}

func (s *TripService) updatePointStatus(ctx context.Context, userID, tripID, pointID uint64, status string, mustBePending bool) error {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip.UserID != userID {
		return ErrTripNoPermission
	}
	point, err := s.tripRepo.GetPointByID(ctx, pointID)
	if err != nil {
		return err
	}
	if point.TripID != tripID {
		return repository.ErrTripPointNotFound
	}
	if mustBePending && point.CheckinStatus != PointStatusPending {
		return ErrPointStatusInvalid
	}
	return s.tripRepo.UpdatePointStatus(ctx, pointID, status)
}

// 行程服务错误
var (
	ErrInvalidTripName  = errors.New("invalid trip name")
	ErrInvalidTripDate  = errors.New("invalid trip date")
	ErrTripCreateLimit  = errors.New("trip create limit exceeded")
	ErrTripNoPermission = errors.New("trip no permission")
	ErrPointStatusInvalid = errors.New("point status invalid")
)

// 占位：避免 fmt 未使用
var _ = fmt.Sprintf
