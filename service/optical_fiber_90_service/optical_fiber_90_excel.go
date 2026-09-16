package optical_fiber_90_service

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// threePointsSheetName 导出 Excel 的 sheet 名称
const threePointsSheetName = "截面三点历史数据"

// threePointsHeaders 为导出 Excel 的表头，顺序与 ThreePointsDataList 各数组字段一一对应
var threePointsHeaders = []string{"时间", "风速(m/s)", "环境温度(°C)", "异常温度(°C)", "光纤-0°(°C)", "光纤-45°(°C)", "光纤-225°(°C)"}

// ExportOpticalFiber90ThreePointsDataExcel 将界面上三个点的历史曲线数据导出为 Excel（xlsx）。
// 参数与 FindOpticalFiber90ThreePointsData 保持一致。
// Excel 中每一行的数据按各自数组的索引对齐，C1/C2/C3 大于本行异常温度时以红色字体显示。
// 返回 xlsx 文件二进制内容、建议文件名与错误。
func (t *OpticalFiber90Service) ExportOpticalFiber90ThreePointsDataExcel(i int, year1, month1, day1, hour1, year2, month2, day2, hour2 string) ([]byte, string, error) {
	list, err := t.FindOpticalFiber90ThreePointsData(i, year1, month1, day1, hour1, year2, month2, day2, hour2)
	if err != nil {
		return nil, "", err
	}
	if list == nil {
		return nil, "", fmt.Errorf("截面三点历史数据为空")
	}

	data, err := buildThreePointsExcel(list, i)
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("截面三点历史数据_位置%d_%s.%s.%s_%s时-%s.%s.%s_%s时.xlsx",
		i, year1, month1, day1, hour1, year2, month2, day2, hour2)
	return data, filename, nil
}

// buildThreePointsExcel 根据三点历史数据生成 xlsx 文件二进制内容。
// 第一行合并整行写入位置与管墩号，第二行为表头，数据从第三行开始。
func buildThreePointsExcel(list *ThreePointsDataList, i int) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := threePointsSheetName
	sheetIndex, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	// 删除 excelize 默认创建的 Sheet1
	if err = f.DeleteSheet("Sheet1"); err != nil {
		return nil, err
	}
	f.SetActiveSheet(sheetIndex)

	// 表头样式：加粗、居中、浅蓝底
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#DCE6F1"}},
	})
	if err != nil {
		return nil, err
	}
	// 数据单元格：保留一位小数
	oneDecimal := "0.0"
	normalStyle, err := f.NewStyle(&excelize.Style{CustomNumFmt: &oneDecimal})
	if err != nil {
		return nil, err
	}
	// 超过本行异常温度：红色加粗
	redStyle, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &oneDecimal,
		Font:         &excelize.Font{Color: "FF0000", Bold: true},
	})
	if err != nil {
		return nil, err
	}

	// 位置i/102和管墩号写入第一行，管墩号用FindPipeSection(i)查询
	titleStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}
	// 合并第一行的 len(threePointsHeaders) 个单元格
	titleEnd, err := excelize.CoordinatesToCellName(len(threePointsHeaders), 1)
	if err != nil {
		return nil, err
	}
	if err = f.MergeCell(sheet, "A1", titleEnd); err != nil {
		return nil, err
	}
	if err = f.SetCellValue(sheet, "A1", fmt.Sprintf("位置:%d/102  %s", i, FindPipeSection(float32(i)))); err != nil {
		return nil, err
	}
	if err = f.SetCellStyle(sheet, "A1", titleEnd, titleStyle); err != nil {
		return nil, err
	}
	if err = f.SetRowHeight(sheet, 1, 22); err != nil {
		return nil, err
	}

	// 写表头（第二行）
	for col, name := range threePointsHeaders {
		axis, err := excelize.CoordinatesToCellName(col+1, 2)
		if err != nil {
			return nil, err
		}
		if err = f.SetCellValue(sheet, axis, name); err != nil {
			return nil, err
		}
		if err = f.SetCellStyle(sheet, axis, axis, headerStyle); err != nil {
			return nil, err
		}
	}
	// 列宽：A 列放时间文本，稍宽一些
	if err = f.SetColWidth(sheet, "A", "A", 18); err != nil {
		return nil, err
	}
	if err = f.SetColWidth(sheet, "B", "G", 16); err != nil {
		return nil, err
	}
	// 冻结标题行与表头行，数据量较大时方便查看
	if err = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		XSplit:      0,
		YSplit:      2,
		TopLeftCell: "A3",
		ActivePane:  "bottomLeft",
	}); err != nil {
		return nil, err
	}

	// 写数据行：每一行按各自数组的索引对齐
	for row := 0; row < len(list.Times); row++ {
		r := row + 3 // 第 1 行为标题、第 2 行为表头
		abnormal := threePointsValueAt(list.AbnormalTemperature, row)
		c1 := threePointsValueAt(list.C1, row)
		c2 := threePointsValueAt(list.C2, row)
		c3 := threePointsValueAt(list.C3, row)

		if err = f.SetCellValue(sheet, fmt.Sprintf("A%d", r), list.Times[row]); err != nil {
			return nil, err
		}
		if err = setThreePointsFloatCell(f, sheet, fmt.Sprintf("B%d", r), threePointsValueAt(list.WindSpeed, row), false, normalStyle, redStyle); err != nil {
			return nil, err
		}
		if err = setThreePointsFloatCell(f, sheet, fmt.Sprintf("C%d", r), threePointsValueAt(list.EnvTemperature, row), false, normalStyle, redStyle); err != nil {
			return nil, err
		}
		if err = setThreePointsFloatCell(f, sheet, fmt.Sprintf("D%d", r), abnormal, false, normalStyle, redStyle); err != nil {
			return nil, err
		}
		if err = setThreePointsFloatCell(f, sheet, fmt.Sprintf("E%d", r), c1, isOverAbnormalTemp(c1, abnormal), normalStyle, redStyle); err != nil {
			return nil, err
		}
		if err = setThreePointsFloatCell(f, sheet, fmt.Sprintf("F%d", r), c2, isOverAbnormalTemp(c2, abnormal), normalStyle, redStyle); err != nil {
			return nil, err
		}
		if err = setThreePointsFloatCell(f, sheet, fmt.Sprintf("G%d", r), c3, isOverAbnormalTemp(c3, abnormal), normalStyle, redStyle); err != nil {
			return nil, err
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// threePointsValueAt 按下标安全取值，越界返回 nil
func threePointsValueAt(arr []*float32, idx int) *float32 {
	if idx < 0 || idx >= len(arr) {
		return nil
	}
	return arr[idx]
}

// isOverAbnormalTemp 判断温度值是否超过同行异常温度，任一为 nil 时视作不超标
func isOverAbnormalTemp(v, abnormal *float32) bool {
	return v != nil && abnormal != nil && *v > *abnormal
}

// setThreePointsFloatCell 写入浮点单元格，v 为 nil 时留空，red 为 true 时使用红色字体样式
func setThreePointsFloatCell(f *excelize.File, sheet, axis string, v *float32, red bool, normalStyle, redStyle int) error {
	if v == nil {
		return f.SetCellValue(sheet, axis, "")
	}
	if err := f.SetCellValue(sheet, axis, *v); err != nil {
		return err
	}
	style := normalStyle
	if red {
		style = redStyle
	}
	return f.SetCellStyle(sheet, axis, axis, style)
}
