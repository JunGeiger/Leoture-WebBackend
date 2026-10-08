-- 切换数据库
\connect leoture_web postgres;

-- Schema 权限（public 默认存在）
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE, CREATE ON SCHEMA public TO leoture_web;

-- 未来对象权限（仅限自己创建的）
ALTER DEFAULT PRIVILEGES FOR ROLE leoture_web IN SCHEMA public GRANT ALL ON TABLES TO leoture_web;
ALTER DEFAULT PRIVILEGES FOR ROLE leoture_web IN SCHEMA public GRANT ALL ON SEQUENCES TO leoture_web;

-- 搜索路径
ALTER ROLE leoture_web SET search_path TO public;

\connect leoture_web gogs;

-- Schema 权限（public 默认存在）
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE, CREATE ON SCHEMA public TO gogs;

-- 未来对象权限（仅限自己创建的）
ALTER DEFAULT PRIVILEGES FOR ROLE gogs IN SCHEMA public GRANT ALL ON TABLES TO gogs;
ALTER DEFAULT PRIVILEGES FOR ROLE gogs IN SCHEMA public GRANT ALL ON SEQUENCES TO gogs;

-- 搜索路径
ALTER ROLE gogs SET search_path TO public;
