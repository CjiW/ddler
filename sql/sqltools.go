package sql

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"main/consts"
	"main/tools"
	"time"
)

func UpdateUserlist() {
	db, _ := gorm.Open(mysql.Open(consts.DATABASE2), &gorm.Config{})
	_ = db.AutoMigrate(&tools.User{})
	openids,names :=tools.GetMember(consts.CHATID)
	for i := 0; i < len(openids); i++ {
		var user tools.User
		errNameIsNotExisting :=db.First(&user,"openid=?",openids[i]).Error
		if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){
		db.Create(&tools.User{Name:names[i],OpenId: openids[i]})
	}else {
		user.Name = names[i]
		db.Save(&user)
	}
	}
}
func GetUserInf(username string) tools.User {
	var user tools.User
	db, _ := gorm.Open(mysql.Open(consts.DATABASE2), &gorm.Config{})
	errNameIsNotExisting :=db.First(&user,"name=?",username).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){
		return tools.User{Name: username,OpenId: ""}
	}else {return user}

}

//AddTask 创建任务
func AddTask(task tools.Task) error {
	db, _ := gorm.Open(mysql.Open(consts.DATABASE1), &gorm.Config{})
	_ = db.AutoMigrate(&tools.TaskData{})
	errNameIsNotExisting :=db.First(&tools.TaskData{},"name=?",task.Name).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){
        taskdata := task.TurnIn()
        _ = db.Create(&taskdata)
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

	db, _ := gorm.Open(mysql.Open(consts.DATABASE1), &gorm.Config{})
	errNameIsNotExisting :=db.First(&taskdata,"name=?",taskName).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){return tools.Task{}}
	task = taskdata.TurnOut()
	return task

}

//FindTask 通过任务ID取出 task (针对H5)
func FindTask(taskId uint) (tools.Task, error) {
	var taskdata tools.TaskData
	var task tools.Task

	db, _ := gorm.Open(mysql.Open(consts.DATABASE1), &gorm.Config{})
	errNameIsNotExisting :=db.First(&taskdata,"id=?", taskId).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){return tools.Task{}, errors.New("name is not existing")}
	task = taskdata.TurnOut()
	return task, nil
}

//FinishTask 执行完成任务后的状态修改及提醒
func FinishTask(taskName, doneId string) tools.Task {
	var (
		taskdata tools.TaskData
		task tools.Task
	)
	db, _ := gorm.Open(mysql.Open(consts.DATABASE1), &gorm.Config{})
	errNameIsNotExisting :=db.First(&taskdata,"name=?",taskName).Error
	if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){

	}else {
		task = taskdata.TurnOut()
		for i :=0;i<len(task.UndoneId);i++ {
			if task.UndoneId[i]==doneId{
				task.UndoneId = append(task.UndoneId[:i], task.UndoneId[i+1:]...)}
		}
		task.DoneId = append(task.DoneId, doneId)
		if len(task.UndoneId)==0{task.Status = 1}
		taskdata = task.TurnIn()
        db.Save(&taskdata)
	}
	return task
}

//AutoRemind 遍历任务，自动提醒
func AutoRemind() {
	var taskdatas []tools.TaskData
	var tasks []tools.Task
	db, _ := gorm.Open(mysql.Open(consts.DATABASE1), &gorm.Config{})
	db.Find(&taskdatas)
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
	tools.SendMsg(tools.RemindMsg(task.SenderId,task.Name,task.End.Format("2006/01/02"),int(task.End.Sub(time.Now()).Hours()/24)),
		"card",
		task.UndoneId)
	return err
}

func ChangeTask(json2 tools.H5json2) error {
	var taskdata tools.TaskData
	db, _ := gorm.Open(mysql.Open(consts.DATABASE1), &gorm.Config{})
	errIdIsNotExisting :=db.First(&taskdata,"id=?",json2.Taskid).Error
	if errors.Is(errIdIsNotExisting, gorm.ErrRecordNotFound){return errIdIsNotExisting}
	start,_ := time.Parse("2006/01/02",json2.Start)
	end,_ := time.Parse("2006/01/02",json2.Ddl)
	var doneids []string
	var undoneids []string
	for _, person := range json2.Persons {
		if person.Status == 0{undoneids = append(undoneids,person.Person)
		}else {doneids = append(doneids,person.Person)}
	}
	task := tools.Task{
		Id:       	json2.Taskid,
		Name:     	json2.Taskname,
		Taskcontent:json2.TaskContent,
		SenderId:	GetUserInf(json2.Sender).OpenId,
		UndoneId: 	undoneids,
		DoneId:   	doneids,
		Start:    	start,
		End:     	 end,
		Status:   	json2.Isdone,
	}
	taskdata = task.TurnIn()
	db.Save(taskdata)
	return nil
}
