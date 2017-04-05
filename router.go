package gobbs

import (
	"log"
	"net/http"
	"reflect"
	"strings"
)

type ControllerInfo struct {
	controllerType reflect.Type

	funcName string //方法名称
}
type ControllerRegister struct {
	Router map[string]*ControllerInfo
}

func (p *ControllerRegister) Add(url, FuncName string, c ControllerInterface) {
	reflectVal := reflect.ValueOf(c)
	t := reflect.Indirect(reflectVal).Type()
	//去掉左边 ／
	if strings.Index(url, "/") != -1 {
		url = strings.TrimLeft(url, "/")
	}
	//检测是否存在对应的方法
	if reflectVal.MethodByName(FuncName).IsValid() == false {
		log.Fatal("'" + FuncName + "' method doesn't exist in the controller " + t.Name())
	}
	//初始化
	route := &ControllerInfo{}
	route.controllerType = t
	route.funcName = FuncName

	p.Router[url] = route

}

func (this Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var url = req.URL.Path

	httpc := &ctx{writer: w, request: req}

	log.Println(req.URL.String())

	//去掉左边 /
	if strings.Index(url, "/") != -1 {
		url = strings.TrimLeft(url, "/")
	}
	//查看是否配置对应路由 如果没有 页面404
	if v, ok := BApp.handle.Router[url]; ok != false {
		var param []reflect.Value
		vc := reflect.New(v.controllerType)
		//设置当前操作构造体的值
		vc.Elem().FieldByName("Ctx").Set(reflect.ValueOf(httpc))
		//反射调用方法
		method := vc.MethodByName(v.funcName)
		method.Call(param)
	} else {
		http.Error(w, "not found page", 404)
	}

}
