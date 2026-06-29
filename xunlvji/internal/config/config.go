package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置
type Config struct {
	App        AppConfig        `yaml:"app"`
	MySQL      MySQLConfig      `yaml:"mysql"`
	Redis      RedisConfig      `yaml:"redis"`
	JWT        JWTConfig        `yaml:"jwt"`
	Volcengine VolcengineConfig `yaml:"volcengine"`
	Aliyun     AliyunConfig     `yaml:"aliyun"`
	AMap       AMapConfig       `yaml:"amap"`
	CPS        CPSConfig        `yaml:"cps"`
	SMS        SMSConfig        `yaml:"sms"`
	Push       PushConfig       `yaml:"push"`
	WeChat     WeChatConfig     `yaml:"wechat"`
	Log        LogConfig        `yaml:"log"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Mode    string `yaml:"mode"`
	Port    int    `yaml:"port"`
}

type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Charset         string `yaml:"charset"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

func (c MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	AccessTTL   int    `yaml:"access_ttl"`
	RefreshTTL  int    `yaml:"refresh_ttl"`
	Issuer      string `yaml:"issuer"`
}

type VolcengineConfig struct {
	AccessKey             string `yaml:"access_key"`
	SecretKey             string `yaml:"secret_key"`
	Region                string `yaml:"region"`
	DoubaoProEndpoint     string `yaml:"doubao_pro_endpoint"`
	DoubaoLiteEndpoint    string `yaml:"doubao_lite_endpoint"`
	ContentSafetyEndpoint string `yaml:"content_safety_endpoint"`
	TOSBucket             string `yaml:"tos_bucket"`
	TOSEndpoint           string `yaml:"tos_endpoint"`
	CDNDomain             string `yaml:"cdn_domain"`
}

type AliyunConfig struct {
	AccessKey             string `yaml:"access_key"`
	SecretKey             string `yaml:"secret_key"`
	QwenEndpoint          string `yaml:"qwen_endpoint"`
	ContentSafetyEndpoint string `yaml:"content_safety_endpoint"`
	OSSBucket             string `yaml:"oss_bucket"`
	OSSEndpoint           string `yaml:"oss_endpoint"`
	CDNDomain             string `yaml:"cdn_domain"`
}

type AMapConfig struct {
	Key string `yaml:"key"`
}

type CPSConfig struct {
	Eleme    CPSProvider `yaml:"eleme"`
	AmapRide CPSProvider `yaml:"amap_ride"`
	Ctrip    CPSProvider `yaml:"ctrip"`
	Meituan  CPSProvider `yaml:"meituan"`
}

type CPSProvider struct {
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
	Enabled   bool   `yaml:"enabled"`
}

type SMSConfig struct {
	Provider     string `yaml:"provider"`
	AccessKey    string `yaml:"access_key"`
	SecretKey    string `yaml:"secret_key"`
	SignName     string `yaml:"sign_name"`
	TemplateCode string `yaml:"template_code"`
}

type PushConfig struct {
	APNsKey       string `yaml:"apns_key"`
	APNsKeyID     string `yaml:"apns_key_id"`
	APNsTeamID    string `yaml:"apns_team_id"`
	FCMKey        string `yaml:"fcm_key"`
	JPushAppKey   string `yaml:"jpush_app_key"`
	JPushMasterSecret string `yaml:"jpush_master_secret"`
}

type WeChatConfig struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

type LogConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path"`
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	return &cfg, nil
}
