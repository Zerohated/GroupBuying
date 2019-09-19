package controllers

import (
	"encoding/json"
	"group_buying/models"
	"strconv"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

// ActivityController Operations about Activity
type ActivityController struct {
	beego.Controller
}

// @Title Create Activity
// @Description Use Json to input name(str),price(float),groupSize(int),limitCount(int),startDate(str),endDate(str)
// @Param body body object	true "<p>{<br>&quot;name&quot;: &quot;name&quot;,<br>&quot;price&quot;: 0.0,<br>&quot;groupSize&quot;: 0,<br>&quot;limitCount&quot;: 0,<br>&quot;startDate&quot;: &quot;2017-01-01&quot;,<br>&quot;endDate&quot;: &quot;2017-12-31&quot;<br>}</p>"
// @Success 200 {"activityId":Activity.Id}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router / [post]
func (controller *ActivityController) CreateActivity() {
	request := make(map[string]interface{})
	json.Unmarshal(controller.Ctx.Input.RequestBody, &request)
	id, err := models.AddActivity(request)
	if err != nil {
		logs.Notice(err.Error())
		controller.CustomAbort(403, err.Error())
	}
	controller.Data["json"] = map[string]int{"activityId": id}
	controller.ServeJSON()
}

// @Title Create ActivityUi
// @Description Use Json to input background,button,detailButton,successButton,description,detail,keyVisual,endNotice,notSuccess,successTop
// @Param	activityId		path 	string	true	"the activityId you want to get"
// @Param body body object	true "<p>{<br>&quot;key&quot;: &quot;value&quot;,<br>...<br><br>&quot;key&quot;: &quot;value&quot;<br>}</p>"
// @Success 200 {"activityId":ActivityUi.Id(is as same as Activity's id)}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /ui/:activityId [post]
func (controller *ActivityController) CreateActivityUI() {
	activityId := controller.Ctx.Input.Param(":activityId")
	id, _ := strconv.Atoi(activityId)
	activityUi := new(models.ActivityUi)
	json.Unmarshal(controller.Ctx.Input.RequestBody, &activityUi)
	activityUi.Id = id
	id, err := models.AddActivityUi(*activityUi)
	if err != nil {
		logs.Notice(err.Error())
		controller.CustomAbort(403, err.Error())
	}
	controller.Data["json"] = map[string]int{"activityId": id}
	controller.ServeJSON()
}

// @Title Get Activity
// @Description Return general information including UI
// @Param	activityId		path 	string	true	"the activityId you want to get"
// @Success 200 {Activity} activityUI,name,price,groupSize,limitCount,existCount,startDate,endDate,isEnd
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /:activityId [get]
func (controller *ActivityController) GetActivity() {
	activityId := controller.Ctx.Input.Param(":activityId")
	if activityId != "" {
		id, _ := strconv.Atoi(activityId)
		result := models.GetActivity(id)
		if _, ok := result["error"]; ok {
			controller.CustomAbort(403, result["error"].(string))
		}
		controller.Data["json"] = result
	} else {
		logs.Notice("ActivityId can't be empty")
		controller.CustomAbort(400, "ActivityId can't be empty")
	}
	controller.ServeJSON()
}

// @Title Get Activity Detail
// @Description Detail information about groups,ticketModels,tickets
// @Param	activityId		path 	string	true	"the activityId you want to get"
// @Success 200 {Activity} groups,ticketModels,tickets,name,price,groupSize
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /detail/:activityId [get]
func (controller *ActivityController) GetActivityDetail() {
	activityId := controller.Ctx.Input.Param(":activityId")
	if activityId != "" {
		id, _ := strconv.Atoi(activityId)
		result := models.GetActivityDetail(id)
		controller.Data["json"] = result
	} else {
		logs.Notice("ActivityId can't be empty")
		controller.CustomAbort(400, "ActivityId can't be empty")
	}
	controller.ServeJSON()
}

// @Title Update Activity
// @Description update the Activity
// @Param	activityId		path 	string	true	"the activityId you want to get"
// @Param body body object	true "<p>{<br>&quot;price&quot;: 0.0,<br>&quot;groupSize&quot;: 0,<br>&quot;limitCount&quot;: 0,<br>&quot;startDate&quot;: &quot;2017-01-01&quot;,<br>&quot;endDate&quot;: &quot;2017-12-31&quot;<br>}</p>"
// @Success 200 {"updated": "succeed"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /:activityId [put]
func (controller *ActivityController) UpdateActivity() {
	request := make(map[string]interface{})
	activityId := controller.Ctx.Input.Param(":activityId")
	if activityId != "" {
		json.Unmarshal(controller.Ctx.Input.RequestBody, &request)
		i, convErr := strconv.Atoi(activityId)
		if convErr != nil {
			logs.Notice("Id is required as a number")
			controller.CustomAbort(400, "Id is required as a number")
		}
		request["id"] = i
		err := models.UpdateActivity(request)
		if err != nil {
			logs.Notice(err.Error())
			controller.CustomAbort(403, err.Error())
		} else {
			controller.Data["json"] = map[string]string{"updated": "succeed"}
			controller.ServeJSON()
		}
	} else {
		logs.Notice("Id is required as a number")
		controller.CustomAbort(400, "Id is required as a number")
	}
}

// @Title Update ActivityUi
// @Description Use Json to input background,button,detailButton,successButton,description,detail,keyVisual,endNotice,notSuccess,successTop
// @Param	activityId		path 	string	true	"the activityId you want to get"
// @Param body body object	true "<p>{<br>&quot;key&quot;: &quot;value&quot;,<br>...<br><br>&quot;key&quot;: &quot;value&quot;<br>}</p>"
// @Success 200 {"activityId":ActivityUi.Id(is as same as Activity's id)}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /ui/:activityId [put]
func (controller *ActivityController) UpdateActivityUI() {
	request := make(map[string]interface{})
	activityId := controller.Ctx.Input.Param(":activityId")
	if activityId != "" {
		json.Unmarshal(controller.Ctx.Input.RequestBody, &request)
		i, convErr := strconv.Atoi(activityId)
		if convErr != nil {
			controller.CustomAbort(400, "Id is required as a number")
		}
		request["id"] = i
		err := models.UpdateActivityUi(request)
		if err != nil {
			controller.CustomAbort(403, err.Error())
		} else {
			controller.Data["json"] = map[string]string{"updated": "succeed"}
			controller.ServeJSON()
		}
	} else {
		controller.CustomAbort(400, "Id is required as a number")
	}
}

// @Title Delete Activity
// @Description Delete the Activity
// @Param	activityId		path 	string	true	"the activityId you want to get"
// @Success 200 {"state": "succeed"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /:activityId [delete]
func (controller *ActivityController) DeleteActivity() {
	activityId := controller.Ctx.Input.Param(":activityId")
	if activityId != "" {
		i, err := strconv.Atoi(activityId)
		if err != nil {
			controller.CustomAbort(403, "Id is required as a number")
		}
		result := models.DeleteActivity(i)
		controller.Data["json"] = result
		controller.ServeJSON()
	}
	controller.CustomAbort(400, "Id is required as a number")
}

// @Title Add or Update Activity's TicketModels
// @Description Add Activity's TicketModels
// @Param	activityId		path 	string	true		"the activityId you want to modify"
// @Param	ticketModelId		path 	string	true		"the ticketModelId you want to modify"
// @Param	body		body 	object	true "{&quot;is_amust&quot;: true, &quot;useDetail&quot;: &quot;blablabla&quot;,&quot;startDate&quot;: &quot;2017-11-11&quot;,&quot;endDate&quot;:&quot;2017-11-11&quot;}"
// @Success 200 {"state": "succeed"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /models/:activityId/:ticketModelId/ [post]
func (controller *ActivityController) AddActivityTicketModels() {
	info := make(map[string]interface{})
	aId := controller.Ctx.Input.Param(":activityId")
	tId := controller.Ctx.Input.Param(":ticketModelId")
	if aId != "" && tId != "" {
		activityId, err := strconv.Atoi(aId)
		if err != nil {
			controller.CustomAbort(403, "Id is required as a number")
		}
		ticketModelId, err := strconv.Atoi(tId)
		if err != nil {
			controller.CustomAbort(403, "Id is required as a number")
		}
		json.Unmarshal(controller.Ctx.Input.RequestBody, &info)
		result := models.AddOrUpdateActivityTicketModels(activityId, ticketModelId, info)
		if _, ok := result["error"]; ok {
			controller.CustomAbort(403, result["error"])
		}
		// result := "empty"
		controller.Data["json"] = result
	} else {
		controller.CustomAbort(400, "ActivityId cannot be null")
	}
	controller.ServeJSON()
}

// @Title Delete Activity's TicketModels
// @Description Delete Activity's TicketModels
// @Param	activityId		path 	string	true		"the activityId you want to modify"
// @Param	ticketModelId		path 	string	true		"the ticketModelId you want to modify"
// @Success 200 {"state": "succeed"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /models/:activityId/:ticketModelId/ [delete]
func (controller *ActivityController) DeleteActivityTicketModels() {
	aId := controller.Ctx.Input.Param(":activityId")
	tId := controller.Ctx.Input.Param(":ticketModelId")
	if aId != "" && tId != "" {
		activityId, err := strconv.Atoi(aId)
		if err != nil {
			controller.CustomAbort(403, "Id is required as a number")
		}
		ticketModelId, err := strconv.Atoi(tId)
		if err != nil {
			controller.CustomAbort(403, "Id is required as a number")
		}
		result := models.DeleteActivityTicketModels(activityId, ticketModelId)
		controller.Data["json"] = result
		controller.ServeJSON()
	}
	controller.CustomAbort(400, "Id is required as a number")
}
