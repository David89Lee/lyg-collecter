package common

const (
	InitCodeNum     = 2
	CodeNum         = 327 // CodeNum/8 作为[]byte的索引，CodeNum%8作为byte内部的索引
	PipeLength      = 102
	EffStartPoint   = 5
	EffEndPoint     = 310 // 目前最终有效点，314~327暂时异常，311~313作为尾部裸露部分
	StatusBytesLong = CodeNum/8 + 1
)
