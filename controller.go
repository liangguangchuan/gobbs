package gobbs

import (
	"fmt"
	"log"
	"path/filepath"

	"strings"

	"github.com/liangguangchuan/gobbs/lib"
)

//控制器 构造体
type Controller struct {
	Ctx            *ctx
	controllerName string                      //控制器名称
	actionName     string                      //方法名称
	method         string                      //请求方式
	Data           map[interface{}]interface{} //控制器数据

}

//控制器接口
type ControllerInterface interface {
	Init(c *ctx, controllerName, actionName string)
	Get()
}

func (this *Controller) Get() {}

//程序停止
func (this *Controller) StopRun() {
	log.Fatal("gobbs stop")
}

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

//页面跳转
func (this *Controller) PageJump(url string) {
	this.Ctx.Redirect(strings.TrimSpace(url))
}

//模板赋值
func (this *Controller) Assign(key, value interface{}) {
	this.Data[key] = value
}

//模板显示
func (this *Controller) Display(tplname ...string) {
	//模板路径 模板名称 模板后缀
	var tpl_path, tpl_filename, tpl_ext string
	//读取配置文件 模板后缀
	tpl_ext = fmt.Sprintf(".%s", BConf.TplExt)
	//如果存在 参数传递
	if len(tplname) > 0 {
		//如果存在传递后缀去掉对应后缀
		if strings.Index(tplname[0], tpl_ext) != -1 {
			tplname[0] = strings.TrimRight(tplname[0], tpl_ext)
		}
		//生成对应目录名称
		tpl_filename = fmt.Sprintf("%s%s", tplname[0], tpl_ext)
		//如果存在 / 说明要跨目录调用对应view
		if strings.Index(tplname[0], "/") == -1 {
			tpl_path = filepath.Join(BConf.TplPATH, this.controllerName, tpl_filename)
		} else {
			tpl_path = filepath.Join(BConf.TplPATH, tpl_filename)
		}

	} else {
		tpl_filename = fmt.Sprintf("%s%s", this.actionName, tpl_ext)
		tpl_path = filepath.Join(BConf.TplPATH, this.controllerName, tpl_filename)
	}

	if lib.FileExists(tpl_path) != true {
		log.Fatal(fmt.Sprintf("`%s` is no exist", tpl_path))
	}
}
