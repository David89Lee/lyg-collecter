package gis_service

import (
	"fmt"
	"image"
	"image/color"
	"lyg-collecter/model"
	"math"
	"strings"

	"github.com/fogleman/gg"
)

const (
	tileSize = 256
)

var (
	//bgColor       = color.RGBA{20, 20, 25, 255}
	//roadColor     = color.RGBA{80, 80, 90, 255}
	//waterColor    = color.RGBA{30, 40, 50, 255}
	//buildingColor = color.RGBA{45, 45, 55, 255}

	// 底图
	bgColor = color.RGBA{26, 51, 102, 255} // #1a3366

	// 道路
	roadColor = color.RGBA{109, 143, 192, 255} // #6D8FC0

	// 水域
	waterColor = color.RGBA{41, 95, 165, 255} // #295FA5

	// 建筑/地块
	buildingColor = color.RGBA{163, 194, 239, 255} // #A3C2EF

	// 湿地
	wetlandColor = color.RGBA{114, 211, 164, 255} // #72D3A4
)

type GisService struct {
}

func (t *GisService) RenderTile(z, x, y int) (image.Image, error) {
	dc := gg.NewContext(tileSize, tileSize)
	dc.SetColor(bgColor)
	dc.Clear()

	// 瓦片范围
	minLng, maxLat := tile2lonlat(x, y, z)
	maxLng, minLat := tile2lonlat(x+1, y+1, z)

	wktList1, err := model.FindWater(minLng, minLat, maxLng, maxLat)
	if err != nil {
		return nil, err
	}
	renderPolygons(dc, wktList1, "water", minLng, minLat, maxLng, maxLat, z, waterColor)

	wktList2, err := model.FindBuilding(minLng, minLat, maxLng, maxLat)
	if err != nil {
		return nil, err
	}
	renderPolygons(dc, wktList2, "building", minLng, minLat, maxLng, maxLat, z, buildingColor)

	wktList3, err := model.FindHighway(minLng, minLat, maxLng, maxLat)
	if err != nil {
		return nil, err
	}
	renderLines(dc, wktList3, minLng, minLat, maxLng, maxLat, z, roadColor)

	wktList4, err := model.FindWetland(minLng, minLat, maxLng, maxLat)
	if err != nil {
		return nil, err
	}
	renderPolygons(dc, wktList4, "wetland", minLng, minLat, maxLng, maxLat, z, wetlandColor)

	return dc.Image(), nil
}

// ------------------------------
// 渲染线
// ------------------------------
func renderLines(dc *gg.Context, wktList []string, minLng, minLat, maxLng, maxLat float64, z int, col color.Color) {

	dc.SetColor(col)
	dc.SetLineWidth(1)

	for _, w := range wktList {
		lines := parseWKTLine(w)
		for _, line := range lines {
			drawLine(dc, line, z, minLng, maxLat, maxLng, minLat)
		}
	}
}

// ------------------------------
// 渲染面
// ------------------------------
func renderPolygons(dc *gg.Context, wktList []string, query string, minLng, minLat, maxLng, maxLat float64, z int, col color.Color) {

	dc.SetColor(col)

	for _, w := range wktList {
		polys := parseWKTPolygon(w)
		for _, poly := range polys {
			drawPolygon(dc, poly, z, minLng, maxLat, maxLng, minLat)
		}
	}
}

// ------------------------------
// 纯原生坐标计算（无任何第三方依赖）
// ------------------------------
func tile2lonlat(x, y, z int) (float64, float64) {
	//n := math.Pi - 2*math.Pi*float64(y)/math.Pow(2, float64(z))
	//lat := 180 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	//lon := float64(x)/math.Pow(2, float64(z))*360 - 180
	//return lon, lat
	n := math.Pow(2, float64(z))
	lon := (float64(x)/n)*360.0 - 180.0
	//lat := (180.0 / math.Pi) * math.Asin(math.Tanh(math.Pi*(1-2*float64(y)/n)))
	lat := math.Atan(math.Sinh(math.Pi*(1-2*float64(y)/n))) * (180.0 / math.Pi)
	return lon, lat
}

// 经纬度转瓦片像素 (标准Web墨卡托，中国完美显示)
func latLonToPixel(lat, lng float64, z int, minLng, maxLat, maxLng, minLat float64) (float64, float64) {
	tileSize := 256.0
	// 计算当前点在瓦片范围内的比例
	x := (lng - minLng) / (maxLng - minLng) * tileSize
	y := (maxLat - lat) / (maxLat - minLat) * tileSize
	return x, y
}

// ------------------------------
// 简易 WKT 解析（无依赖）
// ------------------------------
func parseWKTLine(wkt string) [][]struct{ Lat, Lng float64 } {
	var res [][]struct{ Lat, Lng float64 }
	wkt = strings.TrimPrefix(wkt, "LINESTRING(")
	wkt = strings.TrimSuffix(wkt, ")")
	parts := strings.Split(wkt, ",")
	var line []struct{ Lat, Lng float64 }
	for _, p := range parts {
		var lat, lng float64
		fmt.Sscanf(strings.TrimSpace(p), "%f %f", &lng, &lat)
		line = append(line, struct{ Lat, Lng float64 }{lat, lng})
	}
	res = append(res, line)
	return res
}

func parseWKTPolygon(wkt string) [][]struct{ Lat, Lng float64 } {
	var res [][]struct{ Lat, Lng float64 }
	wkt = strings.TrimPrefix(wkt, "POLYGON((")
	wkt = strings.TrimSuffix(wkt, "))")
	parts := strings.Split(wkt, ",")
	var poly []struct{ Lat, Lng float64 }
	for _, p := range parts {
		var lat, lng float64
		fmt.Sscanf(strings.TrimSpace(p), "%f %f", &lng, &lat)
		poly = append(poly, struct{ Lat, Lng float64 }{lat, lng})
	}
	res = append(res, poly)
	return res
}

func drawLine(dc *gg.Context, line []struct{ Lat, Lng float64 }, z int, minLng, maxLat, maxLng, minLat float64) {
	if len(line) < 2 {
		return
	}

	// 移动到第一个点
	x, y := latLonToPixel(line[0].Lat, line[0].Lng, z, minLng, maxLat, maxLng, minLat)
	dc.MoveTo(x, y)

	// 依次连线
	for _, p := range line[1:] {
		x, y := latLonToPixel(p.Lat, p.Lng, z, minLng, maxLat, maxLng, minLat)
		dc.LineTo(x, y)
	}

	dc.Stroke()
}

func drawPolygon(dc *gg.Context, poly []struct{ Lat, Lng float64 }, z int, minLng, maxLat, maxLng, minLat float64) {
	if len(poly) < 3 {
		return
	}

	x, y := latLonToPixel(poly[0].Lat, poly[0].Lng, z, minLng, maxLat, maxLng, minLat)
	dc.MoveTo(x, y)

	for _, p := range poly[1:] {
		x, y := latLonToPixel(p.Lat, p.Lng, z, minLng, maxLat, maxLng, minLat)
		dc.LineTo(x, y)
	}

	dc.Fill()
}
