package events

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// EventDispatcher 事件分发器
type EventDispatcher struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
	GetEventTypes() []string
	GetHandlerName() string
}

// NewEventDispatcher 创建事件分发器
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[EventType][]EventHandler),
	}
}

// RegisterHandler 注册事件处理器
func (d *EventDispatcher) RegisterHandler(eventType EventType, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch 分发事件
func (d *EventDispatcher) Dispatch(ctx context.Context, event Event) error {
	d.mu.RLock()
	handlers, exists := d.handlers[EventType(event.GetType())]
	d.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no handlers for event type: %s", event.GetType())
	}

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return fmt.Errorf("handler error: %w", err)
		}
	}

	return nil
}

// EventTask 事件任务
type EventTask struct {
	Event      Event
	Context    context.Context
	Dispatcher *EventDispatcher
}

// WorkerPool 工作池
type WorkerPool struct {
	workerCount int
	TaskQueue   chan *EventTask
	workers     []*Worker
	logger      *log.Logger
	mu          sync.RWMutex
	running     bool
	wg          sync.WaitGroup
}

// Worker 工作者
type Worker struct {
	id       int
	taskChan chan *EventTask
	quit     chan bool
	logger   *log.Logger
	metrics  *WorkerMetrics
}

// WorkerMetrics 工作者指标
type WorkerMetrics struct {
	TasksProcessed uint64
	TasksFailed    uint64
	TotalTime      time.Duration
	mu             sync.RWMutex
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workerCount, queueSize int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		TaskQueue:   make(chan *EventTask, queueSize),
		workers:     make([]*Worker, workerCount),
		logger:      log.New(log.Writer(), "[WorkerPool] ", log.LstdFlags),
	}
}

// Start 启动工作池
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.running {
		return
	}

	wp.logger.Printf("Starting worker pool with %d workers", wp.workerCount)

	// 创建并启动工作者
	for i := 0; i < wp.workerCount; i++ {
		worker := &Worker{
			id:       i + 1,
			taskChan: wp.TaskQueue,
			quit:     make(chan bool),
			logger:   log.New(log.Writer(), fmt.Sprintf("[Worker-%d] ", i+1), log.LstdFlags),
			metrics:  &WorkerMetrics{},
		}
		wp.workers[i] = worker

		wp.wg.Add(1)
		go worker.start(&wp.wg)
	}

	wp.running = true
	wp.logger.Println("Worker pool started")
}

// Stop 停止工作池
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.running {
		return
	}

	wp.logger.Println("Stopping worker pool...")

	// 停止所有工作者
	for _, worker := range wp.workers {
		worker.stop()
	}

	// 等待所有工作者完成
	wp.wg.Wait()

	// 关闭任务队列
	close(wp.TaskQueue)

	wp.running = false
	wp.logger.Println("Worker pool stopped")
}

// GetMetrics 获取工作池指标
func (wp *WorkerPool) GetMetrics() map[string]interface{} {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	metrics := map[string]interface{}{
		"worker_count":   wp.workerCount,
		"queue_size":     cap(wp.TaskQueue),
		"queue_length":   len(wp.TaskQueue),
		"running":        wp.running,
		"worker_metrics": make([]map[string]interface{}, len(wp.workers)),
	}

	for i, worker := range wp.workers {
		if worker != nil {
			metrics["worker_metrics"].([]map[string]interface{})[i] = worker.getMetrics()
		}
	}

	return metrics
}

// start 启动工作者
func (w *Worker) start(wg *sync.WaitGroup) {
	defer wg.Done()
	w.logger.Printf("Worker %d started", w.id)

	for {
		select {
		case task := <-w.taskChan:
			if task != nil {
				w.processTask(task)
			}
		case <-w.quit:
			w.logger.Printf("Worker %d stopping", w.id)
			return
		}
	}
}

// stop 停止工作者
func (w *Worker) stop() {
	close(w.quit)
}

// processTask 处理任务
func (w *Worker) processTask(task *EventTask) {
	start := time.Now()
	defer func() {
		w.metrics.mu.Lock()
		w.metrics.TotalTime += time.Since(start)
		w.metrics.mu.Unlock()
	}()

	w.logger.Printf("Processing event: %s (type: %s)", task.Event.GetID(), task.Event.GetType())

	// 处理事件
	if err := task.Dispatcher.Dispatch(task.Context, task.Event); err != nil {
		w.logger.Printf("Failed to process event %s: %v", task.Event.GetID(), err)
		w.metrics.mu.Lock()
		w.metrics.TasksFailed++
		w.metrics.mu.Unlock()
	} else {
		w.logger.Printf("Successfully processed event: %s", task.Event.GetID())
	}

	w.metrics.mu.Lock()
	w.metrics.TasksProcessed++
	w.metrics.mu.Unlock()
}

// getMetrics 获取工作者指标
func (w *Worker) getMetrics() map[string]interface{} {
	w.metrics.mu.RLock()
	defer w.metrics.mu.RUnlock()

	return map[string]interface{}{
		"id":              w.id,
		"tasks_processed": w.metrics.TasksProcessed,
		"tasks_failed":    w.metrics.TasksFailed,
		"total_time":      w.metrics.TotalTime.String(),
		"avg_time":        w.getAverageProcessingTime(),
	}
}

// getAverageProcessingTime 获取平均处理时间
func (w *Worker) getAverageProcessingTime() string {
	if w.metrics.TasksProcessed == 0 {
		return "0s"
	}
	avg := w.metrics.TotalTime / time.Duration(w.metrics.TasksProcessed)
	return avg.String()
}
