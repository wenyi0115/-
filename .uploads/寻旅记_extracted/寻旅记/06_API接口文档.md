# 寻旅记 -- API 接口文档

> 版本：v2.1.7 | 日期：2026年6月 | 状态：正式版
> 平台：iOS App / Android App / 微信小程序

---

## 目录

- [一、通用规范](#一通用规范)
- [二、用户模块](#二用户模块)
- [三、内容模块](#三内容模块)
- [四、行程模块](#四行程模块)
- [五、发布模块](#五发布模块)
- [六、社交模块](#六社交模块)
- [七、攻略问答模块](#七攻略问答模块)
- [八、其他模块](#八其他模块)

---

## 一、通用规范

### 1.1 接口前缀

```
/api/v1/
```

### 1.2 统一返回格式

所有接口统一返回 JSON，结构如下：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码，0 表示成功，非 0 表示错误 |
| message | string | 提示信息 |
| data | object/array/null | 返回数据 |
| request_id | string | 请求追踪 ID，用于日志排查（UUID 格式） |

### 1.3 分页规范

分页接口统一使用以下参数和返回格式：

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码，从 1 开始 |
| page_size | int | 否 | 20 | 每页数量，最大 50 |

**返回格式**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "has_more": true
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| list | array | 数据列表 |
| total | int | 总记录数 |
| page | int | 当前页码 |
| page_size | int | 每页数量 |
| has_more | bool | 是否还有更多数据 |

### 1.4 认证说明

需要登录的接口在标题旁标注 `@auth`。客户端需在请求头中携带 JWT Token：

```
Authorization: Bearer {token}
```

### 1.5 核心错误码

> 完整错误码定义参见 `18_错误码规范.md`，以下为常用错误码速查。

| 错误码 | 说明 | 分类 |
|--------|------|------|
| 0 | 成功 | - |
| 1001 | 参数错误 | 通用 |
| 1002 | 缺少必要参数 | 通用 |
| 1003 | 参数格式不正确 | 通用 |
| 1004 | 数据不存在 | 通用 |
| 1005 | 数据已存在 | 通用 |
| 1006 | 操作过于频繁 | 通用 |
| 2001 | 未登录 | 认证 |
| 2002 | Token 已过期 | 认证 |
| 2003 | Token 无效 | 认证 |
| 2004 | 无操作权限 | 认证 |
| 2005 | 账号已被封禁 | 认证 |
| 2006 | 验证码错误 | 认证 |
| 2007 | 验证码已过期 | 认证 |
| 2008 | 验证码发送过于频繁 | 认证 |
| 2009 | 手机号已注册 | 认证 |
| 2010 | 手机号未注册 | 认证 |
| 2011 | 密码错误 | 认证 |
| 2012 | 第三方登录授权失败 | 认证 |
| 3001 | 用户不存在 | 用户 |
| 3002 | 昵称已被占用 | 用户 |
| 3006 | 操作自己不允许 | 用户 |
| 3007 | 信用分不足，无法操作 | 用户 |
| 4001 | 地点不存在 | 内容 |
| 4002 | 路线不存在 | 内容 |
| 4003 | 攻略不存在 | 内容 |
| 4004 | 打卡记录不存在 | 内容 |
| 4005 | 内容审核已驳回 | 内容 |
| 4007 | 内容包含违规信息 | 内容 |
| 4009 | 24小时内已打卡该地点，24小时后可再次打卡 | 内容 |
| 5001 | 行程不存在 | 行程 |
| 5002 | 行程点位不存在 | 行程 |
| 5003 | 行程状态不允许此操作 | 行程 |
| 5005 | 行程已包含该地点 | 行程 |
| 5009 | 行程已完成，不可修改 | 行程 |
| 6001 | 评论不存在 | 社交 |
| 6003 | 评论内容为空 | 社交 |
| 6007 | 已关注该用户 | 社交 |
| 6008 | 未关注该用户 | 社交 |
| 6009 | 不能关注自己 | 社交 |
| 6010 | 已收藏该内容 | 社交 |
| 6011 | 已点赞 | 社交 |
| 7001 | 文件类型不支持 | 文件 |
| 7002 | 文件大小超限 | 文件 |
| 7003 | 文件上传失败，请重试 | 文件 |
| 8001 | AI 服务暂时不可用 | AI |
| 8004 | AI 生成失败，请重试 | AI |
| 9001 | 服务器内部错误 | 系统 |

### 1.6 通用枚举值

**打卡类型（checkin_type）**：

| 值 | 说明 |
|----|------|
| store | 到店打卡（GPS 验证后） |
| normal | 普通打卡（仅标记完成） |

**行程状态（trip_status）**：

| 值 | 说明 |
|----|------|
| pending | 待出发 |
| active | 进行中 |
| completed | 已完成 |

**打卡状态（checkin_status）**：

| 值 | 说明 |
|----|------|
| pending | 待打卡 |
| checked | 已打卡 |
| skipped | 已跳过 |

**推荐等级（recommend_level）**：

| 值 | 说明 |
|----|------|
| must_go | 必去 |
| can_go | 可去 |
| avoid | 避雷 |

**内容来源（source）**：

> `source` 在不同业务对象中含义不同，使用时注意区分：

| 业务对象 | 值 | 说明 |
|---------|------|------|
| 路线（Route）/ 攻略（Guide） | user | 用户创作 |
| 路线（Route）/ 攻略（Guide） | AI | AI 生成 |
| 路线（Route）/ 攻略（Guide） | official | 官方发布 |
| 打卡记录（CheckInRecord） | trip_checkin | 行程到店打卡产生 |
| 打卡记录（CheckInRecord） | manual_publish | 手动发布产生 |

**可见性（visibility）**：

| 值 | 说明 |
|----|------|
| public | 公开 |
| private | 仅自己可见 |

**地点分类（place_category）**：

| 值 | 说明 |
|----|------|
| food | 吃喝 |
| fun | 玩乐 |
| sightseeing | 逛看 |
| outdoor | 户外 |
| shopping | 逛街 |
| accommodation | 住宿 |

**路线分类（route_category）**：

| 值 | 说明 |
|----|------|
| half_day | 半日游 |
| one_day | 一日游 |
| two_day | 两日游 |
| three_day_plus | 三日及以上 |

**攻略分类（guide_category）**：

| 值 | 说明 |
|----|------|
| nature | 自然生态 |
| history | 历史文化 |
| entertainment | 人工娱乐 |
| urban | 城市公共 |
| transport | 交通枢纽 |
| red_tourism | 红色旅游 |
| religion | 宗教场所 |

### 1.7 接口限频策略（P1-6）

> 为防止接口被滥用、刷量及控制 AI/存储等高成本调用，对所有写接口及高频读接口实施限频。限频基于 Redis 滑动窗口实现，维度为「用户 ID + 接口 + 时间窗口」；未登录接口按「IP + 接口 + 时间窗口」限频。命中限频返回 HTTP 429，业务码 `1006 操作过于频繁`（见 1.5）。

#### 1.7.1 高频端点限频阈值

| 端点 | 接口编号 | 方法 & 路径 | 限频阈值（每用户） | 限频说明 |
|------|----------|-------------|--------------------|----------|
| AI 行程规划 | 4.11 | POST /api/v1/trips/ai-plan | 10 次/天 | 单次调用成本高（豆包 Pro），按天限制避免滥用与成本失控 |
| AI 攻略生成 | （规划中） | POST /api/v1/guides/ai-generate | 10 次/天 | AI 攻略生成能力（按需生成结构化攻略），与 AI 行程规划共享每日 AI 配额上限 |
| 批量图片上传 | 5.1 | POST /api/v1/upload/image | 5 次/小时 | 单次支持批量上传，按小时限制防止存储滥用与带宽占用 |
| 评论发表 | 6.8 | POST /api/v1/social/comment | 30 次/小时 | 防刷评/水军，正常用户互动频次远低于此阈值 |
| 打卡提交 | 4.5 / 4.6 | POST /api/v1/trips/{trip_id}/checkin/store<br>POST /api/v1/trips/{trip_id}/checkin/normal | 20 次/天 | 防虚假刷打卡，含到店打卡与普通打卡合计计数 |

#### 1.7.2 通用限频规则

| 接口类型 | 默认阈值 | 说明 |
|----------|----------|------|
| 验证码发送（2.1） | 1 次/60 秒，5 次/天 | 防短信轰炸（已有 2008 错误码） |
| 登录接口（2.2-2.6） | 5 次/分钟 | 防撞库爆破 |
| 通用写接口 | 60 次/分钟 | 兜底防护，未单独配置的写接口适用 |
| 通用读接口 | 120 次/分钟 | 兜底防护，未单独配置的读接口适用 |

#### 1.7.3 限频响应

```json
{
  "code": 1006,
  "message": "操作过于频繁，请稍后再试",
  "data": {
    "retry_after": 60,
    "limit": 10,
    "remaining": 0,
    "window": "day"
  }
}
```

**响应头**：
- `X-RateLimit-Limit`：窗口内总配额
- `X-RateLimit-Remaining`：窗口内剩余配额
- `X-RateLimit-Reset`：配额重置时间（Unix 秒）
- `Retry-After`：建议重试等待秒数

**说明**：
- AI 类接口（AI 行程规划、AI 攻略生成）超限返回 `8001 AI 服务暂时不可用` 或 `1006`，并提示用户当日 AI 配额已用完，次日重置（按自然日 00:00 重置）。
- 信用分 ≥ 90 的优质创作者可申请提升评论/上传配额（见 04_运营与合规手册 2.4 创作者激励体系）。
- 限频维度优先按用户 ID；未登录场景（如发送验证码）按 IP 限频，IP 限频阈值更严格以防代理刷量。

---

## 二、用户模块

### 2.1 发送验证码

```
POST /api/v1/auth/send-code
```

**描述**：向手机号发送短信验证码，用于登录/注册、换绑手机号、注销等场景的身份验证。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号，格式：11 位中国大陆手机号 |
| scene | string | 否 | 验证场景：`login`（登录/注册，默认）/ `bind`（换绑手机号）/ `delete_account`（注销账号）。不同场景的验证码互不通用 |

**请求示例**：

```json
{
  "phone": "13800138000",
  "scene": "login"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "expire_seconds": 300
  }
}
```

**说明**：
- 同一手机号 60 秒内只能发送一次，超出返回错误码 `2008`。
- 验证码有效期 5 分钟。
- 微信小程序端无需此接口，使用微信授权登录。

---

### 2.2 手机号验证码登录

```
POST /api/v1/auth/login/phone
```

**描述**：使用手机号 + 验证码登录。新用户自动注册。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号 |
| code | string | 是 | 短信验证码 |

**请求示例**：

```json
{
  "phone": "13800138000",
  "code": "123456"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expire_at": "2026-06-24T14:00:00Z",
    "refresh_expire_at": "2026-07-24T12:00:00Z",
    "user": {
      "id": 10001,
      "nickname": "旅行者_10001",
      "avatar": null,
      "phone": "138****8000",
      "credit_score": 100,
      "is_new_user": true
    }
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- `is_new_user` 为 `true` 时表示新注册用户，客户端可引导设置昵称和头像。
- Access Token 有效期 2 小时，Refresh Token 有效期 30 天。
- Access Token 过期后，客户端使用 Refresh Token 调用 `POST /api/v1/auth/refresh-token` 获取新的 Access Token。

---

### 2.3 Apple ID 登录

```
POST /api/v1/auth/login/apple
```

**描述**：Apple ID 登录。新用户自动注册。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| identity_token | string | 是 | Apple 返回的 identityToken |
| authorization_code | string | 是 | Apple 返回的 authorizationCode |
| user_identifier | string | 否 | Apple 用户标识（用于关联已有账号） |
| given_name | string | 否 | 用户名字（首次登录时 Apple 提供） |
| family_name | string | 否 | 用户姓氏（首次登录时 Apple 提供） |

**请求示例**：

```json
{
  "identity_token": "eyJraWQiOi...",
  "authorization_code": "c5f8a9...",
  "user_identifier": "001234.abcdef...",
  "given_name": "小明",
  "family_name": "王"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expire_at": "2026-06-24T14:00:00Z",
    "refresh_expire_at": "2026-07-24T12:00:00Z",
    "user": {
      "id": 10002,
      "nickname": "小明王",
      "avatar": null,
      "email": "w****@privaterelay.appleid.com",
      "credit_score": 100,
      "is_new_user": true
    }
  }
}
```

---

### 2.4 Google 登录

```
POST /api/v1/auth/login/google
```

**描述**：Google 账号登录。新用户自动注册。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id_token | string | 是 | Google 返回的 id_token |
| access_token | string | 是 | Google 返回的 access_token |

**请求示例**：

```json
{
  "id_token": "eyJhbGciOi...",
  "access_token": "ya29.a0Af..."
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expire_at": "2026-06-24T14:00:00Z",
    "refresh_expire_at": "2026-07-24T12:00:00Z",
    "user": {
      "id": 10003,
      "nickname": "John",
      "avatar": "https://lh3.googleusercontent.com/...",
      "email": "j****@gmail.com",
      "credit_score": 100,
      "is_new_user": true
    }
  }
}
```

---

### 2.5 Email 登录

```
POST /api/v1/auth/login/email
```

**描述**：邮箱 + 密码登录。仅用于已注册用户，不支持自动注册（需先通过手机号注册后绑定邮箱）。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码，8-32 位，需包含字母和数字 |

**请求示例**：

```json
{
  "email": "user@example.com",
  "password": "mypassword123"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expire_at": "2026-06-24T14:00:00Z",
    "refresh_expire_at": "2026-07-24T12:00:00Z",
    "user": {
      "id": 10001,
      "nickname": "旅行者_10001",
      "avatar": "https://cdn.example.com/avatars/10001.jpg",
      "email": "u****@example.com",
      "credit_score": 95,
      "is_new_user": false
    }
  }
}
```

---

### 2.6 微信登录

```
POST /api/v1/auth/login/wechat
```

**描述**：微信授权登录。小程序端使用 `wx.login()` 获取 code，App 端使用微信 SDK 获取授权信息。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| code | string | 是 | 微信登录凭证（wx.login 返回） |
| platform | string | 是 | 平台标识：miniprogram（小程序）/ app（App） |
| iv | string | 否 | 加密算法的初始向量（App 端） |
| encrypted_data | string | 否 | 加密数据（App 端） |

**请求示例**：

```json
{
  "code": "081xYz0w3...",
  "platform": "miniprogram"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expire_at": "2026-06-24T14:00:00Z",
    "refresh_expire_at": "2026-07-24T12:00:00Z",
    "user": {
      "id": 10004,
      "nickname": "微信用户",
      "avatar": "https://thirdwx.qlogo.cn/...",
      "credit_score": 100,
      "is_new_user": true
    }
  }
}
```

---

### 2.7 获取用户信息 @auth

```
GET /api/v1/user/profile
```

**描述**：获取当前登录用户的完整信息。

**请求参数**：无（从 Token 中获取用户 ID）。

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 10001,
    "nickname": "旅行者小明",
    "avatar": "https://cdn.example.com/avatars/10001.jpg",
    "phone": "138****8000",
    "email": "u****@example.com",
    "credit_score": 95,
    "prefer_tags": ["美食", "古镇", "户外"],
    "visited_city_count": 12,
    "follower_count": 128,
    "following_count": 56,
    "like_count": 340,
    "language": "zh-CN",
    "created_at": "2026-01-15T10:30:00Z"
  }
}
```

---

### 2.8 获取其他用户信息

```
GET /api/v1/user/{user_id}/profile
```

**描述**：查看指定用户的基本信息（公开信息）。游客也可访问，未登录时 `is_following` 固定返回 `false`。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | int | 是 | 用户 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 10002,
    "nickname": "旅行达人小王",
    "avatar": "https://cdn.example.com/avatars/10002.jpg",
    "credit_score": 88,
    "visited_city_count": 23,
    "follower_count": 520,
    "following_count": 120,
    "like_count": 1200,
    "is_following": false,
    "created_at": "2025-06-01T08:00:00Z"
  }
}
```

---

### 2.9 更新用户资料 @auth

```
PUT /api/v1/user/profile
```

**描述**：更新当前用户的昵称、头像、偏好标签、语言等。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| nickname | string | 否 | 昵称，2-20 个字符 |
| avatar | string | 否 | 头像 URL（需先通过上传图片接口获取） |
| prefer_tags | array[string] | 否 | 偏好标签，最多 5 个 |
| language | string | 否 | 语言偏好，如 zh-CN / en-US |

**请求示例**：

```json
{
  "nickname": "小明爱旅行",
  "prefer_tags": ["美食", "古镇", "户外", "拍照"],
  "language": "zh-CN"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 10001,
    "nickname": "小明爱旅行",
    "avatar": "https://cdn.example.com/avatars/10001.jpg",
    "prefer_tags": ["美食", "古镇", "户外", "拍照"],
    "language": "zh-CN"
  }
}
```

---

### 2.10 绑定手机号 @auth

```
POST /api/v1/user/bind-phone
```

**描述**：为通过第三方登录的用户绑定手机号。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号 |
| code | string | 是 | 短信验证码 |

**请求示例**：

```json
{
  "phone": "13800138000",
  "code": "123456"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 2.11 绑定邮箱 @auth

```
POST /api/v1/user/bind-email
```

**描述**：绑定邮箱并设置密码，用于后续 Email 登录。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码，8-32 位，需包含字母和数字 |

**请求示例**：

```json
{
  "email": "user@example.com",
  "password": "mypassword123"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 2.12 账号注销 @auth

```
POST /api/v1/user/delete-account
```

**描述**：永久注销当前账号。注销需先通过短信验证码确认身份（二次确认），注销后数据保留 30 天，期限内可恢复。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| verification_code | string | 是 | 短信验证码（通过 `2.1` 发送，scene=`delete_account`，6 位），用于二次确认注销操作 |
| reason | string | 否 | 注销原因（选填） |

**请求示例**：

```json
{
  "verification_code": "123456",
  "reason": "不再使用"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "delete_at": "2026-07-22T12:00:00Z",
    "restore_deadline": "2026-07-22T12:00:00Z"
  }
}
```

**说明**：
- 注销前需先调用 `2.1 发送验证码`（scene=`delete_account`）获取验证码，再携带 `verification_code` 调用本接口完成二次确认。
- 验证码错误返回 `2006`，验证码已过期返回 `2007`。
- `delete_at` 为账号实际删除时间（30 天后）。
- `restore_deadline` 为恢复账号的截止时间，在此时间前登录可自动恢复账号。

---

### 2.13 修改密码 @auth

```
PUT /api/v1/user/password
```

**描述**：修改当前登录用户的密码。需验证旧密码。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| old_password | string | 是 | 旧密码 |
| new_password | string | 是 | 新密码，8-32 位，需包含字母和数字 |

**请求示例**：

```json
{
  "old_password": "OldPass123",
  "new_password": "NewPass456"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 旧密码错误返回错误码 `2011`。
- 修改成功后，当前 Access Token 仍然有效，其他设备的 Token 会被失效。

---

### 2.14 用户设置 @auth

#### 2.14.1 获取用户设置

```
GET /api/v1/user/settings
```

**描述**：获取当前用户的设置项，包括推送、隐私、通用等配置。

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "notification_push": 1,
    "notification_interactive": 1,
    "notification_trip": 1,
    "footprint_public": "public",
    "content_public": 1,
    "allow_dm_scope": "all",
    "follow_scope": "all",
    "location_public": 0,
    "history_public": 0,
    "show_favorites": 0,
    "show_following": 0,
    "profile_visibility": "public",
    "personalized_recomm": 1,
    "language": "zh-CN",
    "font_size": 1,
    "dark_mode": 0,
    "daily_push_limit": 10
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| notification_push | integer | 推送开关：0-关闭, 1-开启（总开关，关闭后不发送任何推送） |
| notification_interactive | integer | 互动通知：0-关闭, 1-开启（点赞/评论/关注/收藏提醒） |
| notification_trip | integer | 行程提醒：0-关闭, 1-开启（行程开始/打卡提醒等） |
| footprint_public | string | 足迹可见性：public / friends / private |
| content_public | integer | 发布内容是否公开：0-否, 1-是 |
| allow_dm_scope | string | 私信权限：all / followers / nobody |
| follow_scope | string | 谁可以关注我：all / verified / nobody |
| location_public | integer | 位置信息是否公开：0-否, 1-是 |
| history_public | integer | 历史行程是否公开：0-不公开, 1-公开 |
| show_favorites | integer | 收藏列表是否对外可见：0-不公开, 1-公开 |
| show_following | integer | 关注/粉丝列表是否对外可见：0-不公开, 1-公开 |
| profile_visibility | string | 个人主页可见性：public / friends / private |
| personalized_recomm | integer | 个性化推荐开关：0-关闭, 1-开启 |
| language | string | 语言：zh-CN / en-US 等 |
| font_size | integer | 字体大小：0-小, 1-标准, 2-大 |
| dark_mode | integer | 深色模式：0-跟随系统, 1-关闭, 2-开启 |
| daily_push_limit | integer | 每日推送上限：单用户每日最多接收推送条数 |

#### 2.14.2 更新用户设置

```
PUT /api/v1/user/settings
```

**描述**：更新当前用户的设置项，所有字段均为可选，仅更新传入的字段。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| notification_push | integer | 否 | 推送开关：0-关闭, 1-开启 |
| notification_interactive | integer | 否 | 互动通知：0-关闭, 1-开启 |
| notification_trip | integer | 否 | 行程提醒：0-关闭, 1-开启 |
| footprint_public | string | 否 | 足迹可见性：public / friends / private |
| content_public | integer | 否 | 发布内容是否公开：0-否, 1-是 |
| allow_dm_scope | string | 否 | 私信权限：all / followers / nobody |
| follow_scope | string | 否 | 谁可以关注我：all / verified / nobody |
| location_public | integer | 否 | 位置信息是否公开：0-否, 1-是 |
| history_public | integer | 否 | 历史行程是否公开：0-不公开, 1-公开 |
| show_favorites | integer | 否 | 收藏列表是否对外可见：0-不公开, 1-公开 |
| show_following | integer | 否 | 关注/粉丝列表是否对外可见：0-不公开, 1-公开 |
| profile_visibility | string | 否 | 个人主页可见性：public / friends / private |
| personalized_recomm | integer | 否 | 个性化推荐开关：0-关闭, 1-开启 |
| language | string | 否 | 语言 |
| font_size | integer | 否 | 字体大小：0-小, 1-标准, 2-大 |
| dark_mode | integer | 否 | 深色模式：0-跟随系统, 1-关闭, 2-开启 |
| daily_push_limit | integer | 否 | 每日推送上限 |

**请求示例**：

```json
{
  "notification_push": 0,
  "font_size": 2
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- `trip_reminder` 为 P0 优先级提醒，服务端会忽略将其设为 `false` 的请求。
- 未登录返回错误码 `2001`。

---

## 三、内容模块

### 3.1 首页推荐流

```
GET /api/v1/content/feed/recommend
```

**描述**：首页推荐 Tab 内容流，双列瀑布流混合展示路线、打卡地、攻略。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "type": "route",
        "id": 2001,
        "title": "荆州古城一日游",
        "cover_image": "https://cdn.example.com/routes/2001.jpg",
        "total_duration": 480,
        "total_budget": 350.00,
        "rating": 4.5,
        "creator": {
          "id": 10001,
          "nickname": "小明爱旅行",
          "avatar": "https://cdn.example.com/avatars/10001.jpg"
        },
        "source": "user",
        "card_width": 320,
        "card_height": 240
      },
      {
        "type": "place",
        "id": 3001,
        "name": "荆州古城墙",
        "cover_image": "https://cdn.example.com/places/3001.jpg",
        "category": "sightseeing",
        "recommend_level": "must_go",
        "avg_cost": 0,
        "city": "荆州",
        "distance": 2.5,
        "card_width": 320,
        "card_height": 280
      },
      {
        "type": "guide",
        "id": 4001,
        "title": "荆州古城深度游攻略",
        "cover_image": "https://cdn.example.com/guides/4001.jpg",
        "category": "history",
        "city": "荆州",
        "duration": 360,
        "verified_count": 128,
        "card_width": 380,
        "card_height": 200
      }
    ],
    "total": 350,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

**说明**：
- 卡片类型 `type` 取值：`route`（路线）、`place`（打卡地）、`guide`（攻略）。
- 推荐算法基于用户偏好标签、定位模式、内容质量分和时效性综合排序。
- `card_width` / `card_height` 为客户端瀑布流排版参考尺寸（单位：逻辑像素）。

---

### 3.2 关注流 @auth

```
GET /api/v1/content/feed/following
```

**描述**：首页关注 Tab，仅展示已关注用户发布的内容。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "type": "place",
        "id": 3002,
        "name": "沙市洋码头",
        "cover_image": "https://cdn.example.com/places/3002.jpg",
        "category": "sightseeing",
        "recommend_level": "must_go",
        "city": "荆州",
        "creator": {
          "id": 10002,
          "nickname": "旅行达人小王",
          "avatar": "https://cdn.example.com/avatars/10002.jpg"
        },
        "from_following": true,
        "created_at": "2026-06-20T14:30:00Z"
      }
    ],
    "total": 45,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 3.3 附近流

```
GET /api/v1/content/feed/nearby
```

**描述**：首页附近 Tab，按距离排序展示打卡地和路线。支持三种定位模式：travel（GPS 实时定位，按真实距离排序）、plan / manual（按城市中心计算距离，不要求经纬度）。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| location_mode | string | 否 | travel | 定位模式：travel（GPS 实时定位）/ plan（规划行程城市）/ manual（手动选择城市） |
| city | string | 条件必填 | - | 城市名称，如"荆州"。plan / manual 模式必填，travel 模式不需要 |
| longitude | float | 条件必填 | - | 用户当前经度。travel 模式必填 |
| latitude | float | 条件必填 | - | 用户当前纬度。travel 模式必填 |
| radius | int | 否 | 10000 | 搜索半径，单位：米，默认 10000（10km） |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**请求示例**：

```
GET /api/v1/content/feed/nearby?location_mode=travel&longitude=112.2400&latitude=30.3300&radius=5000&page=1&page_size=20
GET /api/v1/content/feed/nearby?location_mode=manual&city=荆州&page=1&page_size=20
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "type": "place",
        "id": 3001,
        "name": "荆州古城墙",
        "cover_image": "https://cdn.example.com/places/3001.jpg",
        "category": "sightseeing",
        "recommend_level": "must_go",
        "distance": 0.8,
        "card_width": 320,
        "card_height": 280
      }
    ],
    "total": 120,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

