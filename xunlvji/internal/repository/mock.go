// Package repository 数据访问层 - Mock 实现（开发预览用）
package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
	"github.com/xunlvji/xunlvji/internal/model"
)

// === Mock 数据（荆州古城主题）===

// mockUser 测试用户
var mockUser = &model.User{
	ID:         10001,
	Nickname:   "小明爱旅行",
	Avatar:     strPtr("https://cdn.example.com/avatars/10001.jpg"),
	Role:       "creator",
	CreditScore: 95,
	Language:   "zh-CN",
	Status:     1,
	CreatedAt:  time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
}

// mockPlace1 荆州古城墙
var mockPlace1 = model.Place{
	ID:             3001,
	Name:           "荆州古城墙",
	CoverImage:     jsonRaw(`["https://cdn.example.com/places/3001_1.jpg","https://cdn.example.com/places/3001_2.jpg","https://cdn.example.com/places/3001_3.jpg"]`),
	Category:       "sightseeing",
	RecommendLevel: strPtr("must_go"),
	AvgCost:        decimalPtr("0"),
	Rating:         decimalPtr2("4.7"),
	CheckinCount:   520,
	City:           "荆州",
	Address:        strPtr("湖北省荆州市荆州区荆南路"),
	Longitude:      decimalVal("112.2400"),
	Latitude:       decimalVal("30.3300"),
	OpeningHours:   strPtr("08:00-18:00"),
	Phone:          strPtr("0716-1234567"),
	SuggestDuration: intPtr(90),
	Tags:           jsonRaw(`["历史","必去","拍照"]`),
	CreatedAt:      time.Date(2026, 3, 1, 8, 0, 0, 0, time.UTC),
}

// mockPlace2 荆州博物馆
var mockPlace2 = model.Place{
	ID:             3002,
	Name:           "荆州博物馆",
	CoverImage:     jsonRaw(`["https://cdn.example.com/places/3002_1.jpg","https://cdn.example.com/places/3002_2.jpg"]`),
	Category:       "sightseeing",
	RecommendLevel: strPtr("must_go"),
	AvgCost:        decimalPtr("0"),
	Rating:         decimalPtr2("4.5"),
	CheckinCount:   380,
	City:           "荆州",
	Address:        strPtr("湖北省荆州市荆州区荆中路"),
	Longitude:      decimalVal("112.2450"),
	Latitude:       decimalVal("30.3350"),
	OpeningHours:   strPtr("09:00-17:00"),
	Phone:          strPtr("0716-7654321"),
	SuggestDuration: intPtr(120),
	Tags:           jsonRaw(`["博物馆","历史","免费"]`),
	CreatedAt:      time.Date(2026, 3, 5, 9, 0, 0, 0, time.UTC),
}

// mockPlace3 老街美食
var mockPlace3 = model.Place{
	ID:             3003,
	Name:           "老街美食",
	CoverImage:     jsonRaw(`["https://cdn.example.com/places/3003_1.jpg"]`),
	Category:       "food",
	RecommendLevel: strPtr("can_go"),
	AvgCost:        decimalPtr("50"),
	Rating:         decimalPtr2("4.3"),
	CheckinCount:   210,
	City:           "荆州",
	Address:        strPtr("湖北省荆州市沙市区中山路"),
	Longitude:      decimalVal("112.2500"),
	Latitude:       decimalVal("30.3200"),
	OpeningHours:   strPtr("10:00-22:00"),
	Phone:          strPtr("0716-8888888"),
	SuggestDuration: intPtr(60),
	Tags:           jsonRaw(`["美食","小吃","本地"]`),
	CreatedAt:      time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
}

// mockRoute1 荆州古城一日游
var mockRoute1 = model.Route{
	ID:            2001,
	Title:         "荆州古城一日游",
	CoverImage:    jsonRaw(`["https://cdn.example.com/routes/2001_1.jpg","https://cdn.example.com/routes/2001_2.jpg","https://cdn.example.com/routes/2001_3.jpg"]`),
	CreatorID:     10001,
	Category:      "one_day",
	TotalDuration: 480,
	TotalBudget:   decimalPtr("350"),
	SuitableFor:   strPtr("情侣/朋友"),
	Rating:        decimalPtr2("4.5"),
	UsageCount:    320,
	FavoriteCount: 150,
	Source:        "user",
	Tags:          jsonRaw(`["文艺","美食","历史文化"]`),
	Status:        1,
	CommentCount:  12,
	CreatedAt:     time.Date(2026, 5, 20, 9, 0, 0, 0, time.UTC),
	UpdatedAt:     time.Date(2026, 6, 18, 15, 0, 0, 0, time.UTC),
}

