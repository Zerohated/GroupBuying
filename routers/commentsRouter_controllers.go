package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["group_buying/controllers:ActivitiesController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivitiesController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivitiesController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivitiesController"],
		beego.ControllerComments{
			Method: "Dashboard",
			Router: `/dashboard/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivitiesController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivitiesController"],
		beego.ControllerComments{
			Method: "Refund",
			Router: `/refund/:activityId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivitiesController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivitiesController"],
		beego.ControllerComments{
			Method: "Shutdown",
			Router: `/test/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivityController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivityController"],
		beego.ControllerComments{
			Method: "CreateActivity",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivityController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivityController"],
		beego.ControllerComments{
			Method: "GetActivity",
			Router: `/:activityId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivityController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivityController"],
		beego.ControllerComments{
			Method: "UpdateActivity",
			Router: `/:activityId`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivityController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivityController"],
		beego.ControllerComments{
			Method: "DeleteActivity",
			Router: `/:activityId`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivityController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivityController"],
		beego.ControllerComments{
			Method: "GetActivityDetail",
			Router: `/detail/:activityId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivityController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivityController"],
		beego.ControllerComments{
			Method: "AddActivityTicketModels",
			Router: `/models/:activityId/:ticketModelId/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivityController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivityController"],
		beego.ControllerComments{
			Method: "DeleteActivityTicketModels",
			Router: `/models/:activityId/:ticketModelId/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivityController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivityController"],
		beego.ControllerComments{
			Method: "UpdateActivityUI",
			Router: `/ui/:activityId`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:ActivityController"] = append(beego.GlobalControllerRouter["group_buying/controllers:ActivityController"],
		beego.ControllerComments{
			Method: "CreateActivityUI",
			Router: `/ui/:activityId`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:RecordController"] = append(beego.GlobalControllerRouter["group_buying/controllers:RecordController"],
		beego.ControllerComments{
			Method: "AddRecord",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:TicketController"] = append(beego.GlobalControllerRouter["group_buying/controllers:TicketController"],
		beego.ControllerComments{
			Method: "CheckTicket",
			Router: `/:ticketCode/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:TicketController"] = append(beego.GlobalControllerRouter["group_buying/controllers:TicketController"],
		beego.ControllerComments{
			Method: "BurnTicket",
			Router: `/:ticketCode/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:TicketController"] = append(beego.GlobalControllerRouter["group_buying/controllers:TicketController"],
		beego.ControllerComments{
			Method: "GeneratedTickets",
			Router: `/generation/:activityId/:ticketModel/:count/:password`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:TicketController"] = append(beego.GlobalControllerRouter["group_buying/controllers:TicketController"],
		beego.ControllerComments{
			Method: "CreateTicketModel",
			Router: `/model/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:TicketController"] = append(beego.GlobalControllerRouter["group_buying/controllers:TicketController"],
		beego.ControllerComments{
			Method: "GetTicketModel",
			Router: `/model/:ticketModelId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:TicketController"] = append(beego.GlobalControllerRouter["group_buying/controllers:TicketController"],
		beego.ControllerComments{
			Method: "GetAllTicketModels",
			Router: `/models/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:TicketController"] = append(beego.GlobalControllerRouter["group_buying/controllers:TicketController"],
		beego.ControllerComments{
			Method: "GetUsedTickets",
			Router: `/tickets/used/:activityId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:UserController"] = append(beego.GlobalControllerRouter["group_buying/controllers:UserController"],
		beego.ControllerComments{
			Method: "PutUser",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:UserController"] = append(beego.GlobalControllerRouter["group_buying/controllers:UserController"],
		beego.ControllerComments{
			Method: "DeleteUser",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:UserController"] = append(beego.GlobalControllerRouter["group_buying/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUser",
			Router: `/:activityId/:key`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:UserController"] = append(beego.GlobalControllerRouter["group_buying/controllers:UserController"],
		beego.ControllerComments{
			Method: "AddTicketOwner",
			Router: `/:userId/tickets`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:UserController"] = append(beego.GlobalControllerRouter["group_buying/controllers:UserController"],
		beego.ControllerComments{
			Method: "AddNormalUser",
			Router: `/group_user/:groupId/:openId`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:UserController"] = append(beego.GlobalControllerRouter["group_buying/controllers:UserController"],
		beego.ControllerComments{
			Method: "AddStarter",
			Router: `/starter/:activityId/:openId`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:UsersController"] = append(beego.GlobalControllerRouter["group_buying/controllers:UsersController"],
		beego.ControllerComments{
			Method: "GetAllUsers",
			Router: `/:activityId/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["group_buying/controllers:UsersController"] = append(beego.GlobalControllerRouter["group_buying/controllers:UsersController"],
		beego.ControllerComments{
			Method: "GetGroupInfo",
			Router: `/:activityId/group/:groupId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
