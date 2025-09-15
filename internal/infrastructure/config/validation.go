// Package config 配置验证
// Author: MMO Server Team
// Created: 2024

package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if err := validateServerConfig(&config.Server); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	if err := validateDatabaseConfig(&config.Database); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}

	if err := validateSecurityConfig(&config.Security); err != nil {
		return fmt.Errorf("security config validation failed: %w", err)
	}

	if err := validateLoggingConfig(&config.Logging); err != nil {
		return fmt.Errorf("logging config validation failed: %w", err)
	}

	if err := validateGameConfig(&config.Game); err != nil {
		return fmt.Errorf("game config validation failed: %w", err)
	}

	if err := validateNetworkConfig(&config.Network); err != nil {
		return fmt.Errorf("network config validation failed: %w", err)
	}

	if err := validateMessagingConfig(&config.Messaging); err != nil {
		return fmt.Errorf("messaging config validation failed: %w", err)
	}

	if err := validateMonitoringConfig(&config.Monitoring); err != nil {
		return fmt.Errorf("monitoring config validation failed: %w", err)
	}

	if err := validateExcelConfig(&config.Excel); err != nil {
		return fmt.Errorf("excel config validation failed: %w", err)
	}

	return nil
}

// validateServerConfig 验证服务器配置
func validateServerConfig(config *ServerConfig) error {
	if err := validateGatewayConfig(&config.Gateway); err != nil {
		return fmt.Errorf("gateway config: %w", err)
	}

	if err := validateSceneConfig(&config.Scene); err != nil {
		return fmt.Errorf("scene config: %w", err)
	}

	if err := validateBattleConfig(&config.Battle); err != nil {
		return fmt.Errorf("battle config: %w", err)
	}

	if err := validateActivityConfig(&config.Activity); err != nil {
		return fmt.Errorf("activity config: %w", err)
	}

	return nil
}

// validateGatewayConfig 验证网关配置
func validateGatewayConfig(config *GatewayConfig) error {
	if err := validatePort(config.Port); err != nil {
		return fmt.Errorf("gateway port: %w", err)
	}

	if err := validateHost(config.Host); err != nil {
		return fmt.Errorf("gateway host: %w", err)
	}

	if config.MaxConnections <= 0 {
		return fmt.Errorf("max_connections must be positive")
	}

	if config.ReadTimeout <= 0 {
		return fmt.Errorf("read_timeout must be positive")
	}

	if config.WriteTimeout <= 0 {
		return fmt.Errorf("write_timeout must be positive")
	}

	if config.HeartbeatTime <= 0 {
		return fmt.Errorf("heartbeat_time must be positive")
	}

	if config.IdleTimeout <= 0 {
		return fmt.Errorf("idle_timeout must be positive")
	}

	// 验证TLS配置
	if config.TLSEnabled {
		if config.CertFile == "" {
			return fmt.Errorf("cert_file is required when TLS is enabled")
		}
		if config.KeyFile == "" {
			return fmt.Errorf("key_file is required when TLS is enabled")
		}
		if err := validateFileExists(config.CertFile); err != nil {
			return fmt.Errorf("cert_file: %w", err)
		}
		if err := validateFileExists(config.KeyFile); err != nil {
			return fmt.Errorf("key_file: %w", err)
		}
	}

	return nil
}

// validateSceneConfig 验证场景配置
func validateSceneConfig(config *SceneConfig) error {
	if err := validatePort(config.Port); err != nil {
		return fmt.Errorf("scene port: %w", err)
	}

	if err := validateHost(config.Host); err != nil {
		return fmt.Errorf("scene host: %w", err)
	}

	if config.MaxPlayers <= 0 {
		return fmt.Errorf("max_players must be positive")
	}

	if config.TickRate <= 0 || config.TickRate > 120 {
		return fmt.Errorf("tick_rate must be between 1 and 120")
	}

	if config.SyncInterval <= 0 {
		return fmt.Errorf("sync_interval must be positive")
	}

	if config.ViewDistance <= 0 {
		return fmt.Errorf("view_distance must be positive")
	}

	return nil
}

