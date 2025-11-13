package managers

import (
	"context"
	"sync"
	"time"
)

// UpdateManager 更新管理器（单线程游戏主循环）
type UpdateManager struct {
	mu sync.Mutex

	running      bool
	deltaTime    float32
	tickRate     int           // 每秒更新次数
	tickInterval time.Duration // 更新间隔

	// 更新回调列表
	updateCallbacks []UpdateCallback

	// 定时器
	timers []*Timer

	// 任务队列
	taskQueue chan func()
}

// UpdateCallback 更新回调函数
type UpdateCallback func(ctx context.Context, deltaTime float32) error

// Timer 定时器
type Timer struct {
	ID       int64
	Interval float32 // 间隔时间
	Elapsed  float32 // 已过时间
	Repeat   bool    // 是否重复
	Callback func()
	Active   bool
}

var updateManagerInstance *UpdateManager
var updateManagerOnce sync.Once

// GetUpdateManager 获取更新管理器单例
func GetUpdateManager() *UpdateManager {
	updateManagerOnce.Do(func() {
		updateManagerInstance = &UpdateManager{
			tickRate:        30, // 默认30帧/秒
			updateCallbacks: make([]UpdateCallback, 0),
			timers:          make([]*Timer, 0),
			taskQueue:       make(chan func(), 1000),
		}
		updateManagerInstance.tickInterval = time.Second / time.Duration(updateManagerInstance.tickRate)
	})
	return updateManagerInstance
}

// RegisterUpdateCallback 注册更新回调
func (um *UpdateManager) RegisterUpdateCallback(callback UpdateCallback) {
	um.mu.Lock()
	defer um.mu.Unlock()

	um.updateCallbacks = append(um.updateCallbacks, callback)
}

// SetTickRate 设置更新频率
func (um *UpdateManager) SetTickRate(tickRate int) {
	um.mu.Lock()
	defer um.mu.Unlock()

	um.tickRate = tickRate
	um.tickInterval = time.Second / time.Duration(tickRate)
}

// AddTimer 添加定时器
func (um *UpdateManager) AddTimer(interval float32, repeat bool, callback func()) int64 {
	um.mu.Lock()
	defer um.mu.Unlock()

	timerID := int64(len(um.timers) + 1)
	timer := &Timer{
		ID:       timerID,
		Interval: interval,
		Elapsed:  0,
		Repeat:   repeat,
		Callback: callback,
		Active:   true,
	}

	um.timers = append(um.timers, timer)
	return timerID
}

// RemoveTimer 移除定时器
func (um *UpdateManager) RemoveTimer(timerID int64) {
	um.mu.Lock()
	defer um.mu.Unlock()

	for i, timer := range um.timers {
		if timer.ID == timerID {
			timer.Active = false
			um.timers = append(um.timers[:i], um.timers[i+1:]...)
			return
		}
	}
}

// PostTask 投递任务到主线程
func (um *UpdateManager) PostTask(task func()) {
	select {
	case um.taskQueue <- task:
	default:
		// 队列满，丢弃任务或者阻塞等待
	}
}

// Start 启动主循环
func (um *UpdateManager) Start(ctx context.Context) {
	um.mu.Lock()
	if um.running {
		um.mu.Unlock()
		return
	}
	um.running = true
	um.mu.Unlock()

	ticker := time.NewTicker(um.tickInterval)
	defer ticker.Stop()

	lastTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			um.Stop()
			return

		case <-ticker.C:
			now := time.Now()
			deltaTime := float32(now.Sub(lastTime).Seconds())
			lastTime = now

			um.deltaTime = deltaTime
			um.tick(ctx, deltaTime)

		case task := <-um.taskQueue:
			// 执行投递的任务
			if task != nil {
				task()
			}
		}
	}
}

// tick 单次更新
func (um *UpdateManager) tick(ctx context.Context, deltaTime float32) {
	um.mu.Lock()
	callbacks := make([]UpdateCallback, len(um.updateCallbacks))
	copy(callbacks, um.updateCallbacks)
	um.mu.Unlock()

	// 执行所有更新回调
	for _, callback := range callbacks {
		if err := callback(ctx, deltaTime); err != nil {
			// TODO: 记录错误日志
		}
	}

	// 更新定时器
	um.updateTimers(deltaTime)
}

// updateTimers 更新定时器
func (um *UpdateManager) updateTimers(deltaTime float32) {
	um.mu.Lock()
	defer um.mu.Unlock()

	toRemove := make([]int, 0)

	for i, timer := range um.timers {
		if !timer.Active {
			toRemove = append(toRemove, i)
			continue
		}

		timer.Elapsed += deltaTime
		if timer.Elapsed >= timer.Interval {
			// 触发回调
			if timer.Callback != nil {
				timer.Callback()
			}

			if timer.Repeat {
				timer.Elapsed -= timer.Interval
			} else {
				timer.Active = false
				toRemove = append(toRemove, i)
			}
		}
	}

	// 移除失效的定时器
	for i := len(toRemove) - 1; i >= 0; i-- {
		idx := toRemove[i]
		um.timers = append(um.timers[:idx], um.timers[idx+1:]...)
	}
}

// Stop 停止主循环
func (um *UpdateManager) Stop() {
	um.mu.Lock()
	defer um.mu.Unlock()

	um.running = false
}

// IsRunning 是否正在运行
func (um *UpdateManager) IsRunning() bool {
	um.mu.Lock()
	defer um.mu.Unlock()

	return um.running
}

// GetDeltaTime 获取帧间隔时间
func (um *UpdateManager) GetDeltaTime() float32 {
	um.mu.Lock()
	defer um.mu.Unlock()

	return um.deltaTime
}
