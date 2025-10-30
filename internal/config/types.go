package config

import (
	"fmt"
	"strings"
	"time"
)

// Config coordinates all configuration sections for services.
type Config struct {
	App           AppConfig           `yaml:"app"`
	Service       ServiceConfig       `yaml:"service"`
	Server        ServerConfig        `yaml:"server"`
	Database      DatabaseConfig      `yaml:"database"`
	Messaging     MessagingConfig     `yaml:"messaging"`
	Logging       LoggingConfig       `yaml:"logging"`
	Security      SecurityConfig      `yaml:"security"`
	Monitoring    MonitoringConfig    `yaml:"monitoring"`
	Game          GameConfig          `yaml:"game"`
	Domain        DomainConfig        `yaml:"domain"`
	Application   ApplicationConfig   `yaml:"application"`
	Performance   PerformanceConfig   `yaml:"performance"`
	ThirdParty    ThirdPartyConfig    `yaml:"third_party"`
	Session       SessionConfig       `yaml:"session"`
	Gateway       GatewayConfig       `yaml:"gateway"`
	Environment   EnvironmentConfig   `yaml:"environment"`
	Observability ObservabilityConfig `yaml:"observability"`
}

// AppConfig contains global metadata.
type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
	Debug       bool   `yaml:"debug"`
}

// ServiceConfig captures per-service identifiers.
type ServiceConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
	NodeID      string `yaml:"node_id"`
	Region      string `yaml:"region"`
	Cluster     string `yaml:"cluster"`
}

// ServerConfig aggregates protocols served by the process.
type ServerConfig struct {
	HTTP    HTTPServerConfig    `yaml:"http"`
	RPC     RPCServerConfig     `yaml:"rpc"`
	TCP     TCPServerConfig     `yaml:"tcp"`
	GRPC    GRPCServerConfig    `yaml:"grpc"`
	Metrics MetricsServerConfig `yaml:"metrics"`
}

// HTTPServerConfig holds HTTP server details.
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
	EnableSwagger     bool          `yaml:"enable_swagger"`
	RateLimitEnabled  bool          `yaml:"rate_limit_enabled"`
	RateLimitRequests int           `yaml:"rate_limit_requests"`
	RateLimitWindow   time.Duration `yaml:"rate_limit_window"`
}

// RPCServerConfig configures internal RPC endpoints.
type RPCServerConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	MaxConnections  int           `yaml:"max_connections"`
	Timeout         time.Duration `yaml:"timeout"`
	KeepAlive       bool          `yaml:"keep_alive"`
	KeepAlivePeriod time.Duration `yaml:"keep_alive_period"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	TLS             TLSConfig     `yaml:"tls"`
}

// TCPServerConfig configures raw TCP listeners.
type TCPServerConfig struct {
	Host               string        `yaml:"host"`
	Port               int           `yaml:"port"`
	MaxConnections     int           `yaml:"max_connections"`
	ReadTimeout        time.Duration `yaml:"read_timeout"`
	WriteTimeout       time.Duration `yaml:"write_timeout"`
	HeartbeatEnabled   bool          `yaml:"heartbeat_enabled"`
	HeartbeatInterval  time.Duration `yaml:"heartbeat_interval"`
	HeartbeatTimeout   time.Duration `yaml:"heartbeat_timeout"`
	HeartbeatMaxMissed int           `yaml:"heartbeat_max_missed"`
	KeepAlive          bool          `yaml:"keep_alive"`
	KeepAliveInterval  time.Duration `yaml:"keep_alive_interval"`
	NoDelay            bool          `yaml:"no_delay"`
	MaxPacketSize      int           `yaml:"max_packet_size"`
	CompressionEnabled bool          `yaml:"compression_enabled"`
	EncryptionEnabled  bool          `yaml:"encryption_enabled"`
	BufferSize         int           `yaml:"buffer_size"`
}

// GRPCServerConfig configures gRPC endpoints.
type GRPCServerConfig struct {
	Host string    `yaml:"host"`
	Port int       `yaml:"port"`
	TLS  TLSConfig `yaml:"tls"`
}

// MetricsServerConfig configures /metrics endpoint.
type MetricsServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

// DatabaseConfig contains persistence providers.
type DatabaseConfig struct {
	MongoDB MongoDBConfig `yaml:"mongodb"`
	Redis   RedisConfig   `yaml:"redis"`
	SQL     SQLConfig     `yaml:"sql"`
}

// MongoDBConfig defines Mongo connection.
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
	ReplicaSet     string        `yaml:"replica_set"`
	RetryWrites    bool          `yaml:"retry_writes"`
}

// RedisConfig defines redis connection pool.
type RedisConfig struct {
	Addr         string             `yaml:"addr"`
	Password     string             `yaml:"password"`
	DB           int                `yaml:"db"`
	PoolSize     int                `yaml:"pool_size"`
	MinIdleConns int                `yaml:"min_idle_conns"`
	MaxRetries   int                `yaml:"max_retries"`
	DialTimeout  time.Duration      `yaml:"dial_timeout"`
	ReadTimeout  time.Duration      `yaml:"read_timeout"`
	WriteTimeout time.Duration      `yaml:"write_timeout"`
	PoolTimeout  time.Duration      `yaml:"pool_timeout"`
	IdleTimeout  time.Duration      `yaml:"idle_timeout"`
	TLS          TLSConfig          `yaml:"tls"`
	Cluster      RedisClusterConfig `yaml:"cluster"`
}

// RedisClusterConfig holds redis cluster endpoints.
type RedisClusterConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Addresses []string `yaml:"addresses"`
}

// SQLConfig describes relational database.
type SQLConfig struct {
	Driver          string        `yaml:"driver"`
	DSN             string        `yaml:"dsn"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// MessagingConfig contains message bus details.
