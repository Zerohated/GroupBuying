package controllers

import (
	"encoding/json"
	"fmt"
	"group_buying/models"
	"strconv"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

// TicketController about Ticket
type TicketController struct {
	beego.Controller
}

// @Title Create Tickets
// @Description Create Tickets
// @Param	activityId		path 	string	true		"activityId for those Tickets generated"
// @Param	ticketModel		path 	string	true		"ticketModel for those Tickets generated"
// @Param	count		path 	string	true	"How much tickets need to be generated "
// @Param	password		path 	string	true	"password"
// @Success 200 {<br>"count":int,"error":null<br>}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /generation/:activityId/:ticketModel/:count/:password [post]
func (controller *TicketController) GeneratedTickets() {
	activityId := controller.Ctx.Input.Param(":activityId")
	ticketModel := controller.Ctx.Input.Param(":ticketModel")
	count := controller.Ctx.Input.Param(":count")
	password := controller.Ctx.Input.Param(":password")
	if password == "Test123" {
		result := models.GenerateTickets(activityId, ticketModel, count)
		if _, ok := result["error"]; ok {
			controller.CustomAbort(403, result["error"].(string))
		}
		controller.Data["json"] = result
	} else {
		controller.CustomAbort(400, "Wrong Password")
	}
	controller.ServeJSON()
}

// @Title Check Ticket
// @Description Check Ticket
// @Param	ticketCode		path 	string	true		"Ticket.Code"
// @Success 200 {<br>"startDate": "2017/1/1",<br>"endDate": "2017/1/1",<br>"condition": 100,<br>"value": 50,<br>"type": "Voucher",<br>"state": "Legal/Illegal/Used",<br>"error":null<br>}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /:ticketCode/ [get]
func (controller *TicketController) CheckTicket() {
	ticketCode := controller.Ctx.Input.Param(":ticketCode")
	if ticketCode != "" {
		result := models.CheckTicket(ticketCode)
		if _, ok := result["error"]; ok {
			controller.CustomAbort(403, result["error"].(string))
		}
		controller.Data["json"] = result
	} else {
		controller.CustomAbort(400, "ticketCode cannot be null")

	}
	controller.ServeJSON()
}

// @Title Burn Ticket
// @Description Burn Ticket
// @Param	ticketCode		path 	string	true		"Ticket.Code"
// @Success 200 {"state":"succeed"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /:ticketCode/ [delete]
func (controller *TicketController) BurnTicket() {
	ticketCode := controller.Ctx.Input.Param(":ticketCode")
	result := models.BurnTicket(ticketCode)
	if _, ok := result["error"]; ok {
		controller.CustomAbort(403, result["error"].(string))
	}
	controller.Data["json"] = result
	controller.ServeJSON()
}

// @Title CreateOrUpdate TicketModel
// @Description create TicketModel
// @Param	body		body 	object	true "{<br> &quot;type&quot;: &quot;Voucher/Discount/Droit&quot;,<br> &quot;condition&quot;: 100,<br> &quot;picture&quot;: &quot;&quot;,<br> &quot;description&quot;: &quot;&quot;,<br> &quot;startDate&quot;: &quot;2017/10/24&quot;,<br> &quot;endDate&quot;: &quot;2017/10/30&quot;,<br>&quot;value&quot;: int<br>  }"
// @Success 200 {models.TicketModel}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /model/ [post]
func (controller *TicketController) CreateTicketModel() {
	ticketModel := new(models.TicketModel)
	json.Unmarshal(controller.Ctx.Input.RequestBody, &ticketModel)
	result := models.UpdateTicketModel(*ticketModel)
	switch t := result.(type) {
	case error:
		logs.Notice(result.(error).Error())
		controller.CustomAbort(403, result.(error).Error())
	case interface{}:
		controller.Data["json"] = result
	default:
		logs.Notice(fmt.Sprintf("%T", t))
		controller.CustomAbort(403, fmt.Sprintf("%T", t))
	}
	controller.ServeJSON()
}

// @Title Get TicketModel
// @Description find TicketModel by id
// @Param	ticketModelId		path 	string	true		"the ticketModelId you want to get"
// @Success 200 {TicketModel} models.TicketModel
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /model/:ticketModelId [get]
func (controller *TicketController) GetTicketModel() {
	ticketModelId := controller.Ctx.Input.Param(":ticketModelId")
	if ticketModelId != "" {
		activity := models.GetOneTicketModel(ticketModelId)
		controller.Data["json"] = activity
	}
	controller.ServeJSON()
}

// @Title GetAll
// @Description get all TicketModels
// @Success 200  {<br>"ticketModels": [TicketModel1,TicketModel2,...],<br>"count":2,<br>"error":null<br>}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /models/ [get]
func (controller *TicketController) GetAllTicketModels() {
	obs := models.GetAllTicketModels()
	controller.Data["json"] = obs
	controller.ServeJSON()
}

// // @Title Update
// // @Description update the TicketModel
// // @Param	body		body 	object	true "{<br> &quot;id&quot;: 1,<br> &quot;type&quot;: &quot;Voucher/Discount/Droit&quot;,<br> &quot;condition&quot;: 100,<br> &quot;picture&quot;: &quot;&quot;,<br> &quot;description&quot;: &quot;&quot;,<br> &quot;startDate&quot;: &quot;2017/10/24&quot;,<br> &quot;endDate&quot;: &quot;2017/10/30&quot;,<br>&quot;value&quot;: int<br> }"
// // @Success 200 {TicketModel} models.TicketModel
// // @Failure 403 "Error infomation"
// // @Failure 400 "Params Error infomation"
// // @router /ticket_model/ [put]
// func (controller *TicketController) PutTicketModel() {
// 	ticketModel := new(models.TicketModel)
// 	json.Unmarshal(controller.Ctx.Input.RequestBody, &ticketModel)
// 	logs.Debug(ticketModel)
// 	err := models.UpdateTicketModel(*ticketModel)
// 	if err != nil {
// 		controller.CustomAbort(203, err.Error())
// 	} else {
// 		controller.Data["json"] = map[string]string{"updated": "succeed"}
// 	}
// 	controller.ServeJSON()
// }

// @Title GetUsedTickets
// @Description get all Used ticket by activityId, if activityId = 0 then return all used ticket
// @Param	activityId		path	string	true	"the activityId you want to get"
// @Success 200
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /tickets/used/:activityId [get]
func (controller *TicketController) GetUsedTickets() {
	activityId := controller.Ctx.Input.Param(":activityId")
	if activityId == "" {
		controller.CustomAbort(403, "ActivityId can't be null")
	}
	aId, _ := strconv.Atoi(activityId)
	obs := models.GetUsedTickets(aId)
	controller.Data["json"] = obs
	controller.ServeJSON()
}
