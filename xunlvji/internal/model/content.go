package model

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// Place 地点表
// 注意: location POINT(SRID 4326) 字段不在此结构体中映射,需要单独通过原生 SQL 处理空间数据,
// 仅保留 longitude/latitude 两个 DECIMAL 冗余字段便于读取。
type Place struct {
	ID             uint64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name           string          `gorm:"column:name;type:varchar(128);not null" json:"name"`                                    // 地点名称
	Category       string          `gorm:"column:category;type:varchar(32);not null;index:idx_category;index:idx_city_category,priority:2" json:"category"` // 分类(6大):food/fun/sightseeing/outdoor/shopping/accommodation
	Longitude      decimal.Decimal `gorm:"column:longitude;type:decimal(10,7);not null" json:"longitude"`                         // 经度(冗余,便于读取)
	Latitude       decimal.Decimal `gorm:"column:latitude;type:decimal(10,7);not null" json:"latitude"`                           // 纬度(冗余,便于读取)
	Address        *string         `gorm:"column:address;type:varchar(256)" json:"address,omitempty"`                             // 详细地址
	OpeningHours   *string         `gorm:"column:opening_hours;type:varchar(128)" json:"opening_hours,omitempty"`                 // 营业时间
	Phone          *string         `gorm:"column:phone;type:varchar(20)" json:"phone,omitempty"`                                  // 联系电话
	AvgCost        *decimal.Decimal `gorm:"column:avg_cost;type:decimal(10,2)" json:"avg_cost,omitempty"`                          // 人均消费(元)
	CoverImage     json.RawMessage `gorm:"column:cover_image;type:json" json:"cover_image,omitempty"`                              // 封面图URL数组
	RecommendLevel *string         `gorm:"column:recommend_level;type:varchar(16);index:idx_recommend_level" json:"recommend_level,omitempty"` // 推荐等级:must_go/can_go/avoid
	SuggestDuration *int           `gorm:"column:suggest_duration" json:"suggest_duration,omitempty"`                              // 建议时长(分钟)
	Tags           json.RawMessage `gorm:"column:tags;type:json" json:"tags,omitempty"`                                            // 标签
	City           string          `gorm:"column:city;type:varchar(64);not null;index:idx_city;index:idx_city_category,priority:1" json:"city"` // 所在城市
	Rating         decimal.Decimal `gorm:"column:rating;type:decimal(2,1);not null;default:0.0" json:"rating"`                     // 综合评分(0.0-5.0)
	CheckinCount   int             `gorm:"column:checkin_count;not null;default:0" json:"checkin_count"`                           // 打卡总数
	CommentCount   int             `gorm:"column:comment_count;not null;default:0" json:"comment_count"`                           // 评论总数
	ExtraFields    json.RawMessage `gorm:"column:extra_fields;type:json" json:"extra_fields,omitempty"`                            // 专属模板数据(吃喝:菜系/招牌菜;住宿:房型/设施等)
	CPSLinks       json.RawMessage `gorm:"column:cps_links;type:json" json:"cps_links,omitempty"`                                  // CPS分销链接(ctrip_ticket/meituan_deal/eleme_shop/amap_taxi等)
	MultiLang      json.RawMessage `gorm:"column:multi_lang;type:json" json:"multi_lang,omitempty"`                                // 多语言字段,预留入境游({"en":{"name":"...","desc":"..."},"ja":{...}})
	ForeignerInfo  json.RawMessage `gorm:"column:foreigner_info;type:json" json:"foreigner_info,omitempty"`                        // 外国游客专属信息,预留(payment/visa_required/english_guide/signage_lang等)
	TargetAudience string          `gorm:"column:target_audience;type:varchar(20);not null;default:domestic;index:idx_target_audience" json:"target_audience"` // 目标受众:domestic/inbound/outbound
	CreatedAt      time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                      // 创建时间
	UpdatedAt      time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                      // 更新时间
}

