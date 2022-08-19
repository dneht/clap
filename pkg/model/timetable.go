package model

import (
	"time"
)

const TimetableTable = "timetable"

const (
	IdInTimetable         = "id"
	TaskNameInTimetable   = "task_name"
	TaskDescInTimetable   = "task_desc"
	TaskTypeInTimetable   = "task_type"
	TaskInfoInTimetable   = "task_info"
	TaskCronInTimetable   = "task_cron"
	TaskStatusInTimetable = "task_status"
	IsDisableInTimetable  = "is_disable"
	CreatedAtInTimetable  = "created_at"
)

type Timetable struct {
	Id         uint64    `xorm:"not null pk autoincr BIGINT" json:"id"`
	TaskName   string    `xorm:"not null comment('任务名') unique VARCHAR(64)" json:"taskName"`
	TaskDesc   string    `xorm:"not null comment('任务描述') VARCHAR(256)" json:"taskDesc"`
	TaskType   string    `xorm:"not null comment('任务类型，用于找到处理器') VARCHAR(32)" json:"taskType"`
	TaskInfo   string    `xorm:"comment('任务扩展信息') JSON" json:"taskInfo"`
	TaskCron   string    `xorm:"default '' comment('执行计划，cron表达式') VARCHAR(32)" json:"taskCron"`
	TaskStatus int       `xorm:"default 0 comment('执行状态，0等待中、1执行中、3重试中') TINYINT" json:"taskStatus"`
	IsDisable  int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
}