type MessagingConfig struct {
	NATS      NATSConfig      `yaml:"nats"`
	Kafka     KafkaConfig     `yaml:"kafka"`
	SQS       SQSConfig       `yaml:"sqs"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
}

// NATSConfig describes JetStream.
type NATSConfig struct {
	URL           string          `yaml:"url"`
	ClusterID     string          `yaml:"cluster_id"`
	ClientID      string          `yaml:"client_id"`
	MaxReconnect  int             `yaml:"max_reconnect"`
	ReconnectWait time.Duration   `yaml:"reconnect_wait"`
	Timeout       time.Duration   `yaml:"timeout"`
	Credentials   string          `yaml:"credentials"`
	TLS           TLSConfig       `yaml:"tls"`
	JetStream     JetStreamConfig `yaml:"jetstream"`
	Subjects      SubjectsConfig  `yaml:"subjects"`
}

// JetStreamConfig toggles JetStream-specific options.
type JetStreamConfig struct {
	Enabled bool   `yaml:"enabled"`
	Domain  string `yaml:"domain"`
}

// SubjectsConfig enumerates subject keys.
type SubjectsConfig struct {
	PlayerEvents string `yaml:"player_events"`
	GameEvents   string `yaml:"game_events"`
	SystemEvents string `yaml:"system_events"`
	DomainEvents string `yaml:"domain_events"`
}

// KafkaConfig placeholder for future.
type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

// SQSConfig placeholder for AWS SQS.
type SQSConfig struct {
	QueueURL string        `yaml:"queue_url"`
	Timeout  time.Duration `yaml:"timeout"`
}

// SchedulerConfig defines internal scheduler defaults.
type SchedulerConfig struct {
	TickInterval time.Duration `yaml:"tick_interval"`
	MaxWorkers   int           `yaml:"max_workers"`
}

// LoggingConfig controls service logging.
type LoggingConfig struct {
	Level     string            `yaml:"level"`
	Format    string            `yaml:"format"`
	Output    string            `yaml:"output"`
	File      FileLogConfig     `yaml:"file"`
	Fields    map[string]string `yaml:"fields"`
	Sensitive []string          `yaml:"sensitive_fields"`
}

// FileLogConfig used when outputting to file.
type FileLogConfig struct {
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// SecurityConfig holds JWT, TLS, rate limiting.
type SecurityConfig struct {
	JWT            JWTConfig            `yaml:"jwt"`
	RateLimit      RateLimitConfig      `yaml:"rate_limit"`
	Encryption     EncryptionConfig     `yaml:"encryption"`
	PasswordPolicy PasswordPolicyConfig `yaml:"password_policy"`
	DDoSProtection DDoSProtectionConfig `yaml:"ddos_protection"`
	TLS            TLSConfig            `yaml:"tls"`
	CORS           CORSConfig           `yaml:"cors"`
}

// JWTConfig for auth tokens.
type JWTConfig struct {
	Secret          string        `yaml:"secret"`
	Issuer          string        `yaml:"issuer"`
	Audience        string        `yaml:"audience"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
}

// RateLimitConfig describes throttling.
type RateLimitConfig struct {
	Enabled           bool          `yaml:"enabled"`
	RequestsPerMinute int           `yaml:"requests_per_minute"`
	Burst             int           `yaml:"burst"`
	Interval          time.Duration `yaml:"interval"`
	GlobalLimit       int           `yaml:"global_limit"`
	PerIPLimit        int           `yaml:"per_ip_limit"`
}

// EncryptionConfig generic symmetric/asymmetric options.
type EncryptionConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Key       string `yaml:"key"`
	Algorithm string `yaml:"algorithm"`
}

// PasswordPolicyConfig defines password complexity requirements.
type PasswordPolicyConfig struct {
	MinLength        int           `yaml:"min_length"`
	RequireUppercase bool          `yaml:"require_uppercase"`
	RequireLowercase bool          `yaml:"require_lowercase"`
	RequireNumbers   bool          `yaml:"require_numbers"`
	RequireSymbols   bool          `yaml:"require_symbols"`
	MaxAttempts      int           `yaml:"max_attempts"`
	LockoutDuration  time.Duration `yaml:"lockout_duration"`
}

