package internal

// ConfigManagerBase 配置管理器基类
type ConfigManagerBase struct {
}

// Load 加载配置
func (c *ConfigManagerBase) Load() {

}

// Get 获取配置项
func (c *ConfigManagerBase) Get(id uint32) interface{} {
	return nil
}
