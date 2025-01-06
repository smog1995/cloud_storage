package util

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

// http相应的通用数据
type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 生成resp对象
func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// JSONBytes对象转二进制文件
func (resp *RespMsg) JSONBytes() []byte {
	r, err := json.Marshal(resp)
	if err != nil {
		zap.S().Fatalf("JSON转换失败:%s", err.Error())
	}
	return r
}

func (resp *RespMsg) JSONString() string {
	r, err := json.Marshal(resp)
	if err != nil {
		zap.S().Fatalf("JSON转换字符串失败:%s", err.Error())
	}
	return string(r)
}

// GenSimpleRespStream : 只包含code和message的响应体([]byte)
func GenSimpleRespStream(code int, msg string) []byte {
	return []byte(fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg))
}

// GenSimpleRespString : 只包含code和message的响应体(string)
func GenSimpleRespString(code int, msg string) string {
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg)
}
