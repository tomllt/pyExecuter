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

	t.Run("SetTaskTimeout and CheckTimeout", func(t *testing.T) {
		taskID := "timeout_test_task"
		timeoutDuration := 100 * time.Millisecond
		err := timeoutControl.SetTaskTimeout(taskID, timeoutDuration)
		assert.NoError(t, err)

		// 模拟任务运行 (未超时)
		time.Sleep(50 * time.Millisecond)
		isTimeout, err := timeoutControl.CheckTimeout(taskID)
		assert.NoError(t, err)
		assert.False(t, isTimeout)

		// 模拟任务运行 (超时)
		time.Sleep(60 * time.Millisecond)
		isTimeout, err = timeoutControl.CheckTimeout(taskID)
		assert.NoError(t, err)
		assert.True(t, isTimeout)
	})

	t.Run("HandleTimeout", func(t *testing.T) {
		taskID := "timeout_handle_task"
		timeoutDuration := 50 * time.Millisecond
		err := timeoutControl.SetTaskTimeout(taskID, timeoutDuration)
		assert.NoError(t, err)

		time.Sleep(60 * time.Millisecond)
		err = timeoutControl.HandleTimeout(taskID)
		assert.NoError(t, err)
	})

	t.Run("ClearTimeout", func(t *testing.T) {
		taskID := "clear_timeout_task"
		timeoutDuration := 200 * time.Millisecond
		err := timeoutControl.SetTaskTimeout(taskID, timeoutDuration)
		assert.NoError(t, err)

		err = timeoutControl.ClearTimeout(taskID)
		assert.NoError(t, err)

		_, err = timeoutControl.CheckTimeout(taskID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no timeout set for task")
	})

	t.Run("NonExistentTask", func(t *testing.T) {
		_, err := timeoutControl.CheckTimeout("non_existent_task")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no timeout set for task")
	})

	t.Run("ExecuteWithTimeoutControl", func(t *testing.T) {
		task := &pyExecuter.Task{
			ID:      "execute_timeout_task",
			Script:  "import time; time.sleep(0.1); print('Hello, World!')",
			Timeout: 200 * time.Millisecond,
		}

		ctx := context.Background()
		err := pyExecuter.ExecuteWithTimeoutControl(ctx, task, timeoutControl)
		assert.NoError(t, err)

		// Test timeout scenario
		task.Timeout = 50 * time.Millisecond
		err = pyExecuter.ExecuteWithTimeoutControl(ctx, task, timeoutControl)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timed out")
	})
}
