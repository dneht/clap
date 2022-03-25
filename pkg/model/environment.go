package model

import (
	"time"
)

const EnvironmentTable = "environment"

const (
	IdInEnvironment         = "id"
	EnvInEnvironment        = "env"
	EnvNameInEnvironment    = "env_name"
	EnvDescInEnvironment    = "env_desc"
	IsPubInEnvironment      = "is_pub"
	IsSyncInEnvironment     = "is_sync"
	SyncInfoInEnvironment   = "sync_info"
	DeployInfoInEnvironment = "deploy_info"
	FormatInfoInEnvironment = "format_info"
	IsDisableInEnvironment  = "is_disable"
	CreatedAtInEnvironment  = "created_at"
	UpdatedAtInEnvironment  = "updated_at"
)

type Environment struct {
	Id         uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	Env        string    `xorm:"not null comment('环境') VARCHAR(16)" json:"env"`
	EnvName    string    `xorm:"not null comment('环境名') unique VARCHAR(64)" json:"envName"`
	EnvDesc    string    `xorm:"comment('环境描述') VARCHAR(256)" json:"envDesc"`
	IsPub      int       `xorm:"default 0 comment('是否开放访问') TINYINT(1)" json:"isPub"`
	IsSync     int       `xorm:"default 1 comment('是否接收其它环境的同步信息，主要用来推送配置和发布计划等') TINYINT(1)" json:"isSync"`
	SyncInfo   string    `xorm:"comment('数组，其它环境信息，同步到其它环境时需要') JSON" json:"syncInfo"`
	DeployInfo string    `xorm:"comment('部署信息，主要是部署时用到的信息，如cli、git、repo等') JSON" json:"deployInfo"`
	FormatInfo string    `xorm:"comment('规格信息，包含类型对应的仓库、代理、默认启动参数等') JSON" json:"formatInfo"`
	IsDisable  int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
	UpdatedAt  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"updatedAt"`
}
