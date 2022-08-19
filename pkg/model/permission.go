package model

import (
	"time"
)

const PermissionTable = "permission"

const (
	IdInPermission        = "id"
	RoleIdInPermission    = "role_id"
	ResIdInPermission     = "res_id"
	ResPowerInPermission  = "res_power"
	LinkIdInPermission    = "link_id"
	PowerInfoInPermission = "power_info"
	CreatedAtInPermission = "created_at"
	UpdatedAtInPermission = "updated_at"
)

type Permission struct {
	Id        uint64    `xorm:"not null pk autoincr BIGINT" json:"id"`
	RoleId    uint64    `xorm:"not null comment('角色id') index BIGINT" json:"roleId"`
	ResId     uint64    `xorm:"not null comment('资源id') index(idx_res_link_id) BIGINT" json:"resId"`
	ResPower  uint      `xorm:"comment('二进制表示，从右到左的二进制位表示select，update、insert、delete、manage') INT" json:"resPower"`
	LinkId    uint64    `xorm:"default 0 index(idx_res_link_id) BIGINT" json:"linkId"`
	PowerInfo string    `xorm:"comment('权限附加信息') JSON" json:"powerInfo"`
	CreatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
	UpdatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"updatedAt"`
}
