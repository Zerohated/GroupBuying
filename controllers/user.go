package controllers

import (
	"encoding/json"
	"group_buying/models"
	"strconv"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

// UserController about User
type UserController struct {
	beego.Controller
}

// @Title Add Starter
// @Description Add User as Starter
// @Param	activityId	path 	int		true	"the activityId where the Starter belong"
// @Param	openId		path 	string	true	"The Starter's OpenId"
// @Success 200	json:{"userId", "groupId", "newActivityExistCount"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /starter/:activityId/:openId [post]
func (controller *UserController) AddStarter() {
	activityId := controller.Ctx.Input.Param(":activityId")
	openId := controller.Ctx.Input.Param(":openId")
	if activityId != "" && openId != "" {
		aId, err := strconv.Atoi(activityId)
		if err != nil {
			logs.Notice(err.Error())
		}
		result := models.AddStarter(aId, openId)
		if _, ok := result["error"]; ok {
			controller.CustomAbort(403, result["error"])
		}
		controller.Data["json"] = result
	} else {
		switch {
		case activityId == "":
			controller.CustomAbort(400, "ActivityId cannot be null")
		case openId == "":
			controller.CustomAbort(400, "openId cannot be null")
		default:
			controller.CustomAbort(400, "ActivityId & openId cannot be null")
		}
	}
	controller.ServeJSON()
}

// @Title Add NormalUser
// @Description Add User as NormalUser
// @Param	groupId		path 	int		true	"the groupId where the User belong"
// @Param	openId		path 	string	true	"The User's OpenId"
// @Success 200 json{"groupId","newGroupSize","userId"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /group_user/:groupId/:openId [post]
func (controller *UserController) AddNormalUser() {
	gId := controller.Ctx.Input.Param(":groupId")
	openId := controller.Ctx.Input.Param(":openId")
	groupId, err := strconv.Atoi(gId)
	if err == nil && openId != "" {
		result := models.AddNormalUser(groupId, openId)
		if _, ok := result["error"]; ok {
			controller.CustomAbort(403, result["error"])
		}
		controller.Data["json"] = result
	} else {
		switch {
		case err != nil:
			controller.CustomAbort(400, "groupId cannot be null")
		case openId == "":
			controller.CustomAbort(400, "openId cannot be null")
		default:
			controller.CustomAbort(400, "groupId&openId cannot be null")
		}
	}
	controller.ServeJSON()
}

// @Title Get One User
// @Description Get User
// @Param	activityId	path	string	true	"Activity.Id"
// @Param	key			path	string	true	"OpenId/Mobile"
// @Success 200 json{"user","tickets"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /:activityId/:key [get]
func (controller *UserController) GetUser() {
	testkey := controller.Ctx.Input.Param(":key")
	aId := controller.Ctx.Input.Param(":activityId")
	activityId, err := strconv.Atoi(aId)
	if err == nil && testkey != "" {
		result := models.GetUser(activityId, testkey)
		if _, ok := result["error"]; ok {
			controller.CustomAbort(403, result["error"].(string))
		}
		controller.Data["json"] = result
	} else {
		switch {
		case err != nil:
			controller.CustomAbort(203, "ActivityId cannot be null")
		case testkey == "":
			controller.CustomAbort(203, "openId cannot be null")
		default:
			controller.CustomAbort(203, "ActivityId&openId cannot be null")
		}
	}
	controller.ServeJSON()
}

// @Title Update User
// @Description update the User
// @Param	body body 	object	true "<p>{<br>&quot;id&quot;:1,<br>&quot;openId&quot;: &quot;xxxx&quot;,<br>&quot;nickName&quot;: &quot;me&quot;,<br>&quot;headImgUrl&quot;: &quot;yyyy.cn/1.png&quot;,<br>&quot;paid&quot;: false<br>}</p>"
// @Success 200 {"updated": "succeed"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router / [put]
func (controller *UserController) PutUser() {
	test := make(map[string]interface{})
	json.Unmarshal(controller.Ctx.Input.RequestBody, &test)
	err := models.UpdateUser(test)
	if err != nil {
		controller.CustomAbort(403, err.Error())
	} else {
		controller.Data["json"] = map[string]string{"updated": "succeed"}
	}
	controller.ServeJSON()
}

// @Title Delete User
// @Description delete the User
// @Param	body body 	object	true "<p>{<br>&quot;openId&quot;: &quot;xxxx&quot;,<br>&quot;activityId&quot;: 1<br>}</p>"
// @Success 200 {"state": "succeed"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router / [delete]
func (controller *UserController) DeleteUser() {
	req := make(map[string]interface{})
	json.Unmarshal(controller.Ctx.Input.RequestBody, req)
	activityId := int(req["activityId"].(int))
	openId := req["openId"].(string)
	err := models.DeleteUser(activityId, openId)
	if err != nil {
		controller.CustomAbort(403, err.Error())
	} else {
		controller.Data["json"] = map[string]string{"state": "succeed"}
	}
	controller.ServeJSON()
}

// @Title Add Ticket's Owner
// @Description Add Ticket's Owner
// @Param	ticklist		body 	object	true		"{<br>&quot;ticklist&quot;:<br>[1,2,3,4]<br>}"
// @Param	userId		path 	string	true "user's Id"
// @Success 200 {"state":"succeed"}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /:userId/tickets [post]
func (controller *UserController) AddTicketOwner() {
	uId := controller.Ctx.Input.Param(":userId")
	userId, atoiErr := strconv.Atoi(uId)
	ticklistMap := make(map[string][]int)
	err := json.Unmarshal(controller.Ctx.Input.RequestBody, &ticklistMap)
	if err != nil {
		logs.Notice(err.Error())
	}
	if atoiErr == nil && ticklistMap["ticklist"] != nil {
		result := models.AddTicketOwner(userId, ticklistMap["ticklist"])
		if _, ok := result["error"]; ok {
			controller.CustomAbort(403, result["error"].(string))
		}
		controller.Data["json"] = result
	} else {
		controller.CustomAbort(400, "openId or ticklist cannot be empty")
	}
	controller.ServeJSON()
}
