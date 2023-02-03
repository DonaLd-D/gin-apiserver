package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/junhui/gin_demo/routers"
	"github.com/junhui/gin_demo/config"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
)

func main() {
	pflag.Parse()
	// 配置初始化
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}
	gin.SetMode(viper.GetString("runmode"))

	//gin实例
	g := gin.New()
	//加载中间件
	middlewares := []gin.HandlerFunc{}
	//加载路由
	routers.Load(
		g,
		middlewares...,
	)

	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Print("The router has been deployed successfully.")
	}()

	log.Printf("Start to listening the incoming requests on http address: %s", ":8080")
	log.Printf(http.ListenAndServe(":8080", g).Error())
}

func pingServer() error {
	for i := 0; i < 2; i++ {
		resp, err := http.Get("http://127.0.0.1:8080" + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
