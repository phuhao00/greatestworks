// Package main 场景服务主程序
// 负责地图/副本/区域等场景的生命周期与调度
package main

import (
	"context"
	"greatestworks/internal/bootstrap"
	"log"
	"os"
	"os/signal"
	"syscall"

	"greatestworks/internal/config"
	"greatestworks/internal/infrastructure/logging"
)

// SceneServiceConfig 配置别名
type SceneServiceConfig = config.Config

// loadInitialConfig 加载配置
func loadInitialConfig() (*SceneServiceConfig, []string, *config.Loader, error) {
	loader := config.NewLoader(
		config.WithService("scene-service"),
	)

	cfg, files, err := loader.Load()
	if err != nil {
		return nil, nil, nil, err
	}
	return cfg, files, loader, nil
}

// main 主函数
func main() {
	logger := logging.NewBaseLogger(logging.InfoLevel)
	logger.Info("启动场景服务", logging.Fields{})

	cfg, sources, loader, err := loadInitialConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	logger.Info("配置加载成功", logging.Fields{
		"environment": cfg.App.Environment,
		"sources":     sources,
	})

	manager, err := config.NewManager(loader)
	if err != nil {
		log.Fatalf("创建配置管理器失败: %v", err)
	}
	defer func() { _ = manager.Close() }()

	runtimeCfg := manager.Config()
	service := bootstrap.NewSceneBootstrap(runtimeCfg, logger)

	manager.OnChange(func(next *config.Config) {
		if next == nil {
			return
		}
		service.UpdateConfig(next)
		logger.Info("场景服务配置已刷新", logging.Fields{
			"service_version": next.Service.Version,
		})
	})

	watchCtx, watchCancel := context.WithCancel(context.Background())
	defer watchCancel()

	if runtimeCfg != nil && runtimeCfg.Environment.HotReload {
		if err := manager.StartWatching(watchCtx); err != nil {
			logger.Error("启动配置热更新监听失败", err, logging.Fields{})
		} else {
			logger.Info("已启用配置热更新", logging.Fields{})
		}
	}

	if err := service.Start(); err != nil {
		log.Fatalf("启动场景服务失败: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case sig := <-sigChan:
		logger.Info("收到关闭信号", logging.Fields{"signal": sig.String()})
	case <-service.Done():
		logger.Info("上下文已取消", logging.Fields{})
	}

	logger.Info("正在关闭场景服务...", logging.Fields{})
	watchCancel()
	if err := service.Stop(); err != nil {
		logger.Error("关闭场景服务失败", err, logging.Fields{})
		os.Exit(1)
	}

	logger.Info("场景服务已关闭", logging.Fields{})
}