// DDoSProtectionConfig captures advanced network protection thresholds.
type DDoSProtectionConfig struct {
	Enabled      bool          `yaml:"enabled"`
	Threshold    int           `yaml:"threshold"`
	BanDuration  time.Duration `yaml:"ban_duration"`
	IPWhitelist  []string      `yaml:"ip_whitelist"`
	IPBlacklist  []string      `yaml:"ip_blacklist"`
	RateLimitKey string        `yaml:"rate_limit_key"`
}

// TLSConfig reused across sections.
type TLSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
	CAFile     string `yaml:"ca_file"`
	Insecure   bool   `yaml:"insecure"`
	MinVersion string `yaml:"min_version"`
}

// CORSConfig toggles HTTP cross-origin.
type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// MonitoringConfig collects metrics/tracing options.
type MonitoringConfig struct {
	Health    HealthConfig    `yaml:"health"`
	Metrics   MetricsConfig   `yaml:"metrics"`
	Tracing   TracingConfig   `yaml:"tracing"`
	Profiling ProfilingConfig `yaml:"profiling"`
	Alerting  AlertingConfig  `yaml:"alerting"`
	Audit     AuditConfig     `yaml:"audit"`
}

// HealthConfig toggles health endpoint.
type HealthConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// MetricsConfig represents legacy Prometheus settings. Deprecated: Prometheus
// metrics have been removed; keep for backward compatibility in configuration
// files only.
type MetricsConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Namespace string `yaml:"namespace"`
}

// TracingConfig configures distributed tracing.
type TracingConfig struct {
	Enabled     bool    `yaml:"enabled"`
	Endpoint    string  `yaml:"endpoint"`
	SampleRate  float64 `yaml:"sample_rate"`
	ServiceName string  `yaml:"service_name"`
}

// ProfilingConfig toggles pprof.
type ProfilingConfig struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
}

// AlertingConfig external alert endpoints.
type AlertingConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

// AuditConfig audit logs.
type AuditConfig struct {
	Enabled       bool   `yaml:"enabled"`
	LogFile       string `yaml:"log_file"`
	RetentionDays int    `yaml:"retention_days"`
}

// GameConfig domain-specific.
type GameConfig struct {
	Player     PlayerConfig     `yaml:"player"`
	Battle     BattleConfig     `yaml:"battle"`
	Experience ExperienceConfig `yaml:"experience"`
	Chat       ChatConfig       `yaml:"chat"`
	Ranking    RankingConfig    `yaml:"ranking"`
	Weather    WeatherConfig    `yaml:"weather"`
	Plant      PlantConfig      `yaml:"plant"`
}

// PlayerConfig domain defaults.
type PlayerConfig struct {
	MaxLevel          int           `yaml:"max_level"`
	InitialGold       int           `yaml:"initial_gold"`
	InitialExperience int           `yaml:"initial_experience"`
	MaxInventorySlots int           `yaml:"max_inventory_slots"`
	MaxFriends        int           `yaml:"max_friends"`
	SessionTimeout    time.Duration `yaml:"session_timeout"`
}

// BattleConfig domain defaults.
type BattleConfig struct {
	MaxBattleTime      time.Duration `yaml:"max_battle_time"`
	DamageVariance     float64       `yaml:"damage_variance"`
	CriticalRateBase   float64       `yaml:"critical_rate_base"`
	CriticalDamageBase float64       `yaml:"critical_damage_base"`
	MaxParticipants    int           `yaml:"max_participants"`
	TurnTimeout        time.Duration `yaml:"turn_timeout"`
}

// ExperienceConfig domain defaults.
type ExperienceConfig struct {
	BaseExpPerLevel int     `yaml:"base_exp_per_level"`
	ExpMultiplier   float64 `yaml:"exp_multiplier"`
	MaxExpBonus     float64 `yaml:"max_exp_bonus"`
}

// ChatConfig domain defaults.
type ChatConfig struct {
	MaxMessageLength int      `yaml:"max_message_length"`
	RateLimit        int      `yaml:"rate_limit"`
	BannedWords      []string `yaml:"banned_words"`
	SpamProtection   bool     `yaml:"spam_protection"`
}

// RankingConfig domain defaults.
type RankingConfig struct {
	MaxEntries     int           `yaml:"max_entries"`
	UpdateInterval time.Duration `yaml:"update_interval"`
	CacheTTL       time.Duration `yaml:"cache_ttl"`
}

// WeatherConfig domain defaults.
type WeatherConfig struct {
	UpdateInterval  time.Duration `yaml:"update_interval"`
	ForecastDays    int           `yaml:"forecast_days"`
	SeasonalEffects bool          `yaml:"seasonal_effects"`
}

// PlantConfig domain defaults.
type PlantConfig struct {
	GrowthSpeed  float64 `yaml:"growth_speed"`
	HarvestBonus float64 `yaml:"harvest_bonus"`
	MaxFarmSize  int     `yaml:"max_farm_size"`
}

// DomainConfig placeholder for more domain-level sections.
type DomainConfig struct {
	EnabledFeatures []string `yaml:"enabled_features"`
}

