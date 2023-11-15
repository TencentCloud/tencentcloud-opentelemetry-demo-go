package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-opentelemetry-demo/internal/http"
	"go-opentelemetry-demo/pkg/gormclient"
	"go-opentelemetry-demo/pkg/redisclient"
	"go-opentelemetry-demo/pkg/trace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
)

func main() {
	ctx := context.Background()
	//初始化redis
	redisclient.InitRedis()
	//初始化gorm
	gormclient.InitGorm()
	//初始化ot的traceProvider
	tp := trace.InitTraceProvider(ctx)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	r := gin.New()
	// 添加trace中间件
	r.Use(otelgin.Middleware("gin-demo2"))

	InitRouter(r)

	_ = r.Run(":8089")
}

// InitRouter 初始化路由
func InitRouter(router *gin.Engine) {
	router.GET("/postMessage", http.PostMessage)
}
