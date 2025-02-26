package mysqldb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestUser 测试用的用户模型
type TestUser struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"size:255"`
	Email     string    `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// 测试配置
var testConfig = ManagerConfig{
	Instances: []struct {
		Config    InstanceConfig   `yaml:"config" toml:"config"`
		Databases []DatabaseConfig `yaml:"databases" toml:"databases"`
	}{
		{
			Config: InstanceConfig{
				Host:     "localhost",
				Port:     3306,
				Username: "root",
				Password: "",
				Pool: PoolConfig{
					MaxIdleConns: 10,
					MaxOpenConns: 100,
				},
			},
			Databases: []DatabaseConfig{
				{
					Name: "test",
					Mode: ReadOnly,
				},
				{
					Name: "test",
					Mode: ReadWrite,
				},
			},
		},
	},
}

// TestNew 测试创建管理器实例
func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  ManagerConfig
		wantErr bool
	}{
		{
			name:    "Valid configuration",
			config:  testConfig,
			wantErr: false,
		},
		{
			name: "Invalid host configuration",
			config: ManagerConfig{
				Instances: []struct {
					Config    InstanceConfig   `yaml:"config" toml:"config"`
					Databases []DatabaseConfig `yaml:"databases" toml:"databases"`
				}{
					{
						Config: InstanceConfig{
							Host:     "invalid-host",
							Port:     3306,
							Username: "test_user",
							Password: "test_password",
						},
						Databases: []DatabaseConfig{
							{
								Name: "test",
								Mode: ReadOnly,
							},
						},
					},
				},
			},
			wantErr: true, // 无效的host应该返回错误
		},
		{
			name: "Invalid credentials",
			config: ManagerConfig{
				Instances: []struct {
					Config    InstanceConfig   `yaml:"config" toml:"config"`
					Databases []DatabaseConfig `yaml:"databases" toml:"databases"`
				}{
					{
						Config: InstanceConfig{
							Host:     testConfig.Instances[0].Config.Host,
							Port:     testConfig.Instances[0].Config.Port,
							Username: "invalid_user",
							Password: "invalid_password",
						},
						Databases: []DatabaseConfig{
							{
								Name: "test",
								Mode: ReadOnly,
							},
						},
					},
				},
			},
			wantErr: true, // 无效的凭证应该返回错误
		},
		{
			name: "Empty configuration",
			config: ManagerConfig{
				Instances: []struct {
					Config    InstanceConfig   `yaml:"config" toml:"config"`
					Databases []DatabaseConfig `yaml:"databases" toml:"databases"`
				}{},
			},
			wantErr: false, // 空配置是合法的
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr, err := New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "初始化数据库失败") // 错误应该包含初始化失败信息
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, mgr)

				// 如果是有效配置，测试实例是否初始化成功
				if len(tt.config.Instances) > 0 {
					for _, inst := range tt.config.Instances {
						for _, db := range inst.Databases {
							key := DatabaseKey{
								Database: db.Name,
								Mode:     db.Mode,
							}
							// 检查实例是否存在
							val, ok := mgr.instances.Load(key)
							assert.True(t, ok, "数据库实例应该存在: %s", key.String())
							assert.NotNil(t, val, "数据库实例不应为空: %s", key.String())
						}
					}
				}
			}
		})
	}
}

// TestGormInstanceMgr_DB 测试获取数据库实例
func TestGormInstanceMgr_DB(t *testing.T) {
	// 跳过真实数据库连接的测试
	t.Skip("Skipping test that requires database connection")

	mgr, err := New(testConfig)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		database string
		mode     AccessMode
		wantErr  bool
	}{
		{
			name:     "Valid readonly database",
			database: "test",
			mode:     ReadOnly,
			wantErr:  false,
		},
		{
			name:     "Valid readwrite database",
			database: "test",
			mode:     ReadWrite,
			wantErr:  false,
		},
		{
			name:     "Non-existent database",
			database: "non_existent_db",
			mode:     ReadOnly,
			wantErr:  true,
		},
		{
			name:     "Invalid mode",
			database: "test_db",
			mode:     "invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := mgr.DB(tt.database, tt.mode)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				assert.IsType(t, &gorm.DB{}, db)
			}
		})
	}
}

