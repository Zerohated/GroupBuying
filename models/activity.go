package models

import (
	"errors"
	"group_buying/constant"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// Activity define the struct of a activity
type Activity struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Price        float64        `json:"price"`
	GroupSize    int            `json:"groupSize"`
	LimitCount   int            `json:"limitCount"`
	StartDate    time.Time      `orm:"null" json:"startDate"`
	EndDate      time.Time      `orm:"null" json:"endDate"`
	ExistCount   int            `json:"existCount"`
	PrizeNumber  int            `json:"prizeNumber"`
	Locked       bool           `json:"locked"`
	TicketModels []*TicketModel `orm:"null;rel(m2m);rel_through(group_buying/models.ActivityRelTicketModel)" json:"ticketModels"`
	Tickets      []*Ticket      `orm:"null;reverse(many)" json:"tickets"`
	Users        []*User        `orm:"null;reverse(many)" json:"users"`
	Groups       []*Group       `orm:"null;reverse(many)" json:"groups"`
	ActivityUi   *ActivityUi    `orm:"null;rel(one)"`
	Created      time.Time      `orm:"auto_now_add;type(datetime)" json:"-"`
	Updated      time.Time      `orm:"auto_now;type(datetime)" json:"-"`
}

type ActivityRelTicketModel struct {
	Id          int          `json:"id"`
	Activity    *Activity    `orm:"rel(fk);column(activity_id)" json:"activity"`
	TicketModel *TicketModel `orm:"rel(fk);column(ticket_model_id)" json:"ticketModel"`
	IsAmust     bool         `orm:"null"json:"isAmust"`
	UseDetail   string       `orm:"null;type(text)" json:"useDetail"`
	StartDate   time.Time    `orm:"null;type(datetime)" json:"startDate"`
	EndDate     time.Time    `orm:"null;type(datetime)" json:"endDate"`
}

// ActivityUi information, Id is related to Activity's Id
type ActivityUi struct {
	Id            int       `json:"id"`
	Activity      *Activity `orm:"reverse(one)" json:"activity,omitempty"`
	Background    string    `orm:"type(text)" json:"background"`
	Button        string    `orm:"type(text)" json:"button"`
	DetailButton  string    `orm:"type(text)" json:"detailButton"`
	SuccessButton string    `orm:"type(text)" json:"successButton"`
	Description   string    `orm:"type(text)" json:"description"`
	Detail        string    `orm:"type(text)" json:"detail"`
	KeyVisual     string    `orm:"type(text)" json:"keyVisual"`
	EndNotice     string    `orm:"type(text)" json:"endNotice"`
	NotSuccess    string    `orm:"type(text)" json:"notSuccess"`
	SuccessTop    string    `orm:"type(text)" json:"successTop"`
	Created       time.Time `orm:"auto_now_add;type(datetime)" json:"-"`
	Updated       time.Time `orm:"auto_now;type(datetime)" json:"-"`
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(
		new(Activity),
		new(ActivityUi),
		new(ActivityRelTicketModel),
	)
}

// AddActivity returnactivityId and error
func AddActivity(activity map[string]interface{}) (int, error) {
	o := orm.NewOrm()
	temp := Activity{}
	for key, value := range activity {
		switch key {
		case "name":
			temp.Name = value.(string)
		case "price":
			temp.Price = value.(float64)
		case "groupSize":
			temp.GroupSize = int(value.(float64))
		case "limitCount":
			temp.LimitCount = int(value.(float64))
		case "startDate":
			sDate, _ := time.ParseInLocation(constant.Format, value.(string), time.UTC)
			temp.StartDate = sDate
		case "endDate":
			eDate, _ := time.ParseInLocation(constant.Format, value.(string), time.UTC)
			temp.EndDate = eDate
		}
	}
	temp.ExistCount = 0
	temp.Locked = false
	i, err := o.Insert(&temp)
	return int(i), err
}

