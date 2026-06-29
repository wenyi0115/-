package repository

import (
	"context"
	"errors"
	"time"

	"github.com/xunlvji/xunlvji/internal/model"
	"gorm.io/gorm"
)

// 行程仓库错误
var (
	ErrTripNotFound      = errors.New("trip not found")
	ErrTripPointNotFound = errors.New("trip point not found")
	ErrTripStatusInvalid = errors.New("invalid trip status transition")
)

// TripQuery 行程查询条件
type TripQuery struct {
	UserID uint64
	Status string // pending/active/completed，空=全部
}

// TripRepository 行程仓库
type TripRepository interface {
	Create(ctx context.Context, tx *gorm.DB, trip *model.Trip) error
	List(ctx context.Context, q TripQuery, offset, limit int) ([]model.Trip, int64, error)
	GetByID(ctx context.Context, id uint64) (*model.Trip, error)
	Update(ctx context.Context, trip *model.Trip) error
	UpdateStatus(ctx context.Context, id uint64, status string) error
	Delete(ctx context.Context, id uint64) error
	CountByUserToday(ctx context.Context, userID uint64, date time.Time) (int64, error)

	// 点位
	CreatePoints(ctx context.Context, tx *gorm.DB, points []model.TripPoint) error
	ListPoints(ctx context.Context, tripID uint64) ([]model.TripPoint, error)
	GetPointByID(ctx context.Context, pointID uint64) (*model.TripPoint, error)
	AddPoint(ctx context.Context, point *model.TripPoint) error
	DeletePoint(ctx context.Context, pointID uint64) error
	UpdatePointStatus(ctx context.Context, pointID uint64, status string) error
	CountPointsByStatus(ctx context.Context, tripID uint64, status string) (int64, error)

	// 事务
	Transaction(ctx context.Context, fc func(tx *gorm.DB) error) error
	DB() *gorm.DB
}

// === 行程仓库实现 ===

type tripRepo struct {
	db *gorm.DB
}

// NewTripRepository 创建行程仓库
func NewTripRepository(db *gorm.DB) TripRepository {
	return &tripRepo{db: db}
}

func (r *tripRepo) DB() *gorm.DB { return r.db }

func (r *tripRepo) Transaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fc)
}

func (r *tripRepo) Create(ctx context.Context, tx *gorm.DB, trip *model.Trip) error {
	db := r.db.WithContext(ctx)
	if tx != nil {
		db = tx
	}
	return db.Create(trip).Error
}

func (r *tripRepo) List(ctx context.Context, q TripQuery, offset, limit int) ([]model.Trip, int64, error) {
	var list []model.Trip
	var total int64
	tx := r.db.WithContext(ctx).Model(&model.Trip{}).Where("user_id = ?", q.UserID)
	if q.Status != "" {
		tx = tx.Where("status = ?", q.Status)
	}
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("trip_date DESC, created_at DESC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *tripRepo) GetByID(ctx context.Context, id uint64) (*model.Trip, error) {
	var t model.Trip
	if err := r.db.WithContext(ctx).First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *tripRepo) Update(ctx context.Context, trip *model.Trip) error {
	return r.db.WithContext(ctx).Save(trip).Error
}

func (r *tripRepo) UpdateStatus(ctx context.Context, id uint64, status string) error {
	return r.db.WithContext(ctx).Model(&model.Trip{}).Where("id = ?", id).Update("status", status).Error
}

func (r *tripRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删点位
		if err := tx.Where("trip_id = ?", id).Delete(&model.TripPoint{}).Error; err != nil {
			return err
		}
		// 删清单
		if err := tx.Where("trip_id = ?", id).Delete(&model.TripChecklist{}).Error; err != nil {
			return err
		}
		// 删行程
		return tx.Delete(&model.Trip{}, id).Error
	})
}

func (r *tripRepo) CountByUserToday(ctx context.Context, userID uint64, date time.Time) (int64, error) {
	var count int64
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)
	err := r.db.WithContext(ctx).Model(&model.Trip{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, start, end).
		Count(&count).Error
	return count, err
}

// === 点位操作 ===

func (r *tripRepo) CreatePoints(ctx context.Context, tx *gorm.DB, points []model.TripPoint) error {
	db := r.db.WithContext(ctx)
	if tx != nil {
		db = tx
	}
	if len(points) == 0 {
		return nil
	}
	return db.Create(&points).Error
}

func (r *tripRepo) ListPoints(ctx context.Context, tripID uint64) ([]model.TripPoint, error) {
	var points []model.TripPoint
	err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("sort_order ASC").
		Find(&points).Error
	return points, err
}

func (r *tripRepo) GetPointByID(ctx context.Context, pointID uint64) (*model.TripPoint, error) {
	var p model.TripPoint
	if err := r.db.WithContext(ctx).First(&p, pointID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripPointNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *tripRepo) AddPoint(ctx context.Context, point *model.TripPoint) error {
	return r.db.WithContext(ctx).Create(point).Error
}

func (r *tripRepo) DeletePoint(ctx context.Context, pointID uint64) error {
	return r.db.WithContext(ctx).Delete(&model.TripPoint{}, pointID).Error
}

func (r *tripRepo) UpdatePointStatus(ctx context.Context, pointID uint64, status string) error {
	return r.db.WithContext(ctx).Model(&model.TripPoint{}).
		Where("id = ?", pointID).
		Update("checkin_status", status).Error
}

func (r *tripRepo) CountPointsByStatus(ctx context.Context, tripID uint64, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.TripPoint{}).
		Where("trip_id = ? AND checkin_status = ?", tripID, status).
		Count(&count).Error
	return count, err
}
