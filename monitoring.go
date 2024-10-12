package pyExecuter

import (
	"fmt"
	"time"
)

// TaskMonitoring 结构体用于保存每个任务的监控信息
type TaskMonitoring struct {
	TaskID        string
	StartTime     time.Time
	EndTime       time.Time
	Status        string
	ResourceUsage ResourceUsage
}

// ResourceUsage 结构体用于记录资源使用情况
type ResourceUsage struct {
	CPUUsage    float64
	MemoryUsage uint64
	DiskUsage   uint64
}

// TaskMonitor 接口提供监控任务生命周期的基本方法
type TaskMonitor interface {
	StartMonitoring(taskID string) error
	StopMonitoring(taskID string) error
	GetTaskStatus(taskID string) (TaskMonitoring, error)
}

// SimpleTaskMonitor 是 TaskMonitor 接口的一个简单实现
type SimpleTaskMonitor struct {
	tasks map[string]*TaskMonitoring
}

// NewSimpleTaskMonitor 创建一个新的 SimpleTaskMonitor 实例
func NewSimpleTaskMonitor() *SimpleTaskMonitor {
	return &SimpleTaskMonitor{
		tasks: make(map[string]*TaskMonitoring),
	}
}

// StartMonitoring 开始监控某个任务
func (m *SimpleTaskMonitor) StartMonitoring(taskID string) error {
	m.tasks[taskID] = &TaskMonitoring{
		TaskID:    taskID,
		StartTime: time.Now(),
		Status:    "Running",
	}
	return nil
}

// StopMonitoring 停止监控某个任务
func (m *SimpleTaskMonitor) StopMonitoring(taskID string) error {
	if task, exists := m.tasks[taskID]; exists {
		task.EndTime = time.Now()
		task.Status = "Completed"
		return nil
	}
	return fmt.Errorf("task %s not found", taskID)
}

// GetTaskStatus 获取任务状态及资源消耗信息
func (m *SimpleTaskMonitor) GetTaskStatus(taskID string) (TaskMonitoring, error) {
	if task, exists := m.tasks[taskID]; exists {
		return *task, nil
	}
	return TaskMonitoring{}, fmt.Errorf("task %s not found", taskID)
}

// TaskLogger 用于记录任务执行信息
type TaskLogger struct {
	LogFilePath string
}

// TaskLog 结构体用于保存任务日志信息
type TaskLog struct {
	TaskID    string
	StartTime time.Time
	EndTime   time.Time
	Output    string
	Error     string
}

// Logger 接口提供基本的日志记录方法
type Logger interface {
	LogTaskStart(taskID string, startTime time.Time) error
	LogTaskEnd(taskID string, endTime time.Time, output string, err error) error
	FetchLogs(taskID string) ([]TaskLog, error)
}

// SimpleLogger 是 Logger 接口的一个简单实现
type SimpleLogger struct {
	logs map[string][]TaskLog
}

// NewSimpleLogger 创建一个新的 SimpleLogger 实例
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		logs: make(map[string][]TaskLog),
	}
}

// LogTaskStart 记录任务开始
func (l *SimpleLogger) LogTaskStart(taskID string, startTime time.Time) error {
	l.logs[taskID] = append(l.logs[taskID], TaskLog{
		TaskID:    taskID,
		StartTime: startTime,
	})
	return nil
}

// LogTaskEnd 记录任务结束
func (l *SimpleLogger) LogTaskEnd(taskID string, endTime time.Time, output string, err error) error {
	if len(l.logs[taskID]) > 0 {
		lastLog := &l.logs[taskID][len(l.logs[taskID])-1]
		lastLog.EndTime = endTime
		lastLog.Output = output
		if err != nil {
			lastLog.Error = err.Error()
		}
	}
	return nil
}

// FetchLogs 获取特定任务的日志
func (l *SimpleLogger) FetchLogs(taskID string) ([]TaskLog, error) {
	if logs, exists := l.logs[taskID]; exists {
		return logs, nil
	}
	return nil, fmt.Errorf("logs for task %s not found", taskID)
}

// ErrorHandler 负责处理任务的错误情况
type ErrorHandler struct {
	MaxRetryCount int
	RetryInterval time.Duration
}

