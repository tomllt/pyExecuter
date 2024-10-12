package pyExecuter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TimeoutController 任务超时控制结构体
type TimeoutController struct {
	taskTimeouts     map[string]time.Time // 保存任务的截止时间
	taskCancellations map[string]context.CancelFunc // 保存任务的取消函数
	mu               sync.RWMutex
}

// TimeoutControl 任务超时控制接口
type TimeoutControl interface {
	SetTaskTimeout(taskID string, duration time.Duration) error  // 设置任务的超时时间
	CheckTimeout(taskID string) (bool, error)                    // 检查任务是否已超时
	HandleTimeout(taskID string) error                           // 处理超时任务
	ClearTimeout(taskID string) error                            // 清理任务的超时设置
}

// NewTimeoutController 创建 TimeoutController 实例
func NewTimeoutController() *TimeoutController {
	return &TimeoutController{
		taskTimeouts:     make(map[string]time.Time),
		taskCancellations: make(map[string]context.CancelFunc),
	}
}

// SetTaskTimeout 设置任务的超时时间
func (t *TimeoutController) SetTaskTimeout(taskID string, duration time.Duration) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	deadline := time.Now().Add(duration)
	t.taskTimeouts[taskID] = deadline

	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	t.taskCancellations[taskID] = cancel

	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			t.HandleTimeout(taskID)
		}
	}()

	fmt.Printf("Task %s timeout set to %s\n", taskID, duration)
	return nil
}

// CheckTimeout 检查任务是否已经超时
func (t *TimeoutController) CheckTimeout(taskID string) (bool, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	deadline, exists := t.taskTimeouts[taskID]
	if !exists {
		return false, fmt.Errorf("no timeout set for task %s", taskID)
	}

	if time.Now().After(deadline) {
		fmt.Printf("Task %s has timed out\n", taskID)
		return true, nil
	}
	return false, nil
}

// HandleTimeout 处理超时任务
func (t *TimeoutController) HandleTimeout(taskID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	cancel, exists := t.taskCancellations[taskID]
	if !exists {
		return fmt.Errorf("no cancellation function for task %s", taskID)
	}

	cancel() // 取消任务
	delete(t.taskTimeouts, taskID)
	delete(t.taskCancellations, taskID)

	fmt.Printf("Handling timeout for task %s: terminating task\n", taskID)
	return nil
}

// ClearTimeout 清理任务的超时设置
func (t *TimeoutController) ClearTimeout(taskID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	cancel, exists := t.taskCancellations[taskID]
	if !exists {
		return fmt.Errorf("no timeout set for task %s", taskID)
	}

	cancel() // 取消定时器
	delete(t.taskTimeouts, taskID)
	delete(t.taskCancellations, taskID)

	fmt.Printf("Cleared timeout for task %s\n", taskID)
	return nil
}

// ExecuteWithTimeoutControl 演示如何使用 TimeoutControl
func ExecuteWithTimeoutControl(ctx context.Context, task *Task, timeoutControl TimeoutControl) error {
	err := timeoutControl.SetTaskTimeout(task.ID, task.Timeout)
	if err != nil {
		return fmt.Errorf("failed to set timeout: %v", err)
	}
	defer timeoutControl.ClearTimeout(task.ID)

	// 模拟任务执行
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(task.Timeout / 2): // 假设任务在超时之前完成
		// 任务完成，检查是否已超时
		timedOut, err := timeoutControl.CheckTimeout(task.ID)
		if err != nil {
			return fmt.Errorf("failed to check timeout: %v", err)
		}
		if timedOut {
			return fmt.Errorf("task %s timed out", task.ID)
		}
		fmt.Printf("Task %s completed successfully\n", task.ID)
		return nil
	}
}