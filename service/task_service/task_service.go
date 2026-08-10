package task_service

import (
	"encoding/json"
	"io"
	"lyg-collecter/common"
	"lyg-collecter/model"
	"lyg-collecter/scheduler"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type TaskService struct {
}

var schedulerInstance *scheduler.Scheduler

func InitScheduler(s *scheduler.Scheduler) {
	schedulerInstance = s
}

type CollectDataResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data []Zone `json:"data"`
}

type Zone struct {
	NodeId       int            `json:"nodeId"`
	ZoneId       int            `json:"zoneId"`
	NodeCode     string         `json:"nodeCode"`
	NodeName     string         `json:"nodeName"`
	OnlineStatus string         `json:"onlineStatus"`
	AlarmStatus  string         `json:"alermStatus"`
	InsPointCurr []InsPointCurr `json:"insPointCurr"`
}

type InsPointCurr struct {
	Id           int    `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	CurrVal      string `json:"currVal"`
	OnlineStatus string `json:"onlineStatus"`
	AlarmStatus  string `json:"alermStatus"`
	PosElement   string `json:"posElement"`
	Rhf          string `json:"rhf"`
}

func (t *TaskService) InitTask() error {
	tasks, err := model.FindTask()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		collectDataTask := model.Task{
			Name:   "采集数据任务",
			Cron:   "1 */1 * * *",
			Switch: true,
		}
		tasks = []model.Task{collectDataTask}
		err = model.AddTasks(tasks)
		if err != nil {
			return err
		}
	}
	if tasks[0].Switch == true {
		err = schedulerInstance.AddTask(&scheduler.Task{
			ID:       strconv.Itoa(int(tasks[0].ID)),
			Name:     tasks[0].Name,
			Schedule: tasks[0].Cron,
			JobFunc:  t.collectDataScheduleFunc,
		})
		if err != nil {
			zap.L().Error("添加采集数据定时任务失败", zap.Error(err))
			return err
		}

		zap.L().Info("采集数据定时任务已启动")
	}
	return nil
}

func (t *TaskService) collectDataScheduleFunc() {
	zap.L().Info("开始执行采集数据任务")

	dataTaskExecution := model.TaskExecution{
		TaskID: 1,
		Result: true,
	}
	resp, err := http.Get(common.UrlSetting.DataUrl)
	if err != nil {
		zap.L().Error("请求数据接口失败", zap.Error(err))
		dataTaskExecution.Result = false
		dataTaskExecution.Reason = err.Error()
	} else {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			zap.L().Error("读取数据接口返回数据失败", zap.Error(err))
			dataTaskExecution.Result = false
			dataTaskExecution.Reason = err.Error()
			return
		}

		// 解析 resp.Body 内容到 CollectDataResponse 结构体
		var collectDataResp CollectDataResponse
		if err := json.Unmarshal(body, &collectDataResp); err != nil {
			zap.L().Error("解析数据接口返回内容失败", zap.Error(err))
			dataTaskExecution.Result = false
			dataTaskExecution.Reason = err.Error()
		} else {
			// 成功解析后可以使用 collectDataResp 进行业务逻辑处理
			zap.L().Info("数据解析成功", zap.Any("data", collectDataResp))
			opticalFiber90 := model.OpticalFiber90{}
			// 获取当前时间年、月、日、小时
			year, month, day := time.Now().Date()
			hour := time.Now().Hour()
			opticalFiber90.Year = int16(year)
			opticalFiber90.Month = int8(month)
			opticalFiber90.Day = int8(day)
			opticalFiber90.Hour = int8(hour)
			//opticalFiber90.OnlineStatusBytes = make([]byte, common.StatusBytesLong)
			opticalFiber90.AlarmStatusBytes = make([]byte, common.StatusBytesLong)
			for _, insPointCurr := range collectDataResp.Data[0].InsPointCurr {
				changeErr := model.SetCodeField(&opticalFiber90, "Code"+insPointCurr.Code, insPointCurr.CurrVal)
				if changeErr != nil {
					zap.L().Error("设置Code字段失败", zap.Error(changeErr))
				}
				changeErr = model.SetCodeField(&opticalFiber90, "S"+insPointCurr.Code, insPointCurr.OnlineStatus)
				if changeErr != nil {
					zap.L().Error("设置S字段失败", zap.Error(changeErr))
				}
				// insPointCurr.Code转int
				code, changeErr := strconv.Atoi(insPointCurr.Code)
				if changeErr != nil {
					zap.L().Error("Code字段转换int失败", zap.Error(changeErr))
				}
				//if insPointCurr.OnlineStatus == "1" {
				//	opticalFiber90.OnlineStatusBytes[code/8] = opticalFiber90.OnlineStatusBytes[code/8] | (1 << (code % 8))
				//}
				// 二进制存储状态
				if insPointCurr.AlarmStatus == "1" {
					opticalFiber90.AlarmStatusBytes[code/8] = opticalFiber90.AlarmStatusBytes[code/8] | (1 << (code % 8))
				}
			}
			temp, windSpeed, humidity, pressure, err := getWeatherData()
			if err != nil {
				zap.L().Error("获取天气数据失败", zap.Error(err))
			} else {
				opticalFiber90.EnvTemperature = &temp
				opticalFiber90.WindSpeed = &windSpeed
				opticalFiber90.Humidity = &humidity
				opticalFiber90.Pressure = &pressure
			}
			err = model.CreateOpticalFiber90(&opticalFiber90)
			if err != nil {
				zap.L().Error("创建OpticalFiber90数据失败", zap.Error(err))
				dataTaskExecution.Result = false
				dataTaskExecution.Reason = err.Error()
			}
		}
	}
	err = model.CreateTaskExecution(dataTaskExecution)
	if err != nil {
		zap.L().Error("创建任务执行记录失败", zap.Error(err))
		return
	}

}

type WeatherResponse struct {
	Code    int `json:"code"`
	Nowinfo struct {
		Precipitation       float32 `json:"precipitation"`
		Temperature         float32 `json:"temperature"`
		Pressure            float32 `json:"pressure"`
		Humidity            float32 `json:"humidity"`
		WindDirection       string  `json:"windDirection"`
		WindDirectionDegree float32 `json:"windDirectionDegree"`
		WindSpeed           float32 `json:"windSpeed"`
	} `json:"nowinfo"`
}

func getWeatherData() (float32, float32, float32, float32, error) {
	resp, err := http.Get(common.UrlSetting.WeatherUrl)
	if err != nil {
		zap.L().Error("请求天气接口失败", zap.Error(err))
		return 0, 0, 0, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("读取天气接口返回数据失败", zap.Error(err))
		return 0, 0, 0, 0, err
	}
	zap.L().Info("天气接口返回内容", zap.String("weather data", string(body)))
	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		zap.L().Error("解析天气接口返回内容失败", zap.Error(err))
		return 0, 0, 0, 0, err
	}
	return weatherResp.Nowinfo.Temperature, weatherResp.Nowinfo.WindSpeed, weatherResp.Nowinfo.Humidity, weatherResp.Nowinfo.Pressure, nil
}
