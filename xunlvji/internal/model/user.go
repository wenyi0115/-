package model

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// User 用户表
type User struct {
	ID               uint64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nickname         string          `gorm:"column:nickname;type:varchar(64);not null" json:"nickname"`                              // 用户昵称
	Avatar           *string         `gorm:"column:avatar;type:varchar(512)" json:"avatar,omitempty"`                                // 头像URL
	Phone            *string         `gorm:"column:phone;type:varchar(20);uniqueIndex:uk_phone" json:"phone,omitempty"`              // 手机号
	Email            *string         `gorm:"column:email;type:varchar(128);uniqueIndex:uk_email" json:"email,omitempty"`             // 邮箱
	PasswordHash     *string         `gorm:"column:password_hash;type:varchar(256)" json:"-"`                                         // bcrypt密码哈希(第三方登录可为空)
	Role             string          `gorm:"column:role;type:varchar(20);not null;default:user;index:idx_role" json:"role"`          // Deprecated: 【已废弃·仅向后兼容】角色:user/creator/reviewer/operator/admin;新数据统一使用 user_role 关联表查询用户角色,新代码不应再读写此字段,保留仅为兼容历史数据与灰度迁移
	CreditScore      int             `gorm:"column:credit_score;not null;default:100" json:"credit_score"`                            // 信用分(0-100)
	PreferTags       json.RawMessage `gorm:"column:prefer_tags;type:json" json:"prefer_tags,omitempty"`                               // 偏好标签
	VisitedCityCount int             `gorm:"column:visited_city_count;not null;default:0" json:"visited_city_count"`                  // 去过城市数
	FollowerCount    int             `gorm:"column:follower_count;not null;default:0" json:"follower_count"`                          // 粉丝数
	FollowingCount   int             `gorm:"column:following_count;not null;default:0" json:"following_count"`                        // 关注数
	LikeCount        int             `gorm:"column:like_count;not null;default:0" json:"like_count"`                                  // 获赞数
	Bio              *string         `gorm:"column:bio;type:varchar(200)" json:"bio,omitempty"`                                       // 个人简介
	Birthday         *time.Time      `gorm:"column:birthday;type:date" json:"birthday,omitempty"`                                     // 生日
	CreatedAt        time.Time       `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"`                  // 注册时间
	UpdatedAt        time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                      // 更新时间
	DeleteRequestedAt *time.Time     `gorm:"column:delete_requested_at" json:"delete_requested_at,omitempty"`                         // 注销申请时间(30天冷静期)
	Language         string          `gorm:"column:language;type:varchar(16);not null;default:zh-CN" json:"language"`                  // 语言偏好
	Status           int8            `gorm:"column:status;not null;default:1;index:idx_status" json:"status"`                          // 0-禁用 1-正常 2-注销
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// UserOAuth 第三方账号关联表
type UserOAuth struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint64     `gorm:"column:user_id;not null;index:idx_user_id" json:"user_id"`                       // 用户ID
	Provider  string     `gorm:"column:provider;type:varchar(32);not null;uniqueIndex:uk_provider_oauth,priority:1" json:"provider"`  // 提供方:apple/google/wechat
	OauthID   string     `gorm:"column:oauth_id;type:varchar(128);not null;uniqueIndex:uk_provider_oauth,priority:2" json:"oauth_id"` // 第三方用户ID
	UnionID   *string    `gorm:"column:union_id;type:varchar(128);index:idx_union_id" json:"union_id,omitempty"` // 联合ID(微信unionId)
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`                              // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                              // 更新时间
}

// TableName 指定表名
func (UserOAuth) TableName() string {
	return "user_oauth"
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"column:user_id;not null;uniqueIndex:uk_user_role,priority:1" json:"user_id"` // 用户ID
	RoleID    uint64    `gorm:"column:role_id;not null;uniqueIndex:uk_user_role,priority:2;index:idx_role_id" json:"role_id"` // 角色ID
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                         // 创建时间
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_role"
}

