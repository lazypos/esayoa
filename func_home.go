package esayoa

import (
	"bytes"
	"log"
)

type STItem struct {
	Name string
	OK   bool
	Nid  string
}

type STArticlelList struct {
	Date     string
	ArrLines []*STItem //标题->已读状态
}

func NotifyLists(id int, uid string) string {
	arrList, err := LoadArticleList(id, uid)
	if err != nil {
		return ""
	}

	buf := bytes.NewBufferString("")
	for _, arrTitles := range arrList {
		tpl := GMapTPLFile["sub_home_notify.tpl"]
		//log.Println(arrTitles.ArrLines)
		err = tpl.Execute(buf, arrTitles)
		if err != nil {
			log.Println("执行模板sub_home_notify失败", err)
			return ""
		}
	}
	return buf.String()
}

func OrderLists(uid string) string {
	arrList, err := LoadOrderData(uid)
	if err != nil {
		return ""
	}

	buf := bytes.NewBufferString("")
	for _, arrTitles := range arrList {
		tpl := GMapTPLFile["sub_home_order.tpl"]
		err = tpl.Execute(buf, arrTitles)
		if err != nil {
			log.Println("执行模板sub_home_notify失败", err)
			return ""
		}
	}
	return buf.String()
}

func FrameHome(sInfo *SessionInfo) interface{} {
	type STNotify struct {
		ID  int
		UID string
	}
	snfy := &STNotify{ID: sInfo.UserInfo.Id, UID: sInfo.UID}
	return snfy
}
