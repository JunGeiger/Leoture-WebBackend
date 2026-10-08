-- 2. 创建数据库（幂等，Docker 正确姿势）
SELECT 'CREATE DATABASE leoture_web OWNER leoture_web'
    WHERE NOT EXISTS (
    SELECT 1 FROM pg_database WHERE datname = 'leoture_web'
)\gexec

SELECT 'CREATE DATABASE gogs OWNER gogs'
    WHERE NOT EXISTS (
    SELECT 1 FROM pg_database WHERE datname = 'gogs'
)\gexec
