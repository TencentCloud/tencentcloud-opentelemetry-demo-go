
# 0. 背景
该demo模拟了一个gin框架搭建的http服务，其中有使用gorm以及go-redis框架 

还模拟了一个grpc框架搭建的c/s服务模式

# 1. 目录结构

```
.
├── README.md
├── cmd
│   ├── grpc
│   │   └── main.go
│   └── http
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   └── http
│       └── handle.go
└── pkg
    ├── gormclient
    │   └── gorm.go
    ├── redisclient
    │   └── redis.go
    ├── trace
    │   └── trace.go
    └── grpcdemo
        ├── cilent.go
        ├── server.go
        ├── hello-service.proto 
        ├── hello-service.pb.go
        └── hello-service_grpc.pb.go
 
```

1. demo入口见cmd下的两个main.go

2. trace相关配置在trace.go文件中

3. gorm以及redis的埋点sdk设置见pkg下对应文件 

4. gin的埋点设置见main.go函数中

5. 手动埋点的示例见handle.go文件的mockTrace函数

6. grpc框架的客户端和服务端在pkg/grpcdemo/下



# 2. 术语说明
**Span**：一个节点在收到请求以及完成请求的过程是一个 `Span`，`Span` 记录了在这个过程中产生的各种信息。

**Trace**：一条`Trace`（调用链）可以被认为是一个由多个`Span`组成的有向无环图（DAG图）， `Span`与`Span`的关系被命名为`References`。

**Tracer**：`Tracer`接口用来创建`Span`，以及处理如何处理`Inject`(serialize) 和 `Extract` (deserialize)，用于跨进程边界传递。

**Context**:  `Context`是一个非常重要的概念，当我们需要跨服务传播trace数据时，可以在`Context`中存储spanID，traceID等信息，并随请求传输到另一个服务。

# 3. 相关sdk链接

1. [opentelemetry-go-instrument](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation)
2. [go-redis otel example](https://github.com/redis/go-redis/tree/master/example/otel)
3. [gorm otel](https://github.com/go-gorm/opentelemetry)
4. [otelgrpc](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc)

