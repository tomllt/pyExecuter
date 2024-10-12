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
}

// ResourceUsage 资源使用情况
type ResourceUsage struct {
	CPUUsage    float64 // CPU 使用率
	MemoryUsage uint64  // 内存使用量
	DiskUsage   uint64  // 磁盘使用量
}

// getResourceUsage 获取当前系统的资源使用情况
func getResourceUsage() ResourceUsage {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return ResourceUsage{
		CPUUsage:    getCPUUsage(),  // 需要实现的 CPU 使用率获取函数
		MemoryUsage: memStats.Alloc, // 已分配内存
		DiskUsage:   getDiskUsage(), // 需要实现的磁盘使用率获取函数
	}
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
	if _, exists := m.monitorData[taskID]; exists {
		return fmt.Errorf("task %s is already being monitored", taskID)
	}
	m.monitorData[taskID] = &TaskMonitoring{
		TaskID:        taskID,
		StartTime:     time.Now(),
		Status:        "Running",
		ResourceUsage: getResourceUsage(), // 获取初始的资源使用情况
	}
	go m.updateResourceUsage(taskID) // 启动监控更新资源使用情况
	mu.Lock() 
	defer mu.Unlock()
	return nil
}

// updateResourceUsage 动态更新任务的资源使用情况
func (m *BasicTaskMonitor) updateResourceUsage(taskID string) {
	for {
		time.Sleep(1 * time.Second) // 定期更新资源使用情况
		m.mu.Lock()
		monitor, exists := m.monitorData[taskID]
		if !exists || monitor.Status != "Running" {
			m.mu.Unlock()
			break
		}
		monitor.ResourceUsage = getResourceUsage() // 更新资源使用情况
		m.mu.Unlock()
	}
}

// StopMonitoring 实现停止监控任务
func (m *BasicTaskMonitor) StopMonitoring(taskID string) error {
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
	monitor, exists := m.monitorData[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s is not being monitored", taskID)
	}
	return monitor, nil
}