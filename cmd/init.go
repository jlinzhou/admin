package main

import (
	"admin/confs"
	"admin/routers"
	"github.com/gin-gonic/gin"
	"net/http"
)


func initWeb() {
	gin.SetMode(gin.DebugMode) //调试模式
	app := gin.New()
	srv := &http.Server{
		Addr:         ":" + confs.AdminConf.ServerConf.HTTPPort,
		Handler:      app,
		ReadTimeout:  5 * confs.AdminConf.ServerConf.ReadTimeout,
		WriteTimeout: 10 * confs.AdminConf.ServerConf.WriteTimeout,
		IdleTimeout:  15 * confs.AdminConf.ServerConf.IdleTimeout,
	}
	routers.RegisterRouter(app)
	srv.ListenAndServe()
}