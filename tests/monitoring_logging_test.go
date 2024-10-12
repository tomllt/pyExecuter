package pyExecuter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tomllt/pyExecuter"
)

func TestTaskMonitoring(t *testing.T) {
	monitor := pyExecuter.NewBasicTaskMonitor()

	// 开始监控任务
	taskID := "monitor_test_task"
	err := monitor.StartMonitoring(taskID)
	assert.NoError(t, err)

	// 模拟任务运行
	time.Sleep(1 * time.Second)

	// 检查任务状态
	status, err := monitor.GetTaskStatus(taskID)
	assert.NoError(t, err)
	assert.Equal(t, "Running", status.Status)

	// 停止监控任务
	err = monitor.StopMonitoring(taskID)
	assert.NoError(t, err)

	// 检查任务状态
	status, err = monitor.GetTaskStatus(taskID)
	assert.NoError(t, err)
	assert.Equal(t, "Completed", status.Status)
}

func TestTaskLogging(t *testing.T) {
	logger := pyExecuter.NewFileLogger("/tmp/task_log_test.txt")

	// 记录任务开始
	taskID := "log_test_task"
	startTime := time.Now()
	err := logger.LogTaskStart(taskID, startTime)
	assert.NoError(t, err)

	// 模拟任务执行并记录结束
	endTime := time.Now().Add(1 * time.Second)
	output := "Task completed successfully"
	err = logger.LogTaskEnd(taskID, endTime, output, nil)
	assert.NoError(t, err)

	// 获取并验证日志
	logs, err := logger.FetchLogs(taskID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(logs))
	assert.Contains(t, logs[0].Output, "Task completed successfully")
}