// ApplicationConfig cross-cutting app service settings.
type ApplicationConfig struct {
	CommandBus BusConfig `yaml:"command_bus"`
	QueryBus   BusConfig `yaml:"query_bus"`
	EventBus   BusConfig `yaml:"event_bus"`
}

// BusConfig for command/query/event bus.
type BusConfig struct {
	Timeout       time.Duration `yaml:"timeout"`
	RetryAttempts int           `yaml:"retry_attempts"`
	RetryDelay    time.Duration `yaml:"retry_delay"`
	CacheTTL      time.Duration `yaml:"cache_ttl"`
	DeadLetter    bool          `yaml:"dead_letter_queue"`
}

// PerformanceConfig runtime performance knobs.
type PerformanceConfig struct {
	WorkerPool     WorkerPoolConfig     `yaml:"worker_pool"`
	Cache          CacheConfig          `yaml:"cache"`
	RateLimit      RateLimitConfig      `yaml:"rate_limit"`
	ConnectionPool ConnectionPoolConfig `yaml:"connection_pool"`
}

// WorkerPoolConfig concurrency settings.
type WorkerPoolConfig struct {
	Size      int `yaml:"size"`
	QueueSize int `yaml:"queue_size"`
}

// CacheConfig general caching defaults.
type CacheConfig struct {
	DefaultTTL      time.Duration `yaml:"default_ttl"`
	MaxEntries      int           `yaml:"max_entries"`
	CleanupInterval time.Duration `yaml:"cleanup_interval"`
	EvictionPolicy  string        `yaml:"eviction_policy"`
}

// ConnectionPoolConfig defines shared connection pool tuning parameters.
type ConnectionPoolConfig struct {
	MaxIdle     int           `yaml:"max_idle"`
	MaxOpen     int           `yaml:"max_open"`
	MaxLifetime time.Duration `yaml:"max_lifetime"`
}

// ThirdPartyConfig external integration toggles.
type ThirdPartyConfig struct {
	Payment          PaymentConfig          `yaml:"payment"`
	PushNotification PushNotificationConfig `yaml:"push_notification"`
	Email            EmailConfig            `yaml:"email"`
	OAuth            OAuthConfig            `yaml:"oauth"`
}

// PaymentConfig payments provider config.
type PaymentConfig struct {
	Stripe StripeConfig `yaml:"stripe"`
}

// StripeConfig for Stripe integration.
type StripeConfig struct {
	PublicKey     string `yaml:"public_key"`
	SecretKey     string `yaml:"secret_key"`
	WebhookSecret string `yaml:"webhook_secret"`
}

// PushNotificationConfig push provider config.
type PushNotificationConfig struct {
	Firebase FirebaseConfig `yaml:"firebase"`
}

// FirebaseConfig details.
type FirebaseConfig struct {
	ServerKey string `yaml:"server_key"`
}

// EmailConfig email provider config.
type EmailConfig struct {
	SMTP SMTPConfig `yaml:"smtp"`
}

// OAuthConfig stores OAuth provider credentials.
type OAuthConfig struct {
	Providers map[string]OAuthProviderConfig `yaml:"providers"`
}

// OAuthProviderConfig contains individual provider credentials.
type OAuthProviderConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

// SMTPConfig SMTP transport.
type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// EnvironmentConfig environment toggles.
type EnvironmentConfig struct {
	HotReload bool `yaml:"hot_reload"`
	MockData  bool `yaml:"mock_data"`
	TestMode  bool `yaml:"test_mode"`
}

// ObservabilityConfig aggregator for metrics/tracing/logging combos.
type ObservabilityConfig struct {
	Enabled bool   `yaml:"enabled"`
	Backend string `yaml:"backend"`
}

// SessionConfig governs session lifecycle parameters.
type SessionConfig struct {
	MaxSessionsPerUser int           `yaml:"max_sessions_per_user"`
	SessionTimeout     time.Duration `yaml:"session_timeout"`
	CleanupInterval    time.Duration `yaml:"cleanup_interval"`
	StoreType          string        `yaml:"store_type"`
}

// GatewayConfig aggregates gateway-specific knobs used by the edge service.
type GatewayConfig struct {
	GameServices GatewayGameServicesConfig    `yaml:"game_services"`
	AuthService  GatewayExternalServiceConfig `yaml:"auth_service"`
	Connection   GatewayConnectionConfig      `yaml:"connection"`
	Protocol     GatewayProtocolConfig        `yaml:"protocol"`
	Routing      GatewayRoutingConfig         `yaml:"routing"`
}

// GatewayGameServicesConfig captures dependencies on downstream game services.
type GatewayGameServicesConfig struct {
	Discovery    GatewayDiscoveryConfig    `yaml:"discovery"`
	RPC          GatewayRPCConfig          `yaml:"rpc"`
	LoadBalancer GatewayLoadBalancerConfig `yaml:"load_balancer"`
}

// GatewayDiscoveryConfig defines discovery backends.
type GatewayDiscoveryConfig struct {
	Type   string              `yaml:"type"`
	Consul GatewayConsulConfig `yaml:"consul"`
	Etcd   GatewayEtcdConfig   `yaml:"etcd"`
	Static GatewayStaticConfig `yaml:"static"`
}

