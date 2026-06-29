package model

import "time"

// GuideQuestion 攻略提问表
type GuideQuestion struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"column:user_id;not null;index:idx_gq_user" json:"user_id"`                          // 提问用户ID
	Title     string    `gorm:"column:title;type:varchar(128);not null" json:"title"`                              // 问题标题
	Content   *string   `gorm:"column:content;type:text" json:"content,omitempty"`                                 // 问题详细描述
	CityID    *uint32   `gorm:"column:city_id;index:idx_gq_city" json:"city_id,omitempty"`                         // 关联城市ID
	Status    string    `gorm:"column:status;type:varchar(16);not null;default:open" json:"status"`                // open(待回答)/answered(已回答)/closed(已关闭)
	ViewCount uint32    `gorm:"column:view_count;not null;default:0" json:"view_count"`                            // 浏览数
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                 // 更新时间
}

// TableName 指定表名
func (GuideQuestion) TableName() string {
	return "guide_question"
}

// GuideAnswer 攻略回答表
type GuideAnswer struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	QuestionID uint64    `gorm:"column:question_id;not null;index:idx_ga_question" json:"question_id"` // 关联问题ID
	UserID     uint64    `gorm:"column:user_id;not null;index:idx_ga_user" json:"user_id"`            // 回答用户ID
	Content    string    `gorm:"column:content;type:text;not null" json:"content"`                    // 回答内容
	LikeCount  uint32    `gorm:"column:like_count;not null;default:0" json:"like_count"`              // 点赞数
	IsAccepted int8      `gorm:"column:is_accepted;not null;default:0" json:"is_accepted"`            // 是否被采纳:0-否 1-是
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                  // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                   // 更新时间
}

// TableName 指定表名
func (GuideAnswer) TableName() string {
	return "guide_answer"
}

// GuideAnswerLike 攻略回答点赞表
type GuideAnswerLike struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AnswerID  uint64    `gorm:"column:answer_id;not null;uniqueIndex:uk_answer_user,priority:1" json:"answer_id"` // 回答ID
	UserID    uint64    `gorm:"column:user_id;not null;uniqueIndex:uk_answer_user,priority:2;index:idx_gal_user" json:"user_id"` // 点赞用户ID
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                              // 点赞时间
}

// TableName 指定表名
func (GuideAnswerLike) TableName() string {
	return "guide_answer_like"
}
