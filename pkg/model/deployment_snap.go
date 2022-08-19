package model

import (
	"time"
)

const DeploymentSnapTable = "deployment_snap"

const (
	IdInDeploymentSnap           = "id"
	UserIdInDeploymentSnap       = "user_id"
	AppIdInDeploymentSnap        = "app_id"
	EnvIdInDeploymentSnap        = "env_id"
	SpaceIdInDeploymentSnap      = "space_id"
	DeployIdInDeploymentSnap     = "deploy_id"
	FlowIdInDeploymentSnap       = "flow_id"
	BranchNameInDeploymentSnap   = "branch_name"
	DeployStatusInDeploymentSnap = "deploy_status"
	DeployNameInDeploymentSnap   = "deploy_name"
	DeployKindInDeploymentSnap   = "deploy_kind"
	DeployTagInDeploymentSnap    = "deploy_tag"
	DeployRenderInDeploymentSnap = "deploy_render"
	CreatedAtInDeploymentSnap    = "created_at"
)

type DeploymentSnap struct {
	Id           uint64    `xorm:"not null pk autoincr BIGINT" json:"id"`
	UserId       uint64    `xorm:"not null comment('用户id') BIGINT" json:"userId"`
	AppId        uint64    `xorm:"not null comment('项目id') index(idx_app_id) BIGINT" json:"appId"`
	EnvId        uint64    `xorm:"not null comment('环境id') index(idx_app_id) BIGINT" json:"envId"`
	SpaceId      uint64    `xorm:"not null comment('空间id') BIGINT" json:"spaceId"`
	DeployId     uint64    `xorm:"not null comment('部署id') index BIGINT" json:"deployId"`
	FlowId       uint64    `xorm:"default 0 comment('关联的流程id，可以为空') index BIGINT" json:"flowId"`
	BranchName   string    `xorm:"comment('代码分支') VARCHAR(64)" json:"branchName"`
	DeployStatus int       `xorm:"default -1 comment('部署状态，flow_id=0时即直接使用为-1，否则同主表') TINYINT" json:"deployStatus"`
	DeployName   string    `xorm:"comment('打包使用的name') VARCHAR(128)" json:"deployName"`
	DeployKind   string    `xorm:"comment('打包使用的kind') VARCHAR(128)" json:"deployKind"`
	DeployTag    string    `xorm:"comment('打包使用的tag') VARCHAR(24)" json:"deployTag"`
	DeployRender string    `xorm:"comment('部署时的信息快照，合并后的信息') JSON" json:"deployRender"`
	CreatedAt    time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') DATETIME" json:"createdAt"`
}
