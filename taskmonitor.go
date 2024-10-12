package pyExecuter

import (
	"fmt"
	"time"
)

// TaskMonitoring 任务监控结构体
type TaskMonitoring struct {
	TaskID        string        // 任务唯一标识
	StartTime     time.Time     // 任务开始时间
	EndTime       time.Time     // 任务结束时间
	Status        string        // 任务当前状态
	ResourceUsage ResourceUsage // 资源使用情况
	MonitorActive bool          // 标识是否正在监控任务
}

// ResourceUsage 资源使用情况
type ResourceUsage struct {
	CPUUsage    float64 // CPU 使用率
	MemoryUsage uint64  // 内存使用量
	DiskUsage   uint64  // 磁盘使用量
}

// TaskMonitor 任务监控接口
type TaskMonitor interface {
	StartMonitoring(taskID string) error                     // 开始监控某个任务
	StopMonitoring(taskID string) error                      // 停止监控某个任务
	GetTaskStatus(taskID string) (*TaskMonitoring, error)    // 获取任务状态及资源消耗信息
}

// BasicTaskMonitor 任务监控的简单实现
type BasicTaskMonitor struct {
	monitorData map[string]*TaskMonitoring
}

// NewBasicTaskMonitor 创建 BasicTaskMonitor 实例
func NewBasicTaskMonitor() *BasicTaskMonitor {
	return &BasicTaskMonitor{
		monitorData: make(map[string]*TaskMonitoring),
	}
}

// StartMonitoring 实现任务开始监控
func (m *BasicTaskMonitor) StartMonitoring(taskID string) error {
	m.mu.Lock() // 防止竞态条件
	defer m.mu.Unlock()
	if _, exists := m.monitorData[taskID]; exists {
		return fmt.Errorf("task %s is already being monitored", taskID)
	}
	m.monitorData[taskID] = &TaskMonitoring{
		TaskID:    taskID,
		StartTime: time.Now(),
		Status:    "Running",
		ResourceUsage: ResourceUsage{
			CPUUsage:    getCurrentCPUUsage(),
			MemoryUsage: getCurrentMemoryUsage(),
			DiskUsage:   getCurrentDiskUsage(),
		},
		MonitorActive: true,
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.monitorData, taskID) // 结束监控后清理数据
	return nil
}

// StopMonitoring 实现停止监控任务
func (m *BasicTaskMonitor) StopMonitoring(taskID string) error {
	m.mu.Lock() // 防止竞态条件
	defer m.mu.Unlock()
	monitor, exists := m.monitorData[taskID]
	if !exists {
		return fmt.Errorf("task %s is not being monitored", taskID)
	}
	monitor.EndTime = time.Now()
	if err := checkTaskFailure(taskID); err != nil {
		monitor.Status = "Failed"
	} else {
		monitor.Status = "Completed"
	}
	monitor.MonitorActive = false
	return nil
}

// GetTaskStatus 获取任务的监控状态
func (m *BasicTaskMonitor) GetTaskStatus(taskID string) (*TaskMonitoring, error) {
	monitor, exists := m.monitorData[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s is not being monitored", taskID)
	}
	return monitor, nil
}