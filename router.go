package gobbs

import (
	"log"
	"net/http"
	"reflect"
	"strings"
)

//控制器信息
type ControllerInfo struct {
	controllerType reflect.Type
	controllerName string
	funcName       string //方法名称
}

//控制器注册
type ControllerRegister struct {
	Router map[string]*ControllerInfo
}

//控制器注册添加  路由器添加
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
	route.controllerName = t.Name()

	p.Router[url] = route

}

//重写http Handle interface
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
		//使用断言方式 调用对应的init方法进行初始化
		execController, ok := vc.Interface().(ControllerInterface)
		if !ok {
			log.Fatal("controller is not ControllerInterface")
		}
		//调用初始化方法
		execController.Init(httpc, v.funcName, v.controllerName)
		//反射调用运行方法
		method := vc.MethodByName(v.funcName)
		method.Call(param)
	} else {
		http.Error(w, "not found page", 404)
	}

}
