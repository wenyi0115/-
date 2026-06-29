// Package util 公共工具
package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// 分页默认值
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 50
)

// Pagination 分页参数
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Offset   int   `json:"-"`
}

// ParsePagination 从 gin.Context 解析分页参数
func ParsePagination(c *gin.Context) Pagination {
	page := parseInt(c.Query("page"), DefaultPage)
	pageSize := parseInt(c.Query("page_size"), DefaultPageSize)
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Pagination{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	HasMore  bool        `json:"has_more"`
}

// NewPageResult 构造分页结果
func NewPageResult(list interface{}, total int64, p Pagination) *PageResult {
	return &PageResult{
		List:     list,
		Total:    total,
		Page:     p.Page,
		PageSize: p.PageSize,
		HasMore:  int64(p.Offset+p.PageSize) < total,
	}
}

// ParseFloat 从字符串解析浮点数
func ParseFloat(s string, def float64) float64 {
	if s == "" {
		return def
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return f
}

// parseInt 解析整数（失败返回默认值）
func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

// ParseInt 从 query 解析整数
func ParseInt(c *gin.Context, key string, def int) int {
	return parseInt(c.Query(key), def)
}

// ParseInt64 从字符串解析 int64
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