// validateBattleConfig 验证战斗配置
func validateBattleConfig(config *BattleConfig) error {
	if err := validatePort(config.Port); err != nil {
		return fmt.Errorf("battle port: %w", err)
	}

	if err := validateHost(config.Host); err != nil {
		return fmt.Errorf("battle host: %w", err)
	}

	if config.MaxBattles <= 0 {
		return fmt.Errorf("max_battles must be positive")
	}

	if config.TickRate <= 0 || config.TickRate > 120 {
		return fmt.Errorf("tick_rate must be between 1 and 120")
	}

	if config.BattleTime <= 0 {
		return fmt.Errorf("battle_time must be positive")
	}

	if config.MatchTimeout <= 0 {
		return fmt.Errorf("match_timeout must be positive")
	}

	return nil
}

// validateActivityConfig 验证活动配置
func validateActivityConfig(config *ActivityConfig) error {
	if err := validatePort(config.Port); err != nil {
		return fmt.Errorf("activity port: %w", err)
	}

	if err := validateHost(config.Host); err != nil {
		return fmt.Errorf("activity host: %w", err)
	}

	if config.MaxActivities <= 0 {
		return fmt.Errorf("max_activities must be positive")
	}

	if config.UpdateInterval <= 0 {
		return fmt.Errorf("update_interval must be positive")
	}

	if config.CacheTimeout <= 0 {
		return fmt.Errorf("cache_timeout must be positive")
	}

	return nil
}

// validateDatabaseConfig 验证数据库配置
func validateDatabaseConfig(config *DatabaseConfig) error {
	if err := validateMongoDBConfig(&config.MongoDB); err != nil {
		return fmt.Errorf("mongodb: %w", err)
	}

	if err := validateRedisConfig(&config.Redis); err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	return nil
}

// validateMongoDBConfig 验证MongoDB配置
func validateMongoDBConfig(config *MongoDBConfig) error {
	if config.URI == "" {
		return fmt.Errorf("uri is required")
	}

	if _, err := url.Parse(config.URI); err != nil {
		return fmt.Errorf("invalid uri format: %w", err)
	}

	if config.Database == "" {
		return fmt.Errorf("database name is required")
	}

	if config.MaxPoolSize <= 0 {
		return fmt.Errorf("max_pool_size must be positive")
	}

	if config.MinPoolSize < 0 {
		return fmt.Errorf("min_pool_size must be non-negative")
	}

	if config.MinPoolSize > config.MaxPoolSize {
		return fmt.Errorf("min_pool_size cannot be greater than max_pool_size")
	}

	if config.MaxIdleTime <= 0 {
		return fmt.Errorf("max_idle_time must be positive")
	}

	if config.ConnectTimeout <= 0 {
		return fmt.Errorf("connect_timeout must be positive")
	}

	if config.SocketTimeout <= 0 {
		return fmt.Errorf("socket_timeout must be positive")
	}

	validPreferences := []string{"primary", "primaryPreferred", "secondary", "secondaryPreferred", "nearest"}
	if !contains(validPreferences, config.ReadPreference) {
		return fmt.Errorf("invalid read_preference: %s", config.ReadPreference)
	}

	return nil
}

