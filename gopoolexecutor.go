package pyExecuter

import (
	"context"
	"fmt"
	"sync"
	"time"

	gopool "github.com/devchat-ai/gopool"
)

// Task 描述一个需要执行的Python脚本任务
type Task struct {
	ID           string              // 任务的唯一ID
	Script       string              // Python脚本代码（字符串形式）
	Args         []string            // 脚本执行的参数
	Priority     int                 // 任务的优先级（可选）
	Timeout      time.Duration       // 任务超时时间
	RetryCount   int                 // 重试次数
	OnCompletion func(result Result) // 任务完成后的回调函数
}

// Result 描述任务执行的结果
type Result struct {
	TaskID    string    // 对应任务的ID
	Output    string    // 执行的输出结果
	Error     error     // 执行过程中产生的错误
	StartTime time.Time // 任务开始时间
	EndTime   time.Time // 任务结束时间
}

// GopoolExecutor GoPool 的任务执行管理器
type GopoolExecutor struct {
	pool  *gopool.Pool // 使用 devchat-ai/gopool 提供的池
	Queue *TaskQueue   // 任务队列
	mu    sync.Mutex   // 保护任务调度的锁
}

// NewGopoolExecutor 创建一个GopoolExecutor实例
func NewGopoolExecutor(poolSize int, queue *TaskQueue) *GopoolExecutor {
	pool := gopool.NewGoPool(poolSize)
	return &GopoolExecutor{
		pool:  pool,
		Queue: queue,
	}
}

// Start 启动GopoolExecutor，持续从任务队列获取任务并执行
func (e *GopoolExecutor) Start(ctx context.Context) error {
	// 利用 GoPool 并行执行任务，从任务队列获取任务并提交
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				e.mu.Lock()
				task, err := e.Queue.GetTask() // 获取任务
				e.mu.Unlock()
				if err == nil && task != nil {
					e.pool.AddTask(func() (interface{}, error) {
						result := e.ExecuteTask(task)
						if result.Error != nil {
							if task.RetryCount > 0 {
								task.RetryCount--
								e.Queue.AddTask(*task) // 任务失败，重新添加到队列
							} else {
								// 记录失败日志
								fmt.Printf("Task %s failed after retries: %v\n", task.ID, result.Error)
							}
						}
						return nil, result.Error
					})
				}
			}
		}
	}()
	return nil
}

// ExecuteTask 执行单个任务（内部方法）
func (e *GopoolExecutor) ExecuteTask(task *Task) Result {
	result := Result{
		TaskID:    task.ID,
		StartTime: time.Now(),
	}

	executor := &SecurePythonExecutor{}
	executor.SetupEnvironment("task_env") // 设置虚拟环境

	output, err := executor.Execute(task.Script, task.Args, task.Timeout)

	result.EndTime = time.Now()
	result.Output = output
	result.Error = err

	if task.OnCompletion != nil {
		task.OnCompletion(result)
	}

	return result
}
