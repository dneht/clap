package model

import (
	"time"
)

const OperationLogTable = "operation_log"

const (
	IdInOperationLog          = "id"
	UserIdInOperationLog      = "user_id"
	ResIdInOperationLog       = "res_id"
	LogTypeInOperationLog     = "log_type"
	LogInfoInOperationLog     = "log_info"
	RequestFromInOperationLog = "request_from"
	CreatedAtInOperationLog   = "created_at"
)

type OperationLog struct {
	Id          uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	UserId      uint64    `xorm:"not null comment('用户id') index BIGINT(20)" json:"userId"`
	ResId       uint64    `xorm:"not null comment('资源id') index BIGINT(20)" json:"resId"`
	LogType     int       `xorm:"not null comment('操作类型') INT(11)" json:"logType"`
	LogInfo     string    `xorm:"comment('具体内容') JSON" json:"logInfo"`
	RequestFrom string    `xorm:"not null comment('来源信息，如ip、method、path等') JSON" json:"requestFrom"`
	CreatedAt   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') TIMESTAMP" json:"createdAt"`
}
