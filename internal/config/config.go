// Package config 配置管理和热重载
// Author: MMO Server Team
// Created: 2024

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
	// "github.com/phuhao00/netcore-go/core" // 暂时注释掉缺失的包
)

// Logger 简单的日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// Config 全局配置结构
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	NATS     NATSConfig     `json:"nats"`
	Gateway  GatewayConfig  `json:"gateway"`
	Scene    SceneConfig    `json:"scene"`
	Battle   BattleConfig   `json:"battle"`
	Activity ActivityConfig `json:"activity"`
	Login    LoginConfig    `json:"login"`
	Log      LogConfig      `json:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	MaxConns     int    `json:"max_conns"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Database     string `json:"database"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxLifetime  int    `json:"max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Password    string `json:"password"`
	DB          int    `json:"db"`
	PoolSize    int    `json:"pool_size"`
	MinIdleConn int    `json:"min_idle_conn"`
}

// NATSConfig NATS配置
type NATSConfig struct {
	URL            string `json:"url"`
	ClusterID      string `json:"cluster_id"`
	ClientID       string `json:"client_id"`
	MaxReconnect   int    `json:"max_reconnect"`
	ReconnectWait  int    `json:"reconnect_wait"`
	MaxReconnects  int    `json:"max_reconnects"`
	ConnectionName string `json:"connection_name"`
	DrainTimeout   int    `json:"drain_timeout"`
}

// GatewayConfig 网关配置
type GatewayConfig struct {
	Host           string   `json:"host"`
	Port           int      `json:"port"`
	MaxConnections int      `json:"max_connections"`
	Servers        []string `json:"servers"`
}

// SceneConfig 场景服务器配置
type SceneConfig struct {
	SceneID     string `json:"scene_id"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	GatewayAddr string `json:"gateway_addr"`
	MaxPlayers  int    `json:"max_players"`
}

// BattleConfig 战斗服务器配置
type BattleConfig struct {
	ServerID    string `json:"server_id"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	GatewayAddr string `json:"gateway_addr"`
	MaxBattles  int    `json:"max_battles"`
}

// ActivityConfig 活动服务器配置
type ActivityConfig struct {
	ServerID    string `json:"server_id"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	GatewayAddr string `json:"gateway_addr"`
}

// LoginConfig 登录服务器配置
type LoginConfig struct {
	ServerID    string `json:"server_id"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	GatewayAddr string `json:"gateway_addr"`
	JWTSecret   string `json:"jwt_secret"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// Manager 配置管理器
type Manager struct {
	config     *Config
	configPath string
	mutex      sync.RWMutex
	watchers   []func(*Config)
	logger     Logger
	stopChan   chan struct{}
}

// NewManager 创建配置管理器
func NewManager(configPath string, logger Logger) *Manager {
	return &Manager{
		configPath: configPath,
		watchers:   make([]func(*Config), 0),
		logger:     logger,
		stopChan:   make(chan struct{}),
	}
}

// Load 加载配置
func (m *Manager) Load() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	data, err := ioutil.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// 设置默认值
	m.setDefaults(&config)

	m.config = &config
	m.logger.Info("Configuration loaded", "path", m.configPath)

	// 通知观察者
	for _, watcher := range m.watchers {
		go watcher(&config)
	}

	return nil
}

// Get 获取配置
func (m *Manager) Get() *Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.config
}

// Watch 监听配置变化
func (m *Manager) Watch(callback func(*Config)) {
	m.mutex.Lock()
	m.watchers = append(m.watchers, callback)
	m.mutex.Unlock()
}

// StartWatching 开始监听配置文件变化
func (m *Manager) StartWatching() {
	go m.watchConfigFile()
}

// StopWatching 停止监听
func (m *Manager) StopWatching() {
	close(m.stopChan)
}

// watchConfigFile 监听配置文件变化
func (m *Manager) watchConfigFile() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	var lastModTime time.Time
	if stat, err := os.Stat(m.configPath); err == nil {
		lastModTime = stat.ModTime()
	}

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			stat, err := os.Stat(m.configPath)
			if err != nil {
				continue
			}

			if stat.ModTime().After(lastModTime) {
				lastModTime = stat.ModTime()
				m.logger.Info("Config file changed, reloading...")

				if err := m.Load(); err != nil {
					m.logger.Error("Failed to reload config", "error", err)
				}
			}
		}
	}
}

// setDefaults 设置默认值
func (m *Manager) setDefaults(config *Config) {
	// 服务器默认配置
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30
	}
	if config.Server.MaxConns == 0 {
		config.Server.MaxConns = 10000
	}

	// 数据库默认配置
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 27017
	}
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 100
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 10
	}
	if config.Database.MaxLifetime == 0 {
		config.Database.MaxLifetime = 3600
	}

	// Redis默认配置
	if config.Redis.Host == "" {
		config.Redis.Host = "localhost"
	}
	if config.Redis.Port == 0 {
		config.Redis.Port = 6379
	}
	if config.Redis.PoolSize == 0 {
		config.Redis.PoolSize = 10
	}
	if config.Redis.MinIdleConn == 0 {
		config.Redis.MinIdleConn = 5
	}

	// 日志默认配置
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if config.Log.Format == "" {
		config.Log.Format = "json"
	}
	if config.Log.Output == "" {
		config.Log.Output = "stdout"
	}
	if config.Log.MaxSize == 0 {
		config.Log.MaxSize = 100
	}
	if config.Log.MaxBackups == 0 {
		config.Log.MaxBackups = 3
	}
	if config.Log.MaxAge == 0 {
		config.Log.MaxAge = 7
	}
}

// LoadFromEnv 从环境变量加载配置
func (m *Manager) LoadFromEnv() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.config == nil {
		m.config = &Config{}
	}

	// 从环境变量覆盖配置
	if host := os.Getenv("SERVER_HOST"); host != "" {
		m.config.Server.Host = host
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		m.config.Database.Host = dbHost
	}
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		m.config.Redis.Host = redisHost
	}

	m.logger.Info("Environment variables loaded")
}

// Save 保存配置到文件
func (m *Manager) Save() error {
	m.mutex.RLock()
	config := m.config
	m.mutex.RUnlock()

	if config == nil {
		return fmt.Errorf("no config to save")
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := ioutil.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	m.logger.Info("Configuration saved", "path", m.configPath)
	return nil
}
