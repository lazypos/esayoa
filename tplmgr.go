package esayoa

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var GMapTPLFile = make(map[string]*template.Template) //文件名->内容
var GMapHTMFile = make(map[string]string)             //文件名->内容
var GMapTPLFuncs = make(map[string]interface{})       //函数名->地址

func TPLLoad() error {
	err := filepath.Walk("./template/", func(path string, info os.FileInfo, err error) error {
		//模板
		if filepath.Ext(path) == ".tpl" {
			text, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("读取tpl文件错误", err)
				return err
			}
			tpl, err := template.New(info.Name()).Funcs(GMapTPLFuncs).Parse(string(text[:]))
			if err != nil {
				log.Println("解析模板文件失败", err)
				return err
			}
			log.Println("成功解析：", path)
			GMapTPLFile[info.Name()] = tpl
		}
		//非模板
		if filepath.Ext(path) == ".htm" {
			text, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("读取hml文件错误", err)
				return err
			}
			GMapHTMFile[info.Name()] = string(text[:])
		}
		return nil
	})
	return err
}
