package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// WorkerPool 工作�?
type WorkerPool struct {
	workerCount int
	workQueue   chan interface{}
	workers     []*Worker
	processor   WorkerProcessor
	logger      logger.Logger
	ctx         context.Context
	cancel      context.CancelFunc
	stats       *WorkerPoolStats
	mu          sync.RWMutex
}

// WorkerProcessor 工作处理器接�?
type WorkerProcessor func(data interface{}) error

// Worker 工作�?
type Worker struct {
	id        int
	workQueue chan interface{}
	processor WorkerProcessor
	logger    logger.Logger
	ctx       context.Context
	stats     *WorkerStats
	mu        sync.RWMutex
}

// WorkerPoolStats 工作池统计信�?
type WorkerPoolStats struct {
	TotalProcessed int64                `json:"total_processed"`
	TotalFailed    int64                `json:"total_failed"`
	ActiveWorkers  int64                `json:"active_workers"`
	QueueSize      int64                `json:"queue_size"`
	StartTime      time.Time            `json:"start_time"`
	Uptime         time.Duration        `json:"uptime"`
	ByWorker       map[int]*WorkerStats `json:"by_worker"`
}

// WorkerStats 工作者统计信�?
type WorkerStats struct {
	ProcessedCount int64         `json:"processed_count"`
	FailedCount    int64         `json:"failed_count"`
	LastProcessed  time.Time     `json:"last_processed"`
	AvgProcessTime time.Duration `json:"avg_process_time"`
	IsActive       bool          `json:"is_active"`
}

// NewWorkerPool 创建工作�?
func NewWorkerPool(workerCount int, processor WorkerProcessor, logger logger.Logger) *WorkerPool {
	if workerCount <= 0 {
		workerCount = 10 // 默认工作者数�?
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		workerCount: workerCount,
		workQueue:   make(chan interface{}, workerCount*10), // 队列大小为工作者数量的10�?
		workers:     make([]*Worker, workerCount),
		processor:   processor,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
		stats: &WorkerPoolStats{
			StartTime: time.Now(),
			ByWorker:  make(map[int]*WorkerStats),
		},
	}

	// 创建工作�?
	for i := 0; i < workerCount; i++ {
		worker := &Worker{
			id:        i + 1,
			workQueue: pool.workQueue,
			processor: processor,
			logger:    logger,
			ctx:       ctx,
			stats: &WorkerStats{
				IsActive: false,
			},
		}

		pool.workers[i] = worker
		pool.stats.ByWorker[worker.id] = worker.stats
	}

	logger.Info("Worker pool created successfully", "worker_count", workerCount, "queue_capacity", cap(pool.workQueue))
	return pool
}

// Start 启动工作�?
func (p *WorkerPool) Start(ctx context.Context) error {
	p.logger.Info("Starting worker pool", "worker_count", p.workerCount)

	// 启动所有工作�?
	for _, worker := range p.workers {
		go worker.start()
	}

	// 启动统计收集
	go p.collectStats()

	p.logger.Info("Worker pool started successfully")
	return nil
}

// Stop 停止工作�?
func (p *WorkerPool) Stop() error {
	p.logger.Info("Stopping worker pool")

	// 取消上下文，停止所有工作�?
	p.cancel()

	// 关闭工作队列
	close(p.workQueue)

	// 等待所有工作者停�?
	for _, worker := range p.workers {
		worker.stop()
	}

	p.logger.Info("Worker pool stopped successfully")
	return nil
}

// Submit 提交任务
func (p *WorkerPool) Submit(data interface{}) error {
	select {
	case p.workQueue <- data:
		p.logger.Debug("Task submitted to worker pool")
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool is stopped")
	default:
		return fmt.Errorf("worker pool queue is full")
	}
}

// SubmitWithTimeout 带超时的提交任务
func (p *WorkerPool) SubmitWithTimeout(data interface{}, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(p.ctx, timeout)
	defer cancel()

	select {
	case p.workQueue <- data:
		p.logger.Debug("Task submitted to worker pool with timeout")
		return nil
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("submit timeout after %v", timeout)
		}
		return ctx.Err()
	}
}

// GetStats 获取工作池统计信�?
func (p *WorkerPool) GetStats() *WorkerPoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// 计算活跃工作者数�?
	activeWorkers := int64(0)
	for _, worker := range p.workers {
		worker.mu.RLock()
		if worker.stats.IsActive {
			activeWorkers++
		}
		worker.mu.RUnlock()
	}

	// 创建统计信息副本
	stats := &WorkerPoolStats{
		TotalProcessed: p.stats.TotalProcessed,
		TotalFailed:    p.stats.TotalFailed,
		ActiveWorkers:  activeWorkers,
		QueueSize:      int64(len(p.workQueue)),
		StartTime:      p.stats.StartTime,
		Uptime:         time.Since(p.stats.StartTime),
		ByWorker:       make(map[int]*WorkerStats),
	}

	// 复制工作者统计信�?
	for id, workerStats := range p.stats.ByWorker {
		stats.ByWorker[id] = &WorkerStats{
			ProcessedCount: workerStats.ProcessedCount,
			FailedCount:    workerStats.FailedCount,
			LastProcessed:  workerStats.LastProcessed,
			AvgProcessTime: workerStats.AvgProcessTime,
			IsActive:       workerStats.IsActive,
		}
	}

	return stats
}

