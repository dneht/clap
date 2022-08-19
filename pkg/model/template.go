package model

import (
	"time"
)

const TemplateTable = "template"

const (
	IdInTemplate              = "id"
	TemplateNameInTemplate    = "template_name"
	TemplateKindInTemplate    = "template_kind"
	TemplateDescInTemplate    = "template_desc"
	TemplateContentInTemplate = "template_content"
	IsDisableInTemplate       = "is_disable"
	CreatedAtInTemplate       = "created_at"
	UpdatedAtInTemplate       = "updated_at"
)

type Template struct {
	Id              uint64    `xorm:"not null pk autoincr BIGINT" json:"id"`
	TemplateName    string    `xorm:"not null comment('模版名') unique VARCHAR(16)" json:"templateName"`
	TemplateKind    string    `xorm:"not null comment('模版类型') VARCHAR(16)" json:"templateKind"`
	TemplateDesc    string    `xorm:"not null comment('模版描述') VARCHAR(256)" json:"templateDesc"`
	TemplateContent string    `xorm:"comment('模版内容，目前只能是jsonnet') TEXT" json:"templateContent"`
	IsDisable       int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt       time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
	UpdatedAt       time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"updatedAt"`
}
