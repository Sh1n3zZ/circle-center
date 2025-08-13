-- Drop manager tables migration
-- Rollback for creating manager tables

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS request_items;
DROP TABLE IF EXISTS icon_requests;
DROP TABLE IF EXISTS icons;
DROP TABLE IF EXISTS project_api_keys;
DROP TABLE IF EXISTS user_project_roles;
DROP TABLE IF EXISTS user_quotas;
DROP TABLE IF EXISTS projects;
