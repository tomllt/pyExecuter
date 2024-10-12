package pyExecuter

import (
	"fmt"
)

// TaskRecovery 保存与恢复任务状态
type TaskRecovery struct {
	taskStates map[string]string // 记录任务的状态
}

// Recovery 任务恢复接口
type Recovery interface {
	SaveTaskState(taskID string, state string) error      // 保存任务当前状态
	RecoverTaskState(taskID string) (string, error)       // 恢复任务之前的状态
}

// NewTaskRecovery 创建 TaskRecovery 实例
func NewTaskRecovery() *TaskRecovery {
	return &TaskRecovery{
		taskStates: make(map[string]string),
	}
}

// SaveTaskState 保存任务状态
func (r *TaskRecovery) SaveTaskState(taskID string, state string) error {
	r.taskStates[taskID] = state
	fmt.Printf("Task %s state saved: %s\n", taskID, state)
	return nil
}

// RecoverTaskState 恢复任务状态
func (r *TaskRecovery) RecoverTaskState(taskID string) (string, error) {
	state, exists := r.taskStates[taskID]
	if !exists {
		return "", fmt.Errorf("no state found for task %s", taskID)
	}
	fmt.Printf("Recovered task %s to state: %s\n", taskID, state)
	return state, nil
}