// TableName 指定表名
func (Place) TableName() string {
	return "place"
}

// CheckinRecord 打卡记录表(UGC内容表,打卡和发布共用)
type CheckinRecord struct {
	ID              uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID          uint64           `gorm:"column:user_id;not null;index:idx_user_id;index:idx_user_place,priority:1" json:"user_id"` // 用户ID
	PlaceID         uint64           `gorm:"column:place_id;not null;index:idx_place_id;index:idx_user_place,priority:2" json:"place_id"` // 地点ID
	Images          json.RawMessage  `gorm:"column:images;type:json" json:"images,omitempty"`                                          // 图片列表
	Rating          *decimal.Decimal `gorm:"column:rating;type:decimal(2,1)" json:"rating,omitempty"`                                  // 评分(0.0-5.0)
	Description     *string          `gorm:"column:description;type:text" json:"description,omitempty"`                                // 感受描述
	ActualCost      *decimal.Decimal `gorm:"column:actual_cost;type:decimal(8,2)" json:"actual_cost,omitempty"`                        // 实际消费(元)
	PitfallTips     json.RawMessage  `gorm:"column:pitfall_tips;type:json" json:"pitfall_tips,omitempty"`                               // 避雷建议列表(JSON数组)
	CheckinTime     time.Time        `gorm:"column:checkin_time;not null;default:CURRENT_TIMESTAMP;index:idx_checkin_time" json:"checkin_time"` // 打卡时间
	CheckinType     string           `gorm:"column:checkin_type;type:varchar(16);not null;default:normal" json:"checkin_type"`         // store-到店打卡 normal-普通打卡
	Visibility      string           `gorm:"column:visibility;type:varchar(16);not null;default:public;index:idx_visibility" json:"visibility"` // public-公开(发布到打卡地页面) private-仅自己可见
	Status          string           `gorm:"column:status;type:varchar(16);not null;default:pending;index:idx_status" json:"status"`   // pending-审核中 approved-已通过 rejected-已驳回
	RejectReason    *string          `gorm:"column:reject_reason;type:varchar(256)" json:"reject_reason,omitempty"`                    // 审核驳回原因
	Source          string           `gorm:"column:source;type:varchar(20);not null;default:trip_checkin;index:idx_source" json:"source"` // 内容来源:trip_checkin-行程打卡 manual_publish-主动发布
	ExtraData       json.RawMessage  `gorm:"column:extra_data;type:json" json:"extra_data,omitempty"`                                  // 打卡专属内容(吃喝:菜品评价;住宿:房型体验等),按place.category差异化
	ContentRichness string           `gorm:"column:content_richness;type:varchar(10);not null;default:light;index:idx_content_richness" json:"content_richness"` // 内容丰富度:light/standard/rich,系统自动计算
	GPSVerified     bool             `gorm:"column:gps_verified;not null;default:0" json:"gps_verified"`                               // GPS验证标记:1-100米范围内打卡 0-未验证
	CreatedAt       time.Time        `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                        // 创建时间
	UpdatedAt       time.Time        `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                        // 更新时间
}

// TableName 指定表名
func (CheckinRecord) TableName() string {
	return "checkin_record"
}

