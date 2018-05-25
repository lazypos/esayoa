package esayoa

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

//请求管理页
func QueryManagerPage(c *gin.Context) {
	type STManagerTPL struct {
		HeaderPicName string //用户头像
		Account       string //账户
		HomePage      string //主页
	}

	tpl := GMapTPLFile["manager.tpl"]
	info, ok := c.Get("Client")
	if !ok {
		log.Println("无法获取客户端信息")
		c.File("login.html")
		return
	}
	//信息装配
	sInfo := info.(*SessionInfo)
	st := &STManagerTPL{HeaderPicName: "/images/logo.png", Account: sInfo.UserInfo.Uid}
	st.HomePage = GConfig.GetConfig("homepage")
	if len(sInfo.UserInfoEx.HeaderPic) > 0 {
		st.HeaderPicName = sInfo.UserInfoEx.HeaderPic
	}
	if len(sInfo.UserInfoEx.UserName) > 0 {
		st.Account = sInfo.UserInfoEx.UserName
	}

	err := tpl.Execute(c.Writer, st)
	if err != nil {
		log.Println("执行模板失败", err)
	}
}

//frame页分发
func QueryFramePage(c *gin.Context) {
	st, _ := c.Get("Client")
	sInfo := st.(*SessionInfo)

	keyname := c.Param("name")
	var stRst interface{} = nil
	switch keyname {
	case "home":
		stRst = FrameHome(sInfo)
	case "notify":
		stRst = FrameNotifyPage(sInfo, c)
	case "notify.do":
		stRst = FrameNotify(sInfo, c)
	case "employee":
		stRst = FrameEmployee(sInfo)
	}
	if stRst == nil {
		log.Println("处理失败", keyname)
		return
	}

	tpl := GMapTPLFile[fmt.Sprintf(`fm_%s.tpl`, keyname)]
	err := tpl.Execute(c.Writer, stRst)
	if err != nil {
		log.Println("执行模板失败", keyname, err)
	}
}

//POP页分发
func QueryPopPage(c *gin.Context) {
	st, _ := c.Get("Client")
	sInfo := st.(*SessionInfo)

	keyname := c.Param("name")
	stRst := ""
	switch keyname {
	case "newnotify":
		stRst = PopNotify(sInfo)
	case "looknfy":
		stRst = LookNotify(sInfo, c)
	}
	c.String(http.StatusOK, stRst)
}

func Poptest(c *gin.Context) {
	log.Println("Ddsadsadasdsadadasdas======")
	c.File("./template/pop_looknotify.tpl")
}

func ManagerSetCallBack(r *gin.Engine) {
	r.GET("/manager", Intercept(0), QueryManagerPage)
	r.GET("/frame/:name", Intercept(0), QueryFramePage)
	r.GET("/pop/:name", Intercept(0), QueryPopPage)
	r.POST("/notify.do", UploadNotify)
	r.GET("/test", Poptest)

	GMapTPLFuncs["NotifyLists"] = NotifyLists
	GMapTPLFuncs["OrderLists"] = OrderLists
	GMapTPLFuncs["EmployeeList"] = EmployeeList
}