// ErrorHandling 接口提供错误处理和重试的方法
type ErrorHandling interface {
	CaptureError(taskID string, err error) error
	RetryTask(taskID string) error
}

// SimpleErrorHandler 是 ErrorHandling 接口的一个简单实现
type SimpleErrorHandler struct {
	ErrorHandler
	retryCount map[string]int
}

// NewSimpleErrorHandler 创建一个新的 SimpleErrorHandler 实例
func NewSimpleErrorHandler(maxRetryCount int, retryInterval time.Duration) *SimpleErrorHandler {
	return &SimpleErrorHandler{
		ErrorHandler: ErrorHandler{
			MaxRetryCount: maxRetryCount,
			RetryInterval: retryInterval,
		},
		retryCount: make(map[string]int),
	}
}

// CaptureError 捕获任务执行中的异常
func (h *SimpleErrorHandler) CaptureError(taskID string, err error) error {
	fmt.Printf("Error captured for task %s: %v\n", taskID, err)
	return nil
}

// RetryTask 重试任务
func (h *SimpleErrorHandler) RetryTask(taskID string) error {
	if h.retryCount[taskID] < h.MaxRetryCount {
		h.retryCount[taskID]++
		fmt.Printf("Retrying task %s (attempt %d/%d)\n", taskID, h.retryCount[taskID], h.MaxRetryCount)
		time.Sleep(h.RetryInterval)
		return nil
	}
	return fmt.Errorf("max retry count reached for task %s", taskID)
}

// TaskRecovery 用于保存和恢复任务的状态
type TaskRecovery struct {
	TaskID    string
	LastState string
}

// Recovery 接口提供任务状态的保存和恢复方法
type Recovery interface {
	SaveTaskState(taskID string, state string) error
	RecoverTaskState(taskID string) (string, error)
}

// SimpleRecovery 是 Recovery 接口的一个简单实现
type SimpleRecovery struct {
	states map[string]string
}

// NewSimpleRecovery 创建一个新的 SimpleRecovery 实例
func NewSimpleRecovery() *SimpleRecovery {
	return &SimpleRecovery{
		states: make(map[string]string),
	}
}

// SaveTaskState 保存任务当前状态
func (r *SimpleRecovery) SaveTaskState(taskID string, state string) error {
	r.states[taskID] = state
	return nil
}

// RecoverTaskState 恢复任务之前的状态
func (r *SimpleRecovery) RecoverTaskState(taskID string) (string, error) {
	if state, exists := r.states[taskID]; exists {
		return state, nil
	}
	return "", fmt.Errorf("state for task %s not found", taskID)
}

// TimeoutController 负责任务的超时控制
type TimeoutController struct {
	TimeoutDuration time.Duration
}

// TimeoutControl 接口提供超时控制的方法
type TimeoutControl interface {
	SetTaskTimeout(taskID string, duration time.Duration) error
	CheckTimeout(taskID string) bool
	HandleTimeout(taskID string) error
}

// SimpleTimeoutController 是 TimeoutControl 接口的一个简单实现
type SimpleTimeoutController struct {
	TimeoutController
	taskStartTimes map[string]time.Time
}

// NewSimpleTimeoutController 创建一个新的 SimpleTimeoutController 实例
func NewSimpleTimeoutController(defaultTimeout time.Duration) *SimpleTimeoutController {
	return &SimpleTimeoutController{
		TimeoutController: TimeoutController{
			TimeoutDuration: defaultTimeout,
		},
		taskStartTimes: make(map[string]time.Time),
	}
}

// SetTaskTimeout 设置任务的超时时间
func (t *SimpleTimeoutController) SetTaskTimeout(taskID string, duration time.Duration) error {
	t.taskStartTimes[taskID] = time.Now()
	return nil
}

// CheckTimeout 检查任务是否已超时
func (t *SimpleTimeoutController) CheckTimeout(taskID string) bool {
	if startTime, exists := t.taskStartTimes[taskID]; exists {
		return time.Since(startTime) > t.TimeoutDuration
	}
	return false
}

// HandleTimeout 处理超时任务
func (t *SimpleTimeoutController) HandleTimeout(taskID string) error {
	fmt.Printf("Task %s has timed out\n", taskID)
	delete(t.taskStartTimes, taskID)
	return nil
}