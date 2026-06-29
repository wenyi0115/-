package service

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/xunlvji/xunlvji/internal/model"
	"github.com/xunlvji/xunlvji/internal/repository"
	"github.com/xunlvji/xunlvji/pkg/util"
)

// === 分类名称映射 ===

// 攻略分类中文名
var guideCategoryNames = map[string]string{
	"nature":        "自然生态",
	"history":       "历史文化",
	"entertainment": "人工娱乐",
	"urban":         "城市公共",
	"transport":     "交通枢纽",
	"red_tourism":   "红色旅游",
	"religion":      "宗教场所",
}

// 地点分类中文名
var placeCategoryNames = map[string]string{
	"food":           "吃喝",
	"fun":            "玩乐",
	"sightseeing":    "逛看",
	"outdoor":        "户外",
	"shopping":       "逛街",
	"accommodation":  "住宿",
}

// 路线分类中文名
var routeCategoryNames = map[string]string{
	"half_day":         "半日游",
	"one_day":          "一日游",
	"two_day":          "两日游",
	"three_day_plus":   "三日及以上",
}

// 推荐等级中文名
var recommendLevelNames = map[string]string{
	"must_go": "必去",
	"can_go":  "可去",
	"avoid":   "避雷",
}

// GuideCategoryName 获取攻略分类中文名
func GuideCategoryName(code string) string {
	return guideCategoryNames[code]
}

// PlaceCategoryName 获取地点分类中文名
func PlaceCategoryName(code string) string {
	return placeCategoryNames[code]
}

// RouteCategoryName 获取路线分类中文名
func RouteCategoryName(code string) string {
	return routeCategoryNames[code]
}

// RecommendLevelName 获取推荐等级中文名
func RecommendLevelName(code string) string {
	return recommendLevelNames[code]
}

// === DTO 定义 ===

// FeedItem Feed 流卡片项（混合类型）
type FeedItem struct {
	Type           string      `json:"type"` // route/place/guide
	ID             uint64      `json:"id"`
	Title          string      `json:"title,omitempty"`
	Name           string      `json:"name,omitempty"`
	CoverImage     string      `json:"cover_image"`
	CoverImages    []string    `json:"cover_images,omitempty"`
	Category       string      `json:"category,omitempty"`
	CategoryName   string      `json:"category_name,omitempty"`
	RecommendLevel string      `json:"recommend_level,omitempty"`
	AvgCost        float64     `json:"avg_cost,omitempty"`
	Rating         float64     `json:"rating,omitempty"`
	City           string      `json:"city,omitempty"`
	Distance       float64     `json:"distance,omitempty"` // km
	TotalDuration  int         `json:"total_duration,omitempty"`
	TotalBudget    float64     `json:"total_budget,omitempty"`
	PointCount     int         `json:"point_count,omitempty"`
	SuitableFor    string      `json:"suitable_for,omitempty"`
	UsageCount     int         `json:"usage_count,omitempty"`
	FavoriteCount  int         `json:"favorite_count,omitempty"`
	VerifiedCount  int         `json:"verified_count,omitempty"`
	CheckinCount   int         `json:"checkin_count,omitempty"`
	Source         string      `json:"source,omitempty"`
	Tags           []string    `json:"tags,omitempty"`
	ReadDuration   int         `json:"read_duration,omitempty"`
	Duration       int         `json:"duration,omitempty"`
	Price          float64     `json:"price,omitempty"`
	Creator        *CreatorDTO `json:"creator,omitempty"`
	CardWidth      int         `json:"card_width,omitempty"`
	CardHeight     int         `json:"card_height,omitempty"`
	CreatedAt      string      `json:"created_at,omitempty"`
}

// CreatorDTO 创作者信息
type CreatorDTO struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// RoutePointDTO 路线点位
type RoutePointDTO struct {
	SortOrder    int     `json:"sort_order"`
	PointTime    string  `json:"point_time,omitempty"`
	PlaceID      uint64  `json:"place_id,omitempty"`
	Name         string  `json:"name"`
	StayDuration int     `json:"stay_duration,omitempty"`
	Cost         float64 `json:"cost,omitempty"`
	Category     string  `json:"category,omitempty"`
	CategoryName string  `json:"category_name,omitempty"`
	Transport    string  `json:"transport,omitempty"`
}

