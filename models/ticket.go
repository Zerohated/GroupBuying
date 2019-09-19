package models

import (
	"fmt"
	"group_buying/constant"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// TicketModel define the ticket detail
type TicketModel struct {
	// TicketID	Int	Y
	Id             int         `json:"id"`
	Picture        string      `orm:"null" json:"picture"`
	TicketName     string      `json:"ticketName"`
	TicketType     int         `json:"ticketType"`
	TicketTypeName string      `json:"ticketTypeName"`
	DiscountMoney  float64     `json:"discountMoney"`
	DiscountRate   float64     `json:"discountRate"`
	WhenFull       float64     `json:"whenFull"`
	BuyLimitNum    int         `json:"buyLimitNum"`
	Pukaamount     int         `json:"pukaamount"`
	Pukabonus      int         `json:"pukabonus"`
	VIPamount      int         `json:"vIPamount"`
	VIPbonus       int         `json:"vIPbonus"`
	VVIPamount     int         `json:"vVIPamount"`
	VVIPbonus      int         `json:"vVIPbonus"`
	TicketDesc     string      `json:"ticketDesc"`
	DiscountDesc   string      `json:"discountDesc"`
	Remark         string      `json:"remark"`
	Channel        string      `json:"channel"`
	Mac            string      `json:"mac"`
	Activitys      []*Activity `orm:"reverse(many);rel_through(group_buying/models.ActivityRelTicketModel)" json:"activitys"`
	Tickets        []*Ticket   `orm:"reverse(many);null" json:"tickets"`
	Created        time.Time   `orm:"auto_now_add;type(datetime)" json:"created"`
	Updated        time.Time   `orm:"auto_now;type(datetime)" json:"updated"`
}

// Ticket struct defination
// Ticket instance refer to TicketModel
type Ticket struct {
	Code        string       `orm:"pk" json:"code"`
	State       string       `json:"state"`
	TicketModel *TicketModel `orm:"rel(fk);null" json:"ticketModel"`
	Activity    *Activity    `orm:"rel(fk);null" json:"activity"`
	Owner       *User        `orm:"rel(fk);null" json:"owner"`
	Created     time.Time    `orm:"auto_now_add;type(datetime)" json:"created"`
	Updated     time.Time    `orm:"auto_now;type(datetime)" json:"updated"`
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(
		new(TicketModel),
		new(Ticket),
	)
}

// AddOneTicketModel get ticketModel return id and error
func AddOneTicketModel(ticketModel TicketModel) (string, error) {
	o := orm.NewOrm()
	i, err := o.InsertOrUpdate(&ticketModel)
	if err != nil {
		return strconv.FormatInt(i, 10), err
	}
	return strconv.FormatInt(i, 10), nil
}

// GetOneTicketModel return and TicketModel
func GetOneTicketModel(ticketModelId string) *TicketModel {
	id, _ := strconv.Atoi(ticketModelId)
	ticketModel := TicketModel{Id: id}
	o := orm.NewOrm()
	o.Read(&ticketModel)
	return &ticketModel
}

// GetAllTicketModels return all TicketModels
func GetAllTicketModels() map[string]interface{} {
	result := make(map[string]interface{})
	var ticketModels []*TicketModel
	o := orm.NewOrm()
	num, err := o.QueryTable("ticket_model").All(&ticketModels)
	result["count"] = num
	result["error"] = err
	result["ticketModels"] = ticketModels
	return result
}

// UpdateTicketModel function
func UpdateTicketModel(ticketModel TicketModel) interface{} {
	o := orm.NewOrm()
	if o.QueryTable("TicketModel").Filter("Id", ticketModel.Id).Exist() {
		t := TicketModel{Id: ticketModel.Id}
		o.Read(&t)
		ticketModel.Created = t.Created
		_, err := o.Update(&ticketModel)
		if err != nil {
			return err
		}
	} else {
		_, err := o.Insert(&ticketModel)
		if err != nil {
			return err
		}
	}
	logs.Notice(ticketModel)
	return ticketModel
}

// GenerateTickets create tickets
func GenerateTickets(activityId string, ticketModelId string, count string) map[string]interface{} {
	result := make(map[string]interface{})
	o := orm.NewOrm()
	num, _ := strconv.Atoi(count)
	tId, _ := strconv.Atoi(ticketModelId)
	aId, _ := strconv.Atoi(activityId)
	ticketModel := TicketModel{Id: tId}
	activity := Activity{Id: aId}
	if o.QueryTable("ticket_model").Filter("id", tId).Exist() {
		o.Read(&ticketModel)
		// Transaction Begin
		transactionErr := o.Begin()
		// PrepareInsert
		pi, _ := o.QueryTable("ticket").PrepareInsert()
		insertCount := 0
		for i := 1; i <= num; i++ {
			tempStr := fmt.Sprintf("%screate%sticket%d", ticketModelId, count, i)
			temp := Ticket{Code: constant.ConvertMD5(tempStr)}
			temp.State = "Illegal"
			temp.TicketModel = &ticketModel
			temp.Activity = &activity
			_, err := pi.Insert(&temp)
			if err != nil {
				logs.Notice(err.Error())
			} else {
				insertCount++
			}
		}
		pi.Close()
		// Close insert statement
		if insertCount != num {
			transactionErr = o.Rollback()
		} else {
			transactionErr = o.Commit()
		}
		if transactionErr != nil {
			logs.Debug("TransactionErr: %s", transactionErr)
		}
		// Transaction End
		result["count"] = insertCount
	} else {
		result["error"] = "TicketModel not exist"
	}
	return result
}

func CheckTicket(ticketCode string) map[string]interface{} {
	result := make(map[string]interface{})
	o := orm.NewOrm()
	ticket := Ticket{Code: ticketCode}
	if !o.QueryTable("ticket").Filter("code", ticketCode).Exist() {
		result["error"] = "ticket not exists"
		return result
	}
	o.QueryTable("ticket").Filter("code", ticketCode).One(&ticket)
	o.LoadRelated(&ticket, "TicketModel")
	result["ticketModel"] = ticket.TicketModel
	result["state"] = ticket.State
	return result
}

func BurnTicket(ticketCode string) map[string]interface{} {
	result := make(map[string]interface{})
	o := orm.NewOrm()
	ticket := Ticket{Code: ticketCode}
	if !o.QueryTable("ticket").Filter("code", ticketCode).Exist() {
		result["error"] = "ticket not exists"
		return result
	}
	o.QueryTable("ticket").Filter("code", ticketCode).One(&ticket)
	o.LoadRelated(&ticket, "TicketModel")
	// Ticket validation
	if ticket.State == "Legal" {
		ticket.State = "Used"
		o.Update(&ticket)
		result["state"] = "succeed"
	} else {
		result["error"] = "Ticket Illegal"
	}
	return result
}

// GetUsedTickets return all used tickets
func GetUsedTickets(activityId int) map[string]interface{} {
	result := make(map[string]interface{})
	tickets := []orm.Params{}
	o := orm.NewOrm()
	var num int64
	var err error
	if activityId != 0 {
		num, err = o.QueryTable("ticket").Filter("Activity", activityId).Filter("state", "Used").Limit(-1).Values(&tickets, "Code", "State", "TicketModel__TicketName", "TicketModel__TicketTypeName", "Activity__Name")
	} else {
		num, err = o.QueryTable("ticket").Filter("state", "Used").Limit(-1).Values(&tickets, "Code", "State", "TicketModel__TicketName", "TicketModel__TicketTypeName", "Activity__Name")
	}
	ticketSlice := []map[string]interface{}{}
	for _, item := range tickets {
		temp := make(map[string]interface{})
		for k, v := range item {
			str := strings.ToLower(string(k[0]))
			// first := strings.ToLower(k[0])
			other := k[1:len(k)]
			t := fmt.Sprintf("%s%s", str, other)
			temp[strings.Replace(t, "__", "", -1)] = v
		}
		ticketSlice = append(ticketSlice, temp)
	}
	result["count"] = num
	result["error"] = err
	result["tickets"] = ticketSlice
	return result
}
