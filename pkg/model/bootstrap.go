package model

import (
	"time"
)

const BootstrapTable = "bootstrap"

const (
	IdInBootstrap        = "id"
	EnvInBootstrap       = "env"
	PropInBootstrap      = "prop"
	ValueInBootstrap     = "value"
	IsDisableInBootstrap = "is_disable"
	CreatedAtInBootstrap = "created_at"
	UpdatedAtInBootstrap = "updated_at"
)

type Bootstrap struct {
	Id        uint64    `xorm:"not null pk autoincr BIGINT" json:"id"`
	Env       string    `xorm:"not null comment('环境') unique(uk_bootstrap_key) VARCHAR(128)" json:"env"`
	Prop      string    `xorm:"not null comment('属性名') unique(uk_bootstrap_key) VARCHAR(128)" json:"prop"`
	Value     string    `xorm:"not null comment('属性值') VARCHAR(4096)" json:"value"`
	IsDisable int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
	UpdatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"updatedAt"`
}
