package pyExecuter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TaskState 表示任务的状态
type TaskState struct {
	State     string    // 任务的当前状态
	Timestamp time.Time // 状态更新的时间戳
}

// TaskRecovery 保存与恢复任务状态
type TaskRecovery struct {
	taskStates map[string]TaskState // 记录任务的状态
	mu         sync.RWMutex         // 读写锁，用于保护 taskStates
	storageDir string               // 状态持久化存储的目录
}

// Recovery 任务恢复接口
type Recovery interface {
	SaveTaskState(taskID string, state string) error    // 保存任务当前状态
	RecoverTaskState(taskID string) (TaskState, error)  // 恢复任务之前的状态
	PersistStates() error                               // 将所有状态持久化到磁盘
	LoadStates() error                                  // 从磁盘加载所有状态
}

// NewTaskRecovery 创建 TaskRecovery 实例
func NewTaskRecovery(storageDir string) *TaskRecovery {
	return &TaskRecovery{
		taskStates: make(map[string]TaskState),
		storageDir: storageDir,
	}
}

// SaveTaskState 保存任务状态
func (r *TaskRecovery) SaveTaskState(taskID string, state string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.taskStates[taskID] = TaskState{
		State:     state,
		Timestamp: time.Now(),
	}
	fmt.Printf("Task %s state saved: %s\n", taskID, state)
	return r.PersistStates()
}

// RecoverTaskState 恢复任务状态
func (r *TaskRecovery) RecoverTaskState(taskID string) (TaskState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	state, exists := r.taskStates[taskID]
	if !exists {
		return TaskState{}, fmt.Errorf("no state found for task %s", taskID)
	}
	fmt.Printf("Recovered task %s to state: %s\n", taskID, state.State)
	return state, nil
}

// PersistStates 将所有状态持久化到磁盘
func (r *TaskRecovery) PersistStates() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := json.Marshal(r.taskStates)
	if err != nil {
		return fmt.Errorf("failed to marshal task states: %v", err)
	}

	filePath := filepath.Join(r.storageDir, "task_states.json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write task states to file: %v", err)
	}

	return nil
}

// LoadStates 从磁盘加载所有状态
func (r *TaskRecovery) LoadStates() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := filepath.Join(r.storageDir, "task_states.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果文件不存在，不视为错误，只是返回空状态
			return nil
		}
		return fmt.Errorf("failed to read task states from file: %v", err)
	}

	if err := json.Unmarshal(data, &r.taskStates); err != nil {
		return fmt.Errorf("failed to unmarshal task states: %v", err)
	}

	return nil
}