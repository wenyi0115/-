package errcode

// 错误码定义（见 18_错误码规范.md）
// 0=成功，非0=错误

// 通用错误 (1xxx)
const (
	Success              = 0
	ErrInvalidParam      = 1001 // 参数错误
	ErrMissingParam      = 1002 // 缺少必要参数
	ErrInvalidFormat     = 1003 // 参数格式不正确
	ErrNotFound          = 1004 // 数据不存在
	ErrAlreadyExists     = 1005 // 数据已存在
	ErrTooManyRequests   = 1006 // 操作过于频繁
	ErrRequestExpired    = 1007 // 请求已过期
	ErrVersionTooLow     = 1008 // 版本过低
	ErrMaintenance       = 1009 // 服务维护中
)

// 认证与授权 (2xxx)
const (
	ErrNotLoggedIn        = 2001 // 未登录
	ErrTokenExpired       = 2002 // Token已过期
	ErrTokenInvalid       = 2003 // Token无效
	ErrNoPermission       = 2004 // 无操作权限
	ErrAccountBanned      = 2005 // 账号已被封禁
	ErrCaptchaWrong       = 2006 // 验证码错误
	ErrCaptchaExpired     = 2007 // 验证码已过期
	ErrCaptchaTooFrequent = 2008 // 验证码发送过于频繁
	ErrPhoneRegistered    = 2009 // 手机号已注册
	ErrPhoneNotRegistered = 2010 // 手机号未注册
	ErrPasswordWrong      = 2011 // 密码错误
	ErrOAuthFailed        = 2012 // 第三方登录授权失败
	ErrAccountDeleted     = 2013 // 账号已注销
	ErrDeviceAbnormal     = 2014 // 设备登录异常
	ErrRefreshTokenExpired = 2015 // Refresh Token已过期
)

// 用户相关 (3xxx)
const (
	ErrUserNotFound       = 3001 // 用户不存在
	ErrNicknameTaken      = 3002 // 昵称已被占用
	ErrNicknameSensitive  = 3003 // 昵称包含敏感词
	ErrAvatarUploadFailed = 3004 // 头像上传失败
	ErrBioTooLong         = 3005 // 个人简介过长
	ErrOperationSelf      = 3006 // 操作自己不允许
	ErrCreditScoreLow     = 3007 // 信用分不足
	ErrDeleteCooldown     = 3008 // 注销冷静期未结束
)

// 内容相关 (4xxx)
const (
	ErrPlaceNotFound     = 4001 // 地点不存在
	ErrRouteNotFound     = 4002 // 路线不存在
	ErrGuideNotFound     = 4003 // 攻略不存在
	ErrCheckinNotFound   = 4004 // 打卡记录不存在
	ErrContentRejected   = 4005 // 内容审核已驳回
	ErrContentReviewing  = 4006 // 内容审核中
	ErrContentViolation  = 4007 // 内容包含违规信息
	ErrImageLimit        = 4008 // 图片数量超限
	ErrCheckinDuplicate  = 4009 // 24小时内已打卡该地点
	ErrPublishTooFast    = 4010 // 发布频率过高
	ErrContentDeleted    = 4011 // 内容已被删除
	ErrTagLimit          = 4012 // 标签数量超限
	ErrPlaceNameExists   = 4013 // 地点名称已存在
)

// 行程相关 (5xxx)
const (
	ErrTripNotFound       = 5001 // 行程不存在
	ErrTripPointNotFound  = 5002 // 行程点位不存在
	ErrTripStatusInvalid  = 5003 // 行程状态不允许此操作
	ErrTripPointLimit     = 5004 // 行程点位已达上限
	ErrTripPointExists    = 5005 // 行程已包含该地点
	ErrTripInProgress     = 5006 // 不能删除进行中的行程
	ErrTripCompleted      = 5009 // 行程已完成，不可修改
	ErrSkipNeedReason     = 5010 // 跳过点位需要提供理由
)

// 社交相关 (6xxx)
const (
	ErrCommentNotFound   = 6001 // 评论不存在
	ErrCommentDeleted    = 6002 // 评论已被删除
	ErrCommentEmpty      = 6003 // 评论内容为空
	ErrCommentTooLong    = 6004 // 评论内容过长
	ErrCommentSensitive  = 6005 // 评论包含敏感词
	ErrCommentTooFast    = 6006 // 评论过于频繁
	ErrAlreadyFollowed   = 6007 // 已关注该用户
	ErrNotFollowed       = 6008 // 未关注该用户
	ErrFollowSelf        = 6009 // 不能关注自己
	ErrAlreadyFavorited  = 6010 // 已收藏该内容
	ErrAlreadyLiked      = 6011 // 已点赞
	ErrAlreadyVoted      = 6012 // 已投"有用"票
	ErrFollowLimit       = 6013 // 关注数已达上限
	ErrReplyNotFound     = 6014 // 回复的评论不存在
)

