package main

import (
	"com.lu/OnlineTools/controller"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"os"
)

func main() {
	e := echo.New()
	RegisterRoute(e)
	f, err := os.OpenFile("./logs/log.txt", os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("OpenFile err:%v", err)
	}
	e.Logger.SetOutput(f)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}",` +
			`"method":"${method}","uri":"${uri}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output: f,
	}))
	e.Logger.SetLevel(log.INFO)
	e.Logger.Info("server start")
	err = e.Start(":8099")
	if err != nil{
		log.Fatalf("server start err:%v", err)
	}
}

func RegisterRoute(e *echo.Echo) {
	root := e.Group("OnlineTools")
	root.POST("/diff", controller.Diff)
}
