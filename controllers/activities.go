package controllers

import (
	"encoding/json"
	"group_buying/models"

	"github.com/astaxie/beego"
)

// ActivitiesController Operations about Activities
type ActivitiesController struct {
	beego.Controller
}

// @Title GetAll
// @Description get all Activities *Param "state" : "未发布/未开始/进行中/已结束-团满/已结束-过期"
// @Success 200 {<br>"activities": [Activity1,Activity2,...],<br>"count":2,<br>"error":null<br>}
// @Failure 400 "Error infomation"
// @router / [get]
func (controller *ActivitiesController) GetAll() {
	obs := models.GetAllActivities()
	controller.Data["json"] = obs
	if _, ok := obs["error"]; ok {
		controller.CustomAbort(403, obs["error"].(string))
	}
	controller.ServeJSON()
}

// @Title Dashboard
// @Description get all Activities Dashboard. limitCount=拼团上限，successGroupCount=成功团数, userCount=参与人数, successUserCount=成功人数, ticketCount=发放卡券, ticketUsedCount=已核销卡券
// @Success 200 {<br>"activities": [Activity1,Activity2,...],<br>"count":2<br>}
// @Failure 400 "Error infomation"
// @router /dashboard/ [get]
func (controller *ActivitiesController) Dashboard() {
	obs := models.Dashboard()
	controller.Data["json"] = obs
	if _, ok := obs["error"]; ok {
		controller.CustomAbort(403, obs["error"].(string))
	}
	controller.ServeJSON()
}

// @Title Test
// @Description None
// @Param body body object	true "<p>{}</p>"
// @Success 200 {}
// @Failure 400 "Error infomation"
// @router /test/ [post]
func (controller *ActivitiesController) Shutdown() {
	request := make(map[string]interface{})
	json.Unmarshal(controller.Ctx.Input.RequestBody, &request)
	if request["data"] == "panic" {
		panic("test")
	}
	controller.Data["json"] = map[string]string{}
	controller.ServeJSON()

}
