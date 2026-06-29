package model

import "time"

// Comment 评论表
type Comment struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TargetType string    `gorm:"column:target_type;type:varchar(16);not null;index:idx_target,priority:1" json:"target_type"` // 目标类型:route/place/guide
	TargetID   uint64    `gorm:"column:target_id;not null;index:idx_target,priority:2" json:"target_id"`                        // 目标ID
	UserID     uint64    `gorm:"column:user_id;not null;index:idx_user_id" json:"user_id"`                                     // 用户ID
	ParentID   *uint64   `gorm:"column:parent_id;index:idx_parent_id" json:"parent_id,omitempty"`                              // 父评论ID(二级回复)
	Content    string    `gorm:"column:content;type:text;not null" json:"content"`                                            // 评论内容
	LikeCount  int       `gorm:"column:like_count;not null;default:0" json:"like_count"`                                       // 点赞数
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"`                      // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                           // 更新时间
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comment"
}

// CommentLike 评论点赞表
type CommentLike struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CommentID uint64    `gorm:"column:comment_id;not null;uniqueIndex:uk_comment_user,priority:1" json:"comment_id"` // 评论ID
	UserID    uint64    `gorm:"column:user_id;not null;uniqueIndex:uk_comment_user,priority:2;index:idx_user_id" json:"user_id"` // 点赞用户ID
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                  // 点赞时间
}

// TableName 指定表名
func (CommentLike) TableName() string {
	return "comment_like"
}

// Follow 关注关系表
type Follow struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	FollowerID  uint64    `gorm:"column:follower_id;not null;uniqueIndex:uk_follower_following,priority:1;index:idx_follower_id" json:"follower_id"` // 关注者ID
	FollowingID uint64    `gorm:"column:following_id;not null;uniqueIndex:uk_follower_following,priority:2;index:idx_following_id" json:"following_id"` // 被关注者ID
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                                                 // 关注时间
}

// TableName 指定表名
func (Follow) TableName() string {
	return "follow"
}

// Favorite 收藏表
type Favorite struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     uint64    `gorm:"column:user_id;not null;uniqueIndex:uk_user_target,priority:1" json:"user_id"`                              // 用户ID
	TargetType string    `gorm:"column:target_type;type:varchar(16);not null;uniqueIndex:uk_user_target,priority:2;index:idx_target,priority:1" json:"target_type"` // 目标类型:route/place/guide/trip
	TargetID   uint64    `gorm:"column:target_id;not null;uniqueIndex:uk_user_target,priority:3;index:idx_target,priority:2" json:"target_id"` // 目标ID
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                                         // 收藏时间
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorite"
}

// Message 消息表
type Message struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Type       string    `gorm:"column:type;type:varchar(16);not null;index:idx_type" json:"type"`                              // 消息类型:system/private/group
	SenderID   *uint64   `gorm:"column:sender_id;index:idx_sender_receiver,priority:1" json:"sender_id,omitempty"`              // 发送者ID
	ReceiverID uint64    `gorm:"column:receiver_id;not null;index:idx_receiver_read,priority:1;index:idx_sender_receiver,priority:2" json:"receiver_id"` // 接收者ID/群ID(应用层校验)
	Content    string    `gorm:"column:content;type:text;not null" json:"content"`                                              // 消息内容
	IsRead     int8      `gorm:"column:is_read;not null;default:0;index:idx_receiver_read,priority:2" json:"is_read"`            // 0-未读 1-已读
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"`                       // 发送时间
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                            // 更新时间
}

// TableName 指定表名
func (Message) TableName() string {
	return "message"
}

// TipVote 避雷点赞表
type TipVote struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CheckinID uint64    `gorm:"column:checkin_id;not null;uniqueIndex:uk_checkin_tip_user,priority:1;index:idx_checkin_tip,priority:1" json:"checkin_id"` // 打卡记录ID
	TipIndex  int       `gorm:"column:tip_index;not null;uniqueIndex:uk_checkin_tip_user,priority:2;index:idx_checkin_tip,priority:2" json:"tip_index"` // 避雷建议索引(0-based)
	UserID    uint64    `gorm:"column:user_id;not null;uniqueIndex:uk_checkin_tip_user,priority:3" json:"user_id"`             // 点赞用户ID
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                            // 点赞时间
}

// TableName 指定表名
func (TipVote) TableName() string {
	return "tip_vote"
}