// 文件上传 (7xxx)
const (
	ErrFileTypeNotSupport   = 7001 // 文件类型不支持
	ErrFileSizeExceed       = 7002 // 文件大小超限
	ErrFileUploadFailed     = 7003 // 文件上传失败
	ErrFileCountExceed      = 7004 // 文件数量超限
	ErrFileFormatNotSupport = 7005 // 文件格式不支持
	ErrImageResolutionLow   = 7006 // 图片分辨率过低
	ErrImageViolation       = 7007 // 图片包含违规内容
)

// AI服务 (8xxx)
const (
	ErrAIUnavailable     = 8001 // AI服务暂时不可用
	ErrAITimeout         = 8002 // AI生成超时
	ErrAIContentInvalid  = 8003 // AI生成内容不符合规范
	ErrAIFailed          = 8004 // AI生成失败
	ErrAINeedReview      = 8005 // AI内容需要人工审核
	ErrAISortFailed      = 8006 // AI排序失败
	ErrAIChatTimeout     = 8007 // AI对话超时
	ErrAIPlanInvalid     = 8008 // AI规划结果无效
	ErrAIRecognizeFailed = 8009 // AI图片识别失败
	ErrAIReviewReject    = 8010 // AI内容审核不通过
	ErrAIQuotaExhausted  = 8011 // AI API额度不足
	ErrDoubaoTimeout     = 8012 // 豆包API超时
	ErrDoubaoAbnormal    = 8013 // 豆包API返回异常
	ErrAIReviewFailed    = 8014 // AI内容审核失败
)

// 系统错误 (9xxx)
const (
	ErrInternal       = 9001 // 服务器内部错误
	ErrDatabase       = 9002 // 数据库错误
	ErrCache          = 9003 // 缓存服务异常
	ErrSearch         = 9004 // 搜索服务异常
	ErrMessageQueue   = 9005 // 消息队列异常
	ErrThirdParty     = 9006 // 第三方服务异常
	ErrRateLimit      = 9007 // 限流保护
	ErrDegraded       = 9008 // 服务降级中
)

// CPS分销相关 (15000-15099)
const (
	ErrCPSLinkFailed    = 15001 // CPS链接生成失败
	ErrCPSTimeout       = 15002 // CPS平台API超时
	ErrCPSLinkExpired   = 15003 // CPS链接已失效
	ErrCPSNoMatch       = 15004 // CPS匹配无结果
	ErrCPSNotAuthorized = 15005 // CPS平台未授权
)

// 定位相关 (16000-16099)
const (
	ErrLocationDenied   = 16001 // 定位权限被拒绝
	ErrLocationTimeout  = 16002 // 定位超时
	ErrGPSWeak          = 16003 // GPS信号弱
	ErrGeocodeFailed    = 16004 // 逆地理编码失败
	ErrPOINotFound      = 16005 // POI搜索无结果
)

// 行程工具相关 (18000-18099)
const (
	ErrListGenFailed    = 18001 // 清单生成失败
	ErrBillSummaryFail  = 18002 // 账单汇总失败
	ErrMapLoadFailed    = 18003 // 地图加载失败
	ErrWeatherFailed    = 18004 // 天气获取失败
)

