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
      mode: "readonly"
      pool:
        max_idle_conns: 10
        max_open_conns: 100
      log:
        level: "info"
    databases:
      - name: "test"`

	tomlConfig := `[[instances]]
[instances.config]
host = "localhost"
port = 3306
username = "root"
password = ""
mode = "readonly"
[instances.config.pool]
max_idle_conns = 10
max_open_conns = 100
[instances.config.log]
level = "info"

[[instances.databases]]
name = "test"`

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
			assert.Equal(t, ReadOnly, instance.Config.Mode)

			// 验证连接池配置
			assert.Equal(t, 10, instance.Config.Pool.MaxIdleConns)
			assert.Equal(t, 100, instance.Config.Pool.MaxOpenConns)

			// 验证数据库配置
			require.NotEmpty(t, instance.Databases)
			assert.Equal(t, "test", instance.Databases[0].Name)

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

func TestConfigLoader_LoadConfigFromPath(t *testing.T) {
	// 创建临时目录和配置文件
	tempDir := t.TempDir()

	// 创建测试配置文件
	yamlConfig := filepath.Join(tempDir, "config.yaml")
	yamlContent := `instances:
  - config:
      host: "localhost"
      port: 3306
      username: "root"
      password: ""
      mode: "readonly"
      pool:
        max_idle_conns: 10
        max_open_conns: 100
      log:
        level: "info"
    databases:
      - name: "test"`
	require.NoError(t, os.WriteFile(yamlConfig, []byte(yamlContent), 0644))

	tomlConfig := filepath.Join(tempDir, "config.toml")
	tomlContent := `[[instances]]
[instances.config]
host = "localhost"
port = 3306
username = "root"
password = ""
mode = "readonly"
[instances.config.pool]
max_idle_conns = 10  # 最大空闲连接数
max_open_conns = 100 # 最大打开连接数
[instances.config.log]
level = "info"
[[instances.databases]]
name = "test"`
	require.NoError(t, os.WriteFile(tomlConfig, []byte(tomlContent), 0644))

	// 切换到临时目录以测试相对路径
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldWd)
	require.NoError(t, os.Chdir(tempDir))

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "Load YAML config with absolute path",
			filePath: yamlConfig,
			wantErr:  false,
		},
		{
			name:     "Load TOML config with absolute path",
			filePath: tomlConfig,
			wantErr:  false,
		},
		{
			name:     "Load YAML config with relative path",
			filePath: "./config.yaml",
			wantErr:  false,
		},
		{
			name:     "Load TOML config with relative path",
			filePath: "./config.toml",
			wantErr:  false,
		},
		{
			name:     "Load YAML config with parent directory relative path",
			filePath: filepath.Join("..", filepath.Base(tempDir), "config.yaml"),
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
			// 创建配置加载器
			loader := NewConfigLoader()

			// 从指定路径加载配置
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
			assert.Equal(t, ReadOnly, instance.Config.Mode)

			// 验证连接池配置
			assert.Equal(t, 10, instance.Config.Pool.MaxIdleConns)
			assert.Equal(t, 100, instance.Config.Pool.MaxOpenConns)

			// 验证数据库配置
			require.NotEmpty(t, instance.Databases)
			assert.Equal(t, "test", instance.Databases[0].Name)

			// 验证日志配置
			assert.Equal(t, "info", instance.Config.Log.Level)
		})
	}
}
