package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// GetAllActivities return all activities
func Dashboard() map[string]interface{} {
	result := make(map[string]interface{})
	var activities []*Activity
	var activitiesSlice []map[string]interface{}
	o := orm.NewOrm()
	num, err := o.QueryTable("activity").All(&activities)
	if err != nil {
		result["error"] = err.Error()
		return result
	}
	for _, activity := range activities {
		temp := make(map[string]interface{})
		groupSlice := []*Group{}
		temp["id"] = activity.Id
		temp["name"] = activity.Name
		temp["startDate"] = activity.StartDate
		temp["endDate"] = activity.EndDate
		temp["existCount"] = activity.ExistCount
		temp["locked"] = activity.Locked
		temp["limitCount"] = activity.LimitCount
		successGroupCount, _ := o.QueryTable("Group").Filter("Activity", activity.Id).Filter("Success", true).All(&groupSlice)
		userCount, _ := o.QueryTable("User").Filter("Activity", activity).Count()
		temp["successGroupCount"] = int(successGroupCount)
		temp["successUserCount"] = int(successGroupCount) * activity.GroupSize
		temp["userCount"] = int(userCount)
		ticket, _ := o.QueryTable("Ticket").Filter("activity_id", activity.Id).Filter("state", "Legal").Count()
		ticketUsed, _ := o.QueryTable("Ticket").Filter("activity_id", activity.Id).Filter("state", "Used").Count()
		temp["ticketCount"] = int(ticket)
		temp["ticketUsedCount"] = int(ticketUsed)
		// Check State
		if !activity.Locked {
			temp["state"] = "未发布"
		} else {
			// Is locked
			now := time.Now().UTC()
			startDate := activity.StartDate
			endDate := activity.EndDate
			if now.Before(startDate) {
				temp["state"] = "未开始"
			} else {
				// In time range
				if now.After(startDate) && now.Before(endDate) {
					var full bool
					existCount, _ := o.QueryTable("group").Filter("activity_id", activity.Id).Filter("success", true).Count()
					if int(existCount) >= activity.LimitCount {
						full = true
					} else {
						full = false
					}
					if full {
						temp["state"] = "已结束"
					} else {
						temp["state"] = "进行中"
					}
				} else {
					temp["state"] = "已结束"
				}
			}
		}
		activitiesSlice = append(activitiesSlice, temp)
	}
	result["count"] = num
	result["activities"] = activitiesSlice
	return result
}
