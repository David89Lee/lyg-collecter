package response

import "github.com/gin-gonic/gin"

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type PageData struct {
	List       interface{} `json:"list"`
	Pagination `json:"pagination"`
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  GetMsg(errCode),
		Data: data,
	})
}

func (g *Gin) PageDataResponse(httpCode, errCode int, total int, pageSize int, currentPage int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  GetMsg(errCode),
		Data: PageData{List: data, Pagination: Pagination{Page: currentPage, PageSize: pageSize, Total: total}},
	})
}