// validateRedisConfig 验证Redis配置
func validateRedisConfig(config *RedisConfig) error {
	if config.Addr == "" {
		return fmt.Errorf("addr is required")
	}

	if err := validateAddress(config.Addr); err != nil {
		return fmt.Errorf("invalid addr: %w", err)
	}

	if config.DB < 0 || config.DB > 15 {
		return fmt.Errorf("db must be between 0 and 15")
	}

	if config.PoolSize <= 0 {
		return fmt.Errorf("pool_size must be positive")
	}

	if config.MinIdleConns < 0 {
		return fmt.Errorf("min_idle_conns must be non-negative")
	}

	if config.MaxIdleConns < 0 {
		return fmt.Errorf("max_idle_conns must be non-negative")
	}

	if config.MinIdleConns > config.MaxIdleConns {
		return fmt.Errorf("min_idle_conns cannot be greater than max_idle_conns")
	}

	if config.ConnMaxAge <= 0 {
		return fmt.Errorf("conn_max_age must be positive")
	}

	if config.DialTimeout <= 0 {
		return fmt.Errorf("dial_timeout must be positive")
	}

	if config.ReadTimeout <= 0 {
		return fmt.Errorf("read_timeout must be positive")
	}

	if config.WriteTimeout <= 0 {
		return fmt.Errorf("write_timeout must be positive")
	}

	// 验证集群配置
	if config.Cluster.Enabled {
		if len(config.Cluster.Addresses) == 0 {
			return fmt.Errorf("cluster addresses are required when cluster is enabled")
		}
		for i, addr := range config.Cluster.Addresses {
			if err := validateAddress(addr); err != nil {
				return fmt.Errorf("invalid cluster address[%d]: %w", i, err)
			}
		}
	}

	return nil
}

// validateSecurityConfig 验证安全配置
func validateSecurityConfig(config *SecurityConfig) error {
	if err := validateJWTConfig(&config.JWT); err != nil {
		return fmt.Errorf("jwt: %w", err)
	}

	if err := validateEncryptionConfig(&config.Encryption); err != nil {
		return fmt.Errorf("encryption: %w", err)
	}

	if err := validateRateLimitConfig(&config.RateLimit); err != nil {
		return fmt.Errorf("rate_limit: %w", err)
	}

	return nil
}

// validateJWTConfig 验证JWT配置
func validateJWTConfig(config *JWTConfig) error {
	if config.SecretKey == "" {
		return fmt.Errorf("secret_key is required")
	}

	if len(config.SecretKey) < 32 {
		return fmt.Errorf("secret_key must be at least 32 characters")
	}

	if config.TokenDuration <= 0 {
		return fmt.Errorf("token_duration must be positive")
	}

	if config.RefreshTime <= 0 {
		return fmt.Errorf("refresh_time must be positive")
	}

	if config.TokenDuration >= config.RefreshTime {
		return fmt.Errorf("token_duration should be less than refresh_time")
	}

	return nil
}

// validateEncryptionConfig 验证加密配置
func validateEncryptionConfig(config *EncryptionConfig) error {
	validAlgorithms := []string{"AES-128-GCM", "AES-192-GCM", "AES-256-GCM", "ChaCha20-Poly1305"}
	if !contains(validAlgorithms, config.Algorithm) {
		return fmt.Errorf("unsupported algorithm: %s", config.Algorithm)
	}

	if config.KeySize <= 0 {
		return fmt.Errorf("key_size must be positive")
	}

	// 验证密钥大小与算法的匹配
	switch config.Algorithm {
	case "AES-128-GCM":
		if config.KeySize != 16 {
			return fmt.Errorf("AES-128-GCM requires key_size of 16")
		}
	case "AES-192-GCM":
		if config.KeySize != 24 {
			return fmt.Errorf("AES-192-GCM requires key_size of 24")
		}
	case "AES-256-GCM":
		if config.KeySize != 32 {
			return fmt.Errorf("AES-256-GCM requires key_size of 32")
		}
	case "ChaCha20-Poly1305":
		if config.KeySize != 32 {
			return fmt.Errorf("ChaCha20-Poly1305 requires key_size of 32")
		}
	}

	return nil
}

