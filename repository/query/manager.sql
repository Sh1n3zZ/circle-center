-- =============================================================================
-- PROJECTS MANAGEMENT
-- =============================================================================

-- name: CreateProject :execresult
INSERT INTO projects (
  owner_user_id, name, slug, package_name, visibility, description
) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id = ? LIMIT 1;

-- name: GetProjectBySlug :one
SELECT * FROM projects WHERE owner_user_id = ? AND slug = ? LIMIT 1;

-- name: GetProjectByIDAndOwner :one
SELECT * FROM projects WHERE id = ? AND owner_user_id = ? LIMIT 1;

-- name: ListProjectsByOwner :many
SELECT * FROM projects WHERE owner_user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: ListPublicProjects :many
SELECT * FROM projects WHERE visibility = 'public' ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: ListProjectsByVisibility :many
SELECT * FROM projects WHERE visibility = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: UpdateProject :exec
UPDATE projects SET 
  name = ?, 
  slug = ?, 
  package_name = ?, 
  visibility = ?, 
  description = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ? AND owner_user_id = ?;

-- name: UpdateProjectIconCount :exec
UPDATE projects SET 
  icon_count = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = ? AND owner_user_id = ?;

-- name: CountProjectsByOwner :one
SELECT COUNT(*) FROM projects WHERE owner_user_id = ?;

-- name: CountProjectsByVisibility :one
SELECT COUNT(*) FROM projects WHERE visibility = ?;

-- name: GetProjectStats :one
SELECT 
  COUNT(*) as total_projects,
  SUM(icon_count) as total_icons,
  AVG(icon_count) as avg_icons_per_project
FROM projects 
WHERE owner_user_id = ?;

-- =============================================================================
-- USER PROJECT ROLES MANAGEMENT
-- =============================================================================

-- name: CreateUserProjectRole :execresult
INSERT INTO user_project_roles (user_id, project_id, role) VALUES (?, ?, ?);

-- name: GetUserProjectRole :one
SELECT * FROM user_project_roles WHERE user_id = ? AND project_id = ? LIMIT 1;

-- name: ListProjectCollaborators :many
SELECT upr.*, u.username, u.display_name, u.avatar_url
FROM user_project_roles upr
JOIN users u ON upr.user_id = u.id
WHERE upr.project_id = ?
ORDER BY upr.added_at ASC;

-- name: ListUserProjects :many
SELECT upr.*, p.name, p.slug, p.visibility, p.icon_count
FROM user_project_roles upr
JOIN projects p ON upr.project_id = p.id
WHERE upr.user_id = ?
ORDER BY upr.added_at DESC;

-- name: UpdateUserProjectRole :exec
UPDATE user_project_roles SET role = ? WHERE user_id = ? AND project_id = ?;

-- name: DeleteUserProjectRole :exec
DELETE FROM user_project_roles WHERE user_id = ? AND project_id = ?;

-- name: DeleteProjectCollaborators :exec
DELETE FROM user_project_roles WHERE project_id = ?;

-- name: CountProjectCollaborators :one
SELECT COUNT(*) FROM user_project_roles WHERE project_id = ?;

-- =============================================================================
-- PROJECT API KEYS MANAGEMENT
-- =============================================================================

-- name: CreateProjectAPIKey :execresult
INSERT INTO project_api_keys (project_id, name, token_hash) VALUES (?, ?, ?);

-- name: GetProjectAPIKeyByID :one
SELECT * FROM project_api_keys WHERE id = ? LIMIT 1;

-- name: GetProjectAPIKeyByHash :one
SELECT * FROM project_api_keys WHERE token_hash = ? AND active = TRUE LIMIT 1;

-- name: ListProjectAPIKeys :many
SELECT * FROM project_api_keys WHERE project_id = ? ORDER BY created_at DESC;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE project_api_keys SET 
  last_used_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: DeactivateAPIKey :exec
UPDATE project_api_keys SET 
  active = FALSE,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ? AND project_id = ?;

-- name: DeleteAPIKey :exec
DELETE FROM project_api_keys WHERE id = ? AND project_id = ?;

-- name: DeleteProjectAPIKeys :exec
DELETE FROM project_api_keys WHERE project_id = ?;

-- name: CountActiveAPIKeys :one
SELECT COUNT(*) FROM project_api_keys WHERE project_id = ? AND active = TRUE;

-- =============================================================================
-- ICONS MANAGEMENT
-- =============================================================================

