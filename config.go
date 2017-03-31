package gobbs

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/liangguangchuan/gobbs/lib"
)

var (
	//基础配置文件
	BConf *Conf
	//项目访问路径
	AppPath string
	//运行模式 dev prod
	RunMode string
)

type Conf struct {
	Port    int64  `xml:"server_port"`
	AppName string `xml:"app_name"`
	RunMode string `xml:"run_mode"`
}

func init() {
	//创建  Conf
	BConf = newConf()
	var err error
	//获取当前运行的 路径 如果获取失败抛出错误
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	//获取工作目录
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//拼接 conf 路径
	confPath := filepath.Join(workPath, "conf", "conf.xml")
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
		panic(nil)
	}
	//输出 最终构造体值
	log.Fatal(BConf)
}

func newConf() *Conf {
	return &Conf{
		Port:    8080,
		AppName: "xiaochuan",
		RunMode: DEV,
	}
}

//解析 conf.xml
func parseConfig(confPath string) error {

	fileData, err := ioutil.ReadFile(confPath)

	if err != nil {
		return err
	}
	err = xml.Unmarshal(fileData, BConf)
	return err
}
