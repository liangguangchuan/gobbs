package gobbs

import (
	"fmt"
	"log"
	"net/http"
)

const (
	// 版本
	VERSION = "1.0"

	// DEV is for develop
	DEV = "dev"
	// PROD is for production
	PROD = "prod"
)

var (
	BApp *Controller
)

func init() {
	BApp = NewBApp()
}

func NewBApp() *Controller {
	return &Controller{
		handle: &ControllerRegister{
			Router: make(map[string]*ControllerInfo),
		},
	}
}
func Run() {
	var (
		server_listen string = ""
		err           error
	)
	if BConf.Host != "" && BConf.Port != 0 {
		server_listen = fmt.Sprintf("%s:%d", BConf.Host, BConf.Port)
	}
	log.Println("server listn :", server_listen)
	err = http.ListenAndServe(server_listen, Controller{})
	if err != nil {
		log.Fatal(err.Error())
	}
}
func AddRoute(url, FuncName string, c ControllerInterface) {
	BApp.handle.Add(url, FuncName, c)
}
