package models

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/astaxie/beego/logs"

	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/orm"
)

var (
	StarterOpenId []string
	UserOpenId    []string
)

func reset() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:Start123@/dev?charset=utf8&loc=Asia%2FShanghai")
	orm.DefaultTimeLoc = time.UTC
	// Resync Database
	// 数据库别名
	name := "default"
	// drop table 后再建表
	force := false
	// 打印执行过程
	verbose := false
	// 遇到错误立即返回
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	reset()
	o := orm.NewOrm()
	activity := Activity{Id: 111, LimitCount: 1, GroupSize: 2}
	o.Insert(&activity)
	activity = Activity{Id: 222, LimitCount: 1, GroupSize: 2}
	o.Insert(&activity)
	for i := 0; i < 2; i++ {
		StarterOpenId = append(StarterOpenId, fmt.Sprintf("starter%d", i))
	}
	for i := 0; i < 5; i++ {
		UserOpenId = append(UserOpenId, fmt.Sprintf("user%d", i))
	}
}
func Test_AddStarter(t *testing.T) {
	for _, openId := range StarterOpenId {
		result := AddStarter(111, openId)
		t.Log("Add $", openId, "$: ", result)
	}
}
func Test_AddNormalUser(t *testing.T) {
	for _, openId := range UserOpenId {
		result := AddNormalUser(1, openId)
		t.Log("Add $", openId, "$: ", result)
	}
}

func Test_GetUser(t *testing.T) {
	for _, openId := range StarterOpenId {
		result := GetUser(111, openId)
		t.Log("Result $", openId, "$: \n", result)
	}
	for _, openId := range UserOpenId {
		result := GetUser(111, openId)
		t.Log("Result $", openId, "$: \n", result)
	}
}

func Test_GetAllUsers(t *testing.T) {
	result := GetAllUsers(111)
	for key, item := range result {
		t.Log("Key: ", key)
		if key == "users" {
			for _, j := range item.([]map[string]interface{}) {
				t.Log(j)
			}
		} else {
			t.Log("item: ", item)
		}
	}
}

func Test_GetGroupInfo(t *testing.T) {
	groupId := []int{1, 2, 3}
	for _, id := range groupId {
		result := GetGroupInfo(111, id)
		t.Log("Get GroupInfo $", id, "$: ", result)
	}
	result := GetGroupInfo(999, 1)
	t.Log("Get GroupInfo with wrong activityId", result)
}

func Test_UpdateUser(t *testing.T) {
	idSlice := []float64{1, 2, 666}
	for _, id := range idSlice {
		testUser := make(map[string]interface{})
		testUser["id"] = id
		testUser["openId"] = fmt.Sprintf("updated%d", int(id))
		testUser["paid"] = false
		result := UpdateUser(testUser)
		t.Log("UpdateUser: ", result)
	}
}

func Test_AddRecord(t *testing.T) {
	count := 999999
	for _, openId := range UserOpenId {
		paidId := strconv.Itoa(count)
		payment := Record{OpenId: openId, ActivityId: 111, PaidId: paidId, PaidAmount: 99}
		result := AddRecord(payment)
		if result != nil {
			logs.Notice(result.Error())
		}
		count--
	}
}

func Test_GetRefundRecords(t *testing.T) {
	result := GetRefundRecords(111)
	for _, item := range result["records"].([]*Record) {
		t.Log(fmt.Sprintf("%v", item))
	}
}