// UserBadge 用户勋章表
type UserBadge struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     uint64    `gorm:"column:user_id;not null;uniqueIndex:uk_user_badge,priority:1" json:"user_id"` // 用户ID
	BadgeID    uint64    `gorm:"column:badge_id;not null;uniqueIndex:uk_user_badge,priority:2;index:idx_badge_id" json:"badge_id"` // 勋章ID
	UnlockedAt time.Time `gorm:"column:unlocked_at;not null;default:CURRENT_TIMESTAMP" json:"unlocked_at"`     // 解锁时间
}

// TableName 指定表名
func (UserBadge) TableName() string {
	return "user_badge"
}

// UserMedalSummary 勋章汇总表
type UserMedalSummary struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      uint64    `gorm:"column:user_id;not null;uniqueIndex:uk_user_id" json:"user_id"` // 用户ID
	TotalCount  int       `gorm:"column:total_count;not null;default:0" json:"total_count"`      // 勋章总数
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`            // 更新时间
}

// TableName 指定表名
func (UserMedalSummary) TableName() string {
	return "user_medal_summary"
}

// UserSetting 用户设置表
type UserSetting struct {
	ID                      uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID                  uint64    `gorm:"column:user_id;not null;uniqueIndex:uk_user_id" json:"user_id"`                                                        // 用户ID
	FootprintPublic         string    `gorm:"column:footprint_public;type:varchar(16);not null;default:public" json:"footprint_public"`                            // 足迹可见性:public(公开)/friends(仅好友)/private(仅自己)
	ContentPublic           int8      `gorm:"column:content_public;not null;default:1" json:"content_public"`                                                       // 内容是否公开:0-否 1-是
	AllowDMScope            string    `gorm:"column:allow_dm_scope;type:varchar(16);not null;default:all" json:"allow_dm_scope"`                                   // 私信权限:all(所有人)/followers(仅关注者)/nobody(不允许)
	FollowScope             string    `gorm:"column:follow_scope;type:varchar(16);not null;default:all" json:"follow_scope"`                                       // 关注权限:all(所有人)/verify(需验证)/nobody(不允许)
	LocationPublic          int8      `gorm:"column:location_public;not null;default:0" json:"location_public"`                                                     // 位置是否公开:0-否 1-是(默认不公开,用户主动选择公开)
	HistoryPublic           int8      `gorm:"column:history_public;not null;default:0" json:"history_public"`                                                       // 历史行程是否公开:0-否 1-是
	ShowFavorites           int8      `gorm:"column:show_favorites;not null;default:0" json:"show_favorites"`                                                       // 收藏列表是否对外可见:0-不公开 1-公开
	ShowFollowing           int8      `gorm:"column:show_following;not null;default:0" json:"show_following"`                                                       // 关注/粉丝列表是否对外可见:0-不公开 1-公开
	ProfileVisibility       string    `gorm:"column:profile_visibility;type:varchar(16);not null;default:public" json:"profile_visibility"`                         // 个人主页可见性:public(公开)/friends(仅好友)/private(仅自己)
	PersonalizedRecomm      int8      `gorm:"column:personalized_recomm;not null;default:1" json:"personalized_recomm"`                                             // 个性化推荐开关:0-关闭 1-开启
	NotificationPush        int8      `gorm:"column:notification_push;not null;default:1" json:"notification_push"`                                                  // 推送开关:0-关闭 1-开启(总开关)
	NotificationInteractive int8      `gorm:"column:notification_interactive;not null;default:1" json:"notification_interactive"`                                    // 互动通知:0-关闭 1-开启(点赞/评论/关注/私信提醒)
	NotificationTrip        int8      `gorm:"column:notification_trip;not null;default:1" json:"notification_trip"`                                                  // 行程提醒:0-关闭 1-开启(行程开始/打卡提醒等)
	Language                string    `gorm:"column:language;type:varchar(10);not null;default:zh-CN" json:"language"`                                              // 语言:zh-CN/en-US等
	FontSize                int8      `gorm:"column:font_size;not null;default:0" json:"font_size"`                                                                  // 字体大小:0-小 1-中 2-大
	DarkMode                int8      `gorm:"column:dark_mode;not null;default:0" json:"dark_mode"`                                                                  // 深色模式:0-跟随系统 1-关闭 2-开启
	DailyPushLimit          int       `gorm:"column:daily_push_limit;not null;default:10" json:"daily_push_limit"`                                                   // 每日推送上限:单用户每日最多接收推送条数
	UpdatedAt               time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                                                    // 更新时间
}

// TableName 指定表名
func (UserSetting) TableName() string {
	return "user_setting"
}

// UserCityProgress 用户城市进度表
type UserCityProgress struct {
	ID             uint64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID         uint64          `gorm:"column:user_id;not null;uniqueIndex:uk_user_city,priority:1" json:"user_id"`                    // 用户ID
	CityID         uint64          `gorm:"column:city_id;not null;uniqueIndex:uk_user_city,priority:2;index:idx_city_id" json:"city_id"` // 城市ID
	TotalPlaces    int             `gorm:"column:total_places;not null;default:0" json:"total_places"`                                    // 总地点数
	CheckedPlaces  int             `gorm:"column:checked_places;not null;default:0" json:"checked_places"`                                // 已打卡地点数
	Progress       decimal.Decimal `gorm:"column:progress;type:decimal(5,2);not null;default:0.00" json:"progress"`                       // 进度百分比(0-100)
	UnlockedAt     *time.Time      `gorm:"column:unlocked_at" json:"unlocked_at,omitempty"`                                               // 解锁时间
	UpdatedAt      time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                            // 更新时间
}

// TableName 指定表名
func (UserCityProgress) TableName() string {
	return "user_city_progress"
}

// UserProvinceProgress 用户省份进度表
type UserProvinceProgress struct {
	ID         uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     uint64     `gorm:"column:user_id;not null;uniqueIndex:uk_user_province,priority:1" json:"user_id"`                       // 用户ID
	ProvinceID uint64     `gorm:"column:province_id;not null;uniqueIndex:uk_user_province,priority:2;index:idx_province_id" json:"province_id"` // 省份ID
	Unlocked   int8       `gorm:"column:unlocked;not null;default:0" json:"unlocked"`                                                   // 是否解锁:0-未解锁 1-已解锁
	UnlockedAt *time.Time `gorm:"column:unlocked_at" json:"unlocked_at,omitempty"`                                                       // 解锁时间
}

// TableName 指定表名
func (UserProvinceProgress) TableName() string {
	return "user_province_progress"
}

// UserNationalRanking 全国排名表
type UserNationalRanking struct {
	ID            uint64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        uint64          `gorm:"column:user_id;not null;uniqueIndex:uk_user_id" json:"user_id"`            // 用户ID
	CityCount     int             `gorm:"column:city_count;not null;default:0" json:"city_count"`                   // 打卡城市数
	CheckinCount  int             `gorm:"column:checkin_count;not null;default:0" json:"checkin_count"`             // 打卡总数
	TotalScore    decimal.Decimal `gorm:"column:total_score;type:decimal(10,2);not null;default:0.00" json:"total_score"` // 总积分
	Rank          int             `gorm:"column:rank;not null;default:0;index:idx_rank" json:"rank"`                // 排名
	UpdatedAt     time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                       // 更新时间
}

// TableName 指定表名
func (UserNationalRanking) TableName() string {
	return "user_national_ranking"
}

// CreatorApplication 创作者认证申请表
type CreatorApplication struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      uint64     `gorm:"column:user_id;not null;index:idx_user_id" json:"user_id"`                // 用户ID
	RealName    string     `gorm:"column:real_name;type:varchar(64);not null" json:"real_name"`             // 真实姓名
	IDCard      string     `gorm:"column:id_card;type:varchar(18);not null" json:"id_card"`                 // 身份证号
	Portfolio   *string    `gorm:"column:portfolio;type:varchar(512)" json:"portfolio,omitempty"`           // 作品集URL
	Status      string     `gorm:"column:status;type:varchar(16);not null;default:pending;index:idx_status" json:"status"` // 状态:pending/approved/rejected
	ReviewerID  *uint64    `gorm:"column:reviewer_id" json:"reviewer_id,omitempty"`                          // 审核人ID
	ReviewNote  *string    `gorm:"column:review_note;type:varchar(256)" json:"review_note,omitempty"`       // 审核备注
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`                       // 创建时间
	ReviewedAt  *time.Time `gorm:"column:reviewed_at" json:"reviewed_at,omitempty"`                          // 审核时间
}

// TableName 指定表名
func (CreatorApplication) TableName() string {
	return "creator_application"
}
