package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"main/consts"
	"net/http"
	"strconv"
	"strings"
)

//GetToken 获取 token
func GetToken(appId, appSecret string) string {
	url := "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal/"
	method := "POST"

	payload := strings.NewReader("{\"app_id\": \"" + appId + "\",\"app_secret\": \"" + appSecret + "\"}")

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, payload)

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	//req.Header.Add("Cookie", "swp_csrf_token=e3ed77ee-471a-4166-b704-772f4625961f; t_beda37=e1f9e03bb3ba6803dc925874283ba1b50c5c3102c79baf6434076e407c708fa2")

	res, _ := client.Do(req)

	defer func(*http.Response) { res.Body.Close() }(res)

	body, _ := ioutil.ReadAll(res.Body)

	return string(body)[122:164]

}

//SendMsg 发信息，支持text与card类型
func SendMsg(msg, msgType string, recieverIds []string) {
	tenantToken := GetToken(consts.APPID, consts.APPSECRET)
	for i := 0; i < len(recieverIds); i++ {
		var payload *strings.Reader
		var url string
		if msgType == "text" {
			payload = strings.NewReader("{\"receive_id\": \"" + recieverIds[i] + "\",\"content\": \"{\\\"text\\\":\\\"" + msg + "\\\"}\",\"msg_type\": \"text\"}")
			url = consts.FEISHUOPENURLv1
		} else if msgType == "card" {
			payload = strings.NewReader("{\"open_id\": \"" + recieverIds[i] + "\"," + msg + "}")
			url = consts.FEISHUOPENURLv4
		}

		client := &http.Client{}
		req, _ := http.NewRequest("POST", url, payload)

		req.Header.Add("Authorization", "Bearer "+tenantToken)
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		//req.Header.Add("Cookie", "swp_csrf_token=e3ed77ee-471a-4166-b704-772f4625961f; t_beda37=e1f9e03bb3ba6803dc925874283ba1b50c5c3102c79baf6434076e407c708fa2")
		{
			a, _ := client.Do(req)
			b, _ := ioutil.ReadAll(a.Body)
			fmt.Print(string(b))
		}

	}

}

//GetMember 获取群聊成员列表
func GetMember(chatId string) ([]string, []string) {
	tenantToken := GetToken(consts.APPID, consts.APPSECRET)
	url := "https://open.feishu.cn/open-apis/im/v1/chats/" + chatId + "/members?member_id_type=open_id"
	method := "GET"

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)

	req.Header.Add("Authorization", "Bearer "+tenantToken)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	//req.Header.Add("Cookie", "swp_csrf_token=e3ed77ee-471a-4166-b704-772f4625961f; t_beda37=e1f9e03bb3ba6803dc925874283ba1b50c5c3102c79baf6434076e407c708fa2")

	res, _ := client.Do(req)

	defer func(Body io.ReadCloser) { _ = Body.Close() }(res.Body)

	body, _ := ioutil.ReadAll(res.Body)

	var jsonData Chat
	_ = json.NewDecoder(strings.NewReader(string(body))).Decode(&jsonData)
	var allNameList, allIdList []string
	for i := 0; i < jsonData.Data.MemberTotal; i++ {
		allNameList = append(allNameList, jsonData.Data.Items[i].Name)
		allIdList = append(allIdList, jsonData.Data.Items[i].MemberId)
	}

	return allIdList, allNameList
}