**说明**：
- `distance` 单位为公里（km），保留一位小数。
- travel 模式以用户 GPS 经纬度为原点计算真实距离；plan / manual 模式以 `city` 城市中心点计算距离，不要求经纬度。
- travel 模式未提供经纬度（或定位权限被拒绝）时，返回错误码 `16001` 并提示开启定位权限；plan / manual 模式未提供 city 时返回错误码 `1002`。

---

### 3.4 城市内容流

```
GET /api/v1/content/feed/city
```

**描述**：首页城市 Tab，展示指定城市的内容。先展示四个分类卡片入口，再展示瀑布流。

> 该接口保留供未来城市 Tab 扩展使用，当前首页无城市 Tab（首页内容 Tab 为：推荐 / 关注 / 附近 / 攻略 / 打卡地 / 路线）。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| city | string | 是 | - | 城市名称，如 "荆州" |
| category | string | 否 | - | 筛选分类：route（路线）/ place（打卡地）/ guide（攻略），不传则混合 |
| longitude | float | 否 | - | 用户当前经度（用于计算距离） |
| latitude | float | 否 | - | 用户当前纬度（用于计算距离） |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "city": "荆州",
    "city_entries": [
      {
        "type": "ranking",
        "title": "榜单",
        "icon": "trophy",
        "link": "/city/荆州/ranking"
      },
      {
        "type": "hot_routes",
        "title": "热门路线",
        "icon": "route",
        "link": "/city/荆州/routes"
      },
      {
        "type": "hot_places",
        "title": "热门打卡地",
        "icon": "pin",
        "link": "/city/荆州/places"
      },
      {
        "type": "inbound",
        "title": "入境游",
        "icon": "globe",
        "link": "/city/荆州/inbound"
      }
    ],
    "list": [],
    "total": 200,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 3.5 攻略内容流

```
GET /api/v1/content/feed/guide
```

**描述**：首页攻略 Tab，支持子分类筛选。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| category | string | 否 | - | 攻略分类：nature（自然生态）/ history（历史文化）/ entertainment（人工娱乐）/ urban（城市公共）/ transport（交通枢纽）/ red_tourism（红色旅游）/ religion（宗教场所），不传则全部 |
| city | string | 否 | - | 筛选城市 |
| sort | string | 否 | newest | 排序：newest（最新）/ popular（最热） |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 4001,
        "title": "荆州古城深度游攻略",
        "cover_image": "https://cdn.example.com/guides/4001.jpg",
        "category": "history",
        "category_name": "历史文化",
        "city": "荆州",
        "duration": 360,
        "price": 0,
        "read_duration": 8,
        "verified_count": 128,
        "source": "AI",
        "created_at": "2026-06-15T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 3.6 打卡地内容流

```
GET /api/v1/content/feed/place
```

