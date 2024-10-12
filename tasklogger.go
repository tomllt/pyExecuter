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
logEntry := TaskLog{
    TaskID:    taskID,
    StartTime: startTime,
}
f.logs[taskID] = append(f.logs[taskID], logEntry)
if err := f.writeToFile(fmt.Sprintf("Task %s started at %s\n", taskID, startTime)); err != nil {
    return fmt.Errorf("failed to log task start: %v", err)
}
return nil
}

// LogTaskEnd 记录任务结束
func (f *FileLogger) LogTaskEnd(taskID string, endTime time.Time, output string, err error) error {
	taskLogs := f.logs[taskID]
taskLog := &taskLogs[len(taskLogs)-1]
taskLog.EndTime = endTime
taskLog.Output = output
if err != nil {
    taskLog.Error = err.Error()
}
f.logs[taskID] = taskLogs  // 确保更新后的日志写回到映射中
	logMessage := fmt.Sprintf("Task %s ended at %s with output: %s", taskID, endTime, output)
if err != nil {
    logMessage += fmt.Sprintf(" and error: %v", err)
}
if err := f.writeToFile(logMessage + "\n"); err != nil {
    return fmt.Errorf("failed to log task end: %v", err)
}
return nil
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