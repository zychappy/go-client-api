package main

import (
	"github.com/astaxie/beego"
	"github.com/zychappy/go-client-api/common/global"
	_ "github.com/zychappy/go-client-api/routers"
	"github.com/zychappy/go-client-api/service"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	grpcAddress := beego.AppConfig.String("grpcaddress")

	global.TronClient = service.NewGrpcClient(grpcAddress)
	global.TronClient.Start()
	defer global.TronClient.Conn.Close()

	beego.Run()
}