**描述**：首页打卡地 Tab，支持子分类筛选。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| category | string | 否 | - | 分类：food（吃喝）/ fun（玩乐）/ sightseeing（逛看）/ outdoor（户外）/ shopping（逛街）/ accommodation（住宿），不传则全部 |
| city | string | 否 | - | 筛选城市 |
| sort | string | 否 | newest | 排序：newest（最新）/ popular（最热）/ rating（评分最高） |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 3001,
        "name": "荆州古城墙",
        "cover_image": "https://cdn.example.com/places/3001.jpg",
        "category": "sightseeing",
        "category_name": "逛看",
        "recommend_level": "must_go",
        "avg_cost": 0,
        "rating": 4.7,
        "city": "荆州",
        "tags": ["历史", "必去", "拍照"],
        "checkin_count": 520,
        "created_at": "2026-03-01T08:00:00Z"
      }
    ],
    "total": 80,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 3.7 路线内容流

```
GET /api/v1/content/feed/route
```

**描述**：首页路线 Tab，沉浸式整屏滑动。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| category | string | 否 | - | 分类：half_day（半日游）/ one_day（一日游）/ two_day（两日游）/ three_day_plus（三日及以上），不传则全部 |
| city | string | 否 | - | 筛选城市 |
| sort | string | 否 | hot | 排序：hot（最热）/ newest（最新）/ rating（评分最高） |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 10 | 每页数量（路线 Tab 建议每次加载较少） |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 2001,
        "title": "荆州古城一日游",
        "cover_image": "https://cdn.example.com/routes/2001.jpg",
        "cover_images": [
          "https://cdn.example.com/routes/2001_1.jpg",
          "https://cdn.example.com/routes/2001_2.jpg",
          "https://cdn.example.com/routes/2001_3.jpg"
        ],
        "total_duration": 480,
        "total_budget": 350.00,
        "point_count": 8,
        "category": "one_day",
        "category_name": "一日游",
        "suitable_for": "情侣/朋友",
        "rating": 4.5,
        "usage_count": 320,
        "favorite_count": 150,
        "source": "user",
        "creator": {
          "id": 10001,
          "nickname": "小明爱旅行",
          "avatar": "https://cdn.example.com/avatars/10001.jpg"
        },
        "tags": ["文艺", "美食", "历史文化"],
        "created_at": "2026-05-20T09:00:00Z"
      }
    ],
    "total": 200,
    "page": 1,
    "page_size": 10,
    "has_more": true
  }
}
```

---

### 3.8 路线详情

```
GET /api/v1/content/route/{route_id}
```

**描述**：获取路线完整详情，包含时间轴点位列表。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| route_id | int | 是 | 路线 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 2001,
    "title": "荆州古城一日游",
    "cover_images": [
      "https://cdn.example.com/routes/2001_1.jpg",
      "https://cdn.example.com/routes/2001_2.jpg",
      "https://cdn.example.com/routes/2001_3.jpg"
    ],
    "total_duration": 480,
    "total_budget": 350.00,
    "suitable_for": "情侣/朋友",
    "category": "one_day",
    "category_name": "一日游",
    "tags": ["文艺", "美食", "历史文化"],
    "rating": 4.5,
    "usage_count": 320,
    "favorite_count": 150,
    "source": "user",
    "creator": {
      "id": 10001,
      "nickname": "小明爱旅行",
      "avatar": "https://cdn.example.com/avatars/10001.jpg",
      "is_following": false
    },
    "points": [
      {
        "sort_order": 1,
        "point_time": "09:00",
        "place_id": 3001,
        "name": "荆州古城墙",
        "stay_duration": 90,
        "cost": 0,
        "category": "sightseeing",
        "category_name": "逛看",
        "transport": "步行"
      },
      {
        "sort_order": 2,
        "point_time": "10:30",
        "place_id": 3002,
        "name": "荆州博物馆",
        "stay_duration": 120,
        "cost": 0,
        "category": "sightseeing",
        "category_name": "逛看",
        "transport": "步行"
      },
      {
        "sort_order": 3,
        "point_time": "12:30",
        "place_id": 3003,
        "name": "老街美食",
        "stay_duration": 60,
        "cost": 50.00,
        "category": "food",
        "category_name": "吃喝",
        "transport": "步行"
      }
    ],
    "is_favorited": false,
    "comment_count": 12,
    "created_at": "2026-05-20T09:00:00Z",
    "updated_at": "2026-06-18T15:00:00Z"
  }
}
```

---

### 3.9 打卡地详情

```
GET /api/v1/content/place/{place_id}
```

**描述**：获取打卡地完整详情，包含用户评价、避雷区、去过的人图片等。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| place_id | int | 是 | 打卡地（地点）ID |

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| longitude | float | 否 | - | 用户当前经度（计算距离） |
| latitude | float | 否 | - | 用户当前纬度（计算距离） |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 3001,
    "name": "荆州古城墙",
    "cover_images": [
      "https://cdn.example.com/places/3001_1.jpg",
      "https://cdn.example.com/places/3001_2.jpg",
      "https://cdn.example.com/places/3001_3.jpg"
    ],
    "category": "sightseeing",
    "category_name": "逛看",
    "recommend_level": "must_go",
    "recommend_level_name": "必去",
    "avg_cost": 0,
    "rating": 4.7,
    "checkin_count": 520,
    "address": "湖北省荆州市荆州区荆南路",
    "longitude": 112.2400,
    "latitude": 30.3300,
    "distance": 2.5,
    "opening_hours": "08:00-18:00",
    "is_open": true,
    "phone": "0716-XXXXXXX",
    "suggest_duration": 90,
    "tags": ["历史", "必去", "拍照"],
    "description": "荆州古城墙是中国现存最完好的古城墙之一...",
    "pitfall_tips": [
      {
        "index": 0,
        "content": "下午人很多，建议早上8点开门就去",
        "vote_count": 45,
        "is_voted": false
      },
      {
        "index": 1,
        "content": "城墙上没有遮阳，夏天注意防晒",
        "vote_count": 32,
        "is_voted": true
      }
    ],
    "recent_images": [
      "https://cdn.example.com/checkins/img_001.jpg",
      "https://cdn.example.com/checkins/img_002.jpg",
      "https://cdn.example.com/checkins/img_003.jpg"
    ],
    "creator": {
      "id": 10001,
      "nickname": "小明爱旅行",
      "avatar": "https://cdn.example.com/avatars/10001.jpg"
    },
    "is_favorited": false,
    "comment_count": 28,
    "created_at": "2026-03-01T08:00:00Z",
    "cps_links": {
      "ticket": { "platform": "ctrip", "url": "https://m.ctrip.com/...", "price_from": 35 },
      "hotel": { "platform": "ctrip", "url": "https://m.ctrip.com/...", "price_from": 200 },
      "food": { "platform": "eleme", "url": "https://h5.ele.me/..." },
      "taxi": { "platform": "amap", "url": "https://uri.amap.com/..." }
    }
  }
}
```

---

### 3.10 攻略详情

```
GET /api/v1/content/guide/{guide_id}
```

**描述**：获取攻略完整详情，包含分类专属内容模块。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| guide_id | int | 是 | 攻略 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 4001,
    "title": "荆州古城深度游攻略",
    "cover_images": [
      "https://cdn.example.com/guides/4001_1.jpg",
      "https://cdn.example.com/guides/4001_2.jpg"
    ],
    "category": "history",
    "category_name": "历史文化",
    "city": "荆州",
    "duration": 360,
    "price": 0,
    "read_duration": 8,
    "address": "湖北省荆州市荆州区",
    "opening_hours": "08:00-18:00",
    "transport": "公交：乘坐1路、12路到古城站下车；自驾：导航至荆州古城停车场",
    "content": "## 荆州古城简介\n\n荆州古城，又名江陵城...",
    "highlights": "必看：古城墙、宾阳楼、张居正故居",
    "pitfall_reminders": "注意：周末人流量大，建议工作日前往",
    "verified_count": 128,
    "source": "AI",
    "category_content": {
      "history_background": "荆州古城始建于春秋战国时期，是中国南方保存最完好的古城墙之一。作为三国文化的重要发源地，荆州见证了无数历史风云。",
      "must_see_exhibits": [
        {
          "name": "古城墙",
          "description": "全长约11公里，现存城墙为明清时期重建"
        },
        {
          "name": "宾阳楼",
          "description": "古城东门城楼，可俯瞰全城风貌"
        }
      ],
      "guide_service": "提供中英文讲解服务，费用50元/次",
      "tour_route": "建议从宾阳楼出发 → 沿城墙步行至南门 → 参观张居正故居 → 返回"
    },
    "recent_images": [
      "https://cdn.example.com/checkins/guide_img_001.jpg",
      "https://cdn.example.com/checkins/guide_img_002.jpg"
    ],
    "is_favorited": false,
    "comment_count": 15,
    "created_at": "2026-06-15T10:00:00Z"
  }
}
```

**说明**：
- `category_content` 根据攻略分类 `category` 返回不同的结构：
  - `nature`（自然生态）：包含 `difficulty`、`best_season`、`must_see_spots`、`equipment_tips`、`internal_transport`
  - `history`（历史文化）：包含 `history_background`、`must_see_exhibits`、`guide_service`、`tour_route`
  - `entertainment`（人工娱乐）：包含 `must_play_top3`、`queue_estimate`、`show_schedule`、`fast_pass_guide`
  - `urban`（城市公共）：包含 `open_areas`、`surrounding_facilities`、`accessibility`、`event_calendar`
  - `transport`（交通枢纽）：包含 `arrival_methods`、`transfer_guide`、`luggage_storage`、`sim_card_locations`
  - `red_tourism`（红色旅游）：包含 `revolutionary_background`、`memorial_sites`、`visiting_route`、`education_significance`
  - `religion`（宗教场所）：包含 `religious_background`、`worship_etiquette`、`etiquette_tips`、`photography_rules`

---

### 3.11 搜索

```
GET /api/v1/content/search
```

**描述**：全局搜索，支持搜索路线、打卡地、攻略。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| keyword | string | 是 | - | 搜索关键词，最少 1 个字符 |
| type | string | 否 | all | 搜索类型：all（全部）/ route（路线）/ place（打卡地）/ guide（攻略） |
| city | string | 否 | - | 限定城市 |
| sort | string | 否 | relevance | 排序：relevance（相关性）/ newest（最新）/ popular（最热） |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "keyword": "荆州",
    "total_route": 15,
    "total_place": 80,
    "total_guide": 10,
    "list": [
      {
        "type": "place",
        "id": 3001,
        "name": "荆州古城墙",
        "cover_image": "https://cdn.example.com/places/3001.jpg",
        "category": "sightseeing",
        "category_name": "逛看",
        "city": "荆州",
        "rating": 4.7,
        "highlight": "荆州<b>古城</b>墙是中国现存最完好的...",
        "created_at": "2026-03-01T08:00:00Z"
      }
    ],
    "total": 105,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

**说明**：
- `highlight` 字段包含 HTML 标签 `<b>` 标记匹配的关键词，用于搜索结果高亮。
- 搜索支持模糊匹配，同时搜索中文和拼音。

---

### 3.12 搜索建议

```
GET /api/v1/content/search/suggest
```

**描述**：搜索联想词，用户输入时实时展示。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 输入关键词，最少 1 个字符 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "suggestions": [
      {
        "text": "荆州古城墙",
        "type": "place"
      },
      {
        "text": "荆州一日游",
        "type": "route"
      },
      {
        "text": "荆州博物馆",
        "type": "place"
      },
      {
        "text": "荆州古城深度游攻略",
        "type": "guide"
      }
    ]
  }
}
```

---

### 3.13 创建路线 @auth

```
POST /api/v1/content/routes
```

**描述**：用户创建自定义路线，保存为个人路线或发布到社区。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 路线标题，2-64 字符 |
| cover_image | array[string] | 否 | 封面图 URL 数组，最多 3 张 |
| category | string | 是 | 分类：half_day / one_day / two_day / three_day_plus |
| total_duration | int | 是 | 总时长，单位：分钟 |
| total_budget | float | 否 | 总预算，单位：元 |
| suitable_for | string | 否 | 适合人群 |
| tags | array[string] | 否 | 标签，最多 5 个 |
| points | array[object] | 是 | 路线点位列表 |
| points[].sort_order | int | 是 | 排序序号，从 1 开始 |
| points[].place_id | int | 否 | 地点 ID（关联已有地点时传） |
| points[].name | string | 是 | 点位名称 |
| points[].stay_duration | int | 否 | 停留时长，单位：分钟 |
| points[].cost | float | 否 | 预计消费，单位：元 |
| points[].transport | string | 否 | 交通方式 |
| visibility | string | 否 | 可见性：public / private，默认 public |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 2005,
    "title": "荆州古城半日游",
    "status": 1
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 至少包含 2 个点位。
- `place_id` 为空时表示自定义点位，不关联已有地点。
- 信用分 < 60 返回错误码 `3007`。

---

### 3.14 打卡记录详情

```
GET /api/v1/content/checkin/{checkin_id}
```

**描述**：获取指定打卡记录的详情信息。游客可查看公开打卡记录。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| checkin_id | int | 是 | 打卡记录 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 6001,
    "user_id": 10001,
    "user_nickname": "旅行者小明",
    "user_avatar": "https://cdn.example.com/avatars/10001.jpg",
    "place_id": 3004,
    "place_name": "户部巷",
    "images": [
      "https://cdn.example.com/checkin/6001_1.jpg",
      "https://cdn.example.com/checkin/6001_2.jpg"
    ],
    "rating": 4.5,
    "description": "非常热闹的小吃街，推荐热干面和豆皮",
    "actual_cost": 35.00,
    "pitfall_tips": ["周末人很多，建议早上来", "现金和手机支付都可以"],
    "checkin_time": "2026-06-20T08:30:00Z",
    "checkin_type": "store",
    "source": "trip_checkin",
    "content_richness": "standard",
    "gps_verified": true,
    "extra_data": {
      "dish_review": "热干面和豆皮都很正宗"
    },
    "visibility": "public",
    "status": "approved",
    "like_count": 12,
    "comment_count": 3,
    "created_at": "2026-06-20T08:30:00Z"
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 打卡记录不存在返回错误码 `4004`。
- `visibility` 为 `private` 时仅本人可查看，其他用户访问返回错误码 `2004`。

---

### 3.15 获取 CPS 分销链接 @auth

```
GET /api/v1/places/{place_id}/cps
```

**描述**：获取指定打卡地的 CPS 分销链接，支持按类型筛选（门票 / 酒店 / 美食 / 打车）。打卡地详情接口（3.9）返回的 `cps_links` 为聚合预览，本接口可按需获取单类链接。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| place_id | int | 是 | 打卡地（地点）ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 否 | 筛选类型：ticket（门票）/ hotel（酒店）/ food（美食）/ taxi（打车）。不传则返回全部可用类型 |

**请求示例**：

```
GET /api/v1/places/3001/cps?type=ticket
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "platform": "ctrip",
    "url": "https://m.ctrip.com/...",
    "price_from": 35
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- `price_from` 为起价，部分类型（如美食、打车）可能不返回该字段。
- CPS 链接生成失败返回错误码 `15001`；链接已失效返回错误码 `15003`；无匹配结果返回错误码 `15004`。
- 地点不存在返回错误码 `4001`。

