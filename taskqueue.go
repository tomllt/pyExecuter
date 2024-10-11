package pyExecuter

import (
	"fmt"
	"sync"
)

// TaskQueue 任务队列的实现
type TaskQueue struct {
	tasks        []Task
	maxCapacity  int
	mu           sync.RWMutex
	priorityMode string // "FIFO" or "LIFO"
}

// NewTaskQueue 创建一个TaskQueue实例
func NewTaskQueue(maxCapacity int, priorityMode string) *TaskQueue {
	return &TaskQueue{
		tasks:        make([]Task, 0, maxCapacity),
		maxCapacity:  maxCapacity,
		priorityMode: priorityMode,
	}
}

// AddTask 添加任务到队列中
func (q *TaskQueue) AddTask(task Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tasks) >= q.maxCapacity {
		return fmt.Errorf("task queue is full")
	}

	q.tasks = append(q.tasks, task)
	return nil
}

// GetTask 获取一个任务
func (q *TaskQueue) GetTask() (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tasks) == 0 {
		return nil, fmt.Errorf("no tasks available")
	}

	var task Task
	if q.priorityMode == "FIFO" {
		task = q.tasks[0]
		q.tasks = q.tasks[1:]
	} else if q.priorityMode == "LIFO" {
		task = q.tasks[len(q.tasks)-1]
		q.tasks = q.tasks[:len(q.tasks)-1]
	}
	return &task, nil
}