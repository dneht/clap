package model

import (
	"time"
)

const UserPermissionTable = "user_permission"

const (
	UserIdInUserPermission    = "user_id"
	PowIdInUserPermission     = "pow_id"
	PowTypeInUserPermission   = "pow_type"
	ResIdInUserPermission     = "res_id"
	ResTypeInUserPermission   = "res_type"
	CreatedAtInUserPermission = "created_at"
	UpdatedAtInUserPermission = "updated_at"
)

type UserPermission struct {
	UserId    int64     `xorm:"not null pk autoincr BIGINT(20)" json:"userId"`
	PowId     int64     `xorm:"not null pk comment('权限id') BIGINT(20)" json:"powId"`
	PowType   int       `xorm:"not null comment('权限类型、展开res_power') INT(11)" json:"powType"`
	ResId     int64     `xorm:"not null comment('资源id') BIGINT(20)" json:"resId"`
	ResType   int       `xorm:"not null comment('资源类型，如1->plume_menu、10000->inner_project') INT(10)" json:"resType"`
	CreatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') TIMESTAMP" json:"updatedAt"`
}
