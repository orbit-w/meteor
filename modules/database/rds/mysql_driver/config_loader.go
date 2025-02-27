package mysqldb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// ConfigLoader 配置加载器
type ConfigLoader struct {
	searchPaths []string // 配置文件搜索路径
}

// NewConfigLoader 创建配置加载器
func NewConfigLoader(searchPaths ...string) *ConfigLoader {
	if len(searchPaths) == 0 {
		// 默认搜索路径
		searchPaths = []string{
			".",
			"./config",
			"./configs",
			"../config",
			"../configs",
		}
	}
	return &ConfigLoader{searchPaths: searchPaths}
}

// LoadConfigFromPath 从指定路径加载配置文件
func (l *ConfigLoader) LoadConfigFromPath(path string) (*ManagerConfig, error) {
	// 检查文件是否存在
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("配置文件未找到: %s", path)
	}

	// 读取配置文件
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config ManagerConfig

	// 根据文件扩展名选择解析方式
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(content, &config)
	case ".toml":
		err = toml.Unmarshal(content, &config)
	default:
		return nil, fmt.Errorf("不支持的配置文件格式: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}

// LoadConfig 加载配置文件
func (l *ConfigLoader) LoadConfig(filename string) (*ManagerConfig, error) {
	// 首先尝试直接使用提供的路径
	if _, err := os.Stat(filename); err == nil {
		return l.LoadConfigFromPath(filename)
	}

	// 在搜索路径中查找配置文件
	for _, path := range l.searchPaths {
		file := filepath.Join(path, filename)
		if _, err := os.Stat(file); err == nil {
			return l.LoadConfigFromPath(file)
		}
	}

	return nil, fmt.Errorf("配置文件未找到: %s (搜索路径: %v)", filename, l.searchPaths)
}
