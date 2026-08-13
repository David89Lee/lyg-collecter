package model

import (
	"fmt"
	"lyg-collecter/common"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type OpticalFiber90 struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	Year      int16 `gorm:"type:smallint comment '年份';index:idx_date"`
	Month     int8  `gorm:"type:tinyint comment '月份';index:idx_date"`
	Day       int8  `gorm:"type:tinyint comment '日期';index:idx_date"`
	Hour      int8  `gorm:"type:tinyint comment '小时';index:idx_date"`
	//OnlineStatusBytes []byte   `gorm:"type binary(41) comment '在线状态'"`
	AlarmStatusBytes []byte   `gorm:"type binary(41) comment '报警状态'"`
	EnvTemperature   *float32 `gorm:"comment '环境温度'"`
	WindSpeed        *float32 `gorm:"comment '风速'"`
	Humidity         *float32 `gorm:"comment '湿度'"`
	Pressure         *float32 `gorm:"comment '气压'"`
	Code2            *float32 `gorm:"comment '测点2'"`
	Code3            *float32 `gorm:"comment '测点3'"`
	Code4, Code5, Code6, Code7, Code8, Code9, Code10, Code11, Code12, Code13, Code14, Code15, Code16, Code17, Code18, Code19, Code20, Code21, Code22, Code23, Code24, Code25, Code26, Code27, Code28, Code29, Code30, Code31, Code32, Code33, Code34, Code35,
	Code36, Code37, Code38, Code39, Code40, Code41, Code42, Code43, Code44, Code45, Code46, Code47, Code48, Code49, Code50, Code51, Code52, Code53, Code54, Code55, Code56, Code57, Code58, Code59, Code60, Code61, Code62, Code63, Code64, Code65, Code66, Code67, Code68,
	Code69, Code70, Code71, Code72, Code73, Code74, Code75, Code76, Code77, Code78, Code79, Code80, Code81, Code82, Code83, Code84, Code85, Code86, Code87, Code88, Code89, Code90, Code91, Code92, Code93, Code94, Code95, Code96, Code97, Code98, Code99, Code100, Code101,
	Code102, Code103, Code104, Code105, Code106, Code107, Code108, Code109, Code110, Code111, Code112, Code113, Code114, Code115, Code116, Code117, Code118, Code119, Code120, Code121, Code122, Code123, Code124, Code125, Code126, Code127, Code128, Code129, Code130, Code131,
	Code132, Code133, Code134, Code135, Code136, Code137, Code138, Code139, Code140, Code141, Code142, Code143, Code144, Code145, Code146, Code147, Code148, Code149, Code150, Code151, Code152, Code153, Code154, Code155, Code156, Code157, Code158, Code159, Code160, Code161,
	Code162, Code163, Code164, Code165, Code166, Code167, Code168, Code169, Code170, Code171, Code172, Code173, Code174, Code175, Code176, Code177, Code178, Code179, Code180, Code181, Code182, Code183, Code184, Code185, Code186, Code187, Code188, Code189, Code190, Code191,
	Code192, Code193, Code194, Code195, Code196, Code197, Code198, Code199, Code200, Code201, Code202, Code203, Code204, Code205, Code206, Code207, Code208, Code209, Code210, Code211, Code212, Code213, Code214, Code215, Code216, Code217, Code218, Code219, Code220, Code221,
	Code222, Code223, Code224, Code225, Code226, Code227, Code228, Code229, Code230, Code231, Code232, Code233, Code234, Code235, Code236, Code237, Code238, Code239, Code240, Code241, Code242, Code243, Code244, Code245, Code246, Code247, Code248, Code249, Code250, Code251,
	Code252, Code253, Code254, Code255, Code256, Code257, Code258, Code259, Code260, Code261, Code262, Code263, Code264, Code265, Code266, Code267, Code268, Code269, Code270, Code271, Code272, Code273, Code274, Code275, Code276, Code277, Code278, Code279, Code280, Code281,
	Code282, Code283, Code284, Code285, Code286, Code287, Code288, Code289, Code290, Code291, Code292, Code293, Code294, Code295, Code296, Code297, Code298, Code299, Code300, Code301, Code302, Code303, Code304, Code305, Code306, Code307, Code308, Code309, Code310, Code311,
	Code312, Code313, Code314, Code315, Code316, Code317, Code318, Code319, Code320, Code321, Code322, Code323, Code324, Code325, Code326, Code327 *float32
	S2 int8 `gorm:"type:tinyint comment '测点2状态'"`
	S3 int8 `gorm:"type:tinyint comment '测点3状态'"`
	S4, S5, S6, S7, S8, S9, S10, S11, S12, S13, S14, S15, S16, S17, S18, S19, S20, S21, S22, S23, S24, S25, S26, S27, S28, S29, S30, S31, S32, S33, S34, S35,
	S36, S37, S38, S39, S40, S41, S42, S43, S44, S45, S46, S47, S48, S49, S50, S51, S52, S53, S54, S55, S56, S57, S58, S59, S60, S61, S62, S63, S64, S65, S66, S67, S68,
	S69, S70, S71, S72, S73, S74, S75, S76, S77, S78, S79, S80, S81, S82, S83, S84, S85, S86, S87, S88, S89, S90, S91, S92, S93, S94, S95, S96, S97, S98, S99, S100, S101,
	S102, S103, S104, S105, S106, S107, S108, S109, S110, S111, S112, S113, S114, S115, S116, S117, S118, S119, S120, S121, S122, S123, S124, S125, S126, S127, S128, S129, S130, S131,
	S132, S133, S134, S135, S136, S137, S138, S139, S140, S141, S142, S143, S144, S145, S146, S147, S148, S149, S150, S151, S152, S153, S154, S155, S156, S157, S158, S159, S160, S161,
	S162, S163, S164, S165, S166, S167, S168, S169, S170, S171, S172, S173, S174, S175, S176, S177, S178, S179, S180, S181, S182, S183, S184, S185, S186, S187, S188, S189, S190, S191,
	S192, S193, S194, S195, S196, S197, S198, S199, S200, S201, S202, S203, S204, S205, S206, S207, S208, S209, S210, S211, S212, S213, S214, S215, S216, S217, S218, S219, S220, S221,
	S222, S223, S224, S225, S226, S227, S228, S229, S230, S231, S232, S233, S234, S235, S236, S237, S238, S239, S240, S241, S242, S243, S244, S245, S246, S247, S248, S249, S250, S251,
	S252, S253, S254, S255, S256, S257, S258, S259, S260, S261, S262, S263, S264, S265, S266, S267, S268, S269, S270, S271, S272, S273, S274, S275, S276, S277, S278, S279, S280, S281,
	S282, S283, S284, S285, S286, S287, S288, S289, S290, S291, S292, S293, S294, S295, S296, S297, S298, S299, S300, S301, S302, S303, S304, S305, S306, S307, S308, S309, S310, S311,
	S312, S313, S314, S315, S316, S317, S318, S319, S320, S321, S322, S323, S324, S325, S326, S327 int8 `gorm:"type:tinyint"`
}

