package esayoa

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//默认页面
func FrameNotifyPage(sInfo *SessionInfo, exInfo interface{}) interface{} {
	type STNotifyPop struct {
		Display string
		Counts  string
	}

	ctx := exInfo.(*gin.Context)
	del := ctx.Query("del")
	if len(del) > 0 {
		if _, err := GDBpgopt.Execute(fmt.Sprintf(SQL_DEL_NOTIFY, del)); err != nil {
			log.Println("删除文章失败", err, del)
			return ""
		}
	}

	counts, err := GDBpgopt.QueryVal(SQL_COUNT_ARTICL)
	if err != nil {
		log.Println("查询通知总数失败.", err)
		return nil
	}
	dis := "none"
	if sInfo.UserInfo.Auth > 10 {
		dis = "inline"
	}
	nst := &STNotifyPop{Display: dis, Counts: counts}
	return nst
}

//分页内容
func FrameNotify(sInfo *SessionInfo, exInfo interface{}) interface{} {
	type STNotifys struct {
		Del   bool
		Title string
		Date  string
		Nid   string
	}

	//再展示
	ctx := exInfo.(*gin.Context)
	pid, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || pid < 1 {
		pid = 1
	}
	offset := (pid - 1) * 20
	rows, err := GDBpgopt.Query(fmt.Sprintf(SQL_LOAD_ALLARTICL, offset))
	if err != nil {
		log.Println("加载通知列表错误", err)
		return ""
	}
	defer rows.Close()

	d := false
	if sInfo.UserInfo.Auth > 90 {
		d = true
	}

	st := []*STNotifys{}
	for rows.Next() {
		var id string
		var title string
		var date string
		if err = rows.Scan(&id, &title, &date); err != nil {
			log.Println("查询通知总数失败.", err)
			return ""
		}
		nt := &STNotifys{Del: d, Title: title, Date: date[:19], Nid: id}
		st = append(st, nt)
	}

	return st
}

func PopNotify(sInfo *SessionInfo) string {
	text, err := ioutil.ReadFile("./template/pop_newnotify.tpl")
	if err != nil {
		log.Println("读取文件失败", err)
		return ""
	}
	return string(text[:])
}

//发布新通知
func UploadNotify(c *gin.Context) {
	title, _ := c.GetPostForm("title")
	content, _ := c.GetPostForm("content")
	if len(title) == 0 || len(title) > 200 || len(content) == 0 {
		c.String(http.StatusOK, "通知格式不正确！")
		return
	}
	title = FormatInputText(title)
	content = FormatInputText(content)

	fs, err := c.FormFile("fujian")
	filename := ""
	fsize := 0
	if err != nil {
		log.Println("没有附件.")
	} else {
		filename = fmt.Sprintf(`.\\user\\upload\\%v_%v`, time.Now().Unix(), fs.Filename)
		if err = c.SaveUploadedFile(fs, filename); err != nil {
			log.Println("保存附件失败.", err)
			c.String(http.StatusOK, "系统异常！")
			return
		}
	}
	//保存到数据库
	val, err := GDBpgopt.QueryVal(fmt.Sprintf(SQL_SAVE_ARTICLE, title, content, "", filename, fsize/(1024*1024)))
	if err != nil {
		log.Println("保存文件错误", err, val)
		return
	}
	c.String(http.StatusOK, "")
}

//查看文章
func LookNotify(sInfo *SessionInfo, exInfo interface{}) string {
	type STNFY struct {
		Title   string
		Date    string
		Content string
		Fujian  string
		Size    string
		BFujian bool
	}

	ctx := exInfo.(*gin.Context)
	id := ctx.Query("id")
	nid, err := strconv.Atoi(id)
	if err != nil || nid < 0 {
		log.Println("打开通知有误", id)
		return ""
	}
	rows, err := GDBpgopt.Query(fmt.Sprintf(SQL_LOAD_NOTIFY, id))
	if err != nil {
		log.Println("查询文章错误", id)
		return ""
	}
	defer rows.Close()

	for rows.Next() {
		var title string
		var content string
		var date string
		var fujian string
		var fsize string

		if err = rows.Scan(&title, &content, &date, &fujian, &fsize); err != nil {
			log.Println("获取查询结果字段失败", err)
			return ""
		}
		fjname := ""
		bfujian := false
		if len(fujian) > 0 {
			bfujian = true
			fjname = filepath.Base(strings.Replace(fujian, "\\\\", "\\", -1))
			fjname = fjname[11:]
		}

		log.Println("=-/-/--/-/-/-/---", fujian, fjname, bfujian)
		st := &STNFY{Title: title, Date: date[:19], Content: content, Fujian: fjname, Size: fsize, BFujian: bfujian}
		tpl := GMapTPLFile["pop_looknotify.tpl"]
		buf := bytes.NewBufferString("")
		if err = tpl.Execute(buf, st); err != nil {
			log.Println("组模板失败:", err)
			return ""
		}
		return buf.String()
	}
	return ""
}