---

## 四、行程模块

### 4.1 创建行程（手动创建） @auth

```
POST /api/v1/trips/create
```

**描述**：手动创建新的行程，用户自行添加点位。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 行程名称，1-50 个字符 |
| trip_date | date | 是 | 出行日期，格式：YYYY-MM-DD |
| end_date | date | 否 | 行程结束日期，格式：YYYY-MM-DD；不传默认等于 trip_date（单日行程）；支持多日行程（与 DB `trip.end_date` 对应，详见 05_数据库设计文档.md） |
| route_id | int | 否 | 关联的路线 ID（如果基于路线创建） |
| points | array | 否 | 点位列表 |

**point 对象结构**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| place_id | int | 否 | 关联的地点 ID |
| name | string | 是 | 点位名称 |
| point_time | string | 否 | 安排时间，格式：HH:mm |
| stay_duration | int | 否 | 停留时长，单位：分钟 |
| cost | float | 否 | 预计消费，单位：元 |
| category | string | 否 | 分类 |
| transport | string | 否 | 交通方式 |

**请求示例**：

```json
{
  "name": "荆州周末游",
  "trip_date": "2026-07-15",
  "end_date": "2026-07-15",
  "route_id": 2001,
  "points": [
    {
      "place_id": 3001,
      "name": "荆州古城墙",
      "point_time": "09:00",
      "stay_duration": 90,
      "cost": 0,
      "category": "sightseeing",
      "transport": "步行"
    },
    {
      "place_id": 3002,
      "name": "荆州博物馆",
      "point_time": "10:30",
      "stay_duration": 120,
      "cost": 0,
      "category": "sightseeing",
      "transport": "步行"
    }
  ]
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5001,
    "name": "荆州周末游",
    "trip_date": "2026-07-15",
    "end_date": "2026-07-15",
    "status": "pending",
    "route_id": 2001,
    "point_count": 2,
    "created_at": "2026-06-22T10:00:00Z"
  }
}
```

**说明**：
- 每个用户每天最多创建 5 个行程，超出返回错误码 `1006`。
- 如果传入了 `route_id`，会从路线中复制点位数据到行程中。

---

### 4.2 获取行程列表 @auth

```
GET /api/v1/trips/list
```

**描述**：获取当前用户的行程列表，按状态筛选。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| status | string | 否 | - | 行程状态：pending（待出发）/ active（进行中）/ completed（已完成），不传则全部 |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 5001,
        "name": "荆州周末游",
        "trip_date": "2026-07-15",
        "status": "pending",
        "status_name": "待出发",
        "route_id": 2001,
        "route_title": "荆州古城一日游",
        "point_count": 8,
        "checked_count": 0,
        "skipped_count": 0,
        "countdown_days": 23,
        "created_at": "2026-06-22T10:00:00Z"
      },
      {
        "id": 5002,
        "name": "武汉一日游",
        "trip_date": "2026-06-20",
        "status": "active",
        "status_name": "进行中",
        "point_count": 6,
        "checked_count": 3,
        "skipped_count": 1,
        "progress": 60.0,
        "created_at": "2026-06-18T09:00:00Z"
      }
    ],
    "total": 12,
    "page": 1,
    "page_size": 20,
    "has_more": false
  }
}
```

**说明**：
- `progress` 为行程进度百分比，计算公式：已打卡数 / (总点数 - 已跳过数) x 100%。
- `countdown_days` 为距离出行日期剩余天数，仅 `pending` 状态显示。

---

### 4.3 行程详情

```
GET /api/v1/trips/{trip_id}
```

**描述**：获取行程完整详情，包含点位列表和打卡状态。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5002,
    "name": "武汉一日游",
    "trip_date": "2026-06-20",
    "status": "active",
    "status_name": "进行中",
    "route_id": 2002,
    "route_title": "武汉经典一日游",
    "point_count": 6,
    "checked_count": 3,
    "skipped_count": 1,
    "progress": 60.0,
    "elapsed_minutes": 180,
    "estimated_minutes": 420,
    "total_cost": 150.00,
    "current_point": {
      "sort_order": 4,
      "place_id": 3005,
      "name": "黄鹤楼",
      "point_time": "14:00",
      "stay_duration": 90,
      "cost": 70.00,
      "category": "sightseeing",
      "category_name": "逛看",
      "transport": "地铁",
      "checkin_status": "pending",
      "checkin_status_name": "待打卡",
      "place": {
        "id": 3005,
        "name": "黄鹤楼",
        "longitude": 114.3025,
        "latitude": 30.5446,
        "address": "武汉市武昌区蛇山西山坡特1号",
        "opening_hours": "08:00-18:00",
        "is_open": true
      }
    },
    "points": [
      {
        "sort_order": 1,
        "place_id": 3004,
        "name": "户部巷",
        "point_time": "08:00",
        "stay_duration": 60,
        "cost": 30.00,
        "category": "food",
        "category_name": "吃喝",
        "transport": "步行",
        "checkin_status": "checked",
        "checkin_status_name": "已打卡",
        "checkin_record_id": 6001
      },
      {
        "sort_order": 2,
        "place_id": 3006,
        "name": "武汉长江大桥",
        "point_time": "09:30",
        "stay_duration": 30,
        "cost": 0,
        "category": "sightseeing",
        "category_name": "逛看",
        "transport": "步行",
        "checkin_status": "checked",
        "checkin_status_name": "已打卡",
        "checkin_record_id": 6002
      },
      {
        "sort_order": 3,
        "place_id": 3007,
        "name": "昙华林",
        "point_time": "10:30",
        "stay_duration": 60,
        "cost": 0,
        "category": "shopping",
        "category_name": "逛街",
        "transport": "步行",
        "checkin_status": "skipped",
        "checkin_status_name": "已跳过"
      },
      {
        "sort_order": 4,
        "place_id": 3005,
        "name": "黄鹤楼",
        "point_time": "14:00",
        "stay_duration": 90,
        "cost": 70.00,
        "category": "sightseeing",
        "category_name": "逛看",
        "transport": "地铁",
        "checkin_status": "pending",
        "checkin_status_name": "待打卡"
      },
      {
        "sort_order": 5,
        "place_id": 3008,
        "name": "楚河汉街",
        "point_time": "17:00",
        "stay_duration": 120,
        "cost": 50.00,
        "category": "shopping",
        "category_name": "逛街",
        "transport": "地铁",
        "checkin_status": "pending",
        "checkin_status_name": "待打卡"
      },
      {
        "sort_order": 6,
        "place_id": 3009,
        "name": "东湖",
        "point_time": "19:00",
        "stay_duration": 60,
        "cost": 0,
        "category": "outdoor",
        "category_name": "户外",
        "transport": "打车",
        "checkin_status": "pending",
        "checkin_status_name": "待打卡"
      }
    ],
    "created_at": "2026-06-18T09:00:00Z"
  }
}
```

**说明**：
- `current_point` 为第一个 `checkin_status` 为 `pending` 的点位，即当前需要打卡的点位。
- 已打卡的点位包含 `checkin_record_id`，可跳转到打卡记录详情。

---

### 4.4 开始行程 @auth

```
POST /api/v1/trips/{trip_id}/start
```

**描述**：将待出发的行程状态变更为进行中。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5001,
    "status": "active",
    "status_name": "进行中",
    "started_at": "2026-07-15T08:00:00Z"
  }
}
```

**说明**：
- 仅 `pending` 状态的行程可以开始。
- 如果当前日期早于 `trip_date`，仍允许开始（提前出发）。

---

### 4.5 打卡（到店打卡） @auth

```
POST /api/v1/trips/{trip_id}/checkin/store
```

**描述**：到店打卡，需要 GPS 验证在打卡点 100 米范围内。生成打卡记录和 UGC 评价。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| point_id | int | 是 | 行程点位 ID（trip_point.id） |
| longitude | float | 是 | 用户当前经度 |
| latitude | float | 是 | 用户当前纬度 |
| images | array[string] | 是 | 图片 URL 列表，至少 1 张 |
| rating | float | 否 | 评分，范围 0.0-5.0，步长 0.5 |
| description | string | 否 | 感受描述 |
| actual_cost | float | 否 | 实际消费，单位：元 |
| pitfall_tips | array[string] | 否 | 避雷建议，JSON 数组（与数据库一致），如 ["建议早上8点去", "夏天注意防晒"] |
| extra_data | object | 否 | 差异化专属内容（JSON 对象），按 place.category 存储不同字段，如吃喝的菜品评价、住宿的房型体验等 |
| visibility | string | 否 | 可见性：public（公开）/ private（仅自己可见），默认 public |

**请求示例**：

```json
{
  "point_id": 10001,
  "longitude": 112.2400,
  "latitude": 30.3300,
  "images": [
    "https://cdn.example.com/upload/img_001.jpg",
    "https://cdn.example.com/upload/img_002.jpg"
  ],
  "rating": 4.5,
  "description": "古城墙非常壮观，早上去人很少，拍照很好看！",
  "actual_cost": 0,
  "pitfall_tips": ["下午人很多，建议早上8点开门就去"],
  "extra_data": {
    "dish_review": "热干面很正宗"
  },
  "visibility": "public"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "checkin_record_id": 6001,
    "point_id": 10001,
    "checkin_status": "checked",
    "progress": 60.0,
    "credit_change": 5
  }
}
```

**说明**：
- 服务端校验用户 GPS 坐标与打卡点坐标距离，超过 100 米返回错误码 `16003`（GPS 距离超限）。
- `credit_change` 为本次打卡获得的信用分变化（发布内容 +5）。
- **限频**：20 次/天/用户（与 4.6 普通打卡合计计数），防虚假刷打卡，详见 1.7.1。

---

### 4.6 打卡（普通打卡） @auth

```
POST /api/v1/trips/{trip_id}/checkin/normal
```

**描述**：普通打卡，仅标记完成，不生成 UGC 内容。无需 GPS 验证。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| point_id | int | 是 | 行程点位 ID（trip_point.id） |

**请求示例**：

```json
{
  "point_id": 10001
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "point_id": 10001,
    "checkin_status": "checked",
    "progress": 60.0
  }
}
```

**说明**：
- **限频**：20 次/天/用户（与 4.5 到店打卡合计计数），防虚假刷打卡，详见 1.7.1。

---

### 4.7 离线打卡同步 @auth

```
POST /api/v1/trips/{trip_id}/checkin/sync
```

**描述**：同步离线期间的打卡记录。客户端在无网络时本地记录打卡信息，联网后批量提交。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| records | array | 是 | 离线打卡记录列表 |

**record 对象结构**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| point_id | int | 是 | 行程点位 ID |
| type | string | 是 | 打卡类型：store / normal |
| longitude | float | 否 | 经度（到店打卡时必填） |
| latitude | float | 否 | 纬度（到店打卡时必填） |
| checkin_time | string | 是 | 打卡时间，格式：ISO 8601 |
| images | array | 否 | 图片列表（到店打卡时必填） |
| rating | float | 否 | 评分 |
| description | string | 否 | 感受描述 |
| actual_cost | float | 否 | 实际消费 |
| pitfall_tips | array[string] | 否 | 避雷建议，JSON 数组（与数据库一致） |
| visibility | string | 否 | 可见性 |

**请求示例**：

```json
{
  "records": [
    {
      "point_id": 10001,
      "type": "normal",
      "checkin_time": "2026-07-15T09:00:00Z"
    },
    {
      "point_id": 10002,
      "type": "store",
      "longitude": 112.2400,
      "latitude": 30.3300,
      "checkin_time": "2026-07-15T10:30:00Z",
      "images": ["https://cdn.example.com/upload/img_003.jpg"],
      "rating": 4.0,
      "description": "博物馆很值得一看"
    }
  ]
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "synced": 2,
    "failed": 0,
    "details": [
      {
        "point_id": 10001,
        "status": "synced",
        "checkin_status": "checked"
      },
      {
        "point_id": 10002,
        "status": "synced",
        "checkin_record_id": 6003,
        "checkin_status": "checked"
      }
    ],
    "progress": 60.0
  }
}
```

**说明**：
- 超过 24 小时未同步的打卡记录，客户端应提示用户检查网络。
- 到店打卡（type=store）需校验 GPS 坐标，如距离超过 100 米，该条记录同步失败。

**离线打卡同步冲突处理策略**：

| 冲突场景 | 处理策略 | 返回状态 | 客户端行为 |
|---------|---------|---------|-----------|
| 同一点位在线已打卡，离线又打卡（多设备） | **在线记录优先**，离线记录标记为 `conflict`，不覆盖已有打卡数据 | `status: "conflict"` | 弹窗提示"该点位已打卡，离线记录未同步"，提供"查看已有打卡"按钮 |
| 离线打卡时点位已被跳过 | **打卡优先**，取消跳过状态，更新为已打卡 | `status: "synced", checkin_status: "checked"` | Toast 提示"打卡已同步，跳过状态已取消" |
| 离线打卡时点位已从行程删除 | **拒绝同步**，返回错误 | `status: "failed", reason: "point_deleted"` | 提示"该点位已从行程中移除，打卡记录保存为草稿"，可选"保存为独立打卡内容" |
| 离线打卡时行程已完成/删除 | **拒绝同步**，返回错误 | `status: "failed", reason: "trip_closed"` | 提示"行程已结束，打卡记录保存为草稿"，可选"转为主动发布内容" |
| 离线打卡时间与在线打卡时间间隔 < 5 分钟（重复打卡） | **去重处理**，仅保留时间更早的记录 | `status: "deduplicated"` | 静默处理，不提示用户 |
| 离线打卡 GPS 坐标偏差 > 100 米（到店打卡） | **同步但标记未验证**，`gps_verified` 设为 `false` | `status: "synced", gps_verified: false` | Toast 提示"打卡已同步（GPS 未验证）" |

> **设计原则**：数据不丢失——任何无法同步的离线打卡记录均保留在本地，提供"转为主动发布"的降级路径，确保用户离线打卡内容不会因冲突而消失。

---

### 4.8 跳过点位 @auth

```
POST /api/v1/trips/{trip_id}/point/{point_id}/skip
```

**描述**：将当前点位标记为跳过，跳过不计入行程进度。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |
| point_id | int | 是 | 行程点位 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "point_id": 10001,
    "checkin_status": "skipped",
    "progress": 75.0
  }
}
```

**说明**：
- 仅 `pending` 状态的点位可以跳过。
- 跳过点位后，系统自动重新计算剩余路线的最优顺序。

---