// validateRateLimitConfig 验证限流配置
func validateRateLimitConfig(config *RateLimitConfig) error {
	if config.Enabled {
		if config.RPS <= 0 {
			return fmt.Errorf("rps must be positive when rate limiting is enabled")
		}

		if config.Burst <= 0 {
			return fmt.Errorf("burst must be positive when rate limiting is enabled")
		}

		if config.Burst < config.RPS {
			return fmt.Errorf("burst should be greater than or equal to rps")
		}

		if config.WindowSize <= 0 {
			return fmt.Errorf("window_size must be positive when rate limiting is enabled")
		}
	}

	return nil
}

// validateLoggingConfig 验证日志配置
func validateLoggingConfig(config *LoggingConfig) error {
	validLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	if !contains(validLevels, strings.ToLower(config.Level)) {
		return fmt.Errorf("invalid log level: %s", config.Level)
	}

	validFormats := []string{"json", "text", "console"}
	if !contains(validFormats, strings.ToLower(config.Format)) {
		return fmt.Errorf("invalid log format: %s", config.Format)
	}

	validOutputs := []string{"stdout", "stderr", "file"}
	if !contains(validOutputs, strings.ToLower(config.Output)) {
		return fmt.Errorf("invalid log output: %s", config.Output)
	}

	if strings.ToLower(config.Output) == "file" {
		if config.Dir == "" {
			return fmt.Errorf("log dir is required when output is file")
		}

		if config.MaxSize <= 0 {
			return fmt.Errorf("max_size must be positive when output is file")
		}

		if config.MaxBackups < 0 {
			return fmt.Errorf("max_backups must be non-negative")
		}

		if config.MaxAge < 0 {
			return fmt.Errorf("max_age must be non-negative")
		}
	}

	return nil
}

// validateGameConfig 验证游戏配置
func validateGameConfig(config *GameConfig) error {
	if config.MaxLevel <= 0 {
		return fmt.Errorf("max_level must be positive")
	}

	if config.ExpMultiplier <= 0 {
		return fmt.Errorf("exp_multiplier must be positive")
	}

	if config.GoldMultiplier <= 0 {
		return fmt.Errorf("gold_multiplier must be positive")
	}

	if config.DropRate < 0 || config.DropRate > 1 {
		return fmt.Errorf("drop_rate must be between 0 and 1")
	}

	return nil
}

// validateNetworkConfig 验证网络配置
func validateNetworkConfig(config *NetworkConfig) error {
	validProtocols := []string{"tcp", "udp", "websocket"}
	if !contains(validProtocols, strings.ToLower(config.Protocol)) {
		return fmt.Errorf("unsupported protocol: %s", config.Protocol)
	}

	if config.BufferSize <= 0 {
		return fmt.Errorf("buffer_size must be positive")
	}

	if config.MaxPacketSize <= 0 {
		return fmt.Errorf("max_packet_size must be positive")
	}

	if config.BufferSize > config.MaxPacketSize {
		return fmt.Errorf("buffer_size should not exceed max_packet_size")
	}

	validCompressions := []string{"none", "gzip", "lz4", "snappy"}
	if !contains(validCompressions, strings.ToLower(config.CompressionType)) {
		return fmt.Errorf("unsupported compression type: %s", config.CompressionType)
	}

	validEncryptions := []string{"none", "aes", "chacha20"}
	if !contains(validEncryptions, strings.ToLower(config.EncryptionType)) {
		return fmt.Errorf("unsupported encryption type: %s", config.EncryptionType)
	}

	return nil
}

// validateMessagingConfig 验证消息队列配置
func validateMessagingConfig(config *MessagingConfig) error {
	if config.NSQ.Enabled {
		if err := validateNSQConfig(&config.NSQ); err != nil {
			return fmt.Errorf("nsq: %w", err)
		}
	}

	if config.RabbitMQ.Enabled {
		if err := validateRabbitMQConfig(&config.RabbitMQ); err != nil {
			return fmt.Errorf("rabbitmq: %w", err)
		}
	}

	return nil
}

