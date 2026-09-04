package optical_fiber_90_service

import (
	"fmt"
	"lyg-collecter/common"
	"lyg-collecter/model"
	"math"
	"reflect"
	"time"

	"go.uber.org/zap"
)

const historyMonthNum = 5
const showMonthNum = 10

type OpticalFiber90Service struct {
}

type InsulationMonitorData struct {
	MonitorPoints1      []string
	MonitorPoints2      []string
	MonitorPoints3      []string
	Temperatures1       []*float32
	Temperatures2       []*float32
	Temperatures3       []*float32
	Anomalies1          []bool
	Anomalies2          []bool
	Anomalies3          []bool
	HeatLosses1         []*float64
	HeatLosses2         []*float64
	HeatLosses3         []*float64
	DeviceExceptions1   []bool
	DeviceExceptions2   []bool
	DeviceExceptions3   []bool
	AbnormalTemperature float32
	EnvTemperature      float32
	WindSpeed           float32
}

func (t *OpticalFiber90Service) FindInsulationMonitorData() (*InsulationMonitorData, error) {
	opticalFiber90, err := model.GetLastOpticalFiber90()
	if err != nil {
		return nil, err
	}
	return t.FindInsulationMonitorDataByOpticalFiber90(opticalFiber90)
}

func (t *OpticalFiber90Service) FindInsulationMonitorDataByYearMonthDayHour(year, month, day, hour string) (*InsulationMonitorData, error) {
	opticalFiber90, err := model.GetOpticalFiber90ByYearMonthDayHour(year, month, day, hour)
	if err != nil {
		return nil, err
	}
	return t.FindInsulationMonitorDataByOpticalFiber90(opticalFiber90)
}

