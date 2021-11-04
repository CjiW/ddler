package tools

import "time"


type eventHeader struct {
	EventId    string `json:"event_id"`
	Token      string `json:"token"`
	CreateTime string `json:"create_time"`
	EventType  string `json:"event_type"`
	TenantKey  string `json:"tenant_key"`
	AppId      string `json:"app_id"`
}
type personId struct {
	UnionId string `json:"union_id,omitempty"`
	UserId  string `json:"user_id,omitempty"`
	OpenId  string `json:"open_id,omitempty"`
}

type senderData struct {
	SenderId   personId `json:"sender_id"`
	SenderType string   `json:"sender_type,omitempty"`
}
type Mention struct {
	Key  string   `json:"key,omitempty"`
	Id   personId `json:"id"`
	Name string   `json:"name,omitempty"`
}
type messageData struct {
	MessageId   string    `json:"message_id,omitempty"`
	RootId      string    `json:"root_id,omitempty"`
	ParentId    string    `json:"parent_id,omitempty"`
	CreateTime  string    `json:"create_time,omitempty"`
	ChatId      string    `json:"chat_id,omitempty"`
	ChatType    string    `json:"chat_type,omitempty"`
	MessageType string    `json:"message_type,omitempty"`
	Content     string    `json:"content,omitempty"`
	Mentions    []Mention `json:"mentions,omitempty"`
}
type eventData struct {
	Sender  senderData  `json:"sender"`
	Message messageData `json:"message"`
}
type baseModel struct {
	Schema string      `json:"schema,omitempty"`
	Header eventHeader `json:"header"`
}
type Event struct {
	baseModel
	Event eventData `json:"event"`
}


// type CardAction struct {
// 	Open_id string `json:"open_id"`//用户的open_id
// 	User_id string `json:"user_id"`//用户的user_id
// 	Open_message_id string `json:"open_message_id"`//触发交互操作的消息id
// 	Tenant_key string `json:"tenant_key"`//消息归属的租户id
// 	Token string  `json:"token"`//用于更新消息卡片的token（凭证）
// 	Action CAction `json:"action"`//具体的交互信息
// }
// type CAction struct {
// 	Value Value `json:"value"`
// 	Tag string  `json:"tag"`
// }
// type Value struct {
// 	Key int  `json:"key"`
// }


type Member struct {
	MemberId     string `json:"member_id"`
	MemberIdType string `json:"member_id_type"`
	Name         string `json:"name"`
	TenantKey    string `json:"tenant_key"`
}
type Data struct {
	HasMore     bool     `json:"has_more"`
	Items       []Member `json:"items"`
	MemberTotal int      `json:"member_total"`
	PageToken   string   `json:"page_token"`
}
type Chat struct {
	Code int    `json:"code"`
	Data Data   `json:"Data"`
	Msg  string `json:"msg"`
}

type User struct {
	Name   string
	OpenId string
}


type H5json2 struct {
	Option      uint     `json:"option"` //0 - 提醒、1 - 修改
	Taskid      uint     `json:"taskid"` //考虑到 任务名称 可能被修改，以 id 作为唯一标识
	Taskname    string   `json:"taskname"`
	TaskContent string   `json:"taskcontent"`
	Sender      string   `json:"sender"`
	Persons     []Worker `json:"persons"`
	Start       string   `json:"start"`
	Ddl         string   `json:"ddl"`    //这个格式不太清楚
	Isdone      int      `json:"isdone"` //0-未完成，1-已完成
}
type Worker struct {
	Person string `json:"person"`
	Status uint   `json:"status"` //0 - 未完成，1 - 已完成
}



type Task struct {
	Id          uint
	Name        string
	Taskcontent string
	SenderId    string
	UndoneId    []string
	DoneId      []string
	Start       time.Time
	End         time.Time
	Status      int        // 0-进行中，1-已完成
}
type TaskData struct {
	Id       uint
	Name     string
	Taskcontent string
	SenderId string
	UndoneId string
	DoneId   string
	Start    time.Time
	End      time.Time
	Status   int
}

