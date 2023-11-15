package main

import (
	"go-opentelemetry-demo/pkg/grpcdemo"
	"time"
)

func main() {
	go grpcdemo.StartServer()
	time.Sleep(time.Second)
	grpcdemo.StartClient()
	time.Sleep(time.Second)
}
