// Package wechat 微信小程序 SDK 封装
package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/xunlvji/xunlvji/internal/config"
)

// 微信小程序 code2session 接口地址
const code2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"

// 默认 HTTP 超时
const defaultTimeout = 10 * time.Second

var (
	// ErrInvalidCode 无效的 code
	ErrInvalidCode = errors.New("invalid code")
	// ErrInvalidAppID 无效的 appid
	ErrInvalidAppID = errors.New("invalid appid")
	// ErrAPICallFailed 微信 API 调用失败
	ErrAPICallFailed = errors.New("wechat api call failed")
)

// Code2SessionResponse 微信 code2session 接口响应
type Code2SessionResponse struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识（仅在该小程序已绑定到微信开放平台账号时返回）
	ErrCode    int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

// Client 微信小程序客户端
type Client struct {
	appID     string
	appSecret string
	httpCli   *http.Client
}

// NewClient 创建微信小程序客户端
func NewClient(cfg config.WeChatConfig) *Client {
	return &Client{
		appID:     cfg.AppID,
		appSecret: cfg.AppSecret,
		httpCli:   &http.Client{Timeout: defaultTimeout},
	}
}

// Code2Session code 换取 openid 和 session_key
// 文档：https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func (c *Client) Code2Session(ctx context.Context, code string) (*Code2SessionResponse, error) {
	if c.appID == "" || c.appSecret == "" {
		return nil, ErrInvalidAppID
	}
	if code == "" {
		return nil, ErrInvalidCode
	}

	params := url.Values{}
	params.Set("appid", c.appID)
	params.Set("secret", c.appSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	reqURL := code2SessionURL + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAPICallFailed, err)
	}

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAPICallFailed, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: read body: %v", ErrAPICallFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: http status %d, body: %s", ErrAPICallFailed, resp.StatusCode, string(body))
	}

	var result Code2SessionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("%w: unmarshal: %v", ErrAPICallFailed, err)
	}

	// 业务错误码处理
	// 常见错误：40029 invalid code、40011 invalid appid、45011 频率限制
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("%w: errcode=%d errmsg=%s", ErrAPICallFailed, result.ErrCode, result.ErrMsg)
	}

	if result.OpenID == "" {
		return nil, fmt.Errorf("%w: empty openid", ErrAPICallFailed)
	}

	return &result, nil
}
