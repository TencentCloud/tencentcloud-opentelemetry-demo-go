package gormclient

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var GormDB *gorm.DB

// InitGorm 初始化gorm客户端
func InitGorm() {

	var err error
	// 这里可以换成mysql，这里demo使用的sqlite
	GormDB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 添加tracing的插件
	if err = GormDB.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}
}
