package model

type Code int

const (
	ResError = Code(iota + 1001)
	ServerBad
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
}
