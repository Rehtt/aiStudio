package model

type Code int

const (
	OK       = Code(2000)
	ResError = Code(iota + 1001)
	ServerBad
	RequestBad
)

var CodeMap = map[Code]any{
	ResError: &Response{
		Code:  ResError,
		Error: "次数用尽或时间到期",
	},
	ServerBad: &Response{
		Code:  ServerBad,
		Error: "服务器错误，请联系管理员",
	},
	RequestBad: &Response{
		Code:  RequestBad,
		Error: "请求错误",
	},
}