// Msg 错误码对应消息
var Msg = map[int]string{
	Success:              "success",
	ErrInvalidParam:      "参数错误",
	ErrMissingParam:      "缺少必要参数",
	ErrInvalidFormat:     "参数格式不正确",
	ErrNotFound:          "数据不存在",
	ErrAlreadyExists:     "数据已存在",
	ErrTooManyRequests:   "操作过于频繁",
	ErrRequestExpired:    "请求已过期",
	ErrVersionTooLow:     "版本过低，请升级App",
	ErrMaintenance:       "服务维护中",

	ErrNotLoggedIn:        "未登录",
	ErrTokenExpired:       "Token已过期",
	ErrTokenInvalid:       "Token无效",
	ErrNoPermission:       "无操作权限",
	ErrAccountBanned:      "账号已被封禁",
	ErrCaptchaWrong:       "验证码错误",
	ErrCaptchaExpired:     "验证码已过期",
	ErrCaptchaTooFrequent: "验证码发送过于频繁",
	ErrPhoneRegistered:    "手机号已注册",
	ErrPhoneNotRegistered: "手机号未注册",
	ErrPasswordWrong:      "密码错误",
	ErrOAuthFailed:        "第三方登录授权失败",
	ErrAccountDeleted:     "账号已注销",
	ErrDeviceAbnormal:     "设备登录异常",
	ErrRefreshTokenExpired: "Refresh Token已过期",

	ErrUserNotFound:       "用户不存在",
	ErrNicknameTaken:      "昵称已被占用",
	ErrNicknameSensitive:  "昵称包含敏感词",
	ErrAvatarUploadFailed: "头像上传失败",
	ErrBioTooLong:         "个人简介过长（限200字）",
	ErrOperationSelf:      "操作自己不允许",
	ErrCreditScoreLow:     "信用分不足，无法操作",
	ErrDeleteCooldown:     "注销冷静期（30天）未结束",

	ErrPlaceNotFound:    "地点不存在",
	ErrRouteNotFound:    "路线不存在",
	ErrGuideNotFound:    "攻略不存在",
	ErrCheckinNotFound:  "打卡记录不存在",
	ErrContentRejected:  "内容审核已驳回",
	ErrContentReviewing: "内容审核中",
	ErrContentViolation: "内容包含违规信息",
	ErrImageLimit:       "图片数量超限（最多9张）",
	ErrCheckinDuplicate: "24小时内已打卡该地点",
	ErrPublishTooFast:   "发布频率过高（每小时最多5条）",
	ErrContentDeleted:   "内容已被删除",
	ErrTagLimit:         "标签数量超限（最多5个）",
	ErrPlaceNameExists:  "地点名称已存在",

	ErrTripNotFound:      "行程不存在",
	ErrTripPointNotFound: "行程点位不存在",
	ErrTripStatusInvalid: "行程状态不允许此操作",
	ErrTripPointLimit:    "行程点位已达上限（最多20个）",
	ErrTripPointExists:   "行程已包含该地点",
	ErrTripInProgress:    "不能删除进行中的行程",
	ErrTripCompleted:     "行程已完成，不可修改",
	ErrSkipNeedReason:    "跳过点位需要提供理由",

	ErrCommentNotFound:  "评论不存在",
	ErrCommentDeleted:   "评论已被删除",
	ErrCommentEmpty:     "评论内容为空",
	ErrCommentTooLong:   "评论内容过长（限500字）",
	ErrCommentSensitive: "评论包含敏感词",
	ErrCommentTooFast:   "评论过于频繁（每分钟最多3条）",
	ErrAlreadyFollowed:  "已关注该用户",
	ErrNotFollowed:      "未关注该用户",
	ErrFollowSelf:       "不能关注自己",
	ErrAlreadyFavorited: "已收藏该内容",
	ErrAlreadyLiked:     "已点赞",
	ErrAlreadyVoted:     "已投\"有用\"票",
	ErrFollowLimit:      "关注数已达上限（5000）",
	ErrReplyNotFound:    "回复的评论不存在",

	ErrFileTypeNotSupport:   "文件类型不支持",
	ErrFileSizeExceed:       "文件大小超限",
	ErrFileUploadFailed:     "文件上传失败，请重试",
	ErrFileCountExceed:      "文件数量超限（最多9张）",
	ErrFileFormatNotSupport: "文件格式不支持（仅支持JPG/PNG/WebP）",
	ErrImageResolutionLow:   "图片分辨率过低（最小200×200）",
	ErrImageViolation:       "图片包含违规内容",

	ErrAIUnavailable:     "AI服务暂时不可用",
	ErrAITimeout:         "AI生成超时",
	ErrAIContentInvalid:  "AI生成内容不符合规范",
	ErrAIFailed:          "AI生成失败，请重试",
	ErrAINeedReview:      "AI内容需要人工审核",
	ErrAISortFailed:      "AI排序失败",
	ErrAIChatTimeout:     "AI对话超时",
	ErrAIPlanInvalid:     "AI规划结果无效",
	ErrAIRecognizeFailed: "AI图片识别失败",
	ErrAIReviewReject:    "AI内容审核不通过",
	ErrAIQuotaExhausted:  "AI API额度不足",
	ErrDoubaoTimeout:     "豆包API超时",
	ErrDoubaoAbnormal:    "豆包API返回异常",
	ErrAIReviewFailed:    "AI内容审核失败",

	ErrInternal:     "服务器内部错误",
	ErrDatabase:     "数据库错误",
	ErrCache:        "缓存服务异常",
	ErrSearch:       "搜索服务异常",
	ErrMessageQueue: "消息队列异常",
	ErrThirdParty:   "第三方服务异常",
	ErrRateLimit:    "限流保护",
	ErrDegraded:     "服务降级中",

	ErrCPSLinkFailed:    "CPS链接生成失败",
	ErrCPSTimeout:       "CPS平台API超时",
	ErrCPSLinkExpired:   "CPS链接已失效",
	ErrCPSNoMatch:       "CPS匹配无结果",
	ErrCPSNotAuthorized: "CPS平台未授权",

	ErrLocationDenied:  "定位权限被拒绝",
	ErrLocationTimeout: "定位超时",
	ErrGPSWeak:         "GPS信号弱",
	ErrGeocodeFailed:   "逆地理编码失败",
	ErrPOINotFound:     "POI搜索无结果",

	ErrListGenFailed:   "清单生成失败",
	ErrBillSummaryFail: "账单汇总失败",
	ErrMapLoadFailed:   "地图加载失败",
	ErrWeatherFailed:   "天气获取失败",
}

// GetMsg 获取错误码消息
func GetMsg(code int) string {
	if msg, ok := Msg[code]; ok {
		return msg
	}
	return "未知错误"
}
