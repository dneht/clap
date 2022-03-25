package model

import (
	"time"
)

const EnvironmentSpaceTable = "environment_space"

const (
	IdInEnvironmentSpace        = "id"
	EnvIdInEnvironmentSpace     = "env_id"
	SpaceNameInEnvironmentSpace = "space_name"
	SpaceKeepInEnvironmentSpace = "space_keep"
	SpaceDescInEnvironmentSpace = "space_desc"
	SpaceInfoInEnvironmentSpace = "space_info"
	IsViewInEnvironmentSpace    = "is_view"
	IsControlInEnvironmentSpace = "is_control"
	IsDisableInEnvironmentSpace = "is_disable"
	CreatedAtInEnvironmentSpace = "created_at"
	UpdatedAtInEnvironmentSpace = "updated_at"
)

type EnvironmentSpace struct {
	Id        uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	EnvId     uint64    `xorm:"not null comment('环境，一个环境创建时会添加一个默认space') unique(uk_space_name) BIGINT(20)" json:"envId"`
	SpaceName string    `xorm:"not null comment('空间名') unique(uk_space_name) VARCHAR(16)" json:"spaceName"`
	SpaceKeep string    `xorm:"not null comment('空间所处位置，通常是命名空间') VARCHAR(16)" json:"spaceKeep"`
	SpaceDesc string    `xorm:"comment('描述') VARCHAR(256)" json:"spaceDesc"`
	SpaceInfo string    `xorm:"comment('提供项目的缺省信息，主要是conf、code、repo') JSON" json:"spaceInfo"`
	IsView    int       `xorm:"default 0 comment('是否仅查看，会展示全部pod') TINYINT(1)" json:"isView"`
	IsControl int       `xorm:"default 0 comment('是否独占命名空间，独占则deploy后的name不会带上space') TINYINT(1)" json:"isControl"`
	IsDisable int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
	UpdatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"updatedAt"`
}