// RouteDetailDTO 路线详情
type RouteDetailDTO struct {
	ID            uint64          `json:"id"`
	Title         string          `json:"title"`
	CoverImages   []string        `json:"cover_images"`
	TotalDuration int             `json:"total_duration"`
	TotalBudget   float64         `json:"total_budget"`
	SuitableFor   string          `json:"suitable_for"`
	Category      string          `json:"category"`
	CategoryName  string          `json:"category_name"`
	Tags          []string        `json:"tags"`
	Rating        float64         `json:"rating"`
	UsageCount    int             `json:"usage_count"`
	FavoriteCount int             `json:"favorite_count"`
	Source        string          `json:"source"`
	Creator       *CreatorDTO     `json:"creator"`
	Points        []RoutePointDTO `json:"points"`
	IsFavorited   bool            `json:"is_favorited"`
	CommentCount  int             `json:"comment_count"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

// PlaceDetailDTO 地点详情
type PlaceDetailDTO struct {
	ID                uint64      `json:"id"`
	Name              string      `json:"name"`
	CoverImages       []string    `json:"cover_images"`
	Category          string      `json:"category"`
	CategoryName      string      `json:"category_name"`
	RecommendLevel    string      `json:"recommend_level"`
	RecommendLevelName string     `json:"recommend_level_name"`
	AvgCost           float64     `json:"avg_cost"`
	Rating            float64     `json:"rating"`
	CheckinCount      int         `json:"checkin_count"`
	Address           string      `json:"address"`
	Longitude         float64     `json:"longitude"`
	Latitude          float64     `json:"latitude"`
	Distance          float64     `json:"distance,omitempty"`
	OpeningHours      string      `json:"opening_hours"`
	IsOpen            bool        `json:"is_open"`
	Phone             string      `json:"phone"`
	SuggestDuration   int         `json:"suggest_duration"`
	Tags              []string    `json:"tags"`
	Description       string      `json:"description,omitempty"`
	IsFavorited       bool        `json:"is_favorited"`
	CommentCount      int         `json:"comment_count"`
	CreatedAt         string      `json:"created_at"`
}

// GuideDetailDTO 攻略详情
type GuideDetailDTO struct {
	ID           uint64  `json:"id"`
	Title        string  `json:"title"`
	CoverImages  []string `json:"cover_images"`
	Category     string  `json:"category"`
	CategoryName string  `json:"category_name"`
	City         string  `json:"city"`
	Duration     int     `json:"duration"`
	Price        float64 `json:"price"`
	ReadDuration int     `json:"read_duration"`
	Address      string  `json:"address"`
	OpeningHours string  `json:"opening_hours"`
	Transport    string  `json:"transport"`
	Content      string  `json:"content"`
	Highlights   string  `json:"highlights,omitempty"`
	PitfallReminders string `json:"pitfall_reminders,omitempty"`
	VerifiedCount int    `json:"verified_count"`
	Source       string  `json:"source"`
	IsFavorited  bool    `json:"is_favorited"`
	CommentCount int     `json:"comment_count"`
	CreatedAt    string  `json:"created_at"`
}

// === 内容服务 ===

// ContentService 内容服务
type ContentService struct {
	placeRepo   repository.PlaceRepository
	routeRepo   repository.RouteRepository
	guideRepo   repository.GuideRepository
	checkinRepo repository.CheckinRepository
}

// NewContentService 创建内容服务
func NewContentService(
	placeRepo repository.PlaceRepository,
	routeRepo repository.RouteRepository,
	guideRepo repository.GuideRepository,
	checkinRepo repository.CheckinRepository,
) *ContentService {
	return &ContentService{
		placeRepo:   placeRepo,
		routeRepo:   routeRepo,
		guideRepo:   guideRepo,
		checkinRepo: checkinRepo,
	}
}

// ListRecommendFeed 推荐流（混合路线/地点/攻略，简化版：按各自热度取并集）
// currentUserID: 当前登录用户ID（0=游客），用于个性化排序
func (s *ContentService) ListRecommendFeed(ctx context.Context, currentUserID uint64, p util.Pagination) (*util.PageResult, error) {
	// 简化策略：取路线(按热度)、地点(按打卡数)、攻略(按验证数) 各 page_size/3 条混合
	eachCount := p.PageSize / 3
	if eachCount < 1 {
		eachCount = 1
	}

	routes, _, err := s.routeRepo.List(ctx, repository.RouteQuery{Sort: "hot"}, 0, eachCount)
	if err != nil {
		return nil, err
	}
	places, _, err := s.placeRepo.List(ctx, repository.PlaceQuery{Sort: "popular"}, 0, eachCount)
	if err != nil {
		return nil, err
	}
	guides, _, err := s.guideRepo.List(ctx, repository.GuideQuery{Sort: "popular"}, 0, eachCount)
	if err != nil {
		return nil, err
	}

	items := make([]FeedItem, 0, len(routes)+len(places)+len(guides))
	// 交错排列：route, place, guide
	maxLen := maxInt(len(routes), len(places), len(guides))
	for i := 0; i < maxLen; i++ {
		if i < len(routes) {
			items = append(items, s.routeToFeedItem(&routes[i]))
		}
		if i < len(places) {
			items = append(items, s.placeToFeedItem(&places[i], 0))
		}
		if i < len(guides) {
			items = append(items, s.guideToFeedItem(&guides[i]))
		}
	}

	total := int64(len(items))
	return util.NewPageResult(items, total, p), nil
}

// ListGuideFeed 攻略流
func (s *ContentService) ListGuideFeed(ctx context.Context, category, city, sort string, p util.Pagination) (*util.PageResult, error) {
	list, total, err := s.guideRepo.List(ctx, repository.GuideQuery{
		Category: category,
		City:     city,
		Sort:     sort,
	}, p.Offset, p.PageSize)
	if err != nil {
		return nil, err
	}
	items := make([]FeedItem, 0, len(list))
	for i := range list {
		items = append(items, s.guideToFeedItem(&list[i]))
	}
	return util.NewPageResult(items, total, p), nil
}

// ListPlaceFeed 打卡地流
func (s *ContentService) ListPlaceFeed(ctx context.Context, category, city, sort string, p util.Pagination) (*util.PageResult, error) {
	list, total, err := s.placeRepo.List(ctx, repository.PlaceQuery{
		Category: category,
		City:     city,
		Sort:     sort,
	}, p.Offset, p.PageSize)
	if err != nil {
		return nil, err
	}
	items := make([]FeedItem, 0, len(list))
	for i := range list {
		items = append(items, s.placeToFeedItem(&list[i], 0))
	}
	return util.NewPageResult(items, total, p), nil
}

// ListRouteFeed 路线流
func (s *ContentService) ListRouteFeed(ctx context.Context, category, city, sort string, p util.Pagination) (*util.PageResult, error) {
	_ = city // 路线表无 city 字段，城市过滤后续通过点位关联实现
	list, total, err := s.routeRepo.List(ctx, repository.RouteQuery{
		Category: category,
		Sort:     sort,
	}, p.Offset, p.PageSize)
	if err != nil {
		return nil, err
	}
	items := make([]FeedItem, 0, len(list))
	for i := range list {
		item := s.routeToFeedItem(&list[i])
		// 补充点位数量
		points, _ := s.routeRepo.GetPoints(ctx, list[i].ID)
		item.PointCount = len(points)
		// 补充创作者信息
		if creator, err := s.routeRepo.GetCreator(ctx, list[i].CreatorID); err == nil {
			item.Creator = &CreatorDTO{
				ID:       creator.ID,
				Nickname: creator.Nickname,
				Avatar:   ptrToStr(creator.Avatar),
			}
		}
		items = append(items, item)
	}
	return util.NewPageResult(items, total, p), nil
}

// ListNearbyFeed 附近流
func (s *ContentService) ListNearbyFeed(ctx context.Context, locationMode, city string, lng, lat float64, radius int, p util.Pagination) (*util.PageResult, error) {
	// 校验定位模式
	switch locationMode {
	case "travel":
		if lng == 0 || lat == 0 {
			return nil, ErrLocationDenied
		}
	case "plan", "manual":
		if city == "" {
			return nil, ErrMissingCity
		}
		// plan/manual 模式：以城市中心点计算距离，这里简化为查询该城市的地点
		list, total, err := s.placeRepo.List(ctx, repository.PlaceQuery{
			City: city,
			Sort: "popular",
		}, p.Offset, p.PageSize)
		if err != nil {
			return nil, err
		}
		items := make([]FeedItem, 0, len(list))
		for i := range list {
			// 城市模式下距离设为0，或按城市中心计算（这里简化）
			items = append(items, s.placeToFeedItem(&list[i], 0))
		}
		return util.NewPageResult(items, total, p), nil
	default:
		return nil, ErrInvalidLocationMode
	}

	// travel 模式：按距离查询（简化版：查全部地点后在内存计算距离）
	// 生产环境应使用 MySQL 空间索引或高德API
	list, _, err := s.placeRepo.List(ctx, repository.PlaceQuery{
		Sort: "newest",
	}, 0, 200) // 限制查询量
	if err != nil {
		return nil, err
	}

	type placeWithDist struct {
		place model.Place
		dist  float64
	}
	withDist := make([]placeWithDist, 0, len(list))
	for _, p := range list {
		d := haversine(lat, lng, p.Latitude.InexactFloat64(), p.Longitude.InexactFloat64())
		if int(d*1000) <= radius {
			withDist = append(withDist, placeWithDist{place: p, dist: d})
		}
	}

	// 按距离排序
	for i := 0; i < len(withDist); i++ {
		for j := i + 1; j < len(withDist); j++ {
			if withDist[j].dist < withDist[i].dist {
				withDist[i], withDist[j] = withDist[j], withDist[i]
			}
		}
	}

	// 手动分页
	total := int64(len(withDist))
	start := p.Offset
	end := start + p.PageSize
	if start > len(withDist) {
		start = len(withDist)
	}
	if end > len(withDist) {
		end = len(withDist)
	}
	paged := withDist[start:end]

	items := make([]FeedItem, 0, len(paged))
	for _, wd := range paged {
		items = append(items, s.placeToFeedItem(&wd.place, roundTo1Decimal(wd.dist)))
	}
	return util.NewPageResult(items, total, p), nil
}

// GetRouteDetail 路线详情
func (s *ContentService) GetRouteDetail(ctx context.Context, routeID uint64) (*RouteDetailDTO, error) {
	route, err := s.routeRepo.GetByID(ctx, routeID)
	if err != nil {
		return nil, err
	}

	points, err := s.routeRepo.GetPoints(ctx, routeID)
	if err != nil {
		return nil, err
	}

	pointDTOs := make([]RoutePointDTO, 0, len(points))
	for _, p := range points {
		pointDTOs = append(pointDTOs, RoutePointDTO{
			SortOrder:    p.SortOrder,
			PointTime:    ptrToStr(p.PointTime),
			PlaceID:      ptrToUint64(p.PlaceID),
			Name:         p.Name,
			StayDuration: ptrToInt(p.StayDuration),
			Cost:         ptrToFloat(p.Cost),
			Category:     ptrToStr(p.Category),
			CategoryName: PlaceCategoryName(ptrToStr(p.Category)),
			Transport:    ptrToStr(p.Transport),
		})
	}

	dto := &RouteDetailDTO{
		ID:            route.ID,
		Title:         route.Title,
		CoverImages:   parseStringSlice(route.CoverImage),
		TotalDuration: route.TotalDuration,
		TotalBudget:   ptrToFloat(route.TotalBudget),
		SuitableFor:   ptrToStr(route.SuitableFor),
		Category:      route.Category,
		CategoryName:  RouteCategoryName(route.Category),
		Tags:          parseStringSlice(route.Tags),
		Rating:        route.Rating.InexactFloat64(),
		UsageCount:    route.UsageCount,
		FavoriteCount: route.FavoriteCount,
		Source:        route.Source,
		Points:        pointDTOs,
		IsFavorited:   false, // TODO: 接入收藏仓库后填充
		CommentCount:  route.CommentCount,
		CreatedAt:     route.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     route.UpdatedAt.Format(time.RFC3339),
	}

	// 补充创作者信息
	if creator, err := s.routeRepo.GetCreator(ctx, route.CreatorID); err == nil {
		dto.Creator = &CreatorDTO{
			ID:       creator.ID,
			Nickname: creator.Nickname,
			Avatar:   ptrToStr(creator.Avatar),
		}
	}

	return dto, nil
}

// GetPlaceDetail 地点详情
func (s *ContentService) GetPlaceDetail(ctx context.Context, placeID uint64, lng, lat float64) (*PlaceDetailDTO, error) {
	place, err := s.placeRepo.GetByID(ctx, placeID)
	if err != nil {
		return nil, err
	}

	distance := 0.0
	if lng != 0 && lat != 0 {
		distance = haversine(lat, lng, place.Latitude.InexactFloat64(), place.Longitude.InexactFloat64())
		distance = roundTo1Decimal(distance)
	}

	dto := &PlaceDetailDTO{
		ID:                 place.ID,
		Name:               place.Name,
		CoverImages:        parseStringSlice(place.CoverImage),
		Category:           place.Category,
		CategoryName:       PlaceCategoryName(place.Category),
		RecommendLevel:     ptrToStr(place.RecommendLevel),
		RecommendLevelName: RecommendLevelName(ptrToStr(place.RecommendLevel)),
		AvgCost:            ptrToFloat(place.AvgCost),
		Rating:             place.Rating.InexactFloat64(),
		CheckinCount:       place.CheckinCount,
		Address:            ptrToStr(place.Address),
		Longitude:          place.Longitude.InexactFloat64(),
		Latitude:           place.Latitude.InexactFloat64(),
		Distance:           distance,
		OpeningHours:       ptrToStr(place.OpeningHours),
		IsOpen:             isOpenNow(ptrToStr(place.OpeningHours)),
		Phone:              ptrToStr(place.Phone),
		SuggestDuration:    ptrToInt(place.SuggestDuration),
		Tags:               parseStringSlice(place.Tags),
		IsFavorited:        false, // TODO: 接入收藏仓库后填充
		CommentCount:       place.CommentCount,
		CreatedAt:          place.CreatedAt.Format(time.RFC3339),
	}
	return dto, nil
}

// GetGuideDetail 攻略详情
func (s *ContentService) GetGuideDetail(ctx context.Context, guideID uint64) (*GuideDetailDTO, error) {
	guide, err := s.guideRepo.GetByID(ctx, guideID)
	if err != nil {
		return nil, err
	}

	dto := &GuideDetailDTO{
		ID:               guide.ID,
		Title:            guide.Title,
		CoverImages:      parseStringSlice(guide.CoverImage),
		Category:         guide.Category,
		CategoryName:     GuideCategoryName(guide.Category),
		City:             guide.City,
		Duration:         ptrToInt(guide.Duration),
		Price:            ptrToFloat(guide.Price),
		ReadDuration:     ptrToInt(guide.ReadDuration),
		Address:          ptrToStr(guide.Address),
		OpeningHours:     ptrToStr(guide.OpeningHours),
		Transport:        ptrToStr(guide.Transport),
		Content:          guide.Content,
		Highlights:       ptrToStr(guide.Highlights),
		PitfallReminders: ptrToStr(guide.PitfallReminders),
		VerifiedCount:    guide.VerifiedCount,
		Source:           guide.Source,
		IsFavorited:      false, // TODO: 接入收藏仓库后填充
		CommentCount:     guide.CommentCount,
		CreatedAt:        guide.CreatedAt.Format(time.RFC3339),
	}
	return dto, nil
}

// === 转换辅助函数 ===

func (s *ContentService) routeToFeedItem(r *model.Route) FeedItem {
	item := FeedItem{
		Type:          "route",
		ID:            r.ID,
		Title:         r.Title,
		CoverImage:    firstString(parseStringSlice(r.CoverImage)),
		CoverImages:   parseStringSlice(r.CoverImage),
		TotalDuration: r.TotalDuration,
		TotalBudget:   ptrToFloat(r.TotalBudget),
		SuitableFor:   ptrToStr(r.SuitableFor),
		Category:      r.Category,
		CategoryName:  RouteCategoryName(r.Category),
		Rating:        r.Rating.InexactFloat64(),
		UsageCount:    r.UsageCount,
		FavoriteCount: r.FavoriteCount,
		Source:        r.Source,
		Tags:          parseStringSlice(r.Tags),
		CardWidth:     320,
		CardHeight:    240,
		CreatedAt:     r.CreatedAt.Format(time.RFC3339),
	}
	return item
}

func (s *ContentService) placeToFeedItem(p *model.Place, distance float64) FeedItem {
	item := FeedItem{
		Type:           "place",
		ID:             p.ID,
		Name:           p.Name,
		CoverImage:     firstString(parseStringSlice(p.CoverImage)),
		CoverImages:    parseStringSlice(p.CoverImage),
		Category:       p.Category,
		CategoryName:   PlaceCategoryName(p.Category),
		RecommendLevel: ptrToStr(p.RecommendLevel),
		AvgCost:        ptrToFloat(p.AvgCost),
		Rating:         p.Rating.InexactFloat64(),
		City:           p.City,
		Distance:       distance,
		CheckinCount:   p.CheckinCount,
		Tags:           parseStringSlice(p.Tags),
		CardWidth:      320,
		CardHeight:     280,
		CreatedAt:      p.CreatedAt.Format(time.RFC3339),
	}
	return item
}

func (s *ContentService) guideToFeedItem(g *model.Guide) FeedItem {
	item := FeedItem{
		Type:          "guide",
		ID:            g.ID,
		Title:         g.Title,
		CoverImage:    firstString(parseStringSlice(g.CoverImage)),
		CoverImages:   parseStringSlice(g.CoverImage),
		Category:      g.Category,
		CategoryName:  GuideCategoryName(g.Category),
		City:          g.City,
		Duration:      ptrToInt(g.Duration),
		Price:         ptrToFloat(g.Price),
		ReadDuration:  ptrToInt(g.ReadDuration),
		VerifiedCount: g.VerifiedCount,
		Source:        g.Source,
		CardWidth:     380,
		CardHeight:    200,
		CreatedAt:     g.CreatedAt.Format(time.RFC3339),
	}
	return item
}

// === 通用辅助函数 ===

// parseStringSlice 解析 JSON 字符串数组
func parseStringSlice(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var arr []string
	if err := json.Unmarshal(raw, &arr); err != nil {
		return []string{}
	}
	return arr
}

// firstString 取切片第一个元素
func firstString(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	return arr[0]
}

// ptrToInt 指针转 int
func ptrToInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// ptrToFloat decimal.Decimal 指针转 float64
func ptrToFloat(p *decimal.Decimal) float64 {
	if p == nil {
		return 0
	}
	return p.InexactFloat64()
}

// ptrToUint64 指针转 uint64
func ptrToUint64(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}

// roundTo1Decimal 保留一位小数
func roundTo1Decimal(f float64) float64 {
	return math.Round(f*10) / 10
}

// maxInt 取多个整数的最大值
func maxInt(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	m := nums[0]
	for _, n := range nums[1:] {
		if n > m {
			m = n
		}
	}
	return m
}

// haversine 计算两点间距离（单位：km）
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const r = 6371 // 地球半径 km
	dLat := toRad(lat2 - lat1)
	dLng := toRad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return r * c
}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// isOpenNow 简单判断当前是否营业（格式 "08:00-18:00"）
func isOpenNow(hours string) bool {
	if hours == "" || !strings.Contains(hours, "-") {
		return false
	}
	parts := strings.SplitN(hours, "-", 2)
	if len(parts) != 2 {
		return false
	}
	now := time.Now()
	nowMin := now.Hour()*60 + now.Minute()
	startMin := parseTimeToMinutes(strings.TrimSpace(parts[0]))
	endMin := parseTimeToMinutes(strings.TrimSpace(parts[1]))
	if startMin < 0 || endMin < 0 {
		return false
	}
	return nowMin >= startMin && nowMin <= endMin
}

func parseTimeToMinutes(s string) int {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return -1
	}
	h := atoiSafe(parts[0])
	m := atoiSafe(parts[1])
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return -1
	}
	return h*60 + m
}

func atoiSafe(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}

// 内容服务错误
var (
	ErrLocationDenied     = errors.New("location denied")
	ErrMissingCity        = errors.New("missing city")
	ErrInvalidLocationMode = errors.New("invalid location mode")
)
