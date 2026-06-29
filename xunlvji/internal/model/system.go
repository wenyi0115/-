package model

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// Role 角色表
type Role struct {
	ID          uint64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"column:name;type:varchar(64);not null" json:"name"`                       // 角色名称
	Code        string          `gorm:"column:code;type:varchar(32);not null;uniqueIndex:uk_code" json:"code"`   // 角色编码(唯一)
	Description *string         `gorm:"column:description;type:varchar(256)" json:"description,omitempty"`      // 角色描述
	Permissions json.RawMessage `gorm:"column:permissions;type:json" json:"permissions,omitempty"`               // 权限列表(JSON数组)
	CreatedAt   time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`                      // 创建时间
}

// TableName 指定表名
func (Role) TableName() string {
	return "role"
}

// BadgeDefinition 勋章定义表
type BadgeDefinition struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code           string    `gorm:"column:code;type:varchar(64);not null;uniqueIndex:uk_code" json:"code"`   // 勋章编码(唯一)
	Name           string    `gorm:"column:name;type:varchar(64);not null" json:"name"`                       // 勋章名称
	Description    *string   `gorm:"column:description;type:varchar(256)" json:"description,omitempty"`       // 勋章描述
	Icon           string    `gorm:"column:icon;type:varchar(512);not null" json:"icon"`                      // 勋章图标URL
	Category       string    `gorm:"column:category;type:varchar(32);not null;index:idx_category" json:"category"` // 分类:checkin/social/creator/system
	ConditionType  string    `gorm:"column:condition_type;type:varchar(32);not null" json:"condition_type"`   // 条件类型:city_count/checkin_count/follower_count等
	ConditionValue string    `gorm:"column:condition_value;type:varchar(64);not null" json:"condition_value"` // 条件值
	SortOrder      int       `gorm:"column:sort_order;not null;default:1" json:"sort_order"`                  // 排序序号
	Status         int8      `gorm:"column:status;not null;default:1;index:idx_status" json:"status"`         // 0-禁用 1-启用
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                       // 创建时间
}

// TableName 指定表名
func (BadgeDefinition) TableName() string {
	return "badge_definition"
}

// AuditLog 审计日志表
type AuditLog struct {
	ID         uint64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     *uint64         `gorm:"column:user_id;index:idx_user_id" json:"user_id,omitempty"`               // 操作用户ID
	Action     string          `gorm:"column:action;type:varchar(64);not null;index:idx_action" json:"action"`  // 操作类型
	TargetType string          `gorm:"column:target_type;type:varchar(32);not null;index:idx_target,priority:1" json:"target_type"` // 目标类型
	TargetID   *uint64         `gorm:"column:target_id;index:idx_target,priority:2" json:"target_id,omitempty"` // 目标ID
	Detail     json.RawMessage `gorm:"column:detail;type:json" json:"detail,omitempty"`                          // 操作详情
	IP         *string         `gorm:"column:ip;type:varchar(45)" json:"ip,omitempty"`                           // IP地址
	CreatedAt  time.Time       `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"`  // 创建时间
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_log"
}

// Report 举报表
type Report struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ReporterID  uint64     `gorm:"column:reporter_id;not null;index:idx_reporter_id" json:"reporter_id"`                // 举报人ID
	TargetType  string     `gorm:"column:target_type;type:varchar(20);not null;index:idx_target,priority:1" json:"target_type"` // 举报对象类型:checkin/comment/place/user
	TargetID    uint64     `gorm:"column:target_id;not null;index:idx_target,priority:2" json:"target_id"`              // 举报对象ID
	Reason      string     `gorm:"column:reason;type:varchar(50);not null" json:"reason"`                              // 举报原因:spam/inappropriate/wrong_info/copyright/other
	Description *string    `gorm:"column:description;type:text" json:"description,omitempty"`                          // 详细描述
	Status      string     `gorm:"column:status;type:varchar(20);not null;default:pending;index:idx_status" json:"status"` // 状态:pending-待处理 resolved-已处理 rejected-已驳回
	HandledBy   *uint64    `gorm:"column:handled_by" json:"handled_by,omitempty"`                                       // 处理人ID(管理员)
	HandledAt   *time.Time `gorm:"column:handled_at" json:"handled_at,omitempty"`                                       // 处理时间
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"`             // 创建时间
}

// TableName 指定表名
func (Report) TableName() string {
	return "report"
}

// City 城市信息表
type City struct {
	ID          uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string           `gorm:"column:name;type:varchar(50);not null;uniqueIndex:uk_name" json:"name"` // 城市名称
	ProvinceID  *uint64          `gorm:"column:province_id;index:idx_province_id" json:"province_id,omitempty"` // 省份ID(预留)
	Province    *string          `gorm:"column:province;type:varchar(50)" json:"province,omitempty"`            // 省份名称(冗余,便于展示)
	Latitude    *decimal.Decimal `gorm:"column:latitude;type:decimal(10,7)" json:"latitude,omitempty"`          // 纬度
	Longitude   *decimal.Decimal `gorm:"column:longitude;type:decimal(10,7)" json:"longitude,omitempty"`        // 经度
	BestSeason  *string          `gorm:"column:best_season;type:varchar(100)" json:"best_season,omitempty"`      // 最佳旅行季节
	Description *string          `gorm:"column:description;type:text" json:"description,omitempty"`             // 城市简介
	CoverImage  *string          `gorm:"column:cover_image;type:varchar(500)" json:"cover_image,omitempty"`     // 封面图URL
	IsHot       bool             `gorm:"column:is_hot;not null;default:0;index:idx_is_hot" json:"is_hot"`       // 是否热门城市:0-否 1-是
	SortOrder   int              `gorm:"column:sort_order;not null;default:0;index:idx_sort_order" json:"sort_order"` // 排序序号
	Status      int8             `gorm:"column:status;not null;default:1;index:idx_status" json:"status"`       // 0-禁用 1-启用
	CreatedAt   time.Time        `gorm:"column:created_at;autoCreateTime" json:"created_at"`                     // 创建时间
}

// TableName 指定表名
func (City) TableName() string {
	return "city"
}

// ChatGroup 群聊表
type ChatGroup struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(64);not null" json:"name"`            // 群名称
	Avatar      *string   `gorm:"column:avatar;type:varchar(512)" json:"avatar,omitempty"`     // 群头像
	OwnerID     uint64    `gorm:"column:owner_id;not null;index:idx_owner_id" json:"owner_id"` // 群主ID
	MemberCount int       `gorm:"column:member_count;not null;default:1" json:"member_count"`  // 成员数
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`          // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`          // 更新时间
}

// TableName 指定表名
func (ChatGroup) TableName() string {
	return "chat_group"
}

// GroupMember 群成员表
type GroupMember struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GroupID   uint64    `gorm:"column:group_id;not null;uniqueIndex:uk_group_user,priority:1" json:"group_id"` // 群ID
	UserID    uint64    `gorm:"column:user_id;not null;uniqueIndex:uk_group_user,priority:2;index:idx_user_id" json:"user_id"` // 用户ID
	Role      string    `gorm:"column:role;type:varchar(16);not null;default:member" json:"role"`             // 角色:admin/member
	JoinedAt  time.Time `gorm:"column:joined_at;not null;default:CURRENT_TIMESTAMP" json:"joined_at"`        // 加入时间
}

// TableName 指定表名
func (GroupMember) TableName() string {
	return "group_member"
}
