package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logger"

	"gopkg.in/yaml.v3"
)

// ConfigLoader 配置加载器
type ConfigLoader struct {
	logger     logger.Logger
	mu         sync.RWMutex
	configs    map[string]interface{}
	watchers   map[string]*FileWatcher
	validators map[string]Validator
	env        string
	basePath   string
}

// LoaderConfig 加载器配置
type LoaderConfig struct {
	Environment    string   `json:"environment" yaml:"environment"`
	BasePath       string   `json:"base_path" yaml:"base_path"`
	ConfigFiles    []string `json:"config_files" yaml:"config_files"`
	WatchFiles     bool     `json:"watch_files" yaml:"watch_files"`
	EnvPrefix      string   `json:"env_prefix" yaml:"env_prefix"`
	ValidateOnLoad bool     `json:"validate_on_load" yaml:"validate_on_load"`
}

// Validator 配置验证器接口
type Validator interface {
	// Validate 验证配置
	Validate(config interface{}) error

	// GetConfigName 获取配置名称
	GetConfigName() string
}

// Loader 配置加载器接口
type Loader interface {
	// LoadConfig 加载配置
	LoadConfig(name string, config interface{}) error

	// LoadConfigFromFile 从文件加载配置
	LoadConfigFromFile(filePath string, config interface{}) error

	// LoadConfigFromEnv 从环境变量加载配置
	LoadConfigFromEnv(prefix string, config interface{}) error

	// SaveConfig 保存配置
	SaveConfig(name string, config interface{}) error

	// GetConfig 获取配置
	GetConfig(name string) (interface{}, error)

	// SetConfig 设置配置
	SetConfig(name string, config interface{}) error

	// RegisterValidator 注册验证器
	RegisterValidator(validator Validator) error

	// WatchConfig 监听配置变化
	WatchConfig(name string, callback func(interface{})) error

	// GetEnvironment 获取环境
	GetEnvironment() string

	// SetEnvironment 设置环境
	SetEnvironment(env string)

	// Reload 重新加载所有配置
	Reload() error

	// Close 关闭加载器
	Close() error
}

// NewConfigLoader 创建配置加载器
func NewConfigLoader(config *LoaderConfig, logger logger.Logger) Loader {
	if config == nil {
		config = &LoaderConfig{
			Environment:    getEnvOrDefault("APP_ENV", "development"),
			BasePath:       "./configs",
			ConfigFiles:    []string{"app.yaml", "database.yaml", "redis.yaml"},
			WatchFiles:     true,
			EnvPrefix:      "APP",
			ValidateOnLoad: true,
		}
	}

	loader := &ConfigLoader{
		logger:     logger,
		configs:    make(map[string]interface{}),
		watchers:   make(map[string]*FileWatcher),
		validators: make(map[string]Validator),
		env:        config.Environment,
		basePath:   config.BasePath,
	}

	// 创建配置目录
	if err := os.MkdirAll(config.BasePath, 0755); err != nil {
		logger.Error("Failed to create config directory", "error", err, "path", config.BasePath)
	}

	// 加载默认配置文件
	for _, configFile := range config.ConfigFiles {
		filePath := filepath.Join(config.BasePath, configFile)
		if _, err := os.Stat(filePath); err == nil {
			name := strings.TrimSuffix(configFile, filepath.Ext(configFile))
			var configData interface{}
			if err := loader.LoadConfigFromFile(filePath, &configData); err != nil {
				logger.Error("Failed to load config file", "error", err, "file", configFile)
			} else {
				loader.SetConfig(name, configData)

				// 启动文件监听
				if config.WatchFiles {
					loader.WatchConfig(name, func(newConfig interface{}) {
						logger.Info("Config file changed", "name", name)
					})
				}
			}
		}
	}

	logger.Info("Config loader initialized successfully", "environment", config.Environment, "base_path", config.BasePath)
	return loader
}

