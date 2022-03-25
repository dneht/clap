package model

import (
	"time"
)

const ResourceTable = "resource"

const (
	IdInResource        = "id"
	ResNameInResource   = "res_name"
	ResOrderInResource  = "res_order"
	ResInfoInResource   = "res_info"
	CreatedAtInResource = "created_at"
	UpdatedAtInResource = "updated_at"
)

type Resource struct {
	Id        uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	ResName   string    `xorm:"not null comment('资源名') unique VARCHAR(128)" json:"resName"`
	ResOrder  int       `xorm:"default 0 comment('资源排序，在同一个parent_id下有效') INT(11)" json:"resOrder"`
	ResInfo   string    `xorm:"comment('资源附加信息') JSON" json:"resInfo"`
	CreatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
	UpdatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"updatedAt"`
}
