// Package config 统一配置管理
// Author: MMO Server Team
// Created: 2024

package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 主配置结构
type Config struct {
	App        AppConfig        `yaml:"app"`
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Cache      CacheConfig      `yaml:"cache"`
	Security   SecurityConfig   `yaml:"security"`
	Logging    LoggingConfig    `yaml:"logging"`
	Game       GameConfig       `yaml:"game"`
	Network    NetworkConfig    `yaml:"network"`
	Messaging  MessagingConfig  `yaml:"messaging"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
}

// ConfigLoader 配置加载器
type ConfigLoader struct {
	configPath string
}

// NewConfigLoader 创建配置加载器
func NewConfigLoader(configPath string) *ConfigLoader {
	return &ConfigLoader{
		configPath: configPath,
	}
}

// Load 加载配置
func (cl *ConfigLoader) Load() (*Config, error) {
	return LoadConfig(cl.configPath)
}

// AppConfig 应用配置
type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
	Debug       bool   `yaml:"debug"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	HTTP      HTTPServerConfig      `yaml:"http"`
	WebSocket WebSocketServerConfig `yaml:"websocket"`
	TCP       TCPServerConfig       `yaml:"tcp"`
	Metrics   MetricsServerConfig   `yaml:"metrics"`
}

// HTTPServerConfig HTTP服务器配置
type HTTPServerConfig struct {
	Host              string        `yaml:"host"`
	Port              int           `yaml:"port"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	MaxHeaderBytes    int           `yaml:"max_header_bytes"`
	EnableCORS        bool          `yaml:"enable_cors"`
	EnableMetrics     bool          `yaml:"enable_metrics"`
	EnableRequestID   bool          `yaml:"enable_request_id"`
	EnableLogging     bool          `yaml:"enable_logging"`
	EnableRecovery    bool          `yaml:"enable_recovery"`
	RateLimitEnabled  bool          `yaml:"rate_limit_enabled"`
	RateLimitRequests int           `yaml:"rate_limit_requests"`
	RateLimitDuration time.Duration `yaml:"rate_limit_duration"`
}

// WebSocketServerConfig WebSocket服务器配置
type WebSocketServerConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	ReadBufferSize  int    `yaml:"read_buffer_size"`
	WriteBufferSize int    `yaml:"write_buffer_size"`
	CheckOrigin     bool   `yaml:"check_origin"`
}

// TCPServerConfig TCP服务器配置
type TCPServerConfig struct {
	Host               string        `yaml:"host"`
	Port               int           `yaml:"port"`
	MaxConnections     int           `yaml:"max_connections"`
	ReadTimeout        time.Duration `yaml:"read_timeout"`
	WriteTimeout       time.Duration `yaml:"write_timeout"`
	HeartbeatInterval  time.Duration `yaml:"heartbeat_interval"`
	KeepAliveInterval  time.Duration `yaml:"keep_alive_interval"`
	MaxPacketSize      int           `yaml:"max_packet_size"`
	CompressionEnabled bool          `yaml:"compression_enabled"`
	EncryptionEnabled  bool          `yaml:"encryption_enabled"`
	EnableMetrics      bool          `yaml:"enable_metrics"`
	BufferSize         int           `yaml:"buffer_size"`
}

// MetricsServerConfig 指标服务器配置
type MetricsServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
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
	Username       string        `yaml:"username"`
	Password       string        `yaml:"password"`
	AuthSource     string        `yaml:"auth_source"`
	MaxPoolSize    int           `yaml:"max_pool_size"`
	MinPoolSize    int           `yaml:"min_pool_size"`
	MaxIdleTime    time.Duration `yaml:"max_idle_time"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
	SocketTimeout  time.Duration `yaml:"socket_timeout"`
	RetryWrites    bool          `yaml:"retry_writes"`
	ReadPreference string        `yaml:"read_preference"`
	ReplicaSet     string        `yaml:"replica_set"`
	WriteConcern   WriteConcern  `yaml:"write_concern"`
}

