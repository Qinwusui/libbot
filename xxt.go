package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var jar *cookiejar.Jar
var client *http.Client

func OpenJw(xuehao, pwd string) bool {
	request, _ := http.NewRequest("POST", "http://jwxt.wut.edu.cn/admin/login",
		strings.NewReader(fmt.Sprintf("username=%s&password=%s&rememberMe=", xuehao, pwd)),
	)
	request.Header = http.Header{
		"Accept":          {"text/html,application/xhtml+xml,application/xml;"},
		"Accept-Encoding": {"gzip, deflate"},
		"Connection":      {"keep-alive"},
		"Content-Type":    {"application/x-www-form-urlencoded"},
		"User-Agent":      {"Chrome/111.0.0.0"},
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return false
	}
	if url := resp.Header["Location"]; url != nil {
		return Redirect2Jw(fmt.Sprintf("http://jwxt.wut.edu.cn%s", url[0]))
	}
	return false
}
func Redirect2Jw(url string) bool {
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	checkErr(err)
	return resp.StatusCode == 200
}
func cleanCookie() {
	lock := sync.Mutex{}
	lock.Lock()
	jar, _ = cookiejar.New(nil)

	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}
	lock.Unlock()

}

// 获取课表
func GetCourse(
	week int,
	uname, pwd string,
	login func(string, string) bool,
	stuInfo func() (string, string),
	sendMsg func(string),
) *Course {
	b := login(uname, pwd)
	if !b {
		sendMsg("登录失败了哦~检查一下用户名和密码吧")
		return nil
	}
	id, xnxq := stuInfo()
	if id == "" || xnxq == "" {
		sendMsg("登录失败了哦~检查一下用户名和密码吧")
		return nil
	}
	url := "http://jwxt.wut.edu.cn/admin/api/getXskb?xnxq=" + xnxq + "&userId=" + id + "&xqid=&role=xs"
	if week >= 0 {
		url += fmt.Sprintf("&week=%v", week)
	}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header = http.Header{
		"Accept":          {"text/html,application/json;"},
		"Accept-Encoding": {"none"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {"Chrome/111.0.0.0"},
	}
	resp, err := client.Do(request)
	if err != nil {
		sendMsg("获取课表失败了哦~")
		return nil
	}
	bts, _ := io.ReadAll(resp.Body)

	course := new(Course)
	err = json.Unmarshal(bts, &course)
	if err != nil {
		sendMsg("获取课表失败了哦~")
		return nil
	}
	defer resp.Body.Close()
	kbData := course.Data.KckbData
	sort.SliceStable(kbData, func(i, j int) bool {
		if kbData[i].Xingqi < kbData[j].Xingqi {
			return true
		} else if kbData[i].Xingqi > kbData[j].Xingqi {
			return false
		}
		return kbData[i].Djc < kbData[j].Djc
	})
	course.Data.KckbData = kbData
	return course
}

func getStuInfo() (string, string) {
	header := http.Header{
		"Accept": {"text/html"},
	}
	request, _ := http.NewRequest("GET", "http://jwxt.wut.edu.cn/admin//xjgl/xsjbxx/xskp", nil)
	request.Header = header
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
	}

	id := doc.Find(`#id`).AttrOr("value", "")
	var t = time.Now()
	year := t.Year()
	month := t.Month()
	var xnxq string
	if month >= 2 && month < 9 {
		//表示在上半年，上半年的意思就是第二学期
		xnxq = fmt.Sprintf("%d-%d-2", year-1, year)
	}
	if month < 2 || month >= 9 {
		xnxq = fmt.Sprintf("%d-%d-1", year, year+1)
	}
	if id == "" || xnxq == "" {
		log.Println("Get ID And Xnxq Err")
	}
	return id, xnxq
}
