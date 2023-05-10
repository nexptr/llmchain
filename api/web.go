package api

import (
	"context"
	"net/http"
	"os"

	"github.com/exppii/llmchain"
	"github.com/exppii/llmchain/api/conf"
	"github.com/exppii/llmchain/api/model"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap/zapcore"
)

type App struct {
	mng *model.Manager

	srv *http.Server //http server

	ctx context.Context
}

// NewAPPWithConfig with config
func NewAPPWithConfig(cf *conf.Config) *App {

	LogI(`llmchain version: `, llmchain.VERSION)
	LogI(`git commit hash: `, llmchain.GitHash)
	LogI(`UTC build time: `, llmchain.BuildStamp)
	LogI(`HTTP server address: `, cf.APIAddr)

	manager := model.NewModelManager(cf)

	router := initGinRoute(manager, cf.LogLevel)

	srv := &http.Server{Addr: cf.APIAddr, Handler: router}

	return &App{
		mng: manager,
		srv: srv,
	}

}

// StartContext 启动
func (m *App) StartContext(ctx context.Context) error {

	m.ctx = ctx

	// m.mng.Load() may be slow，in order not to block the main process，
	// goroutine is used here, so we can use ctrl+c to terminate it
	go func() {
		if err := m.mng.Load(); err != nil {
			LogE(`load model failed: `, err.Error())
			os.Exit(1)
		}
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

func initGinRoute(manager *model.Manager, level Level) *gin.Engine {

	if level == zapcore.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// log.SetFlags(log.LstdFlags) // gin will disable log flags

	router := gin.Default()

	// openAI compatible API endpoint
	router.POST("/v1/chat/completions", chatEndpointHandler(manager))
	router.POST("/chat/completions", chatEndpointHandler(manager))

	router.POST("/v1/edits", editEndpointHandler(manager))
	router.POST("/edits", editEndpointHandler(manager))

	router.POST("/v1/completions", completionEndpointHandler(manager))
	router.POST("/completions", completionEndpointHandler(manager))

	router.POST("/v1/embeddings", embeddingsEndpointHandler(manager))
	router.POST("/embeddings", embeddingsEndpointHandler(manager))

	// /v1/engines/{engine_id}/embeddings

	router.POST("/v1/engines/:model/embeddings", embeddingsEndpointHandler(manager))

	router.GET("/v1/models", listModelsHandler(manager))
	router.GET("/models", listModelsHandler(manager))

	//这样设置默认可能是不安全的，因为头部字段可以伪造，需求前置的反向代理的xff 确保是对的
	router.SetTrustedProxies([]string{"0.0.0.0", "::"})

	return router
}
