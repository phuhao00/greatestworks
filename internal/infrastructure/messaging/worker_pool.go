package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// WorkerPool 工作池
type WorkerPool struct {
	workerCount int
	workQueue   chan interface{}
	workers     []*Worker
	processor   WorkerProcessor
	logger      logging.Logger
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	mu          sync.RWMutex
	running     bool
}

// Worker 工作者
type Worker struct {
	id        int
	workQueue chan interface{}
	processor WorkerProcessor
	logger    logging.Logger
	ctx       context.Context
	wg        *sync.WaitGroup
}

// WorkerProcessor 工作处理器接口
type WorkerProcessor interface {
	Process(ctx context.Context, work interface{}) error
}

// NewWorkerPool 创建工作者池
func NewWorkerPool(workerCount int, processor WorkerProcessor, logger logging.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		workerCount: workerCount,
		workQueue:   make(chan interface{}, workerCount*2), // 缓冲区大小为工作者数量的2倍
		processor:   processor,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
		workers:     make([]*Worker, 0, workerCount),
	}
}

// Start 启动工作者池
func (wp *WorkerPool) Start() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.running {
		return fmt.Errorf("worker pool is already running")
	}

	// 创建工作线程
	for i := 0; i < wp.workerCount; i++ {
		worker := &Worker{
			id:        i,
			workQueue: wp.workQueue,
			processor: wp.processor,
			logger:    wp.logger,
			ctx:       wp.ctx,
			wg:        &wp.wg,
		}

		wp.workers = append(wp.workers, worker)
		wp.wg.Add(1)
		go worker.start()
	}

	wp.running = true
	wp.logger.Info("Worker pool started", logging.Fields{
		"worker_count": wp.workerCount,
	})
	return nil
}

// Stop 停止工作者池
func (wp *WorkerPool) Stop() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.running {
		return nil
	}

	// 取消上下文
	wp.cancel()

	// 关闭工作队列
	close(wp.workQueue)

	// 等待所有工作者完成
	wp.wg.Wait()

	wp.running = false
	wp.logger.Info("Worker pool stopped")
	return nil
}

// Submit 提交工作
func (wp *WorkerPool) Submit(work interface{}) error {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if !wp.running {
		return fmt.Errorf("worker pool is not running")
	}

	select {
	case wp.workQueue <- work:
		wp.logger.Debug("Work submitted", logging.Fields{
			"work": work,
		})
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("worker pool is stopping")
	default:
		return fmt.Errorf("work queue is full")
	}
}

// SubmitWithTimeout 带超时的提交工作
func (wp *WorkerPool) SubmitWithTimeout(work interface{}, timeout time.Duration) error {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if !wp.running {
		return fmt.Errorf("worker pool is not running")
	}

	select {
	case wp.workQueue <- work:
		wp.logger.Debug("Work submitted with timeout", logging.Fields{
			"work":    work,
			"timeout": timeout,
		})
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("work submission timeout")
	case <-wp.ctx.Done():
		return fmt.Errorf("worker pool is stopping")
	}
}

// IsRunning 检查是否正在运行
func (wp *WorkerPool) IsRunning() bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.running
}

// GetWorkerCount 获取工作者数量
func (wp *WorkerPool) GetWorkerCount() int {
	return wp.workerCount
}

// GetQueueSize 获取队列大小
func (wp *WorkerPool) GetQueueSize() int {
	return len(wp.workQueue)
}

// GetStats 获取统计信息
func (wp *WorkerPool) GetStats() map[string]interface{} {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["running"] = wp.running
	stats["worker_count"] = wp.workerCount
	stats["queue_size"] = len(wp.workQueue)
	stats["queue_capacity"] = cap(wp.workQueue)

	return stats
}

// 工作者方法

// start 启动工作者
func (w *Worker) start() {
	defer w.wg.Done()

	w.logger.Debug("Worker started", logging.Fields{
		"worker_id": w.id,
	})

	for {
		select {
		case work, ok := <-w.workQueue:
			if !ok {
				w.logger.Debug("Work queue closed, worker stopping", logging.Fields{
					"worker_id": w.id,
				})
				return
			}

			w.processWork(work)

		case <-w.ctx.Done():
			w.logger.Debug("Worker context cancelled, stopping", logging.Fields{
				"worker_id": w.id,
			})
			return
		}
	}
}

// processWork 处理工作
func (w *Worker) processWork(work interface{}) {
	start := time.Now()

	w.logger.Debug("Worker processing work", logging.Fields{
		"worker_id": w.id,
		"work":      work,
	})

	// 处理工作
	if err := w.processor.Process(w.ctx, work); err != nil {
		w.logger.Error("Worker failed to process work", err, logging.Fields{
			"worker_id": w.id,
			"work":      work,
		})
		return
	}

	duration := time.Since(start)
	w.logger.Debug("Worker completed work", logging.Fields{
		"worker_id": w.id,
		"work":      work,
		"duration":  duration,
	})
}
