package mysqldb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigLoader_LoadConfig(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()

	// 准备测试配置文件
	yamlConfig := `instances:
  - config:
      host: "localhost"
      port: 3306
      username: "root"
      password: ""
      pool:
        maxIdleConns: 10
        maxOpenConns: 100
      log:
        level: "info"
    databases:
      - name: "test"
        mode: "readonly"`

	tomlConfig := `[[instances]]
config.host = "localhost"
config.port = 3306
config.username = "root"
config.password = ""
config.pool.maxIdleConns = 10
config.pool.maxOpenConns = 100
config.log.level = "info"

[[instances.databases]]
name = "test"
mode = "readonly"`

	// 创建测试文件
	yamlPath := filepath.Join(tempDir, "config.yaml")
	tomlPath := filepath.Join(tempDir, "config.toml")
	require.NoError(t, os.WriteFile(yamlPath, []byte(yamlConfig), 0644))
	require.NoError(t, os.WriteFile(tomlPath, []byte(tomlConfig), 0644))

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "Load config.yaml",
			filePath: yamlPath,
			wantErr:  false,
		},
		{
			name:     "Load config.toml",
			filePath: tomlPath,
			wantErr:  false,
		},
		{
			name:     "Non-existent file",
			filePath: filepath.Join(tempDir, "non_existent.yaml"),
			wantErr:  true,
		},
		{
			name:     "Invalid file extension",
			filePath: filepath.Join(tempDir, "config.txt"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewConfigLoader()
			config, err := loader.LoadConfig(tt.filePath)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, config)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, config)

			// 验证配置内容
			require.NotEmpty(t, config.Instances)
			instance := config.Instances[0]

			// 验证基本配置
			assert.Equal(t, "localhost", instance.Config.Host)
			assert.Equal(t, 3306, instance.Config.Port)
			assert.Equal(t, "root", instance.Config.Username)
			assert.Equal(t, "", instance.Config.Password)

			// 验证连接池配置
			assert.Equal(t, 10, instance.Config.Pool.MaxIdleConns)
			assert.Equal(t, 100, instance.Config.Pool.MaxOpenConns)

			// 验证数据库配置
			require.NotEmpty(t, instance.Databases)
			assert.Equal(t, "test", instance.Databases[0].Name)
			assert.Equal(t, ReadOnly, instance.Databases[0].Mode)

			// 验证日志配置
			assert.Equal(t, "info", instance.Config.Log.Level)
		})
	}
}

func TestConfigLoader_DirectPath(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()

	// 创建配置文件
	configPath := filepath.Join(tempDir, "mysql.yaml")
	yamlConfig := `instances:
  - config:
      host: "localhost"
      port: 3306`
	require.NoError(t, os.WriteFile(configPath, []byte(yamlConfig), 0644))

	// 测试直接使用文件路径
	loader := NewConfigLoader()
	config, err := loader.LoadConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "localhost", config.Instances[0].Config.Host)
	assert.Equal(t, 3306, config.Instances[0].Config.Port)
}
