package api

import (
	"fmt"
	"main/consts"
	"main/service"
	"main/sql"
	"main/tools"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleMsg(ctx *gin.Context)  {
	var eventJson tools.Event

	err := ctx.ShouldBindJSON(&eventJson)
	if err != nil {fmt.Print(err)}
	sql.UpdateUserlist(eventJson.Event.Message.ChatId)
	if len(eventJson.Event.Message.Content)<=6 || eventJson.Event.Message.Content[2:6] != "text"{return}

	str := strconv.Quote(eventJson.Event.Message.Content)
	str = str[13:len(str)-4]
	func (str string){
		isMatch,_ := regexp.MatchString(`^【任务发布】*`, str)
		if isMatch {service.TaskPublish(eventJson)}
		isMatch,_ = regexp.MatchString(`^【任务完成】*`, str)
		if isMatch {service.TaskFinish(eventJson)}
		isMatch,_ = regexp.MatchString(`^【任务提醒】*`, str)
		if isMatch {service.TaskRemind(eventJson)}}(str)

}

func HandleH5get(ctx *gin.Context)  {
	var taskdatas []tools.TaskData
	var tasks []tools.Taskforh5
	consts.GlobalDB.Find(&taskdatas)
	for _, taskdata := range taskdatas {
		task := taskdata.TurnOut()
		fmt.Print(task,"\n")
		done := []string{}
		for _, id := range task.DoneId {
			name :=sql.GetUserInf2(id).Name
			for _, n := range done {if n==name{name = ""}}
			if  name == ""{continue}
			done = append(done, name)
		}
		notdone := []string{}
		for _, id := range task.UndoneId {
			name :=sql.GetUserInf2(id).Name
			for _, n := range notdone {if n==name{name = ""}}
			if  name == ""{continue}
			
			notdone = append(notdone, name)			
		}
		taskforh5 :=tools.Taskforh5{
			Taskid:      int(task.Id),
			Taskname:    task.Name,
			Taskcontent: task.Taskcontent,
			Sender:      sql.GetUserInf2(task.SenderId).Name,
			Notdone:     notdone,
			Done:        done,
			Start:       task.Start.Format("2006/01/02"),
			Ddl:         task.End.Format("2006/01/02"),
			Isdone:      task.Status,
		}
		
		tasks = append(tasks, taskforh5)
		
	}
	
	ctx.JSON(200,gin.H{"tasks":tasks})
}
func HandleH5post(ctx *gin.Context) {
	var jsonData tools.H5json2
	_ = ctx.ShouldBindJSON(&jsonData)
	fmt.Print(jsonData,"\n")
	switch jsonData.Option {
	case 0:
		err := sql.Remind(uint(jsonData.Taskid))
		if err != nil {
			ctx.JSON(200,gin.H{"msg":err.Error()})
			return
		}else {ctx.JSON(200,gin.H{"msg":"ok"})}

	case 1:
		err := sql.ChangeTask(jsonData)
		if err != nil {
			ctx.JSON(200,gin.H{"msg":err.Error()})
		}else {ctx.JSON(200,gin.H{"msg":"ok"})}
	}
}
