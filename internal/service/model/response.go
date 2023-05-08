package model

type Response struct {
	Code  Code   `json:"code"`
	Msg   string `json:"msg,omitempty"`
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}
