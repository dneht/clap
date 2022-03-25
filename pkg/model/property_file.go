package model

import (
	"time"
)

const PropertyFileTable = "property_file"

const (
	IdInPropertyFile          = "id"
	ResIdInPropertyFile       = "res_id"
	LinkIdInPropertyFile      = "link_id"
	FileNameInPropertyFile    = "file_name"
	FileReadmeInPropertyFile  = "file_readme"
	FileContentInPropertyFile = "file_content"
	FileHashInPropertyFile    = "file_hash"
	IsDisableInPropertyFile   = "is_disable"
	CreatedAtInPropertyFile   = "created_at"
	UpdatedAtInPropertyFile   = "updated_at"
)

type PropertyFile struct {
	Id          uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	ResId       uint64    `xorm:"not null comment('资源id') unique(uk_link_res_id) BIGINT(20)" json:"resId"`
	LinkId      uint64    `xorm:"not null comment('关联id') unique(uk_link_res_id) BIGINT(20)" json:"linkId"`
	FileName    string    `xorm:"not null comment('文件名，不包含文件路径') unique(uk_link_res_id) VARCHAR(64)" json:"fileName"`
	FileReadme  string    `xorm:"not null comment('配置文件说明') VARCHAR(256)" json:"fileReadme"`
	FileContent string    `xorm:"not null comment('配置文件文本') TEXT" json:"fileContent"`
	FileHash    string    `xorm:"not null comment('根据file_content计算的hash') VARCHAR(64)" json:"fileHash"`
	IsDisable   int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
	UpdatedAt   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"updatedAt"`
}
