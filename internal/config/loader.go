package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"greatestworks/application/services"

	"gopkg.in/yaml.v3"
)

// AppConfig 应用程序完整配置
type AppConfig struct {
	Server      services.ServerConfig `yaml:"server"`
	MongoDB     MongoDBConfig         `yaml:"mongodb"`
	NATS        NATSConfig            `yaml:"nats"`
	Redis       RedisConfig           `yaml:"redis"`
	Game        GameConfig            `yaml:"game"`
	Logging     LoggingConfig         `yaml:"logging"`
	Monitoring  MonitoringConfig      `yaml:"monitoring"`
	Security    SecurityConfig        `yaml:"security"`
	ThirdParty  ThirdPartyConfig      `yaml:"third_party"`
	Development EnvironmentConfig     `yaml:"development"`
	Production  EnvironmentConfig     `yaml:"production"`
	Performance PerformanceConfig     `yaml:"performance"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
	Timeout  int    `yaml:"timeout"`
}

// 使用config.go中定义的类型

// GameConfig 游戏配置
type GameConfig struct {
	Player PlayerConfig `yaml:"player"`
	Bag    BagConfig    `yaml:"bag"`
	Pet    PetConfig    `yaml:"pet"`
	VIP    VIPConfig    `yaml:"vip"`
	Chat   ChatConfig   `yaml:"chat"`
	Shop   ShopConfig   `yaml:"shop"`
}

// PlayerConfig 玩家配置
type PlayerConfig struct {
	MaxLevel        int    `yaml:"max_level"`
	InitialGold     uint32 `yaml:"initial_gold"`
	InitialDiamonds uint32 `yaml:"initial_diamonds"`
	MaxFriends      int    `yaml:"max_friends"`
}

// BagConfig 背包配置
type BagConfig struct {
	DefaultCapacity int    `yaml:"default_capacity"`
	MaxCapacity     int    `yaml:"max_capacity"`
	ExpandCostBase  uint32 `yaml:"expand_cost_base"`
}

// PetConfig 宠物配置
type PetConfig struct {
	MaxPetsPerPlayer          int `yaml:"max_pets_per_player"`
	MaxActivePets             int `yaml:"max_active_pets"`
	EvolutionLevelRequirement int `yaml:"evolution_level_requirement"`
}

// VIPConfig VIP配置
type VIPConfig struct {
	MaxLevel      int     `yaml:"max_level"`
	ExpMultiplier float64 `yaml:"exp_multiplier"`
}

// ChatConfig 聊天配置
type ChatConfig struct {
	MessageHistoryDays int  `yaml:"message_history_days"`
	MaxMessageLength   int  `yaml:"max_message_length"`
	SpamProtection     bool `yaml:"spam_protection"`
}

// ShopConfig 商店配置
type ShopConfig struct {
	RefreshInterval   time.Duration `yaml:"refresh_interval"`
	MaxPurchasePerDay int           `yaml:"max_purchase_per_day"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled             bool          `yaml:"enabled"`
	MetricsPort         int           `yaml:"metrics_port"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	JWTSecret  string          `yaml:"jwt_secret"`
	JWTExpiry  time.Duration   `yaml:"jwt_expiry"`
	BcryptCost int             `yaml:"bcrypt_cost"`
	RateLimit  RateLimitConfig `yaml:"rate_limit"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerMinute int `yaml:"requests_per_minute"`
	Burst             int `yaml:"burst"`
}

// ThirdPartyConfig 第三方服务配置
type ThirdPartyConfig struct {
	Payment          PaymentConfig          `yaml:"payment"`
	PushNotification PushNotificationConfig `yaml:"push_notification"`
	Email            EmailConfig            `yaml:"email"`
}

// PaymentConfig 支付配置
type PaymentConfig struct {
	Stripe StripeConfig `yaml:"stripe"`
}

// StripeConfig Stripe配置
type StripeConfig struct {
	PublicKey     string `yaml:"public_key"`
	SecretKey     string `yaml:"secret_key"`
	WebhookSecret string `yaml:"webhook_secret"`
}

// PushNotificationConfig 推送通知配置
type PushNotificationConfig struct {
	Firebase FirebaseConfig `yaml:"firebase"`
}

