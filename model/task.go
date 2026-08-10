package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name   string `gorm:"type:varchar(30) comment '任务名称'"`
	Cron   string `gorm:"type:varchar(60) comment 'cron表达式'"`
	Switch bool   `gorm:"comment '定时任务开关, true: 开, false: 关'"`
}

type TaskExecution struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	TaskID    uint
	Result    bool   `gorm:"comment '任务执行结果, true: 成功, false: 失败'"`
	Reason    string `gorm:"type:varchar(400) comment '任务执行失败原因'"`
}

func AddTasks(tasks []Task) error {
	return DB.Create(&tasks).Error
}

func FindTask() ([]Task, error) {
	var tasks []Task
	err := DB.Find(&tasks).Error
	return tasks, err
}

func CreateTaskExecution(taskExecution TaskExecution) error {
	return DB.Create(&taskExecution).Error
}
