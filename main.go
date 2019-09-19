package main

import (
	_ "group_buying/routers"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// orm.RegisterDataBase("default", "mysql", "root:Start123@/dev?charset=utf8&loc=Asia%2FShanghai")
	orm.RegisterDataBase("default", "mysql", "root:Start123@/dev?charset=utf8&loc=Asia%2FShanghai")
	orm.DefaultTimeLoc = time.UTC
}

func main() {
	logs.SetLogger("console")
	logs.Async()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/dev"] = "swagger"
		logs.SetLogger(logs.AdapterFile, `{"filename":"logs/dev.log","level":7,"daily":true,"maxdays":10}`)
	}
	if beego.BConfig.RunMode == "prod" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/dev"] = "swagger"
		logs.SetLogger(logs.AdapterFile, `{"filename":"logs/project.log","level":7,"daily":true,"maxdays":10}`)
	}
	orm.RunCommand() //orm command
	beego.Run()
	orm.Debug = true
}
