package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	DBConfig `mapstructure:"db"`
}
type DBConfig struct {
	DB_USERNAME string `mapstructure:"username"`
	DB_PASSWORD string `mapstructure:"password"`
	DB_NAME     string `mapstructure:"dbname"`
	DB_HOST     string `mapstructure:"host"`
	DB_PORT     string `mapstructure:"port"`
	CHAR_SET    string `mapstructure:"charset"`
}

var Config AppConfig

func InitConfig() {
	viper.SetConfigName("config") // 配置文件名（不带扩展名）
	viper.AddConfigPath(".")      // 读取配置文件当前目录
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read config file: " + err.Error())
	}
	if err := viper.Unmarshal(&Config); err != nil {
		panic("Failed to unmarshal config: " + err.Error())
	}
	fmt.Println(Config)
}

func GetDBDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		Config.DB_USERNAME,
		Config.DB_PASSWORD,
		Config.DB_HOST,
		Config.DB_PORT,
		Config.DB_NAME,
		Config.CHAR_SET)
}
