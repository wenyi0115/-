-- ============================================================
-- 文件名称: 001_init.sql
-- 项目名称: 寻旅记 (Xunlvji) 后端
-- 文件描述: 数据库初始化迁移脚本,包含全部 39 张表的建表语句及种子数据
-- 版本号: v1.0.0
-- 创建日期: 2026-06-28
-- 数据库: MySQL 8.0+
-- 字符集: utf8mb4 / utf8mb4_unicode_ci
-- 使用说明:
--   1. 确保已创建目标数据库(如 xunlvji),或取消下方 CREATE DATABASE 注释后执行
--   2. 执行命令: mysql -u root -p xunlvji < migrations/001_init.sql
--   3. 本脚本未做 DROP TABLE IF EXISTS 处理,重复执行前请先清理已有表
--   4. place 表含 SPATIAL INDEX(location POINT SRID 4326),需 MySQL 8.0+
-- 表数量: 39
-- ============================================================

-- 设置字符集与外键检查
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 数据库创建(默认注释,若数据库已预创建则无需打开)
-- CREATE DATABASE IF NOT EXISTS `xunlvji` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- USE `xunlvji`;


-- ============================================================
-- 用户相关表
-- ============================================================

-- ==================== 用户表 ====================
CREATE TABLE `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `nickname` VARCHAR(64) NOT NULL COMMENT '用户昵称',
    `avatar` VARCHAR(512) DEFAULT NULL COMMENT '头像URL',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    `email` VARCHAR(128) DEFAULT NULL COMMENT '邮箱',
    `password_hash` VARCHAR(256) DEFAULT NULL COMMENT 'bcrypt密码哈希(第三方登录可为空)',
    `role` VARCHAR(20) NOT NULL DEFAULT 'user' COMMENT '【已废弃·仅向后兼容】角色:user/creator/reviewer/operator/admin;新数据统一使用user_role关联表查询用户角色(详见3.26角色表/3.27用户角色关联表),新代码不应再读写此字段,保留仅为兼容历史数据与灰度迁移',
    `credit_score` INT NOT NULL DEFAULT 100 COMMENT '信用分(0-100)',
    `prefer_tags` JSON DEFAULT NULL COMMENT '偏好标签',
    `visited_city_count` INT NOT NULL DEFAULT 0 COMMENT '去过城市数',
    `follower_count` INT NOT NULL DEFAULT 0 COMMENT '粉丝数',
    `following_count` INT NOT NULL DEFAULT 0 COMMENT '关注数',
    `like_count` INT NOT NULL DEFAULT 0 COMMENT '获赞数',
    `bio` VARCHAR(200) DEFAULT NULL COMMENT '个人简介',
    `birthday` DATE DEFAULT NULL COMMENT '生日',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_requested_at` DATETIME DEFAULT NULL COMMENT '注销申请时间(30天冷静期)',
    `language` VARCHAR(16) NOT NULL DEFAULT 'zh-CN' COMMENT '语言偏好',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '0-禁用 1-正常 2-注销',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_phone` (`phone`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `idx_status` (`status`),
    KEY `idx_role` (`role`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `chk_user_credit` CHECK (`credit_score` >= 0 AND `credit_score` <= 100)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ========== 第三方账号关联表 ==========
CREATE TABLE `user_oauth` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `provider` VARCHAR(32) NOT NULL COMMENT '提供方:apple/google/wechat',
    `oauth_id` VARCHAR(128) NOT NULL COMMENT '第三方用户ID',
    `union_id` VARCHAR(128) DEFAULT NULL COMMENT '联合ID(微信unionId)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_provider_oauth` (`provider`, `oauth_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_union_id` (`union_id`),
    CONSTRAINT `fk_oauth_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='第三方账号关联表';

-- ========== 角色表 ==========
CREATE TABLE `role` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '角色名称',
    `code` VARCHAR(32) NOT NULL COMMENT '角色编码(唯一)',
    `description` VARCHAR(256) DEFAULT NULL COMMENT '角色描述',
    `permissions` JSON DEFAULT NULL COMMENT '权限列表(JSON数组)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ========== 用户角色关联表 ==========
CREATE TABLE `user_role` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
    KEY `idx_role_id` (`role_id`),
    CONSTRAINT `fk_ur_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_ur_role` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ========== 创作者认证申请表 ==========
CREATE TABLE `creator_application` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `real_name` VARCHAR(64) NOT NULL COMMENT '真实姓名',
    `id_card` VARCHAR(18) NOT NULL COMMENT '身份证号',
    `portfolio` VARCHAR(512) DEFAULT NULL COMMENT '作品集URL',
    `status` VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT '状态:pending/approved/rejected',
    `reviewer_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '审核人ID',
    `review_note` VARCHAR(256) DEFAULT NULL COMMENT '审核备注',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `reviewed_at` DATETIME DEFAULT NULL COMMENT '审核时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_ca_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_ca_reviewer` FOREIGN KEY (`reviewer_id`) REFERENCES `user` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='创作者认证申请表';

-- ========== 用户设置表 ==========
CREATE TABLE `user_setting` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `footprint_public` VARCHAR(16) NOT NULL DEFAULT 'public' COMMENT '足迹可见性:public(公开)/friends(仅好友)/private(仅自己)',
    `content_public` TINYINT NOT NULL DEFAULT 1 COMMENT '内容是否公开:0-否 1-是',
    `allow_dm_scope` VARCHAR(16) NOT NULL DEFAULT 'all' COMMENT '私信权限:all(所有人)/followers(仅关注者)/nobody(不允许)',
    `follow_scope` VARCHAR(16) NOT NULL DEFAULT 'all' COMMENT '关注权限:all(所有人)/verify(需验证)/nobody(不允许)',
    `location_public` TINYINT NOT NULL DEFAULT 0 COMMENT '位置是否公开:0-否 1-是(默认不公开,用户主动选择公开)',
    `history_public` TINYINT NOT NULL DEFAULT 0 COMMENT '历史行程是否公开:0-否 1-是',
    `show_favorites` TINYINT NOT NULL DEFAULT 0 COMMENT '收藏列表是否对外可见:0-不公开 1-公开',
    `show_following` TINYINT NOT NULL DEFAULT 0 COMMENT '关注/粉丝列表是否对外可见:0-不公开 1-公开',
    `profile_visibility` VARCHAR(16) NOT NULL DEFAULT 'public' COMMENT '个人主页可见性:public(公开)/friends(仅好友)/private(仅自己)',
    `personalized_recomm` TINYINT NOT NULL DEFAULT 1 COMMENT '个性化推荐开关:0-关闭 1-开启',
    `notification_push` TINYINT NOT NULL DEFAULT 1 COMMENT '推送开关:0-关闭 1-开启(总开关)',
    `notification_interactive` TINYINT NOT NULL DEFAULT 1 COMMENT '互动通知:0-关闭 1-开启(点赞/评论/关注/私信提醒)',
    `notification_trip` TINYINT NOT NULL DEFAULT 1 COMMENT '行程提醒:0-关闭 1-开启(行程开始/打卡提醒等)',
    `language` VARCHAR(10) NOT NULL DEFAULT 'zh-CN' COMMENT '语言:zh-CN/en-US等',
    `font_size` TINYINT NOT NULL DEFAULT 0 COMMENT '字体大小:0-小 1-中 2-大',
    `dark_mode` TINYINT NOT NULL DEFAULT 0 COMMENT '深色模式:0-跟随系统 1-关闭 2-开启',
    `daily_push_limit` INT NOT NULL DEFAULT 10 COMMENT '每日推送上限:单用户每日最多接收推送条数',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`),
    CONSTRAINT `fk_us_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户设置表';


-- ============================================================
-- 地点与打卡相关表
-- ============================================================

-- ==================== 地点表 ====================
CREATE TABLE `place` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(128) NOT NULL COMMENT '地点名称',
    `category` VARCHAR(32) NOT NULL COMMENT '分类(6大):food/fun/sightseeing/outdoor/shopping/accommodation',
    `location` POINT NOT NULL SRID 4326 COMMENT '地理位置(SRID 4326)',
    `longitude` DECIMAL(10,7) NOT NULL COMMENT '经度(冗余,便于读取)',
    `latitude` DECIMAL(10,7) NOT NULL COMMENT '纬度(冗余,便于读取)',
    `address` VARCHAR(256) DEFAULT NULL COMMENT '详细地址',
    `opening_hours` VARCHAR(128) DEFAULT NULL COMMENT '营业时间',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
    `avg_cost` DECIMAL(10,2) DEFAULT NULL COMMENT '人均消费(元)',
    `cover_image` JSON DEFAULT NULL COMMENT '封面图URL数组',
    `recommend_level` VARCHAR(16) DEFAULT NULL COMMENT '推荐等级:must_go/can_go/avoid',
    `suggest_duration` INT DEFAULT NULL COMMENT '建议时长(分钟)',
    `tags` JSON DEFAULT NULL COMMENT '标签',
    `city` VARCHAR(64) NOT NULL COMMENT '所在城市',
    `rating` DECIMAL(2,1) NOT NULL DEFAULT 0.0 COMMENT '综合评分(0.0-5.0)',
    `checkin_count` INT NOT NULL DEFAULT 0 COMMENT '打卡总数',
    `comment_count` INT NOT NULL DEFAULT 0 COMMENT '评论总数',
    `extra_fields` JSON DEFAULT NULL COMMENT '专属模板数据(吃喝:菜系/招牌菜;住宿:房型/设施等)',
    `cps_links` JSON DEFAULT NULL COMMENT 'CPS分销链接(ctrip_ticket/meituan_deal/eleme_shop/amap_taxi等)',
    `multi_lang` JSON DEFAULT NULL COMMENT '多语言字段,预留入境游({"en":{"name":"...","desc":"..."},"ja":{...}})',
    `foreigner_info` JSON DEFAULT NULL COMMENT '外国游客专属信息,预留(payment/visa_required/english_guide/signage_lang等)',
    `target_audience` VARCHAR(20) NOT NULL DEFAULT 'domestic' COMMENT '目标受众:domestic/inbound/outbound',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_category` (`category`),
    KEY `idx_city` (`city`),
    KEY `idx_city_category` (`city`, `category`),
    KEY `idx_recommend_level` (`recommend_level`),
    KEY `idx_target_audience` (`target_audience`),
    SPATIAL INDEX `sp_idx_location` (`location`),
    CONSTRAINT `chk_place_rating` CHECK (`rating` >= 0.0 AND `rating` <= 5.0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='地点表';

-- ==================== 打卡记录表(UGC内容表,打卡和发布共用) ====================
CREATE TABLE `checkin_record` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `place_id` BIGINT UNSIGNED NOT NULL COMMENT '地点ID',
    `images` JSON DEFAULT NULL COMMENT '图片列表',
    `rating` DECIMAL(2,1) DEFAULT NULL COMMENT '评分(0.0-5.0)',
    `description` TEXT DEFAULT NULL COMMENT '感受描述',
    `actual_cost` DECIMAL(8,2) DEFAULT NULL COMMENT '实际消费(元)',
    `pitfall_tips` JSON DEFAULT NULL COMMENT '避雷建议列表(JSON数组)',
    `checkin_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '打卡时间',
    `checkin_type` VARCHAR(16) NOT NULL DEFAULT 'normal' COMMENT 'store-到店打卡 normal-普通打卡',
    `visibility` VARCHAR(16) NOT NULL DEFAULT 'public' COMMENT 'public-公开(发布到打卡地页面) private-仅自己可见',
    `status` VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT 'pending-审核中(含PRD的pending_review/ai_reviewing/ai_approved/ai_rejected/human_review中间态) approved-已通过 rejected-已驳回',
    `reject_reason` VARCHAR(256) DEFAULT NULL COMMENT '审核驳回原因',
    `source` VARCHAR(20) NOT NULL DEFAULT 'trip_checkin' COMMENT '内容来源:trip_checkin-行程打卡 manual_publish-主动发布',
    `extra_data` JSON DEFAULT NULL COMMENT '打卡专属内容(吃喝:菜品评价;住宿:房型体验等),按place.category差异化',
    `content_richness` VARCHAR(10) NOT NULL DEFAULT 'light' COMMENT '内容丰富度:light/standard/rich,系统自动计算',
    `gps_verified` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'GPS验证标记:1-100米范围内打卡 0-未验证',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_place_id` (`place_id`),
    KEY `idx_user_place` (`user_id`, `place_id`),
    KEY `idx_checkin_time` (`checkin_time`),
    KEY `idx_status` (`status`),
    KEY `idx_visibility` (`visibility`),
    KEY `idx_source` (`source`),
    KEY `idx_content_richness` (`content_richness`),
    CONSTRAINT `chk_checkin_rating` CHECK (`rating` IS NULL OR (`rating` >= 0.0 AND `rating` <= 5.0)),
    CONSTRAINT `fk_checkin_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_checkin_place` FOREIGN KEY (`place_id`) REFERENCES `place` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='打卡记录表(UGC内容表)';

-- ==================== 避雷点赞表 ====================
CREATE TABLE `tip_vote` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `checkin_id` BIGINT UNSIGNED NOT NULL COMMENT '打卡记录ID',
    `tip_index` INT NOT NULL COMMENT '避雷建议索引(0-based)',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '点赞用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_checkin_tip_user` (`checkin_id`, `tip_index`, `user_id`),
    KEY `idx_checkin_tip` (`checkin_id`, `tip_index`),
    CONSTRAINT `fk_tipvote_checkin` FOREIGN KEY (`checkin_id`) REFERENCES `checkin_record` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_tipvote_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='避雷点赞表';


-- ============================================================
-- 路线相关表
-- ============================================================

-- ==================== 路线表 ====================
CREATE TABLE `route` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(128) NOT NULL COMMENT '路线标题',
    `cover_image` JSON DEFAULT NULL COMMENT '封面图URL数组',
    `total_duration` INT NOT NULL COMMENT '总时长(分钟)',
    `total_budget` DECIMAL(10,2) DEFAULT NULL COMMENT '总预算(元)',
    `suitable_for` VARCHAR(64) DEFAULT NULL COMMENT '适合人群',
    `category` VARCHAR(32) NOT NULL COMMENT '分类:half_day/one_day/two_day/three_day_plus',
    `tags` JSON DEFAULT NULL COMMENT '标签',
    `creator_id` BIGINT UNSIGNED NOT NULL COMMENT '创建者ID',
    `source` VARCHAR(16) NOT NULL DEFAULT 'user' COMMENT 'user/AI/official',
    `usage_count` INT NOT NULL DEFAULT 0 COMMENT '使用人数',
    `favorite_count` INT NOT NULL DEFAULT 0 COMMENT '收藏数',
    `comment_count` INT NOT NULL DEFAULT 0 COMMENT '评论总数',
    `rating` DECIMAL(2,1) NOT NULL DEFAULT 0.0 COMMENT '综合评分(0.0-5.0)',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '0-下架 1-正常 2-草稿',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_creator_id` (`creator_id`),
    KEY `idx_category` (`category`),
    KEY `idx_source` (`source`),
    KEY `idx_rating` (`rating`),
    KEY `idx_usage_count` (`usage_count`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_status` (`status`),
    CONSTRAINT `chk_route_rating` CHECK (`rating` >= 0.0 AND `rating` <= 5.0),
    CONSTRAINT `fk_route_creator` FOREIGN KEY (`creator_id`) REFERENCES `user` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路线表';

-- ==================== 路线点位表 ====================
CREATE TABLE `route_point` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `route_id` BIGINT UNSIGNED NOT NULL COMMENT '路线ID',
    `sort_order` INT NOT NULL DEFAULT 1 COMMENT '排序序号(从1开始)',
    `point_time` TIME DEFAULT NULL COMMENT '预计到达时间',
    `place_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '地点ID(可选)',
    `name` VARCHAR(128) NOT NULL COMMENT '点位名称',
    `stay_duration` INT DEFAULT NULL COMMENT '停留时长(分钟)',
    `cost` DECIMAL(10,2) DEFAULT NULL COMMENT '预计消费(元)',
    `category` VARCHAR(32) DEFAULT NULL COMMENT '分类',
    `transport` VARCHAR(32) DEFAULT NULL COMMENT '交通方式',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_route_id` (`route_id`),
    KEY `idx_route_sort` (`route_id`, `sort_order`),
    CONSTRAINT `fk_rp_route` FOREIGN KEY (`route_id`) REFERENCES `route` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rp_place` FOREIGN KEY (`place_id`) REFERENCES `place` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路线点位表';

-- ==================== 路线关联表 ====================
CREATE TABLE `route_related` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `route_id` BIGINT UNSIGNED NOT NULL COMMENT '路线ID',
    `related_route_id` BIGINT UNSIGNED NOT NULL COMMENT '关联路线ID',
    `sort_order` INT NOT NULL DEFAULT 1 COMMENT '排序序号',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_route_related` (`route_id`, `related_route_id`),
    KEY `idx_related_route_id` (`related_route_id`),
    CONSTRAINT `fk_rr_route` FOREIGN KEY (`route_id`) REFERENCES `route` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rr_related_route` FOREIGN KEY (`related_route_id`) REFERENCES `route` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路线关联表';


-- ============================================================
-- 攻略相关表
-- ============================================================

-- ==================== 攻略表 ====================
CREATE TABLE `guide` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(128) NOT NULL COMMENT '攻略标题',
    `cover_image` JSON DEFAULT NULL COMMENT '封面图URL数组',
    `category` VARCHAR(32) NOT NULL COMMENT '分类(7大):nature/history/entertainment/urban/transport/red_tourism/religion',
    `city` VARCHAR(64) NOT NULL COMMENT '所在城市',
    `duration` INT DEFAULT NULL COMMENT '游玩时长(分钟)',
    `price` DECIMAL(10,2) DEFAULT NULL COMMENT '价格(元)',
    `read_duration` INT DEFAULT NULL COMMENT '阅读时长(分钟)',
    `address` VARCHAR(256) DEFAULT NULL COMMENT '地址',
    `opening_hours` VARCHAR(128) DEFAULT NULL COMMENT '开放时间',
    `transport` VARCHAR(256) DEFAULT NULL COMMENT '交通方式',
    `content` MEDIUMTEXT NOT NULL COMMENT '正文内容(Markdown)',
    `highlights` TEXT DEFAULT NULL COMMENT '必看亮点',
    `pitfall_reminders` TEXT DEFAULT NULL COMMENT '避雷提醒',
    `verified_count` INT NOT NULL DEFAULT 0 COMMENT '验证人数',
    `comment_count` INT NOT NULL DEFAULT 0 COMMENT '评论总数',
    `source` VARCHAR(16) NOT NULL DEFAULT 'AI' COMMENT '来源标签(固定AI)',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '0-下架 1-正常 2-草稿',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_category` (`category`),
    KEY `idx_city` (`city`),
    KEY `idx_city_category` (`city`, `category`),
    KEY `idx_verified_count` (`verified_count`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='攻略表';

-- ==================== 攻略提问表 ====================
CREATE TABLE `guide_question` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '提问用户ID',
    `title` VARCHAR(128) NOT NULL COMMENT '问题标题',
    `content` TEXT COMMENT '问题详细描述',
    `city_id` INT UNSIGNED DEFAULT NULL COMMENT '关联城市ID',
    `status` VARCHAR(16) NOT NULL DEFAULT 'open' COMMENT 'open(待回答)/answered(已回答)/closed(已关闭)',
    `view_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '浏览数',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_gq_user` (`user_id`),
    KEY `idx_gq_city` (`city_id`),
    CONSTRAINT `fk_gq_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='攻略提问表';

-- ==================== 攻略回答表 ====================
CREATE TABLE `guide_answer` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `question_id` BIGINT UNSIGNED NOT NULL COMMENT '关联问题ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '回答用户ID',
    `content` TEXT NOT NULL COMMENT '回答内容',
    `like_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
    `is_accepted` TINYINT NOT NULL DEFAULT 0 COMMENT '是否被采纳:0-否 1-是',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_ga_question` (`question_id`),
    KEY `idx_ga_user` (`user_id`),
    CONSTRAINT `fk_ga_question` FOREIGN KEY (`question_id`) REFERENCES `guide_question` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_ga_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='攻略回答表';

-- ========== 攻略回答点赞表 ==========
CREATE TABLE `guide_answer_like` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `answer_id` BIGINT UNSIGNED NOT NULL COMMENT '回答ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '点赞用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_answer_user` (`answer_id`, `user_id`),
    KEY `idx_gal_user` (`user_id`),
    CONSTRAINT `fk_gal_answer` FOREIGN KEY (`answer_id`) REFERENCES `guide_answer` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_gal_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='攻略回答点赞表';


-- ============================================================
-- 行程相关表
-- ============================================================

-- ==================== 行程表 ====================
CREATE TABLE `trip` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `name` VARCHAR(128) NOT NULL COMMENT '行程名称',
    `trip_date` DATE NOT NULL COMMENT '出行日期(开始日期)',
    `end_date` DATE DEFAULT NULL COMMENT '结束日期(多日行程)',
    `budget` DECIMAL(10,2) DEFAULT NULL COMMENT '行程预算,单位:元',
    `visibility` VARCHAR(16) NOT NULL DEFAULT 'private' COMMENT '可见性:public(公开)/private(仅自己可见)/friends(仅好友可见)',
    `status` VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT 'pending-待出发 active-进行中 completed-已完成',
    `route_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联路线ID(可选)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_user_status` (`user_id`, `status`),
    KEY `idx_trip_date` (`trip_date`),
    KEY `idx_route_id` (`route_id`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_trip_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_trip_route` FOREIGN KEY (`route_id`) REFERENCES `route` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='行程表';

-- ==================== 行程点位表 ====================
CREATE TABLE `trip_point` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `trip_id` BIGINT UNSIGNED NOT NULL COMMENT '行程ID',
    `sort_order` INT NOT NULL DEFAULT 1 COMMENT '排序序号(从1开始)',
    `place_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '地点ID',
    `name` VARCHAR(128) NOT NULL COMMENT '点位名称',
    `point_time` TIME DEFAULT NULL COMMENT '安排时间',
    `stay_duration` INT DEFAULT NULL COMMENT '停留时长(分钟)',
    `cost` DECIMAL(10,2) DEFAULT NULL COMMENT '预计消费(元)',
    `category` VARCHAR(32) DEFAULT NULL COMMENT '分类',
    `transport` VARCHAR(32) DEFAULT NULL COMMENT '交通方式',
    `checkin_status` VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT 'pending-待打卡 checked-已打卡 skipped-已跳过',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_trip_id` (`trip_id`),
    KEY `idx_trip_sort` (`trip_id`, `sort_order`),
    KEY `idx_checkin_status` (`checkin_status`),
    CONSTRAINT `fk_tp_trip` FOREIGN KEY (`trip_id`) REFERENCES `trip` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_tp_place` FOREIGN KEY (`place_id`) REFERENCES `place` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='行程点位表';

-- ==================== 行程交通表 ====================
CREATE TABLE `route_transport` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `trip_id` BIGINT UNSIGNED NOT NULL COMMENT '行程ID',
    `from_place_id` BIGINT UNSIGNED NOT NULL COMMENT '起点地点ID',
    `to_place_id` BIGINT UNSIGNED NOT NULL COMMENT '终点地点ID',
    `transport_type` VARCHAR(20) NOT NULL COMMENT '交通方式:walk/bus/taxi/drive',
    `distance` INT DEFAULT NULL COMMENT '距离(米)',
    `duration` INT DEFAULT NULL COMMENT '预计时长(秒)',
    `cost` DECIMAL(8,2) DEFAULT NULL COMMENT '预计费用(元)',
    `cps_link` VARCHAR(500) DEFAULT NULL COMMENT '打车CPS链接(高德)',
    `sort_order` INT NOT NULL DEFAULT 1 COMMENT '排序序号(从1开始)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_trip_id` (`trip_id`),
    KEY `idx_trip_sort` (`trip_id`, `sort_order`),
    KEY `idx_from_place_id` (`from_place_id`),
    KEY `idx_to_place_id` (`to_place_id`),
    CONSTRAINT `fk_rt_trip` FOREIGN KEY (`trip_id`) REFERENCES `trip` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rt_from_place` FOREIGN KEY (`from_place_id`) REFERENCES `place` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_rt_to_place` FOREIGN KEY (`to_place_id`) REFERENCES `place` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='行程交通表';

-- ==================== 行程准备清单表 ====================
CREATE TABLE `trip_checklist` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `trip_id` BIGINT UNSIGNED NOT NULL COMMENT '行程ID',
    `item_name` VARCHAR(100) NOT NULL COMMENT '清单项名称',
    `item_category` VARCHAR(20) NOT NULL DEFAULT 'other' COMMENT '分类:document/electronic/clothing/other',
    `is_checked` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已勾选:0-否 1-是',
    `is_ai_generated` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否AI生成:0-否 1-是',
    `sort_order` INT NOT NULL DEFAULT 1 COMMENT '排序序号(从1开始)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_trip_id` (`trip_id`),
    KEY `idx_trip_sort` (`trip_id`, `sort_order`),
    CONSTRAINT `fk_tc_trip` FOREIGN KEY (`trip_id`) REFERENCES `trip` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='行程准备清单表';

-- ==================== 行程账单表（P1：P0阶段不创建，P1开发时建表） ====================
CREATE TABLE `trip_expense` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `trip_id` BIGINT UNSIGNED NOT NULL COMMENT '行程ID',
    `place_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联地点ID(可空,手动添加不关联)',
    `checkin_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联打卡记录ID(可空,从打卡自动来)',
    `expense_type` VARCHAR(20) NOT NULL COMMENT '消费类型:food/transport/ticket/hotel/other',
    `amount` DECIMAL(8,2) NOT NULL COMMENT '金额(元)',
    `description` VARCHAR(200) DEFAULT NULL COMMENT '描述',
    `is_auto` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否自动生成:0-否 1-是(从打卡记录来)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_trip_id` (`trip_id`),
    KEY `idx_trip_type` (`trip_id`, `expense_type`),
    KEY `idx_place_id` (`place_id`),
    KEY `idx_checkin_id` (`checkin_id`),
    CONSTRAINT `fk_te_trip` FOREIGN KEY (`trip_id`) REFERENCES `trip` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_te_place` FOREIGN KEY (`place_id`) REFERENCES `place` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_te_checkin` FOREIGN KEY (`checkin_id`) REFERENCES `checkin_record` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='行程账单表';


-- ============================================================
-- 社交互动相关表
-- ============================================================

-- ==================== 评论表 ====================
CREATE TABLE `comment` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `target_type` VARCHAR(16) NOT NULL COMMENT '目标类型:route/place/guide',
    `target_id` BIGINT UNSIGNED NOT NULL COMMENT '目标ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `parent_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '父评论ID(二级回复)',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `like_count` INT NOT NULL DEFAULT 0 COMMENT '点赞数',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_target` (`target_type`, `target_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_comment_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_comment_parent` FOREIGN KEY (`parent_id`) REFERENCES `comment` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- ==================== 评论点赞表 ====================
CREATE TABLE `comment_like` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `comment_id` BIGINT UNSIGNED NOT NULL COMMENT '评论ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '点赞用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_comment_user` (`comment_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_cl_comment` FOREIGN KEY (`comment_id`) REFERENCES `comment` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_cl_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论点赞表';

-- ==================== 关注关系表 ====================
CREATE TABLE `follow` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `follower_id` BIGINT UNSIGNED NOT NULL COMMENT '关注者ID',
    `following_id` BIGINT UNSIGNED NOT NULL COMMENT '被关注者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '关注时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_follower_following` (`follower_id`, `following_id`),
    KEY `idx_following_id` (`following_id`),
    KEY `idx_follower_id` (`follower_id`),
    CONSTRAINT `fk_follow_follower` FOREIGN KEY (`follower_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_follow_following` FOREIGN KEY (`following_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='关注关系表';

-- ==================== 收藏表 ====================
CREATE TABLE `favorite` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `target_type` VARCHAR(16) NOT NULL COMMENT '目标类型:route/place/guide/trip',
    `target_id` BIGINT UNSIGNED NOT NULL COMMENT '目标ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_target` (`user_id`, `target_type`, `target_id`),
    KEY `idx_target` (`target_type`, `target_id`),
    CONSTRAINT `fk_favorite_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏表';

-- ==================== 消息表 ====================
CREATE TABLE `message` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `type` VARCHAR(16) NOT NULL COMMENT '消息类型:system/private/group',
    `sender_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '发送者ID',
    `receiver_id` BIGINT UNSIGNED NOT NULL COMMENT '接收者ID/群ID(应用层校验)',
    `content` TEXT NOT NULL COMMENT '消息内容',
    `is_read` TINYINT NOT NULL DEFAULT 0 COMMENT '0-未读 1-已读',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_receiver_read` (`receiver_id`, `is_read`),
    KEY `idx_sender_receiver` (`sender_id`, `receiver_id`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_type` (`type`),
    CONSTRAINT `fk_msg_sender` FOREIGN KEY (`sender_id`) REFERENCES `user` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息表';


-- ============================================================
-- 群聊相关表
-- ============================================================

-- ========== 群聊表 ==========
CREATE TABLE `chat_group` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '群名称',
    `avatar` VARCHAR(512) DEFAULT NULL COMMENT '群头像',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '群主ID',
    `member_count` INT NOT NULL DEFAULT 1 COMMENT '成员数',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    CONSTRAINT `fk_group_owner` FOREIGN KEY (`owner_id`) REFERENCES `user` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群聊表';

-- ========== 群成员表 ==========
CREATE TABLE `group_member` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `group_id` BIGINT UNSIGNED NOT NULL COMMENT '群ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `role` VARCHAR(16) NOT NULL DEFAULT 'member' COMMENT '角色:admin/member',
    `joined_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_user` (`group_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_gm_group` FOREIGN KEY (`group_id`) REFERENCES `chat_group` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_gm_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群成员表';


-- ============================================================
-- 勋章相关表
-- ============================================================

-- ========== 勋章定义表 ==========
CREATE TABLE `badge_definition` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(64) NOT NULL COMMENT '勋章编码(唯一)',
    `name` VARCHAR(64) NOT NULL COMMENT '勋章名称',
    `description` VARCHAR(256) DEFAULT NULL COMMENT '勋章描述',
    `icon` VARCHAR(512) NOT NULL COMMENT '勋章图标URL',
    `category` VARCHAR(32) NOT NULL COMMENT '分类:checkin/social/creator/system',
    `condition_type` VARCHAR(32) NOT NULL COMMENT '条件类型:city_count/checkin_count/follower_count等',
    `condition_value` VARCHAR(64) NOT NULL COMMENT '条件值',
    `sort_order` INT NOT NULL DEFAULT 1 COMMENT '排序序号',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '0-禁用 1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_category` (`category`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='勋章定义表';

-- ========== 用户勋章表 ==========
CREATE TABLE `user_badge` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `badge_id` BIGINT UNSIGNED NOT NULL COMMENT '勋章ID',
    `unlocked_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '解锁时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_badge` (`user_id`, `badge_id`),
    KEY `idx_badge_id` (`badge_id`),
    CONSTRAINT `fk_ub_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_ub_badge` FOREIGN KEY (`badge_id`) REFERENCES `badge_definition` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户勋章表';

-- ========== 勋章汇总表 ==========
CREATE TABLE `user_medal_summary` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `total_count` INT NOT NULL DEFAULT 0 COMMENT '勋章总数',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`),
    CONSTRAINT `fk_ums_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='勋章汇总表';


-- ============================================================
-- AI 对话相关表
-- ============================================================

-- ========== AI对话会话表（P1：P0阶段不创建，P1开发时建表） ==========
CREATE TABLE `ai_chat_session` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `title` VARCHAR(128) NOT NULL COMMENT '会话标题',
    `context` JSON DEFAULT NULL COMMENT '上下文信息',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_acs_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI对话会话表';

-- ========== AI对话消息表（P1：P0阶段不创建，P1开发时建表） ==========
CREATE TABLE `ai_chat_message` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `session_id` BIGINT UNSIGNED NOT NULL COMMENT '会话ID',
    `role` VARCHAR(16) NOT NULL COMMENT '角色:user/assistant',
    `content` TEXT NOT NULL COMMENT '消息内容',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_session_id` (`session_id`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_acm_session` FOREIGN KEY (`session_id`) REFERENCES `ai_chat_session` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI对话消息表';


-- ============================================================
-- 城市与排名相关表
-- ============================================================

-- ========== 城市信息表 ==========
CREATE TABLE `city` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL COMMENT '城市名称',
    `province_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '省份ID(预留)',
    `province` VARCHAR(50) DEFAULT NULL COMMENT '省份名称(冗余,便于展示)',
    `latitude` DECIMAL(10,7) DEFAULT NULL COMMENT '纬度',
    `longitude` DECIMAL(10,7) DEFAULT NULL COMMENT '经度',
    `best_season` VARCHAR(100) DEFAULT NULL COMMENT '最佳旅行季节',
    `description` TEXT DEFAULT NULL COMMENT '城市简介',
    `cover_image` VARCHAR(500) DEFAULT NULL COMMENT '封面图URL',
    `is_hot` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否热门城市:0-否 1-是',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序序号',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '0-禁用 1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`),
    KEY `idx_province_id` (`province_id`),
    KEY `idx_is_hot` (`is_hot`),
    KEY `idx_sort_order` (`sort_order`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='城市信息表';

-- ========== 用户城市进度表 ==========
CREATE TABLE `user_city_progress` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `city_id` BIGINT UNSIGNED NOT NULL COMMENT '城市ID',
    `total_places` INT NOT NULL DEFAULT 0 COMMENT '总地点数',
    `checked_places` INT NOT NULL DEFAULT 0 COMMENT '已打卡地点数',
    `progress` DECIMAL(5,2) NOT NULL DEFAULT 0.00 COMMENT '进度百分比(0-100)',
    `unlocked_at` DATETIME DEFAULT NULL COMMENT '解锁时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_city` (`user_id`, `city_id`),
    KEY `idx_city_id` (`city_id`),
    CONSTRAINT `fk_ucp_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_ucp_city` FOREIGN KEY (`city_id`) REFERENCES `city` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户城市进度表';

-- ========== 用户省份进度表 ==========
CREATE TABLE `user_province_progress` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `province_id` BIGINT UNSIGNED NOT NULL COMMENT '省份ID',
    `unlocked` TINYINT NOT NULL DEFAULT 0 COMMENT '是否解锁:0-未解锁 1-已解锁',
    `unlocked_at` DATETIME DEFAULT NULL COMMENT '解锁时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_province` (`user_id`, `province_id`),
    KEY `idx_province_id` (`province_id`),
    CONSTRAINT `fk_upp_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户省份进度表';

-- ========== 全国排名表 ==========
CREATE TABLE `user_national_ranking` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `city_count` INT NOT NULL DEFAULT 0 COMMENT '打卡城市数',
    `checkin_count` INT NOT NULL DEFAULT 0 COMMENT '打卡总数',
    `total_score` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '总积分',
    `rank` INT NOT NULL DEFAULT 0 COMMENT '排名',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`),
    KEY `idx_rank` (`rank`),
    KEY `idx_total_score` (`total_score`),
    CONSTRAINT `fk_unr_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='全国排名表';


-- ============================================================
-- 系统相关表
-- ============================================================

-- ========== 审计日志表 ==========
CREATE TABLE `audit_log` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '操作用户ID',
    `action` VARCHAR(64) NOT NULL COMMENT '操作类型',
    `target_type` VARCHAR(32) NOT NULL COMMENT '目标类型',
    `target_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '目标ID',
    `detail` JSON DEFAULT NULL COMMENT '操作详情',
    `ip` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_action` (`action`),
    KEY `idx_target` (`target_type`, `target_id`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';

-- ==================== 举报表 ====================
CREATE TABLE `report` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `reporter_id` BIGINT UNSIGNED NOT NULL COMMENT '举报人ID',
    `target_type` VARCHAR(20) NOT NULL COMMENT '举报对象类型:checkin/comment/place/user',
    `target_id` BIGINT UNSIGNED NOT NULL COMMENT '举报对象ID',
    `reason` VARCHAR(50) NOT NULL COMMENT '举报原因:spam/inappropriate/wrong_info/copyright/other',
    `description` TEXT DEFAULT NULL COMMENT '详细描述',
    `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态:pending-待处理 resolved-已处理 rejected-已驳回',
    `handled_by` BIGINT UNSIGNED DEFAULT NULL COMMENT '处理人ID(管理员)',
    `handled_at` DATETIME DEFAULT NULL COMMENT '处理时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_reporter_id` (`reporter_id`),
    KEY `idx_target` (`target_type`, `target_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_report_reporter` FOREIGN KEY (`reporter_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_report_handler` FOREIGN KEY (`handled_by`) REFERENCES `user` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='举报表';


-- ============================================================
-- 种子数据 (Seed Data)
-- ============================================================

-- ---------- 默认角色 ----------
INSERT INTO `role` (`name`, `code`, `description`, `permissions`) VALUES
('普通用户', 'user', '普通注册用户', '["content:view","content:publish","trip:manage","social:interact"]'),
('创作者', 'creator', '认证创作者', '["content:view","content:publish","trip:manage","social:interact","content:priority"]'),
('审核员', 'reviewer', '内容审核员', '["content:view","content:review","report:handle"]'),
('运营人员', 'operator', '平台运营', '["content:view","content:manage","city:manage","badge:manage"]'),
('管理员', 'admin', '系统管理员', '["*"]');

-- ---------- 勋章定义(对应原型勋章页,共 8 枚) ----------
INSERT INTO `badge_definition` (`code`, `name`, `description`, `icon`, `category`, `condition_type`, `condition_value`, `sort_order`, `status`) VALUES
('first_checkin', '初出茅庐', '完成首次打卡', '', 'checkin', 'checkin_count', '1', 1, 1),
('city_explorer', '城市探索者', '打卡5个不同城市', '', 'checkin', 'city_count', '5', 2, 1),
('food_master', '美食达人', '完成10次吃喝类打卡', '', 'checkin', 'checkin_count_food', '10', 3, 1),
('mountain_climber', '登山爱好者', '完成5次户外类打卡', '', 'checkin', 'checkin_count_outdoor', '5', 4, 1),
('social_star', '社交之星', '获得100个关注', '', 'social', 'follower_count', '100', 5, 1),
('route_creator', '路线规划师', '创建3条路线', '', 'creator', 'route_count', '3', 6, 1),
('checkin_master', '打卡达人', '完成50次打卡', '', 'checkin', 'checkin_count', '50', 7, 1),
('travel_expert', '旅行家', '打卡20个不同城市', '', 'checkin', 'city_count', '20', 8, 1);


-- ============================================================
-- 结束:恢复外键检查
-- ============================================================
SET FOREIGN_KEY_CHECKS = 1;