// TestGormInstanceMgr_Table 测试获取表级别的DB实例
func TestGormInstanceMgr_Table(t *testing.T) {
	// 跳过真实数据库连接的测试
	t.Skip("Skipping test that requires database connection")

	mgr, err := New(testConfig)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		database string
		mode     AccessMode
		table    string
		wantErr  bool
	}{
		{
			name:     "Valid table in readonly database",
			database: "test",
			mode:     ReadOnly,
			table:    "users",
			wantErr:  false,
		},
		{
			name:     "Valid table in readwrite database",
			database: "test",
			mode:     ReadWrite,
			table:    "test_users",
			wantErr:  false,
		},
		{
			name:     "Non-existent database",
			database: "non_existent_db",
			mode:     ReadOnly,
			table:    "test_users",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := mgr.Table(tt.database, tt.mode, tt.table)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				assert.IsType(t, &gorm.DB{}, db)
			}
		})
	}
}

// TestGormInstanceMgr_ConcurrentAccess 测试并发访问
func TestGormInstanceMgr_ConcurrentAccess(t *testing.T) {
	// 跳过真实数据库连接的测试
	t.Skip("Skipping test that requires database connection")

	mgr, err := New(testConfig)
	assert.NoError(t, err)

	// 并发测试
	concurrency := 10
	done := make(chan bool)

	for i := 0; i < concurrency; i++ {
		go func() {
			// 测试只读实例
			db1, err := mgr.DB("test", ReadOnly)
			assert.NoError(t, err)
			assert.NotNil(t, db1)

			// 测试读写实例
			db2, err := mgr.DB("test", ReadWrite)
			assert.NoError(t, err)
			assert.NotNil(t, db2)

			// 测试表级别访问
			db3, err := mgr.Table("test", ReadOnly, "test_users")
			assert.NoError(t, err)
			assert.NotNil(t, db3)

			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < concurrency; i++ {
		<-done
	}
}

// TestGormInstanceMgr_ConnectionFailure 测试连接失败的情况
func TestGormInstanceMgr_ConnectionFailure(t *testing.T) {
	// 创建一个无效的配置
	invalidConfig := ManagerConfig{
		Instances: []struct {
			Config    InstanceConfig   `yaml:"config" toml:"config"`
			Databases []DatabaseConfig `yaml:"databases" toml:"databases"`
		}{
			{
				Config: InstanceConfig{
					Host:     "invalid-host",
					Port:     3306,
					Username: "test_user",
					Password: "test_password",
				},
				Databases: []DatabaseConfig{
					{
						Name: "test",
						Mode: ReadOnly,
					},
				},
			},
		},
	}

	// 创建管理器实例
	_, err := New(invalidConfig)
	assert.Error(t, err) // 创建管理器应该成功
}

type TestUserTable struct {
	ID       int64 `json:"id" gorm:"column:id;type:int(11);notnull;autoincrement;primary_key"`
	MemberID int64 `json:"member_id" gorm:"column:member_id;index;type:int(11);notnull;default:0;comment:用户 ID"`
}

func (t TestUserTable) TableName() string {
	return "test_users"
}

// TestCreateTable 测试创建表
func TestCreateTable(t *testing.T) {
	// 创建管理器实例
	mgr, err := New(testConfig)
	assert.NoError(t, err)

	// 获取读写实例
	db, err := mgr.DB("test", ReadWrite)
	assert.NoError(t, err)

	// 删除表（如果存在）
	//err = db.Migrator().DropTable(&TestUserTable{})
	//assert.NoError(t, err)

	// 创建表
	err = db.AutoMigrate(&TestUserTable{})
	assert.NoError(t, err)

	// 插入测试数据
	testUsers := []TestUserTable{
		{MemberID: 1001},
		{MemberID: 1002},
		{MemberID: 1003},
	}

	// 批量插入
	result := db.Create(&testUsers)
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(3), result.RowsAffected)

	// 查询数据
	var users []TestUserTable
	result = db.Find(&users)
	assert.NoError(t, result.Error)
	assert.Equal(t, 3, len(users))

	// 验证数据
	for i, user := range users {
		assert.Equal(t, testUsers[i].MemberID, user.MemberID)
		assert.NotZero(t, user.ID) // ID应该自增
	}

	// 测试只读实例
	readOnlyDB, err := mgr.DB("test", ReadOnly)
	assert.NoError(t, err)

	// 尝试从只读实例查询
	var readOnlyUsers []TestUserTable
	result = readOnlyDB.Find(&readOnlyUsers)
	assert.NoError(t, result.Error)
	assert.Equal(t, len(testUsers), len(readOnlyUsers))
}
