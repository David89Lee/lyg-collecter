package v1

import (
	"fmt"
	"lyg-collecter/router/response"
	"lyg-collecter/service/optical_fiber_90_service"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary 获取绝热异常及定位数据
// @Tags 光纤测温
// @Produce  json
// @Success 200 {object} response.Response{data=optical_fiber_90_service.InsulationMonitorData} "成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response
// @Router /api/v1/insulation-monitor-data [get]
func FindInsulationMonitorData(c *gin.Context) {
	resG := response.Gin{C: c}
	service := optical_fiber_90_service.OpticalFiber90Service{}

	data, err := service.FindInsulationMonitorData()
	if err != nil {
		zap.L().Error("获取绝热异常及定位数据失败", zap.Error(err))
		resG.Response(http.StatusOK, response.ErrorFindInsulationMonitorData, nil)
		return
	}
	resG.Response(http.StatusOK, response.SUCCESS, data)
}

// @Summary 测试参数
// @Tags 光纤测温
// @Produce  json
// @Success 200 {object} response.Response{data=optical_fiber_90_service.InsulationMonitorData} "成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response
// @Router /api/v1/insulation-monitor-data/{year}/{month}/{day}/{hour} [get]
func FindInsulationMonitorDataByYearMonthDayHour(c *gin.Context) {
	resG := response.Gin{C: c}
	// 打印c
	zap.L().With(zap.Any("c", c)).Info("c")
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	dayStr := c.Query("day")
	hourStr := c.Query("hour")
	zap.L().With(zap.String("year", yearStr), zap.String("month", monthStr), zap.String("day", dayStr), zap.String("hour", hourStr)).Info("query")

	service := optical_fiber_90_service.OpticalFiber90Service{}
	data, err := service.FindInsulationMonitorDataByYearMonthDayHour(yearStr, monthStr, dayStr, hourStr)
	if err != nil {
		zap.L().Error("获取绝热异常及定位历史数据失败", zap.Error(err))
		resG.Response(http.StatusOK, response.ErrorFindInsulationMonitorDataByYearMonthDayHour, nil)
		return
	}
	resG.Response(http.StatusOK, response.SUCCESS, data)
}

// @Summary 获取管道散热损失分析
// @Tags 光纤测温
// @Produce  json
// @Success 200 {object} response.Response{data=optical_fiber_90_service.HeatLossData} "成功"
// @Failure 500 {object} response.Response
// @Router /api/v1/heat-loss-data [get]
func FindHeatLossData(c *gin.Context) {
	resG := response.Gin{C: c}
	service := optical_fiber_90_service.OpticalFiber90Service{}

	data, err := service.FindHeatLossData()
	if err != nil {
		zap.L().Error("获取管道散热损失分析失败", zap.Error(err))
		resG.Response(http.StatusOK, response.ErrorFindHeatLossData, nil)
		return
	}
	resG.Response(http.StatusOK, response.SUCCESS, data)
}

// @Summary 获取绝热性能监控及预测
// @Tags 光纤测温
// @Produce  json
// @Success 200 {object} response.Response{data=optical_fiber_90_service.InsulationPerformanceData} "成功"
// @Failure 500 {object} response.Response
// @Router /api/v1/insulation-performance-data [get]
func FindInsulationPerformanceData(c *gin.Context) {
	resG := response.Gin{C: c}
	service := optical_fiber_90_service.OpticalFiber90Service{}

	//data, err := service.FindInsulationPerformanceData()
	data, err := service.FindInsulationPerformanceData2()
	if err != nil {
		zap.L().Error("获取绝热性能监控及预测失败", zap.Error(err))
		resG.Response(http.StatusOK, response.ErrorFindInsulationPerformanceData, nil)
		return
	}
	resG.Response(http.StatusOK, response.SUCCESS, data)
}

// @Summary 获取天气数据
// @Tags 光纤测温
// @Produce  json
// @Success 200 {object} response.Response{data=optical_fiber_90_service.WeatherData} "成功"
// @Failure 500 {object} response.Response
// @Router /api/v1/weather [get]
func GetLastWeather(c *gin.Context) {
	resG := response.Gin{C: c}
	service := optical_fiber_90_service.OpticalFiber90Service{}

	data, err := service.GetLastWeather()
	if err != nil {
		zap.L().Error("获取天气数据失败", zap.Error(err))
		resG.Response(http.StatusOK, response.ErrorGetWeatherData, nil)
		return
	}
	resG.Response(http.StatusOK, response.SUCCESS, data)
}

// @Summary 获取界面上三个点历史曲线数据
// @Tags 光纤测温
// @Produce  json
// @Success 200 {object} response.Response{data=optical_fiber_90_service.ThreePointsDataList} "成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response
// @Router /api/v1/three-points-data [get]
func FindOpticalFiber90ThreePointsData(c *gin.Context) {
	resG := response.Gin{C: c}
	iStr := c.Query("i")
	year1Str := c.Query("year1")
	month1Str := c.Query("month1")
	day1Str := c.Query("day1")
	hour1Str := c.Query("hour1")
	year2Str := c.Query("year2")
	month2Str := c.Query("month2")
	day2Str := c.Query("day2")
	hour2Str := c.Query("hour2")
	i, err := strconv.Atoi(iStr)
	if err != nil {
		zap.L().Error("参数i转换int失败", zap.Error(err))
		resG.Response(http.StatusOK, response.InvalidParams, nil)
		return
	}
	service := optical_fiber_90_service.OpticalFiber90Service{}
	data, err := service.FindOpticalFiber90ThreePointsData(i, year1Str, month1Str, day1Str, hour1Str, year2Str, month2Str, day2Str, hour2Str)
	if err != nil {
		zap.L().Error("获取界面上三个点历史曲线数据失败", zap.Error(err))
		resG.Response(http.StatusOK, response.ErrorFindOpticalFiber90ThreePointsData, nil)
		return
	}
	resG.Response(http.StatusOK, response.SUCCESS, data)
}

// @Summary 导出界面上三个点历史曲线数据为 Excel
// @Tags 光纤测温
// @Produce octet-stream
// @Success 200 {file} binary "xlsx 文件"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response
// @Router /api/v1/three-points-data/export [get]
func ExportOpticalFiber90ThreePointsDataExcel(c *gin.Context) {
	resG := response.Gin{C: c}
	iStr := c.Query("i")
	year1Str := c.Query("year1")
	month1Str := c.Query("month1")
	day1Str := c.Query("day1")
	hour1Str := c.Query("hour1")
	year2Str := c.Query("year2")
	month2Str := c.Query("month2")
	day2Str := c.Query("day2")
	hour2Str := c.Query("hour2")
	i, err := strconv.Atoi(iStr)
	if err != nil {
		zap.L().Error("参数i转换int失败", zap.Error(err))
		resG.Response(http.StatusOK, response.InvalidParams, nil)
		return
	}
	service := optical_fiber_90_service.OpticalFiber90Service{}
	data, filename, err := service.ExportOpticalFiber90ThreePointsDataExcel(i, year1Str, month1Str, day1Str, hour1Str, year2Str, month2Str, day2Str, hour2Str)
	if err != nil {
		zap.L().Error("导出界面上三个点历史曲线数据失败", zap.Error(err))
		resG.Response(http.StatusOK, response.ErrorExportOpticalFiber90ThreePointsDataExcel, nil)
		return
	}
	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(filename)))
	c.Data(http.StatusOK, contentType, data)
}
