package main

import (
	"lyg-collecter/common"
	"lyg-collecter/middleware"
	"lyg-collecter/model"
	"lyg-collecter/router"
	"lyg-collecter/scheduler"
	"lyg-collecter/service/task_service"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var s *scheduler.Scheduler

func init() {
	common.SetUpSetting("conf")
	middleware.InitLogger()
	model.Setup()

	// 初始化定时任务
	zap.L().Info("初始化定时任务……")
	s = scheduler.NewScheduler()
	task_service.InitScheduler(s)
	taskService := task_service.TaskService{}
	err := taskService.InitTask()
	if err != nil {
		zap.L().Error("初始化定时任务失败", zap.Error(err))
		panic(err)
	}
	zap.L().Info("初始化定时任务完成")
}

// @title 连云港二期
// @version 1.0
// @description 连云港二期
// @contact.name LYG
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:8080
func main() {
	defer func(l *zap.Logger) {
		err := l.Sync()
		if err != nil {
			zap.L().Error("zap日志同步失败", zap.Error(err))
		}
	}(zap.L())

	s.Start()

	r := router.InitRouter()

	// start server in a goroutine
	go func() {
		if err := r.Run(common.ServerSetting.HttpPort); err != nil {
			zap.L().Error("启动服务失败", zap.Error(err))
		}
	}()

	// 处理ctrl+c信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	s.Stop()

	zap.L().Info("Shutting down gracefully...")
}