// FirebaseConfig Firebase配置
type FirebaseConfig struct {
	ServerKey string `yaml:"server_key"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTP SMTPConfig `yaml:"smtp"`
}

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// EnvironmentConfig 环境配置
type EnvironmentConfig struct {
	Debug     bool `yaml:"debug"`
	HotReload bool `yaml:"hot_reload"`
	MockData  bool `yaml:"mock_data"`
	TestMode  bool `yaml:"test_mode"`
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	DBPool      DBPoolConfig      `yaml:"db_pool"`
	Cache       CacheConfig       `yaml:"cache"`
	Concurrency ConcurrencyConfig `yaml:"concurrency"`
}

// DBPoolConfig 数据库连接池配置
type DBPoolConfig struct {
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	DefaultTTL      time.Duration `yaml:"default_ttl"`
	CleanupInterval time.Duration `yaml:"cleanup_interval"`
	MaxSize         int           `yaml:"max_size"`
}

// ConcurrencyConfig 并发配置
type ConcurrencyConfig struct {
	MaxGoroutines   int `yaml:"max_goroutines"`
	WorkerPoolSize  int `yaml:"worker_pool_size"`
	QueueBufferSize int `yaml:"queue_buffer_size"`
}

// ConfigLoader 配置加载器
type ConfigLoader struct {
	configPath string
	config     *AppConfig
}

// NewConfigLoader 创建配置加载器
func NewConfigLoader(configPath string) *ConfigLoader {
	return &ConfigLoader{
		configPath: configPath,
	}
}

// Load 加载配置文件
func (cl *ConfigLoader) Load() (*AppConfig, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(cl.configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", cl.configPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(cl.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析YAML
	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 应用环境变量覆盖
	if err := cl.applyEnvironmentOverrides(&config); err != nil {
		return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
	}

	// 验证配置
	if err := cl.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	cl.config = &config
	return &config, nil
}

// GetConfig 获取当前配置
func (cl *ConfigLoader) GetConfig() *AppConfig {
	return cl.config
}

// Reload 重新加载配置
func (cl *ConfigLoader) Reload() error {
	_, err := cl.Load()
	return err
}

// applyEnvironmentOverrides 应用环境变量覆盖
func (cl *ConfigLoader) applyEnvironmentOverrides(config *AppConfig) error {
	// MongoDB配置覆盖
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		config.MongoDB.URI = uri
	}
	if db := os.Getenv("MONGODB_DATABASE"); db != "" {
		config.MongoDB.Database = db
	}

	// NATS配置覆盖
	if url := os.Getenv("NATS_URL"); url != "" {
		config.NATS.URL = url
	}
	if clusterID := os.Getenv("NATS_CLUSTER_ID"); clusterID != "" {
		config.NATS.ClusterID = clusterID
	}
	if clientID := os.Getenv("NATS_CLIENT_ID"); clientID != "" {
		config.NATS.ClientID = clientID
	}

	// Redis配置覆盖
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		config.Redis.Addr = addr
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}

	// 服务器配置覆盖
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		// 这里可以添加端口解析逻辑
	}

	// 安全配置覆盖
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.Security.JWTSecret = jwtSecret
	}

	// 第三方服务配置覆盖
	if stripeSecret := os.Getenv("STRIPE_SECRET_KEY"); stripeSecret != "" {
		config.ThirdParty.Payment.Stripe.SecretKey = stripeSecret
	}
	if stripePublic := os.Getenv("STRIPE_PUBLIC_KEY"); stripePublic != "" {
		config.ThirdParty.Payment.Stripe.PublicKey = stripePublic
	}

	return nil
}

// validateConfig 验证配置
func (cl *ConfigLoader) validateConfig(config *AppConfig) error {
	// 验证必需的配置项
	if config.MongoDB.URI == "" {
		return fmt.Errorf("MongoDB URI is required")
	}
	if config.MongoDB.Database == "" {
		return fmt.Errorf("MongoDB database name is required")
	}
	if config.NATS.URL == "" {
		return fmt.Errorf("NATS URL is required")
	}
	if config.Security.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	// 验证端口范围
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}
	if config.Monitoring.MetricsPort < 1 || config.Monitoring.MetricsPort > 65535 {
		return fmt.Errorf("invalid metrics port: %d", config.Monitoring.MetricsPort)
	}

	// 验证游戏配置
	if config.Game.Player.MaxLevel < 1 {
		return fmt.Errorf("player max level must be greater than 0")
	}
	if config.Game.Bag.DefaultCapacity < 1 {
		return fmt.Errorf("bag default capacity must be greater than 0")
	}
	if config.Game.Pet.MaxPetsPerPlayer < 1 {
		return fmt.Errorf("max pets per player must be greater than 0")
	}

	return nil
}

// ToServiceConfig 转换为服务配置
func (config *AppConfig) ToServiceConfig() services.Config {
	return services.Config{
		MongoDB: database.MongoConfig{
			URI:      config.MongoDB.URI,
			Database: config.MongoDB.Database,
			Timeout:  config.MongoDB.Timeout,
		},
		NATS: messaging.NATSConfig{
			URL:            config.NATS.URL,
			ClusterID:      config.NATS.ClusterID,
			ClientID:       config.NATS.ClientID,
			ReconnectWait:  config.NATS.ReconnectWait,
			MaxReconnects:  config.NATS.MaxReconnects,
			ConnectionName: config.NATS.ConnectionName,
			DrainTimeout:   config.NATS.DrainTimeout,
		},
		Server: config.Server,
	}
}

// GetEnvironment 获取当前环境
func GetEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	return strings.ToLower(env)
}

// IsProduction 是否为生产环境
func IsProduction() bool {
	return GetEnvironment() == "production"
}

// IsDevelopment 是否为开发环境
func IsDevelopment() bool {
	return GetEnvironment() == "development"
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	// 默认配置文件路径
	env := GetEnvironment()
	configFile := fmt.Sprintf("config.%s.yaml", env)

	// 查找配置文件
	searchPaths := []string{
		".",
		"./config",
		"../config",
		"/etc/mmo-game",
	}

	for _, searchPath := range searchPaths {
		fullPath := filepath.Join(searchPath, configFile)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	// 如果找不到环境特定的配置文件，使用默认配置文件
	for _, searchPath := range searchPaths {
		fullPath := filepath.Join(searchPath, "config.yaml")
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	return "config/config.yaml"
}
