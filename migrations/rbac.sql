-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(64),
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(24) UNIQUE,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE  users IS '用户表';
COMMENT ON COLUMN users.id IS '用户ID';
COMMENT ON COLUMN users.username IS '用户名称';
COMMENT ON COLUMN users.password IS '密码';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.email IS '邮箱';
COMMENT ON COLUMN users.phone IS '手机号';
COMMENT ON COLUMN users.status IS '状态，1:正常, 2:禁用';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';
COMMENT ON COLUMN users.deleted_at IS '删除时间，默认为空，非空则说明已经被软删除';

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE,
    code VARCHAR(32) NOT NULL UNIQUE,
    remark TEXT,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE  roles IS '角色表';
COMMENT ON COLUMN roles.id IS '角色ID';
COMMENT ON COLUMN roles.name IS '角色名称';
COMMENT ON COLUMN roles.code IS '角色编码';
COMMENT ON COLUMN roles.remark IS '备注';
COMMENT ON COLUMN roles.status IS '状态，1:正常, 2:禁用';
COMMENT ON COLUMN roles.created_at IS '创建时间';
COMMENT ON COLUMN roles.updated_at IS '更新时间';
COMMENT ON COLUMN roles.deleted_at IS '删除时间，默认为空，非空则说明已经被软删除';

-- 菜单表
CREATE TABLE IF NOT EXISTS menus (
    id UUID PRIMARY KEY,
    parent_id UUID,
    type SMALLINT NOT NULL, -- 1:目录 2:菜单 3:按钮
    name VARCHAR(32) NOT NULL UNIQUE,
    title VARCHAR(64) NOT NULL,
    path VARCHAR(255),
    component VARCHAR(255),
    icon VARCHAR(32),
    keep_alive BOOLEAN NOT NULL DEFAULT FALSE,
    order_no INT NOT NULL DEFAULT 0,
    external_url VARCHAR(255),
    is_hidden BOOLEAN NOT NULL DEFAULT TRUE,
    perm_code VARCHAR(32) UNIQUE,
    api_method VARCHAR(8),
    api_path VARCHAR(255),
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE (api_method, api_path)
);
COMMENT ON TABLE  menus IS '菜单表';
COMMENT ON COLUMN menus.id IS '菜单ID';
COMMENT ON COLUMN menus.parent_id IS '父菜单ID，为NULL则说明是根节点';
COMMENT ON COLUMN menus.type IS '类型：1-目录, 2-菜单, 3-按钮';
COMMENT ON COLUMN menus.name IS '名称';
COMMENT ON COLUMN menus.title IS '标题';
COMMENT ON COLUMN menus.path IS '前端路由地址';
COMMENT ON COLUMN menus.component IS '前端路由组件';
COMMENT ON COLUMN menus.icon IS '图标';
COMMENT ON COLUMN menus.keep_alive IS '是否开启页面缓存';
COMMENT ON COLUMN menus.order_no IS '排序';
COMMENT ON COLUMN menus.external_url IS '外部链接';
COMMENT ON COLUMN menus.is_hidden IS '是否隐藏菜单';
COMMENT ON COLUMN menus.perm_code IS '权限标识，如 system:user:add，按钮级权限必填';
COMMENT ON COLUMN menus.api_method IS 'HTTP方法';
COMMENT ON COLUMN menus.api_path IS 'API地址';
COMMENT ON COLUMN menus.status IS '状态：1:正常, 2:禁用';
COMMENT ON COLUMN menus.created_at IS '创建时间';
COMMENT ON COLUMN menus.updated_at IS '更新时间';
COMMENT ON COLUMN menus.deleted_at IS '删除时间，默认为空，非空则说明已经被软删除';

-- 用户-角色关联表
CREATE TABLE IF NOT EXISTS user_role (
    user_id UUID NOT NULL REFERENCES users(id),
    role_id UUID NOT NULL REFERENCES roles(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id)
);
COMMENT ON TABLE  user_role IS '用户角色关联表';
COMMENT ON COLUMN user_role.user_id IS '用户ID';
COMMENT ON COLUMN user_role.role_id IS '角色ID';
COMMENT ON COLUMN user_role.created_at IS '创建时间';

-- 角色-菜单关联表
CREATE TABLE IF NOT EXISTS role_menu (
    role_id UUID NOT NULL REFERENCES roles(id),
    menu_id UUID NOT NULL REFERENCES menus(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, menu_id)
);
COMMENT ON TABLE  role_menu IS '角色菜单关联表';
COMMENT ON COLUMN role_menu.role_id IS '角色ID';
COMMENT ON COLUMN role_menu.menu_id IS '菜单ID';
COMMENT ON COLUMN role_menu.created_at IS '创建时间';

