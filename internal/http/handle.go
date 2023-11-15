package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-opentelemetry-demo/pkg/gormclient"
	"go-opentelemetry-demo/pkg/redisclient"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// PostMessage 发送消息
func PostMessage(ginCtx *gin.Context) {
	span := trace.SpanFromContext(ConvertToGinContext(ginCtx))
	span.SetAttributes(attribute.KeyValue{
		Key:   "component",
		Value: attribute.StringValue("http")})
	//Redis调用，需要对gin context进行转换
	redisRequest(ConvertToGinContext(ginCtx))
	// gorm调用，需要对gin context进行转换
	gormRequest(ConvertToGinContext(ginCtx))
	// 手动埋点示例
	mockTrace()
	var context *gin.Context
	_ = context.Err()
	ginCtx.JSON(200, gin.H{
		"Hello": "otel-gin-demo",
	})
}

// ConvertToGinContext convert to ginContext 由于gincontext对context进行了封装，需要拆出来trace需要的context
func ConvertToGinContext(c context.Context) context.Context {
	var ct context.Context
	ctx, ok := c.(*gin.Context)
	if ok {
		ct = ctx.Request.Context()
	} else {
		ct = c
	}
	return ct
}

func gormRequest(ctx context.Context) {
	var num int
	if err := gormclient.GormDB.WithContext(ctx).Raw("SELECT 42").Scan(&num).Error; err != nil {
		panic(err)
	}
}

func redisRequest(ctx context.Context) {
	err := redisclient.Rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := redisclient.Rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := redisclient.Rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

}

// 手动埋点，模拟trace
func mockTrace() {

	tracer := otel.Tracer("example.com/basic")

	ctx0 := context.Background()

	ctx1, finish1 := tracer.Start(ctx0, "foo")
	defer finish1.End()

	ctx2, finish2 := tracer.Start(ctx1, "bar")
	defer finish2.End()

	ctx3, finish3 := tracer.Start(ctx2, "baz")
	defer finish3.End()

	ctx := ctx3
	getSpan(ctx)
	addAttribute(ctx)
	addEvent(ctx)
	recordException(ctx)
	createChild(ctx, tracer)
}

// example of getting the current span
// 获取当前的Span。
func getSpan(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	fmt.Printf("current span: %v\n", span)
}

// example of adding an attribute to a span
// 向Span中添加属性值。
func addAttribute(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.KeyValue{
		Key:   "label-key-1",
		Value: attribute.StringValue("label-value-1")})
}

// example of adding an event to a span
// 向Span中添加事件。
func addEvent(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("event1", trace.WithAttributes(
		attribute.String("event-attr1", "event-string1"),
		attribute.Int64("event-attr2", 10)))
}

// example of recording an exception
// 记录Span结果以及错误信息。
func recordException(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(errors.New("exception has occurred"))
	span.SetStatus(codes.Error, "pkg error")
}

// example of creating a child span
// 创建子Span。
func createChild(ctx context.Context, tracer trace.Tracer) {
	// span := gindemo.SpanFromContext(ctx)
	_, childSpan := tracer.Start(ctx, "child")
	childSpan.SetStatus(codes.Error, "1")
	defer childSpan.End()
}