-- name: CreateIcon :execresult
INSERT INTO icons (
  project_id, name, pkg, component_info, drawable, status, metadata
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetIconByID :one
SELECT * FROM icons WHERE id = ? LIMIT 1;

-- name: GetIconByComponent :one
SELECT * FROM icons WHERE project_id = ? AND component_info = ? LIMIT 1;

-- name: ListProjectIcons :many
SELECT * FROM icons WHERE project_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: ListIconsByStatus :many
SELECT * FROM icons WHERE project_id = ? AND status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: ListIconsByPackage :many
SELECT * FROM icons WHERE project_id = ? AND pkg = ? ORDER BY name ASC;

-- name: UpdateIcon :exec
UPDATE icons SET 
  name = ?, 
  pkg = ?, 
  component_info = ?, 
  drawable = ?, 
  status = ?, 
  metadata = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ? AND project_id = ?;

-- name: UpdateIconStatus :exec
UPDATE icons SET 
  status = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ? AND project_id = ?;

-- name: DeleteIcon :exec
DELETE FROM icons WHERE id = ? AND project_id = ?;

-- name: DeleteProjectIcons :exec
DELETE FROM icons WHERE project_id = ?;

-- name: CountProjectIcons :one
SELECT COUNT(*) FROM icons WHERE project_id = ?;

-- name: CountIconsByStatus :one
SELECT COUNT(*) FROM icons WHERE project_id = ? AND status = ?;

-- name: GetIconStats :one
SELECT 
  COUNT(*) as total_icons,
  COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_count,
  COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_count,
  COUNT(CASE WHEN status = 'published' THEN 1 END) as published_count,
  COUNT(CASE WHEN status = 'rejected' THEN 1 END) as rejected_count
FROM icons 
WHERE project_id = ?;

-- =============================================================================
-- ICON REQUESTS MANAGEMENT
-- =============================================================================

-- name: CreateIconRequest :execresult
INSERT INTO icon_requests (
  project_id, requested_by_user_id, source, apps_json, archive_path
) VALUES (?, ?, ?, ?, ?);

-- name: GetIconRequestByID :one
SELECT * FROM icon_requests WHERE id = ? LIMIT 1;

-- name: GetIconRequestByIDAndProject :one
SELECT * FROM icon_requests WHERE id = ? AND project_id = ? LIMIT 1;

-- name: ListProjectRequests :many
SELECT * FROM icon_requests WHERE project_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: ListRequestsByStatus :many
SELECT * FROM icon_requests WHERE project_id = ? AND status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: UpdateRequestStatus :exec
UPDATE icon_requests SET 
  status = ?,
  message = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ? AND project_id = ?;

-- name: UpdateRequestArchivePath :exec
UPDATE icon_requests SET 
  archive_path = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ? AND project_id = ?;

-- name: DeleteIconRequest :exec
DELETE FROM icon_requests WHERE id = ? AND project_id = ?;

-- name: DeleteProjectRequests :exec
DELETE FROM icon_requests WHERE project_id = ?;

-- name: CountProjectRequests :one
SELECT COUNT(*) FROM icon_requests WHERE project_id = ?;

-- name: CountRequestsByStatus :one
SELECT COUNT(*) FROM icon_requests WHERE project_id = ? AND status = ?;

-- name: GetRequestStats :one
SELECT 
  COUNT(*) as total_requests,
  COUNT(CASE WHEN status = 'received' THEN 1 END) as received_count,
  COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing_count,
  COUNT(CASE WHEN status = 'done' THEN 1 END) as done_count,
  COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count
FROM icon_requests 
WHERE project_id = ?;

-- =============================================================================
-- REQUEST ITEMS MANAGEMENT
-- =============================================================================

-- name: CreateRequestItem :execresult
INSERT INTO request_items (
  request_id, project_id, name, pkg, component_info, drawable
) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetRequestItemByID :one
SELECT * FROM request_items WHERE id = ? LIMIT 1;

-- name: GetRequestItemByComponent :one
SELECT * FROM request_items WHERE request_id = ? AND component_info = ? LIMIT 1;

-- name: ListRequestItems :many
SELECT * FROM request_items WHERE request_id = ? ORDER BY created_at ASC;

-- name: ListProjectRequestItems :many
SELECT * FROM request_items WHERE project_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: ListItemsByResolution :many
SELECT * FROM request_items WHERE request_id = ? AND resolution = ? ORDER BY created_at ASC;

-- name: UpdateRequestItem :exec
UPDATE request_items SET 
  name = ?, 
  pkg = ?, 
  component_info = ?, 
  drawable = ?, 
  matched_icon_id = ?, 
  resolution = ?, 
  notes = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ? AND request_id = ?;

-- name: UpdateItemResolution :exec
UPDATE request_items SET 
  resolution = ?,
  matched_icon_id = ?,
  notes = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ? AND request_id = ?;

-- name: DeleteRequestItem :exec
DELETE FROM request_items WHERE id = ? AND request_id = ?;

-- name: DeleteRequestItems :exec
DELETE FROM request_items WHERE request_id = ?;

-- name: DeleteProjectRequestItems :exec
DELETE FROM request_items WHERE project_id = ?;

-- name: CountRequestItems :one
SELECT COUNT(*) FROM request_items WHERE request_id = ?;

-- name: CountItemsByResolution :one
SELECT COUNT(*) FROM request_items WHERE request_id = ? AND resolution = ?;

-- name: GetItemStats :one
SELECT 
  COUNT(*) as total_items,
  COUNT(CASE WHEN resolution = 'pending' THEN 1 END) as pending_count,
  COUNT(CASE WHEN resolution = 'created' THEN 1 END) as created_count,
  COUNT(CASE WHEN resolution = 'duplicate' THEN 1 END) as duplicate_count,
  COUNT(CASE WHEN resolution = 'rejected' THEN 1 END) as rejected_count
FROM request_items 
WHERE request_id = ?;

-- =============================================================================
-- USER QUOTAS MANAGEMENT
-- =============================================================================

-- name: CreateUserQuota :execresult
INSERT INTO user_quotas (user_id, max_projects) VALUES (?, ?);

-- name: GetUserQuota :one
SELECT * FROM user_quotas WHERE user_id = ? LIMIT 1;

-- name: UpdateUserQuota :exec
UPDATE user_quotas SET 
  max_projects = ?
WHERE user_id = ?;

-- name: DeleteUserQuota :exec
DELETE FROM user_quotas WHERE user_id = ?;

-- name: CheckUserQuota :one
SELECT 
  uq.max_projects,
  COALESCE(p.project_count, 0) as current_projects,
  CASE 
    WHEN uq.max_projects = 0 THEN TRUE
    WHEN uq.max_projects IS NULL THEN TRUE
    ELSE COALESCE(p.project_count, 0) < uq.max_projects
  END as can_create_project
FROM user_quotas uq
LEFT JOIN (
  SELECT owner_user_id, COUNT(*) as project_count 
  FROM projects 
  WHERE owner_user_id = ? 
  GROUP BY owner_user_id
) p ON uq.user_id = p.owner_user_id
WHERE uq.user_id = ?;

-- =============================================================================
-- COMPLEX QUERIES AND JOINS
-- =============================================================================

-- name: GetProjectWithStats :one
SELECT 
  p.*,
  COALESCE(icon_stats.total_icons, 0) as total_icons,
  COALESCE(icon_stats.published_icons, 0) as published_icons,
  COALESCE(request_stats.total_requests, 0) as total_requests,
  COALESCE(collaborator_stats.collaborator_count, 0) as collaborator_count
FROM projects p
LEFT JOIN (
  SELECT 
    project_id,
    COUNT(*) as total_icons,
    COUNT(CASE WHEN status = 'published' THEN 1 END) as published_icons
  FROM icons 
  GROUP BY project_id
) icon_stats ON p.id = icon_stats.project_id
LEFT JOIN (
  SELECT 
    project_id,
    COUNT(*) as total_requests
  FROM icon_requests 
  GROUP BY project_id
) request_stats ON p.id = request_stats.project_id
LEFT JOIN (
  SELECT 
    project_id,
    COUNT(*) as collaborator_count
  FROM user_project_roles 
  GROUP BY project_id
) collaborator_stats ON p.id = collaborator_stats.project_id
WHERE p.id = ?;

-- name: GetIconWithRequestInfo :one
SELECT 
  i.*,
  ri.request_id,
  ri.resolution as request_resolution,
  ri.notes as request_notes
FROM icons i
LEFT JOIN request_items ri ON i.id = ri.matched_icon_id
WHERE i.id = ?;

-- name: ListRecentActivity :many
SELECT 
  'icon_request' as activity_type,
  ir.id as activity_id,
  ir.project_id,
  ir.status,
  ir.created_at,
  p.name as project_name,
  u.username as requested_by
FROM icon_requests ir
JOIN projects p ON ir.project_id = p.id
LEFT JOIN users u ON ir.requested_by_user_id = u.id
WHERE ir.project_id = ?
UNION ALL
SELECT 
  'icon_status_change' as activity_type,
  i.id as activity_id,
  i.project_id,
  i.status,
  i.updated_at as created_at,
  p.name as project_name,
  NULL as requested_by
FROM icons i
JOIN projects p ON i.project_id = p.id
WHERE i.project_id = ? AND i.updated_at > i.created_at
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: SearchIcons :many
SELECT * FROM icons 
WHERE project_id = ? 
  AND (name LIKE ? OR pkg LIKE ? OR component_info LIKE ?)
ORDER BY 
  CASE WHEN name LIKE ? THEN 1 
       WHEN pkg LIKE ? THEN 2 
       ELSE 3 END,
  created_at DESC
LIMIT ? OFFSET ?;

-- name: GetDuplicateIcons :many
SELECT 
  i1.*,
  i2.id as duplicate_id,
  i2.name as duplicate_name,
  i2.created_at as duplicate_created_at
FROM icons i1
JOIN icons i2 ON i1.component_info = i2.component_info 
  AND i1.project_id = i2.project_id 
  AND i1.id != i2.id
WHERE i1.project_id = ?
ORDER BY i1.component_info, i1.created_at ASC;
