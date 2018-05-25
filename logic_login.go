package esayoa

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	ERR_OK             = 0
	ERR_LOGIN_USER_FMT = 1 //用户名格式不正确
	ERR_LOGIN_PASS_FMT = 2 //密码格式不正确
	ERR_LOGIN_USERPASS = 3 //用户或密码错误
	ERR_LOGIN_SYSTEM   = 5 //系统繁忙
)

type STUserBaseInfo struct {
	Id         int    //唯一ID
	Uid        string //用户名
	Pass       string //密码
	Auth       int    //权限
	RegistTime string // 注册时间
}

type STUserInfoEx struct {
	HeaderPic string //头像名字
	UserName  string //名字
	Phone     string //联系方式
	Number    string //员工编号
}

//检查登陆用户密码
func CheckLogin(user, pass string) (int, string, *STUserBaseInfo) {
	//处理用户和密码
	if user != "adminesayoa" && !regexp.MustCompile(`^1[0-9]{10}$`).MatchString(user) {
		return ERR_LOGIN_USER_FMT, "用户格式不正确!", nil
	}
	if !regexp.MustCompile(`^[\x21-\x7e]{6,16}$`).MatchString(pass) {
		return ERR_LOGIN_PASS_FMT, "密码格式不正确!", nil
	}

	pass = fmt.Sprintf(`%x`, md5.Sum([]byte(pass)))
	pass = fmt.Sprintf(`%x`, md5.Sum([]byte("zcw"+pass+"esayoa")))

	//根据user查询出个人信息
	uinfo, err := LoadUserBaseInfo(user)
	log.Println(uinfo, err)
	if err != nil {
		return ERR_LOGIN_SYSTEM, "系统忙，请稍后再试.", nil
	}
	if uinfo != nil && len(uinfo.Pass) == 0 {
		return ERR_OK, "", uinfo
	}
	if uinfo == nil || pass != uinfo.Pass {
		return ERR_LOGIN_USERPASS, "用户或密码不正确!", nil
	}
	return ERR_OK, "", uinfo
}

//加载用户其他信息
func LoadUserInfoEx(id int) (*STUserInfoEx, error) {
	return &STUserInfoEx{}, nil
}

//用户登陆
func UserLogin(c *gin.Context) {
	user := c.Query("user")
	pass := c.Query("pass")
	remote := strings.Split(c.Request.RemoteAddr, ":")[0]
	log.Println("用户登陆", user, pass, remote)
	//登陆
	errCode, errStr, uinfo := CheckLogin(user, pass)
	if errCode != 0 {
		log.Println("登陆失败,错误代码：", errCode)
		c.JSON(http.StatusOK, gin.H{
			"error":   errCode,
			"message": errStr,
		})
		return
	}
	//加载扩展信息
	exinfo, err := LoadUserInfoEx(uinfo.Id)
	if err != nil {
		log.Println("加载用户扩展信息失败.", err)
		c.JSON(http.StatusOK, gin.H{
			"error":   ERR_LOGIN_SYSTEM,
			"message": "系统忙，请稍后再试.",
		})

		return
	}
	//保存会话
	cookieStr := GSessions.AddSession(user, remote, uinfo, exinfo)
	http.SetCookie(c.Writer, &http.Cookie{Name: "oakey", Value: cookieStr, Path: "/"})

	c.JSON(http.StatusOK, gin.H{
		"error":   ERR_OK,
		"message": "/manager",
	})
}

func LoginSetCallBack(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.File("./login.html")
	})
	r.GET("/login.do", UserLogin)
}
