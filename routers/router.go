// @APIVersion 2.0.0
// @Title GroupBuying Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact david.chou93@gmail.com
package routers

import (
	"group_buying/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	ns := beego.NewNamespace("/v2",
		beego.NSNamespace("/activity",
			beego.NSInclude(
				&controllers.ActivityController{},
			),
		),
		beego.NSNamespace("/record",
			beego.NSInclude(
				&controllers.RecordController{},
			),
		),
		beego.NSNamespace("/activities",
			beego.NSInclude(
				&controllers.ActivitiesController{},
			),
		),
		beego.NSNamespace("/ticket",
			beego.NSInclude(
				&controllers.TicketController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/users",
			beego.NSInclude(
				&controllers.UsersController{},
			),
		),
	)
	beego.AddNamespace(ns)
	// support ajax across domains
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
}
