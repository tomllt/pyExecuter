package pyExecuter

import (
	"context"
	"fmt"
	"time"
)

// TimeoutController 任务超时控制结构体
type TimeoutController struct {
	taskTimeouts map[string]time.Duration // 保存任务的超时时间
}

// TimeoutControl 任务超时控制接口
type TimeoutControl interface {
	SetTaskTimeout(taskID string, duration time.Duration) error  // 设置任务的超时时间
	CheckTimeout(taskID string, startTime time.Time) bool        // 检查任务是否已超时
	HandleTimeout(taskID string) error                          // 处理超时任务
}

// NewTimeoutController 创建 TimeoutController 实例
func NewTimeoutController() *TimeoutController {
	return &TimeoutController{
		taskTimeouts: make(map[string]time.Duration),
	}
}

// SetTaskTimeout 设置任务的超时时间
func (t *TimeoutController) SetTaskTimeout(taskID string, duration time.Duration) error {
	t.taskTimeouts[taskID] = duration
	fmt.Printf("Task %s timeout set to %s\n", taskID, duration)
	return nil
}

// CheckTimeout 检查任务是否已经超时
func (t *TimeoutController) CheckTimeout(taskID string, startTime time.Time) bool {
	timeout, exists := t.taskTimeouts[taskID]
	if !exists {
		return false
	}
	if time.Since(startTime) > timeout {
		fmt.Printf("Task %s has timed out\n", taskID)
		return true
	}
	return false
}

// HandleTimeout 处理超时任务
func (t *TimeoutController) HandleTimeout(taskID string) error {
	// 假设有机制来取消执行中的任务
	fmt.Printf("Handling timeout for task %s: terminating task\n", taskID)
	// 在实际实现中，这里可以包括终止任务、释放资源等操作
	// 使用 context.Cancel 或直接调用任务执行器中的取消逻辑
	return fmt.Errorf("task %s timed out", taskID)
}

// Example usage during task execution
func ExecuteWithTimeoutControl(ctx context.Context, task *Task, timeoutControl TimeoutControl) error {
	startTime := time.Now()
	timeoutControl.SetTaskTimeout(task.ID, task.Timeout)

	// 执行任务前的超时检查
	select {
	case <-time.After(task.Timeout):
		return timeoutControl.HandleTimeout(task.ID)
	case <-ctx.Done():
		return ctx.Err()
	default:
		// 模拟任务执行逻辑
	}

	// 检查是否超时
	if timeoutControl.CheckTimeout(task.ID, startTime) {
		return timeoutControl.HandleTimeout(task.ID)
	}

	return nil
}