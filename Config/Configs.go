package Config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

var configFile []byte

type AppConfig struct {
	App APP `yaml:"app"`
}

type MySQLConfig struct {
	MySQL MySQL `yaml:"mysql"`
}

type APP struct {
	APPID string `yaml:"app_id"`
	APPSecret string `yaml:"app_secret"`
	RunPort string `yaml:"run_port"`
	JWTKey string `yaml:"jwt_key"`
}

type MySQL struct {
	Host 	 string `yaml:"mysql_host"`
	Port     string `yaml:"mysql_port"`
	UserName string `yaml:"mysql_username"`
	PassWord string `yaml:"mysql_password"`
	Database string `yaml:"mysql_database"`
}

func GetAppConfig() (Result *AppConfig, err error) {
	err = yaml.Unmarshal(configFile, &Result)
	return Result, err
}

func GetMySQLConfig() (Result *MySQLConfig, err error) {
	err = yaml.Unmarshal(configFile, &Result)
	return Result, err
}

func init() {
	var err error
	configFile, err = ioutil.ReadFile("Res/config.yaml")

	if err != nil {
		log.Println(err)
	}
}