// GatewayConsulConfig configures Consul discovery.
type GatewayConsulConfig struct {
	Address     string `yaml:"address"`
	Datacenter  string `yaml:"datacenter"`
	ServiceName string `yaml:"service_name"`
}

// GatewayEtcdConfig configures Etcd discovery.
type GatewayEtcdConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

// GatewayStaticConfig provides static endpoints.
type GatewayStaticConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

// GatewayRPCConfig describes outbound RPC connectivity.
type GatewayRPCConfig struct {
	Protocol       string                      `yaml:"protocol"`
	Timeout        time.Duration               `yaml:"timeout"`
	RetryAttempts  int                         `yaml:"retry_attempts"`
	RetryDelay     time.Duration               `yaml:"retry_delay"`
	CircuitBreaker GatewayCircuitBreakerConfig `yaml:"circuit_breaker"`
}

// GatewayCircuitBreakerConfig contains circuit breaker knobs shared between rpc and auth.
type GatewayCircuitBreakerConfig struct {
	Enabled          bool          `yaml:"enabled"`
	FailureThreshold int           `yaml:"failure_threshold"`
	Timeout          time.Duration `yaml:"timeout"`
	MaxRequests      int           `yaml:"max_requests"`
}

// GatewayLoadBalancerConfig configures gateway load balancing strategy.
type GatewayLoadBalancerConfig struct {
	Strategy    string                   `yaml:"strategy"`
	HealthCheck GatewayHealthCheckConfig `yaml:"health_check"`
}

// GatewayHealthCheckConfig defines health probes for downstream services.
type GatewayHealthCheckConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
	Path     string        `yaml:"path"`
}

// GatewayExternalServiceConfig represents HTTP dependencies such as auth service.
type GatewayExternalServiceConfig struct {
	BaseURL        string                      `yaml:"base_url"`
	Timeout        time.Duration               `yaml:"timeout"`
	RetryAttempts  int                         `yaml:"retry_attempts"`
	RetryDelay     time.Duration               `yaml:"retry_delay"`
	CircuitBreaker GatewayCircuitBreakerConfig `yaml:"circuit_breaker"`
}

// GatewayConnectionConfig covers local connection pool and messaging details.
type GatewayConnectionConfig struct {
	MaxConnections    int                       `yaml:"max_connections"`
	ConnectionTimeout time.Duration             `yaml:"connection_timeout"`
	IdleTimeout       time.Duration             `yaml:"idle_timeout"`
	CleanupInterval   time.Duration             `yaml:"cleanup_interval"`
	Session           GatewaySessionConfig      `yaml:"session"`
	MessageQueue      GatewayMessageQueueConfig `yaml:"message_queue"`
}

// GatewaySessionConfig describes gateway session caches.
type GatewaySessionConfig struct {
	Timeout         time.Duration `yaml:"timeout"`
	CleanupInterval time.Duration `yaml:"cleanup_interval"`
	StoreType       string        `yaml:"store_type"`
}

// GatewayMessageQueueConfig configures message queue integration.
type GatewayMessageQueueConfig struct {
	Enabled  bool                       `yaml:"enabled"`
	Provider string                     `yaml:"provider"`
	Topics   GatewayMessageTopicsConfig `yaml:"topics"`
}

// GatewayMessageTopicsConfig enumerates queue topics.
type GatewayMessageTopicsConfig struct {
	PlayerEvents string `yaml:"player_events"`
	GameEvents   string `yaml:"game_events"`
	SystemEvents string `yaml:"system_events"`
}

// GatewayProtocolConfig defines protocol bridging.
type GatewayProtocolConfig struct {
	Client GatewayProtocolEndpointConfig `yaml:"client"`
	Game   GatewayProtocolEndpointConfig `yaml:"game"`
}

// GatewayProtocolEndpointConfig captures per endpoint protocol settings.
type GatewayProtocolEndpointConfig struct {
	Type        string `yaml:"type"`
	Codec       string `yaml:"codec"`
	Compression bool   `yaml:"compression"`
	Encryption  bool   `yaml:"encryption"`
}

// GatewayRoutingConfig configures message routing rules.
type GatewayRoutingConfig struct {
	Rules        []GatewayRoutingRule             `yaml:"rules"`
	LoadBalancer GatewayRoutingLoadBalancerConfig `yaml:"load_balancer"`
}

// GatewayRoutingRule describes a single routing entry.
type GatewayRoutingRule struct {
	Pattern string `yaml:"pattern"`
	Target  string `yaml:"target"`
	Method  string `yaml:"method"`
}

// GatewayRoutingLoadBalancerConfig tunes routing load balancing.
type GatewayRoutingLoadBalancerConfig struct {
	Strategy    string `yaml:"strategy"`
	HealthCheck bool   `yaml:"health_check"`
	Failover    bool   `yaml:"failover"`
}

