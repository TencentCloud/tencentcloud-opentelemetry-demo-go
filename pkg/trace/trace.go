package trace

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
)

// InitTraceProvider configures an OpenTelemetry exporter and gindemo provider
func InitTraceProvider(ctx context.Context) *sdktrace.TracerProvider {

	//New exporter
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint("ap-guangzhou.apm.tencentcs.com:4317"), // 替换成apm上报地址
		otlptracegrpc.WithInsecure(),
	}
	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		log.Fatal(err)
	}

	//设置Token，也可以设置环境变量：OTEL_RESOURCE_ATTRIBUTES=token=xxxxxxxxx
	r, err := resource.New(ctx, []resource.Option{
		resource.WithAttributes(
			attribute.KeyValue{Key: "token", Value: attribute.StringValue("xxxxxx")}, // 替换成上报实例的Token
			attribute.KeyValue{Key: "service.name", Value: attribute.StringValue("Test-service")},
			attribute.KeyValue{Key: "service.instance", Value: attribute.StringValue("127.0.0.1")}, // 替换成本机的IP地址
		),
	}...)
	if err != nil {
		log.Fatal(err)
	}

	//New TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