func (t *OpticalFiber90Service) FindInsulationMonitorDataByOpticalFiber90(opticalFiber90 *model.OpticalFiber90) (*InsulationMonitorData, error) {
	//opticalFiber90, err := model.GetLastOpticalFiber90()
	//if err != nil {
	//	return nil, err
	//}
	// 目前用裸露出的光纤点作为环境温度
	//opticalFiber90.EnvTemperature = opticalFiber90.Code2

	// opticalFiber90.EnvTemperature和opticalFiber90.Code2取大的，大的更接近测点温度，热损更低
	if *opticalFiber90.EnvTemperature < *opticalFiber90.Code2 {
		opticalFiber90.EnvTemperature = opticalFiber90.Code2
	}

	// 初始化 insulationMonitorData
	insulationMonitorData := &InsulationMonitorData{}

	// 有三段，分三次处理，以高位为基准根据管道长度往前推，默认1米一个测点
	var position float32 = 0.0
	// 第一段22~123
	// 第一段改为 5~106
	for i := common.EffEndPoint - common.PipeLength*3 + 1; i <= common.EffEndPoint-common.PipeLength*2; i++ {
		position++
		insulationMonitorData.MonitorPoints1 = append(insulationMonitorData.MonitorPoints1, FindPipeSection(position))
		// 用反射获取opticalFiber90的非指针属性值s%d
		exception := reflect.ValueOf(opticalFiber90).Elem().FieldByName(fmt.Sprintf("S%d", i)).Interface().(int8) == 0
		insulationMonitorData.DeviceExceptions1 = append(insulationMonitorData.DeviceExceptions1, exception)
		if exception {
			insulationMonitorData.Temperatures1 = append(insulationMonitorData.Temperatures1, new(float32))
			insulationMonitorData.Anomalies1 = append(insulationMonitorData.Anomalies1, false)
			insulationMonitorData.HeatLosses1 = append(insulationMonitorData.HeatLosses1, new(float64))
			continue
		}
		// 用反射获取opticalFiber90的属性值
		field := reflect.ValueOf(opticalFiber90).Elem().FieldByName(fmt.Sprintf("Code%d", i)).Interface().(*float32)
		var offsetVal float32
		if field != nil {
			offsetVal = *field + common.GeneralSetting.Offset
		}
		insulationMonitorData.Temperatures1 = append(insulationMonitorData.Temperatures1, &offsetVal)
		insulationMonitorData.Anomalies1 = append(insulationMonitorData.Anomalies1, offsetVal-*opticalFiber90.EnvTemperature > common.GeneralSetting.AbnormalTemperatureDifference)
		insulationMonitorData.HeatLosses1 = append(insulationMonitorData.HeatLosses1, common.CalHeatLoss(1.0, &offsetVal, opticalFiber90.EnvTemperature, opticalFiber90.WindSpeed))
	}
	position = 0.0
	// 第二段208~107
	for i := common.EffEndPoint - common.PipeLength; i >= common.EffEndPoint-common.PipeLength*2+1; i-- {
		position++
		insulationMonitorData.MonitorPoints2 = append(insulationMonitorData.MonitorPoints2, FindPipeSection(position))
		// 用反射获取opticalFiber90的非指针属性值s%d
		exception := reflect.ValueOf(opticalFiber90).Elem().FieldByName(fmt.Sprintf("S%d", i)).Interface().(int8) == 0
		insulationMonitorData.DeviceExceptions2 = append(insulationMonitorData.DeviceExceptions2, exception)
		if exception {
			insulationMonitorData.Temperatures2 = append(insulationMonitorData.Temperatures2, new(float32))
			insulationMonitorData.Anomalies2 = append(insulationMonitorData.Anomalies2, false)
			insulationMonitorData.HeatLosses2 = append(insulationMonitorData.HeatLosses2, new(float64))
			continue
		}
		// 用反射获取opticalFiber90的属性值
		field := reflect.ValueOf(opticalFiber90).Elem().FieldByName(fmt.Sprintf("Code%d", i)).Interface().(*float32)
		var offsetVal float32
		if field != nil {
			offsetVal = *field + common.GeneralSetting.Offset
		}
		insulationMonitorData.Temperatures2 = append(insulationMonitorData.Temperatures2, &offsetVal)
		insulationMonitorData.Anomalies2 = append(insulationMonitorData.Anomalies2, offsetVal-*opticalFiber90.EnvTemperature > common.GeneralSetting.AbnormalTemperatureDifference)
		insulationMonitorData.HeatLosses2 = append(insulationMonitorData.HeatLosses2, common.CalHeatLoss(1.0, &offsetVal, opticalFiber90.EnvTemperature, opticalFiber90.WindSpeed))
	}
	position = 0.0
	// 第三段209~310， 314~327暂时异常
	for i := common.EffEndPoint - common.PipeLength + 1; i <= common.EffEndPoint; i++ {
		position++
		insulationMonitorData.MonitorPoints3 = append(insulationMonitorData.MonitorPoints3, FindPipeSection(position))
		// 用反射获取opticalFiber90的非指针属性值s%d
		exception := reflect.ValueOf(opticalFiber90).Elem().FieldByName(fmt.Sprintf("S%d", i)).Interface().(int8) == 0
		insulationMonitorData.DeviceExceptions3 = append(insulationMonitorData.DeviceExceptions3, exception)
		if exception {
			insulationMonitorData.Temperatures3 = append(insulationMonitorData.Temperatures3, new(float32))
			insulationMonitorData.Anomalies3 = append(insulationMonitorData.Anomalies3, false)
			insulationMonitorData.HeatLosses3 = append(insulationMonitorData.HeatLosses3, new(float64))
			continue
		}
		// 用反射获取opticalFiber90的属性值
		field := reflect.ValueOf(opticalFiber90).Elem().FieldByName(fmt.Sprintf("Code%d", i)).Interface().(*float32)
		var offsetVal float32
		if field != nil {
			offsetVal = *field + common.GeneralSetting.Offset
		}
		insulationMonitorData.Temperatures3 = append(insulationMonitorData.Temperatures3, &offsetVal)
		insulationMonitorData.Anomalies3 = append(insulationMonitorData.Anomalies3, offsetVal-*opticalFiber90.EnvTemperature > common.GeneralSetting.AbnormalTemperatureDifference)
		insulationMonitorData.HeatLosses3 = append(insulationMonitorData.HeatLosses3, common.CalHeatLoss(1.0, &offsetVal, opticalFiber90.EnvTemperature, opticalFiber90.WindSpeed))
	}
	insulationMonitorData.AbnormalTemperature = *opticalFiber90.EnvTemperature + common.GeneralSetting.AbnormalTemperatureDifference
	insulationMonitorData.EnvTemperature = *opticalFiber90.EnvTemperature
	insulationMonitorData.WindSpeed = *opticalFiber90.WindSpeed
	return insulationMonitorData, nil
}