// mockRoute2 荆州两日深度游
var mockRoute2 = model.Route{
	ID:            2002,
	Title:         "荆州两日深度游",
	CoverImage:    jsonRaw(`["https://cdn.example.com/routes/2002_1.jpg","https://cdn.example.com/routes/2002_2.jpg"]`),
	CreatorID:     10001,
	Category:      "two_day",
	TotalDuration: 960,
	TotalBudget:   decimalPtr("680"),
	SuitableFor:   strPtr("家庭/亲子"),
	Rating:        decimalPtr2("4.8"),
	UsageCount:    180,
	FavoriteCount: 95,
	Source:        "user",
	Tags:          jsonRaw(`["深度游","亲子","文化"]`),
	Status:        1,
	CommentCount:  8,
	CreatedAt:     time.Date(2026, 5, 25, 10, 0, 0, 0, time.UTC),
	UpdatedAt:     time.Date(2026, 6, 20, 11, 0, 0, 0, time.UTC),
}

// mockRoutePoints 路线点位
var mockRoutePoints = []model.RoutePoint{
	{ID: 1, RouteID: 2001, SortOrder: 1, PointTime: strPtr("09:00"), PlaceID: uint64Ptr(3001), Name: "荆州古城墙", StayDuration: intPtr(90), Cost: decimalPtr("0"), Category: strPtr("sightseeing"), Transport: strPtr("步行")},
	{ID: 2, RouteID: 2001, SortOrder: 2, PointTime: strPtr("10:30"), PlaceID: uint64Ptr(3002), Name: "荆州博物馆", StayDuration: intPtr(120), Cost: decimalPtr("0"), Category: strPtr("sightseeing"), Transport: strPtr("步行")},
	{ID: 3, RouteID: 2001, SortOrder: 3, PointTime: strPtr("12:30"), PlaceID: uint64Ptr(3003), Name: "老街美食", StayDuration: intPtr(60), Cost: decimalPtr("50"), Category: strPtr("food"), Transport: strPtr("步行")},
}

// mockGuide1 荆州古城深度游攻略
var mockGuide1 = model.Guide{
	ID:               4001,
	Title:            "荆州古城深度游攻略",
	CoverImage:       jsonRaw(`["https://cdn.example.com/guides/4001_1.jpg","https://cdn.example.com/guides/4001_2.jpg"]`),
	Category:         "history",
	City:             "荆州",
	Duration:         intPtr(360),
	Price:            decimalPtr("0"),
	ReadDuration:     intPtr(8),
	Address:          strPtr("湖北省荆州市荆州区"),
	OpeningHours:     strPtr("08:00-18:00"),
	Transport:        strPtr("公交：乘坐1路、12路到古城站下车；自驾：导航至荆州古城停车场"),
	Content:          "## 荆州古城简介\n\n荆州古城，又名江陵城，是中国南方保存最完好的古城墙之一。作为三国文化的重要发源地，荆州见证了无数历史风云。\n\n## 必看景点\n\n1. **古城墙**：全长约11公里，现存城墙为明清时期重建\n2. **宾阳楼**：古城东门城楼，可俯瞰全城风貌\n3. **张居正故居**：明代首辅张居正的故居\n\n## 游览建议\n\n建议安排一天时间，从宾阳楼出发，沿城墙步行至南门，参观张居正故居后返回。夏季注意防晒，城墙上无遮阳。",
	Highlights:       strPtr("必看：古城墙、宾阳楼、张居正故居"),
	PitfallReminders: strPtr("注意：周末人流量大，建议工作日前往"),
	VerifiedCount:    128,
	Source:           "AI",
	Status:           1,
	CommentCount:     15,
	CreatedAt:        time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC),
}

// mockGuide2 自然生态攻略
var mockGuide2 = model.Guide{
	ID:            4002,
	Title:         "荆州周边自然生态游",
	CoverImage:    jsonRaw(`["https://cdn.example.com/guides/4002_1.jpg"]`),
	Category:      "nature",
	City:          "荆州",
	Duration:      intPtr(240),
	Price:         decimalPtr("0"),
	ReadDuration:  intPtr(6),
	Address:       strPtr("湖北省荆州市周边"),
	OpeningHours:  strPtr("全天"),
	Transport:     strPtr("自驾推荐"),
	Content:       "## 荆州周边自然生态\n\n荆州周边有多处自然景观，适合周末短途游。\n\n## 推荐路线\n\n1. **洪湖湿地**：观鸟胜地\n2. **石首麋鹿保护区**：国宝麋鹿\n3. **松滋洈水**：山水风光",
	Highlights:    strPtr("必看：洪湖湿地、石首麋鹿"),
	VerifiedCount: 56,
	Source:         "AI",
	Status:         1,
	CommentCount:   6,
	CreatedAt:      time.Date(2026, 6, 10, 14, 0, 0, 0, time.UTC),
}

