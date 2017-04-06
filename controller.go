package gobbs

import (
	"strings"
)

//控制器 构造体
type Controller struct {
	Ctx            *ctx
	controllerName string                      //控制器名称
	actionName     string                      //方法名称
	method         string                      //请求方式
	Data           map[interface{}]interface{} //控制器数据
	TplName        string                      //视图目录
}

//控制器接口
type ControllerInterface interface {
	Init(c *ctx, controllerName, actionName string)
	Get()
}

func (this *Controller) Get() {}

//初始化 控制器
func (this *Controller) Init(c *ctx, controllerName, actionName string) {
	this.Ctx = c
	this.controllerName = strings.ToLower(controllerName)
	this.actionName = strings.ToLower(actionName)
	this.method = c.request.Method
	this.Data = make(map[interface{}]interface{})
}

//向页面写入内容字符串
func (this *Controller) WriterString(msg string) {
	this.Ctx.Echo(msg)
}

//json 页面输出
func (this *Controller) ServeJSON() {
	var hasIndent = true
	//如果运行模式为生产环境不缩进输出
	if BConf.RunMode == PROD {
		hasIndent = false
	}

	this.Ctx.JSON(this.Data["json"], hasIndent)
}
