package model

import (
	"time"
)

const PropertySnapTable = "property_snap"

const (
	IdInPropertySnap          = "id"
	UserIdInPropertySnap      = "user_id"
	ResIdInPropertySnap       = "res_id"
	LinkIdInPropertySnap      = "link_id"
	PropIdInPropertySnap      = "prop_id"
	FileNameInPropertySnap    = "file_name"
	FileContentInPropertySnap = "file_content"
	CreatedAtInPropertySnap   = "created_at"
)

type PropertySnap struct {
	Id          uint64    `xorm:"not null pk autoincr BIGINT" json:"id"`
	UserId      uint64    `xorm:"not null comment('用户id') BIGINT" json:"userId"`
	ResId       uint64    `xorm:"not null comment('资源id') index(idx_link_res_id) BIGINT" json:"resId"`
	LinkId      uint64    `xorm:"not null comment('关联id') index(idx_link_res_id) BIGINT" json:"linkId"`
	PropId      uint64    `xorm:"not null comment('配置id') index BIGINT" json:"propId"`
	FileName    string    `xorm:"not null comment('文件名，不包含文件路径') index(idx_link_res_id) VARCHAR(64)" json:"fileName"`
	FileContent string    `xorm:"not null comment('配置文件文本') TEXT" json:"fileContent"`
	CreatedAt   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
}