// === Mock 仓库实现 ===

// mockPlaceRepo Mock 地点仓库
type mockPlaceRepo struct{}

// NewMockPlaceRepository 创建 Mock 地点仓库
func NewMockPlaceRepository() PlaceRepository {
	return &mockPlaceRepo{}
}

func (r *mockPlaceRepo) List(ctx context.Context, q PlaceQuery, offset, limit int) ([]model.Place, int64, error) {
	all := []model.Place{mockPlace1, mockPlace2, mockPlace3}
	var filtered []model.Place
	for _, p := range all {
		if q.Category != "" && p.Category != q.Category {
			continue
		}
		if q.City != "" && p.City != q.City {
			continue
		}
		filtered = append(filtered, p)
	}
	total := int64(len(filtered))
	end := offset + limit
	if offset > len(filtered) {
		offset = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func (r *mockPlaceRepo) GetByID(ctx context.Context, id uint64) (*model.Place, error) {
	for _, p := range []model.Place{mockPlace1, mockPlace2, mockPlace3} {
		if p.ID == id {
			cp := p
			return &cp, nil
		}
	}
	return nil, ErrPlaceNotFound
}

// mockRouteRepo Mock 路线仓库
type mockRouteRepo struct{}

// NewMockRouteRepository 创建 Mock 路线仓库
func NewMockRouteRepository() RouteRepository {
	return &mockRouteRepo{}
}

func (r *mockRouteRepo) List(ctx context.Context, q RouteQuery, offset, limit int) ([]model.Route, int64, error) {
	all := []model.Route{mockRoute1, mockRoute2}
	var filtered []model.Route
	for _, rt := range all {
		if q.Category != "" && rt.Category != q.Category {
			continue
		}
		filtered = append(filtered, rt)
	}
	total := int64(len(filtered))
	end := offset + limit
	if offset > len(filtered) {
		offset = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func (r *mockRouteRepo) GetByID(ctx context.Context, id uint64) (*model.Route, error) {
	for _, rt := range []model.Route{mockRoute1, mockRoute2} {
		if rt.ID == id {
			cp := rt
			return &cp, nil
		}
	}
	return nil, ErrRouteNotFound
}

func (r *mockRouteRepo) GetPoints(ctx context.Context, routeID uint64) ([]model.RoutePoint, error) {
	var result []model.RoutePoint
	for _, p := range mockRoutePoints {
		if p.RouteID == routeID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *mockRouteRepo) GetCreator(ctx context.Context, userID uint64) (*model.User, error) {
	if mockUser.ID == userID {
		cp := *mockUser
		return &cp, nil
	}
	return nil, ErrUserNotFound
}

// mockGuideRepo Mock 攻略仓库
type mockGuideRepo struct{}

// NewMockGuideRepository 创建 Mock 攻略仓库
func NewMockGuideRepository() GuideRepository {
	return &mockGuideRepo{}
}

func (r *mockGuideRepo) List(ctx context.Context, q GuideQuery, offset, limit int) ([]model.Guide, int64, error) {
	all := []model.Guide{mockGuide1, mockGuide2}
	var filtered []model.Guide
	for _, g := range all {
		if q.Category != "" && g.Category != q.Category {
			continue
		}
		if q.City != "" && g.City != q.City {
			continue
		}
		filtered = append(filtered, g)
	}
	total := int64(len(filtered))
	end := offset + limit
	if offset > len(filtered) {
		offset = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func (r *mockGuideRepo) GetByID(ctx context.Context, id uint64) (*model.Guide, error) {
	for _, g := range []model.Guide{mockGuide1, mockGuide2} {
		if g.ID == id {
			cp := g
			return &cp, nil
		}
	}
	return nil, ErrGuideNotFound
}

// mockCheckinRepo Mock 打卡记录仓库
type mockCheckinRepo struct{}

// NewMockCheckinRepository 创建 Mock 打卡记录仓库
func NewMockCheckinRepository() CheckinRepository {
	return &mockCheckinRepo{}
}

func (r *mockCheckinRepo) ListByPlace(ctx context.Context, placeID uint64, status string, offset, limit int) ([]model.CheckinRecord, int64, error) {
	return []model.CheckinRecord{}, 0, nil
}

func (r *mockCheckinRepo) GetByID(ctx context.Context, id uint64) (*model.CheckinRecord, error) {
	return nil, ErrCheckinNotFound
}

// === 辅助函数 ===

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
func uint64Ptr(i uint64) *uint64 { return &i }
func decimalPtr(s string) *decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return &d
}
func decimalVal(s string) decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return d
}
func decimalPtr2(s string) decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return d
}
func jsonRaw(s string) json.RawMessage { return json.RawMessage(s) }
