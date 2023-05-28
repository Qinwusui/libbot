package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
)

func PrivateMsgCheck(qq int64, cqm string) {
	if len(cqm) > 6 && cqm[:6] == "/start" {
		setupTime, err := time.ParseInLocation("/start 2006 1 2", cqm, time.FixedZone("UTC", 8*24*60))
		nowTime := time.Now()
		diffTime := setupTime.Sub(nowTime)

		if err != nil {
			log.Println(err)
			go SendMsg(PrivateTextMsg(qq, "开学日期输入不正确,请重新输入"+err.Error()))
			return
		}
		var strMsg string
		diffDay := diffTime.Hours() / 24
		if SaveStartTime(&setupTime, qq) {

			if diffDay < 0 {
				strMsg = fmt.Sprintf("设置开学日期成功,已经开学%.2f天", math.Abs(diffDay))
			} else {
				strMsg = fmt.Sprintf("设置开学日期成功,距离开学还有%.2f天", diffDay)
			}
		} else {
			if diffDay < 0 {
				strMsg = fmt.Sprintf("覆写开学日期成功,已经开学%.2f天", math.Abs(diffDay))
			} else {
				strMsg = fmt.Sprintf("覆写开学日期成功,距离开学还有%.2f天", diffDay)

			}
		}
		go SendMsg(PrivateTextMsg(qq, strMsg))

	}
	if strings.Contains(cqm, "/setup") {

		regex := regexp.MustCompile(`/setup \d{11} \w{6,18}`)
		if regex.MatchString(cqm) {
			var userId, pwd string
			_, _ = fmt.Sscanf(cqm, "/setup %s %s", &userId, &pwd)
			var msg string
			if SaveXXTInfo(userId, pwd, qq) {
				//首次保存
				msg = "信息已经保存在机器人所在服务器中！除登录外机器人不会用于其他用途。若填写错误请修改后重新发送即可"
			} else {
				msg = "已更新信息"
			}
			go SendMsg(PrivateTextMsg(qq, msg))
		} else {
			go SendMsg(PrivateTextMsg(qq, "指令不正确，指令格式：/setup 学号 密码\n且必须在http://jwxt.wut.edu.cn/admin中修改密码才能正常登录"))
		}
	}
	//reg := regexp.MustCompile(`/notify 每周\d \d{1,2}时\d{1,2}分`)
	//if reg.MatchString(cqm) {
	//	//设置定时推送时间，只能是每周一次
	//	var d, h, m int
	//	_, err := fmt.Sscanf(cqm, "/notify 每周%d %d时%d分", &d, &h, &m)
	//	if err != nil || d < 1 || d > 7 || h < 0 || h > 24 || m < 0 || m > 60 {
	//		go SendMsg("请检查输入的时间是否正确")
	//		return
	//	}
	//
	//	msg := "设置"
	//	if !checkNeedScheduler(qq, true, d, h, m) {
	//		msg = "覆写"
	//	}
	//
	//	SendMsg(PrivateTextMsg(qq, msg+fmt.Sprintf("提醒日期成功，在每周%d会通过短信推送您当周的课表", d)))
	//}
	//if cqm == "/notnotify" {
	//	checkNeedScheduler(qq, false, 0, 0, 0)
	//	go SendMsg("你已取消每周课表推送~\n如需再次使用只需再次使用/notify指令即可")
	//}
	// if cqm == "/开启定时推送课表" {

	// }
	// if cqm == "/关闭定时推送课表" {

	// }
}

func SaveXXTInfo(phone, pwd string, qq int64) bool {
	m := GetMemberMap()
	check := false
	qqStr := fmt.Sprintf("%d", qq)
	v := m[qqStr]
	if v == nil {
		//说明这个QQ还没有信息保存
		m[qqStr] = map[string]interface{}{
			"phone": phone,
			"pwd":   pwd,
		}
		check = true
	} else {
		p := m[qqStr].(map[string]interface{})["phone"]
		w := m[qqStr].(map[string]interface{})["pwd"]
		check = (p == nil) && (w == nil)
		m[qqStr].(map[string]interface{})["phone"] = phone
		m[qqStr].(map[string]interface{})["pwd"] = pwd
	}
	bytes, err := json.MarshalIndent(m, " ", "  ")
	checkErr(err)
	_ = os.WriteFile(timeFile, bytes, 0644)
	return check
}
func ReadMemberInfo(qq int64) (*string, *string) {
	m := GetMemberMap()
	qqStr := fmt.Sprintf("%d", qq)
	v := m[qqStr]
	if v == nil {
		return nil, nil
	}
	pi := m[qqStr].(map[string]interface{})["phone"]
	var p string
	var w string
	if pi == nil {
		p = ""
	} else {
		p = pi.(string)
	}

	wi := m[qqStr].(map[string]interface{})["pwd"]
	if wi == nil {
		w = ""
	} else {
		w = wi.(string)
	}
	return &p, &w
}
