INSERT INTO users (id, username, password, nickname, email, phone)
VALUES ('019f8d61-b8c2-73f9-8623-6907dc275fce'::uuid, 'leoture', '$2a$10$m5brK.a096htnVxzd2YdHe8irs1KAVe.6yKjzW1lqIgmbFM3q90yO', 'GoWeb', '123456789@leoture.com', '18888888888');

INSERT INTO roles (id, name, code, remark)
VALUES ('019fcc22-31fd-7559-bbfc-8c6d5079557d'::uuid, '通用', 'common', '通用角色，系统运行必要');
INSERT INTO roles (id, name, code, remark)
VALUES ('019f8d61-b8c2-73f9-8623-6cb937ccb8ee'::uuid, '管理员', 'admin', '管理员角色，系统运行必要');

INSERT INTO menus (id, parent_id, type, name, title, path, component, icon, keep_alive, order_no, external_url, is_hidden, perm_code, api_method, api_path)
VALUES ('019f8d61-b8c2-73f9-8623-729314d73a16'::uuid, NULL, 1, 'system', '系统管理', '/system', NULL, 'system-uicons:grid', false, 1, NULL, false, NULL, NULL, NULL);
INSERT INTO menus (id, parent_id, type, name, title, path, component, icon, keep_alive, order_no, external_url, is_hidden, perm_code, api_method, api_path)
VALUES ('019f8d61-b8c2-73f9-8623-7677d581d813'::uuid, '019f8d61-b8c2-73f9-8623-729314d73a16'::uuid, 2,'system_user', '用户管理', '/system/user', 'layout.base$view.user', 'system-uicons:users', true, 1, NULL, false, NULL, NULL, '/api/v1/user/list');
INSERT INTO menus (id, parent_id, type, name, title, path, component, icon, keep_alive, order_no, external_url, is_hidden, perm_code, api_method, api_path)
VALUES ('019f8d61-b8c2-73f9-8623-7c7bf7338027'::uuid, '019f8d61-b8c2-73f9-8623-729314d73a16'::uuid, 2,'system_role', '角色管理', '/system/role', 'system-uicons:venn', 'system-uicons:user-male-circle', true, 2, NULL, false, NULL, NULL, '/api/v1/role/list');
INSERT INTO menus (id, parent_id, type, name, title, path, component, icon, keep_alive, order_no, external_url, is_hidden, perm_code, api_method, api_path)
VALUES ('019f8d61-b8c2-73f9-8623-80d1c970d611'::uuid, '019f8d61-b8c2-73f9-8623-729314d73a16'::uuid, 2,'system_menu', '菜单管理', '/system/menu', 'layout.base$view.menu', 'system-uicons:browser-alt', true, 3, NULL, false, NULL, NULL, '/api/v1/menu/list');
INSERT INTO menus (id, parent_id, type, name, title, path, component, icon, keep_alive, order_no, external_url, is_hidden, perm_code, api_method, api_path)
VALUES ('019f8d61-b8c2-73f9-8623-8486c5d3e00a'::uuid, '019f8d61-b8c2-73f9-8623-729314d73a16'::uuid, 2,'system_dict', '字典管理', '/system/dict', 'layout.base$view.dict', 'system-uicons:document-list', true, 4, NULL, false, NULL, NULL, '/api/v1/dict/list');

INSERT INTO user_role (user_id, role_id)
VALUES ('019f8d61-b8c2-73f9-8623-6907dc275fce'::uuid, '019fcc22-31fd-7559-bbfc-8c6d5079557d'::uuid);
INSERT INTO user_role (user_id, role_id)
VALUES ('019f8d61-b8c2-73f9-8623-6907dc275fce'::uuid, '019f8d61-b8c2-73f9-8623-6cb937ccb8ee'::uuid);

INSERT INTO role_menu (role_id, menu_id)
VALUES ('019f8d61-b8c2-73f9-8623-6cb937ccb8ee'::uuid, '019f8d61-b8c2-73f9-8623-729314d73a16'::uuid);
INSERT INTO role_menu (role_id, menu_id)
VALUES ('019f8d61-b8c2-73f9-8623-6cb937ccb8ee'::uuid, '019f8d61-b8c2-73f9-8623-7677d581d813'::uuid);
INSERT INTO role_menu (role_id, menu_id)
VALUES ('019f8d61-b8c2-73f9-8623-6cb937ccb8ee'::uuid, '019f8d61-b8c2-73f9-8623-7c7bf7338027'::uuid);
INSERT INTO role_menu (role_id, menu_id)
VALUES ('019f8d61-b8c2-73f9-8623-6cb937ccb8ee'::uuid, '019f8d61-b8c2-73f9-8623-80d1c970d611'::uuid);
INSERT INTO role_menu (role_id, menu_id)
VALUES ('019f8d61-b8c2-73f9-8623-6cb937ccb8ee'::uuid, '019f8d61-b8c2-73f9-8623-8486c5d3e00a'::uuid);

INSERT INTO dictionaries (id, name, category, key, val, order_no, is_default, remark)
VALUES ('019f8d8b-81da-76fd-8371-54434ad85659'::uuid, '数据状态', 'STATUS', '正常', '1', '1', 'false', '系统数据状态');
INSERT INTO dictionaries (id, name, category, key, val, order_no, is_default, remark)
VALUES ('019f8d8b-81da-76fd-8371-5b6b6023bab3'::uuid, '数据状态', 'STATUS', '禁用', '2', '2', 'false', '系统数据状态');
INSERT INTO dictionaries (id, name, category, key, val, order_no, is_default, remark)
VALUES ('019f8d8b-81da-76fd-8371-6783587fd875'::uuid, '菜单类型', 'MENU_TYPE', '目录', '1', '1', 'false', '系统菜单类型');
INSERT INTO dictionaries (id, name, category, key, val, order_no, is_default, remark)
VALUES ('019f8d8b-81da-76fd-8371-6fef25ff307f'::uuid, '菜单类型', 'MENU_TYPE', '菜单', '2', '2', 'false', '系统菜单类型');
INSERT INTO dictionaries (id, name, category, key, val, order_no, is_default, remark)
VALUES ('019f8d8b-81da-76fd-8371-77c364692bd3'::uuid, '菜单类型', 'MENU_TYPE', '按钮', '3', '3', 'false', '系统菜单类型');

INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g', 'leoture', 'common');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'common', '/api/v1/auth/getUserInfo', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'common', '/api/v1/auth/updatePassword', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'common', '/api/v1/auth/updateUser', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'common', '/api/v1/auth/refreshToken', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'common', '/api/v1/route/getConstantRoutes', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'common', '/api/v1/route/isRouteExis', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'common', '/api/v1/auth/logout/:id', '.*', 'allow');

INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g', 'leoture', 'admin');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/menu/create', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/menu/update', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/menu/delete/:id', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/menu/:id', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/menu/list', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/menu/page', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/menu/tree', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/role/create', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/role/update', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/role/delete/:id', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/role/:id', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/role/code/:code', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/role/list', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/role/page', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/role-menu/update', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/role-menu/menus/:roleID', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/user/create', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/user/update', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/user/delete/:id', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/user/:id', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/user/username/:username', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/user/list', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/user/page', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/user-role/update', '.*', 'allow');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'admin', '/api/v1/user-role/roles/:userID', '.*', 'allow');
