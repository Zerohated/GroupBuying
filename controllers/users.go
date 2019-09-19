package controllers

import (
	"group_buying/models"
	"strconv"

	"github.com/astaxie/beego"
)

// UsersController about Users
type UsersController struct {
	beego.Controller
}

// @Title GetAll
// @Description get all User
// @Param	activityId	path	int	true		"the activityId"
// @Success 200 {<br>"users": [User1,User2,...],<br>"count":2}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /:activityId/ [get]
func (controller *UsersController) GetAllUsers() {
	aId := controller.Ctx.Input.Param(":activityId")
	if aId != "" {
		activityId, err := strconv.Atoi(aId)
		if err != nil {
			controller.CustomAbort(400, "ActivityId must be int")
		}
		result := models.GetAllUsers(activityId)
		if _, ok := result["error"]; ok {
			controller.CustomAbort(403, result["error"].(string))
		}
		controller.Data["json"] = result
	} else {
		controller.CustomAbort(400, "ActivityId cannot be null")
	}
	controller.ServeJSON()
}

// @Title Get Group
// @Description Get Group
// @Param	activityId		path 	string	true		"the activityId"
// @Param	groupId		path 	string	true	"groupId"
// @Success 200 {<br>"group": Group1,<br>"error":null<br>}
// @Failure 403 "Error infomation"
// @Failure 400 "Params Error infomation"
// @router /:activityId/group/:groupId [get]
func (controller *UsersController) GetGroupInfo() {
	activityId := controller.Ctx.Input.Param(":activityId")
	groupId := controller.Ctx.Input.Param(":groupId")
	if activityId != "" && groupId != "" {
		aId, err := strconv.Atoi(activityId)
		if err != nil {
			controller.CustomAbort(400, err.Error())
		}
		gId, err := strconv.Atoi(groupId)
		if err != nil {
			controller.CustomAbort(400, err.Error())
		}
		result := models.GetGroupInfo(aId, gId)
		if _, ok := result["error"]; ok {
			controller.CustomAbort(403, result["error"].(string))
		}
		controller.Data["json"] = result
	} else {
		switch {
		case activityId == "":
			controller.CustomAbort(203, "ActivityId cannot be null")
		case groupId == "":
			controller.CustomAbort(203, "groupId cannot be null")
		default:
			controller.CustomAbort(203, "ActivityId&groupId cannot be null")
		}
	}
	controller.ServeJSON()
}