func CreateOpticalFiber90(opticalFiber90 *OpticalFiber90) error {
	return DB.Create(&opticalFiber90).Error
}

func SetCodeField(opticalFiber90 *OpticalFiber90, fieldName string, value string) error {
	// 用反射实现动态设置字段
	field := reflect.ValueOf(opticalFiber90).Elem().FieldByName(fieldName)
	if field.IsValid() && field.CanSet() {
		switch field.Kind() {
		case reflect.Ptr:
			// 指针类型，需要解引用
			if field.Type().Elem().Kind() == reflect.Float32 {
				if value != "" {
					floatValue, err := strconv.ParseFloat(value, 32)
					if err != nil {
						zap.L().Error("Failed to parse float value", zap.Error(err))
						return err
					}
					// 如果指针为 nil，需要先分配内存
					if field.IsNil() {
						field.Set(reflect.New(field.Type().Elem()))
					}
					// 解引用并赋值
					field.Elem().SetFloat(floatValue)
				}
			}
		case reflect.Float32:
			if value != "" {
				floatValue, err := strconv.ParseFloat(value, 32)
				if err != nil {
					zap.L().Error("Failed to parse float value", zap.Error(err))
					return err
				}
				field.SetFloat(floatValue)
			}
		case reflect.Int8:
			if value != "" {
				intValue, err := strconv.Atoi(value)
				if err != nil {
					zap.L().Error("Failed to parse int value", zap.Error(err))
					return err
				}
				field.SetInt(int64(intValue))
			}
		default:
			zap.L().Error("Unsupported field type", zap.String("fieldName", fieldName), zap.String("fieldType", field.Type().String()))
		}
	}
	return nil
}