// LoadConfig 加载配置
func (l *ConfigLoader) LoadConfig(name string, config interface{}) error {
	// 首先尝试从文件加载
	filePath := l.getConfigFilePath(name)
	if _, err := os.Stat(filePath); err == nil {
		if err := l.LoadConfigFromFile(filePath, config); err != nil {
			l.logger.Error("Failed to load config from file", "error", err, "name", name, "file", filePath)
			return err
		}
	} else {
		// 文件不存在，使用默认配置
		l.logger.Debug("Config file not found, using defaults", "name", name, "file", filePath)
	}

	// 从环境变量覆盖配置
	if err := l.LoadConfigFromEnv(strings.ToUpper(name), config); err != nil {
		l.logger.Error("Failed to load config from environment", "error", err, "name", name)
		return err
	}

	// 验证配置
	if validator, exists := l.validators[name]; exists {
		if err := validator.Validate(config); err != nil {
			l.logger.Error("Config validation failed", "error", err, "name", name)
			return fmt.Errorf("config validation failed for %s: %w", name, err)
		}
	}

	// 存储配置
	l.SetConfig(name, config)

	l.logger.Info("Config loaded successfully", "name", name)
	return nil
}

// LoadConfigFromFile 从文件加载配置
func (l *ConfigLoader) LoadConfigFromFile(filePath string, config interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		err = json.Unmarshal(data, config)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, config)
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to unmarshal config file %s: %w", filePath, err)
	}

	l.logger.Debug("Config loaded from file", "file", filePath)
	return nil
}

// LoadConfigFromEnv 从环境变量加载配置
func (l *ConfigLoader) LoadConfigFromEnv(prefix string, config interface{}) error {
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// 跳过非导出字段
		if !field.CanSet() {
			continue
		}

		// 获取环境变量名
		envName := l.getEnvName(prefix, fieldType)
		envValue := os.Getenv(envName)

		if envValue == "" {
			continue
		}

		// 设置字段值
		if err := l.setFieldValue(field, envValue); err != nil {
			l.logger.Error("Failed to set field value from env", "error", err, "field", fieldType.Name, "env", envName)
			continue
		}

		l.logger.Debug("Field set from environment", "field", fieldType.Name, "env", envName)
	}

	return nil
}

// SaveConfig 保存配置
func (l *ConfigLoader) SaveConfig(name string, config interface{}) error {
	filePath := l.getConfigFilePath(name)

	// 创建目录
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// 序列化配置
	var data []byte
	var err error

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		data, err = json.MarshalIndent(config, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(config)
	default:
		// 默认使用YAML格式
		filePath = strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".yaml"
		data, err = yaml.Marshal(config)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 写入文件
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	l.logger.Info("Config saved successfully", "name", name, "file", filePath)
	return nil
}

// GetConfig 获取配置
func (l *ConfigLoader) GetConfig(name string) (interface{}, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	config, exists := l.configs[name]
	if !exists {
		return nil, fmt.Errorf("config %s not found", name)
	}

	return config, nil
}

// SetConfig 设置配置
func (l *ConfigLoader) SetConfig(name string, config interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.configs[name] = config

	l.logger.Debug("Config set successfully", "name", name)
	return nil
}

// RegisterValidator 注册验证器
func (l *ConfigLoader) RegisterValidator(validator Validator) error {
	name := validator.GetConfigName()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.validators[name] = validator

	l.logger.Info("Validator registered successfully", "config_name", name)
	return nil
}

// WatchConfig 监听配置变化
func (l *ConfigLoader) WatchConfig(name string, callback func(interface{})) error {
	filePath := l.getConfigFilePath(name)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("config file %s does not exist", filePath)
	}

	// 创建文件监听器
	watcher, err := NewFileWatcher(filePath, func(path string) {
		l.logger.Info("Config file changed, reloading", "name", name, "file", path)

		// 重新加载配置
		var newConfig interface{}
		if err := l.LoadConfigFromFile(filePath, &newConfig); err != nil {
			l.logger.Error("Failed to reload config", "error", err, "name", name)
			return
		}

		// 更新配置
		l.SetConfig(name, newConfig)

		// 调用回调
		if callback != nil {
			callback(newConfig)
		}
	}, l.logger)

	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	// 启动监听
	if err := watcher.Start(); err != nil {
		return fmt.Errorf("failed to start file watcher: %w", err)
	}

	l.mu.Lock()
	l.watchers[name] = watcher
	l.mu.Unlock()

	l.logger.Info("Config watcher started", "name", name, "file", filePath)
	return nil
}

// GetEnvironment 获取环境
func (l *ConfigLoader) GetEnvironment() string {
	return l.env
}

