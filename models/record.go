package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/orm"
)

// Record is the paidment record of users, one user may have more than one record
type Record struct {
	Id         int       `json:"id"`
	OpenId     string    `json:"openId"`
	ActivityId int       `orm:null" json:"activityId"`
	PaidId     string    `orm:"null" json:"paidId"`
	PaidState  string    `orm:"null" json:"paidState"`
	PaidAmount float64   `orm:"null" json:"paidAmount"`
	Created    time.Time `orm:"auto_now_add;type(datetime)"`
	Updated    time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	// Regist models
	orm.RegisterModel(
		new(Record),
	)
}

// AddRecord add one record
func AddRecord(record map[string]interface{}) error {
	o := orm.NewOrm()
	temp := Record{ActivityId: int(record["activityId"].(float64)), OpenId: record["openId"].(string)}
	// break if activity not exist already existed
	if !o.QueryTable("Activity").Filter("Id", temp.ActivityId).Exist() {
		result := errors.New("Activity not existed")
		return result
	}
	// break if user not exist
	if !o.QueryTable("User").Filter("Activity__Id", temp.ActivityId).Filter("OpenId", temp.OpenId).Exist() {
		result := errors.New("User not exitst")
		return result
	}
	for key, value := range record {
		switch key {
		case "paidId":
			temp.PaidId = value.(string)
		case "paidState":
			temp.PaidState = value.(string)
		case "paidAmount":
			temp.PaidAmount = value.(float64)
		}
	}
	_, err := o.Insert(&temp)
	if err != nil {
		return err
	}
	return nil
}

func GetRefundRecords(activityId int) map[string]interface{} {
	result := make(map[string]interface{})
	records := []*Record{}
	o := orm.NewOrm()
	if !o.QueryTable("activity").Filter("id", activityId).Exist() {
		result["error"] = "Activity not exists"
		return result
	}
	sql := o.Raw("SELECT * FROM record WHERE activity_id = ? AND open_id NOT IN (SELECT open_id from `user` WHERE group_id IN (SELECT id from `group` WHERE activity_id = ? AND success = 1))", activityId, activityId)
	num, err := sql.QueryRows(&records)
	result["count"] = num
	if err != nil {
		logs.Notice(err.Error())
		result["error"] = err
		return result
	}
	result["records"] = records
	return result
}
