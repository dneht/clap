package model

import (
	"time"
)

const DeploymentLogTable = "deployment_log"

const (
	IdInDeploymentLog           = "id"
	AppIdInDeploymentLog        = "app_id"
	EnvIdInDeploymentLog        = "env_id"
	SpaceIdInDeploymentLog      = "space_id"
	DeployIdInDeploymentLog     = "deploy_id"
	PlanIdInDeploymentLog       = "plan_id"
	DeployTimeInDeploymentLog   = "deploy_time"
	DeployStatusInDeploymentLog = "deploy_status"
	SnapshotInfoInDeploymentLog = "snapshot_info"
	CreatedAtInDeploymentLog    = "created_at"
)

type DeploymentLog struct {
	Id           uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	AppId        uint64    `xorm:"not null comment('项目') index(idx_app_id) BIGINT(20)" json:"appId"`
	EnvId        uint64    `xorm:"not null comment('环境') index(idx_app_id) BIGINT(20)" json:"envId"`
	SpaceId      uint64    `xorm:"not null comment('环境空间') BIGINT(20)" json:"spaceId"`
	DeployId     uint64    `xorm:"not null comment('部署') BIGINT(20)" json:"deployId"`
	PlanId       uint64    `xorm:"comment('关联的发布计划id，可以为空') BIGINT(20)" json:"planId"`
	DeployTime   uint      `xorm:"comment('部署用时，单位秒') INT(10)" json:"deployTime"`
	DeployStatus int       `xorm:"default 0 comment('部署状态，0未知、1成功、2等待、9失败') TINYINT(4)" json:"deployStatus"`
	SnapshotInfo string    `xorm:"comment('部署时项目、打包及生成的信息快照') JSON" json:"snapshotInfo"`
	CreatedAt    time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') TIMESTAMP" json:"createdAt"`
}
