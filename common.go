package esayoa

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"
)

//拦截器,权限验证
func Intercept(level int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie("oakey"); err == nil {
			info, ok := GSessions.CheckCookie(cookie.Value)
			//判断权限
			if ok && info.UserInfo.Auth >= level &&
				strings.Split(c.Request.RemoteAddr, ":")[0] == info.RemoteIP {
				//保存信息
				c.Set("Client", info)
				c.Next()
				return
			}
		}
		//不通过
		c.File("jump.html")
		c.Abort()
	}
}

//格式化文本
func FormatInputText(input string) string {
	input = template.HTMLEscapeString(input)
	input = strings.Replace(input, "-", "_", -1)
	input = strings.Replace(input, " ", "&nbsp;", -1)
	input = strings.Replace(input, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)
	input = strings.Replace(input, "\r\n", "</br>", -1)
	return input
}
