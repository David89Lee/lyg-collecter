package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "lyg-collecter/docs"
	v1 "lyg-collecter/router/api/v1"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	// 使用zap记录日志
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		zap.L().Info("HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiv1 := r.Group("/api/v1")

	apiv1.GET("/insulation-monitor-data", v1.FindInsulationMonitorData)
	apiv1.GET("/insulation-monitor-data/:year/:month/:day/:hour", v1.FindInsulationMonitorDataByYearMonthDayHour)
	apiv1.GET("/heat-loss-data", v1.FindHeatLossData)
	apiv1.GET("/insulation-performance-data", v1.FindInsulationPerformanceData)
	apiv1.GET("/weather", v1.GetLastWeather)
	apiv1.GET("/three-points-data", v1.FindOpticalFiber90ThreePointsData)

	apiv1.GET("/tile/:z/:x/:y", v1.FindTile)

	return r
}