// 计算管墩号区间
func FindPipeSection(position float32) string {
	switch {
	case position < 3.4:
		return "- G4-243墩子间"
	case position >= 3.4 && position < 18.4:
		return "G4-243墩子 - G4-242墩子间"
	case position >= 18.4 && position < 33.4:
		return "G4-242墩子 - G4-241墩子间"
	case position >= 33.4 && position < 37.4:
		return "G4-241墩子 - G4-240墩子间"
	case position >= 37.4 && position < 52.4:
		return "G4-240墩子 - G4-239墩子间"
	case position >= 52.4 && position < 56.4:
		return "G4-239墩子 - G4-238墩子间"
	case position >= 56.4 && position < 71.4:
		return "G4-238墩子 - G4-237墩子间"
	case position >= 71.4 && position < 86.4:
		return "G4-237墩子 - G4-236墩子间"
	case position >= 86.4 && position < 99.4:
		return "G4-236墩子 - G4-235墩子间"
	case position >= 99.4:
		return "G4-235墩子 -"
	}
	return ""
}

type PointsHistoryData struct {
	Temperatures1 []*float32
	Temperatures2 []*float32
	Temperatures3 []*float32
}

type HeatLossData struct {
	PipeSection string
	HeatLoss    *float64
	Standard    float32
	TotalHeat   float64
}

func (t *OpticalFiber90Service) FindHeatLossData() ([]*HeatLossData, error) {
	var heatLossData HeatLossData
	heatLossData.PipeSection = "瑞恒线"
	avg, env, windSpeed, err := model.GetLastOpticalFiber90AvgAndEnvTemperature()
	if err != nil {
		return nil, err
	}
	if avg != nil {
		*avg = *avg + common.GeneralSetting.Offset
	}
	heatLossData.HeatLoss = common.CalHeatLoss(1.0, avg, env, windSpeed)
	heatLossData.Standard = common.GeneralSetting.StandardHeatLoss
	heatLossData.TotalHeat = common.CalTotalHeat(heatLossData.HeatLoss, 102, 1.15)
	return []*HeatLossData{&heatLossData}, nil
}

type InsulationPerformanceData struct {
	Dates             []int
	StandardLoss      []float32
	ActualLoss        []*float64
	PredictLoss       []*float64
	PredictedLifespan *float64
}