// GetQueueSize 获取队列大小
func (p *WorkerPool) GetQueueSize() int {
	return len(p.workQueue)
}

// GetWorkerCount 获取工作者数�?
func (p *WorkerPool) GetWorkerCount() int {
	return p.workerCount
}

// IsRunning 检查工作池是否运行�?
func (p *WorkerPool) IsRunning() bool {
	select {
	case <-p.ctx.Done():
		return false
	default:
		return true
	}
}

// 私有方法

// collectStats 收集统计信息
func (p *WorkerPool) collectStats() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := p.GetStats()
			p.logger.Debug("Worker pool metrics",
				"total_processed", stats.TotalProcessed,
				"total_failed", stats.TotalFailed,
				"active_workers", stats.ActiveWorkers,
				"queue_size", stats.QueueSize,
				"uptime", stats.Uptime)
		case <-p.ctx.Done():
			return
		}
	}
}

// updateStats 更新统计信息
func (p *WorkerPool) updateStats(success bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if success {
		p.stats.TotalProcessed++
	} else {
		p.stats.TotalFailed++
	}
}

// Worker 方法

// start 启动工作�?
func (w *Worker) start() {
	w.logger.Debug("Worker started", "worker_id", w.id)

	for {
		select {
		case data := <-w.workQueue:
			if data != nil {
				w.processTask(data)
			}
		case <-w.ctx.Done():
			w.logger.Debug("Worker stopped", "worker_id", w.id)
			return
		}
	}
}

// stop 停止工作�?
func (w *Worker) stop() {
	w.mu.Lock()
	w.stats.IsActive = false
	w.mu.Unlock()

	w.logger.Debug("Worker stopping", "worker_id", w.id)
}

// processTask 处理任务
func (w *Worker) processTask(data interface{}) {
	start := time.Now()

	// 标记为活�?
	w.mu.Lock()
	w.stats.IsActive = true
	w.mu.Unlock()

	defer func() {
		// 标记为非活跃
		w.mu.Lock()
		w.stats.IsActive = false
		w.mu.Unlock()
	}()

	w.logger.Debug("Worker processing task", "worker_id", w.id)

	// 处理任务
	err := w.processor(data)
	processTime := time.Since(start)

	// 更新统计信息
	w.updateStats(err == nil, processTime)

	if err != nil {
		w.logger.Error("Worker task processing failed", "error", err, "worker_id", w.id, "process_time", processTime)
	} else {
		w.logger.Debug("Worker task processed successfully", "worker_id", w.id, "process_time", processTime)
	}
}

// updateStats 更新工作者统计信�?
func (w *Worker) updateStats(success bool, processTime time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if success {
		w.stats.ProcessedCount++
		w.stats.LastProcessed = time.Now()

		// 更新平均处理时间
		if w.stats.AvgProcessTime == 0 {
			w.stats.AvgProcessTime = processTime
		} else {
			w.stats.AvgProcessTime = (w.stats.AvgProcessTime + processTime) / 2
		}
	} else {
		w.stats.FailedCount++
	}
}

// GetStats 获取工作者统计信�?
func (w *Worker) GetStats() *WorkerStats {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return &WorkerStats{
		ProcessedCount: w.stats.ProcessedCount,
		FailedCount:    w.stats.FailedCount,
		LastProcessed:  w.stats.LastProcessed,
		AvgProcessTime: w.stats.AvgProcessTime,
		IsActive:       w.stats.IsActive,
	}
}

// GetID 获取工作者ID
func (w *Worker) GetID() int {
	return w.id
}

// IsActive 检查工作者是否活�?
func (w *Worker) IsActive() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.stats.IsActive
}

// 工作池配�?
type WorkerPoolConfig struct {
	WorkerCount     int           `json:"worker_count" yaml:"worker_count"`
	QueueSize       int           `json:"queue_size" yaml:"queue_size"`
	TaskTimeout     time.Duration `json:"task_timeout" yaml:"task_timeout"`
	EnableMetrics   bool          `json:"enable_metrics" yaml:"enable_metrics"`
	MetricsInterval time.Duration `json:"metrics_interval" yaml:"metrics_interval"`
}

// NewWorkerPoolFromConfig 从配置创建工作池
func NewWorkerPoolFromConfig(config *WorkerPoolConfig, processor WorkerProcessor, logger logger.Logger) *WorkerPool {
	if config == nil {
		config = &WorkerPoolConfig{
			WorkerCount:     10,
			QueueSize:       100,
			TaskTimeout:     30 * time.Second,
			EnableMetrics:   true,
			MetricsInterval: 30 * time.Second,
		}
	}

	pool := NewWorkerPool(config.WorkerCount, processor, logger)

	// 如果指定了队列大小，重新创建队列
	if config.QueueSize > 0 {
		pool.workQueue = make(chan interface{}, config.QueueSize)
		// 更新所有工作者的队列引用
		for _, worker := range pool.workers {
			worker.workQueue = pool.workQueue
		}
	}

	logger.Info("Worker pool created from config",
		"worker_count", config.WorkerCount,
		"queue_size", config.QueueSize,
		"task_timeout", config.TaskTimeout,
		"enable_metrics", config.EnableMetrics)

	return pool
}
