# Logger 模块

该模块为 `meteor` 项目提供日志功能。它支持 `glog` 和 `zap` 日志库。

## 安装

使用以下命令安装 logger 模块：

```sh
go get github.com/orbit-w/meteor/modules/net/transport/logger
```

## 使用方法

### 导入包

```go
import "github.com/orbit-w/meteor/modules/net/transport/logger"
```

### 创建 Logger

#### ZapLogger

`ZapLogger` 使用 `zap` 进行日志记录，并带有指定的前缀。

```go
zapLogger := logger.NewLogger("MyPrefix")

zapLogger.Info("这是一条信息日志")
zapLogger.Infof("这是一条格式化的信息日志: %s", "格式化内容")
zapLogger.Error("这是一条错误日志")
zapLogger.Errorf("这是一条格式化的错误日志: %s", "格式化内容")
```

### 配置

logger 模块使用 `viper` 从文件中读取配置。包含以下字段Key：
```yaml
v: debug        # 日志等级
log_dir: log_dir # 日志目录；如果为空，日志输出到 stdout
```

### 示例配置

```yaml
log:
  v: debug
  log_dir: /var/log/meteor
```

### 设置基础 Logger

可以使用 `SetBaseLogger` 函数设置自定义基础 logger。

```go
customLogger := zap.NewExample()
logger.SetBaseLogger(customLogger)
```

### 停止基础 Logger

要正确刷新并关闭基础 logger，请使用 `StopBaseLogger` 函数。

```go
logger.StopBaseLogger()
```

## 许可证

该项目使用 MIT 许可证。
```