func GetLastOpticalFiber90() (*OpticalFiber90, error) {
	var opticalFiber90 OpticalFiber90
	err := DB.Order("id desc").Limit(1).Find(&opticalFiber90).Error
	if err != nil {
		return nil, err
	}
	return &opticalFiber90, nil
}

func GetOpticalFiber90ByYearMonthDayHour(year, month, day, hour string) (*OpticalFiber90, error) {
	var opticalFiber90 OpticalFiber90
	err := DB.Where("year = ? AND month = ? AND day = ? AND hour = ?", year, month, day, hour).Order("id desc").Limit(1).Find(&opticalFiber90).Error
	if err != nil {
		return nil, err
	}
	return &opticalFiber90, nil
}

func GetLastOpticalFiber90AvgAndEnvTemperature() (avg, env, windSpeed *float32, err error) {
	type Result struct {
		Avg            *float32
		EnvTemperature *float32
		WindSpeed      *float32
	}
	var result Result
	// code2作环境温度
	err = DB.Table("optical_fiber90").Select(fmt.Sprintf("%s/%s", selectCodeSumSql(), selectSSumSql())+" as avg", "case when env_temperature > code2 then env_temperature else code2 end as env_temperature", "wind_speed").Order("id desc").Limit(1).Scan(&result).Error
	if err != nil {
		return nil, nil, nil, err
	}
	return result.Avg, result.EnvTemperature, result.WindSpeed, nil
}

type YearAvgDataResult struct {
	Year           int
	Temperature    *float32
	EnvTemperature *float32
	WindSpeed      *float32
}

func FindOpticalFiber90YearAvg() ([]*YearAvgDataResult, error) {
	var yearAvgDataResult []*YearAvgDataResult
	err := DB.Table("optical_fiber90").
		// code2作环境温度
		Select("year", fmt.Sprintf("AVG(%s/%s) as temperature", selectCodeSumSql(), selectSSumSql()), "AVG(case when env_temperature > code2 then env_temperature else code2 end) as env_temperature", "AVG(wind_speed) as wind_speed").
		Where("code2 is not null"). //code2是null大概率别的点位都是null，现场断电的状态，排除计算，不然会按0进行计算，虽然整月都没数据的概率很小，但保险起见还是加上
		Group("year").
		Scan(&yearAvgDataResult).Error
	return yearAvgDataResult, err
}

// 不足一年的取2026年3、4月份的数据作为2025年的数据
func FindOpticalFiber90AvgByYearMonths(year int, months []int) (*YearAvgDataResult, error) {
	var yearAvgData YearAvgDataResult
	err := DB.Table("optical_fiber90").
		// code2作环境温度
		Select(fmt.Sprintf("AVG(%s/%s) as temperature", selectCodeSumSql(), selectSSumSql()), "AVG(case when env_temperature > code2 then env_temperature else code2 end) as env_temperature", "AVG(wind_speed) as wind_speed").
		Where("code2 is not null and year = ? and month in ?", year, months). //code2是null大概率别的点位都是null，现场断电的状态，排除计算，不然会按0进行计算，虽然整月都没数据的概率很小，但保险起见还是加上
		Take(&yearAvgData).Error
	if err != nil {
		return nil, err
	}
	yearAvgData.Year = year - 1
	return &yearAvgData, err
}

type MonthAvgDataResult struct {
	Year           int
	Month          int
	Temperature    *float32
	EnvTemperature *float32
	WindSpeed      *float32
}

func FindOpticalFiber90MonthAvg(monthNum int) ([]*MonthAvgDataResult, error) {
	dateSqlArray := make([]string, monthNum)
	// 获取最近monthNum 个月的年月
	for index := range monthNum {
		year, month, _ := time.Now().AddDate(0, -index, 0).Date()
		dateSqlArray[index] = fmt.Sprintf(" (year = %d and month = %d) ", year, month)
	}
	var avgDataResult []*MonthAvgDataResult
	err := DB.Table("optical_fiber90").
		// code2作环境温度
		Select("year", "month", fmt.Sprintf("AVG(%s/%s) as temperature", selectCodeSumSql(), selectSSumSql()), "AVG(case when env_temperature > code2 then env_temperature else code2 end) as env_temperature", "AVG(wind_speed) as wind_speed").
		Where(strings.Join(dateSqlArray, " or ")).
		Where("code2 is not null"). //code2是null大概率别的点位都是null，现场断电的状态，排除计算，不然会按0进行计算，虽然整月都没数据的概率很小，但保险起见还是加上
		Group("year, month").
		Scan(&avgDataResult).Error
	return avgDataResult, err
}

