package mysqldb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigLoader_LoadConfig(t *testing.T) {
	// 获取当前目录
	currentDir, err := os.Getwd()
	require.NoError(t, err)

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "Load config.yaml",
			filename: "config.yaml",
			wantErr:  false,
		},
		{
			name:     "Load config.toml",
			filename: "config.toml",
			wantErr:  false,
		},
		{
			name:     "Non-existent file",
			filename: "non_existent.yaml",
			wantErr:  true,
		},
		{
			name:     "Invalid file extension",
			filename: "config.txt",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建配置加载器，使用当前目录
			loader := NewConfigLoader(currentDir)

			// 加载配置
			config, err := loader.LoadConfig(tt.filename)

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

func TestConfigLoader_SearchPaths(t *testing.T) {
	// 创建嵌套的测试目录结构
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	// 在 config 目录中创建配置文件
	configPath := filepath.Join(configDir, "mysql.yaml")
	yamlConfig := `instances:
  - config:
      host: "localhost"
      port: 3306`
	require.NoError(t, os.WriteFile(configPath, []byte(yamlConfig), 0644))

	// 切换到临时目录
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldWd)
	require.NoError(t, os.Chdir(tempDir))

	// 测试搜索路径
	loader := NewConfigLoader()
	config, err := loader.LoadConfig("mysql.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, config)
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
      pool:
        max_idle_conns: 10
        max_open_conns: 100
      log:
        level: "info"
    databases:
      - name: "test"
        mode: "readonly"`
	require.NoError(t, os.WriteFile(yamlConfig, []byte(yamlContent), 0644))

	tomlConfig := filepath.Join(tempDir, "config.toml")
	tomlContent := `[[instances]]
[instances.config]
host = "localhost"
port = 3306
username = "root"
password = ""
[instances.config.pool]
max_idle_conns = 10  # 最大空闲连接数
max_open_conns = 100 # 最大打开连接数
[instances.config.log]
level = "info"
[[instances.databases]]
name = "test"
mode = "readonly"`
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
			config, err := loader.LoadConfigFromPath(tt.filePath)

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
