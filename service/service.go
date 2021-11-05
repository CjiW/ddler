package service

import (
	"main/sql"
	"main/tools"
	"strings"
	"time"
)

// 针对飞书事件

func TaskPublish(eventJson tools.Event)  {
	var name, start, end, publisherid string
	var receivers []string
	for i:=0;i<len(eventJson.Event.Message.Mentions);i++{
		receivers = append(receivers, eventJson.Event.Message.Mentions[i].Id.OpenId)
	}
	publisherid = eventJson.Event.Sender.SenderId.OpenId
	alist := strings.Split(eventJson.Event.Message.Content,"\\n")
	if len(alist)<4 {return}
	name = alist[2]
	start = strings.Split(alist[3], "-")[0]
	end = strings.Split(alist[3], "-")[1]
    end = end[:len(end)-2]
	endTime,_ := time.Parse("2006/01/02", end)
	startTime,_ := time.Parse("2006/01/02", start)
	err :=sql.AddTask(tools.Task{Name: name,SenderId: publisherid,Taskcontent: "", UndoneId: receivers,Start: startTime,End: endTime,Status: 0})
	if err != nil {
		return
	}
	nowTime := time.Now()
	msg := tools.NewMsg(publisherid,name,end,int(endTime.Sub(nowTime).Hours()/24))
    tools.SendMsg(msg,"card", receivers)
}
func TaskFinish(eventJson tools.Event) {
	var taskName string
	var doneId string
	alist := strings.Split(eventJson.Event.Message.Content,"\\n")
	if len(alist)<2 {return}
    taskName = alist[1][:len(alist[1])-2]
	doneId = eventJson.Event.Sender.SenderId.OpenId
	task,msg := sql.FinishTask(taskName,doneId)
    if	msg == "已完成："+taskName{
		endTime := task.End
		nowTime := time.Now()
		tools.SendMsg(
			tools.FinishMsg(doneId, taskName,
				task.End.Format("2006/01/02"),
				int((endTime.Sub(nowTime)).Hours()/24)),
			"card",task.UndoneId)
		tools.SendMsg(
			tools.FinishMsg(doneId, taskName,
				task.End.Format("2006/01/02"),
				int((endTime.Sub(nowTime)).Hours()/24)),
			"card",[]string{task.SenderId})}
	tools.SendMsg(msg,"text",[]string{doneId})
}
func TaskRemind(eventJson tools.Event) {
	alist := strings.Split(eventJson.Event.Message.Content,"\\n")
	if len(alist)<2 {
		return
	}
	content := alist[1]
	taskName := content[:len(content)-2]
	senderId := eventJson.Event.Sender.SenderId.OpenId

	task := sql.GetTaskInf(taskName)
	tools.SendMsg(
		tools.RemindMsg(senderId,taskName,
			task.End.Format("2006/01/02"),
			int(task.End.Sub(time.Now()).Hours()/24)),
		"card",task.UndoneId)
		tools.SendMsg("提醒已发出！","text",[]string{senderId})

}