// validateNSQConfig 验证NSQ配置
func validateNSQConfig(config *NSQConfig) error {
	if config.NSQDAddress == "" {
		return fmt.Errorf("nsqd_address is required")
	}

	if err := validateAddress(config.NSQDAddress); err != nil {
		return fmt.Errorf("invalid nsqd_address: %w", err)
	}

	if len(config.LookupdHTTP) == 0 {
		return fmt.Errorf("lookupd_http addresses are required")
	}

	for i, addr := range config.LookupdHTTP {
		if _, err := url.Parse("http://" + addr); err != nil {
			return fmt.Errorf("invalid lookupd_http[%d]: %w", i, err)
		}
	}

	if config.MaxInFlight <= 0 {
		return fmt.Errorf("max_in_flight must be positive")
	}

	return nil
}

// validateRabbitMQConfig 验证RabbitMQ配置
func validateRabbitMQConfig(config *RabbitMQConfig) error {
	if config.URL == "" {
		return fmt.Errorf("url is required")
	}

	if _, err := url.Parse(config.URL); err != nil {
		return fmt.Errorf("invalid url format: %w", err)
	}

	if config.Exchange == "" {
		return fmt.Errorf("exchange is required")
	}

	if config.Queue == "" {
		return fmt.Errorf("queue is required")
	}

	return nil
}

// validateMonitoringConfig 验证监控配置
func validateMonitoringConfig(config *MonitoringConfig) error {
	if config.Enabled {
		if err := validatePort(config.Port); err != nil {
			return fmt.Errorf("monitoring port: %w", err)
		}

		if config.Path == "" {
			return fmt.Errorf("path is required when monitoring is enabled")
		}

		if !strings.HasPrefix(config.Path, "/") {
			return fmt.Errorf("path must start with /")
		}

		if config.Prometheus.Enabled {
			if config.Prometheus.Namespace == "" {
				return fmt.Errorf("prometheus namespace is required")
			}

			if config.Prometheus.Subsystem == "" {
				return fmt.Errorf("prometheus subsystem is required")
			}
		}
	}

	return nil
}

// validateExcelConfig 验证Excel配置
func validateExcelConfig(config *ExcelConfig) error {
	if config.Path == "" {
		return fmt.Errorf("path is required")
	}

	// 检查路径是否存在
	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		return fmt.Errorf("excel path does not exist: %s", config.Path)
	}

	// 验证Excel文件配置
	excelFiles := map[string]string{
		"activity":   config.Activity,
		"battlepass": config.BattlePass,
		"pet":        config.Pet,
		"npc":        config.Npc,
		"plant":      config.Plant,
		"shop":       config.Shop,
		"task":       config.Task,
		"skill":      config.Skill,
		"vip":        config.Vip,
		"building":   config.Building,
		"condition":  config.Condition,
		"synthetise": config.Synthetise,
		"minigame":   config.MiniGame,
		"email":      config.Email,
	}

	for name, file := range excelFiles {
		if file != "" {
			fullPath := BuildRealPath(config.Path, file)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return fmt.Errorf("%s excel file does not exist: %s", name, fullPath)
			}
			if !strings.HasSuffix(strings.ToLower(file), ".xlsx") && !strings.HasSuffix(strings.ToLower(file), ".xls") {
				return fmt.Errorf("%s file must be an Excel file (.xlsx or .xls): %s", name, file)
			}
		}
	}

	return nil
}

// 辅助验证函数

// validatePort 验证端口号
func validatePort(port int) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}
	return nil
}

