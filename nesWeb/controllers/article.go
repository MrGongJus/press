package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"nesWeb/models"
	"math"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"bytes"
)

type ArticleControllers struct {
	beego.Controller
}
// index.html 记录数分页
//指定信息类型
//登录匡的限制
func (this *ArticleControllers)ShowIndex(){
	//登录校验
	userName := this.GetSession("userName")
	if userName == nil{
		this.Redirect("/login",302)
	}
	//获取所有文章数据
	o := orm.NewOrm()
	//获取所有文章
	qs := o.QueryTable("Article")
	var articles []models.Article
	//qs.All(&articles)

    //每两行为以页
	pageSize := 2
	//处理首行内容
	pageIndex,err := this.GetInt("pageIndex")
	if err != nil{
		pageIndex = 1
	}
	//获取数据库中部分数据
	start := pageSize * (pageIndex-1)

	//获取所有代类型的数据传递给前段展示
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"] = articleTypes

	//把类型存入redis中
	//连接函数
	conn,err := redis.Dial("tcp","6379")
	if err != nil{
		beego.Error("Dail err",err)
		return
	}
	defer conn.Close()

	//操作函数
	//conn.Do("set","articleTypes",articleTypes)

	data,err := redis.Bytes(conn.Do("set","articleType",))
	if len(data)==0{
		//序列化和反序列化
		//缓存容器
		var buffer bytes.Buffer
		//解码器
		enc := gob.NewEncoder(&buffer)
		//编码
		enc.Encode(&articleTypes)

		//更新操作
		conn.Do("set","articleTypes",buffer.Bytes())
	}else {
		//解码器
		data2 ,_ := redis.Bytes(conn.Do("get","articleTypes"))

		dec := gob.NewDecoder(bytes.NewReader(data2))
		//解码
		dec.Decode(&articleTypes)
	}

	//返回值COUNT为INT64
	var count int64
	//获取下拉匡数据
	typeName := this.GetString("select")
	//判断下拉匡中的值为空
	if typeName == ""{
		//RelatedSel为惰性查询。
		//指定查询带有类型数据
		count,_ = qs.RelatedSel("ArticleType").Count()
		//截取数据中 所有带有类型数据的数据
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)
	}else{
		//filter 过了， 参一  数据类型 - 类型名， 参二 ： 下拉匡中的类型数据
		//过滤以后得到的数据
		count,_ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
		//截取后过滤的数据
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}
	//ceil天花版函数。
	//总页数
	pageCount := math.Ceil(float64(count)/float64(pageSize))
	//通过key值把数据指定到前段
	this.Data["count"] =count
	this.Data["pageCount"] = pageCount

	this.Data["articles"] = articles
	this.Data["pageIndex"] = pageIndex

	this.TplName = "index.html"
}
//获取所有数据类型数据，展示添加文章页面 ， 主要把数据传递给前段展示
func (this *ArticleControllers)ShowAdd(){
	//获取所有类型的数据传递给前段
	o := orm.NewOrm()
	//切片存储所有数据类型。
	var articleType []models.ArticleType
	//获取数据库中所有带类型数据值
	o.QueryTable("ArticleType").All(&articleType)
	//传递给前段
	this.Data["articleType"] = articleType
	this.TplName = "add.html"
}
//获取add.html中标题，内容，图片数据， 校验图片大小，后缀，是否重名，把数据插入数据库
//文章类型分类，把分类内容写入文章中
func (this *ArticleControllers)AddArticle(){
	//获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	//获取图片数据 ， GetFile
	file,mul,err := this.GetFile("uploadname")
	//校验数据
	if  articleName == "" || content == "" || err != nil {
		beego.Error("file err",err)
		this.TplName = "add.html"
		return
	}
	//关闭数据
	defer file.Close()
	//文件大小
	if mul.Size > 5000000{
		beego.Error("文件太大，无法加载")
		this.TplName = "add.html"
		return
	}

	//文件后缀
	ext := path.Ext(mul.Filename)
	if ext != ".jpg" && ext != ".png" && ext != "jpeg"{
		beego.Error("无法识别")
		this.TplName = "add.html"
		return
	}
	//文件重名
	filepath := time.Now().Format("2016-01-02 15:04:05")
	//SaveToFile , 图片name值 ，  存储路径
	//保存数据
	this.SaveToFile("uploadname","./static/img/"+filepath+ext)
	//把数据插入到数据库
	o := orm.NewOrm()
	var article models.Article
	article.Content = content
	article.Img = "/static/img/"+filepath+ext
	article.Title = articleName
	article.ArticleType=new(models.ArticleType)
	//获取下拉框的数据
	TypeName := this.GetString("select")
	//获取类型对象
	var articleType models.ArticleType
	articleType.TypeName = TypeName
	//读取指定内容
	o.Read(&articleType,"TypeName")
	//把类型对象插入到文章中
	article.ArticleType = &articleType
	//插入
	_,err =o.Insert(&article)
	if err != nil{
		beego.Info(err)
		this.TplName = "add.html"
		return
	}
	this.Redirect("/article/index",302)
}
//展示文章详情页
func (this *ArticleControllers)ShowContent(){
	//通过id获取数据
	id,err := this.GetInt("id")
	if err != nil{
		beego.Error("id err",err)
		this.TplName = "index.html"
		return
	}
	//查询数据库，获取文章信息
	o := orm.NewOrm()
	var article models.Article
	article.Id2 = id
	//查询
	o.Read(&article)

	//多对多查询两种方式
	//不能去重
	//o.LoadRelated(&article,"Users")
	//能去重
	//存储所有用户信息
	var users []models.User
	//把所有用户信息过滤去重，返回前端
	o.QueryTable("User").Filter("Articles__Article__Id2",article.Id2).Distinct().All(&users)
	this.Data["users"] = users

	//每次访问数据阅读次数加一
	article.ReadCount += 1
	//更新数据库数据
	o.Update(&article)
	//返回数据
	this.Data["article"] = article

	//添加浏览记录， 多对多插入
	m2m:=o.QueryM2M(&article,"Users")
	//向表里面插入对象指针
	var user models.User
	userName := this.GetSession("userName")
	//类型断言
	user.Name = userName.(string)
	o.Read(&user,"Name")
	//插入对象
	m2m.Add(user)
	this.TplName = "content.html"
}
//展示编辑文章页面
func (this *ArticleControllers) ShowUpdateArticle() {
	//填充文章原来的数据
	id,err := this.GetInt("id")
	if err != nil{
		beego.Error("showUpdateArticle err",err)
		this.TplName = "update.html"
		return
	}
	//获取对象
	o := orm.NewOrm()
	var article models.Article
	article.Id2 = id
	//查询
	o.Read(&article)
	//返回数据
	this.Data["article"] = article
	this.TplName = "update.html"
}
//封装函数， 写接口
func UploadFile(this *ArticleControllers,fileName string)string{
	file,mul,err := this.GetFile(fileName)
	defer file.Close()
	if err != nil{
		beego.Error("file err",err)
		this.TplName = "add.html"
		return ""
	}
	//文件大小
	if mul.Size > 5000000{
		beego.Error("文件太大，无法加载")
		this.TplName = "add.html"
		return ""
	}

	//文件后缀
	ext := path.Ext(mul.Filename)
	if ext != ".jpg" && ext != ".png" && ext != "jpeg"{
		beego.Error("无法识别")
		this.TplName = "add.html"
		return ""
	}
	//文件重名
	filePath := time.Now().Format("2016-01-02 15:04:05")

	this.SaveToFile(fileName,"./static/img/"+filePath+ext)

	return "/static/img/"+filePath+ext
}
//处理编辑数据
func(this*ArticleControllers)HandelUpdate(){
	//获取数据
	id,err := this.GetInt("id")
	title := this.GetString("articleName")
	content := this.GetString("content")
	filepath := UploadFile(this,"uploadname")
	//校验数据
	if err != nil || title == "" || content == "" || filepath == ""{
		beego.Error("handelUpdate err",err)
		this.TplName = "update.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id2 = id
	//读取数据库是否有数据，在更新
	err = o.Read(&article)
	if err != nil{
		beego.Error("read err",err)
		this.TplName = "update.html"
		return
	}
	//赋值
	article.Content = content
	article.Title = title
	article.Img = filepath
	//更新所有数据
	o.Update(&article)
	//跳转到首页
	this.Redirect("/article/index",302)
}
//处理删除文章
func (this *ArticleControllers)ShowDeleteHandle(){
	//获取要删除文章id
	id,err := this.GetInt("id")
	if err != nil{
		beego.Error("delect err",err)
		this.TplName = "index.html"
		return
	}
	//删除操作
	o := orm.NewOrm()
	var article models.Article
	article.Id2 = id
	//删除
	_,err = o.Delete(&article)

	if err != nil{
		beego.Error("Id2 err",err)
		this.TplName = "index.html"
		return
	}
	//跳转到首页
	this.Redirect("/article/index",302)
}
//展示添加文章类型页面
func(this *ArticleControllers)ShowAddType(){
	//获取所有类型
	o := orm.NewOrm()
	qs := o.QueryTable("ArticleType")
	var ArticleTypes []models.ArticleType
	qs.All(&ArticleTypes)
	//数据传递给前端
	this.Data["ArticleTypes"] = ArticleTypes
	this.TplName = "addType.html"
}
//处理添加类型数据
func (this *ArticleControllers)HandleAddType(){
	//获取数据
	typeName := this.GetString("typeName")
	if typeName == ""{
		beego.Error("typeName err ")
		this.TplName = "addType.html"
		return
	}
	//向数据库中插入数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)
	//跳转
	this.Redirect("/article/addType", 302)
}
//删除文章类型
func (this *ArticleControllers)DeleteType(){
	//通过id获取数据
	id,err := this.GetInt("id")
	if err != nil{
		beego.Error("delete err",err)
		this.TplName = "addType.html"
		return
	}
	//删除数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType)
	//跳转
	this.Redirect("/article/addType",302)
}