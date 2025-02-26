package mysqldb

import "fmt"

// 自定义错误类型
type DatabaseNotFoundError struct {
	Database string
	Mode     AccessMode
}

func (e *DatabaseNotFoundError) Error() string {
	return fmt.Sprintf("未找到数据库配置 (请求: %s/%s)",
		e.Database, e.Mode)
}
