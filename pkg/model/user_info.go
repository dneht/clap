package model

import (
	"time"
)

const UserInfoTable = "user_info"

const (
	IdInUserInfo          = "id"
	UserNameInUserInfo    = "user_name"
	UserFromInUserInfo    = "user_from"
	NicknameInUserInfo    = "nickname"
	PasswordInUserInfo    = "password"
	AvatarInUserInfo      = "avatar"
	AccessTokenInUserInfo = "access_token"
	RoleListInUserInfo    = "role_list"
	IsDisableInUserInfo   = "is_disable"
	CreatedAtInUserInfo   = "created_at"
	UpdatedAtInUserInfo   = "updated_at"
)

type UserInfo struct {
	Id          uint64    `xorm:"not null pk autoincr BIGINT(20)" json:"id"`
	UserName    string    `xorm:"not null comment('用户名') unique(uk_user_from_name) VARCHAR(64)" json:"userName"`
	UserFrom    uint      `xorm:"default 0 comment('用户来源、本系统0') unique(uk_user_from_name) INT(10)" json:"userFrom"`
	Nickname    string    `xorm:"not null comment('昵称') VARCHAR(128)" json:"nickname"`
	Password    string    `xorm:"not null comment('密码') VARCHAR(128)" json:"password"`
	Avatar      string    `xorm:"not null comment('用户头像') VARCHAR(256)" json:"avatar"`
	AccessToken string    `xorm:"not null comment('访问token') VARCHAR(128)" json:"accessToken"`
	RoleList    string    `xorm:"comment('用户加入的角色') JSON" json:"roleList"`
	IsDisable   int       `xorm:"default 0 comment('是否已被禁用') TINYINT(1)" json:"isDisable"`
	CreatedAt   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('添加时间') TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') TIMESTAMP" json:"updatedAt"`
}