### 4.9 取消跳过 @auth

```
POST /api/v1/trips/{trip_id}/point/{point_id}/unskip
```

**描述**：将已跳过的点位恢复为待打卡状态。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |
| point_id | int | 是 | 行程点位 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "point_id": 10001,
    "checkin_status": "pending",
    "progress": 50.0
  }
}
```

**说明**：
- 仅 `skipped` 状态的点位可以取消跳过。

---

### 4.10 删除行程 @auth

```
DELETE /api/v1/trips/{trip_id}
```

**描述**：删除指定行程，级联删除所有行程点位。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

**说明**：
- 进行中的行程（active）不能删除，返回错误码 `5006`（不能删除进行中的行程）；仅待出发（pending）或已完成（completed）的行程可删除。
- 删除已完成行程时，已产生的打卡记录保留。
- 删除操作不可逆，客户端需二次确认。

---

### 4.11 AI 行程规划 @auth（P1）

```
POST /api/v1/trips/ai-plan
```

**描述**：AI 智能规划行程，根据用户输入自动生成行程建议。（P1 功能，P0 阶段不开发此接口）

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| city | string | 是 | 目标城市 |
| trip_date | date | 是 | 出行日期，格式：YYYY-MM-DD |
| duration | int | 是 | 行程天数，1-7 天 |
| budget | float | 否 | 预算，单位：元 |
| prefer_tags | array[string] | 否 | 偏好标签，如 ["美食", "历史", "户外"] |
| with_people | string | 否 | 同行人群：solo（独自）/ couple（情侣）/ family（亲子）/ friends（朋友） |
| start_time | string | 否 | 每天出发时间，格式：HH:mm，默认 "09:00" |
| additional_notes | string | 否 | 补充说明，如 "老人同行，需要轻松路线" |

**请求示例**：

```json
{
  "city": "荆州",
  "trip_date": "2026-07-15",
  "duration": 1,
  "budget": 500.00,
  "prefer_tags": ["美食", "历史文化"],
  "with_people": "friends",
  "start_time": "09:00",
  "additional_notes": "希望多安排一些美食打卡"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "name": "荆州美食文化一日游",
    "trip_date": "2026-07-15",
    "total_duration": 480,
    "total_budget": 350.00,
    "source": "AI",
    "ai_disclaimer": "此行程由 AI 生成，仅供参考。实地打卡后欢迎验证和纠错。",
    "points": [
      {
        "sort_order": 1,
        "point_time": "09:00",
        "place_id": 3001,
        "name": "荆州古城墙",
        "stay_duration": 90,
        "cost": 0,
        "category": "sightseeing",
        "category_name": "逛看",
        "transport": "步行",
        "reason": "荆州必打卡地标，早上去人少适合拍照"
      },
      {
        "sort_order": 2,
        "point_time": "10:30",
        "place_id": 3010,
        "name": "荆州博物馆",
        "stay_duration": 120,
        "cost": 0,
        "category": "sightseeing",
        "category_name": "逛看",
        "transport": "步行",
        "reason": "了解荆州历史文化的最佳场所"
      },
      {
        "sort_order": 3,
        "point_time": "12:30",
        "place_id": 3011,
        "name": "老街美食街",
        "stay_duration": 90,
        "cost": 80.00,
        "category": "food",
        "category_name": "吃喝",
        "transport": "打车",
        "reason": "荆州本地美食聚集地，推荐荆州鱼糕、早堂面"
      }
    ]
  }
}
```

**说明**：
- **限频**：10 次/天/用户，超出返回错误码 `1006`（与 AI 攻略生成共享每日 AI 配额），详见 1.7.1。
- 简单行程（1-2 天）预计响应时间 ≤ 3 秒，复杂行程（3-7 天）≤ 5 秒。
- AI 生成内容需在前端明确标注 `ai_disclaimer` 提示文案。
- 返回的数据可直接用于调用创建行程接口（4.1）。

---

### 4.12 添加点位到行程 @auth

```
POST /api/v1/trips/{trip_id}/points
```

**描述**：向已有行程中新增点位。AI 实时提示最佳插入位置。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| place_id | int | 是 | 地点 ID |
| insert_after | int | 否 | 插入到指定点位序号之后，不传则追加到末尾 |

**请求示例**：

```json
{
  "place_id": 3012,
  "insert_after": 2
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "point_id": 10010,
    "sort_order": 3,
    "ai_suggestion": "建议安排在 11:00-12:00，这样不会影响下午的行程"
  }
}
```

---

### 4.13 删除行程点位 @auth

```
DELETE /api/v1/trips/{trip_id}/point/{point_id}
```

**描述**：删除行程中的指定点位。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |
| point_id | int | 是 | 行程点位 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "progress": 50.0
  }
}
```

**说明**：
- 删除点位后，剩余点位自动重新排序。
- 已打卡的点位不可删除。

---

### 4.14 编辑行程基本信息 @auth

```
PUT /api/v1/trips/{trip_id}
```

**描述**：编辑行程的基本信息（标题、日期、可见性等）。仅行程创建者可操作。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 否 | 行程标题 |
| trip_date | string | 否 | 行程日期，格式 YYYY-MM-DD |
| visibility | string | 否 | 可见性：public / private |
| city | string | 否 | 目标城市 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5001,
    "title": "荆州周末游（已更新）",
    "updated_at": "2026-06-24T10:00:00Z"
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 行程不存在返回错误码 `5001`。
- 非创建者编辑返回错误码 `2004`。
- 行程已完成（status=completed）时返回错误码 `5009`。

---

### 4.15 结束行程 @auth

```
POST /api/v1/trips/{trip_id}/complete
```

**描述**：手动结束行程，将状态从 `active` 改为 `completed`。结束后不可再修改点位。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5001,
    "status": "completed",
    "completed_at": "2026-06-24T18:00:00Z",
    "summary": {
      "total_points": 6,
      "checked_points": 5,
      "skipped_points": 1,
      "total_cost": 128.50
    }
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 行程不存在返回错误码 `5001`。
- 非创建者操作返回错误码 `2004`。
- 行程状态不为 `active` 时返回错误码 `5003`。

---

### 4.16 行程清单管理 @auth

#### 4.16.1 获取行程清单

```
GET /api/v1/trips/{trip_id}/checklist
```

**描述**：获取指定行程的清单项列表，包含手动添加和 AI 自动生成的项。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 7001,
        "trip_id": 5001,
        "item_name": "身份证",
        "item_category": "document",
        "is_checked": true,
        "is_ai_generated": true,
        "sort_order": 1,
        "created_at": "2026-06-24T10:00:00Z"
      },
      {
        "id": 7002,
        "trip_id": 5001,
        "item_name": "充电宝",
        "item_category": "electronic",
        "is_checked": false,
        "is_ai_generated": false,
        "sort_order": 2,
        "created_at": "2026-06-24T10:05:00Z"
      }
    ]
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 清单项 ID |
| trip_id | int | 所属行程 ID |
| item_name | string | 清单项名称 |
| item_category | string | 分类：document（证件）/ electronic（电子设备）/ clothing（衣物）/ other（其他） |
| is_checked | boolean | 是否已勾选 |
| is_ai_generated | boolean | 是否为 AI 自动生成 |
| sort_order | int | 排序序号 |
| created_at | string | 创建时间 |

#### 4.16.2 添加清单项

```
POST /api/v1/trips/{trip_id}/checklist
```

**描述**：手动添加一个行程清单项。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| item_name | string | 是 | 清单项名称 |
| item_category | string | 否 | 分类：document（证件）/ electronic（电子设备）/ clothing（衣物）/ other（其他），默认 other |

**请求示例**：

```json
{
  "item_name": "防晒霜",
  "item_category": "other"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 7003,
    "trip_id": 5001,
    "item_name": "防晒霜",
    "item_category": "other",
    "is_checked": false,
    "is_ai_generated": false,
    "sort_order": 3
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

#### 4.16.3 勾选/取消勾选清单项

```
PUT /api/v1/trips/{trip_id}/checklist/{item_id}
```

**描述**：勾选或取消勾选指定清单项。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |
| item_id | int | 是 | 清单项 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| is_checked | bool | 是 | 是否勾选 |

**请求示例**：

```json
{
  "is_checked": true
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 7003,
    "is_checked": true,
    "checked_at": "2026-06-24T18:30:00Z"
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

#### 4.16.4 删除清单项

```
DELETE /api/v1/trips/{trip_id}/checklist/{item_id}
```

**描述**：删除指定的行程清单项。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |
| item_id | int | 是 | 清单项 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 行程不存在返回错误码 `5001`。
- 清单项不存在返回错误码 `1004`。

---

### 4.17 行程账单管理 @auth（P1）

> **P0/P1 说明**：P0 阶段仅展示行程总预算与已消费汇总（读取概览数据），账单明细管理（手动添加/删除消费记录、分类统计）为 P1 功能。

#### 4.17.1 获取行程账单

```
GET /api/v1/trips/{trip_id}/expenses
```

**描述**：获取指定行程的账单汇总，包含实际消费、预估消费、预算及按分类的明细。支持按分类筛选。（P0 仅返回总预算与已消费汇总，明细分类为 P1）

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category | string | 否 | 按分类筛选：food / ticket / transport / shopping / other。不传则返回全部分类 |

**请求示例**：

```
GET /api/v1/trips/5001/expenses?category=food
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_actual": 68.00,
    "total_estimated": 80.00,
    "budget": 200.00,
    "categories": [
      {
        "category": "food",
        "actual_amount": 68.00,
        "estimated_amount": 80.00,
        "items": [
          {
            "id": 8001,
            "category": "food",
            "amount": 35.00,
            "description": "户部巷早餐",
            "source": "auto",
            "created_at": "2026-06-24T08:30:00Z"
          }
        ]
      }
    ]
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

#### 4.17.2 手动添加消费记录

```
POST /api/v1/trips/{trip_id}/expenses
```

**描述**：手动添加一条消费记录到行程账单。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category | string | 是 | 分类：food / ticket / transport / shopping / other |
| amount | number | 是 | 金额，单位：元 |
| description | string | 否 | 消费描述 |

**请求示例**：

```json
{
  "category": "food",
  "amount": 35.00,
  "description": "户部巷早餐"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 8002,
    "trip_id": 5001,
    "category": "food",
    "amount": 35.00,
    "description": "户部巷早餐",
    "source": "manual"
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

#### 4.17.3 删除消费记录

```
DELETE /api/v1/trips/{trip_id}/expenses/{expense_id}
```

**描述**：删除指定的消费记录。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |
| expense_id | int | 是 | 消费记录 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 行程不存在返回错误码 `5001`。
- 消费记录不存在返回错误码 `1004`。
- 账单汇总失败返回错误码 `18002`。

---

### 4.18 加入行程（从路线模板创建） @auth

```
POST /api/v1/trips/from-route
```

**描述**：基于现有路线模板创建行程，将路线点位复制为行程点位。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| route_id | int | 是 | 路线 ID |
| name | string | 是 | 行程名称 |
| trip_date | string | 是 | 出行日期，格式：YYYY-MM-DD |
| end_date | string | 否 | 结束日期，格式：YYYY-MM-DD |
| budget | number | 否 | 预算，单位：元 |

**请求示例**：

```json
{
  "route_id": 2001,
  "name": "端午节荆州游",
  "trip_date": "2026-06-25",
  "end_date": "2026-06-26",
  "budget": 500.00
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5002,
    "user_id": 10001,
    "name": "端午节荆州游",
    "trip_date": "2026-06-25",
    "end_date": "2026-06-26",
    "budget": 500.00,
    "visibility": "private",
    "status": "pending",
    "route_id": 2001,
    "point_count": 5,
    "created_at": "2026-06-22T10:00:00Z"
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 路线不存在返回错误码 `4002`。
- 路线点位自动复制到 `trip_point` 表，状态为 `pending`。
- 多日路线自动按 `day_number` 分组复制。

---

### 4.19 批量添加点位到行程 @auth

```
POST /api/v1/trips/{trip_id}/points/batch
```

**描述**：批量添加多个地点到行程中。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| points | array | 是 | 点位列表 |
| points[].place_id | int | 是 | 地点 ID |
| points[].day_number | int | 否 | 天数序号，默认 1 |
| points[].point_time | string | 否 | 安排时间，格式：HH:mm |
| points[].stay_duration | int | 否 | 停留时长，单位：分钟 |
| points[].transport | string | 否 | 交通方式 |

**请求示例**：

```json
{
  "points": [
    { "place_id": 1001, "day_number": 1, "point_time": "09:00", "stay_duration": 120 },
    { "place_id": 1002, "day_number": 1, "point_time": "14:00", "stay_duration": 90 }
  ]
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "added_count": 2,
    "total_points": 5
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 行程不存在返回错误码 `5001`。
- 行程点位已达上限返回错误码 `5004`。
- 行程已完成返回错误码 `5009`。

---

### 4.20 重新排序行程点位 @auth

```
PUT /api/v1/trips/{trip_id}/points/reorder
```

**描述**：调整行程点位的排序顺序。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| point_ids | array | 是 | 按新顺序排列的点位 ID 列表 |

**请求示例**：

```json
{
  "point_ids": [3003, 3001, 3002]
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 行程不存在返回错误码 `5001`。
- 点位 ID 列表必须包含行程的所有点位，不允许遗漏。

---

### 4.21 获取行程进度 @auth

```
GET /api/v1/trips/{trip_id}/progress
```

**描述**：获取行程的打卡进度信息。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "trip_id": 5001,
    "total_points": 6,
    "checked_count": 3,
    "skipped_count": 1,
    "pending_count": 2,
    "progress_pct": 50.00,
    "status": "active"
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 分母为零保护：当总点位数量为 0 时，`progress_pct` 显示为 0.00。
- 行程不存在返回错误码 `5001`。

---

### 4.22 分享行程 @auth

```
POST /api/v1/trips/{trip_id}/share
```

**描述**：生成行程分享海报/链接。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 行程 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "share_url": "https://xunlvji.com/share/trip/5001",
    "poster_url": "https://cdn.example.com/posters/trip_5001.png",
    "title": "端午节荆州游",
    "description": "3天2晚 · 5个打卡地 · 预算500元"
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 行程不存在返回错误码 `5001`。
- 非创建者操作返回错误码 `2004`。

---

### 4.23 复制他人行程 @auth

```
POST /api/v1/trips/{trip_id}/copy
```

**描述**：复制他人公开的行程为自己的行程。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| trip_id | int | 是 | 源行程 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 新行程名称，不传则使用原名称 |
| trip_date | string | 是 | 出行日期，格式：YYYY-MM-DD |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5003,
    "name": "端午节荆州游（副本）",
    "trip_date": "2026-07-01",
    "status": "pending",
    "point_count": 5,
    "copied_from": 5001
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 只能复制 `visibility` 为 `public` 的行程，否则返回错误码 `2004`。
- 源行程不存在返回错误码 `5001`。

---

## 五、发布模块

### 5.1 上传图片 @auth

```
POST /api/v1/upload/image
```

**描述**：上传图片文件，返回 CDN URL。使用 `multipart/form-data`。

**请求头**：

```
Content-Type: multipart/form-data
```

**表单参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 图片文件，支持 JPG / PNG / WebP / HEIC，最大 10MB |
| scene | string | 否 | 上传场景：checkin（打卡）/ place（打卡地）/ route（路线）/ avatar（头像），默认 checkin |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "url": "https://cdn.example.com/upload/2026/06/22/abc123.jpg",
    "width": 1920,
    "height": 1080,
    "size": 204800,
    "format": "jpg"
  }
}
```

**说明**：
- 服务端自动转换为 WebP 格式并压缩。
- 单张图片最大 10MB，超过限制返回错误码 `7002`。
- 支持批量上传（多个 file 字段），返回数组格式。
- **限频**：5 次/小时/用户（批量上传按 1 次计），防存储滥用与带宽占用，详见 1.7.1。

**批量上传返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "images": [
      {
        "url": "https://cdn.example.com/upload/2026/06/22/abc123.jpg",
        "width": 1920,
        "height": 1080,
        "size": 204800,
        "format": "jpg"
      },
      {
        "url": "https://cdn.example.com/upload/2026/06/22/def456.jpg",
        "width": 1080,
        "height": 1920,
        "size": 180500,
        "format": "jpg"
      }
    ]
  }
}
```

---

### 5.2 发布打卡地 @auth

```
POST /api/v1/content/publish/place
```

**描述**：用户发布新的打卡地（地点 + 打卡记录）。需经过 AI 预审，信用分低于 60 禁止发布。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 地点名称 |
| category | string | 是 | 分类：food / fun / sightseeing / outdoor / shopping / accommodation |
| longitude | float | 是 | 经度 |
| latitude | float | 是 | 纬度 |
| city | string | 是 | 所在城市，如"荆州"。服务端根据经纬度反查城市，若与传入值不一致以服务端为准 |
| address | string | 否 | 详细地址 |
| opening_hours | string | 否 | 营业时间，如 "09:00-22:00" |
| phone | string | 否 | 联系电话 |
| avg_cost | float | 是 | 人均消费，单位：元 |
| images | array[string] | 是 | 图片 URL 列表，至少 1 张 |
| rating | float | 否 | 评分，范围 0.0-5.0 |
| description | string | 是 | 感受描述 |
| actual_cost | float | 否 | 实际消费 |
| pitfall_tips | array[string] | 是 | 避雷建议，JSON 数组（与数据库一致） |
| recommend_level | string | 否 | 推荐等级：must_go / can_go / avoid |
| suggest_duration | int | 是 | 建议游玩时长，单位：分钟 |
| tags | array[string] | 否 | 标签，最多 5 个 |
| visibility | string | 否 | 可见性：public / private，默认 public |
| extra_data | object | 否 | 差异化专属内容（JSON 对象），按 place.category 存储不同字段，如吃喝的菜品评价、住宿的房型体验等 |

**请求示例**：

```json
{
  "name": "荆州老街咖啡馆",
  "category": "food",
  "longitude": 112.2450,
  "latitude": 30.3350,
  "city": "荆州",
  "address": "湖北省荆州市荆州区老街88号",
  "opening_hours": "09:00-22:00",
  "phone": "0716-1234567",
  "avg_cost": 35.00,
  "images": [
    "https://cdn.example.com/upload/img_001.jpg",
    "https://cdn.example.com/upload/img_002.jpg"
  ],
  "rating": 4.0,
  "description": "老街里的宝藏咖啡馆，手冲咖啡很棒，环境也很安静适合拍照。",
  "actual_cost": 38.00,
  "pitfall_tips": ["周末下午人比较多，建议工作日去"],
  "recommend_level": "must_go",
  "suggest_duration": 60,
  "tags": ["咖啡", "拍照", "安静"],
  "visibility": "public",
  "extra_data": {
    "dish_review": "手冲咖啡很棒"
  }
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "place_id": 3013,
    "checkin_record_id": 6010,
    "status": "pending",
    "audit_status": "审核中",
    "credit_change": 5
  }
}
```

**说明**：
- `audit_status` 为 `审核中` 时，内容仅自己可见，审核通过后公开。
- 信用分审核规则（完整，初始分 100）：
  - 信用分 ≥ 90：免审直接发布。
  - 80 ≤ 信用分 < 90（80-89）：快速过审（30 分钟内完成审核）。
  - 70 ≤ 信用分 < 80（70-79）：抽审。
  - 60 ≤ 信用分 < 70（60-69）：强制审核。
  - 信用分 < 60：返回错误码 `3007`（信用分过低，禁止发布）。
- 内容包含敏感词时返回错误码 `4007`。

---

### 5.3 获取我的发布列表 @auth

```
GET /api/v1/content/publish/my-list
```

**描述**：获取当前用户发布的所有打卡地记录。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| status | string | 否 | - | 审核状态：pending（审核中）/ approved（已通过）/ rejected（已驳回），不传则全部 |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "checkin_record_id": 6010,
        "place_id": 3013,
        "place_name": "荆州老街咖啡馆",
        "cover_image": "https://cdn.example.com/upload/img_001.jpg",
        "category": "food",
        "category_name": "吃喝",
        "rating": 4.0,
        "status": "pending",
        "status_name": "审核中",
        "visibility": "public",
        "reject_reason": null,
        "created_at": "2026-06-22T12:00:00Z"
      },
      {
        "checkin_record_id": 6011,
        "place_id": 3014,
        "place_name": "沙市洋码头夜市",
        "cover_image": "https://cdn.example.com/upload/img_003.jpg",
        "category": "food",
        "category_name": "吃喝",
        "rating": 3.5,
        "status": "rejected",
        "status_name": "已驳回",
        "visibility": "public",
        "reject_reason": "图片模糊，请重新上传清晰图片",
        "created_at": "2026-06-20T18:00:00Z"
      }
    ],
    "total": 15,
    "page": 1,
    "page_size": 20,
    "has_more": false
  }
}
```

