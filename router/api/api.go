package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
	"main/sql"
	"main/tools"
	"regexp"
	"strconv"
)

func HandleMsg(ctx *gin.Context)  {
	var eventJson tools.Event

	err := ctx.ShouldBindJSON(&eventJson)
	if err != nil {fmt.Print(err)}
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

}
func HandleH5post(ctx *gin.Context) {
	var jsonData tools.H5json2
	_ = ctx.ShouldBindJSON(&jsonData)
	fmt.Print(jsonData.Option)
	switch jsonData.Option {
	case 0:
		err := sql.Remind(jsonData.Taskid)
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
