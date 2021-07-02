package model

import (
	"time"
)

const RolePermissionTable = "role_permission"

const (
	RoleIdInRolePermission    = "role_id"
	PowIdInRolePermission     = "pow_id"
	PowTypeInRolePermission   = "pow_type"
	ResIdInRolePermission     = "res_id"
	ResTypeInRolePermission   = "res_type"
	CreatedAtInRolePermission = "created_at"
	UpdatedAtInRolePermission = "updated_at"
)

type RolePermission struct {
	RoleId    int64     `xorm:"not null pk autoincr BIGINT(20)" json:"roleId"`
	PowId     int64     `xorm:"not null pk comment('权限id') BIGINT(20)" json:"powId"`
	PowType   int       `xorm:"not null comment('权限类型、展开res_power') INT(11)" json:"powType"`
	ResId     int64     `xorm:"not null comment('资源id') BIGINT(20)" json:"resId"`
	ResType   int       `xorm:"not null comment('资源类型，如1->plume_menu、10000->inner_project') INT(10)" json:"resType"`
	CreatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') TIMESTAMP" json:"updatedAt"`
}
