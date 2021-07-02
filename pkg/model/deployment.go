package model

import (
	"time"
)

const DeploymentTable = "deployment"

const (
	IdInDeployment           = "id"
	AppIdInDeployment        = "app_id"
	EnvIdInDeployment        = "env_id"
	SpaceIdInDeployment      = "space_id"
	PlanIdInDeployment       = "plan_id"
	BranchNameInDeployment   = "branch_name"
	DeployNameInDeployment   = "deploy_name"
	DeployStatusInDeployment = "deploy_status"
	DeployTagInDeployment    = "deploy_tag"
	AppInfoInDeployment      = "app_info"
	IsPackageInDeployment    = "is_package"
	IsDisableInDeployment    = "is_disable"
	CreatedAtInDeployment    = "created_at"
	UpdatedAtInDeployment    = "updated_at"
)

type Deployment struct {
	Id           uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	AppId        uint64    `xorm:"not null comment('项目') index BIGINT(20)" json:"appId"`
	EnvId        uint64    `xorm:"not null comment('环境') index BIGINT(20)" json:"envId"`
	SpaceId      uint64    `xorm:"not null comment('环境空间') index unique(uk_deploy_name) BIGINT(20)" json:"spaceId"`
	PlanId       int64     `xorm:"default 0 comment('计划id，关联使用，如果不为0则包含在发布计划中') index BIGINT(20)" json:"planId"`
	BranchName   string    `xorm:"comment('代码分支') VARCHAR(64)" json:"branchName"`
	DeployName   string    `xorm:"not null comment('部署名') unique(uk_deploy_name) VARCHAR(64)" json:"deployName"`
	DeployStatus int       `xorm:"default 0 comment('部署状态，修改需要加锁。0默认、1打包中、2打包完成、3打包失败、6已发布') TINYINT(4)" json:"deployStatus"`
	DeployTag    string    `xorm:"comment('记录当前或者上次一打包使用的tag') VARCHAR(24)" json:"deployTag"`
	AppInfo      string    `xorm:"comment('创建部署时覆盖原始的项目信息') JSON" json:"appInfo"`
	IsPackage    int       `xorm:"default 1 comment('是否能打包，默认能') TINYINT(1)" json:"isPackage"`
	IsDisable    int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt    time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') TIMESTAMP" json:"createdAt"`
	UpdatedAt    time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') TIMESTAMP" json:"updatedAt"`
}
