package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsConn *websocket.Conn

type GroupMsg struct {
	Time       int64       `json:"time"`
	Sender     interface{} `json:"sender"`
	GroupId    int64       `json:"group_id"`
	Message    string      `json:"message"`
	RawMessage string      `json:"raw_message"`
	MessageId  int64       `json:"message_id"`
	Seq        int64       `json:"message_seq"`
	//Sender
	/*
		"sender": {
				"age": 0,
				"area": "",
				"card": "terraria-save the awp",
				"level": "",
				"nickname": "echo flower",
				"role": "admin",
				"sex": "unknown",
				"title": "?",
				"user_id": 472203400
			},
	*/
}

type PrivateMsg struct {
	Time       int64       `json:"time"`
	Sender     interface{} `json:"sender"`
	TargetId   int64       `json:"target_id"`
	UserId     int64       `json:"user_id"`
	Message    string      `json:"message"`
	RawMessage string      `json:"raw_message"`
	MessageId  int64       `json:"message_id"`
	/**
	"sender": {
			"age": 0,
			"nickname": "name",
			"sex": "unknown",
			"user_id": 2064508450
		},
	*/
}

// NoticeMsgDispatcher 通知消息分发处理
func NoticeMsgDispatcher(obj map[string]interface{}) {
	noticeType := obj["notice_type"]
	postType := obj["post_type"]
	subType := obj["sub_type"]
	if noticeType == nil {
		return
	}
	groupIdObj := obj["group_id"]
	userIdObj := obj["user_id"]
	messageIdObj := obj["message_id"]

	switch noticeType {
	case "group_recall":
		//群撤回消息
		{
			if groupIdObj == nil || userIdObj == nil || messageIdObj == nil {
				return
			}
			messageId := int64(messageIdObj.(float64))
			groupId := int64(groupIdObj.(float64))
			if groupId == 859260403 || groupId == 796993590 {
				return
			}
			send := map[string]interface{}{
				"action": "get_msg",
				"params": map[string]interface{}{
					"message_id": messageId,
				},
			}
			if wsConn == nil {
				log.Println("wsConn为空")
				return
			}
			err := wsConn.WriteJSON(send)
			checkErr(err)
		}
	case "group_increase":
		{
			groupId := int64(groupIdObj.(float64))
			userId := int64(obj["user_id"].(float64))
			switch groupId {
			case 329972361:
				{
					go SendMsg(AtMsg(groupId, userId, "欢迎加入本群，请修改群名片，格式terraria-xxx  并仔细阅读本群公告"))
				}
			}

		}
	case "notify":
		{
			if postType == "notice" {
				switch subType {
				case "poke":
					{
						//戳一戳

					}
				}
			}
		}
	}

}

// MsgDispatcher 消息分发，所有消息到此处分发
func MsgDispatcher(msgByte []byte) {
	obj := make(map[string]interface{})
	_ = json.Unmarshal(msgByte, &obj)
	msgType := obj["message_type"]
	metaEventType := obj["meta_event_type"]
	subType := obj["sub_type"]
	dataObj := obj["data"]
	retcodeObj := obj["retcode"]
	if dataObj != nil {
		data := dataObj.(map[string]interface{})
		if data == nil {
			return
		}
		msg := data["message"]
		if msg == nil {
			return
		}
		if retcodeObj != nil && int64(retcodeObj.(float64)) == 0 {
			GetMsgRes(data, func(qq, groupId int64, msg string) {
				switch groupId {
				case 226273990: //VV的群
					{
						if qq != 312137484 {
							SendMsg(PrivateTextMsg(312137484, fmt.Sprintf("[%v]在[%v]撤回了一条消息\n%v", qq, groupId, msg)))
						}
					}
				default: //默认处理
					{
						if qq == 2064508450 {
							return
						}
						SendMsg(PrivateTextMsg(2064508450, fmt.Sprintf("[%v]在[%v]撤回了一条消息\n%v", qq, groupId, msg)))
					}
				}

			})
			return
		}
	}

	NoticeMsgDispatcher(obj)
	if metaEventType != nil {
		return
	}
	//事件分发
	//群聊消息或者私聊消息
	switch msgType {
	case "group":
		switch subType {
		case "normal":
			GroupMsgHandler(obj)
		}

	case "private":
		switch subType {
		case "group":
			//群临时消息
			PrivateTmpMsgHandler(obj)
		default:
			PrivateMsgHandler(obj)
		}

	}

}

