package pyExecuter

import (
	"fmt"
	"time"
)

// ErrorHandler 处理任务错误及重试机制
type ErrorHandler struct {
	MaxRetryCount int           // 最大重试次数
	RetryInterval time.Duration // 重试间隔
}

// ErrorHandling 错误处理接口
type ErrorHandling interface {
	CaptureError(taskID string, err error) error  // 捕获任务执行中的异常
	RetryTask(taskID string) error                // 重试任务
}

// BasicErrorHandler 简单的错误处理实现
type BasicErrorHandler struct {
	MaxRetryCount  int           // 最大重试次数
	RetryInterval  time.Duration // 重试间隔
	retryCount     map[string]int // 记录任务的重试次数
	queue          *TaskQueue     // 用于重新将任务添加到队列
}

// NewBasicErrorHandler 创建 BasicErrorHandler 实例
func NewBasicErrorHandler(maxRetry int, retryInterval time.Duration, queue *TaskQueue) *BasicErrorHandler {
	return &BasicErrorHandler{
		MaxRetryCount:  maxRetry,
		RetryInterval:  retryInterval,
		retryCount:     make(map[string]int),
		queue:          queue,
	}
}

// CaptureError 处理任务执行中的异常
func (h *BasicErrorHandler) CaptureError(taskID string, err error) error {
	if h.retryCount[taskID] >= h.MaxRetryCount {
		return fmt.Errorf("task %s exceeded max retry count with error: %v", taskID, err)
	}

	fmt.Printf("Task %s encountered an error: %v. Retrying...\n", taskID, err)
	time.Sleep(h.RetryInterval)
	return h.RetryTask(taskID)
}

// RetryTask 重试任务
func (h *BasicErrorHandler) RetryTask(taskID string) error {
	h.retryCount[taskID]++
	task, err := h.queue.GetTaskByID(taskID) // 通过ID获取任务的逻辑需要实现
	if err != nil {
		return fmt.Errorf("failed to retrieve task %s for retry: %v", taskID, err)
	}

	err = h.queue.AddTask(task)
	if err != nil {
		return fmt.Errorf("failed to re-add task %s to queue: %v", taskID, err)
	}

	return nil
}