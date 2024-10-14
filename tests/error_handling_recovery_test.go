package pyExecuter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tomllt/pyExecuter"
)

func TestErrorHandling(t *testing.T) {
	queue := pyExecuter.NewTaskQueue(10, "FIFO")
	errorHandler := pyExecuter.NewBasicErrorHandler(3, 1*time.Second, queue)

	// 模拟任务执行错误
	taskID := "error_handling_test_task"
	err := errorHandler.CaptureError(taskID, fmt.Errorf("sample error"))
	assert.NoError(t, err)

	// 验证重试次数
	for i := 0; i < 3; i++ {
		err := errorHandler.CaptureError(taskID, fmt.Errorf("retry error %d", i+1))
		assert.NoError(t, err)
	}

	// 验证超过最大重试次数
	err = errorHandler.CaptureError(taskID, fmt.Errorf("final error"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeded max retry count")
}

func TestTaskRecovery(t *testing.T) {
	recovery := pyExecuter.NewTaskRecovery()

	// 保存任务状态
	taskID := "recovery_test_task"
	state := "InProgress"
	err := recovery.SaveTaskState(taskID, state)
	assert.NoError(t, err)

	// 恢复任务状态
	recoveredState, err := recovery.RecoverTaskState(taskID)
	assert.NoError(t, err)
	assert.Equal(t, state, recoveredState)
}
