-- 1. 创建用户（幂等）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_roles WHERE rolname = 'leoture_web'
    ) THEN
        EXECUTE 'CREATE ROLE leoture_web WITH LOGIN PASSWORD ''leotureWeb_service''';
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_roles WHERE rolname = 'gogs'
    ) THEN
        EXECUTE 'CREATE ROLE gogs WITH LOGIN PASSWORD ''gogs_database''';
    END IF;
END
$$;
