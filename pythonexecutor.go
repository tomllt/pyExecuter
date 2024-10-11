package pyExecuter

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// PythonExecutor Python脚本执行的接口
type PythonExecutor interface {
	Execute(script string, args []string, timeout time.Duration) (string, error) // 执行Python脚本
	SetupEnvironment(envName string) error                                       // 设置Python虚拟环境
}

// SecurePythonExecutor 实现了PythonExecutor接口，具有虚拟环境管理和安全机制
type SecurePythonExecutor struct {
	Environment string
}

// SetupEnvironment 设置Python虚拟环境
func (p *SecurePythonExecutor) SetupEnvironment(envName string) error {
	// 创建虚拟环境
	cmd := exec.Command("python", "-m", "venv", envName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create virtual environment: %v", err)
	}

	p.Environment = envName
	return nil
}

// Execute 执行Python脚本，返回输出或者错误
func (p *SecurePythonExecutor) Execute(script string, args []string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 创建一个临时文件来存储脚本
	tmpFile, err := createTempPythonFile(script)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer removeTempPythonFile(tmpFile)

	// 准备命令
	pythonPath := filepath.Join(p.Environment, "bin", "python")
	cmdArgs := append([]string{tmpFile}, args...)
	cmd := exec.CommandContext(ctx, pythonPath, cmdArgs...)

	// 设置环境变量
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("VIRTUAL_ENV=%s", p.Environment),
		fmt.Sprintf("PATH=%s:%s", filepath.Join(p.Environment, "bin"), os.Getenv("PATH")),
	)

	// 设置输出缓冲
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// 执行命令
	err = cmd.Run()

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("execution timed out after %v", timeout)
		}
		return "", fmt.Errorf("execution failed: %v", err)
	}

	return out.String(), nil
}

// createTempPythonFile 创建一个临时的Python文件
func createTempPythonFile(script string) (string, error) {
	// 创建一个临时目录
	tempDir, err := os.MkdirTemp("", "python_script_")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// 在临时目录中创建一个Python文件
	tempFile := filepath.Join(tempDir, "script.py")
	err = os.WriteFile(tempFile, []byte(script), 0600)
	if err != nil {
		os.RemoveAll(tempDir) // 清理临时目录
		return "", fmt.Errorf("failed to write script to temp file: %v", err)
	}

	return tempFile, nil
}

// removeTempPythonFile 删除临时的Python文件
func removeTempPythonFile(filePath string) error {
	// 获取临时文件所在的目录
	tempDir := filepath.Dir(filePath)

	// 删除整个临时目录及其内容
	err := os.RemoveAll(tempDir)
	if err != nil {
		return fmt.Errorf("failed to remove temp directory: %v", err)
	}

	return nil
}
