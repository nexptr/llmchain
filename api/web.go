package api

import (
	"context"
	"net/http"

	"github.com/exppii/llmchain"
	"github.com/exppii/llmchain/api/conf"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap/zapcore"
)

type App struct {
	srv    *http.Server //http server
	router *gin.Engine  //http 路由表

	cf *conf.Config

	ctx context.Context
}

// NewAPPWithConfig with config
func NewAPPWithConfig(cf *conf.Config) *App {

	return &App{
		cf: cf,
	}

}

// StartContext 启动
func (m *App) StartContext(ctx context.Context) error {

	m.ctx = ctx
	//加载路由
	m.router = initGinRoute(m.cf.LogLevel)

	LogI(`llmchain version: `, llmchain.VERSION)
	LogI(`git commit hash: `, llmchain.GitHash)
	LogI(`utc build time: `, llmchain.BuildStamp)

	go func() {

		//这样设置默认可能是不安全的，因为头部字段可以伪造，需求前置的反向代理的xff 确保是对的
		m.router.SetTrustedProxies([]string{"0.0.0.0", "::"})

		m.srv = &http.Server{Addr: m.cf.APIAddr, Handler: m.router}
		LogI(`start HTTP server：`, m.cf.APIAddr)
		m.srv.ListenAndServe()

	}()

	return nil

}

// GracefulStop 退出，每个模块实现stop
func (m *App) GracefulStop() {
	if m.srv != nil {
		LogD(`quit http server...`)
		m.srv.Shutdown(m.ctx)
	}

}

func initGinRoute(level Level) *gin.Engine {

	if level == zapcore.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// log.SetFlags(log.LstdFlags) // gin will disable log flags

	router := gin.Default()

	// openAI compatible API endpoint
	router.POST("/v1/chat/completions", chatEndpointHandler)
	router.POST("/chat/completions", chatEndpointHandler)

	router.POST("/v1/edits", editEndpointHandler)
	router.POST("/edits", editEndpointHandler)

	router.POST("/v1/completions", completionEndpointHandler)
	router.POST("/completions", completionEndpointHandler)

	router.POST("/v1/embeddings", embeddingsEndpointHandler)
	router.POST("/embeddings", embeddingsEndpointHandler)

	// /v1/engines/{engine_id}/embeddings

	router.POST("/v1/engines/:model/embeddings", embeddingsEndpointHandler)

	router.GET("/v1/models", listModelsHandler)
	router.GET("/models", listModelsHandler)

	return router
}