**说明**：
- `reject_reason` 仅在 `status` 为 `rejected` 时有值，展示已驳回原因。
- 已驳回的内容可修改后重新提交。

---

### 5.4 重新提交审核 @auth

```
POST /api/v1/content/publish/{checkin_id}/resubmit
```

**描述**：修改已驳回的打卡记录后重新提交审核。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| checkin_id | int | 是 | 打卡记录 ID |

**请求参数**：同 5.2 发布打卡地，所有字段均为可选，仅更新传入的字段。

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "checkin_record_id": 6011,
    "status": "pending",
    "audit_status": "审核中"
  }
}
```

---

## 六、社交模块

### 6.1 关注用户 @auth

```
POST /api/v1/social/follow
```

**描述**：关注指定用户。不能关注自己。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | int | 是 | 要关注的用户 ID |

**请求示例**：

```json
{
  "user_id": 10002
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_following": true,
    "follower_count": 521
  }
}
```

**说明**：
- 已关注则不做任何操作，仍返回成功。
- 关注自己返回错误码 `6009`。

---

### 6.2 取消关注 @auth

```
POST /api/v1/social/unfollow
```

**描述**：取消关注指定用户。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | int | 是 | 要取消关注的用户 ID |

**请求示例**：

```json
{
  "user_id": 10002
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_following": false,
    "follower_count": 520
  }
}
```

---

### 6.3 获取关注列表 @auth

```
GET /api/v1/social/following/list
```

**描述**：获取当前用户关注的人列表。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 10002,
        "nickname": "旅行达人小王",
        "avatar": "https://cdn.example.com/avatars/10002.jpg",
        "follower_count": 520,
        "is_following": true,
        "followed_at": "2026-05-10T15:00:00Z"
      }
    ],
    "total": 56,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 6.4 获取粉丝列表 @auth

```
GET /api/v1/social/follower/list
```

**描述**：获取当前用户的粉丝列表。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 10003,
        "nickname": "John",
        "avatar": "https://cdn.example.com/avatars/10003.jpg",
        "follower_count": 120,
        "is_following": false,
        "followed_at": "2026-06-01T08:00:00Z"
      }
    ],
    "total": 128,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 6.5 收藏 @auth

```
POST /api/v1/social/favorite
```

**描述**：收藏路线、打卡地、攻略或行程。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_type | string | 是 | 目标类型：route / place / guide / trip |
| target_id | int | 是 | 目标 ID |

**请求示例**：

```json
{
  "target_type": "route",
  "target_id": 2001
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_favorited": true,
    "favorite_count": 151
  }
}
```

---

### 6.6 取消收藏 @auth

```
POST /api/v1/social/unfavorite
```

**描述**：取消收藏。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_type | string | 是 | 目标类型：route / place / guide / trip |
| target_id | int | 是 | 目标 ID |

**请求示例**：

```json
{
  "target_type": "route",
  "target_id": 2001
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_favorited": false,
    "favorite_count": 150
  }
}
```

---

### 6.7 获取收藏列表 @auth

```
GET /api/v1/social/favorite/list
```

**描述**：获取当前用户的收藏列表。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| target_type | string | 否 | - | 筛选类型：route / place / guide / trip，不传则全部 |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "target_type": "route",
        "target_id": 2001,
        "title": "荆州古城一日游",
        "cover_image": "https://cdn.example.com/routes/2001.jpg",
        "category": "one_day",
        "category_name": "一日游",
        "rating": 4.5,
        "favorited_at": "2026-06-20T15:00:00Z"
      }
    ],
    "total": 30,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 6.8 发布评论 @auth

```
POST /api/v1/social/comment
```

**描述**：发布一级评论或二级回复（楼中楼）。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_type | string | 是 | 目标类型：route / place / guide |
| target_id | int | 是 | 目标 ID |
| content | string | 是 | 评论内容，1-500 个字符 |
| parent_id | int | 否 | 父评论 ID（二级回复时必填，NULL 表示一级评论） |

**请求示例**：

```json
{
  "target_type": "route",
  "target_id": 2001,
  "content": "这条路线很棒，上周末刚走过！",
  "parent_id": null
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 7001,
    "content": "这条路线很棒，上周末刚走过！",
    "user": {
      "id": 10001,
      "nickname": "小明爱旅行",
      "avatar": "https://cdn.example.com/avatars/10001.jpg"
    },
    "parent_id": null,
    "created_at": "2026-06-22T14:00:00Z"
  }
}
```

**说明**：
- 评论内容不能为空，返回错误码 `6003`。
- 仅支持二级评论，`parent_id` 指向的评论如果已有 `parent_id`（即已是二级回复），返回错误码 `1001`。
- 内容包含敏感词时返回错误码 `4007`。
- **限频**：30 次/小时/用户，防刷评/水军，详见 1.7.1。

---

### 6.9 获取评论列表

```
GET /api/v1/social/comment/list
```

**描述**：获取指定目标的评论列表（楼中楼结构）。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| target_type | string | 是 | - | 目标类型：route / place / guide |
| target_id | int | 是 | - | 目标 ID |
| sort | string | 否 | newest | 排序：newest（最新）/ hottest（最热） |
| page | int | 否 | 1 | 页码（一级评论分页） |
| page_size | int | 否 | 20 | 每页数量（一级评论数量） |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 7001,
        "content": "这条路线很棒，上周末刚走过！",
        "user": {
          "id": 10001,
          "nickname": "小明爱旅行",
          "avatar": "https://cdn.example.com/avatars/10001.jpg"
        },
        "reply_count": 2,
        "replies": [
          {
            "id": 7002,
            "content": "是的，我也走过，推荐！",
            "user": {
              "id": 10002,
              "nickname": "旅行达人小王",
              "avatar": "https://cdn.example.com/avatars/10002.jpg"
            },
            "reply_to": {
              "id": 10001,
              "nickname": "小明爱旅行"
            },
            "created_at": "2026-06-22T15:00:00Z"
          },
          {
            "id": 7003,
            "content": "请问周末人多吗？",
            "user": {
              "id": 10003,
              "nickname": "John",
              "avatar": "https://cdn.example.com/avatars/10003.jpg"
            },
            "reply_to": {
              "id": 10001,
              "nickname": "小明爱旅行"
            },
            "created_at": "2026-06-22T16:00:00Z"
          }
        ],
        "created_at": "2026-06-22T14:00:00Z"
      }
    ],
    "total": 12,
    "page": 1,
    "page_size": 20,
    "has_more": false
  }
}
```

**说明**：
- 分页以一级评论为单位，每页返回指定数量的一级评论。
- 每条一级评论附带最多 3 条二级回复预览，更多回复需单独请求。

---

### 6.10 获取二级回复列表

```
GET /api/v1/social/comment/replies
```

**描述**：获取指定一级评论的所有二级回复。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| comment_id | int | 是 | - | 一级评论 ID |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 7002,
        "content": "是的，我也走过，推荐！",
        "user": {
          "id": 10002,
          "nickname": "旅行达人小王",
          "avatar": "https://cdn.example.com/avatars/10002.jpg"
        },
        "reply_to": {
          "id": 10001,
          "nickname": "小明爱旅行"
        },
        "created_at": "2026-06-22T15:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20,
    "has_more": false
  }
}
```

---

### 6.11 删除评论 @auth

```
DELETE /api/v1/social/comment/{comment_id}
```

**描述**：删除自己的评论。删除一级评论会级联删除所有二级回复。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| comment_id | int | 是 | 评论 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

**说明**：
- 只能删除自己的评论，删除他人评论返回错误码 `2004`。

---

### 6.12 避雷"有用"点赞 @auth

```
POST /api/v1/social/tip-vote
```

**描述**：对打卡记录中的某条避雷建议点"有用"。每个用户对每条避雷只能点赞一次。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| checkin_id | int | 是 | 打卡记录 ID |
| tip_index | int | 是 | 避雷建议的索引（0-based），对应打卡记录中 pitfall_tips 的索引 |

**请求示例**：

```json
{
  "checkin_id": 6001,
  "tip_index": 0
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_voted": true,
    "vote_count": 46
  }
}
```

**说明**：
- 已点赞则不做任何操作，仍返回成功。
- 不能给自己的避雷建议点赞，返回错误码 `3006`。

---

### 6.13 取消避雷点赞 @auth

```
POST /api/v1/social/tip-unvote
```

**描述**：取消对避雷建议的"有用"点赞。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| checkin_id | int | 是 | 打卡记录 ID |
| tip_index | int | 是 | 避雷建议的索引（0-based） |

**请求示例**：

```json
{
  "checkin_id": 6001,
  "tip_index": 0
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_voted": false,
    "vote_count": 45
  }
}
```

---

### 6.14 获取消息列表 @auth

```
GET /api/v1/social/message/list
```

**描述**：获取当前用户的消息列表，包含系统通知、私信、群聊，按时间倒序排列。P0 阶段仅返回 system 类型消息，private/group 类型在 P1 阶段启用。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| type | string | 否 | - | 消息类型：system（系统通知）/ private（私信）/ group（群聊），不传则全部 |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 8001,
        "type": "system",
        "type_name": "系统通知",
        "icon": "system",
        "title": "审核结果通知",
        "content": "您发布的「荆州老街咖啡馆」已通过审核",
        "is_read": false,
        "created_at": "2026-06-22T13:00:00Z"
      },
      {
        "id": 8002,
        "type": "private",
        "type_name": "私信",
        "icon": "chat",
        "sender": {
          "id": 10002,
          "nickname": "旅行达人小王",
          "avatar": "https://cdn.example.com/avatars/10002.jpg"
        },
        "content": "你好，请问荆州古城墙需要门票吗？",
        "is_read": true,
        "created_at": "2026-06-21T20:00:00Z"
      },
      {
        "id": 8003,
        "type": "group",
        "type_name": "群聊",
        "icon": "group",
        "title": "荆州旅行爱好者群",
        "content": "小明爱旅行：周末有人一起去荆州吗？",
        "is_read": false,
        "created_at": "2026-06-22T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20,
    "has_more": true,
    "unread_count": 3
  }
}
```

