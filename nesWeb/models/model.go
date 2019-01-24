package models

import (

	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	//名字唯一，不能重复
	Name string   `orm:"unique"`
	Pwd string
	Articles []*Article `orm:"rel(m2m)"`
}

type Article struct {
	//设置主建，自动增长，beego默认规定，没有设置主建是，名为ID，数据类型为int作为默认主建
	Id2 int  `orm:"pk;auto"`
	Title  string  `orm:"size(40)"`
	Content string   `orm:"size(100)"`
	ReadCount  int   `orm:"default(0)"`
	Time  time.Time   `orm:"type(datetime);auto_now_add"`
	Img  string     `orm:"null"`
	ArticleType   *ArticleType  `orm:"rel(fk);auto_delete(set_null);null"`
	Users  []*User `orm:"reverse(many)"`
}
type ArticleType struct {
	Id int
	TypeName string
	Articles []*Article `orm:"reverse(many)"`
}
func init(){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/news?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}