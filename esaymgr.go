package esayoa

import (
	"github.com/gin-gonic/gin"
	"log"
)

func EsayOAInit(r *gin.Engine) error {
	log.Println("系统初始化中......")
	//连接并初始化数据库
	if err := GDBpgopt.Init(); err != nil {
		log.Println("数据库连接错误", err)
		return err
	}
	if err := InitDB(); err != nil {
		log.Println("数据库初始化错误", err)
		return err
	}
	//初始化会话管理
	if err := GSessions.Init(); err != nil {
		log.Println("会话初始化错误", err)
		return err
	}
	if err := GConfig.InitConfig(); err != nil {
		log.Println("出释化配置错误", err)
		return err
	}

	//注册各种回调
	LoginSetCallBack(r)
	ManagerSetCallBack(r)

	//预加载所有模板
	if err := TPLLoad(); err != nil {
		log.Println("模板加载错误", err)
		return err
	}
	return nil
}
