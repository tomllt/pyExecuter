package pyExecuter

import (
	"fmt"
	"sort"
	"sync"
)

// TaskQueue 任务队列的实现
type TaskQueue struct {
	tasks        []*Task
	maxCapacity  int
	mu           sync.RWMutex
	priorityMode string // "FIFO" or "LIFO"
}

// NewTaskQueue 创建一个TaskQueue实例
func NewTaskQueue(maxCapacity int, priorityMode string) *TaskQueue {
	if priorityMode != "FIFO" && priorityMode != "LIFO" {
		priorityMode = "FIFO" // 默认使用FIFO
	}
	return &TaskQueue{
		tasks:        make([]*Task, 0, maxCapacity),
		maxCapacity:  maxCapacity,
		priorityMode: priorityMode,
	}
}

// AddTask 添加任务到队列中
func (q *TaskQueue) AddTask(task *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tasks) >= q.maxCapacity {
		return fmt.Errorf("task queue is full")
	}

	// 根据优先级插入任务
	index := sort.Search(len(q.tasks), func(i int) bool {
		return q.tasks[i].Priority <= task.Priority
	})
	q.tasks = append(q.tasks, nil)
	copy(q.tasks[index+1:], q.tasks[index:])
	q.tasks[index] = task
	return nil
}

// GetTask 获取一个任务
func (q *TaskQueue) GetTask() (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tasks) == 0 {
		return nil, fmt.Errorf("no tasks available")
	}

	var task *Task
	if q.priorityMode == "FIFO" {
		task = q.tasks[0]
		q.tasks = q.tasks[1:]
	} else if q.priorityMode == "LIFO" {
		task = q.tasks[len(q.tasks)-1]
		q.tasks = q.tasks[:len(q.tasks)-1]
	}
	return task, nil
}

// Size 返回队列中的任务数量
func (q *TaskQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.tasks)
}

func (q *TaskQueue) GetTaskByID(taskID string) (*Task, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	for _, task := range q.tasks {
		if task.ID == taskID {
			return task, nil
		}
	}
	return nil, fmt.Errorf("task with ID %s not found", taskID)
}