package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

const codeUser = "./code.json"

var rdb *redis.Client
var ctx, _ = context.WithCancel(context.Background())

// var c *cron.Cron

func Main() {
	log.Println("YikuYo!")
	log.Println("载入用户验证码列表")
	if checkFileIsExist(codeUser) {
		bytes, err := os.ReadFile(codeUser)
		checkErr(err)
		cu := new(CodeUser)
		_ = json.Unmarshal(bytes, &cu)
		codeUserData = append(codeUserData, cu)
	} else {
		log.Println("不存在用户获取验证码列表文件")
	}
	InitRedis()
	InitWs()
	// InitImageRec()
	go receiveHandler()
	go CFListener()
	select {}
}
