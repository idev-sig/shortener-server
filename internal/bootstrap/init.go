package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "toml"
)

func init() {
	// 检查是否通过环境变量指定了配置文件
	customConfigFile := os.Getenv("SHORTENER_CONFIG_FILE")

	if customConfigFile != "" {
		// 使用指定的配置文件
		viper.SetConfigFile(customConfigFile)
	} else {
		// 使用默认配置文件搜索路径
		viper.SetConfigName(configName)
		viper.AddConfigPath("./config/dev")
		viper.AddConfigPath("./config/prod")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
	}

	initDefaultConfig()

	if err := viper.ReadInConfig(); err != nil {
		// 处理配置文件不存在的情况
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 如果指定了自定义配置文件但不存在，则报错
			if customConfigFile != "" {
				panic(fmt.Errorf("指定的配置文件不存在: %s", customConfigFile))
			}
			// 创建并保存默认配置
			configFile := filepath.Join(configName + "." + configType)
			if err := viper.SafeWriteConfigAs(configFile); err != nil {
				panic(fmt.Errorf("write config failed: %s\n%v", configFile, err))
			}
		} else {
			// 其他类型的配置错误
			panic(
				fmt.Errorf("fatal error config file: %w", err),
			)
		}
	}

	// log.Printf("config: %+v\n", viper.AllSettings())
	bootstrap()
}

// 初始化
func bootstrap() {
	// init shared config
	initSharedConfig()

	// init db
	initDB()

	// init cache
	initCache()

	// init geoip
	initGeoIP()
}
