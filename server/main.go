package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/tongruirenye/OrgICSX5/server/config"
	"github.com/tongruirenye/OrgICSX5/server/ics"
	"github.com/tongruirenye/OrgICSX5/server/storage"
)

func main() {
	configFlag := pflag.StringP("config", "c", "", "web proxy config path.")
	parseFlag := pflag.BoolP("parse", "p", false, "parse")
	localFlag := pflag.StringP("local", "l", "", "local file path")
	pflag.Parse()

	if err := config.InitConfig(*configFlag); err != nil {
		fmt.Println(err)
		return
	}

	logfile, err := os.OpenFile(config.AppPath+"app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	storage.InitStorage()

	icsParser := ics.NewIcs(log.Default())
	if *parseFlag {
		if *localFlag != "" {
			if e := icsParser.DoLocal(config.AppPath + *localFlag); e != nil {
				fmt.Println(e)
			}
		} else {
			icsParser.Do()
		}
		return
	}

	if config.AppConfig.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(ics.UseIcs(icsParser))
	r.Use(func(c *gin.Context) {
		c.Next()
		log.Printf("path=%s ip=%s method=%s", c.Request.URL.Path, c.ClientIP(), c.Request.Method)
	})
	r.POST("/ics/gen", ics.GenIcs)

	r.Static("/public", "./public")

	http_server := &http.Server{
		Addr:    config.AppConfig.Port,
		Handler: r,
	}

	go func() {
		if err := http_server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("http server error:%v\n", err)
		}
	}()

	go func() {
		icsParser.Run()
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sc
	log.Println("shutdow server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := http_server.Shutdown(ctx); err != nil {
		log.Printf("http server shutdown error:%v\n", err)
	}
	log.Println("shutdown http")

	icsParser.Close()
	log.Println("shutdow ics")
}
