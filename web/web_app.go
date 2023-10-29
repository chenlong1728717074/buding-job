package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type WebApp struct {
	engine *gin.Engine
	server *http.Server
}

func NewWebApp() *WebApp {
	return &WebApp{
		engine: gin.Default(),
	}
}
func (app *WebApp) Start() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8081),
		Handler: app.engine, // 使用 Gin 引擎作为处理程序
	}
	//路由
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	app.server = server
}
func (app *WebApp) Stop(ctx context.Context) {
	if err := app.server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown err:", err)
	}
}
