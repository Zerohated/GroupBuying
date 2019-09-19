package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/orm"
)

// User is the instance of user in certain Activity
type User struct {
	Id         int       `json:"id"`
	OpenId     string    `json:"openId"`
	NickName   string    `json:"nickName"`
	HeadImgUrl string    `json:"headImgUrl"`
	Name       string    `json:"name" orm:"null"`
	Mobile     string    `json:"mobile" orm:"null"`
	Activity   *Activity `orm:"rel(fk);null" json:"activity"`
	Group      *Group    `orm:"rel(fk);null" json:"group"`
	IsStarter  bool      `json:"isStarter"`
	Paid       bool      `orm:"default(false)" json:"paid"`
	PrePaidId  string    `orm:"null" json:"prePaidId"`
	PaidId     string    `orm:"null" json:"paidId"`
	Tickets    []*Ticket `orm:"null;reverse(many)" json:"tickets"`
	Created    time.Time `orm:"auto_now_add;type(datetime)"`
	Updated    time.Time `orm:"auto_now;type(datetime)"`
}

// Group struct defination
type Group struct {
	Id       int       `json:"id"`
	Activity *Activity `orm:"rel(fk);null" json:"activity"`
	Users    []*User   `orm:"null;reverse(many)" json:"users"`
	Size     int       `json:"size"`
	Success  bool      `json:"success"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	// Regist models
	orm.RegisterModel(
		new(User),
		new(Group),
	)
}

// AddStarter Use openId to create a starter for activity refer to activityId
func AddStarter(activityId int, openId string) map[string]string {
	result := make(map[string]string)
	o := orm.NewOrm()
	// judge whether activityId existed
	if !o.QueryTable("activity").Filter("id", activityId).Exist() {
		result["error"] = "activityId not existed"
		return result
	}
	// judge whether openId existed
	if o.QueryTable("User").Filter("Activity__Id", activityId).Filter("open_id", openId).Exist() {
		result["error"] = "openId already existed"
		return result
	}
	// judge whether activity's group is full
	activity := Activity{Id: activityId}
	o.Read(&activity)
	var full bool
	existCount, _ := o.QueryTable("group").Filter("activity_id", activityId).Filter("success", true).Count()
	if int(existCount) >= activity.LimitCount {
		full = true
	} else {
		full = false
	}
	if full {
		// if activity is full
		result["error"] = "Activity is full"
		return result
	}
	// if activity is not full
	// Start Transaction
	o.Begin()
	// update exist count
	_, err := o.QueryTable("activity").Filter("id", activityId).Update(orm.Params{
		"exist_count": orm.ColValue(orm.ColAdd, 1)})
	if err != nil {
		result["error"] = "addGroupError:" + err.Error()
		logs.Notice("insertErr: ", err.Error())
		// Transaction Rollback
		if transactionErr := o.Rollback(); transactionErr != nil {
			logs.Notice("transactionErr: ", transactionErr.Error())
		}
		return result
	}
	result["newActivityExistCount"] = strconv.Itoa(activity.ExistCount + 1)
	// create new group
	group := Group{Activity: &activity}
	group.Size = 1
	num, err := o.Insert(&group)
	if err != nil {
		result["error"] = "addGroupError:" + err.Error()
		logs.Notice(err.Error())
		// Transaction Rollback
		if transactionErr := o.Rollback(); transactionErr != nil {
			logs.Notice("transactionErr: ", transactionErr.Error())
		}
		return result
	}
	result["groupId"] = strconv.FormatInt(num, 10)
	user := User{OpenId: openId, Activity: &activity}
	user.Group = &group
	user.IsStarter = true
	user.Paid = false
	num, err = o.Insert(&user)
	if err != nil {
		result["error"] = "addStarterError:" + err.Error()
		logs.Notice(err.Error())
		// Transaction Rollback
		if transactionErr := o.Rollback(); transactionErr != nil {
			logs.Notice("transactionErr: ", transactionErr.Error())
		}
		return result
	}
	result["userId"] = strconv.FormatInt(num, 10)
	commitErr := o.Commit()
	if commitErr != nil {
		logs.Notice(commitErr.Error())
	}
	return result
}

// AddNormalUser add one normal user to a group refer to groupId
func AddNormalUser(groupId int, openId string) map[string]string {
	result := make(map[string]string)
	o := orm.NewOrm()
	// break if groupId not exist
	if !o.QueryTable("group").Filter("id", groupId).Exist() {
		result["error"] = "Group not exitst"
		return result
	}
	// break if group not exist already existed
	search := o.QueryTable("activity").Filter("Groups__Id", groupId).Limit(1)
	if !search.Exist() {
		result["error"] = "Activity not existed"
		return result
	}
	activity := Activity{}
	err := search.One(&activity)
	if err != nil {
		logs.Notice(err.Error())
	}
	// break if openId already existed
	if o.QueryTable("User").Filter("Activity__Id", activity.Id).Filter("open_id", openId).Exist() {
		result["error"] = "OpenId already existed"
		return result
	}
	group := Group{Id: groupId}
	o.Read(&group)
	if group.Size >= activity.GroupSize {
		result["error"] = "Group is full"
		return result
	}
	// proceed if groupId exists, activity exists, openId not exists, group is not full

	// Set groupId
	result["groupId"] = strconv.Itoa(groupId)
	o.Begin()
	user := User{OpenId: openId, Group: &group}
	user.Activity = group.Activity
	user.IsStarter = false
	user.Paid = false
	group.Size++
	num, err := o.Insert(&user)
	if err != nil {
		result["error"] = "addUserError:" + err.Error()
		logs.Notice("insertErr: ", err.Error())
		if transactionErr := o.Rollback(); transactionErr != nil {
			logs.Notice("transactionErr: ", transactionErr.Error())
		}
		return result
	}
	_, err = o.Update(&group)
	if err != nil {
		result["error"] = "UpdateUserError:" + err.Error()
		logs.Notice("UpdateUserError: ", err.Error())
		if transactionErr := o.Rollback(); transactionErr != nil {
			logs.Notice("transactionErr: ", transactionErr.Error())
		}
		return result
	}
	commitErr := o.Commit()
	if commitErr != nil {
		logs.Notice(commitErr.Error())
	}
	result["userId"] = strconv.FormatInt(num, 10)
	return result
}

// GetUser use activityId and key to get User's info
func GetUser(id int, key string) map[string]interface{} {
	result := make(map[string]interface{})
	o := orm.NewOrm()
	if o.QueryTable("activity").Filter("id", id).Exist() {
		userMap := []orm.Params{}
		switch {
		case o.QueryTable("User").Filter("Activity__Id", id).Filter("mobile", key).Exist():
			_, err := o.QueryTable("User").Filter("Activity__Id", id).Filter("mobile", key).Values(&userMap, "Id", "OpenId", "NickName", "HeadImgUrl", "Name", "Mobile", "IsStarter", "Paid", "PaidId", "Group")
			if err != nil {
				logs.Notice(err.Error())
			}
		case o.QueryTable("User").Filter("Activity__Id", id).Filter("open_id", key).Exist():
			_, err := o.QueryTable("User").Filter("Activity__Id", id).Filter("open_id", key).Values(&userMap, "Id", "OpenId", "NickName", "HeadImgUrl", "Name", "Mobile", "IsStarter", "Paid", "PaidId", "Group")
			if err != nil {
				logs.Notice(err.Error())
			}
		default:
			result["error"] = "User not existed"
			return result
		}
		tempUser := make(map[string]interface{})
		for k, v := range userMap[0] {
			switch {
			case k == "Group__Group":
				tempUser["group"] = map[string]interface{}{"id": v}
			default:
				str := strings.ToLower(string(k[0]))
				// first := strings.ToLower(k[0])
				other := k[1:len(k)]
				tempUser[fmt.Sprintf("%s%s", str, other)] = v
			}
		}
		group := Group{Id: int(tempUser["group"].(map[string]interface{})["id"].(int64))}
		o.Read(&group)
		tempUser["success"] = group.Success
		result["user"] = tempUser

		ticketMap := []orm.Params{}
		_, ticketErr := o.QueryTable("ticket").Filter("owner_id", tempUser["id"]).OrderBy("-updated").Values(&ticketMap, "Code", "State", "TicketModel__ID", "TicketModel__Picture", "TicketModel__TicketName")
		if ticketErr != nil {
			result["error"] = ticketErr.Error()
			return result
		}
		ticketSlice := []map[string]interface{}{}
		for _, item := range ticketMap {
			temp := make(map[string]interface{})
			for k, v := range item {
				switch {
				case k == "TicketModel__Picture":
					temp["picture"] = v
				case k == "TicketModel__Id":
					temp["ticketModelId"] = v
				case k == "TicketModel__TicketName":
					temp["ticketModelName"] = v
				default:
					str := strings.ToLower(string(k[0]))
					// first := strings.ToLower(k[0])
					other := k[1:len(k)]
					temp[fmt.Sprintf("%s%s", str, other)] = v
				}
			}
			ticketSlice = append(ticketSlice, temp)
		}

		result["tickets"] = ticketSlice
	} else {
		result["error"] = "activityId not existed"
	}
	return result
}

func GetAllUsers(activityId int) map[string]interface{} {
	result := make(map[string]interface{})
	// var users []*User
	users := []orm.Params{}
	o := orm.NewOrm()
	if !o.QueryTable("activity").Filter("id", activityId).Exist() {
		result["error"] = "Activity not exists"
		return result
	}
	num, err := o.QueryTable("User").Filter("Activity__Id", activityId).Limit(-1).Values(&users, "Id", "OpenId", "NickName", "Name", "Mobile", "HeadImgUrl", "IsStarter", "Paid", "PaidId")
	result["count"] = num
	if err != nil {
		logs.Notice(err.Error())
		result["error"] = err
		return result
	}
	usersSlice := []map[string]interface{}{}
	for _, item := range users {
		temp := make(map[string]interface{})
		for k, v := range item {
			str := strings.ToLower(string(k[0]))
			// first := strings.ToLower(k[0])
			other := k[1:len(k)]
			temp[fmt.Sprintf("%s%s", str, other)] = v
		}
		usersSlice = append(usersSlice, temp)
	}
	result["users"] = usersSlice
	return result
}

func GetGroupInfo(activityId int, groupId int) map[string]interface{} {
	result := make(map[string]interface{})
	o := orm.NewOrm()
	if o.QueryTable("activity").Filter("id", activityId).Exist() {
		count, _ := o.QueryTable("group").Filter("Activity__Id", activityId).Filter("id", groupId).Count()
		if int(count) == 0 {
			result["error"] = "No such groupId in this activity"
			return result
		}
		temp := Group{}
		groupInfo := make(map[string]interface{})
		err := o.QueryTable("group").Filter("Activity__Id", activityId).Filter("id", groupId).Limit(1).One(&temp)
		if err != nil {
			result["error"] = err.Error()
		} else {
			users := []orm.Params{}
			count, usersErr := o.QueryTable("User").Filter("Activity__Id", activityId).Filter("Group__Id", groupId).Values(&users, "Id", "OpenId", "NickName", "HeadImgUrl", "Name", "Mobile", "IsStarter", "Paid", "PaidId", "Created")
			if usersErr != nil {
				result["error"] = usersErr.Error()
				return result
			}
			usersSlice := []map[string]interface{}{}
			for _, item := range users {
				temp := make(map[string]interface{})
				for k, v := range item {
					str := strings.ToLower(string(k[0]))
					// first := strings.ToLower(k[0])
					other := k[1:len(k)]
					temp[fmt.Sprintf("%s%s", str, other)] = v
				}
				usersSlice = append(usersSlice, temp)
			}
			groupInfo["id"] = temp.Id
			groupInfo["activity"] = temp.Activity.Id
			groupInfo["users"] = usersSlice
			groupInfo["size"] = temp.Size
			result["group"] = groupInfo
			// Group size
			activity := Activity{Id: activityId}
			o.Read(&activity)
			if int(count) >= activity.GroupSize {
				groupInfo["groupSucceed"] = "True"
			} else {
				groupInfo["groupSucceed"] = "False"
			}
		}

	} else {
		result["error"] = "activityId not existed"
	}
	return result
}

func UpdateUser(user map[string]interface{}) error {
	o := orm.NewOrm()
	temp := User{Id: int(user["id"].(float64))}
	readErr := o.Read(&temp)
	o.LoadRelated(&temp, "Group")
	o.LoadRelated(&temp, "Activity")
	if readErr != nil {
		logs.Notice("User with this id not found:", temp.Id)
		return readErr
	}
	for key, value := range user {
		switch key {
		case "openId":
			temp.OpenId = value.(string)
		case "nickName":
			temp.NickName = value.(string)
		case "headImgUrl":
			temp.HeadImgUrl = value.(string)
		case "name":
			temp.Name = value.(string)
		case "mobile":
			temp.Mobile = value.(string)
		case "prePaidId":
			temp.PrePaidId = value.(string)
		case "paidId":
			temp.PaidId = value.(string)
		case "paid":
			temp.Paid = value.(bool)
		}
	}
	if temp.Paid == true {
		if temp.Group.Size >= temp.Activity.GroupSize {
			temp.Group.Success = true
		}
		o.Update(temp.Group)
	}
	if _, err := o.Update(&temp); err != nil {
		return err
	}
	return nil
}

// TODO: AddTicketOwner needs rewrite
func AddTicketOwner(userId int, ticklist []int) map[string]interface{} {
	result := make(map[string]interface{})
	o := orm.NewOrm()
	if o.QueryTable("user").Filter("id", userId).Exist() {
		transactionErr := o.Begin()
		ticket := Ticket{}
		failedCount := 0
		succeedCount := 0
		queryErr := ""
		for _, i := range ticklist {
			err := o.QueryTable("ticket").Filter("ticket_model_id", i).Filter("state", "Illegal").One(&ticket)
			if err != nil {
				logs.Debug("No ticket avaliable: ticketModel - ", i)
				queryErr += fmt.Sprintf("No ticket avaliable: ticketModel - %d;", i)
				failedCount++
			} else {
				activityUser := User{Id: userId}
				o.Read(&activityUser, "open_id")
				ticket.State = "Legal"
				ticket.Owner = &activityUser
				o.Update(&ticket, "owner", "updated", "state")
				succeedCount++
			}
		}
		if failedCount != 0 {
			transactionErr = o.Rollback()
			result["error"] = queryErr
			return result
		}
		transactionErr = o.Commit()
		result["state"] = "succeed"
		if transactionErr != nil {
			logs.Debug(transactionErr)
		}
	}
	return result
}

func DeleteUser(activityId int, openId string) error {
	o := orm.NewOrm()
	queryString := o.QueryTable("User").Filter("Activity", activityId).Filter("OpenId", openId)
	if queryString.Exist() {
		user := new(User)
		err := queryString.One(&user)
		if err != nil {
			return err
		}
		_, err = o.Delete(&user)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Activity not exist")
}
