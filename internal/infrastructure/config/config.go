// Package config 统一配置管理
// Author: MMO Server Team
// Created: 2024

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 主配置结构
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Cache      CacheConfig      `yaml:"cache"`
	Security   SecurityConfig   `yaml:"security"`
	Logging    LoggingConfig    `yaml:"logging"`
	Game       GameConfig       `yaml:"game"`
	Network    NetworkConfig    `yaml:"network"`
	Messaging  MessagingConfig  `yaml:"messaging"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	Excel      ExcelConfig      `yaml:"excel"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Gateway  GatewayConfig  `yaml:"gateway"`
	Scene    SceneConfig    `yaml:"scene"`
	Battle   BattleConfig   `yaml:"battle"`
	Activity ActivityConfig `yaml:"activity"`
}

// GatewayConfig 网关服务器配置
type GatewayConfig struct {
	Port            int           `yaml:"port"`
	Host            string        `yaml:"host"`
	MaxConnections  int           `yaml:"max_connections"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	HeartbeatTime   time.Duration `yaml:"heartbeat_time"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	TLSEnabled      bool          `yaml:"tls_enabled"`
	CertFile        string        `yaml:"cert_file"`
	KeyFile         string        `yaml:"key_file"`
}

// SceneConfig 场景服务器配置
type SceneConfig struct {
	Port         int           `yaml:"port"`
	Host         string        `yaml:"host"`
	MaxPlayers   int           `yaml:"max_players"`
	TickRate     int           `yaml:"tick_rate"`
	SyncInterval time.Duration `yaml:"sync_interval"`
	ViewDistance float64       `yaml:"view_distance"`
}

// BattleConfig 战斗服务器配置
type BattleConfig struct {
	Port         int           `yaml:"port"`
	Host         string        `yaml:"host"`
	MaxBattles   int           `yaml:"max_battles"`
	TickRate     int           `yaml:"tick_rate"`
	BattleTime   time.Duration `yaml:"battle_time"`
	MatchTimeout time.Duration `yaml:"match_timeout"`
}

// ActivityConfig 活动服务器配置
type ActivityConfig struct {
	Port           int           `yaml:"port"`
	Host           string        `yaml:"host"`
	MaxActivities  int           `yaml:"max_activities"`
	UpdateInterval time.Duration `yaml:"update_interval"`
	CacheTimeout   time.Duration `yaml:"cache_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MongoDB MongoDBConfig `yaml:"mongodb"`
	Redis   RedisConfig   `yaml:"redis"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI            string        `yaml:"uri"`
	Database       string        `yaml:"database"`
	MaxPoolSize    int           `yaml:"max_pool_size"`
	MinPoolSize    int           `yaml:"min_pool_size"`
	MaxIdleTime    time.Duration `yaml:"max_idle_time"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
	SocketTimeout  time.Duration `yaml:"socket_timeout"`
	RetryWrites    bool          `yaml:"retry_writes"`
	ReadPreference string        `yaml:"read_preference"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string        `yaml:"addr"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	MaxIdleConns int           `yaml:"max_idle_conns"`
	ConnMaxAge   time.Duration `yaml:"conn_max_age"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	Cluster      ClusterConfig `yaml:"cluster"`
}

// ClusterConfig Redis集群配置
type ClusterConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Addresses []string `yaml:"addresses"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	DefaultTTL    time.Duration `yaml:"default_ttl"`
	CleanupTime   time.Duration `yaml:"cleanup_time"`
	MaxSize       int64         `yaml:"max_size"`
	EvictionPolicy string       `yaml:"eviction_policy"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	JWT        JWTConfig        `yaml:"jwt"`
	Encryption EncryptionConfig `yaml:"encryption"`
	RateLimit  RateLimitConfig  `yaml:"rate_limit"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey     string        `yaml:"secret_key"`
	TokenDuration time.Duration `yaml:"token_duration"`
	RefreshTime   time.Duration `yaml:"refresh_time"`
	Issuer        string        `yaml:"issuer"`
	Audience      string        `yaml:"audience"`
}

// EncryptionConfig 加密配置
type EncryptionConfig struct {
	Algorithm string `yaml:"algorithm"`
	KeySize   int    `yaml:"key_size"`
	Salt      string `yaml:"salt"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled    bool          `yaml:"enabled"`
	RPS        int           `yaml:"rps"`
	Burst      int           `yaml:"burst"`
	WindowSize time.Duration `yaml:"window_size"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	Dir        string `yaml:"dir"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
	Prefix     string `yaml:"prefix"`
}