// ApplyDefaults populates zero-value fields with opinionated defaults.
func (c *Config) ApplyDefaults() {
	if c.App.Name == "" {
		c.App.Name = "GreatestWorks"
	}
	if c.App.Version == "" {
		c.App.Version = "1.0.0"
	}
	if c.App.Environment == "" {
		c.App.Environment = "development"
	}

	if c.Service.Name == "" {
		c.Service.Name = c.App.Name
	}
	if c.Service.Version == "" {
		c.Service.Version = c.App.Version
	}
	if c.Service.Environment == "" {
		c.Service.Environment = c.App.Environment
	}
	if c.Service.Cluster == "" {
		c.Service.Cluster = c.App.Environment
	}

	if c.Server.HTTP.Host == "" {
		c.Server.HTTP.Host = "0.0.0.0"
	}
	if c.Server.HTTP.Port == 0 {
		c.Server.HTTP.Port = 8080
	}
	if c.Server.HTTP.ReadTimeout == 0 {
		c.Server.HTTP.ReadTimeout = 30 * time.Second
	}
	if c.Server.HTTP.WriteTimeout == 0 {
		c.Server.HTTP.WriteTimeout = 30 * time.Second
	}
	if c.Server.HTTP.IdleTimeout == 0 {
		c.Server.HTTP.IdleTimeout = 60 * time.Second
	}

	if c.Server.RPC.Host == "" {
		c.Server.RPC.Host = c.Server.HTTP.Host
	}
	if c.Server.RPC.Port == 0 {
		c.Server.RPC.Port = 8081
	}
	if c.Server.RPC.Timeout == 0 {
		c.Server.RPC.Timeout = 30 * time.Second
	}
	if c.Server.RPC.KeepAlivePeriod == 0 {
		c.Server.RPC.KeepAlivePeriod = 30 * time.Second
	}

	if c.Server.TCP.Host == "" {
		c.Server.TCP.Host = c.Server.HTTP.Host
	}
	if c.Server.TCP.Port == 0 {
		c.Server.TCP.Port = 9090
	}
	if c.Server.TCP.MaxConnections == 0 {
		c.Server.TCP.MaxConnections = 10000
	}
	if c.Server.TCP.HeartbeatInterval == 0 {
		c.Server.TCP.HeartbeatInterval = 30 * time.Second
	}
	if c.Server.TCP.HeartbeatTimeout == 0 {
		c.Server.TCP.HeartbeatTimeout = 10 * time.Second
	}
	if c.Server.TCP.HeartbeatMaxMissed == 0 {
		c.Server.TCP.HeartbeatMaxMissed = 3
	}
	if c.Server.TCP.KeepAliveInterval == 0 {
		c.Server.TCP.KeepAliveInterval = 30 * time.Second
	}

	if c.Server.Metrics.Host == "" {
		c.Server.Metrics.Host = "0.0.0.0"
	}
	if c.Server.Metrics.Port == 0 {
		c.Server.Metrics.Port = 9000
	}
	if c.Server.Metrics.Path == "" {
		c.Server.Metrics.Path = "/metrics"
	}

	if c.Database.MongoDB.URI == "" {
		c.Database.MongoDB.URI = "mongodb://localhost:27017"
	}
	if c.Database.MongoDB.Database == "" {
		c.Database.MongoDB.Database = "mmo_game"
	}
	if c.Database.MongoDB.MaxPoolSize == 0 {
		c.Database.MongoDB.MaxPoolSize = 100
	}
	if c.Database.MongoDB.MinPoolSize == 0 {
		c.Database.MongoDB.MinPoolSize = 5
	}
	if c.Database.MongoDB.ConnectTimeout == 0 {
		c.Database.MongoDB.ConnectTimeout = 10 * time.Second
	}
	if c.Database.MongoDB.SocketTimeout == 0 {
		c.Database.MongoDB.SocketTimeout = 30 * time.Second
	}

	if c.Database.Redis.Addr == "" {
		c.Database.Redis.Addr = "localhost:6379"
	}
	if c.Database.Redis.PoolSize == 0 {
		c.Database.Redis.PoolSize = 100
	}
	if c.Database.Redis.MinIdleConns == 0 {
		c.Database.Redis.MinIdleConns = 10
	}
	if c.Database.Redis.DialTimeout == 0 {
		c.Database.Redis.DialTimeout = 5 * time.Second
	}
	if c.Database.Redis.ReadTimeout == 0 {
		c.Database.Redis.ReadTimeout = 3 * time.Second
	}
	if c.Database.Redis.WriteTimeout == 0 {
		c.Database.Redis.WriteTimeout = 3 * time.Second
	}
	if c.Database.Redis.PoolTimeout == 0 {
		c.Database.Redis.PoolTimeout = 4 * time.Second
	}
	if c.Database.Redis.IdleTimeout == 0 {
		c.Database.Redis.IdleTimeout = 5 * time.Minute
	}

	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.Logging.Format == "" {
		c.Logging.Format = "json"
	}
	if c.Logging.Output == "" {
		c.Logging.Output = "stdout"
	}
	if c.Logging.Fields == nil {
		c.Logging.Fields = make(map[string]string)
	}

	if c.Security.JWT.Secret == "" {
		c.Security.JWT.Secret = "dev-secret-change-me"
	}
	if c.Security.JWT.AccessTokenTTL == 0 {
		c.Security.JWT.AccessTokenTTL = 15 * time.Minute
	}
	if c.Security.JWT.RefreshTokenTTL == 0 {
		c.Security.JWT.RefreshTokenTTL = 168 * time.Hour
	}
	if c.Security.RateLimit.RequestsPerMinute == 0 {
		c.Security.RateLimit.RequestsPerMinute = 1000
	}
	if c.Security.RateLimit.Burst == 0 {
		c.Security.RateLimit.Burst = 100
	}
	if c.Security.RateLimit.Interval == 0 {
		c.Security.RateLimit.Interval = time.Minute
	}
	if c.Security.RateLimit.GlobalLimit == 0 {
		c.Security.RateLimit.GlobalLimit = 1000
	}
	if c.Security.RateLimit.PerIPLimit == 0 {
		c.Security.RateLimit.PerIPLimit = 100
	}
	if c.Security.Encryption.Algorithm == "" {
		c.Security.Encryption.Algorithm = "AES-256-GCM"
	}
	if c.Security.PasswordPolicy.MinLength == 0 {
		c.Security.PasswordPolicy.MinLength = 8
	}
	if c.Security.PasswordPolicy.MaxAttempts == 0 {
		c.Security.PasswordPolicy.MaxAttempts = 5
	}
	if c.Security.PasswordPolicy.LockoutDuration == 0 {
		c.Security.PasswordPolicy.LockoutDuration = 15 * time.Minute
	}
	if c.Security.DDoSProtection.Threshold == 0 {
		c.Security.DDoSProtection.Threshold = 1000
	}
	if c.Security.DDoSProtection.BanDuration == 0 {
		c.Security.DDoSProtection.BanDuration = time.Hour
	}
	if c.Security.CORS.AllowedMethods == nil {
		c.Security.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}

	if c.Monitoring.Health.Path == "" {
		c.Monitoring.Health.Path = "/healthz"
	}
	if c.Monitoring.Metrics.Namespace == "" {
		c.Monitoring.Metrics.Namespace = "greatestworks"
	}
	if c.Monitoring.Profiling.Host == "" {
		c.Monitoring.Profiling.Host = "0.0.0.0"
	}
	if c.Monitoring.Profiling.Enabled && c.Monitoring.Profiling.Port == 0 {
		c.Monitoring.Profiling.Port = 6060
	}

	if c.Messaging.NATS.MaxReconnect == 0 {
		c.Messaging.NATS.MaxReconnect = 10
	}
	if c.Messaging.NATS.ReconnectWait == 0 {
		c.Messaging.NATS.ReconnectWait = 2 * time.Second
	}
	if c.Messaging.NATS.Timeout == 0 {
		c.Messaging.NATS.Timeout = 5 * time.Second
	}
	if c.Messaging.NATS.Subjects.PlayerEvents == "" {
		c.Messaging.NATS.Subjects.PlayerEvents = "player.events.>"
	}
	if c.Messaging.NATS.Subjects.GameEvents == "" {
		c.Messaging.NATS.Subjects.GameEvents = "game.events.>"
	}
	if c.Messaging.NATS.Subjects.SystemEvents == "" {
		c.Messaging.NATS.Subjects.SystemEvents = "system.events.>"
	}

	if c.Game.Player.MaxLevel == 0 {
		c.Game.Player.MaxLevel = 100
	}
	if c.Game.Player.MaxInventorySlots == 0 {
		c.Game.Player.MaxInventorySlots = 100
	}
	if c.Game.Battle.MaxBattleTime == 0 {
		c.Game.Battle.MaxBattleTime = 10 * time.Minute
	}
	if c.Game.Battle.TurnTimeout == 0 {
		c.Game.Battle.TurnTimeout = 30 * time.Second
	}
	if c.Game.Experience.BaseExpPerLevel == 0 {
		c.Game.Experience.BaseExpPerLevel = 100
	}
	if c.Game.Experience.ExpMultiplier == 0 {
		c.Game.Experience.ExpMultiplier = 1.2
	}
	if c.Game.Chat.MaxMessageLength == 0 {
		c.Game.Chat.MaxMessageLength = 500
	}

	if c.Application.CommandBus.Timeout == 0 {
		c.Application.CommandBus.Timeout = 5 * time.Second
	}
	if c.Application.QueryBus.Timeout == 0 {
		c.Application.QueryBus.Timeout = 5 * time.Second
	}
	if c.Application.EventBus.Timeout == 0 {
		c.Application.EventBus.Timeout = 5 * time.Second
	}

	if c.Performance.WorkerPool.Size == 0 {
		c.Performance.WorkerPool.Size = 100
	}
	if c.Performance.WorkerPool.QueueSize == 0 {
		c.Performance.WorkerPool.QueueSize = 1000
	}
	if c.Performance.Cache.DefaultTTL == 0 {
		c.Performance.Cache.DefaultTTL = time.Hour
	}
	if c.Performance.Cache.CleanupInterval == 0 {
		c.Performance.Cache.CleanupInterval = 10 * time.Minute
	}
	if c.Performance.Cache.EvictionPolicy == "" {
		c.Performance.Cache.EvictionPolicy = "lfu"
	}
	if c.Performance.ConnectionPool.MaxIdle == 0 {
		c.Performance.ConnectionPool.MaxIdle = 100
	}
	if c.Performance.ConnectionPool.MaxOpen == 0 {
		c.Performance.ConnectionPool.MaxOpen = 200
	}
	if c.Performance.ConnectionPool.MaxLifetime == 0 {
		c.Performance.ConnectionPool.MaxLifetime = time.Hour
	}

	if c.Session.MaxSessionsPerUser == 0 {
		c.Session.MaxSessionsPerUser = 3
	}
	if c.Session.SessionTimeout == 0 {
		c.Session.SessionTimeout = 24 * time.Hour
	}
	if c.Session.CleanupInterval == 0 {
		c.Session.CleanupInterval = time.Hour
	}
	if c.Session.StoreType == "" {
		c.Session.StoreType = "memory"
	}

	if c.Gateway.Connection.Session.Timeout == 0 {
		c.Gateway.Connection.Session.Timeout = 24 * time.Hour
	}
	if c.Gateway.Connection.Session.CleanupInterval == 0 {
		c.Gateway.Connection.Session.CleanupInterval = time.Hour
	}
}