**说明**：
- `unread_count` 为当前用户所有未读消息总数。
- 消息列表不分栏不分组，用 `icon` 区分类型。

---

### 6.15 标记消息已读 @auth

```
POST /api/v1/social/message/read
```

**描述**：标记指定消息为已读，或全部已读。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| message_ids | array[int] | 否 | 消息 ID 列表，不传则全部已读 |

**请求示例**：

```json
{
  "message_ids": [8001, 8003]
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "read_count": 2
  }
}
```

---

### 6.16 发送私信 @auth（P1）

```
POST /api/v1/social/message/send
```

**描述**：向指定用户发送私信。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| receiver_id | int | 是 | 接收者用户 ID |
| content | string | 是 | 消息内容，1-500 个字符 |

**请求示例**：

```json
{
  "receiver_id": 10002,
  "content": "你好，请问荆州古城墙需要门票吗？"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 8002,
    "type": "private",
    "content": "你好，请问荆州古城墙需要门票吗？",
    "created_at": "2026-06-21T20:00:00Z"
  }
}
```

---

### 6.17 获取私信对话 @auth（P1）

```
GET /api/v1/social/message/conversation
```

**描述**：获取与指定用户的私信对话记录。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| target_user_id | int | 是 | - | 对方用户 ID |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 30 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "target_user": {
      "id": 10002,
      "nickname": "旅行达人小王",
      "avatar": "https://cdn.example.com/avatars/10002.jpg"
    },
    "list": [
      {
        "id": 8004,
        "sender_id": 10001,
        "content": "你好",
        "created_at": "2026-06-21T19:00:00Z"
      },
      {
        "id": 8005,
        "sender_id": 10002,
        "content": "你好！有什么可以帮你的吗？",
        "created_at": "2026-06-21T19:05:00Z"
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 30,
    "has_more": false
  }
}
```

---

### 6.18 发送群聊消息 @auth（P1）

```
POST /api/v1/social/message/group-send
```

**描述**：向群聊发送消息。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| group_id | int | 是 | 群聊 ID |
| content | string | 是 | 消息内容，1-500 个字符 |

**请求示例**：

```json
{
  "group_id": 1001,
  "content": "周末有人一起去荆州吗？"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 8003,
    "type": "group",
    "content": "周末有人一起去荆州吗？",
    "created_at": "2026-06-22T10:00:00Z"
  }
}
```

---

### 6.19 获取群聊消息 @auth（P1）

```
GET /api/v1/social/message/group-conversation
```

**描述**：获取指定群聊的聊天记录。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| group_id | int | 是 | - | 群聊 ID |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 30 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "group": {
      "id": 1001,
      "name": "荆州旅行爱好者群",
      "member_count": 128
    },
    "list": [
      {
        "id": 8003,
        "sender": {
          "id": 10001,
          "nickname": "小明爱旅行",
          "avatar": "https://cdn.example.com/avatars/10001.jpg"
        },
        "content": "周末有人一起去荆州吗？",
        "created_at": "2026-06-22T10:00:00Z"
      }
    ],
    "total": 200,
    "page": 1,
    "page_size": 30,
    "has_more": true
  }
}
```

---

### 6.20 删除消息 @auth

```
DELETE /api/v1/social/message/{message_id}
```

**描述**：删除指定的消息（仅删除自己的记录，对方仍可见）。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| message_id | int | 是 | 消息 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 6.21 评论点赞 @auth

```
POST /api/v1/social/comments/{id}/like
```

**描述**：对评论进行点赞或取消点赞（Toggle 机制），已点赞则取消，未点赞则点赞。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 评论 ID |

**请求体**：无（空请求体）。

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "liked": true,
    "like_count": 42
  }
}
```

**说明**：
- 点赞失败时返回错误码 `1001`。
- 评论不存在或已删除时返回错误码 `6001`。

---

### 6.22 路线相关推荐

```
GET /api/v1/content/routes/{id}/related
```

**描述**：获取与指定路线相关的推荐路线列表，基于路线主题、区域、风格等维度进行关联推荐。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 路线 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 2002,
        "title": "荆州古城深度游",
        "cover_image": "https://cdn.example.com/routes/2002.jpg",
        "category": "half_day",
        "category_name": "半日游",
        "distance": 5.2,
        "duration": 180,
        "point_count": 6,
        "rating": 4.5,
        "city": "荆州"
      },
      {
        "id": 2003,
        "title": "荆州博物馆文化之旅",
        "cover_image": "https://cdn.example.com/routes/2003.jpg",
        "category": "one_day",
        "category_name": "一日游",
        "distance": 3.8,
        "duration": 150,
        "point_count": 4,
        "rating": 4.3,
        "city": "荆州"
      }
    ]
  }
}
```

---

### 6.23 路线附近打卡地

```
GET /api/v1/content/routes/{id}/nearby
```

**描述**：获取指定路线附近的打卡地列表，方便用户发现沿途可探访的地点。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 路线 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| longitude | float | 是 | - | 当前经度 |
| latitude | float | 是 | - | 当前纬度 |
| radius | int | 否 | 5000 | 搜索半径（米），默认 5000 米 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 3005,
        "name": "古城墙咖啡馆",
        "cover_image": "https://cdn.example.com/places/3005.jpg",
        "category": "food",
        "category_name": "好吃",
        "rating": 4.6,
        "distance": 320,
        "address": "荆州市荆州区古城墙旁",
        "avg_cost": 35
      },
      {
        "id": 3006,
        "name": "三国文化书店",
        "cover_image": "https://cdn.example.com/places/3006.jpg",
        "category": "shopping",
        "category_name": "好买",
        "rating": 4.2,
        "distance": 580,
        "address": "荆州市荆州区东大街",
        "avg_cost": 0
      }
    ]
  }
}
```

---

### 6.24 打卡地附近推荐

```
GET /api/v1/content/places/{id}/nearby
```

**描述**：获取指定打卡地附近的其他打卡地推荐，帮助用户发现周边可探访的地点。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 打卡地 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| radius | int | 否 | 3000 | 搜索半径（米），默认 3000 米 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 3007,
        "name": "老字号荆州鱼糕",
        "cover_image": "https://cdn.example.com/places/3007.jpg",
        "category": "food",
        "category_name": "好吃",
        "rating": 4.8,
        "distance": 150,
        "address": "荆州市荆州区解放路 88 号",
        "avg_cost": 50
      },
      {
        "id": 3008,
        "name": "荆州手工艺坊",
        "cover_image": "https://cdn.example.com/places/3008.jpg",
        "category": "shopping",
        "category_name": "好买",
        "rating": 4.3,
        "distance": 420,
        "address": "荆州市荆州区文化街 12 号",
        "avg_cost": 0
      }
    ]
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

---

### 6.25 举报内容/用户 @auth

```
POST /api/v1/social/report
```

**描述**：举报违规内容或用户，提交后进入审核队列。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_type | string | 是 | 举报目标类型：place / route / guide / checkin / comment / user |
| target_id | int | 是 | 目标 ID |
| reason | string | 是 | 举报原因：spam（垃圾广告）/ abuse（辱骂攻击）/ porn（色情低俗）/ fraud（欺诈）/ other（其他） |
| description | string | 否 | 补充说明，最多 500 字 |
| evidence | array[string] | 否 | 证据截图 URL 数组，最多 5 张 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "report_id": 9001,
    "status": "pending"
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- 同一用户对同一目标 24 小时内只能举报一次。
- 举报提交后由后台审核团队处理，处理结果通过消息通知用户。

---

### 6.26 获取未读消息数 @auth

```
GET /api/v1/social/unread-count
```

**描述**：获取当前用户的未读消息总数，按类型分类统计。

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "system": 3,
    "private": 5,
    "group": 7,
    "like": 0,
    "comment": 0,
    "follow": 0
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- `total` 为所有未读消息总数。
- `like` / `comment` / `follow` 为互动类未读通知数。
- 客户端可使用此接口实现消息红点提示，建议每 60 秒轮询一次或使用 WebSocket 推送。

---

## 七、攻略问答模块

### 7.1 攻略问答列表

```
GET /api/v1/content/guides/{id}/qa
```

**描述**：获取指定攻略下的问答列表，支持分页。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 攻略 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 5001,
        "question": "这条路线适合带老人一起走吗？",
        "answer": "适合的，路线整体比较平缓，沿途也有休息点。",
        "answerer": {
          "id": 10002,
          "nickname": "荆州本地人",
          "avatar": "https://cdn.example.com/avatars/10002.jpg"
        },
        "like_count": 12,
        "is_liked": false,
        "created_at": "2026-06-20T09:00:00Z"
      },
      {
        "id": 5002,
        "question": "沿途有推荐的餐厅吗？",
        "answer": "古城墙附近有几家老字号，推荐荆楚风味馆。",
        "answerer": {
          "id": 10003,
          "nickname": "吃货小分队",
          "avatar": "https://cdn.example.com/avatars/10003.jpg"
        },
        "like_count": 8,
        "is_liked": true,
        "created_at": "2026-06-19T15:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 7.2 提问 @auth

```
POST /api/v1/content/guides/{id}/qa
```

**描述**：在指定攻略下提出一个问题。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 攻略 ID |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question | string | 是 | 问题内容，1-500 个字符 |

**请求示例**：

```json
{
  "question": "这条路线适合带老人一起走吗？"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5001,
    "question": "这条路线适合带老人一起走吗？",
    "created_at": "2026-06-22T14:30:00Z"
  }
}
```

**说明**：
- 问题内容不能为空，返回错误码 `1002`。
- 攻略不存在或已下架时返回错误码 `4003`。
- 内容包含敏感词时返回错误码 `4007`。

---

### 7.3 回答 @auth

```
POST /api/v1/content/guides/{id}/qa/{question_id}/answer
```

**描述**：对攻略下的某个问题作出回答。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 攻略 ID |
| question_id | int | 是 | 问题 ID（guide_question 表） |

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| answer | string | 是 | 回答内容，1-1000 个字符 |

**请求示例**：

```json
{
  "answer": "适合的，路线整体比较平缓，沿途也有休息点，老人走起来不会太累。"
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5001,
    "answer": "适合的，路线整体比较平缓，沿途也有休息点，老人走起来不会太累。",
    "answered_by": {
      "id": 10002,
      "nickname": "荆州本地人",
      "avatar": "https://cdn.example.com/avatars/10002.jpg"
    },
    "answer_at": "2026-06-22T15:00:00Z"
  }
}
```

**说明**：
- 回答内容不能为空，返回错误码 `1002`。
- 问答已被删除时返回错误码 `1004`。

---

### 7.4 点赞回答 @auth

```
POST /api/v1/content/guides/{id}/qa/{answer_id}/like
```

**描述**：对攻略问答的回答进行点赞或取消点赞（Toggle 机制）。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 攻略 ID |
| answer_id | int | 是 | 回答 ID（guide_answer 表） |

**请求体**：无（空请求体）。

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "liked": true,
    "like_count": 13
  }
}
```

**说明**：
- 问答已被删除时返回错误码 `1004`。

---

## 八、其他模块

### 8.1 城市列表

```
GET /api/v1/common/cities
```

**描述**：获取系统支持的城市列表，用于城市选择。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| keyword | string | 否 | - | 搜索关键词，支持拼音和中文 |
| level | string | 否 | hot | 城市级别：hot（热门城市）/ all（全部） |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "hot_cities": [
      {
        "name": "荆州",
        "pinyin": "jingzhou",
        "place_count": 120,
        "route_count": 45,
        "guide_count": 10
      },
      {
        "name": "武汉",
        "pinyin": "wuhan",
        "place_count": 500,
        "route_count": 200,
        "guide_count": 50
      }
    ],
    "all_cities": []
  }
}
```

---

### 8.2 榜单（P1）

```
GET /api/v1/common/ranking
```

**描述**：获取打卡榜、好评榜、消费榜。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| type | string | 是 | - | 榜单类型：checkin（打卡榜）/ rating（好评榜）/ expense（消费榜） |
| city | string | 否 | - | 筛选城市，不传则全国 |
| category | string | 否 | - | 筛选分类，不传则全部 |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "type": "checkin",
    "type_name": "打卡榜",
    "city": "荆州",
    "list": [
      {
        "rank": 1,
        "place_id": 3001,
        "name": "荆州古城墙",
        "cover_image": "https://cdn.example.com/places/3001.jpg",
        "category": "sightseeing",
        "category_name": "逛看",
        "rating": 4.7,
        "checkin_count": 520,
        "avg_cost": 0
      },
      {
        "rank": 2,
        "place_id": 3002,
        "name": "荆州博物馆",
        "cover_image": "https://cdn.example.com/places/3002.jpg",
        "category": "sightseeing",
        "category_name": "逛看",
        "rating": 4.5,
        "checkin_count": 380,
        "avg_cost": 0
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

**说明**：
- 打卡榜：按 `checkin_count` 降序排列。
- 好评榜：按 `rating` 降序排列，同分按 `checkin_count` 降序。
- 消费榜：按 `avg_cost` 降序排列。

---

### 8.3 足迹地图数据 @auth（P1）

```
GET /api/v1/common/footprint
```

**描述**：获取当前用户的足迹数据，用于足迹地图展示。

**请求参数**：无

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "visited_cities": [
      {
        "city": "荆州",
        "checkin_count": 15,
        "last_visit": "2026-06-20"
      },
      {
        "city": "武汉",
        "checkin_count": 8,
        "last_visit": "2026-05-15"
      }
    ],
    "total_cities": 12,
    "total_checkins": 45,
    "total_places": 38,
    "checkin_points": [
      {
        "place_id": 3001,
        "name": "荆州古城墙",
        "longitude": 112.2400,
        "latitude": 30.3300,
        "category": "sightseeing",
        "checkin_time": "2026-06-20T09:00:00Z"
      },
      {
        "place_id": 3005,
        "name": "黄鹤楼",
        "longitude": 114.3025,
        "latitude": 30.5446,
        "category": "sightseeing",
        "checkin_time": "2026-05-15T14:00:00Z"
      }
    ],
    "unlocked_badges": [
      {
        "id": 1,
        "name": "荆州探索者",
        "icon": "https://cdn.example.com/badges/jingzhou.png",
        "description": "在荆州打卡 10 个地点",
        "unlocked_at": "2026-06-20T12:00:00Z"
      }
    ]
  }
}
```

