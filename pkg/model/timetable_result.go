package model

import (
	"time"
)

const TimetableResultTable = "timetable_result"

const (
	IdInTimetableResult         = "id"
	TaskIdInTimetableResult     = "task_id"
	LastStatusInTimetableResult = "last_status"
	LastResultInTimetableResult = "last_result"
	CreatedAtInTimetableResult  = "created_at"
)

type TimetableResult struct {
	Id         uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	TaskId     uint64    `xorm:"not null comment('任务id') index BIGINT(20)" json:"taskId"`
	LastStatus int       `xorm:"default 1 comment('执行状态，true成功、false失败') TINYINT(1)" json:"lastStatus"`
	LastResult string    `xorm:"comment('上次执行结果，有的任务可能有输出') JSON" json:"lastResult"`
	CreatedAt  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
}