// Validate ensures essential configuration values are present and acceptable.
func (c *Config) Validate() error {
	var problems []string

	if c.Database.MongoDB.URI == "" {
		problems = append(problems, "database.mongodb.uri is required")
	}
	if c.Database.MongoDB.Database == "" {
		problems = append(problems, "database.mongodb.database is required")
	}
	if c.Security.JWT.Secret == "" {
		problems = append(problems, "security.jwt.secret is required")
	}

	if !validPort(c.Server.HTTP.Port) {
		problems = append(problems, fmt.Sprintf("server.http.port out of range: %d", c.Server.HTTP.Port))
	}
	if !validPort(c.Server.RPC.Port) {
		problems = append(problems, fmt.Sprintf("server.rpc.port out of range: %d", c.Server.RPC.Port))
	}
	if !validPort(c.Server.TCP.Port) {
		problems = append(problems, fmt.Sprintf("server.tcp.port out of range: %d", c.Server.TCP.Port))
	}
	if !validPort(c.Server.Metrics.Port) {
		problems = append(problems, fmt.Sprintf("server.metrics.port out of range: %d", c.Server.Metrics.Port))
	}

	if len(problems) > 0 {
		return fmt.Errorf("config validation failed: %s", strings.Join(problems, "; "))
	}
	return nil
}

