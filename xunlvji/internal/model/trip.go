package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Trip 行程表
type Trip struct {
	ID         uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     uint64           `gorm:"column:user_id;not null;index:idx_user_id;index:idx_user_status,priority:1" json:"user_id"` // 用户ID
	Name       string           `gorm:"column:name;type:varchar(128);not null" json:"name"`                                       // 行程名称
	TripDate   time.Time        `gorm:"column:trip_date;type:date;not null;index:idx_trip_date" json:"trip_date"`                  // 出行日期(开始日期)
	EndDate    *time.Time       `gorm:"column:end_date;type:date" json:"end_date,omitempty"`                                      // 结束日期(多日行程)
	Budget     *decimal.Decimal `gorm:"column:budget;type:decimal(10,2)" json:"budget,omitempty"`                                 // 行程预算,单位:元
	Visibility string           `gorm:"column:visibility;type:varchar(16);not null;default:private" json:"visibility"`            // 可见性:public(公开)/private(仅自己可见)/friends(仅好友可见)
	Status     string           `gorm:"column:status;type:varchar(16);not null;default:pending;index:idx_user_status,priority:2" json:"status"` // pending-待出发 active-进行中 completed-已完成
	RouteID    *uint64          `gorm:"column:route_id;index:idx_route_id" json:"route_id,omitempty"`                              // 关联路线ID(可选)
	CreatedAt  time.Time        `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"`                   // 创建时间
	UpdatedAt  time.Time        `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                        // 更新时间
}

// TableName 指定表名
func (Trip) TableName() string {
	return "trip"
}

// TripPoint 行程点位表
type TripPoint struct {
	ID            uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TripID        uint64           `gorm:"column:trip_id;not null;index:idx_trip_id;index:idx_trip_sort,priority:1" json:"trip_id"` // 行程ID
	SortOrder     int              `gorm:"column:sort_order;not null;default:1;index:idx_trip_sort,priority:2" json:"sort_order"`   // 排序序号(从1开始)
	PlaceID       *uint64          `gorm:"column:place_id" json:"place_id,omitempty"`                                                // 地点ID
	Name          string           `gorm:"column:name;type:varchar(128);not null" json:"name"`                                       // 点位名称
	PointTime     *string          `gorm:"column:point_time;type:time" json:"point_time,omitempty"`                                  // 安排时间(HH:MM:SS)
	StayDuration  *int             `gorm:"column:stay_duration" json:"stay_duration,omitempty"`                                      // 停留时长(分钟)
	Cost          *decimal.Decimal `gorm:"column:cost;type:decimal(10,2)" json:"cost,omitempty"`                                     // 预计消费(元)
	Category      *string          `gorm:"column:category;type:varchar(32)" json:"category,omitempty"`                               // 分类
	Transport     *string          `gorm:"column:transport;type:varchar(32)" json:"transport,omitempty"`                             // 交通方式
	CheckinStatus string           `gorm:"column:checkin_status;type:varchar(16);not null;default:pending;index:idx_checkin_status" json:"checkin_status"` // pending-待打卡 checked-已打卡 skipped-已跳过
	CreatedAt     time.Time        `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                        // 创建时间
	UpdatedAt     time.Time        `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                        // 更新时间
}

// TableName 指定表名
func (TripPoint) TableName() string {
	return "trip_point"
}

// TripChecklist 行程准备清单表
type TripChecklist struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TripID        uint64    `gorm:"column:trip_id;not null;index:idx_trip_id;index:idx_trip_sort,priority:1" json:"trip_id"`    // 行程ID
	ItemName      string    `gorm:"column:item_name;type:varchar(100);not null" json:"item_name"`                                // 清单项名称
	ItemCategory  string    `gorm:"column:item_category;type:varchar(20);not null;default:other" json:"item_category"`          // 分类:document/electronic/clothing/other
	IsChecked     bool      `gorm:"column:is_checked;not null;default:0" json:"is_checked"`                                       // 是否已勾选:0-否 1-是
	IsAIGenerated bool      `gorm:"column:is_ai_generated;not null;default:0" json:"is_ai_generated"`                             // 是否AI生成:0-否 1-是
	SortOrder     int       `gorm:"column:sort_order;not null;default:1;index:idx_trip_sort,priority:2" json:"sort_order"`       // 排序序号(从1开始)
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                           // 创建时间
}

// TableName 指定表名
func (TripChecklist) TableName() string {
	return "trip_checklist"
}

// TripExpense 行程账单表
// 注: 原 DDL 标注 "P1: P0阶段不创建,P1开发时建表"。本结构体已就绪,实际建表由迁移层在 P1 阶段处理。
type TripExpense struct {
	ID          uint64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TripID      uint64          `gorm:"column:trip_id;not null;index:idx_trip_id;index:idx_trip_type,priority:1" json:"trip_id"` // 行程ID
	PlaceID     *uint64         `gorm:"column:place_id;index:idx_place_id" json:"place_id,omitempty"`                            // 关联地点ID(可空,手动添加不关联)
	CheckinID   *uint64         `gorm:"column:checkin_id;index:idx_checkin_id" json:"checkin_id,omitempty"`                      // 关联打卡记录ID(可空,从打卡自动来)
	ExpenseType string          `gorm:"column:expense_type;type:varchar(20);not null;index:idx_trip_type,priority:2" json:"expense_type"` // 消费类型:food/transport/ticket/hotel/other
	Amount      decimal.Decimal `gorm:"column:amount;type:decimal(8,2);not null" json:"amount"`                                  // 金额(元)
	Description *string         `gorm:"column:description;type:varchar(200)" json:"description,omitempty"`                       // 描述
	IsAuto      bool            `gorm:"column:is_auto;not null;default:0" json:"is_auto"`                                         // 是否自动生成:0-否 1-是(从打卡记录来)
	CreatedAt   time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                       // 创建时间
}

// TableName 指定表名
func (TripExpense) TableName() string {
	return "trip_expense"
}
