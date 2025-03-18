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
}

// NewConfigLoader 创建配置加载器
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

// LoadConfig 加载配置文件
func (l *ConfigLoader) LoadConfig(filePath string) (*ManagerConfig, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("配置文件未找到: %s", filePath)
	}

	// 读取配置文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config ManagerConfig

	// 根据文件扩展名选择解析方式
	ext := strings.ToLower(filepath.Ext(filePath))
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
