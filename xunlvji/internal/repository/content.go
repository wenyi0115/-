package repository

import (
	"context"
	"errors"

	"github.com/xunlvji/xunlvji/internal/model"
	"gorm.io/gorm"
)

// 内容查询错误
var (
	ErrPlaceNotFound   = errors.New("place not found")
	ErrRouteNotFound   = errors.New("route not found")
	ErrGuideNotFound   = errors.New("guide not found")
	ErrCheckinNotFound = errors.New("checkin not found")
)

// PlaceQuery 地点查询条件
type PlaceQuery struct {
	Category string // 分类筛选
	City     string // 城市筛选
	Sort     string // 排序：newest/popular/rating
	// 附近查询
	HasLocation bool
	Longitude   float64
	Latitude    float64
	Radius      int // 米
}

// PlaceRepository 地点仓库
type PlaceRepository interface {
	List(ctx context.Context, q PlaceQuery, offset, limit int) ([]model.Place, int64, error)
	GetByID(ctx context.Context, id uint64) (*model.Place, error)
}

// RouteQuery 路线查询条件
type RouteQuery struct {
	Category string // 分类：half_day/one_day/two_day/three_day_plus
	City     string // 城市筛选（关联点位城市，这里简化按无城市过滤，路线表无city字段）
	Sort     string // hot/newest/rating
}

// RouteRepository 路线仓库
type RouteRepository interface {
	List(ctx context.Context, q RouteQuery, offset, limit int) ([]model.Route, int64, error)
	GetByID(ctx context.Context, id uint64) (*model.Route, error)
	GetPoints(ctx context.Context, routeID uint64) ([]model.RoutePoint, error)
	GetCreator(ctx context.Context, userID uint64) (*model.User, error)
}

// GuideQuery 攻略查询条件
type GuideQuery struct {
	Category string // 分类
	City     string // 城市
	Sort     string // newest/popular
}

// GuideRepository 攻略仓库
type GuideRepository interface {
	List(ctx context.Context, q GuideQuery, offset, limit int) ([]model.Guide, int64, error)
	GetByID(ctx context.Context, id uint64) (*model.Guide, error)
}

// CheckinRepository 打卡记录仓库
type CheckinRepository interface {
	ListByPlace(ctx context.Context, placeID uint64, status string, offset, limit int) ([]model.CheckinRecord, int64, error)
	GetByID(ctx context.Context, id uint64) (*model.CheckinRecord, error)
}

// === 地点仓库实现 ===

type placeRepo struct {
	db *gorm.DB
}

// NewPlaceRepository 创建地点仓库
func NewPlaceRepository(db *gorm.DB) PlaceRepository {
	return &placeRepo{db: db}
}

func (r *placeRepo) List(ctx context.Context, q PlaceQuery, offset, limit int) ([]model.Place, int64, error) {
	var list []model.Place
	var total int64

	tx := r.db.WithContext(ctx).Model(&model.Place{}).Where("status = ?", 1)
	if q.Category != "" {
		tx = tx.Where("category = ?", q.Category)
	}
	if q.City != "" {
		tx = tx.Where("city = ?", q.City)
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	order := "created_at DESC"
	switch q.Sort {
	case "popular":
		order = "checkin_count DESC"
	case "rating":
		order = "rating DESC"
	}
	if err := tx.Order(order).Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *placeRepo) GetByID(ctx context.Context, id uint64) (*model.Place, error) {
	var p model.Place
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlaceNotFound
		}
		return nil, err
	}
	return &p, nil
}

// === 路线仓库实现 ===

type routeRepo struct {
	db *gorm.DB
}

// NewRouteRepository 创建路线仓库
func NewRouteRepository(db *gorm.DB) RouteRepository {
	return &routeRepo{db: db}
}

func (r *routeRepo) List(ctx context.Context, q RouteQuery, offset, limit int) ([]model.Route, int64, error) {
	var list []model.Route
	var total int64

	tx := r.db.WithContext(ctx).Model(&model.Route{}).Where("status = ?", 1)
	if q.Category != "" {
		tx = tx.Where("category = ?", q.Category)
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "usage_count DESC, created_at DESC"
	switch q.Sort {
	case "newest":
		order = "created_at DESC"
	case "rating":
		order = "rating DESC"
	}
	if err := tx.Order(order).Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *routeRepo) GetByID(ctx context.Context, id uint64) (*model.Route, error) {
	var rt model.Route
	if err := r.db.WithContext(ctx).First(&rt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRouteNotFound
		}
		return nil, err
	}
	return &rt, nil
}

func (r *routeRepo) GetPoints(ctx context.Context, routeID uint64) ([]model.RoutePoint, error) {
	var points []model.RoutePoint
	err := r.db.WithContext(ctx).
		Where("route_id = ?", routeID).
		Order("sort_order ASC").
		Find(&points).Error
	return points, err
}

func (r *routeRepo) GetCreator(ctx context.Context, userID uint64) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).First(&u, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// === 攻略仓库实现 ===

type guideRepo struct {
	db *gorm.DB
}

// NewGuideRepository 创建攻略仓库
func NewGuideRepository(db *gorm.DB) GuideRepository {
	return &guideRepo{db: db}
}

func (r *guideRepo) List(ctx context.Context, q GuideQuery, offset, limit int) ([]model.Guide, int64, error) {
	var list []model.Guide
	var total int64

	tx := r.db.WithContext(ctx).Model(&model.Guide{}).Where("status = ?", 1)
	if q.Category != "" {
		tx = tx.Where("category = ?", q.Category)
	}
	if q.City != "" {
		tx = tx.Where("city = ?", q.City)
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at DESC"
	if q.Sort == "popular" {
		order = "verified_count DESC"
	}
	if err := tx.Order(order).Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *guideRepo) GetByID(ctx context.Context, id uint64) (*model.Guide, error) {
	var g model.Guide
	if err := r.db.WithContext(ctx).First(&g, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGuideNotFound
		}
		return nil, err
	}
	return &g, nil
}

// === 打卡记录仓库实现 ===

type checkinRepo struct {
	db *gorm.DB
}

// NewCheckinRepository 创建打卡记录仓库
func NewCheckinRepository(db *gorm.DB) CheckinRepository {
	return &checkinRepo{db: db}
}

func (r *checkinRepo) ListByPlace(ctx context.Context, placeID uint64, status string, offset, limit int) ([]model.CheckinRecord, int64, error) {
	var list []model.CheckinRecord
	var total int64

	tx := r.db.WithContext(ctx).Model(&model.CheckinRecord{}).Where("place_id = ?", placeID)
	if status != "" {
		tx = tx.Where("status = ?", status)
	} else {
		tx = tx.Where("status = ?", "approved")
	}
	tx = tx.Where("visibility = ?", "public")

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("checkin_time DESC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *checkinRepo) GetByID(ctx context.Context, id uint64) (*model.CheckinRecord, error) {
	var c model.CheckinRecord
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCheckinNotFound
		}
		return nil, err
	}
	return &c, nil
}
