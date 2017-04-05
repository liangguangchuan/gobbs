package gobbs

import (
	"net/http"
)

type Controller struct {
	handle *ControllerRegister
	Ctx    *ctx
}
type ctx struct {
	writer  http.ResponseWriter
	request *http.Request
}

type ControllerInterface interface {
	Get()
}

func (this Controller) Get() {}

func (this Controller) WriterString(msg string) {
	this.Ctx.writer.Write([]byte(msg))
}