// WriteConcern 写关注配置
type WriteConcern struct {
	W        string        `yaml:"w"`
	J        bool          `yaml:"j"`
	WTimeout time.Duration `yaml:"wtimeout"`
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
	PoolTimeout  time.Duration `yaml:"pool_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	MaxRetries   int           `yaml:"max_retries"`
	Cluster      ClusterConfig `yaml:"cluster"`
}

// ClusterConfig Redis集群配置
type ClusterConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Addresses []string `yaml:"addresses"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	DefaultTTL      time.Duration `yaml:"default_ttl"`
	MaxEntries      int64         `yaml:"max_entries"`
	CleanupInterval time.Duration `yaml:"cleanup_interval"`
	EvictionPolicy  string        `yaml:"eviction_policy"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	JWT        JWTConfig        `yaml:"jwt"`
	Encryption EncryptionConfig `yaml:"encryption"`
	CORS       CORSConfig       `yaml:"cors"`
	TLS        TLSConfig        `yaml:"tls"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string        `yaml:"secret"`
	Issuer          string        `yaml:"issuer"`
	Audience        string        `yaml:"audience"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
}

// EncryptionConfig 加密配置
type EncryptionConfig struct {
	Key       string `yaml:"key"`
	Algorithm string `yaml:"algorithm"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// TLSConfig TLS配置
type TLSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
	MinVersion string `yaml:"min_version"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level     string            `yaml:"level"`
	Format    string            `yaml:"format"`
	Output    string            `yaml:"output"`
	File      FileLogConfig     `yaml:"file"`
	Fields    map[string]string `yaml:"fields"`
	Sensitive []string          `yaml:"sensitive_fields"`
}

// FileLogConfig 文件日志配置
type FileLogConfig struct {
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// GameConfig 游戏配置
type GameConfig struct {
	Player     PlayerConfig     `yaml:"player"`
	Battle     BattleConfig     `yaml:"battle"`
	Experience ExperienceConfig `yaml:"experience"`
	Chat       ChatConfig       `yaml:"chat"`
}

// PlayerConfig 玩家配置
type PlayerConfig struct {
	MaxLevel          int   `yaml:"max_level"`
	InitialGold       int64 `yaml:"initial_gold"`
	InitialExperience int64 `yaml:"initial_experience"`
	MaxInventorySlots int   `yaml:"max_inventory_slots"`
}

// BattleConfig 战斗配置
type BattleConfig struct {
	MaxBattleTime      time.Duration `yaml:"max_battle_time"`
	DamageVariance     float64       `yaml:"damage_variance"`
	CriticalRateBase   float64       `yaml:"critical_rate_base"`
	CriticalDamageBase float64       `yaml:"critical_damage_base"`
}

// ExperienceConfig 经验配置
type ExperienceConfig struct {
	BaseExpPerLevel int     `yaml:"base_exp_per_level"`
	ExpMultiplier   float64 `yaml:"exp_multiplier"`
	MaxExpBonus     float64 `yaml:"max_exp_bonus"`
}

// ChatConfig 聊天配置
type ChatConfig struct {
	MaxMessageLength int      `yaml:"max_message_length"`
	RateLimit        int      `yaml:"rate_limit"`
	BannedWords      []string `yaml:"banned_words"`
	ProfanityFilter  bool     `yaml:"profanity_filter"`
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
	NATS NATSConfig `yaml:"nats"`
}

// NATSConfig NATS配置
type NATSConfig struct {
	URL           string          `yaml:"url"`
	ClusterID     string          `yaml:"cluster_id"`
	ClientID      string          `yaml:"client_id"`
	MaxReconnect  int             `yaml:"max_reconnect"`
	ReconnectWait time.Duration   `yaml:"reconnect_wait"`
	Timeout       time.Duration   `yaml:"timeout"`
	TLS           TLSConfig       `yaml:"tls"`
	JetStream     JetStreamConfig `yaml:"jetstream"`
	Subjects      SubjectsConfig  `yaml:"subjects"`
}

// JetStreamConfig JetStream配置
type JetStreamConfig struct {
	Enabled bool   `yaml:"enabled"`
	Domain  string `yaml:"domain"`
}

// SubjectsConfig 主题配置
type SubjectsConfig struct {
	PlayerEvents string `yaml:"player_events"`
	GameEvents   string `yaml:"game_events"`
	SystemEvents string `yaml:"system_events"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Health    HealthConfig    `yaml:"health"`
	Metrics   MetricsConfig   `yaml:"metrics"`
	Tracing   TracingConfig   `yaml:"tracing"`
	Profiling ProfilingConfig `yaml:"profiling"`
	Alerting  AlertingConfig  `yaml:"alerting"`
	Audit     AuditConfig     `yaml:"audit"`
}

// HealthConfig 健康检查配置
type HealthConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Namespace string `yaml:"namespace"`
}

