package esayoa

import (
	"fmt"
	"hash/crc32"
	"log"
	"sync"
	"time"
)

const (
	SESSION_TIMEOUT = 600
)

type SessionInfo struct {
	LastTime   int64
	RemoteIP   string
	UID        string
	UserInfo   *STUserBaseInfo
	UserInfoEx *STUserInfoEx
}

type Sessions struct {
	MapSession map[string]*SessionInfo //cookie->info
	SessMutex  sync.Mutex
}

var GSessions = &Sessions{}

func (this *Sessions) Init() error {
	this.MapSession = make(map[string]*SessionInfo)

	go this.SessionSchuld()
	return nil
}

func (this *Sessions) ClearSession() {
	this.SessMutex.Lock()
	defer this.SessMutex.Unlock()

	timeNow := time.Now().Unix()
	for k, st := range this.MapSession {
		if timeNow-st.LastTime > SESSION_TIMEOUT {
			log.Println("登陆超时：", st.UID)
			delete(this.MapSession, k)
		}
	}
}

func (this *Sessions) SessionSchuld() {
	ticket := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticket.C:
			this.ClearSession()
		}
	}
}

//添加新会话
func (this *Sessions) AddSession(uid, rempteip string, info *STUserBaseInfo, exInfo *STUserInfoEx) string {
	nowTime := time.Now().Unix()
	cookie := fmt.Sprintf(`%v%v`, crc32.ChecksumIEEE([]byte(uid)), nowTime)
	log.Println(rempteip, "Cookie:", cookie)
	sinfo := &SessionInfo{
		LastTime:   nowTime,
		RemoteIP:   rempteip,
		UID:        uid,
		UserInfo:   info,
		UserInfoEx: exInfo}

	this.SessMutex.Lock()
	defer this.SessMutex.Unlock()
	this.MapSession[cookie] = sinfo
	return cookie
}

//检查并获取session
func (this *Sessions) CheckCookie(cookie string) (*SessionInfo, bool) {
	this.SessMutex.Lock()
	defer this.SessMutex.Unlock()
	info := this.MapSession[cookie]
	if info != nil {
		info.LastTime = time.Now().Unix()
		return info, true
	}
	return nil, false
}