// 私聊消息转换
func PrivateTmpMsgHandler(obj map[string]interface{}) {
	privateMsg := TmpMsg2PrivateMsg(obj)
	message := privateMsg.Message
	qq := int64(privateMsg.Sender.(map[string]interface{})["user_id"].(float64))
	PrivateMsgCheck(qq, message)
}

// GetMsgRes 撤回消息后，通过getMsg获取到的消息
func GetMsgRes(obj map[string]interface{}, recall func(qq, groupId int64, msg string)) {
	groupIdObj := obj["group_id"]
	messageObj := obj["message"]
	senderObj := obj["sender"]
	if senderObj == nil || groupIdObj == nil || messageObj == nil {
		return
	}
	groupId := int64(groupIdObj.(float64))
	message := messageObj.(string)
	if message == "" {
		log.Println("message为空")
		return
	}
	userIdObj := senderObj.(map[string]interface{})["user_id"]
	if userIdObj == nil {
		log.Println("userId为空")
		return
	}
	userId := int64(userIdObj.(float64))
	recall(userId, groupId, message)
}

func PrivateMsgHandler(obj map[string]interface{}) {
	privateMsg := Msg2PrivateMsg(obj)
	message := privateMsg.Message
	qq := int64(privateMsg.Sender.(map[string]interface{})["user_id"].(float64))
	PrivateMsgCheck(qq, message)
}

func GroupMsgHandler(obj map[string]interface{}) {
	groupMsg := Msg2GroupMsg(obj)
	message := groupMsg.Message
	messageId := groupMsg.MessageId
	qqInterface := groupMsg.Sender.(map[string]interface{})["user_id"]
	nickname := groupMsg.Sender.(map[string]interface{})["nickname"]
	senderCard := groupMsg.Sender.(map[string]interface{})["card"].(string)
	groupId := groupMsg.GroupId
	seq := groupMsg.Seq
	if qqInterface == nil {
		log.Println("qq为空")
		return
	}
	if nickname == nil {
		log.Println("昵称为空")
		return
	}
	nick := nickname.(string)
	qq := int64(qqInterface.(float64))
	if groupId == 259893170 { //259893170 637081363
		CFEvent(qq, groupId, messageId, seq, nick, message)
	}
	CheckGroupMsg(qq, groupId, messageId, seq, nick, senderCard, message)
}

// 群临时消息转私聊消息
func TmpMsg2PrivateMsg(obj map[string]interface{}) *PrivateMsg {
	privateMsg := new(PrivateMsg)
	privateMsg.Sender = obj["sender"]
	privateMsg.Time = int64(obj["time"].(float64))
	privateMsg.RawMessage = obj["raw_message"].(string)
	privateMsg.Message = obj["message"].(string)
	privateMsg.MessageId = int64(obj["message_id"].(float64))
	return privateMsg
}

// 消息转私聊消息
func Msg2PrivateMsg(obj map[string]interface{}) *PrivateMsg {
	privateMsg := new(PrivateMsg)
	privateMsg.Sender = obj["sender"]
	privateMsg.Time = int64(obj["time"].(float64))
	privateMsg.TargetId = int64(obj["target_id"].(float64))
	privateMsg.RawMessage = obj["raw_message"].(string)
	privateMsg.Message = obj["message"].(string)
	privateMsg.MessageId = int64(obj["message_id"].(float64))
	return privateMsg
}

// Msg2GroupMsg 消息转群聊消息
func Msg2GroupMsg(obj map[string]interface{}) *GroupMsg {
	groupMsg := new(GroupMsg)
	groupMsg.GroupId = int64(obj["group_id"].(float64))
	groupMsg.Message = obj["message"].(string)
	groupMsg.Sender = obj["sender"]
	groupMsg.Time = int64(obj["time"].(float64))
	groupMsg.RawMessage = obj["raw_message"].(string)
	groupMsg.MessageId = int64(obj["message_id"].(float64))
	groupMsg.Seq = int64(obj["message_seq"].(float64))
	return groupMsg
}

// 消息捕获
func receiveHandler() {
	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			log.Println(err)
			wsConn.Close()
			wsConn = nil
			InitWs()
			break
		}
		MsgDispatcher(msg)
	}
}

func InitWs() {
	socketUrl := "ws://47.115.227.17:8080"
	pwd := base64.StdEncoding.EncodeToString([]byte("Qinsansui233..."))
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, http.Header{
		"Authorization": {fmt.Sprintf("Bearer %v", pwd)},
	})
	if err != nil {
		log.Println(err)
		return
	}
	wsConn = conn
	log.Println("建立Websocket成功")
}
