package pyExecuter_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tomllt/pyExecuter"
)

func TestGopoolExecutor(t *testing.T) {
	// 创建任务队列
	queue := pyExecuter.NewTaskQueue(100, "FIFO")

	// 创建执行器
	executor := pyExecuter.NewGopoolExecutor(5, queue)

	// 启动执行器
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := executor.Start(ctx)
	assert.NoError(t, err)

	// 创建测试任务
	task := &pyExecuter.Task{
		ID:      "test_task",
		Script:  "print('Hello, World!')",
		Timeout: 5 * time.Second,
	}

	// 添加任务到队列
	err = queue.AddTask(task)
	assert.NoError(t, err)

	// 等待任务执行完成
	time.Sleep(2 * time.Second)

	// 检查执行器状态
	stats := executor.GetStats()
	assert.Equal(t, 0, stats["queue_size"])
	assert.LessOrEqual(t, stats["running_workers"], 5)

	// 测试任务执行
	result := executor.ExecuteTask(task)
	assert.NoError(t, result.Error)
	assert.Contains(t, result.Output, "Hello, World!")
}

func TestPythonExecutor(t *testing.T) {
	executor := &pyExecuter.SecurePythonExecutor{}

	// 设置环境
	err := executor.SetupEnvironment("test_env")
	assert.NoError(t, err)

	// 执行 Python 脚本
	script := `
import sys
print(f"Python version: {sys.version}")
print("Hello from Python!")
`
	output, err := executor.Execute(script, []string{}, 5*time.Second)
	assert.NoError(t, err)
	assert.Contains(t, output, "Python version:")
	assert.Contains(t, output, "Hello from Python!")
}

func TestTaskQueue(t *testing.T) {
	queue := pyExecuter.NewTaskQueue(10, "FIFO")

	task1 := &pyExecuter.Task{ID: "task1", Priority: 1}
	task2 := &pyExecuter.Task{ID: "task2", Priority: 2}

	err := queue.AddTask(task1)
	assert.NoError(t, err)
	err = queue.AddTask(task2)
	assert.NoError(t, err)

	assert.Equal(t, 2, queue.Size())

	retrievedTask, err := queue.GetTask()
	assert.NoError(t, err)
	assert.Equal(t, "task2", retrievedTask.ID) // Higher priority task should be retrieved first

	retrievedTask, err = queue.GetTask()
	assert.NoError(t, err)
	assert.Equal(t, "task1", retrievedTask.ID)

	_, err = queue.GetTask()
	assert.Error(t, err) // Queue should be empty now
}

func TestTaskRecovery(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	
	recovery := pyExecuter.NewTaskRecovery(tempDir)

	// 测试保存任务状态
	err := recovery.SaveTaskState("task1", "running")
	assert.NoError(t, err)

	// 测试恢复任务状态
	state, err := recovery.RecoverTaskState("task1")
	assert.NoError(t, err)
	assert.Equal(t, "running", state.State)

	// 测试持久化状态
	err = recovery.PersistStates()
	assert.NoError(t, err)

	// 创建新的 TaskRecovery 实例来测试加载状态
	newRecovery := pyExecuter.NewTaskRecovery(tempDir)
	err = newRecovery.LoadStates()
	assert.NoError(t, err)

	// 验证加载的状态
	state, err = newRecovery.RecoverTaskState("task1")
	assert.NoError(t, err)
	assert.Equal(t, "running", state.State)

	// 测试恢复不存在的任务状态
	_, err = newRecovery.RecoverTaskState("non_existent_task")
	assert.Error(t, err)

	// 测试更新任务状态
	err = newRecovery.SaveTaskState("task1", "completed")
	assert.NoError(t, err)

	state, err = newRecovery.RecoverTaskState("task1")
	assert.NoError(t, err)
	assert.Equal(t, "completed", state.State)
}

func TestTaskMonitor(t *testing.T) {
	monitor := pyExecuter.NewBasicTaskMonitor()

	err := monitor.StartMonitoring("task1")
	assert.NoError(t, err)

	time.Sleep(1 * time.Second) // Wait for a bit to simulate task running

	status, err := monitor.GetTaskStatus("task1")
	assert.NoError(t, err)
	assert.Equal(t, "Running", status.Status)

	err = monitor.StopMonitoring("task1")
	assert.NoError(t, err)

	status, err = monitor.GetTaskStatus("task1")
	assert.NoError(t, err)
	assert.Equal(t, "Completed", status.Status)
}

func TestErrorHandler(t *testing.T) {
	queue := pyExecuter.NewTaskQueue(10, "FIFO")
	handler := pyExecuter.NewBasicErrorHandler(3, 1*time.Second, queue)

	task := &pyExecuter.Task{ID: "errorTask"}
	err := queue.AddTask(task)
	assert.NoError(t, err)

	err = handler.CaptureError("errorTask", assert.AnError)
	assert.NoError(t, err) // First retry

	err = handler.CaptureError("errorTask", assert.AnError)
	assert.NoError(t, err) // Second retry

	err = handler.CaptureError("errorTask", assert.AnError)
	assert.NoError(t, err) // Third retry

	err = handler.CaptureError("errorTask", assert.AnError)
	assert.Error(t, err) // Should exceed max retry count
}