type HourAvgDataResult struct {
	Year           int
	Month          int
	Day            int
	Hour           int
	Temperature    *float32
	EnvTemperature *float32
	WindSpeed      *float32
}

func FindOpticalFiber90HourAvg(monthNum int) ([]*HourAvgDataResult, error) {
	dateSqlArray := make([]string, monthNum)
	// 获取最近monthNum 个月的年月
	for index := range monthNum {
		year, month, _ := time.Now().AddDate(0, -index, 0).Date()
		dateSqlArray[index] = fmt.Sprintf(" (year = %d and month = %d) ", year, month)
	}
	var hourAvgDataResult []*HourAvgDataResult
	err := DB.Table("optical_fiber90").
		// code2作环境温度
		Select("year", "month", "day", "hour", fmt.Sprintf("AVG(%s/%s) as temperature", selectCodeSumSql(), selectSSumSql()), "AVG(case when env_temperature > code2 then env_temperature else code2 end) as env_temperature", "AVG(wind_speed) as wind_speed").
		Where(strings.Join(dateSqlArray, " or ")).
		Where("code2 is not null"). // code2是null大概率别的点位都是null，现场断电的状态，排除计算，不然会按0进行计算
		Group("year, month, day, hour").
		Scan(&hourAvgDataResult).Error
	return hourAvgDataResult, err
}

type ThreePointsData struct {
	Year           int
	Month          int
	Day            int
	Hour           int
	WindSpeed      *float32
	EnvTemperature *float32
	C1             *float32
	C2             *float32
	C3             *float32
}

func GetOpticalFiber90IDByYearMonthDayHour(year, month, day, hour string) (id uint, err error) {
	err = DB.Table("optical_fiber90").Where("year = ? AND month = ? AND day = ? AND hour = ?", year, month, day, hour).Limit(1).Pluck("id", &id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func FindOpticalFiber90ThreePointsDataOpen(c1, c2, c3, s1, s2, s3, year1, month1, day1, hour1 string) ([]*ThreePointsData, error) {
	var threePointsData []*ThreePointsData
	err := DB.Table("optical_fiber90").
		Select("year", "month", "day", "hour", "wind_speed", "case when env_temperature > code2 then env_temperature else code2 end as env_temperature",
			fmt.Sprintf("case when %s=1 then %s else null end as c1", s1, c1),
			fmt.Sprintf("case when %s=1 then %s else null end as c2", s2, c2),
			fmt.Sprintf("case when %s=1 then %s else null end as c3", s3, c3),
		).Where("id >= (?)",
		DB.Table("optical_fiber90").Select("id").Where("year = ? and month = ? and day = ? and hour = ?", year1, month1, day1, hour1),
	).
		Scan(&threePointsData).Error
	return threePointsData, err
}

func FindOpticalFiber90ThreePointsDataClosed(c1, c2, c3, s1, s2, s3, year1, month1, day1, hour1 string, id2 uint) ([]*ThreePointsData, error) {
	var threePointsData []*ThreePointsData
	err := DB.Table("optical_fiber90").
		Select("year", "month", "day", "hour", "wind_speed", "case when env_temperature > code2 then env_temperature else code2 end as env_temperature",
			fmt.Sprintf("case when %s=1 then %s else null end as c1", s1, c1),
			fmt.Sprintf("case when %s=1 then %s else null end as c2", s2, c2),
			fmt.Sprintf("case when %s=1 then %s else null end as c3", s3, c3),
		).Where("id >= (?) and id <= (?)",
		DB.Table("optical_fiber90").Select("id").Where("year = ? and month = ? and day = ? and hour = ?", year1, month1, day1, hour1),
		id2,
	).
		Scan(&threePointsData).Error
	return threePointsData, err
}

func selectCodeSumSql() string {
	var codeArray = make([]string, common.EffEndPoint-common.EffStartPoint+1)
	for i := common.EffStartPoint; i <= common.EffEndPoint; i++ {
		codeArray[i-common.EffStartPoint] = fmt.Sprintf("code%d*s%d", i, i)
	}
	return "(" + strings.Join(codeArray, " + ") + ")"
}

func selectSSumSql() string {
	var codeArray = make([]string, common.EffEndPoint-common.EffStartPoint+1)
	for i := common.EffStartPoint; i <= common.EffEndPoint; i++ {
		codeArray[i-common.EffStartPoint] = fmt.Sprintf("s%d", i)
	}
	return "(" + strings.Join(codeArray, " + ") + ")"
}
