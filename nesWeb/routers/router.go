package routers

import (
	"nesWeb/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"

)

func init() {
    beego.InsertFilter("/article/*", beego.BeforeExec,FilterFunc)
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserControllers{},"get:ShowTpl;post:Handel")
    beego.Router("/login",&controllers.UserControllers{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/article/index",&controllers.ArticleControllers{},"get:ShowIndex")
    beego.Router("/article/addArticle",&controllers.ArticleControllers{},"get:ShowAdd;post:AddArticle")
    beego.Router("/article/content",&controllers.ArticleControllers{},"get:ShowContent")
    beego.Router("/article/UpdateArticle",&controllers.ArticleControllers{},"get:ShowUpdateArticle;post:HandelUpdate")
    beego.Router("/article/deleteAticle",&controllers.ArticleControllers{},"get:ShowDeleteHandle")
    beego.Router("/article/addType",&controllers.ArticleControllers{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/article/logout",&controllers.UserControllers{},"get:Logout")
    beego.Router("/article/deleteType",&controllers.ArticleControllers{},"get:DeleteType")
    beego.Router("/redis",&controllers.RedisGit{},"get:ShowRedis")
}
func FilterFunc(ctx*context.Context){
    userName := ctx.Input.Session("userName")
    if userName == nil {
        ctx.Redirect(302,"/login")
    }
}