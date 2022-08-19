package model

import (
	"time"
)

const RoleInfoTable = "role_info"

const (
	IdInRoleInfo         = "id"
	RoleNameInRoleInfo   = "role_name"
	RoleFromInRoleInfo   = "role_from"
	RoleRemarkInRoleInfo = "role_remark"
	IsManageInRoleInfo   = "is_manage"
	IsSuperInRoleInfo    = "is_super"
	IsDisableInRoleInfo  = "is_disable"
	CreatedAtInRoleInfo  = "created_at"
	UpdatedAtInRoleInfo  = "updated_at"
)

type RoleInfo struct {
	Id         uint64    `xorm:"not null pk autoincr BIGINT" json:"id"`
	RoleName   string    `xorm:"not null comment('角色名') unique(uk_role_from_name) VARCHAR(64)" json:"roleName"`
	RoleFrom   uint      `xorm:"default 0 comment('角色来源、本系统0') unique(uk_role_from_name) INT" json:"roleFrom"`
	RoleRemark string    `xorm:"not null comment('备注信息') VARCHAR(256)" json:"roleRemark"`
	IsManage   int       `xorm:"default 0 comment('是否是管理角色') TINYINT(1)" json:"isManage"`
	IsSuper    int       `xorm:"default 0 comment('是否是超级管理角色') TINYINT(1)" json:"isSuper"`
	IsDisable  int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
	UpdatedAt  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"updatedAt"`
}
