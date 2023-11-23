package web

import (
	"buding-job/common/log"
	"buding-job/web/api"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebApp struct {
	engine *gin.Engine
	server *http.Server
}

func NewWebApp() *WebApp {
	// 设置 Gin 使用 logrus 的日志库
	gin.DefaultWriter = log.GetLog().Writer()
	gin.DefaultErrorWriter = log.GetLog().Writer()
	engine := gin.Default()
	return &WebApp{
		engine: engine,
	}
}
func (app *WebApp) Start() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8081),
		Handler: app.engine, // 使用 Gin 引擎作为处理程序
	}
	//路由
	app.router()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.GetLog().Fatalf("listen: %s\n", err)
		}
	}()
	app.server = server
}
func (app *WebApp) Stop(ctx context.Context) {
	if err := app.server.Shutdown(ctx); err != nil {
		log.GetLog().Println("Server Shutdown err:", err)
	}
}

func (app *WebApp) router() {
	api.NewJobInfoApi(app.engine).Router()
	api.NewJobManagementApi(app.engine).Router()
}