func (t *OpticalFiber90Service) FindInsulationPerformanceData() (*InsulationPerformanceData, error) {
	// 获取每个月的平均值
	avgDataResults, err := model.FindOpticalFiber90MonthAvg(historyMonthNum)
	if err != nil {
		return nil, err
	}
	var insulationPerformanceData InsulationPerformanceData

	var offsetVal float32
	for _, avgDataResult := range avgDataResults {
		// avgDataResult 中的月份设置到Dates中
		insulationPerformanceData.Dates = append(insulationPerformanceData.Dates, avgDataResult.Month)
		// 设置标准值
		insulationPerformanceData.StandardLoss = append(insulationPerformanceData.StandardLoss, common.GeneralSetting.StandardHeatLoss)
		// 计算实际损耗并设置
		if avgDataResult.Temperature != nil {
			offsetVal = *avgDataResult.Temperature + common.GeneralSetting.Offset
		}
		zap.L().Info("Month, Coefficient and temperature",
			zap.Int("Month", avgDataResult.Month),
			zap.Float32("Coefficient", common.GeneralSetting.TemperatureCoefficients[avgDataResult.Month-1]),
			zap.Float32("Temperature", offsetVal))
		heatLoss := common.CalHeatLoss(common.GeneralSetting.TemperatureCoefficients[avgDataResult.Month-1], &offsetVal, avgDataResult.EnvTemperature, avgDataResult.WindSpeed)
		insulationPerformanceData.ActualLoss = append(insulationPerformanceData.ActualLoss, heatLoss)
		// 预测折线前面部分一样
		insulationPerformanceData.PredictLoss = append(insulationPerformanceData.PredictLoss, heatLoss)

	}
	length := len(avgDataResults)
	// 拟合公式，为了生成预测折线y=ax+b
	var a, b float64
	hasPredictLossData := true
	// 当历史数据小于两个月，用每个小时的数据去拟合
	if length < 2 {
		zap.L().Info("历史数据不足2个月，用每个小时数据去拟合")
		// 获取最近2个月小时数据
		hourAvgDataResults, err := model.FindOpticalFiber90HourAvg(2)
		if err != nil {
			return nil, err
		}
		if len(hourAvgDataResults) < 2 { // 如果小时数据都没有2条，预测折线后面部分不设置
			hasPredictLossData = false
		} else {
			// 生成小时的x,y合集
			var x, y []float64
			// 跟两个月前第一个小时比小时的差值，计算x（相对月份的比例），以防数据缺失
			firstDate := time.Date(hourAvgDataResults[0].Year, time.Month(hourAvgDataResults[0].Month), 1, 0, 0, 0, 0, time.Local)
			// firstDate月份的总小时数
			firstMonthHours := firstDate.AddDate(0, 1, 0).Sub(firstDate).Hours()
			secondDate := time.Date(hourAvgDataResults[0].Year, time.Month(hourAvgDataResults[0].Month), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
			// secondDate月份的总小时数
			secondMonthHours := secondDate.AddDate(0, 1, 0).Sub(secondDate).Hours()
			for _, hourAvgDataResult := range hourAvgDataResults {
				date := time.Date(hourAvgDataResult.Year, time.Month(hourAvgDataResult.Month), hourAvgDataResult.Day, hourAvgDataResult.Hour, 0, 0, 0, time.Local)
				// 如果是第一个月，(date.Hour()-firstDate.Hour())/firstMonthHours
				if date.Month() == firstDate.Month() {
					x = append(x, date.Sub(firstDate).Hours()/firstMonthHours)
				} else {
					// 如果是第二个月，在计算当月小时比例的基础上再加1
					x = append(x, 1+date.Sub(secondDate).Hours()/secondMonthHours)
				}
				if hourAvgDataResult.Temperature != nil {
					offsetVal = *hourAvgDataResult.Temperature + common.GeneralSetting.Offset
				}
				heatLoss := common.CalHeatLoss(1.0, &offsetVal, hourAvgDataResult.EnvTemperature, hourAvgDataResult.WindSpeed)
				if heatLoss != nil {
					y = append(y, *heatLoss)
				} else {
					// 处理 heatLoss 为 nil 的情况，例如跳过该数据点或赋予默认值
					zap.L().Warn("heatLoss is nil, skipping this data point")
					y = append(y, 0)
				}
			}
			a, b, err = common.LinearRegressionWithXY(x, y)
			if err != nil {
				return nil, err
			}
		}
	} else { // 大于两个月的时候直接用月的数据去拟合
		a, b, err = common.LinearRegressionWithY(insulationPerformanceData.ActualLoss)
	}

	// 打印a、b
	zap.L().Info("a,b", zap.Float64("a", a), zap.Float64("b", b))

	// 补齐10个月
	curDate := time.Date(avgDataResults[length-1].Year, time.Month(avgDataResults[length-1].Month+1), 1, 0, 0, 0, 0, time.Local)
	for i := range showMonthNum - length {
		// 补齐日期
		insulationPerformanceData.Dates = append(insulationPerformanceData.Dates, int(curDate.AddDate(0, i, 0).Month()))
		// 补齐标准值
		insulationPerformanceData.StandardLoss = append(insulationPerformanceData.StandardLoss, common.GeneralSetting.StandardHeatLoss)
		if hasPredictLossData {
			// 补齐预测折线y
			predictValue := a*float64(i+length) + b
			insulationPerformanceData.PredictLoss = append(insulationPerformanceData.PredictLoss, &predictValue)
		}
	}

	if hasPredictLossData {
		// 计算预测寿命PredictedLifespan，y=StandardHeatLoss倒推x
		var predictedLifespan *float64
		if a == 0 {
			predictedLifespan = nil
		} else {
			p := (float64(common.GeneralSetting.StandardHeatLoss)-b)/a - float64(length-1)
			predictedLifespan = &p
		}
		insulationPerformanceData.PredictedLifespan = predictedLifespan
	} else {
		insulationPerformanceData.PredictedLifespan = nil
	}
	return &insulationPerformanceData, nil
}

type InsulationPerformanceUnit2 struct {
	Year              int // 当前年
	Loss              *float64
	PredictedLifespan *float64
}

func (t *OpticalFiber90Service) FindInsulationPerformanceData2() ([]InsulationPerformanceUnit2, error) {
	// 获取每年的平均值
	yearAvgDataResults, err := model.FindOpticalFiber90YearAvg()
	if err != nil {
		return nil, err
	}

	var insulationPerformanceData = make([]InsulationPerformanceUnit2, 0)

	var y []*float64

	// 如果只有一年的数据，由于是26年3月份开始采集的，取3、4月份作为25年数据
	if len(yearAvgDataResults) == 1 {
		year2025AvgData, err := model.FindOpticalFiber90AvgByYearMonths(2026, []int{3, 4})
		if err != nil {
			return nil, err
		}
		// year2025AvgData加到yearAvgDataResults的开头
		yearAvgDataResults = append([]*model.YearAvgDataResult{year2025AvgData}, yearAvgDataResults...)
	}

	for _, yearAvgDataResult := range yearAvgDataResults {
		// 计算热损
		heatLoss := common.CalHeatLoss(1, yearAvgDataResult.Temperature, yearAvgDataResult.EnvTemperature, yearAvgDataResult.WindSpeed)
		y = append(y, heatLoss)
		insulationPerformanceData = append(insulationPerformanceData, InsulationPerformanceUnit2{
			Year:              yearAvgDataResult.Year,
			Loss:              heatLoss,
			PredictedLifespan: nil,
		})
	}
	a, b, err := common.LinearRegressionWithY(y)
	// 打印a、b
	zap.L().Info("a,b", zap.Float64("a", a), zap.Float64("b", b))

	// 计算预测寿命PredictedLifespan，y=StandardHeatLoss倒推x
	var predictedLifespan *float64
	var strandardHeatLoss = float64(common.GeneralSetting.StandardHeatLoss)
	lifespan := float64(30)
	if a == 0 { // a=0说明每年热损都是一样的，预测寿命无穷大，设置封顶30年
		predictedLifespan = &lifespan
	} else {
		//p := (strandardHeatLoss - b) / a
		//predictedLifespan = &p
		predictedLifespan = &lifespan // 根据客户要求，第一年预测寿命30年
	}
	for i := range insulationPerformanceData {
		lifespan := *predictedLifespan - float64(i)
		insulationPerformanceData[i].PredictedLifespan = &lifespan
	}
	// 最后加上热损到标准值时候的点，在x轴上
	zeroLifespan := 0.0
	insulationPerformanceData = append(insulationPerformanceData, InsulationPerformanceUnit2{
		Year:              int(2025 + math.Ceil(*predictedLifespan)),
		Loss:              &strandardHeatLoss,
		PredictedLifespan: &zeroLifespan,
	})
	// 这边不按热损排序了，前端按热损插值x轴的位置，最后一个值用虚线
	return insulationPerformanceData, nil
}

type WeatherData struct {
	Temperature float32
	WindSpeed   float32
	Humidity    float32
}

func (t *OpticalFiber90Service) GetLastWeather() (*WeatherData, error) {
	opticalFiber90, err := model.GetLastOpticalFiber90()
	if err != nil {
		return nil, err
	}

	var temperature float32
	if *opticalFiber90.EnvTemperature > *opticalFiber90.Code2 {
		temperature = *opticalFiber90.EnvTemperature
	} else {
		temperature = *opticalFiber90.Code2
	}

	return &WeatherData{
		Temperature: temperature,
		WindSpeed:   *opticalFiber90.WindSpeed,
		Humidity:    *opticalFiber90.Humidity,
	}, nil
}

type ThreePointsDataList struct {
	Times               []string
	WindSpeed           []*float32
	EnvTemperature      []*float32
	AbnormalTemperature []*float32
	C1                  []*float32
	C2                  []*float32
	C3                  []*float32
}

func (t *OpticalFiber90Service) FindOpticalFiber90ThreePointsData(i int, year1, month1, day1, hour1, year2, month2, day2, hour2 string) (*ThreePointsDataList, error) {
	// 根据i生成code和s，作为查询字段
	c1 := fmt.Sprintf("%s%d", "code", common.EffEndPoint-common.PipeLength*3+1+i)
	s1 := fmt.Sprintf("%s%d", "s", common.EffEndPoint-common.PipeLength*3+1+i)
	c2 := fmt.Sprintf("%s%d", "code", common.EffEndPoint-common.PipeLength-i)
	s2 := fmt.Sprintf("%s%d", "s", common.EffEndPoint-common.PipeLength-i)
	c3 := fmt.Sprintf("%s%d", "code", common.EffEndPoint-common.PipeLength+1+i)
	s3 := fmt.Sprintf("%s%d", "s", common.EffEndPoint-common.PipeLength+1+i)

	var threePointsDatas []*model.ThreePointsData

	// 检查截止时间有没有记录
	id2, err := model.GetOpticalFiber90IDByYearMonthDayHour(year2, month2, day2, hour2)
	if err != nil {
		return nil, err
	} else {
		if id2 == 0 {
			threePointsDatas, err = model.FindOpticalFiber90ThreePointsDataOpen(c1, c2, c3, s1, s2, s3, year1, month1, day1, hour1)
		} else {
			threePointsDatas, err = model.FindOpticalFiber90ThreePointsDataClosed(c1, c2, c3, s1, s2, s3, year1, month1, day1, hour1, id2)
			if err != nil {
				return nil, err
			}
		}

	}

	var threePointsDataList = &ThreePointsDataList{
		Times:               make([]string, 0),
		WindSpeed:           make([]*float32, 0),
		EnvTemperature:      make([]*float32, 0),
		AbnormalTemperature: make([]*float32, 0),
		C1:                  make([]*float32, 0),
		C2:                  make([]*float32, 0),
		C3:                  make([]*float32, 0),
	}

	for _, threePointsData := range threePointsDatas {
		t := fmt.Sprintf("%d.%d.%d %d时", threePointsData.Year, threePointsData.Month, threePointsData.Day, threePointsData.Hour)
		threePointsDataList.Times = append(threePointsDataList.Times, t)
		threePointsDataList.WindSpeed = append(threePointsDataList.WindSpeed, threePointsData.WindSpeed)
		threePointsDataList.EnvTemperature = append(threePointsDataList.EnvTemperature, threePointsData.EnvTemperature)
		if threePointsData.EnvTemperature == nil {
			threePointsDataList.AbnormalTemperature = append(threePointsDataList.AbnormalTemperature, nil)
		} else {
			diff := *threePointsData.EnvTemperature + common.GeneralSetting.AbnormalTemperatureDifference
			threePointsDataList.AbnormalTemperature = append(threePointsDataList.AbnormalTemperature, &diff)
		}
		threePointsDataList.C1 = append(threePointsDataList.C1, threePointsData.C1)
		threePointsDataList.C2 = append(threePointsDataList.C2, threePointsData.C2)
		threePointsDataList.C3 = append(threePointsDataList.C3, threePointsData.C3)
	}
	return threePointsDataList, nil
}
