package gobbs

import (
	"html/template"
	"log"
	"net/url"
	"path/filepath"
	"strings"
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
	tpl_ext = "." + BConf.TplExt
	//如果存在 参数传递
	if len(tplname) > 0 {
		//如果存在传递后缀去掉对应后缀
		if strings.Index(tplname[0], tpl_ext) != -1 {
			tplname[0] = strings.TrimRight(tplname[0], tpl_ext)
		}
		//生成对应目录名称
		tpl_filename = tplname[0] + tpl_ext
		//如果存在 / 说明要跨目录调用对应view
		if strings.Index(tplname[0], "/") == -1 {
			tpl_path = filepath.Join(BConf.TplPATH, this.controllerName, tpl_filename)
		} else {
			tpl_path = filepath.Join(BConf.TplPATH, tpl_filename)
		}

	} else {
		tpl_filename = this.actionName + tpl_ext
		tpl_path = filepath.Join(BConf.TplPATH, this.controllerName, tpl_filename)
	}
	//拼接工作目录生成最终 view路径
	tpl_path = filepath.Join(WorkPath, tpl_path)
	//执行模板文件
	this.executeTemplatFile(tpl_filename, tpl_path)

}

//执行模板文件
func (this *Controller) executeTemplatFile(tpl_name, tpl_path string) {
	//new里面的参数 不能随便传递
	//如果要操作文件必须与文件名同名 如果需要解析多个文件 new 值为第一个
	tpl, err := template.New(tpl_name).ParseFiles(tpl_path)

	if err != nil {
		this.Ctx.RunError(err)
		return
	}
	//模板渲染输出
	err = tpl.Execute(this.Ctx.writer, this.Data)
	if err != nil {
		this.Ctx.RunError(err)
	}
}

//获取get 输入字符串
func (this *Controller) GetInputString(key string) []string {
	var default_value []string
	//获取输入
	get_value := this.GetInput()
	//如果获取输入到的值为 nil 说明并没有传递 get参数
	if get_value == nil {
		return default_value
		//如果获取到对应的key 返回对应key的值
	} else if v, ok := get_value[key]; ok == true {
		return v
	}

	return default_value
}

//获取输入
func (this *Controller) GetInput() url.Values {
	//如果 资源输入为空 说明没有解析
	if this.Ctx.request.Form == nil {
		//解析URL中的查询字符串，并将解析结果更新到r.Form字段
		this.Ctx.request.ParseForm()
	}
	return this.Ctx.request.Form
}
