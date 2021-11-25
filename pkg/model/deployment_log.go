package model

import (
	"time"
)

const DeploymentLogTable = "deployment_log"

const (
	IdInDeploymentLog           = "id"
	UserIdInDeploymentLog       = "user_id"
	AppIdInDeploymentLog        = "app_id"
	EnvIdInDeploymentLog        = "env_id"
	SpaceIdInDeploymentLog      = "space_id"
	DeployIdInDeploymentLog     = "deploy_id"
	PlanIdInDeploymentLog       = "plan_id"
	PropIdsInDeploymentLog      = "prop_ids"
	BranchNameInDeploymentLog   = "branch_name"
	DeployTagInDeploymentLog    = "deploy_tag"
	SnapshotInfoInDeploymentLog = "snapshot_info"
	CreatedAtInDeploymentLog    = "created_at"
)

type DeploymentLog struct {
	Id           uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	UserId       uint64    `xorm:"not null comment('用户id') BIGINT(20)" json:"userId"`
	AppId        uint64    `xorm:"not null comment('项目id') index(idx_app_id) BIGINT(20)" json:"appId"`
	EnvId        uint64    `xorm:"not null comment('环境id') index(idx_app_id) BIGINT(20)" json:"envId"`
	SpaceId      uint64    `xorm:"not null comment('空间id') BIGINT(20)" json:"spaceId"`
	DeployId     uint64    `xorm:"not null comment('部署id') BIGINT(20)" json:"deployId"`
	PlanId       uint64    `xorm:"comment('关联的发布计划id，可以为空') BIGINT(20)" json:"planId"`
	PropIds      string    `xorm:"comment('关联的配置快照id列表，可以为空') JSON" json:"propIds"`
	BranchName   string    `xorm:"comment('代码分支') VARCHAR(64)" json:"branchName"`
	DeployTag    string    `xorm:"comment('打包使用的tag') VARCHAR(24)" json:"deployTag"`
	SnapshotInfo string    `xorm:"comment('部署时的信息快照，合并后的信息') JSON" json:"snapshotInfo"`
	CreatedAt    time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') TIMESTAMP" json:"createdAt"`
}
