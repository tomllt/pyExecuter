package pyExecuter

import (
	"fmt"
	"os"
	"time"
)

// TaskLog 任务日志记录结构体
type TaskLog struct {
	TaskID    string
	StartTime time.Time
	EndTime   time.Time
	Output    string
	Error     string
}

// Logger 日志记录接口
type Logger interface {
	LogTaskStart(taskID string, startTime time.Time) error    // 记录任务开始
	LogTaskEnd(taskID string, endTime time.Time, output string, err error) error  // 记录任务结束
	FetchLogs(taskID string) ([]TaskLog, error)  // 获取特定任务的日志
}

// FileLogger 基于文件的任务日志记录器
type FileLogger struct {
	LogFilePath string   // 日志文件路径
	logs        map[string][]TaskLog
}

// NewFileLogger 创建 FileLogger 实例
func NewFileLogger(logFilePath string) *FileLogger {
	return &FileLogger{
		LogFilePath: logFilePath,
		logs:        make(map[string][]TaskLog),
	}
}

// LogTaskStart 记录任务开始
func (f *FileLogger) LogTaskStart(taskID string, startTime time.Time) error {
	f.logs[taskID] = append(f.logs[taskID], TaskLog{
		TaskID:    taskID,
		StartTime: startTime,
	})
	return f.writeToFile(fmt.Sprintf("Task %s started at %s\n", taskID, startTime))
}

// LogTaskEnd 记录任务结束
func (f *FileLogger) LogTaskEnd(taskID string, endTime time.Time, output string, err error) error {
	taskLogs := f.logs[taskID]
	taskLogs[len(taskLogs)-1].EndTime = endTime
	taskLogs[len(taskLogs)-1].Output = output
	if err != nil {
		taskLogs[len(taskLogs)-1].Error = err.Error()
	}
	return f.writeToFile(fmt.Sprintf("Task %s ended at %s with output: %s and error: %v\n", taskID, endTime, output, err))
}

// FetchLogs 获取指定任务的日志
func (f *FileLogger) FetchLogs(taskID string) ([]TaskLog, error) {
	taskLogs, exists := f.logs[taskID]
	if !exists {
		return nil, fmt.Errorf("no logs found for task %s", taskID)
	}
	return taskLogs, nil
}

// writeToFile 将日志写入文件
func (f *FileLogger) writeToFile(logEntry string) error {
	file, err := os.OpenFile(f.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(logEntry); err != nil {
		return fmt.Errorf("failed to write log entry: %v", err)
	}
	return nil
}