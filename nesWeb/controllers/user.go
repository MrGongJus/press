package controllers

import ("github.com/astaxie/beego"
	_ "nesWeb/models"
	"github.com/astaxie/beego/orm"
	"nesWeb/models"
)
type UserControllers struct {
	beego.Controller
}
//展示注册页面
func (this *UserControllers) ShowTpl() {
	this.TplName = "register.html"
}
//处理注册数据
func (this *UserControllers) Handel() {

	//把数据插入到数据库
	UserName := this.GetString("userName")
	passwd := this.GetString("password")
	if UserName == "" || passwd == "" {
		beego.Info("输入不能为空")
		this.TplName = "register.html"
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = UserName
	user.Pwd = passwd
	//插入数据
	_,err := o.Insert(&user)
	if err != nil{
		beego.Error("注册失败")
		this.TplName = "register.html"
		return
	}
	this.Ctx.WriteString("注册成功")
}
//展示登录页面
func (this *UserControllers)ShowLogin(){
	//获取cookie
	userName := this.Ctx.GetCookie("userName")
	//判断cookie为空，用户名为空
	if userName == ""{
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}else {
		//否则写入数据
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}

	this.TplName = "login.html"
}
//处理登录数据
func (this *UserControllers)HandleLogin(){
	//获取数据
	userNmae := this.GetString("userName")
	pwd := this.GetString("password")
	if userNmae == "" || pwd == "" {
		beego.Error("用户名和密码不能为空")
		this.TplName = "login.html"
		return
	}
	//查询数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userNmae
	//查询， 要指定靠什么查询
	err := o.Read(&user,"Name")
	if err != nil{
		beego.Error("输入错误")
		this.TplName = "login.html"
		return
	}
	//判断密码是否正确
	if user.Pwd != pwd{
		beego.Error("登录失败")
		this.TplName = "login.html"
		return
	}

	//登录成功的情况下，选中复选匡，把用户名存储到Setcookie中
	remember := this.GetString("remember")
	if remember == "on"{
		// 为 on 时 ， 参一， 参二 为建值， 参三是指定的时间
		this.Ctx.SetCookie("userName",userNmae,60*60)
	}else {
		//为空时，为 -1
		this.Ctx.SetCookie("userName",userNmae,-1)
	}
	//获取服务器端的定时器，关闭时失效
	this.SetSession("userName",userNmae)

	//this.Ctx.WriteString("登录成功")
	this.Redirect("/article/index",302)
}
//退出登录状态
func(this*UserControllers)Logout(){
	//推出时, 删除Session
	this.DelSession("userName")
	this.Redirect("/login",302)
}