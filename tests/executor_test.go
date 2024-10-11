package pyExecuter_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/pyExecuter"
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
		ID:     "test_task",
		Script: "print('Hello, World!')",
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