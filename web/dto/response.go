package dto

import "buding-job/common/constant"

type Response[T any] struct {
	Code int32  `json:"code" form:"code" json:"code" uri:"code" xml:"code" yaml:"code" `
	Msg  string `json:"msg" form:"msg" json:"msg" uri:"msg" xml:"msg" yaml:"msg" `
	Data T      `json:"data" form:"data" json:"data" uri:"data" xml:"data" yaml:"data" `
}

func NewResponse[T any](code int32, msg string, t T) *Response[T] {
	resp := &Response[T]{
		Code: code,
		Msg:  msg,
		Data: t,
	}
	return resp
}

func NewOkResponse[T any](t T) *Response[T] {
	resp := &Response[T]{
		Code: constant.HttpOk,
		Msg:  "ok",
		Data: t,
	}
	return resp
}
func NewErrResponse[T any](msg string, t T) *Response[T] {
	resp := &Response[T]{
		Code: constant.HttpErr,
		Msg:  msg,
	}
	return resp
}
