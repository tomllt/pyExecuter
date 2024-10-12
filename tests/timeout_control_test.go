package pyExecuter_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tomllt/pyExecuter"
)

func TestTimeoutControl(t *testing.T) {
	timeoutControl := pyExecuter.NewTimeoutController()

	// 设置任务超时时间
	taskID := "timeout_test_task"
	timeoutDuration := 2 * time.Second
	err := timeoutControl.SetTaskTimeout(taskID, timeoutDuration)
	assert.NoError(t, err)

	// 模拟任务运行
	startTime := time.Now()
	time.Sleep(3 * time.Second)

	// 检查任务是否超时
	isTimeout := timeoutControl.CheckTimeout(taskID, startTime)
	assert.True(t, isTimeout)

	// 处理超时任务
	err = timeoutControl.HandleTimeout(taskID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timed out")
}
