package controllers

import ("github.com/astaxie/beego"
		"github.com/gomodule/redigo/redis"
)

type RedisGit struct {
	beego.Controller
}
func (this *RedisGit)ShowRedis(){
	//连接函数
	conn,err := redis.Dial("tcp","127.0.0.1:6379")
	if err != nil{
		beego.Error("Dial err",err)
		return
	}
	//操作函数
	//把数据存起来
	//conn.Send("set","aa","bb")
	//刷一下数据
	//conn.Flush()
	//调用执行函数
	//conn.Receive()

	//接口类型
	resh,err := conn.Do("set","aa","bb")

	//回复助手函数
	//获取value为字符串类型
	//result,_ := redis.String(resh,err)
	//beego.Error(result)
	//获取为字符类型
	result,_:=redis.Values(resh,err)

	//获取不同类型数据 用scan 函数
	var v1,v2 string
	var v3 int
	redis.Scan(result,&v1,&v2,&v3)

	//经常访问，不长修改的数据，存入redis
}