// SetEnvironment 设置环境
func (l *ConfigLoader) SetEnvironment(env string) {
	l.env = env
	l.logger.Info("Environment changed", "environment", env)
}

// Reload 重新加载所有配置
func (l *ConfigLoader) Reload() error {
	l.mu.RLock()
	configNames := make([]string, 0, len(l.configs))
	for name := range l.configs {
		configNames = append(configNames, name)
	}
	l.mu.RUnlock()

	for _, name := range configNames {
		var config interface{}
		if err := l.LoadConfig(name, &config); err != nil {
			l.logger.Error("Failed to reload config", "error", err, "name", name)
			continue
		}
	}

	l.logger.Info("All configs reloaded successfully")
	return nil
}

// Close 关闭加载器
func (l *ConfigLoader) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 停止所有文件监听器
	for name, watcher := range l.watchers {
		if err := watcher.Stop(); err != nil {
			l.logger.Error("Failed to stop file watcher", "error", err, "name", name)
		}
	}

	// 清空数据
	l.configs = make(map[string]interface{})
	l.watchers = make(map[string]*FileWatcher)
	l.validators = make(map[string]Validator)

	l.logger.Info("Config loader closed successfully")
	return nil
}

// 私有方法

// getConfigFilePath 获取配置文件路径
func (l *ConfigLoader) getConfigFilePath(name string) string {
	// 尝试不同的文件扩展名
	extensions := []string{".yaml", ".yml", ".json"}

	for _, ext := range extensions {
		filePath := filepath.Join(l.basePath, l.env, name+ext)
		if _, err := os.Stat(filePath); err == nil {
			return filePath
		}
	}

	// 如果环境特定的文件不存在，尝试通用文件
	for _, ext := range extensions {
		filePath := filepath.Join(l.basePath, name+ext)
		if _, err := os.Stat(filePath); err == nil {
			return filePath
		}
	}

	// 默认返回YAML格式的环境特定文件路径
	return filepath.Join(l.basePath, l.env, name+".yaml")
}

// getEnvName 获取环境变量名
func (l *ConfigLoader) getEnvName(prefix string, field reflect.StructField) string {
	// 检查是否有env标签
	if envTag := field.Tag.Get("env"); envTag != "" {
		return envTag
	}

	// 使用字段名生成环境变量名
	fieldName := field.Name

	// 转换为大写并添加下划线
	envName := strings.ToUpper(prefix + "_" + fieldName)

	// 处理驼峰命名
	envName = strings.ReplaceAll(envName, "ID", "_ID")
	envName = strings.ReplaceAll(envName, "URL", "_URL")
	envName = strings.ReplaceAll(envName, "API", "_API")

	return envName
}

// setFieldValue 设置字段值
func (l *ConfigLoader) setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			// 处理time.Duration类型
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("failed to parse duration: %w", err)
			}
			field.SetInt(int64(duration))
		} else {
			intVal, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse int: %w", err)
			}
			field.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %w", err)
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float: %w", err)
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("failed to parse bool: %w", err)
		}
		field.SetBool(boolVal)
	case reflect.Slice:
		// 处理字符串切片
		if field.Type().Elem().Kind() == reflect.String {
			sliceVal := strings.Split(value, ",")
			for i, v := range sliceVal {
				sliceVal[i] = strings.TrimSpace(v)
			}
			field.Set(reflect.ValueOf(sliceVal))
		} else {
			return fmt.Errorf("unsupported slice type: %s", field.Type())
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// getEnvOrDefault 获取环境变量或默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 配置结构体示例
type AppConfig struct {
	Name       string        `json:"name" yaml:"name" env:"APP_NAME"`
	Version    string        `json:"version" yaml:"version" env:"APP_VERSION"`
	Port       int           `json:"port" yaml:"port" env:"APP_PORT"`
	Host       string        `json:"host" yaml:"host" env:"APP_HOST"`
	Debug      bool          `json:"debug" yaml:"debug" env:"APP_DEBUG"`
	Timeout    time.Duration `json:"timeout" yaml:"timeout" env:"APP_TIMEOUT"`
	Features   []string      `json:"features" yaml:"features" env:"APP_FEATURES"`
	MaxWorkers int           `json:"max_workers" yaml:"max_workers" env:"APP_MAX_WORKERS"`
}
