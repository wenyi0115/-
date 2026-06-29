package model

import (
	"encoding/json"
	"time"
)

// AIChatSession AI对话会话表
// 注: 原 DDL 标注 "P1: P0阶段不创建,P1开发时建表"。本结构体已就绪,实际建表由迁移层在 P1 阶段处理。
type AIChatSession struct {
	ID        uint64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint64          `gorm:"column:user_id;not null;index:idx_user_id" json:"user_id"`     // 用户ID
	Title     string          `gorm:"column:title;type:varchar(128);not null" json:"title"`        // 会话标题
	Context   json.RawMessage `gorm:"column:context;type:json" json:"context,omitempty"`           // 上下文信息
	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`          // 更新时间
}

// TableName 指定表名
func (AIChatSession) TableName() string {
	return "ai_chat_session"
}

// AIChatMessage AI对话消息表
// 注: 原 DDL 标注 "P1: P0阶段不创建,P1开发时建表"。本结构体已就绪,实际建表由迁移层在 P1 阶段处理。
type AIChatMessage struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SessionID uint64    `gorm:"column:session_id;not null;index:idx_session_id" json:"session_id"` // 会话ID
	Role      string    `gorm:"column:role;type:varchar(16);not null" json:"role"`                 // 角色:user/assistant
	Content   string    `gorm:"column:content;type:text;not null" json:"content"`                  // 消息内容
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"` // 创建时间
}

// TableName 指定表名
func (AIChatMessage) TableName() string {
	return "ai_chat_message"
}