// validateHost 验证主机地址
func validateHost(host string) error {
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	// 检查是否为有效的IP地址
	if ip := net.ParseIP(host); ip != nil {
		return nil
	}

	// 检查是否为有效的主机名
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`, host); matched {
		return nil
	}

	return fmt.Errorf("invalid host format: %s", host)
}

// validateAddress 验证地址格式 (host:port)
func validateAddress(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address format: %w", err)
	}

	if err := validateHost(host); err != nil {
		return fmt.Errorf("invalid host in address: %w", err)
	}

	// 验证端口
	if port == "" {
		return fmt.Errorf("port is required")
	}

	return nil
}

// validateFileExists 验证文件是否存在
func validateFileExists(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if !filepath.IsAbs(filePath) {
		return fmt.Errorf("file path must be absolute: %s", filePath)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	return nil
}

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

// DefaultValidator 默认配置验证器
type DefaultValidator struct{}

// Validate 实现ConfigValidator接口
func (dv *DefaultValidator) Validate(config *Config) error {
	return validateConfig(config)
}

// BusinessValidator 业务规则验证器
type BusinessValidator struct{}

// Validate 实现ConfigValidator接口
func (bv *BusinessValidator) Validate(config *Config) error {
	// 业务规则验证

	// 检查服务器端口冲突
	ports := []int{
		config.Server.Gateway.Port,
		config.Server.Scene.Port,
		config.Server.Battle.Port,
		config.Server.Activity.Port,
		config.Monitoring.Port,
	}

	portMap := make(map[int]bool)
	for _, port := range ports {
		if port > 0 {
			if portMap[port] {
				return fmt.Errorf("port conflict detected: %d", port)
			}
			portMap[port] = true
		}
	}

	// 检查超时配置的合理性
	if config.Server.Gateway.ReadTimeout > config.Server.Gateway.IdleTimeout {
		return fmt.Errorf("gateway read_timeout should not exceed idle_timeout")
	}

	if config.Server.Gateway.WriteTimeout > config.Server.Gateway.IdleTimeout {
		return fmt.Errorf("gateway write_timeout should not exceed idle_timeout")
	}

	// 检查数据库连接池配置
	if config.Database.MongoDB.MinPoolSize > config.Database.MongoDB.MaxPoolSize/2 {
		return fmt.Errorf("mongodb min_pool_size should not exceed half of max_pool_size")
	}

	if config.Database.Redis.MinIdleConns > config.Database.Redis.PoolSize/2 {
		return fmt.Errorf("redis min_idle_conns should not exceed half of pool_size")
	}

	// 检查游戏配置的合理性
	if config.Game.ExpMultiplier > 10.0 {
		return fmt.Errorf("exp_multiplier seems too high: %f", config.Game.ExpMultiplier)
	}

	if config.Game.GoldMultiplier > 10.0 {
		return fmt.Errorf("gold_multiplier seems too high: %f", config.Game.GoldMultiplier)
	}

	return nil
}

// SecurityValidator 安全配置验证器
type SecurityValidator struct{}

// Validate 实现ConfigValidator接口
func (sv *SecurityValidator) Validate(config *Config) error {
	// 安全相关验证

	// 检查JWT密钥强度
	if len(config.Security.JWT.SecretKey) < 64 {
		return fmt.Errorf("jwt secret_key should be at least 64 characters for better security")
	}

	// 检查是否使用了默认密钥
	defaultKeys := []string{
		"your-secret-key",
		"default-secret",
		"123456",
		"password",
		"secret",
	}

	for _, defaultKey := range defaultKeys {
		if strings.Contains(strings.ToLower(config.Security.JWT.SecretKey), defaultKey) {
			return fmt.Errorf("jwt secret_key appears to use default or weak value")
		}
	}

	// 检查生产环境的安全配置
	if getEnvironment() == "prod" || getEnvironment() == "production" {
		if !config.Server.Gateway.TLSEnabled {
			return fmt.Errorf("TLS should be enabled in production environment")
		}

		if !config.Security.RateLimit.Enabled {
			return fmt.Errorf("rate limiting should be enabled in production environment")
		}

		if config.Logging.Level == "debug" || config.Logging.Level == "trace" {
			return fmt.Errorf("debug/trace logging should not be used in production")
		}
	}

	return nil
}