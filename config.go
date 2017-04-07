package gobbs

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/liangguangchuan/gobbs/lib"
)

var (
	//基础配置文件
	BConf *Conf
	//项目访问路径
	AppPath string
	//运行模式 dev prod
	RunMode string
	//项目目录
	WorkPath string
	//支持 view 层解析格式
	TplExt = []string{"tpl", "html", "htm"}
)

//配置构造体
type Conf struct {
	Host      string            `xml:"server_host"` //运行域名
	Port      int64             `xml:"server_port"` //运行端口
	AppName   string            `xml:"app_name"`    //项目名称
	RunMode   string            `xml:"run_mode"`    //运行模块
	TplPATH   string            `xml:"tpl_path"`    //模板路径
	TplExt    string            `xml:"tpl_ext"`     //模板后缀
	StaticDir map[string]string `xml:"static_dir"`  //静态文件目录
	Db        confDB            `xml:"db"`          //db 数据
}

//db配置 可能直接使用第三方orm
type confDB struct {
	Host     string //请求地址
	Port     int64  //端口
	Username string //登录用户
	Userpass string //登录密码
	Datebase string //请求数据库
	TablePre string //表前缀
}

//配置初始化
func init() {
	//创建  Conf
	BConf = newConf()
	var err error
	//获取当前运行的 路径 如果获取失败抛出错误
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		log.Fatal(err)
	}

	//获取工作目录
	WorkPath, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//拼接 conf 路径
	confPath := filepath.Join(WorkPath, "conf", "conf.xml")
	//如果项目目录拼接conf/conf.xml 不存在对应文件
	if !lib.FileExists(confPath) {
		confPath = filepath.Join(AppPath, "conf", "conf.xml")
		// 根据运行文件目录拼接conf/conf.xml 不存在对应文件
		if !lib.FileExists(confPath) {
			return
		}
	}

	//读取文件并赋值 conf
	if err = parseConfig(confPath); err != nil {
		log.Fatal(err)
	}

	//tpl 配置 后缀检查
	if TplExtCheck(BConf.TplExt) == false {
		log.Fatal("`tpl_ext` can only be html,htm,tpl")
	}
	//如果当前运行为 develop 输出对应的 conf配置
	if BConf.RunMode == DEV {
		log.Println(BConf)
	}
}

//生产conf
func newConf() *Conf {
	return &Conf{
		Host:      "127.0.0.1",
		Port:      8080,
		AppName:   "xiaochuan",
		RunMode:   DEV,
		TplPATH:   "view",
		TplExt:    "tpl",
		StaticDir: map[string]string{"public": "public"},
		Db:        confDB{},
	}
}

//解析 conf.xml
func parseConfig(confPath string) error {
	//文件读取
	fileData, err := ioutil.ReadFile(confPath)

	if err != nil {
		return err
	}
	//xml 解析到对应构造体
	err = xml.Unmarshal(fileData, BConf)

	return err
}

//模板后缀检查
func TplExtCheck(ext string) bool {

	for _, v := range TplExt {

		if ext == v {
			return true
		}
	}
	return false
}
func (c *Conf) getConf(key string) interface{} {
	val := reflect.ValueOf(c)
	v := val.Elem().FieldByName(key)
	//如果存在对应的字段
	if v.IsValid() {
		return v.Interface()
	} else {
		return nil
	}

}
