package mysqldb

import (
	"errors"
	"fmt"
)

// 自定义错误类型
type DatabaseNotFoundError struct {
	Database string
	Mode     AccessMode
}

func (e *DatabaseNotFoundError) Error() string {
	return fmt.Sprintf("未找到数据库配置 (请求: %s/%s)",
		e.Database, e.Mode)
}

var (
	ErrConfigNotFound = errors.New("配置文件未找到")
	ErrInvalidConfig  = errors.New("无效的配置文件格式")
)

// NewConnectionError 创建带上下文的连接错误
func NewConnectionError(host string, port int, user string, err error) error {
	return fmt.Errorf("mysql连接失败 [host=%s port=%d user=%s]: %w",
		host, port, user, err)
}

// NewPingError 创建探活错误
func NewPingError(host string, port int, err error) error {
	return fmt.Errorf("mysql探活失败 [host=%s port=%d]: %w",
		host, port, err)
}

// NewSQLDBError 创建底层连接错误
func NewSQLDBError(host string, port int, err error) error {
	return fmt.Errorf("mysql获取底层连接失败 [host=%s port=%d]: %w",
		host, port, err)
}
