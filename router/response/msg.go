package response

var MsgFlags = map[int]string{
	SUCCESS:                    "ok",
	ERROR:                      "fail",
	InvalidParams:              "服务端错误",
	ServerError:                "请求参数错误",
	ErrorAuthCheckTokenFail:    "token不合法",
	ErrorAuthCheckTokenTimeout: "token超时",
	ErrorGenerateToken:         "生成token失败",
	ErrorConnectWebSocket:      "连接websocket失败",
	ErrorReadWebSocket:         "读取websocket失败",

	ErrorFindInsulationMonitorData:                   "获取绝热异常及定位数据失败",
	ErrorFindHeatLossData:                            "获取热损数据失败",
	ErrorFindInsulationPerformanceData:               "获取绝热性能监控及预测数据失败",
	ErrorGetWeatherData:                              "获取天气数据失败",
	ErrorFindInsulationMonitorDataByYearMonthDayHour: "获取绝热异常及定位历史数据失败",
	ErrorFindOpticalFiber90ThreePointsData:           "获取截面三点历史数据失败",

	ErrorFindTile: "获取瓦片数据失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