func NewMsg(senderid, task, ddl string, left int) string {
	return `"msg_type": "interactive", "card": 
	{ "config": { 
		"wide_screen_mode": true 
		}, 
		"elements": [ 
			{ 
				"tag": "hr" 
			}, 
			{ 
				"fields": [
					{ 
						"is_short": true, 
						"text": { 
							"content": "**发布人：**\n<at id=` + senderid + `></at>", 
							"tag": "lark_md" 
						} 
					}, 
					{ 
						"is_short": true, 
						"text": { 
							"content": "**任务内容：**\n` + task + `", 
							"tag": "lark_md" 
						} 
					}, 
					{ 
						"is_short": false, 
						"text": { 
							"content": "", 
							"tag": "lark_md" 
						} 
					}, 
					{ 
						"is_short": true, 
						"text": { 
							"content": "**DDL：**\n` + ddl + `", 
							"tag": "lark_md" 
						} 
					}, 
					{ 
						"is_short": true, 
						"text": { 
							"content": "**剩余：**\n` + strconv.Itoa(left) + `天", 
							"tag": "lark_md" 
						} 
					} ,
        			{
          				"is_short": true,
          				"text": {
           			 		"content": "***尽快完成吧！！***",
            				"tag": "lark_md"
          				}
        			}
					], 
				"tag": "div" 
			}, 
			{ 
				"tag": "hr" 
			}, 
			{ 
				"tag": "note", 
				"elements": [ 
					{ 
						"is_short": true,
						"tag": "lark_md", 
						"content": "[任务列表](` + consts.LISTURL + `)" 
					},
					{ 
						"is_short": true,
						"tag": "lark_md", 
						"content": "[任务详情](` + consts.DETAILsURL + `)" 
					}
				] 
			}
		], 
		"header": { 
			"template": "Turquoise", 
			"title": { 
				"content": "你有新任务啦~~", 
				"tag": "plain_text" 
			} 
		} 
	}`

}
func RemindMsg(senderid, task, ddl string, left int) string {
	return `"msg_type": "interactive", "card":{
		"config": {
			"wide_screen_mode": true
		},
		"elements": [
	{
	"tag": "hr"
	},
	{
	"fields": [
	{
	"is_short": true,
	"text": {
	"content": "**发布人：**\n<at id=` + senderid + `></at>",
	"tag": "lark_md"
	}
	},
	{
	"is_short": true,
	"text": {
	"content": "**任务内容：**\n` + task + `",
	"tag": "lark_md"
	}
	},
	{
	"is_short": false,
	"text": {
	"content": "",
	"tag": "lark_md"
	}
	},
	{
	"is_short": true,
	"text": {
	"content": "**DDL：**\n` + ddl + `",
	"tag": "lark_md"
	}
	},
	{
	"is_short": true,
	"text": {
	"content": "**剩余：**\n` + strconv.Itoa(left) + `天",
	"tag": "lark_md"
	}
	},
	{
	"is_short": false,
	"text": {
	"content": "",
	"tag": "lark_md"
	}
	},
	{
	"is_short": true,
	"text": {
	"content": "***尽快完成吧！！***",
	"tag": "lark_md"
	}
	}
	],
	"tag": "div"
	},
	{
	"tag": "hr"
	},
	{
	"tag": "note",
	"elements": [
	{
	"is_short": true,
	"tag": "lark_md",
	"content": "[任务列表](` + consts.LISTURL + `)"
	},
	{
	"is_short": true,
	"tag": "lark_md",
	"content": "[任务详情](` + consts.DETAILsURL + `)"
	}
	]
	}
	],
	"header": {
	"template": "Orange",
	"title": {
	"content": "任务要尽快完成哦！",
	"tag": "plain_text"
	}
	}
}`
}
func FinishMsg(doneId, task, ddl string, left int) string {
	return `"msg_type": "interactive", "card":{
  	"config": {
    	"wide_screen_mode": true
 	 },
  	"elements": [
    	{
			"tag": "hr"
    	},
    	{
      		"fields": [
        	{
         	 	"is_short": true,
          		"text": {
            		"content": "**发布人：**\n<at id=` + doneId + `></at>",
            		"tag": "lark_md"
          		}
        	},
        	{
          		"is_short": true,
          		"text": {
            		"content": "**任务内容：**\n` + task + `",
            		"tag": "lark_md"
          		}
        	},
        	{
          		"is_short": false,
          		"text": {
            		"content": "",
            		"tag": "lark_md"
          		}
        	},
        	{
          		"is_short": true,
          		"text": {
            		"content": "**DDL：**\n` + ddl + `",
            		"tag": "lark_md"
          		}
			},
        	{
          		"is_short": true,
          		"text": {
            		"content": "**提前：**\n` + strconv.Itoa(left) + `天",
            		"tag": "lark_md"
          		}
        	}
      		],
      		"tag": "div"
    	},
    	{
      		"tag": "hr"
    	},
    	{
      		"tag": "note",
      		"elements": [
			{
          		"tag": "lark_md",
          		"content": "[任务列表](` + consts.DETAILsURL + `)"
        	}
      		]
    	}
		
  		],
  	"header": {
    	"template": "Wathet",
    	"title": {
      		"content": "我完成了！！",
      		"tag": "plain_text"
    	}
  	}
}`
}

//TurnIn 存入数据库的转换
//       切片 --> string
func (task Task) TurnIn() TaskData {
	return TaskData{
		Id:          task.Id,
		Name:        task.Name,
		Taskcontent: task.Taskcontent,
		SenderId:    task.SenderId,
		UndoneId:    strings.Join(task.UndoneId, ","),
		DoneId:      strings.Join(task.DoneId, ","),
		Start:       task.Start,
		End:         task.End,
		Status:      task.Status,
	}
}

//TurnOut 将结构体中 string 转换为易操作类型
func (data TaskData) TurnOut() Task {
	return Task{
		Id:          data.Id,
		Name:        data.Name,
		Taskcontent: data.Taskcontent,
		SenderId:    data.SenderId,
		UndoneId:    strings.Split(data.UndoneId, ","),
		DoneId:      strings.Split(data.DoneId, ","),
		Start:       data.Start,
		End:         data.End,
		Status:      data.Status,
	}
}
