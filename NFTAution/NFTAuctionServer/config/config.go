package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	AppConfig        AppConfig        `mapstructure:"app"`
	DBConfig         DBConfig         `mapstructure:"db"`
	BlockchainConfig BlockchainConfig `mapstructure:"blockchain"`
}
type AppConfig struct {
	Port string `mapstructure:"port"`
}
type DBConfig struct {
	DB_USERNAME string `mapstructure:"username"`
	DB_PASSWORD string `mapstructure:"password"`
	DB_NAME     string `mapstructure:"dbname"`
	DB_HOST     string `mapstructure:"host"`
	DB_PORT     string `mapstructure:"port"`
	CHAR_SET    string `mapstructure:"charset"`
}
type BlockchainConfig struct {
	RPC_URL                 string `mapstructure:"rpc_url"`
	Private_KEY             string `mapstructure:"private_key"`
	Aution_Contract_Address string `mapstructure:"aution_contract_address"`
	NFT_Contract_Address    string `mapstructure:"nft_contract_address"`
	Buyer                   string `mapstructure:"buyer"`
}

var Config ServerConfig

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}

func GetDBDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		Config.DBConfig.DB_USERNAME,
		Config.DBConfig.DB_PASSWORD,
		Config.DBConfig.DB_HOST,
		Config.DBConfig.DB_PORT,
		Config.DBConfig.DB_NAME,
		Config.DBConfig.CHAR_SET)
}
