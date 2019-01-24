package main

import (
	_ "nesWeb/routers"
	"github.com/astaxie/beego"
	_ "nesWeb/models"
)

func main() {
	beego.AddFuncMap("prePage",ShowprePage)
	beego.AddFuncMap("nextPage",ShownextPage)
	beego.Run()
}

func ShowprePage(pageIndex int)int {
	if pageIndex <= 1{
		return 1
	}
	return pageIndex - 1
}
func ShownextPage(pageIndex int,pageCount float64)int{
	if pageIndex >= int(pageCount){
		return int(pageCount)
	}
	return pageIndex + 1
}