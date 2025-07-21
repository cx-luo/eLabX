import { requestClient } from '#/api/request';

/**
 * Get the list of projects
 * @param params Query parameters for filtering, pagination, etc.
 */
export function getProjectList(params: Record<string, any>) {
  return requestClient.post('/admin/project/list', params);
}

/**
 * Get the details of a specific project
 * @param params Parameters to identify the project (e.g., projectId)
 */
export function getProjectDetail(params: Record<string, any>) {
  return requestClient.post('/admin/project/detail', params);
}

/**
 * Create a new project
 * @param data Data for the new project
 */
export function createProject(data: Record<string, any>) {
  return requestClient.post('/admin/project/create', data);
}

/**
 * Update an existing project
 * @param data Data to update the project (should include project identifier)
 */
export function updateProject(data: Record<string, any>) {
  return requestClient.post('/admin/project/update', data);
}

/**
 * Delete a project
 * @param params Parameters to identify the project to delete (e.g., projectId)
 */
export function deleteProject(params: Record<string, any>) {
  return requestClient.post('/admin/project/delete', params);
}
