package common

import (
	"errors"
	"math"

	"go.uber.org/zap"
)

const (
	// 材料表面辐射系数，单位为瓦每平方米4次方开
	Dr = 0.85 * (5.67 * 0.00000001)
	// 保温外径，单位mm
	De = 1.15
)

func CalHeatLoss(temperatureCoefficient float32, temperature, envTemperature, windSpeed *float32) *float64 {
	// 打印temperature、envTemperature
	zap.L().Info("temperature", zap.Float32("temperature", *temperature))
	zap.L().Info("envTemperature", zap.Float32("envTemperature", *envTemperature))
	alpha := CalAlpha(temperatureCoefficient, temperature, envTemperature, windSpeed)
	if alpha == nil {
		return nil
	}

	// 对流放热系数
	//var alphaC = 8.9 * (math.Pow(float64(*windSpeed), 0.9) / math.Pow(De, 0.1))

	// 温度转换为K
	temperatureK := *temperature*temperatureCoefficient + 273.15
	envTemperatureK := *envTemperature*temperatureCoefficient + 273.15

	// 温度因子
	//var ar = 4 * math.Pow(float64((*temperature+*envTemperature)/2), 3)

	// 辐射放热系数
	//var alphaR = ar * Dr

	heatLoss := *alpha * float64(temperatureK-envTemperatureK)

	return &heatLoss
}

func CalAlpha(temperatureCoefficient float32, temperature, envTemperature, windSpeed *float32) *float64 {
	if temperature == nil || envTemperature == nil || windSpeed == nil {
		return nil
	}

	// 对流放热系数
	var alphaC = 8.9 * (math.Pow(float64(*windSpeed), 0.9) / math.Pow(De, 0.1))

	// 温度转换为K
	temperatureK := *temperature*temperatureCoefficient + 273.15
	envTemperatureK := *envTemperature*temperatureCoefficient + 273.15

	// 温度因子
	var ar = 4 * math.Pow(float64((temperatureK+envTemperatureK)/2), 3)

	// 辐射放热系数
	var alphaR = ar * Dr

	alpha := alphaR + alphaC

	zap.L().Info("alphaR", zap.Float64("alphaR", alphaR))
	zap.L().Info("alphaC", zap.Float64("alphaC", alphaC))

	return &alpha
}

func CalTotalHeat(heatLoss *float64, length int, diameter float32) float64 {
	return *heatLoss * float64(length) * math.Pi * float64(diameter)
}

// LinearRegression 最小二乘法拟合一次函数 y = ax + b，用于历史数据少于2个月时，预测未来数据
// 输入：x和y的数值切片（长度必须相等且大于1）
// 输出：斜率a、截距b，以及可能的错误
func LinearRegressionWithXY(x, y []float64) (a, b float64, err error) {
	// 校验输入数据合法性
	if len(x) != len(y) {
		return 0, 0, errors.New("x和y的长度必须相等")
	}
	n := len(x)
	if n < 2 {
		return 0, 0, errors.New("数据点数量至少需要2个")
	}

	// 计算所需的累加值
	var sumX, sumY, sumXY, sumX2 float64
	for i := 0; i < n; i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	// 计算分母（避免除以0）
	denominator := float64(n)*sumX2 - sumX*sumX
	if denominator == 0 {
		return 0, 0, errors.New("分母为0，无法计算（x值全部相同）")
	}

	// 计算斜率a
	a = (float64(n)*sumXY - sumX*sumY) / denominator

	// 计算截距b（b = 平均值y - a*平均值x）
	meanX := sumX / float64(n)
	meanY := sumY / float64(n)
	b = meanY - a*meanX

	return a, b, nil
}

// 历史数据大于2个月，x规律增长0,1,2,3...，不需要这个参数
func LinearRegressionWithY(y []*float64) (a, b float64, err error) {
	n := len(y)
	if n < 2 {
		return 0, 0, errors.New("数据点数量至少需要2个")
	}

	// 计算所需的累加值
	var sumX, sumX2 int
	var sumY, sumXY float64
	for i := 0; i < n; i++ {
		sumX += i
		sumY += *y[i]
		sumXY += float64(i) * *y[i]
		sumX2 += i * i
	}

	denominator := float64(n*sumX2 - sumX*sumX) // 不会等于0

	// 计算斜率a
	a = (float64(n)*sumXY - float64(sumX)*sumY) / denominator

	// 计算截距b（b = 平均值y - a*平均值x）
	meanX := sumX / n
	meanY := sumY / float64(n)
	b = meanY - a*float64(meanX)

	return a, b, nil
}
