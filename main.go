package main

import (
	"github.com/sunbelife/Prelook_Strobe_Backend/Config"
	"github.com/sunbelife/Prelook_Strobe_Backend/Model"
	"github.com/sunbelife/Prelook_Strobe_Backend/Router"
	"log"
)

func main() {
	// 初始化路由器
	r := Router.InitRouter()
	// 初始化数据库
	Model.Init()

	AppConfig, err := Config.GetAppConfig()

	if err != nil {
		log.Println(err)
	}

	// 启动端口
	r.Run("0.0.0.0:" + AppConfig.App.RunPort)
}
