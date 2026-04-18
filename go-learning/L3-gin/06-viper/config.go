package main

import "github.com/spf13/viper"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

func NewConfig() *Config {
	return &Config{}
}
func (c *Config) LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}