// Route 路线表
type Route struct {
	ID             uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title          string           `gorm:"column:title;type:varchar(128);not null" json:"title"`                                    // 路线标题
	CoverImage     json.RawMessage  `gorm:"column:cover_image;type:json" json:"cover_image,omitempty"`                               // 封面图URL数组
	TotalDuration  int              `gorm:"column:total_duration;not null" json:"total_duration"`                                    // 总时长(分钟)
	TotalBudget    *decimal.Decimal `gorm:"column:total_budget;type:decimal(10,2)" json:"total_budget,omitempty"`                    // 总预算(元)
	SuitableFor    *string          `gorm:"column:suitable_for;type:varchar(64)" json:"suitable_for,omitempty"`                      // 适合人群
	Category       string           `gorm:"column:category;type:varchar(32);not null;index:idx_category" json:"category"`            // 分类:half_day/one_day/two_day/three_day_plus
	Tags           json.RawMessage  `gorm:"column:tags;type:json" json:"tags,omitempty"`                                             // 标签
	CreatorID      uint64           `gorm:"column:creator_id;not null;index:idx_creator_id" json:"creator_id"`                        // 创建者ID
	Source         string           `gorm:"column:source;type:varchar(16);not null;default:user;index:idx_source" json:"source"`     // user/AI/official
	UsageCount     int              `gorm:"column:usage_count;not null;default:0;index:idx_usage_count" json:"usage_count"`          // 使用人数
	FavoriteCount  int              `gorm:"column:favorite_count;not null;default:0" json:"favorite_count"`                           // 收藏数
	CommentCount   int              `gorm:"column:comment_count;not null;default:0" json:"comment_count"`                             // 评论总数
	Rating         decimal.Decimal  `gorm:"column:rating;type:decimal(2,1);not null;default:0.0;index:idx_rating" json:"rating"`     // 综合评分(0.0-5.0)
	Status         int8             `gorm:"column:status;not null;default:1;index:idx_status" json:"status"`                          // 0-下架 1-正常 2-草稿
	CreatedAt      time.Time        `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"`                  // 创建时间
	UpdatedAt      time.Time        `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                       // 更新时间
}

// TableName 指定表名
func (Route) TableName() string {
	return "route"
}

// RoutePoint 路线点位表
type RoutePoint struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RouteID      uint64           `gorm:"column:route_id;not null;index:idx_route_id;index:idx_route_sort,priority:1" json:"route_id"` // 路线ID
	SortOrder    int              `gorm:"column:sort_order;not null;default:1;index:idx_route_sort,priority:2" json:"sort_order"`     // 排序序号(从1开始)
	PointTime    *string          `gorm:"column:point_time;type:time" json:"point_time,omitempty"`                                    // 预计到达时间(HH:MM:SS)
	PlaceID      *uint64          `gorm:"column:place_id" json:"place_id,omitempty"`                                                   // 地点ID(可选)
	Name         string           `gorm:"column:name;type:varchar(128);not null" json:"name"`                                          // 点位名称
	StayDuration *int             `gorm:"column:stay_duration" json:"stay_duration,omitempty"`                                         // 停留时长(分钟)
	Cost         *decimal.Decimal `gorm:"column:cost;type:decimal(10,2)" json:"cost,omitempty"`                                        // 预计消费(元)
	Category     *string          `gorm:"column:category;type:varchar(32)" json:"category,omitempty"`                                  // 分类
	Transport    *string          `gorm:"column:transport;type:varchar(32)" json:"transport,omitempty"`                                // 交通方式
	CreatedAt    time.Time        `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                          // 创建时间
	UpdatedAt    time.Time        `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                          // 更新时间
}

// TableName 指定表名
func (RoutePoint) TableName() string {
	return "route_point"
}

// RouteRelated 路线关联表
type RouteRelated struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RouteID        uint64    `gorm:"column:route_id;not null;uniqueIndex:uk_route_related,priority:1" json:"route_id"`                    // 路线ID
	RelatedRouteID uint64    `gorm:"column:related_route_id;not null;uniqueIndex:uk_route_related,priority:2;index:idx_related_route_id" json:"related_route_id"` // 关联路线ID
	SortOrder      int       `gorm:"column:sort_order;not null;default:1" json:"sort_order"`                                              // 排序序号
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                                   // 创建时间
}

// TableName 指定表名
func (RouteRelated) TableName() string {
	return "route_related"
}

// RouteTransport 行程交通表
type RouteTransport struct {
	ID            uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TripID        uint64           `gorm:"column:trip_id;not null;index:idx_trip_id;index:idx_trip_sort,priority:1" json:"trip_id"` // 行程ID
	FromPlaceID   uint64           `gorm:"column:from_place_id;not null;index:idx_from_place_id" json:"from_place_id"`              // 起点地点ID
	ToPlaceID     uint64           `gorm:"column:to_place_id;not null;index:idx_to_place_id" json:"to_place_id"`                    // 终点地点ID
	TransportType string           `gorm:"column:transport_type;type:varchar(20);not null" json:"transport_type"`                   // 交通方式:walk/bus/taxi/drive
	Distance      *int             `gorm:"column:distance" json:"distance,omitempty"`                                               // 距离(米)
	Duration      *int             `gorm:"column:duration" json:"duration,omitempty"`                                               // 预计时长(秒)
	Cost          *decimal.Decimal `gorm:"column:cost;type:decimal(8,2)" json:"cost,omitempty"`                                     // 预计费用(元)
	CPSLink       *string          `gorm:"column:cps_link;type:varchar(500)" json:"cps_link,omitempty"`                             // 打车CPS链接(高德)
	SortOrder     int              `gorm:"column:sort_order;not null;default:1;index:idx_trip_sort,priority:2" json:"sort_order"`   // 排序序号(从1开始)
	CreatedAt     time.Time        `gorm:"column:created_at;autoCreateTime" json:"created_at"`                                      // 创建时间
}

// TableName 指定表名
func (RouteTransport) TableName() string {
	return "route_transport"
}

// Guide 攻略表
type Guide struct {
	ID               uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title            string           `gorm:"column:title;type:varchar(128);not null" json:"title"`                                     // 攻略标题
	CoverImage       json.RawMessage  `gorm:"column:cover_image;type:json" json:"cover_image,omitempty"`                                // 封面图URL数组
	Category         string           `gorm:"column:category;type:varchar(32);not null;index:idx_category;index:idx_city_category,priority:2" json:"category"` // 分类(7大):nature/history/entertainment/urban/transport/red_tourism/religion
	City             string           `gorm:"column:city;type:varchar(64);not null;index:idx_city;index:idx_city_category,priority:1" json:"city"` // 所在城市
	Duration         *int             `gorm:"column:duration" json:"duration,omitempty"`                                                // 游玩时长(分钟)
	Price            *decimal.Decimal `gorm:"column:price;type:decimal(10,2)" json:"price,omitempty"`                                   // 价格(元)
	ReadDuration     *int             `gorm:"column:read_duration" json:"read_duration,omitempty"`                                      // 阅读时长(分钟)
	Address          *string          `gorm:"column:address;type:varchar(256)" json:"address,omitempty"`                                // 地址
	OpeningHours     *string          `gorm:"column:opening_hours;type:varchar(128)" json:"opening_hours,omitempty"`                    // 开放时间
	Transport        *string          `gorm:"column:transport;type:varchar(256)" json:"transport,omitempty"`                            // 交通方式
	Content          string           `gorm:"column:content;type:mediumtext;not null" json:"content"`                                   // 正文内容(Markdown)
	Highlights       *string          `gorm:"column:highlights;type:text" json:"highlights,omitempty"`                                  // 必看亮点
	PitfallReminders *string          `gorm:"column:pitfall_reminders;type:text" json:"pitfall_reminders,omitempty"`                    // 避雷提醒
	VerifiedCount    int              `gorm:"column:verified_count;not null;default:0;index:idx_verified_count" json:"verified_count"` // 验证人数
	CommentCount     int              `gorm:"column:comment_count;not null;default:0" json:"comment_count"`                             // 评论总数
	Source           string           `gorm:"column:source;type:varchar(16);not null;default:AI" json:"source"`                         // 来源标签(固定AI)
	Status           int8             `gorm:"column:status;not null;default:1;index:idx_status" json:"status"`                           // 0-下架 1-正常 2-草稿
	CreatedAt        time.Time        `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"`                   // 创建时间
	UpdatedAt        time.Time        `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                                        // 更新时间
}

// TableName 指定表名
func (Guide) TableName() string {
	return "guide"
}