// GameConfig 游戏配置
type GameConfig struct {
	MaxLevel        int     `yaml:"max_level"`
	ExpMultiplier   float64 `yaml:"exp_multiplier"`
	GoldMultiplier  float64 `yaml:"gold_multiplier"`
	DropRate        float64 `yaml:"drop_rate"`
	PKEnabled       bool    `yaml:"pk_enabled"`
	GuildEnabled    bool    `yaml:"guild_enabled"`
	TradeEnabled    bool    `yaml:"trade_enabled"`
	MaintenanceMode bool    `yaml:"maintenance_mode"`
}

// NetworkConfig 网络配置
type NetworkConfig struct {
	Protocol        string `yaml:"protocol"`
	BufferSize      int    `yaml:"buffer_size"`
	MaxPacketSize   int    `yaml:"max_packet_size"`
	CompressionType string `yaml:"compression_type"`
	EncryptionType  string `yaml:"encryption_type"`
	KeepAlive       bool   `yaml:"keep_alive"`
	NoDelay         bool   `yaml:"no_delay"`
}

// MessagingConfig 消息队列配置
type MessagingConfig struct {
	NSQ      NSQConfig      `yaml:"nsq"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

// NSQConfig NSQ配置
type NSQConfig struct {
	Enabled     bool     `yaml:"enabled"`
	NSQDAddress string   `yaml:"nsqd_address"`
	LookupdHTTP []string `yaml:"lookupd_http"`
	MaxInFlight int      `yaml:"max_in_flight"`
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Enabled  bool   `yaml:"enabled"`
	URL      string `yaml:"url"`
	Exchange string `yaml:"exchange"`
	Queue    string `yaml:"queue"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled    bool           `yaml:"enabled"`
	Port       int            `yaml:"port"`
	Path       string         `yaml:"path"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
}

// PrometheusConfig Prometheus配置
type PrometheusConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Namespace string `yaml:"namespace"`
	Subsystem string `yaml:"subsystem"`
}

// ExcelConfig Excel配置
type ExcelConfig struct {
	Path       string `yaml:"path"`
	Activity   string `yaml:"activity"`
	BattlePass string `yaml:"battle_pass"`
	Pet        string `yaml:"pet"`
	Npc        string `yaml:"npc"`
	Plant      string `yaml:"plant"`
	Shop       string `yaml:"shop"`
	Task       string `yaml:"task"`
	Skill      string `yaml:"skill"`
	Vip        string `yaml:"vip"`
	Building   string `yaml:"building"`
	Condition  string `yaml:"condition"`
	Synthetise string `yaml:"synthetise"`
	MiniGame   string `yaml:"mini_game"`
	Email      string `yaml:"email"`
}

// 全局配置实例
var (
	globalConfig *Config
	configMutex  sync.RWMutex
	configOnce   sync.Once
)

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	// 确定配置文件路径
	if configPath == "" {
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "dev"
		}
		configPath = fmt.Sprintf("config.%s.yaml", env)
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// 环境变量替换
	configContent := os.ExpandEnv(string(data))

	// 解析YAML
	var config Config
	if err := yaml.Unmarshal([]byte(configContent), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 设置默认值
	setDefaults(&config)

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// InitConfig 初始化全局配置
func InitConfig(configPath string) error {
	var err error
	configOnce.Do(func() {
		globalConfig, err = LoadConfig(configPath)
	})
	return err
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return globalConfig
}

// ReloadConfig 重新加载配置
func ReloadConfig(configPath string) error {
	newConfig, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	configMutex.Lock()
	globalConfig = newConfig
	configMutex.Unlock()

	return nil
}

// BuildRealPath 构建真实路径
func BuildRealPath(base string, module string) string {
	if filepath.IsAbs(module) {
		return module
	}
	return filepath.Join(base, module)
}

// GetExcelPath 获取Excel文件路径
func GetExcelPath(module string) string {
	config := GetConfig()
	if config == nil {
		return ""
	}

	var excelFile string
	switch strings.ToLower(module) {
	case "activity":
		excelFile = config.Excel.Activity
	case "battlepass":
		excelFile = config.Excel.BattlePass
	case "pet":
		excelFile = config.Excel.Pet
	case "npc":
		excelFile = config.Excel.Npc
	case "plant":
		excelFile = config.Excel.Plant
	case "shop":
		excelFile = config.Excel.Shop
	case "task":
		excelFile = config.Excel.Task
	case "skill":
		excelFile = config.Excel.Skill
	case "vip":
		excelFile = config.Excel.Vip
	case "building":
		excelFile = config.Excel.Building
	case "condition":
		excelFile = config.Excel.Condition
	case "synthetise":
		excelFile = config.Excel.Synthetise
	case "minigame":
		excelFile = config.Excel.MiniGame
	case "email":
		excelFile = config.Excel.Email
	default:
		return ""
	}

	return BuildRealPath(config.Excel.Path, excelFile)
}