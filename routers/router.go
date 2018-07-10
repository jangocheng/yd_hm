package routers

import (
	"yd_hm/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"fmt"
	"strings"
)

//登录验证/api路由权限验证
var FilterUser = func(ctx *context.Context) {
	//用户在线id
	uid,ok := ctx.GetSecureCookie(beego.AppConfig.String("secertkey"),"uid")
	fmt.Println("用户id---->"+uid)
	if !ok {//&& ctx.Request.RequestURI != "/login/login" {
		ctx.Redirect(302, "/login/login")	//重定向返回json{status:0}
	}else if ctx.Request.RequestURI != ("/permission/getpermissionbyid") {
		//用户路由权限判断
		url := ctx.Input.URL()
		//fmt.Println(ctx.Input.URL())	//uri可以获得路由后面的参数,url则不包含
		apilist,_ := ctx.GetSecureCookie(beego.AppConfig.String("secertkey"),"apipermissionbyid")
		//fmt.Println(apilist)
		//var p []*models.Permission
		//json.Unmarshal([]byte(apilist),&p)
		b := strings.Contains(apilist,url) //路由权限判断
		if b {
			fmt.Println("uid---->"+uid+"====>"+url+"路由权限通过")
		}else {
			fmt.Println("uid---->"+uid+"====>"+url+"路由权限未通过")
			ctx.Redirect(302, "/error")
		}
	}
}

func init() {
	beego.InsertFilter("/*",beego.BeforeRouter,FilterUser)
	beego.AutoRouter(&controllers.LoginController{})
	beego.AutoRouter(&controllers.JmlistController{})
	beego.AutoRouter(&controllers.JmuserController{})
	beego.AutoRouter(&controllers.RoleController{})
	beego.AutoRouter(&controllers.ApiController{})
	beego.AutoRouter(&controllers.PermissionController{})
}
