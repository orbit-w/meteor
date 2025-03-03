package mysqldb

type (
	// AccessMode 访问模式
	AccessMode string

	// PoolConfig 连接池配置
	PoolConfig struct {
		MaxIdleConns int `yaml:"max_idle_conns" toml:"max_idle_conns"` // 最大空闲连接数
		MaxOpenConns int `yaml:"max_open_conns" toml:"max_open_conns"` // 最大打开连接数
	}

	// LogLevel 日志级别
	LogLevel string

	// LogConfig 日志配置
	LogConfig struct {
		Level string `yaml:"level" toml:"level"` // 日志级别：silent, error, warn, info
	}

	// InstanceConfig MySQL实例配置
	InstanceConfig struct {
		Host     string     `yaml:"host" toml:"host"`         // MySQL服务器地址
		Port     int        `yaml:"port" toml:"port"`         // 端口号
		Username string     `yaml:"username" toml:"username"` // 数据库用户名
		Password string     `yaml:"password" toml:"password"` // 数据库密码
		Mode     AccessMode `yaml:"mode" toml:"mode"`         // 访问模式
		Pool     PoolConfig `yaml:"pool" toml:"pool"`         // 连接池配置
		Log      LogConfig  `yaml:"log" toml:"log"`           // 日志配置
	}

	// DatabaseConfig 数据库配置
	DatabaseConfig struct {
		Name string `yaml:"name" toml:"name"` // 数据库名称
	}

	// ManagerConfig 集中式连接配置
	ManagerConfig struct {
		Instances []struct {
			Config    InstanceConfig   `yaml:"config" toml:"config"`       // 实例配置
			Databases []DatabaseConfig `yaml:"databases" toml:"databases"` // 数据库配置
		} `yaml:"instances" toml:"instances"`
	}
)

const (
	ReadOnly  AccessMode = "readonly"  // 只读模式
	ReadWrite AccessMode = "readwrite" // 读写模式
)

// DefaultPoolConfig 返回默认的连接池配置
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxIdleConns: 10,  // 默认最大空闲连接数
		MaxOpenConns: 100, // 默认最大打开连接数
	}
}

// DefaultInstanceConfig 返回默认的实例配置
func DefaultInstanceConfig() InstanceConfig {
	return InstanceConfig{
		Host:     "localhost", // 默认主机地址
		Port:     3306,        // 默认MySQL端口
		Username: "root",      // 默认用户名
		Password: "",          // 默认密码为空
		Pool:     DefaultPoolConfig(),
		Mode:     ReadOnly, // 默认为只读模式
		Log: LogConfig{
			Level: "silent", // 默认为静默模式
		},
	}
}

// DefaultManagerConfig 返回默认的管理器配置
func DefaultManagerConfig() ManagerConfig {
	return ManagerConfig{
		Instances: []struct {
			Config    InstanceConfig   `yaml:"config" toml:"config"`
			Databases []DatabaseConfig `yaml:"databases" toml:"databases"`
		}{
			{
				Config: DefaultInstanceConfig(),
				Databases: []DatabaseConfig{
					{
						Name: "test",
					},
					{
						Name: "test",
					},
				},
			},
		},
	}
}
