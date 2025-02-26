package mysqldb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseKey 数据库键
type DatabaseKey struct {
	Database string
	Mode     AccessMode
}

// String 返回数据库键的字符串表示
func (k DatabaseKey) String() string {
	return fmt.Sprintf("%s:%s", k.Database, k.Mode)
}

// GormInstanceMgr 实现数据库连接管理器
// GormInstanceMgr 管理 MySQL 数据库连接实例
// 特性：
// - 支持读写分离
// - 自动连接池管理
// - 连接健康检查
// - 线程安全
// 使用示例：
// mgr := New(config)
// db, err := mgr.DB("mydb", ReadWrite)
type GormInstanceMgr struct {
	config    ManagerConfig
	instances sync.Map // map[DatabaseKey]*gorm.DB
}

// New 创建管理器实例
func New(cfg ManagerConfig) (*GormInstanceMgr, error) {
	mgr := &GormInstanceMgr{
		config: cfg,
	}

	// 预初始化所有实例的数据库连接
	for _, instance := range cfg.Instances {
		// 初始化每个数据库
		for _, db := range instance.Databases {
			key := DatabaseKey{
				Database: db.Name,
				Mode:     db.Mode,
			}

			if err := mgr.initializeDB(key, instance.Config); err != nil {
				return nil, fmt.Errorf("初始化数据库失败(%s): %w", key.String(), err)
			}
		}
	}

	return mgr, nil
}

// initializeDB 初始化单个数据库连接
func (m *GormInstanceMgr) initializeDB(key DatabaseKey, instanceCfg InstanceConfig) error {
	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		instanceCfg.Username, instanceCfg.Password, instanceCfg.Host, instanceCfg.Port, key.Database)

	// 创建GORM实例
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层连接失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(instanceCfg.Pool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(instanceCfg.Pool.MaxOpenConns)

	// 探活检查
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 使用 Ping 探活
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return fmt.Errorf("探活检查失败: %w", err)
	}

	// 存储实例
	m.instances.Store(key, db)

	return nil
}

// DB 获取数据库实例
func (m *GormInstanceMgr) DB(database string, mode AccessMode) (*gorm.DB, error) {
	if database == "" {
		return nil, fmt.Errorf("database名称不能为空")
	}

	key := DatabaseKey{
		Database: database,
		Mode:     mode,
	}

	// 获取实例
	val, ok := m.instances.Load(key)
	if !ok {
		return nil, &DatabaseNotFoundError{
			Database: database,
			Mode:     mode,
		}
	}

	return val.(*gorm.DB), nil
}

// Table 获取表级别的DB实例
func (m *GormInstanceMgr) Table(database string, mode AccessMode, table string) (*gorm.DB, error) {
	db, err := m.DB(database, mode)
	if err != nil {
		return nil, err
	}
	return db.Table(table), nil
}

// NewFromFile 从配置文件创建管理器实例
// filename 可以是相对路径，会自动搜索常见的配置目录
func NewFromFile(filename string) (*GormInstanceMgr, error) {
	loader := NewConfigLoader()
	config, err := loader.LoadConfig(filename)
	if err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %w", err)
	}
	return New(*config)
}

// MustNew 创建管理器实例，如果出错则panic
func MustNew(config ManagerConfig) *GormInstanceMgr {
	mgr, err := New(config)
	if err != nil {
		panic(err)
	}
	return mgr
}

// MustNewFromFile 从配置文件创建管理器实例，如果出错则panic
func MustNewFromFile(filename string) *GormInstanceMgr {
	mgr, err := NewFromFile(filename)
	if err != nil {
		panic(err)
	}
	return mgr
}