// AddActivityUi returnactivityId and error
func AddActivityUi(activityUI ActivityUi) (int, error) {
	o := orm.NewOrm()
	if o.QueryTable("Activity").Filter("Id", activityUI.Id).Exist() {
		activity := Activity{}
		o.QueryTable("Activity").Filter("Id", activityUI.Id).One(&activity)
		activityUI.Activity = &activity
		activity.ActivityUi = &activityUI
		o.Update(&activity)
		i, err := o.Insert(&activityUI)
		return int(i), err
	}
	return 0, errors.New("This activity not exist")
}

// GetActivity return an Activity information
func GetActivity(activityId int) map[string]interface{} {
	result := make(map[string]interface{})
	activity := Activity{Id: activityId}
	o := orm.NewOrm()
	if !o.QueryTable("activity").Filter("id", activityId).Exist() {
		result["error"] = "Activity not exist"
		return result
	}
	o.QueryTable("activity").Filter("id", activityId).One(&activity)
	if activity.ActivityUi != nil {
		o.Read(activity.ActivityUi)
		defer func() {
			result["activityUI"] = activity.ActivityUi
		}()
	}
	defer func() {
		result["name"] = activity.Name
		result["price"] = activity.Price
		result["groupSize"] = activity.GroupSize
		result["limitCount"] = activity.LimitCount
		result["existCount"] = activity.ExistCount
		result["startDate"] = activity.StartDate
		result["endDate"] = activity.EndDate
		result["locked"] = activity.Locked
		result["prizeNumber"] = activity.PrizeNumber
	}()
	// IsEnd Validation Start
	isEnd := false
	full := false
	expired := false
	// -isFull Validation Start

	existCount, _ := o.QueryTable("group").Filter("activity_id", activityId).Filter("success", true).Count()
	if int(existCount) >= activity.LimitCount {
		full = true
	} else {
		full = false
	}
	// -isFull Validation End
	// -inTimeRange Validation Start
	now := time.Now().UTC()
	startDate := activity.StartDate
	endDate := activity.EndDate
	if now.After(startDate) && now.Before(endDate) {
		expired = false
	} else {
		expired = true
	}
	// -inTimeRange Validation End
	if full || expired {
		isEnd = true
	} else {
		isEnd = false
	}
	defer func() {
		result["isEnd"] = isEnd
	}()
	// IsEnd Validation End
	return result
}

func GetActivityDetail(activityId int) map[string]interface{} {
	result := make(map[string]interface{})
	activity := Activity{Id: activityId}
	o := orm.NewOrm()
	o.QueryTable("activity").Filter("id", activityId).One(&activity)
	defer func() {
		result["name"] = activity.Name
		result["price"] = activity.Price
		result["groupSize"] = activity.GroupSize
		result["locked"] = activity.Locked
		result["prizeNumber"] = activity.PrizeNumber
	}()
	// tickets := []*Ticket{}
	// o.LoadRelated(&activity, "Tickets")
	// for _, ticket := range activity.Tickets {
	// 	o.LoadRelated(ticket, "Owner")
	// 	o.LoadRelated(ticket, "TicketModel")
	// 	tickets = append(tickets, ticket)
	// }
	// result["tickets"] = tickets
	// Load TicketModels
	activityRelTicketModel := []*ActivityRelTicketModel{}
	_, err := o.QueryTable("ActivityRelTicketModel").Filter("activity_id", activityId).All(&activityRelTicketModel)
	if err != nil {
		logs.Notice(err)
	}
	temp := []map[string]interface{}{}
	for _, item := range activityRelTicketModel {
		t := make(map[string]interface{})
		o.LoadRelated(item, "ticket_model_id")
		t["ticketModel"] = item.TicketModel
		t["isAmust"] = item.IsAmust
		t["useDetail"] = item.UseDetail
		t["startDate"] = item.StartDate
		t["endDate"] = item.EndDate
		temp = append(temp, t)
	}
	result["activityRelTicketModel"] = temp
	// Load Groups
	groups := []*Group{}
	o.LoadRelated(&activity, "Groups")
	for _, group := range activity.Groups {
		o.LoadRelated(group, "Users")
		groups = append(groups, group)
	}
	result["groups"] = groups
	return result
}

