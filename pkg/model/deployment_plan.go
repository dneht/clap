package model

import (
	"time"
)

const DeploymentPlanTable = "deployment_plan"

const (
	IdInDeploymentPlan           = "id"
	EnvIdInDeploymentPlan        = "env_id"
	UserIdInDeploymentPlan       = "user_id"
	DeployListInDeploymentPlan   = "deploy_list"
	DeployStatusInDeploymentPlan = "deploy_status"
	IsDisableInDeploymentPlan    = "is_disable"
	CreatedAtInDeploymentPlan    = "created_at"
)

type DeploymentPlan struct {
	Id           uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	EnvId        uint64    `xorm:"not null comment('环境') index BIGINT(20)" json:"envId"`
	UserId       uint64    `xorm:"not null comment('创建者') BIGINT(20)" json:"userId"`
	DeployList   string    `xorm:"comment('部署列表，里面是deploy_id，只能在同一个env下') JSON" json:"deployList"`
	DeployStatus int       `xorm:"default 0 comment('部署状态，0未知、1成功、2等待、9失败') TINYINT(4)" json:"deployStatus"`
	IsDisable    int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt    time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') TIMESTAMP" json:"createdAt"`
}