// Clone produces a shallow copy of the configuration with safe duplicated slices and maps.
func (c *Config) Clone() *Config {
	if c == nil {
		return nil
	}
	clone := *c
	clone.Logging.Fields = copyStringMap(c.Logging.Fields)
	clone.Logging.Sensitive = copyStringSlice(c.Logging.Sensitive)
	clone.Security.CORS.AllowedOrigins = copyStringSlice(c.Security.CORS.AllowedOrigins)
	clone.Security.CORS.AllowedMethods = copyStringSlice(c.Security.CORS.AllowedMethods)
	clone.Security.CORS.AllowedHeaders = copyStringSlice(c.Security.CORS.AllowedHeaders)
	clone.Security.CORS.ExposeHeaders = copyStringSlice(c.Security.CORS.ExposeHeaders)
	clone.Security.DDoSProtection.IPWhitelist = copyStringSlice(c.Security.DDoSProtection.IPWhitelist)
	clone.Security.DDoSProtection.IPBlacklist = copyStringSlice(c.Security.DDoSProtection.IPBlacklist)
	clone.Game.Chat.BannedWords = copyStringSlice(c.Game.Chat.BannedWords)
	clone.Domain.EnabledFeatures = copyStringSlice(c.Domain.EnabledFeatures)
	clone.ThirdParty.OAuth.Providers = copyOAuthProviders(c.ThirdParty.OAuth.Providers)
	clone.Gateway.GameServices.Discovery.Etcd.Endpoints = copyStringSlice(c.Gateway.GameServices.Discovery.Etcd.Endpoints)
	clone.Gateway.GameServices.Discovery.Static.Endpoints = copyStringSlice(c.Gateway.GameServices.Discovery.Static.Endpoints)
	clone.Gateway.Routing.Rules = copyGatewayRoutingRules(c.Gateway.Routing.Rules)
	return &clone
}

func validPort(port int) bool {
	return port >= 0 && port <= 65535
}

func copyStringSlice(src []string) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

func copyStringMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func copyOAuthProviders(src map[string]OAuthProviderConfig) map[string]OAuthProviderConfig {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]OAuthProviderConfig, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func copyGatewayRoutingRules(src []GatewayRoutingRule) []GatewayRoutingRule {
	if len(src) == 0 {
		return nil
	}
	dst := make([]GatewayRoutingRule, len(src))
	copy(dst, src)
	return dst
}