CREATE TABLE IF NOT EXISTS dictionaries (
    id UUID PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    category VARCHAR(32) NOT NULL,
    key VARCHAR(255) NOT NULL,
    val VARCHAR(255) NOT NULL,
    order_no INT NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    remark TEXT,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE (category, key)
);
COMMENT ON TABLE  dictionaries IS '字典表';
COMMENT ON COLUMN dictionaries.id IS '字典ID';
COMMENT ON COLUMN dictionaries.name IS '名称';
COMMENT ON COLUMN dictionaries.category IS '类别';
COMMENT ON COLUMN dictionaries.key IS '字典key';
COMMENT ON COLUMN dictionaries.val IS '字典value';
COMMENT ON COLUMN dictionaries.order_no IS '排序';
COMMENT ON COLUMN dictionaries.is_default IS '是否选项默认值，true / false';
COMMENT ON COLUMN dictionaries.remark IS '备注';
COMMENT ON COLUMN dictionaries.status IS '状态，1:正常, 2:禁用';
COMMENT ON COLUMN dictionaries.created_at IS '创建时间';
COMMENT ON COLUMN dictionaries.updated_at IS '更新时间';
COMMENT ON COLUMN dictionaries.deleted_at IS '删除时间，默认为空，非空则说明已经被软删除';

CREATE TABLE IF NOT EXISTS casbin_rule (
    id BIGSERIAL PRIMARY KEY,
    ptype VARCHAR(8) NOT NULL,
    v0 VARCHAR(255),
    v1 VARCHAR(255),
    v2 VARCHAR(255),
    v3 VARCHAR(255),
    v4 VARCHAR(255),
    v5 VARCHAR(255),
    UNIQUE (ptype, v0, v1, v3, v4, v5)
);
COMMENT ON TABLE  casbin_rule IS 'Casbin权限表';
COMMENT ON COLUMN casbin_rule.id IS '自动递增ID字段';
COMMENT ON COLUMN casbin_rule.ptype IS '策略类型:“p”(策略)或“g”(分组)';
COMMENT ON COLUMN casbin_rule.v0 IS 'Subject (username / roleCode)';
COMMENT ON COLUMN casbin_rule.v1 IS 'Object  (roleCode / url)';
COMMENT ON COLUMN casbin_rule.v2 IS 'Action  (GET/PUT/POST/DELETE/OPTIONS)';
COMMENT ON COLUMN casbin_rule.v3 IS 'Effect  (allow/deny)';
COMMENT ON COLUMN casbin_rule.v4 IS '可选字段';
COMMENT ON COLUMN casbin_rule.v5 IS '可选字段';

CREATE TABLE IF NOT EXISTS login_log (
    id UUID PRIMARY KEY,
    username VARCHAR(32) NOT NULL,
    ip_address INET NOT NULL,
    user_agent VARCHAR(512) NOT NULL DEFAULT '',
    origin_url VARCHAR(255),
    location VARCHAR(255),
    result VARCHAR(512) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE  login_log IS '用户登录日志表';
COMMENT ON COLUMN login_log.id IS '自动递增ID字段';
COMMENT ON COLUMN login_log.username IS '登录用户名称';
COMMENT ON COLUMN login_log.ip_address IS '客户端IP地址';
COMMENT ON COLUMN login_log.user_agent IS '客户端浏览器标识';
COMMENT ON COLUMN login_log.origin_url IS '登录请求来源页面URL';
COMMENT ON COLUMN login_log.location IS '用户自报地点信息(可选)';
COMMENT ON COLUMN login_log.result IS '操作结果：成功/失败';
COMMENT ON COLUMN login_log.created_at IS '登录请求日志创建时间';

CREATE TABLE IF NOT EXISTS user_session (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    refresh_jti UUID NOT NULL,
    login_log_id UUID NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE  user_session IS '用户会话Session表';
COMMENT ON COLUMN user_session.id IS '会话唯一标识ID';
COMMENT ON COLUMN user_session.user_id IS '关联用户ID';
COMMENT ON COLUMN user_session.refresh_jti IS 'Refresh Token唯一标识';
COMMENT ON COLUMN user_session.login_log_id IS '关联登录日志ID';
COMMENT ON COLUMN user_session.status IS '会话状态: 1-活跃 2-主动注销 3-已被替换(已有新登录) 4-风控/封禁';
COMMENT ON COLUMN user_session.created_at IS '会话创建时间';
COMMENT ON COLUMN user_session.updated_at IS '会话最后更新时间';
