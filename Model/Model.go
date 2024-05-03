package Model

import (
	"github.com/sunbelife/Prelook_Strobe_Backend/Config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"fmt"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"
const FALSE = "false"
const TRUE = "true"

var db *gorm.DB

type Data struct {
	IDs []int
	Mode string
}

func Init() {
	MySQLConfig, err := Config.GetMySQLConfig()

	if err != nil {
		log.Println("mysql config err")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", MySQLConfig.MySQL.UserName, MySQLConfig.MySQL.PassWord, MySQLConfig.MySQL.Host, MySQLConfig.MySQL.Port, MySQLConfig.MySQL.Database)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 如果数据库链接异常则抛出
	if err != nil {
		log.Fatal(err)
	}
}
