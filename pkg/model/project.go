package model

import (
	"time"
)

const ProjectTable = "project"

const (
	IdInProject         = "id"
	AppKeyInProject     = "app_key"
	AppNameInProject    = "app_name"
	AppDescInProject    = "app_desc"
	AppTypeInProject    = "app_type"
	AppInfoInProject    = "app_info"
	SourceInfoInProject = "source_info"
	InjectInfoInProject = "inject_info"
	IsIngressInProject  = "is_ingress"
	IsDisableInProject  = "is_disable"
	CreatedAtInProject  = "created_at"
	UpdatedAtInProject  = "updated_at"
)

type Project struct {
	Id         uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	AppKey     string    `xorm:"not null comment('项目key') unique VARCHAR(32)" json:"appKey"`
	AppName    string    `xorm:"not null comment('项目名') unique VARCHAR(64)" json:"appName"`
	AppDesc    string    `xorm:"not null comment('项目描述') VARCHAR(256)" json:"appDesc"`
	AppType    int       `xorm:"not null comment('项目类型，5nginx、10java、11tomcat、20go、60python、90node') INT(11)" json:"appType"`
	AppInfo    string    `xorm:"comment('附加信息，包含项目打包，运行等信息') JSON" json:"appInfo"`
	SourceInfo string    `xorm:"comment('加密信息，包含资源、密钥信息等，secret应该存放在不同的地方') JSON" json:"sourceInfo"`
	InjectInfo string    `xorm:"comment('注入信息，包含运行时注入信息、如收集日志、链路追踪等') JSON" json:"injectInfo"`
	IsIngress  int       `xorm:"default 1 comment('是否允许进入执行命令') TINYINT(1)" json:"isIngress"`
	IsDisable  int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
	UpdatedAt  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"updatedAt"`
}
