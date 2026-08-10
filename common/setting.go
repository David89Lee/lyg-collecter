package common

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

type Server struct {
	HttpPort string
}

var ServerSetting = &Server{}

type Log struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

var LogSetting = &Log{}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

var DatabaseSetting = &Database{}

type Postgresql struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

var PostgresqlSetting = &Postgresql{}

type Url struct {
	DataUrl    string
	WeatherUrl string
}

var UrlSetting = &Url{}

type General struct {
	AbnormalTemperatureDifference float32
	StandardHeatLoss              float32
	Offset                        float32
	TemperatureCoefficients       []float32
}

var GeneralSetting = &General{}

var commonCfg *ini.File

func SetUpSetting(CfgFileDir string) {
	var err error

	commonCfgFile := fmt.Sprintf("%v/common.ini", CfgFileDir)

	commonCfg, err = ini.Load(commonCfgFile)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/common.ini': %v", err)
	}

	mapTo(commonCfg, "server", ServerSetting)
	mapTo(commonCfg, "log", LogSetting)
	mapTo(commonCfg, "database", DatabaseSetting)
	mapTo(commonCfg, "postgresql", PostgresqlSetting)
	mapTo(commonCfg, "url", UrlSetting)
	//mapTo(commonCfg, "general", GeneralSetting)

	generalSection := commonCfg.Section("general")
	GeneralSetting.AbnormalTemperatureDifference = float32(generalSection.Key("AbnormalTemperatureDifference").MustFloat64())
	GeneralSetting.StandardHeatLoss = float32(generalSection.Key("StandardHeatLoss").MustFloat64())
	GeneralSetting.Offset = float32(generalSection.Key("Offset").MustFloat64())

	// 手动解析 TemperatureCoefficients
	tempCoeffStr := generalSection.Key("TemperatureCoefficients").String()
	if tempCoeffStr != "" {
		coeffs := strings.Split(tempCoeffStr, ",")
		GeneralSetting.TemperatureCoefficients = make([]float32, 0, len(coeffs))
		for _, coeff := range coeffs {
			coeff = strings.TrimSpace(coeff)
			if val, err := strconv.ParseFloat(coeff, 32); err == nil {
				GeneralSetting.TemperatureCoefficients = append(GeneralSetting.TemperatureCoefficients, float32(val))
			} else {
				log.Printf("Warning: failed to parse temperature coefficient '%s': %v", coeff, err)
			}
		}
	}
}

func mapTo(cfg *ini.File, section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
