package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func basicConfig() {
	viper.SetConfigName("config") // 配置文件名（不带扩展名）
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 配置文件路径
	viper.AddConfigPath("$HOME/.app")
	if err := viper.ReadInConfig(); err != nil {
		panic("读取配置文件失败: " + err.Error())
	}
}
func loadConfigToEnv() {
	viper.AutomaticEnv()      // 从环境变量加载配置
	viper.SetEnvPrefix("APP") // 环境变量前缀，例如 APP_SERVER_PORT,设置前缀后，Viper 会自动将配置键转换为环境变量名，例如 server.port 会转换为 APP_SERVER_PORT
}
func main() {
	loadConfigToEnv()
	port := viper.GetString("server.port")
	dbhost := viper.GetString("database.host")

	fmt.Printf("Server will run on port: %s\n", port)
	fmt.Printf("Database host: %s\n", dbhost)

}
