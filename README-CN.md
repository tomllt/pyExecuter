# pyExecuter

pyExecuter 是一个用 Go 语言编写的工具，用于执行 Python 代码并返回结果。

## 功能

- 执行 Python 代码片段
- 返回执行结果或错误信息

## 安装

确保您的系统中已安装 Go 和 Python。然后，您可以通过以下命令安装 pyExecuter：

```
go get github.com/yourusername/pyExecuter
```

## 使用方法

在您的 Go 代码中引入 pyExecuter：

```go
import "github.com/yourusername/pyExecuter"
```

然后，您可以使用 `ExecutePython` 函数来执行 Python 代码：

```go
result, err := pyExecuter.ExecutePython("print('Hello, World!')")
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

## 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何参与项目开发。

## 许可证

本项目采用 MIT 许可证。详情请见 [LICENSE](LICENSE) 文件。