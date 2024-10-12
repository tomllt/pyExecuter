package pyExecuter

import (
	"fmt"
	"sync"
	"time"
)

// TaskMonitoring 任务监控结构体
type TaskMonitoring struct {
	TaskID        string        // 任务唯一标识
	StartTime     time.Time     // 任务开始时间
	EndTime       time.Time     // 任务结束时间
	Status        string        // 任务当前状态
	ResourceUsage ResourceUsage // 资源使用情况
}

// ResourceUsage 资源使用情况
type ResourceUsage struct {
	CPUUsage     float64 // CPU 使用率
	MemoryUsage  uint64  // 内存使用量
	DiskUsage    uint64  // 磁盘使用量
	NetworkUsage uint64  // 网络使用量
}

// TaskMonitor 任务监控接口
type TaskMonitor interface {
	StartMonitoring(taskID string) error                     // 开始监控某个任务
	StopMonitoring(taskID string) error                      // 停止监控某个任务
	GetTaskStatus(taskID string) (*TaskMonitoring, error)    // 获取任务状态及资源消耗信息
	UpdateTaskStatus(taskID string, status string, usage ResourceUsage) error // 更新任务状态和资源使用情况
}

// BasicTaskMonitor 任务监控的简单实现
type BasicTaskMonitor struct {
	monitorData map[string]*TaskMonitoring
	mu          sync.RWMutex
	updateTicker *time.Ticker
}

// NewBasicTaskMonitor 创建 BasicTaskMonitor 实例
func NewBasicTaskMonitor() *BasicTaskMonitor {
	monitor := &BasicTaskMonitor{
		monitorData:  make(map[string]*TaskMonitoring),
		updateTicker: time.NewTicker(5 * time.Second), // 每5秒更新一次
	}
	go monitor.periodicallyUpdateResourceUsage()
	return monitor
}

// StartMonitoring 实现任务开始监控
func (m *BasicTaskMonitor) StartMonitoring(taskID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.monitorData[taskID]; exists {
		return fmt.Errorf("task %s is already being monitored", taskID)
	}
	m.monitorData[taskID] = &TaskMonitoring{
		TaskID:    taskID,
		StartTime: time.Now(),
		Status:    "Running",
		ResourceUsage: ResourceUsage{}, // 初始化资源使用情况
	}
	return nil
}

// StopMonitoring 实现停止监控任务
func (m *BasicTaskMonitor) StopMonitoring(taskID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	monitor, exists := m.monitorData[taskID]
	if !exists {
		return fmt.Errorf("task %s is not being monitored", taskID)
	}
	monitor.EndTime = time.Now()
	monitor.Status = "Completed"
	return nil
}

// GetTaskStatus 获取任务的监控状态
func (m *BasicTaskMonitor) GetTaskStatus(taskID string) (*TaskMonitoring, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	monitor, exists := m.monitorData[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s is not being monitored", taskID)
	}
	return monitor, nil
}

// UpdateTaskStatus 更新任务状态和资源使用情况
func (m *BasicTaskMonitor) UpdateTaskStatus(taskID string, status string, usage ResourceUsage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	monitor, exists := m.monitorData[taskID]
	if !exists {
		return fmt.Errorf("task %s is not being monitored", taskID)
	}
	monitor.Status = status
	monitor.ResourceUsage = usage
	return nil
}

// periodicallyUpdateResourceUsage 定期更新所有任务的资源使用情况
func (m *BasicTaskMonitor) periodicallyUpdateResourceUsage() {
	for range m.updateTicker.C {
		m.mu.RLock()
		for taskID := range m.monitorData {
			// 这里应该实现实际的资源使用情况更新逻辑
			// 例如，调用系统API获取CPU、内存、磁盘和网络使用情况
			usage := ResourceUsage{
				CPUUsage:     0.5,  // 示例值
				MemoryUsage:  1024, // 示例值
				DiskUsage:    2048, // 示例值
				NetworkUsage: 512,  // 示例值
			}
			m.UpdateTaskStatus(taskID, m.monitorData[taskID].Status, usage)
		}
		m.mu.RUnlock()
	}
}