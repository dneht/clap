package model

import (
	"time"
)

const PropertyTable = "property"

const (
	IdInProperty        = "id"
	EnvInProperty       = "env"
	PropInProperty      = "prop"
	ValueInProperty     = "value"
	IsDisableInProperty = "is_disable"
	CreatedAtInProperty = "created_at"
	UpdatedAtInProperty = "updated_at"
)

type Property struct {
	Id        uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	Env       string    `xorm:"not null comment('环境') unique(uk_property_key) VARCHAR(128)" json:"env"`
	Prop      string    `xorm:"not null comment('属性名') unique(uk_property_key) VARCHAR(128)" json:"prop"`
	Value     string    `xorm:"not null comment('属性值') VARCHAR(512)" json:"value"`
	IsDisable int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') TIMESTAMP" json:"updatedAt"`
}
