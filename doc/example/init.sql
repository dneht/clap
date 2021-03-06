create
database if not exists clap default character set utf8mb4;
use
clap;


drop table if exists bootstrap;
create table bootstrap
(
    id         bigint unsigned not null auto_increment,
    env        varchar(128) not null comment '环境',
    prop       varchar(128) not null comment '属性名',
    value      varchar(512) not null comment '属性值',
    is_disable boolean               default false comment '是否已被禁用',
    created_at timestamp    not null default current_timestamp comment '添加时间',
    updated_at timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    unique uk_bootstrap_key (env, prop),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '启动信息';


drop table if exists environment;
create table environment
(
    id          bigint unsigned not null auto_increment,
    env         varchar(16) not null comment '环境',
    env_name    varchar(64) not null comment '环境名',
    env_desc    varchar(256) comment '环境描述',
    is_pub      boolean              default false comment '是否开放访问',
    is_sync     boolean              default true comment '是否接收其它环境的同步信息，主要用来推送配置和发布计划等',
    sync_info   json comment '数组，其它环境信息，同步到其它环境时需要',
    deploy_info json comment '部署信息，主要是部署时用到的信息，如cli、git、repo等',
    format_info json comment '规格信息，包含类型对应的仓库、代理、默认启动参数等',
    is_disable  boolean              default false comment '是否已被禁用',
    created_at  timestamp   not null default current_timestamp comment '添加时间',
    updated_at  timestamp   not null default current_timestamp on update current_timestamp comment '更新时间',
    unique uk_env_name (env_name),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '环境信息';


drop table if exists environment_space;
create table environment_space
(
    id         bigint unsigned not null auto_increment,
    env_id     bigint unsigned not null comment '环境，一个环境创建时会添加一个默认space',
    space_name varchar(16) not null comment '空间名',
    space_keep varchar(16) not null comment '空间所处位置，通常是命名空间',
    space_desc varchar(256) comment '描述',
    space_info json comment '提供项目的缺省信息，主要是conf、code、repo',
    is_view    boolean              default false comment '是否仅查看，会展示全部pod',
    is_control boolean              default false comment '是否独占命名空间，独占则deploy后的name不会带上space',
    is_disable boolean              default false comment '是否已被禁用',
    created_at timestamp   not null default current_timestamp comment '添加时间',
    updated_at timestamp   not null default current_timestamp on update current_timestamp comment '更新时间',
    unique uk_space_name (env_id, space_name),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '环境空间';


drop table if exists project;
create table project
(
    id          bigint unsigned not null auto_increment,
    app_key     varchar(32)  not null comment '项目key',
    app_name    varchar(64)  not null comment '项目名',
    app_desc    varchar(256) not null comment '项目描述',
    app_type    int          not null comment '项目类型',
    app_info    json comment '附加信息，包含项目打包，运行等信息',
    source_info json comment '加密信息，包含资源、密钥信息等，secret应该存放在不同的地方',
    inject_info json comment '注入信息，包含运行时注入信息、如收集日志、链路追踪等',
    is_ingress  boolean               default true comment '是否允许进入执行命令',
    is_disable  boolean               default false comment '是否已被禁用',
    created_at  timestamp    not null default current_timestamp comment '添加时间',
    updated_at  timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    unique uk_app_key (app_key),
    unique uk_app_name (app_name),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '项目信息';


drop table if exists template;
create table template
(
    id               bigint unsigned not null auto_increment,
    template_name    varchar(16)  not null comment '模版名',
    template_kind    varchar(16)  not null comment '模版类型',
    template_desc    varchar(256) not null comment '模版描述',
    template_content text comment '模版内容，目前只能是jsonnet',
    is_disable       boolean               default false comment '是否已被禁用',
    created_at       timestamp    not null default current_timestamp comment '添加时间',
    updated_at       timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    unique uk_template_name (template_name),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '模版信息';


drop table if exists deployment;
create table deployment
(
    id            bigint unsigned not null auto_increment,
    app_id        bigint unsigned not null comment '项目',
    env_id        bigint unsigned not null comment '环境',
    space_id      bigint unsigned not null comment '环境空间',
    branch_name   varchar(64) comment '代码分支',
    deploy_name   varchar(64) not null comment '部署名',
    deploy_status tinyint              default 0 comment '部署状态，修改需要加锁。0默认、1打包中、2打包完成、3打包失败、6已发布',
    deploy_tag    varchar(24) comment '记录当前或者上次一打包使用的tag',
    app_info      json comment '创建部署时覆盖原始的项目信息',
    is_package    boolean              default true comment '是否能打包，默认能',
    is_disable    boolean              default false comment '是否已被禁用',
    created_at    timestamp   not null default current_timestamp comment '添加时间',
    updated_at    timestamp   not null default current_timestamp on update current_timestamp comment '更新时间',
    index         idx_deploy_app (app_id),
    index         idx_deploy_env (env_id),
    index         idx_deploy_space (space_id),
    index         idx_deploy_plan (plan_id),
    unique uk_deploy_name (space_id, deploy_name),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '发布信息';


drop table if exists deployment_plan;
create table deployment_plan
(
    id            bigint unsigned not null auto_increment,
    env_id        bigint unsigned not null comment '环境',
    user_id       bigint unsigned not null comment '创建者',
    deploy_id     bigint unsigned not null comment '部署id',
    deploy_status tinyint            default 0 comment '部署状态，0未知、1成功、2等待、9失败',
    is_disable    boolean            default false comment '是否已被禁用',
    created_at    timestamp not null default current_timestamp comment '添加时间',
    index         idx_user_id (user_id),
    primary key (id, deploy_id)
) engine = innodb
  default charset = utf8mb4 comment = '发布计划';


drop table if exists deployment_log;
create table deployment_log
(
    id            bigint unsigned not null auto_increment,
    app_id        bigint unsigned not null comment '项目',
    env_id        bigint unsigned not null comment '环境',
    space_id      bigint unsigned not null comment '环境空间',
    deploy_id     bigint unsigned not null comment '部署',
    plan_id       bigint unsigned comment '关联的发布计划id，可以为空',
    prop_ids      json comment '关联的配置快照id列表，可以为空',
    branch_name   varchar(64) comment '代码分支',
    deploy_tag    varchar(24) comment '打包使用的tag',
    snapshot_info json comment '部署时的信息快照，合并后的信息',
    created_at    timestamp not null default current_timestamp comment '添加时间',
    index         idx_app_id (env_id, app_id),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '发布日志';


drop table if exists property_file;
create table property_file
(
    id           bigint unsigned not null auto_increment,
    res_id       bigint unsigned not null comment '资源id',
    link_id      bigint unsigned not null comment '关联id',
    file_name    varchar(64)  not null comment '文件名，不包含文件路径',
    file_readme  varchar(256) not null comment '配置文件说明',
    file_content text         not null comment '配置文件文本',
    file_hash    varchar(64)  not null comment '根据file_content计算的hash',
    is_disable   boolean               default false comment '是否已被禁用',
    created_at   timestamp    not null default current_timestamp comment '添加时间',
    updated_at   timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    unique uk_link_res_id (link_id, res_id, file_name),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '配置文件';


drop table if exists property_snap;
create table property_snap
(
    id           bigint unsigned not null auto_increment,
    res_id       bigint unsigned not null comment '资源id',
    link_id      bigint unsigned not null comment '关联id',
    prop_id      bigint unsigned not null comment '配置id',
    file_name    varchar(64) not null comment '文件名，不包含文件路径',
    file_content text        not null comment '配置文件文本',
    created_at   timestamp   not null default current_timestamp comment '添加时间',
    index        idx_link_res_id (link_id, res_id, file_name),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '配置快照';


drop table if exists resource;
create table resource
(
    id         bigint unsigned not null auto_increment,
    sys_id     int unsigned not null default 0 comment '系统id，默认0本系统',
    res_name   varchar(128) not null comment '资源名',
    res_order  int                   default 0 comment '资源排序，在同一个parent_id下有效',
    res_info   json comment '资源附加信息',
    created_at timestamp    not null default current_timestamp comment '添加时间',
    updated_at timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    unique uk_res_name (res_name),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '资源信息';


drop table if exists permission;
create table permission
(
    id         bigint unsigned not null auto_increment,
    sys_id     int unsigned not null default 0 comment '系统id，默认0本系统',
    role_id    bigint unsigned not null comment '角色id',
    res_id     bigint unsigned not null comment '资源id',
    res_power  int unsigned comment '二进制表示，从右到左的二进制位表示select，update、insert、delete、grant x4',
    link_id    bigint unsigned default 0 comment '关联id',
    power_info json comment '权限附加信息',
    created_at timestamp not null default current_timestamp comment '添加时间',
    updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
    index      idx_role_id (role_id),
    index      idx_res_link_id (res_id, link_id),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '权限信息';


drop table if exists role_info;
create table role_info
(
    id          bigint unsigned not null auto_increment,
    role_name   varchar(64)  not null comment '角色名',
    role_from   int unsigned default 0 comment '角色来源、本系统0，自动创建1',
    role_remark varchar(256) not null comment '备注信息',
    is_manage   boolean               default false comment '是否是管理角色',
    is_super    boolean               default false comment '是否是超级管理角色',
    is_disable  boolean               default false comment '是否已被禁用',
    created_at  timestamp    not null default current_timestamp comment '添加时间',
    updated_at  timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    unique uk_role_from_name (role_from, role_name),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '角色信息';


drop table if exists user_info;
create table user_info
(
    id           bigint unsigned not null auto_increment,
    user_name    varchar(64)  not null comment '用户名',
    user_from    int unsigned default 0 comment '用户来源、本系统0',
    nickname     varchar(128) not null comment '昵称',
    password     varchar(128) not null comment '密码',
    avatar       varchar(256) not null comment '用户头像',
    access_token varchar(128) not null comment '访问token',
    role_list    json comment '用户加入的角色',
    is_disable   boolean               default false comment '是否已被禁用',
    created_at   timestamp    not null default current_timestamp comment '添加时间',
    updated_at   timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    unique uk_user_from_name (user_from, user_name),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '用户信息';


drop table if exists operation_log;
create table operation_log
(
    id           bigint unsigned not null auto_increment,
    user_id      bigint unsigned not null comment '用户id',
    res_id       bigint unsigned not null comment '资源id',
    log_type     int       not null comment '操作类型',
    log_info     json comment '具体内容',
    request_from json      not null comment '来源信息，如ip、method、path等',
    created_at   timestamp not null default current_timestamp comment '添加时间',
    index        idx_user_id (user_id),
    index        idx_res_id (res_id),
    primary key (id)
) engine = innodb
  default charset = utf8mb4 comment = '操作日志';
