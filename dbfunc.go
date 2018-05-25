package esayoa

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	SQL_LOAD_BASEINFO  = `SELECT * FROM userbase WHERE UID='%s'`
	SQL_LOAD_ARTICLE   = `SELECT notifyid, title, updatetime FROM notifys WHERE sendto like '%%,%v,%%' or sendto='' order by updatetime desc LIMIT 10`
	SQL_LOAD_ARTIREAD  = `SELECT nfyreadids FROM notifyread WHERE uid='%s'`
	SQL_LOAD_LEAVEINFO = `SELECT createtime FROM leave WHERE uid='%s' and state<>0`
	SQL_LOAD_ALLARTICL = `SELECT notifyid, title, updatetime FROM notifys order by notifyid desc LIMIT 20 OFFSET %d`
	SQL_COUNT_ARTICL   = `SELECT count(*) FROM notifys`
	SQL_COUNT_USERINFO = `SELECT count(*) FROM userbase WHERE uid<>'adminesayoa'`
	SQL_LOAD_USERINFO  = `SELECT * FROM userinfo order by id limit 15`
	//标题 内容 发送目标 文件路径 文件大小
	SQL_SAVE_ARTICLE  = `SELECT CreateNotify('%v','%v','%v','%v','%vMB')`
	SQL_SEARCH_NOTIFY = `SELECT notifyid, title, updatetime FROM notifys WHERE title like '%%%v%%' or content like  '%%%v%%' order by notifyid desc`
	SQL_DEL_NOTIFY    = `Delete from notifys WHERE notifyid=%v`
	//查询通知
	SQL_LOAD_NOTIFY = `SELECT title, content, updatetime,fujian,fjszie FROM notifys WHERE notifyid=%v`
)

//初始化数据
func InitDB() error {
	//先判断是否需要初始化 1初始化文件 2程序调用参数
	if len(os.Args) == 1 || os.Args[1] != "init" {
		return nil
	}
	if _, err := os.Stat("esayoa.dat"); err == nil {
		return nil
	}
	for _, sql := range GSQLArray {
		if _, err := GDBpgopt.Execute(sql); err != nil {
			log.Println("初始化数据库失败:", sql)
			return err
		}
	}
	//存文件
	if err := ioutil.WriteFile("esayoa.dat", []byte("ok"), 0x666); err != nil {
		log.Println("初始化数据库失败->文件", err)
		return err
	}
	log.Println("初始化数据库成功.")
	return nil
}

//加载用户基本信息
func LoadUserBaseInfo(uid string) (*STUserBaseInfo, error) {
	row, err := GDBpgopt.Query(fmt.Sprintf(SQL_LOAD_BASEINFO, uid))
	if err != nil {
		return nil, err
	}
	defer row.Close()
	//无记录
	if !row.Next() {
		return nil, nil
	}
	st := &STUserBaseInfo{}
	err = row.Scan(&st.Id, &st.Uid, &st.Pass, &st.Auth, &st.RegistTime)
	if err != nil {
		return nil, err
	}
	return st, nil
}

//加载文章信息
func LoadArticleList(id int, uid string) ([]*STArticlelList, error) {
	arrArticles := []*STArticlelList{}

	readids, err := GDBpgopt.QueryVal(fmt.Sprintf(SQL_LOAD_ARTIREAD, uid))
	if err != nil {
		log.Println("查询已读文章ID错误", uid, err)
		return arrArticles, err
	}
	//查询文章
	row, err := GDBpgopt.Query(fmt.Sprintf(SQL_LOAD_ARTICLE, id))
	if err != nil {
		return arrArticles, err
	}
	defer row.Close()

	for row.Next() {
		var nid string
		var title string
		var date string
		if err = row.Scan(&nid, &title, &date); err != nil {
			log.Println("获取通知表字段失败", err)
			return arrArticles, err
		}
		//log.Println(nid, title, date)
		bnew := !strings.Contains(readids, fmt.Sprintf(`,%s,`, nid))
		//不需要新建
		if len(arrArticles) > 0 && arrArticles[len(arrArticles)-1].Date == date[:10] {
			st := arrArticles[len(arrArticles)-1]
			st.ArrLines = append(st.ArrLines, &STItem{Nid: nid, Name: date[11:19] + " " + title, OK: bnew})
		} else {
			st := &STArticlelList{}
			st.Date = date[:10]
			st.ArrLines = append(st.ArrLines, &STItem{Nid: nid, Name: date[11:19] + " " + title, OK: bnew})
			arrArticles = append(arrArticles, st)
		}
	}
	return arrArticles, nil
}

//加载流程数据
func LoadOrderData(uid string) ([]*STArticlelList, error) {
	arrArticles := []*STArticlelList{}

	row, err := GDBpgopt.Query(fmt.Sprintf(SQL_LOAD_LEAVEINFO, uid))
	if err != nil {
		return arrArticles, err
	}
	defer row.Close()

	for row.Next() {
		var createtime string
		if err = row.Scan(&createtime); err != nil {
			log.Println("获取请假表字段失败", err)
			return arrArticles, err
		}
		//不需要新建
		if len(arrArticles) > 0 && arrArticles[len(arrArticles)-1].Date == createtime[:10] {
			st := arrArticles[len(arrArticles)-1]
			st.ArrLines = append(st.ArrLines, &STItem{Name: createtime[11:] + " 请假申请", OK: true})
		} else {
			st := &STArticlelList{}
			st.Date = createtime[:10]
			st.ArrLines = append(st.ArrLines, &STItem{Name: createtime[11:] + " 请假申请", OK: true})
			arrArticles = append(arrArticles, st)
		}
	}
	return arrArticles, nil
}
