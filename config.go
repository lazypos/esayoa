package esayoa

import (
	"sync"
)

type STConfig struct {
	MapConfig map[string]string
	MuxCfg    sync.Mutex
}

var GConfig = &STConfig{}

//初始化全局设置
func (this *STConfig) InitConfig() error {
	this.MapConfig = make(map[string]string)
	return nil
}

func (this *STConfig) GetConfig(key string) string {
	this.MuxCfg.Lock()
	defer this.MuxCfg.Unlock()
	return this.MapConfig[key]
}