// GetAllActivities return all activities
func GetAllActivities() map[string]interface{} {
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
		o.LoadRelated(activity, "ActivityUi")
		temp := make(map[string]interface{})
		temp["id"] = activity.Id
		temp["name"] = activity.Name
		temp["price"] = activity.Price
		temp["groupSize"] = activity.GroupSize
		temp["limitCount"] = activity.LimitCount
		temp["startDate"] = activity.StartDate
		temp["endDate"] = activity.EndDate
		temp["prizeNumber"] = activity.PrizeNumber
		temp["existCount"] = activity.ExistCount
		temp["locked"] = activity.Locked
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
					temp["state"] = "进行中"
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

// UpdateActivity function
func UpdateActivity(activity map[string]interface{}) error {
	o := orm.NewOrm()
	temp := Activity{Id: activity["id"].(int)}
	readErr := o.Read(&temp)
	if readErr == nil {
		if temp.Locked {
			return errors.New("Activity is locked")
		}
		for key, value := range activity {
			switch key {
			case "name":
				temp.Name = value.(string)
			case "price":
				temp.Price = value.(float64)
			case "groupSize":
				temp.GroupSize = int(value.(float64))
			case "limitCount":
				temp.LimitCount = int(value.(float64))
			case "prizeNumber":
				temp.PrizeNumber = int(value.(float64))
			case "locked":
				temp.Locked = value.(bool)
			case "startDate":
				sDate, _ := time.ParseInLocation(constant.Format, value.(string), time.UTC)
				temp.StartDate = sDate
			case "endDate":
				eDate, _ := time.ParseInLocation(constant.Format, value.(string), time.UTC)
				temp.EndDate = eDate
			}
		}
		_, err := o.Update(&temp)
		return err
	} else {
		logs.Notice(readErr.Error())
		return errors.New("ActivityId not existed")
	}
}

// UpdateActivityUi function
func UpdateActivityUi(activityUi map[string]interface{}) error {
	o := orm.NewOrm()
	temp := Activity{Id: activityUi["id"].(int)}
	au := ActivityUi{Id: activityUi["id"].(int)}
	readErr := o.Read(&temp)
	if readErr == nil {
		if err := o.Read(&au); err == nil {
			for key, value := range activityUi {
				switch key {
				case "background":
					au.Background = value.(string)
				case "button":
					au.Button = value.(string)
				case "detailButton":
					au.DetailButton = value.(string)
				case "successButton":
					au.SuccessButton = value.(string)
				case "description":
					au.Description = value.(string)
				case "detail":
					au.Detail = value.(string)
				case "keyVisual":
					au.KeyVisual = value.(string)
				case "endNotice":
					au.EndNotice = value.(string)
				case "notSuccess":
					au.NotSuccess = value.(string)
				case "successTop":
					au.SuccessTop = value.(string)
				}
			}
			_, err := o.Update(&au)
			return err
		}
		return errors.New("This activity not Existed")
	}
	return errors.New("Activity with this Id not existed")
}

func DeleteActivity(activityId int) map[string]string {
	o := orm.NewOrm()
	if o.QueryTable("Activity").Filter("Id", activityId).Exist() {
		activity := Activity{Id: activityId}
		_, err := o.Delete(&activity)
		if err != nil {
			return map[string]string{"error": err.Error()}
		}
		return map[string]string{"state": "succeed"}
	}
	return map[string]string{"error": "Activity not exist"}
}

// AddOrUpdateActivityTicketModels Relation between Activity and TicketModels
func AddOrUpdateActivityTicketModels(activityId int, ticketModelId int, info map[string]interface{}) map[string]string {
	result := make(map[string]string)
	o := orm.NewOrm()
	if !o.QueryTable("activity").Filter("id", activityId).Exist() {
		result["error"] = "Activity not exist"
		return result
	}
	if !o.QueryTable("ticket_model").Filter("id", ticketModelId).Exist() {
		result["error"] = "TicketModel not exist"
		return result
	}
	activity := Activity{Id: activityId}
	ticketModel := TicketModel{Id: ticketModelId}
	o.Read(&activity)
	o.Read(&ticketModel)
	activityRelTicketModel := ActivityRelTicketModel{Activity: &activity, TicketModel: &ticketModel}
	if created, id, err := o.ReadOrCreate(&activityRelTicketModel, "Activity", "TicketModel"); err == nil {
		logs.Notice("created: ", created)
		logs.Notice("id: ", id)
		activityRelTicketModel.Id = int(id)
		for key, value := range info {
			switch key {
			case "isAmust":
				activityRelTicketModel.IsAmust = value.(bool)
			case "useDetail":
				activityRelTicketModel.UseDetail = value.(string)
			case "startDate":
				sDate, _ := time.Parse(constant.Format, value.(string))
				activityRelTicketModel.StartDate = sDate
			case "endDate":
				eDate, _ := time.Parse(constant.Format, value.(string))
				activityRelTicketModel.EndDate = eDate
			}
			_, insertErr := o.Update(&activityRelTicketModel)
			if insertErr != nil {
				result["error"] = err.Error()
				return result
			}
			result["state"] = "succeed"
		}

	}
	return result
}

func DeleteActivityTicketModels(activityId int, ticketModelId int) map[string]string {
	result := make(map[string]string)
	o := orm.NewOrm()
	if !o.QueryTable("activity").Filter("id", activityId).Exist() {
		result["error"] = "Activity not exist"
		return result
	}
	if !o.QueryTable("ticket_model").Filter("id", ticketModelId).Exist() {
		result["error"] = "TicketModel not exist"
		return result
	}
	activity := Activity{Id: activityId}
	ticketModel := TicketModel{Id: ticketModelId}
	o.Read(&activity)
	o.Read(&ticketModel)
	if !o.QueryM2M(&activity, "TicketModels").Exist(&TicketModel{Id: ticketModelId}) {
		result["error"] = "this relation not exist"
		return result
	}
	activityRelTicketModel := ActivityRelTicketModel{}
	err := o.QueryTable("ActivityRelTicketModel").Filter("activity_id", activityId).Filter("ticket_model_id", ticketModelId).One(&activityRelTicketModel)
	if err != nil {
		result["error"] = err.Error()
		return result
	}
	_, deleteErr := o.Delete(&activityRelTicketModel)
	if deleteErr != nil {
		result["error"] = err.Error()
		return result
	}
	result["state"] = "succeed"
	return result

}

// _,err := o.QueryTable("ActivityRelTicketModel").Filter("activity_id", activityId).Filter("ticket_model_id",ticketModelId)
// One(&activityRelTicketModel)
// res, err := o.Raw("INSERT INTO `activity_rel_ticket_model` (`activity_id`, `ticket_model_id`, `is_amust`, `use_detail`,`start_date`,`end_date`) VALUES (?, ?, ?, ?, ?, ?)",
// 	1, 1, true, "", time.Now(), time.Now()).Exec()
// if err == nil {
// 	num, _ := res.RowsAffected()
// 	fmt.Println("mysql row affected nums: ", num)
// }
// if addErr == nil {
// 	result["inserted"] = strconv.FormatInt(num, 10)
// } else {
// 	result["error"] = "insertErr:" + addErr.Error()
// }

// nums, removeErr := m2m.Clear()
// if removeErr == nil {
// 	result["removed"] = strconv.FormatInt(nums, 10)
// } else {
// 	result["error"] = "removeErr:" + removeErr.Error()
// }
