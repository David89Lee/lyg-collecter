package response

const (
	SUCCESS                    = 200
	ERROR                      = 500
	InvalidParams              = 400
	ServerError                = 600
	ErrorAuthCheckTokenFail    = 601
	ErrorAuthCheckTokenTimeout = 602
	ErrorGenerateToken         = 603
	ErrorConnectWebSocket      = 610
	ErrorReadWebSocket         = 611

	ErrorFindInsulationMonitorData                   = 10000
	ErrorFindHeatLossData                            = 10001
	ErrorFindInsulationPerformanceData               = 10002
	ErrorGetWeatherData                              = 10003
	ErrorFindInsulationMonitorDataByYearMonthDayHour = 10004
	ErrorFindOpticalFiber90ThreePointsData           = 10005
	ErrorExportOpticalFiber90ThreePointsDataExcel    = 10006

	ErrorFindTile = 10010
)