// TracingConfig 链路追踪配置
type TracingConfig struct {
	Enabled        bool    `yaml:"enabled"`
	JaegerEndpoint string  `yaml:"jaeger_endpoint"`
	SampleRate     float64 `yaml:"sample_rate"`
}

// ProfilingConfig 性能分析配置
type ProfilingConfig struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
}

// AlertingConfig 告警配置
type AlertingConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

// AuditConfig 审计配置
type AuditConfig struct {
	Enabled       bool   `yaml:"enabled"`
	LogFile       string `yaml:"log_file"`
	RetentionDays int    `yaml:"retention_days"`
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
			env = "development"
		}
		configPath = fmt.Sprintf("configs/config.%s.yaml", env)
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

// setDefaults 设置默认值
func setDefaults(config *Config) {
	// 应用默认配置
	if config.App.Name == "" {
		config.App.Name = "GreatestWorks MMO Server"
	}
	if config.App.Version == "" {
		config.App.Version = "1.0.0"
	}
	if config.App.Environment == "" {
		config.App.Environment = "development"
	}

	// HTTP服务器默认配置
	if config.Server.HTTP.Host == "" {
		config.Server.HTTP.Host = "0.0.0.0"
	}
	if config.Server.HTTP.Port == 0 {
		config.Server.HTTP.Port = 8080
	}
	if config.Server.HTTP.ReadTimeout == 0 {
		config.Server.HTTP.ReadTimeout = 30 * time.Second
	}
	if config.Server.HTTP.WriteTimeout == 0 {
		config.Server.HTTP.WriteTimeout = 30 * time.Second
	}
	if config.Server.HTTP.IdleTimeout == 0 {
		config.Server.HTTP.IdleTimeout = 60 * time.Second
	}

	// TCP服务器默认配置
	if config.Server.TCP.Host == "" {
		config.Server.TCP.Host = "0.0.0.0"
	}
	if config.Server.TCP.Port == 0 {
		config.Server.TCP.Port = 8082
	}
	if config.Server.TCP.MaxConnections == 0 {
		config.Server.TCP.MaxConnections = 10000
	}

	// MongoDB默认配置
	if config.Database.MongoDB.URI == "" {
		config.Database.MongoDB.URI = "mongodb://localhost:27017"
	}
	if config.Database.MongoDB.Database == "" {
		config.Database.MongoDB.Database = "mmo_game"
	}
	if config.Database.MongoDB.MaxPoolSize == 0 {
		config.Database.MongoDB.MaxPoolSize = 100
	}
	if config.Database.MongoDB.MinPoolSize == 0 {
		config.Database.MongoDB.MinPoolSize = 10
	}

	// Redis默认配置
	if config.Database.Redis.Addr == "" {
		config.Database.Redis.Addr = "localhost:6379"
	}
	if config.Database.Redis.PoolSize == 0 {
		config.Database.Redis.PoolSize = 100
	}
	if config.Database.Redis.MinIdleConns == 0 {
		config.Database.Redis.MinIdleConns = 10
	}

	// 日志默认配置
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}
	if config.Logging.Output == "" {
		config.Logging.Output = "stdout"
	}

	// 游戏默认配置
	if config.Game.Player.MaxLevel == 0 {
		config.Game.Player.MaxLevel = 100
	}
	if config.Game.Player.InitialGold == 0 {
		config.Game.Player.InitialGold = 1000
	}
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	// 验证必需的配置项
	if config.Database.MongoDB.URI == "" {
		return fmt.Errorf("MongoDB URI is required")
	}
	if config.Database.MongoDB.Database == "" {
		return fmt.Errorf("MongoDB database name is required")
	}
	if config.Database.Redis.Addr == "" {
		return fmt.Errorf("Redis address is required")
	}
	if config.Security.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	// 验证端口范围
	if config.Server.HTTP.Port < 1 || config.Server.HTTP.Port > 65535 {
		return fmt.Errorf("invalid HTTP server port: %d", config.Server.HTTP.Port)
	}
	if config.Server.TCP.Port < 1 || config.Server.TCP.Port > 65535 {
		return fmt.Errorf("invalid TCP server port: %d", config.Server.TCP.Port)
	}

	// 验证游戏配置
	if config.Game.Player.MaxLevel < 1 {
		return fmt.Errorf("player max level must be greater than 0")
	}

	return nil
}

// GetEnvironment 获取当前环境
func GetEnvironment() string {
	env := os.Getenv("APP_ENV")
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
