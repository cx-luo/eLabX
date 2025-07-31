import { requestClient } from '#/api/request';
import type { UserApi } from '../core';

/**
 * Get the list of all permissions
 * @param params Optional filter parameters
 */
export function getPermissionListApi(params?: Record<string, any>) {
  return requestClient.post('/admin/permission/list', params);
}

/**
 * Create a new permission
 * @param data Permission data
 */
export function createPermissionApi(data: Record<string, any>) {
  return requestClient.post('/admin/permission/create', data);
}

/**
 * Update an existing permission
 * @param data Permission data with id
 */
export function updatePermissionApi(data: Record<string, any>) {
  return requestClient.post('/admin/permission/update', data);
}

/**
 * Delete a permission
 * @param params Object containing permission id
 */
export function deletePermissionApi(params: { permissionId: string | number }) {
  return requestClient.post('/admin/permission/delete', params);
}

/**
 * Get filtered permission list (for remote search, etc.)
 * @param params Filter parameters
 * @returns Promise resolving to filtered permission list and total count
 */
export function getPermsListWithFilterApi(params: Record<string, any>) {
  return requestClient.post<{ items: UserApi.ElnPermission[]; total: number }>('/admin/permission/list/filter', params);
}
