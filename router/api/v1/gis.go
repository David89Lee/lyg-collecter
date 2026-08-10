package v1

import (
	"bytes"
	"image/png"
	"lyg-collecter/router/response"
	"lyg-collecter/service/gis_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary 获取底图
// @Tags gis地图
// @Produce  json
// @Success 200 {object} response.Response{data=optical_fiber_90_service.InsulationMonitorData} "成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response
// @Router /api/v1/tile/{z}/{x}/{y} [get]
func FindTile(c *gin.Context) {
	resG := response.Gin{C: c}
	service := gis_service.GisService{}

	z, _ := strconv.Atoi(c.Param("z"))
	x, _ := strconv.Atoi(c.Param("x"))
	y, _ := strconv.Atoi(c.Param("y"))

	img, err := service.RenderTile(z, x, y)
	if err != nil {
		zap.L().Error("获取Gis底图切片失败", zap.Error(err))
		resG.Response(http.StatusOK, response.ErrorFindInsulationMonitorData, nil)
		return
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	c.Data(http.StatusOK, "image/png", buf.Bytes())
}