**说明**：
- `checkin_points` 为所有打卡点的坐标信息，用于地图标注。
- `unlocked_badges` 为已解锁的勋章列表。

---

### 8.4 AI 助手对话 @auth（P1）

```
POST /api/v1/ai/chat
```

**描述**：与 AI 助手「寻旅小助手」对话，支持旅行相关问题咨询。（P1 功能，P0 阶段不开发此接口）

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| message | string | 是 | 用户消息，1-1000 个字符 |
| session_id | string | 否 | 会话 ID，首次对话不传，后续传上次返回的 session_id 以保持上下文 |
| context | object | 否 | 额外上下文信息 |

**context 对象结构**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| city | string | 否 | 当前所在/选中的城市 |
| longitude | float | 否 | 当前经度 |
| latitude | float | 否 | 当前纬度 |
| trip_id | int | 否 | 当前行程 ID（如在行程中提问） |

**请求示例**：

```json
{
  "message": "荆州有哪些必去的景点？",
  "session_id": "sess_abc123",
  "context": {
    "city": "荆州"
  }
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "session_id": "sess_abc123",
    "reply": "荆州作为历史文化名城，必去的景点有：\n\n1. **荆州古城墙** - 中国现存最完好的古城墙之一，建议游玩 1.5 小时\n2. **荆州博物馆** - 馆藏丰富，可了解荆州历史，建议游玩 2 小时\n3. **张居正故居** - 明代名臣故居，感受历史底蕴\n\n需要我帮你规划一条荆州一日游路线吗？",
    "suggestions": [
      "帮我规划一条荆州一日游路线",
      "荆州有哪些好吃的美食？",
      "荆州古城墙门票多少钱？"
    ],
    "created_at": "2026-06-22T12:00:00Z"
  }
}
```

**说明**：
- `session_id` 用于保持对话上下文，有效期 24 小时。
- `suggestions` 为 AI 推荐的后续问题，帮助用户继续对话。
- AI 生成内容需在前端标注"AI 生成，仅供参考"。

---

### 8.5 获取 AI 对话历史 @auth（P1）

```
GET /api/v1/ai/chat/history
```

**描述**：获取与 AI 助手的对话历史。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| session_id | string | 否 | - | 会话 ID，不传则返回所有会话列表 |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**会话列表返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "session_id": "sess_abc123",
        "title": "荆州必去景点咨询",
        "last_message": "需要我帮你规划一条荆州一日游路线吗？",
        "message_count": 6,
        "updated_at": "2026-06-22T12:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20,
    "has_more": false
  }
}
```

**对话详情返回示例**（传入 session_id）：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "session_id": "sess_abc123",
    "title": "荆州必去景点咨询",
    "messages": [
      {
        "role": "user",
        "content": "荆州有哪些必去的景点？",
        "created_at": "2026-06-22T11:58:00Z"
      },
      {
        "role": "assistant",
        "content": "荆州作为历史文化名城...",
        "suggestions": ["帮我规划一条荆州一日游路线"],
        "created_at": "2026-06-22T11:58:05Z"
      }
    ],
    "created_at": "2026-06-22T11:58:00Z"
  }
}
```

---

### 8.6 删除 AI 对话会话 @auth（P1）

```
DELETE /api/v1/ai/chat/session/{session_id}
```

**描述**：删除指定的 AI 对话会话及所有历史消息。

**路径参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| session_id | string | 是 | 会话 ID |

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

### 8.7 获取首页配置

```
GET /api/v1/common/home-config
```

**描述**：获取首页全局配置，如 Tab 列表、默认城市、Banner 等。

**请求参数**：无

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tabs": [
      {"key": "recommend", "name": "推荐"},
      {"key": "following", "name": "关注"},
      {"key": "nearby", "name": "附近"},
      {"key": "guide", "name": "攻略"},
      {"key": "place", "name": "打卡地"},
      {"key": "route", "name": "路线"}
    ],
    "default_city": "荆州",
    "banners": [
      {
        "image": "https://cdn.example.com/banners/1.jpg",
        "link": "/route/2001",
        "link_type": "route"
      }
    ],
    "guide_categories": [
      {"key": "", "name": "全部"},
      {"key": "nature", "name": "自然生态"},
      {"key": "history", "name": "历史文化"},
      {"key": "entertainment", "name": "人工娱乐"},
      {"key": "urban", "name": "城市公共"},
      {"key": "transport", "name": "交通枢纽"},
      {"key": "red_tourism", "name": "红色旅游"},
      {"key": "religion", "name": "宗教场所"}
    ],
    "place_categories": [
      {"key": "", "name": "全部"},
      {"key": "food", "name": "吃喝"},
      {"key": "fun", "name": "玩乐"},
      {"key": "sightseeing", "name": "逛看"},
      {"key": "outdoor", "name": "户外"},
      {"key": "shopping", "name": "逛街"},
      {"key": "accommodation", "name": "住宿"}
    ]
  }
}
```

---

### 8.8 刷新 Token @auth

```
POST /api/v1/auth/refresh-token
```

**描述**：使用 Refresh Token 刷新 Access Token，避免频繁登录。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| refresh_token | string | 是 | 登录时返回的 Refresh Token |

**请求示例**：

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expire_at": "2026-06-24T14:00:00Z",
    "refresh_expire_at": "2026-07-24T12:00:00Z"
  },
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**说明**：
- Access Token 有效期 2 小时，Refresh Token 有效期 30 天。
- Refresh Token 过期后需重新登录。

---

## 附录 A：接口索引

| 序号 | 接口 | 方法 | 路径 | 认证 |
|------|------|------|------|------|
| 2.1 | 发送验证码 | POST | /api/v1/auth/send-code | - |
| 2.2 | 手机号验证码登录 | POST | /api/v1/auth/login/phone | - |
| 2.3 | Apple ID 登录 | POST | /api/v1/auth/login/apple | - |
| 2.4 | Google 登录 | POST | /api/v1/auth/login/google | - |
| 2.5 | Email 登录 | POST | /api/v1/auth/login/email | - |
| 2.6 | 微信登录 | POST | /api/v1/auth/login/wechat | - |
| 2.7 | 获取用户信息 | GET | /api/v1/user/profile | @auth |
| 2.8 | 获取其他用户信息 | GET | /api/v1/user/{user_id}/profile | - |
| 2.9 | 更新用户资料 | PUT | /api/v1/user/profile | @auth |
| 2.10 | 绑定手机号 | POST | /api/v1/user/bind-phone | @auth |
| 2.11 | 绑定邮箱 | POST | /api/v1/user/bind-email | @auth |
| 2.12 | 账号注销 | POST | /api/v1/user/delete-account | @auth |
| 2.13 | 修改密码 | PUT | /api/v1/user/password | @auth |
| 2.14 | 用户设置 | GET/PUT | /api/v1/user/settings | @auth |
| 3.1 | 首页推荐流 | GET | /api/v1/content/feed/recommend | - |
| 3.2 | 关注流 | GET | /api/v1/content/feed/following | @auth |
| 3.3 | 附近流 | GET | /api/v1/content/feed/nearby | - |
| 3.4 | 城市内容流 | GET | /api/v1/content/feed/city | - |
| 3.5 | 攻略内容流 | GET | /api/v1/content/feed/guide | - |
| 3.6 | 打卡地内容流 | GET | /api/v1/content/feed/place | - |
| 3.7 | 路线内容流 | GET | /api/v1/content/feed/route | - |
| 3.8 | 路线详情 | GET | /api/v1/content/route/{route_id} | - |
| 3.9 | 打卡地详情 | GET | /api/v1/content/place/{place_id} | - |
| 3.10 | 攻略详情 | GET | /api/v1/content/guide/{guide_id} | - |
| 3.11 | 搜索 | GET | /api/v1/content/search | - |
| 3.12 | 搜索建议 | GET | /api/v1/content/search/suggest | - |
| 3.13 | 创建路线 | POST | /api/v1/content/routes | @auth |
| 3.14 | 打卡记录详情 | GET | /api/v1/content/checkin/{checkin_id} | - |
| 3.15 | 获取 CPS 分销链接 | GET | /api/v1/places/{place_id}/cps | @auth |
| 4.1 | 创建行程 | POST | /api/v1/trips/create | @auth |
| 4.2 | 获取行程列表 | GET | /api/v1/trips/list | @auth |
| 4.3 | 行程详情 | GET | /api/v1/trips/{trip_id} | - |
| 4.4 | 开始行程 | POST | /api/v1/trips/{trip_id}/start | @auth |
| 4.5 | 打卡（到店打卡） | POST | /api/v1/trips/{trip_id}/checkin/store | @auth |
| 4.6 | 打卡（普通打卡） | POST | /api/v1/trips/{trip_id}/checkin/normal | @auth |
| 4.7 | 离线打卡同步 | POST | /api/v1/trips/{trip_id}/checkin/sync | @auth |
| 4.8 | 跳过点位 | POST | /api/v1/trips/{trip_id}/point/{point_id}/skip | @auth |
| 4.9 | 取消跳过 | POST | /api/v1/trips/{trip_id}/point/{point_id}/unskip | @auth |
| 4.10 | 删除行程 | DELETE | /api/v1/trips/{trip_id} | @auth |
| 4.11 | AI 行程规划 | POST | /api/v1/trips/ai-plan | @auth（P1） |
| 4.12 | 添加点位 | POST | /api/v1/trips/{trip_id}/points | @auth |
| 4.13 | 删除点位 | DELETE | /api/v1/trips/{trip_id}/point/{point_id} | @auth |
| 4.14 | 编辑行程基本信息 | PUT | /api/v1/trips/{trip_id} | @auth |
| 4.15 | 结束行程 | POST | /api/v1/trips/{trip_id}/complete | @auth |
| 4.16 | 行程清单管理 | GET/POST/PUT/DELETE | /api/v1/trips/{trip_id}/checklist | @auth |
| 4.17 | 行程账单管理 | GET/POST/DELETE | /api/v1/trips/{trip_id}/expenses | @auth（P1） |
| 5.1 | 上传图片 | POST | /api/v1/upload/image | @auth |
| 5.2 | 发布打卡地 | POST | /api/v1/content/publish/place | @auth |
| 5.3 | 获取我的发布列表 | GET | /api/v1/content/publish/my-list | @auth |
| 5.4 | 重新提交审核 | POST | /api/v1/content/publish/{checkin_id}/resubmit | @auth |
| 6.1 | 关注用户 | POST | /api/v1/social/follow | @auth |
| 6.2 | 取消关注 | POST | /api/v1/social/unfollow | @auth |
| 6.3 | 获取关注列表 | GET | /api/v1/social/following/list | @auth |
| 6.4 | 获取粉丝列表 | GET | /api/v1/social/follower/list | @auth |
| 6.5 | 收藏 | POST | /api/v1/social/favorite | @auth |
| 6.6 | 取消收藏 | POST | /api/v1/social/unfavorite | @auth |
| 6.7 | 获取收藏列表 | GET | /api/v1/social/favorite/list | @auth |
| 6.8 | 发布评论 | POST | /api/v1/social/comment | @auth |
| 6.9 | 获取评论列表 | GET | /api/v1/social/comment/list | - |
| 6.10 | 获取二级回复 | GET | /api/v1/social/comment/replies | - |
| 6.11 | 删除评论 | DELETE | /api/v1/social/comment/{comment_id} | @auth |
| 6.12 | 避雷点赞 | POST | /api/v1/social/tip-vote | @auth |
| 6.13 | 取消避雷点赞 | POST | /api/v1/social/tip-unvote | @auth |
| 6.14 | 获取消息列表 | GET | /api/v1/social/message/list | @auth |
| 6.15 | 标记消息已读 | POST | /api/v1/social/message/read | @auth |
| 6.16 | 发送私信 | POST | /api/v1/social/message/send | @auth（P1） |
| 6.17 | 获取私信对话 | GET | /api/v1/social/message/conversation | @auth（P1） |
| 6.18 | 发送群聊消息 | POST | /api/v1/social/message/group-send | @auth（P1） |
| 6.19 | 获取群聊消息 | GET | /api/v1/social/message/group-conversation | @auth（P1） |
| 6.20 | 删除消息 | DELETE | /api/v1/social/message/{message_id} | @auth |
| 6.21 | 评论点赞 | POST | /api/v1/social/comments/{id}/like | @auth |
| 6.22 | 路线相关推荐 | GET | /api/v1/content/routes/{id}/related | - |
| 6.23 | 路线附近打卡地 | GET | /api/v1/content/routes/{id}/nearby | - |
| 6.24 | 打卡地附近推荐 | GET | /api/v1/content/places/{id}/nearby | - |
| 6.25 | 举报内容/用户 | POST | /api/v1/social/report | @auth |
| 6.26 | 获取未读消息数 | GET | /api/v1/social/unread-count | @auth |
| 7.1 | 攻略问答列表 | GET | /api/v1/content/guides/{id}/qa | - |
| 7.2 | 提问 | POST | /api/v1/content/guides/{id}/qa | @auth |
| 7.3 | 回答 | POST | /api/v1/content/guides/{id}/qa/{question_id}/answer | @auth |
| 7.4 | 点赞回答 | POST | /api/v1/content/guides/{id}/qa/{answer_id}/like | @auth |
| 8.1 | 城市列表 | GET | /api/v1/common/cities | - |
| 8.2 | 榜单 | GET | /api/v1/common/ranking | @auth（P1） |
| 8.3 | 足迹地图数据 | GET | /api/v1/common/footprint | @auth（P1） |
| 8.4 | AI 助手对话 | POST | /api/v1/ai/chat | @auth（P1） |
| 8.5 | 获取 AI 对话历史 | GET | /api/v1/ai/chat/history | @auth（P1） |
| 8.6 | 删除 AI 对话会话 | DELETE | /api/v1/ai/chat/session/{session_id} | @auth（P1） |
| 8.7 | 获取首页配置 | GET | /api/v1/common/home-config | - |
| 8.8 | 刷新 Token | POST | /api/v1/auth/refresh-token | @auth |
