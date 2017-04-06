package gobbs

import (
	"net/http"
)

type ctx struct {
	writer  http.ResponseWriter //http 回应写入
	request *http.Request       //http 请求
}

func (this *ctx) Header(key, val string) {
	this.writer.Header().Set(key, val)
}

func (this *ctx) Echo(result string) {
	this.writer.Write([]byte(result))
}
