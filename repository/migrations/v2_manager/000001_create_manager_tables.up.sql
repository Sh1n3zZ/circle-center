-- Create manager tables migration
-- Multi-user, multi-project (icon pack), multi-icon system with quota management

-- Projects table: one icon pack = one project
CREATE TABLE projects (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  owner_user_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(255) NOT NULL COMMENT 'Project display name',
  slug VARCHAR(100) NOT NULL COMMENT 'URL-friendly identifier, unique per user',
  package_name VARCHAR(255) NULL COMMENT 'Android package name if used as identifier',
  visibility ENUM('private', 'public') NOT NULL DEFAULT 'private' COMMENT 'Project visibility level',
  description TEXT NULL COMMENT 'Project description',
  icon_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Cached icon count for performance',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  
  -- Indexes and constraints
  PRIMARY KEY (id),
  UNIQUE INDEX idx_unique_user_slug (owner_user_id, slug),
  INDEX idx_owner_user_id (owner_user_id),
  INDEX idx_visibility (visibility),
  INDEX idx_created_at (created_at),
  INDEX idx_updated_at (updated_at),
  
  -- Foreign key constraint
  CONSTRAINT fk_projects_owner_user_id FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
  COMMENT='Icon pack projects table - one project per icon pack';

-- User project roles table: optional RBAC for project collaboration
CREATE TABLE user_project_roles (
  user_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NOT NULL,
  role ENUM('owner', 'admin', 'editor', 'viewer') NOT NULL COMMENT 'User role in the project',
  added_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  
  -- Indexes and constraints
  PRIMARY KEY (user_id, project_id),
  INDEX idx_project_id (project_id),
  INDEX idx_user_id (user_id),
  INDEX idx_role (role),
  
  -- Foreign key constraints
  CONSTRAINT fk_user_project_roles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_project_roles_project_id FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
  COMMENT='User roles and permissions for project collaboration';

-- Project API keys table: for mobile client authentication
CREATE TABLE project_api_keys (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  project_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(100) NOT NULL COMMENT 'Label for key identification and rotation',
  token_hash VARCHAR(255) NOT NULL COMMENT 'Hashed API token, never store plaintext',
  active BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'Whether the key is active',
  last_used_at TIMESTAMP(6) NULL COMMENT 'Last time this key was used',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  
  -- Indexes and constraints
  PRIMARY KEY (id),
  UNIQUE INDEX idx_unique_project_key_name (project_id, name),
  UNIQUE INDEX idx_unique_token_hash (token_hash),
  INDEX idx_project_id (project_id),
  INDEX idx_active (active),
  INDEX idx_last_used_at (last_used_at),
  
  -- Foreign key constraint
  CONSTRAINT fk_project_api_keys_project_id FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
  COMMENT='API keys for project authentication';

-- Icons table: individual icons with status bound to icon level
CREATE TABLE icons (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  project_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(255) NOT NULL COMMENT 'Human-friendly icon label',
  pkg VARCHAR(255) NOT NULL COMMENT 'Package name',
  component_info VARCHAR(500) NOT NULL COMMENT 'Component identifier e.g. com.app/.MainActivity',
  drawable VARCHAR(255) NOT NULL COMMENT 'Expected drawable name inside pack',
  status ENUM('pending', 'in_progress', 'published', 'rejected') NOT NULL DEFAULT 'pending' COMMENT 'Icon processing status',
  metadata JSON NULL COMMENT 'Optional metadata as JSON',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  
  -- Indexes and constraints
  PRIMARY KEY (id),
  UNIQUE INDEX idx_unique_icon_component (project_id, component_info),
  INDEX idx_project_id (project_id),
  INDEX idx_project_status (project_id, status),
  INDEX idx_project_pkg (project_id, pkg),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at),
  INDEX idx_updated_at (updated_at),
  
  -- Foreign key constraint
  CONSTRAINT fk_icons_project_id FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
  COMMENT='Individual icons with status tracking';

-- Icon requests table: incoming request batches tied to projects
CREATE TABLE icon_requests (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  project_id BIGINT UNSIGNED NOT NULL,
  requested_by_user_id BIGINT UNSIGNED NULL COMMENT 'User who made the request if authenticated',
  source ENUM('api', 'email', 'manual') NOT NULL COMMENT 'Request source',
  apps_json JSON NOT NULL COMMENT 'Raw JSON payload from client',
  archive_path VARCHAR(500) NULL COMMENT 'Temporary storage path/URL for zip files (TTL)',
  status ENUM('received', 'processing', 'done', 'failed') NOT NULL DEFAULT 'received' COMMENT 'Request processing status',
  message TEXT NULL COMMENT 'Status message or error details',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  
  -- Indexes and constraints
  PRIMARY KEY (id),
  INDEX idx_project_id (project_id),
  INDEX idx_project_created (project_id, created_at DESC),
  INDEX idx_requested_by_user_id (requested_by_user_id),
  INDEX idx_source (source),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at),
  INDEX idx_updated_at (updated_at),
  
  -- Foreign key constraints
  CONSTRAINT fk_icon_requests_project_id FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_icon_requests_requested_by_user_id FOREIGN KEY (requested_by_user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
  COMMENT='Icon request batches from clients';

-- Request items table: individual items within a request batch
CREATE TABLE request_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  request_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(255) NOT NULL COMMENT 'Component name',
  pkg VARCHAR(255) NOT NULL COMMENT 'Package name',
  component_info VARCHAR(500) NOT NULL COMMENT 'Component identifier',
  drawable VARCHAR(255) NOT NULL COMMENT 'Drawable name',
  matched_icon_id BIGINT UNSIGNED NULL COMMENT 'Reference to matched icon if found',
  resolution ENUM('pending', 'created', 'duplicate', 'rejected') NOT NULL DEFAULT 'pending' COMMENT 'Item resolution status',
  notes TEXT NULL COMMENT 'Processing notes',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  
  -- Indexes and constraints
  PRIMARY KEY (id),
  UNIQUE INDEX idx_unique_req_item (request_id, component_info),
  INDEX idx_request_id (request_id),
  INDEX idx_project_id (project_id),
  INDEX idx_project_component (project_id, component_info),
  INDEX idx_matched_icon_id (matched_icon_id),
  INDEX idx_resolution (resolution),
  INDEX idx_created_at (created_at),
  INDEX idx_updated_at (updated_at),
  
  -- Foreign key constraints
  CONSTRAINT fk_request_items_request_id FOREIGN KEY (request_id) REFERENCES icon_requests(id) ON DELETE CASCADE,
  CONSTRAINT fk_request_items_project_id FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_request_items_matched_icon_id FOREIGN KEY (matched_icon_id) REFERENCES icons(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
  COMMENT='Individual request items within a batch';

-- User quotas table: limit project count per user
CREATE TABLE user_quotas (
  user_id BIGINT UNSIGNED NOT NULL,
  max_projects INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Maximum number of projects user can create (0 = unlimited)',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  
  -- Indexes and constraints
  PRIMARY KEY (user_id),
  
  -- Foreign key constraint
  CONSTRAINT fk_user_quotas_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
  COMMENT='User project creation quotas';

-- Note: Triggers for quota enforcement and icon count updates
-- are implemented in application layer for better compatibility with sqlc
-- 
-- Quota enforcement: Use CheckUserQuota query before creating projects
-- Icon count updates: Use UpdateProjectIconCount query after icon operations
