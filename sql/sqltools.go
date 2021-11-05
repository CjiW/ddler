package sql

import (
	"errors"
	"main/consts"
	"main/tools"
	"time"

	"gorm.io/gorm"
)

func UpdateUserlist(chatid string) {

	openids,names :=tools.GetMember(chatid)
	for i := 0; i < len(openids); i++ {
		var user tools.User
		errNameIsNotExisting :=consts.GlobalDB.First(&user,"open_id=?",openids[i]).Error
		if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){
		consts.GlobalDB.Create(&tools.User{Name:names[i],OpenId: openids[i]})
	}else {
		user.Name = names[i]
		consts.GlobalDB.Save(&user)
	}
	}
}
func GetUserInf1(username string) tools.User {
	var user tools.User
	errNameIsNotExisting :=consts.GlobalDB.First(&user,"name=?",username).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){
		return tools.User{Name: username,OpenId: ""}
	}else {return user}

}
func GetUserInf2(userid string) tools.User {
	var user tools.User
	errNameIsNotExisting :=consts.GlobalDB.First(&user,"open_id=?",userid).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){
		return tools.User{Name: "",OpenId: userid}
	}else {return user}

}


//AddTask 创建任务
func AddTask(task tools.Task) error {
	_ = consts.GlobalDB.AutoMigrate(&tools.TaskData{})
	errNameIsNotExisting :=consts.GlobalDB.First(&tools.TaskData{},"name=?",task.Name).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){
        taskdata := task.TurnIn()
        _ = consts.GlobalDB.Create(&taskdata)
		tools.SendMsg("任务发布成功:"+task.Name,"text",[]string{task.SenderId})
		return nil
	}else {
		tools.SendMsg("任务已存在:"+task.Name,"text",[]string{task.SenderId})		
		return errors.New("alreadyExisting")
	}
}

//GetTaskInf 通过任务名取出 task
func GetTaskInf(taskName string) tools.Task {
	var taskdata tools.TaskData
	var task tools.Task

	errNameIsNotExisting :=consts.GlobalDB.First(&taskdata,"name=?",taskName).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){return tools.Task{}}
	task = taskdata.TurnOut()
	return task

}

//FindTask 通过任务ID取出 task (针对H5)
func FindTask(taskId uint) (tools.Task, error) {
	var taskdata tools.TaskData
	var task tools.Task

	errNameIsNotExisting :=consts.GlobalDB.First(&taskdata,"id=?", taskId).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){return tools.Task{}, errors.New("name is not existing")}
	task = taskdata.TurnOut()
	return task, nil
}

//FinishTask 执行完成任务后的状态修改及提醒
func FinishTask(taskName, doneId string) (tools.Task,string) {
	var (
		taskdata tools.TaskData
		task tools.Task
		msg string
	)
	errNameIsNotExisting :=consts.GlobalDB.First(&taskdata,"name=?",taskName).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){
		msg="未找到任务："+taskName
	}else {
		task = taskdata.TurnOut()
		msg = "您不在未完成列表中"
		for i :=0;i<len(task.UndoneId);i++ {
			if task.UndoneId[i]==doneId{
				task.UndoneId = append(task.UndoneId[:i], task.UndoneId[i+1:]...)
				task.DoneId = append(task.DoneId, doneId)
				msg = "已完成："+taskName}
		}
		
		if len(task.UndoneId)==0{task.Status = 1}
		taskdata = task.TurnIn()
        consts.GlobalDB.Save(&taskdata)
	}
	return task, msg
}

//AutoRemind 遍历任务，自动提醒
func AutoRemind() {
	var taskdatas []tools.TaskData
	var tasks []tools.Task
	consts.GlobalDB.Find(&taskdatas)
	for i:=0;i<len(taskdatas);i++{
		tasks = append(tasks, taskdatas[i].TurnOut())
	}
	for i:=0;i<len(tasks);i++{
		if tasks[i].Status == 1{continue}
		leftDay := int(tasks[i].End.Sub(time.Now()).Hours()/24)
		switch leftDay {
		case 7:
			tools.SendMsg(tools.RemindMsg(tasks[i].SenderId,tasks[i].Name,tasks[i].End.Format("2006/01/02"),7),"text",tasks[i].UndoneId)
		case 3:
			tools.SendMsg(tools.RemindMsg(tasks[i].SenderId,tasks[i].Name,tasks[i].End.Format("2006/01/02"),3),"text",tasks[i].UndoneId)
		case 1:
			tools.SendMsg(tools.RemindMsg(tasks[i].SenderId,tasks[i].Name,tasks[i].End.Format("2006/01/02"),1),"text",tasks[i].UndoneId)
		default:
			continue
		}

}}
func Remind(taskId uint) error{
	task,err := FindTask(taskId)
	if task.Status == 0 {
	tools.SendMsg(tools.RemindMsg(task.SenderId,task.Name,task.End.Format("2006/01/02"),int(task.End.Sub(time.Now()).Hours()/24)),
		"card",
		task.UndoneId)}
	return err
}

func ChangeTask(json2 tools.H5json2) error {
	var taskdata tools.TaskData
	errIdIsNotExisting :=consts.GlobalDB.First(&taskdata,"id=?",json2.Taskid).Error
	if errors.Is(errIdIsNotExisting, gorm.ErrRecordNotFound){return errIdIsNotExisting}
	start,_ := time.Parse("2006/01/02",json2.Start)
	end,_ := time.Parse("2006/01/02",json2.Ddl)
	var doneids []string
	var undoneids []string
	for _, person := range json2.Persons {
		if person.Status == 0{undoneids = append(undoneids,GetUserInf1(person.Person).OpenId)
		}else {doneids = append(doneids,GetUserInf1(person.Person).OpenId)}
	}
	task := tools.Task{
		Id:       	uint(json2.Taskid),
		Name:     	json2.Taskname,
		Taskcontent:json2.TaskContent,
		SenderId:	GetUserInf1(json2.Sender).OpenId,
		UndoneId: 	undoneids,
		DoneId:   	doneids,
		Start:    	start,
		End:     	 end,
		Status:   	json2.Isdone,
	}
	taskdata = task.TurnIn()
	consts.GlobalDB.Save(taskdata)
	return nil
}