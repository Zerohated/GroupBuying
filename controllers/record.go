package controllers

import (
	"encoding/json"
	"group_buying/models"
	"strconv"

	"github.com/astaxie/beego"
)

// RecordController about User
type RecordController struct {
	beego.Controller
}

// @Title Add Record
// @Description Add one record
// @Param	body body 	object	true "<p>{<br>&quot;activityId&quot;:1,<br>&quot;openId&quot;: &quot;xxxx&quot;,<br>&quot;paidId&quot;: &quot;me&quot;,<br>&quot;paidState&quot;: &quot;yyyy&quot;,<br>&quot;paidAmount&quot;: &quot;yyyy&quot;<br>}</p>"
// @Success 200	{"state":"success"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router / [post]
func (controller *RecordController) AddRecord() {
	record := make(map[string]interface{})
	json.Unmarshal(controller.Ctx.Input.RequestBody, &record)
	err := models.AddRecord(record)
	if err != nil {
		controller.CustomAbort(403, err.Error())
	}
	controller.Data["json"] = map[string]string{"state": "success"}
	controller.ServeJSON()
}

// @Title Refund
// @Description get all refund record
// @Param	activityId		path 	string	true	"the activityId you want to get"
// @Success 200 {<br>"record": [record1,record2,...],<br>"count":2<br>}
// @Failure 400 "Error infomation"
// @router /refund/:activityId [get]
func (controller *ActivitiesController) Refund() {
	activityId := controller.Ctx.Input.Param(":activityId")
	aId, _ := strconv.Atoi(activityId)
	obs := models.GetRefundRecords(aId)
	controller.Data["json"] = obs
	if _, ok := obs["error"]; ok {
		controller.CustomAbort(403, obs["error"].(string))
	}
	controller.ServeJSON()
}
