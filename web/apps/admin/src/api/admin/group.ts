import { requestClient } from '#/api/request';
import type { TreeKey } from 'element-plus/es/components/tree/src/tree.type';

/**
 * Get group list with pagination
 * @param params Query parameters for group list
 */
export function getGroupListApi(params: Record<string, any>) {
  return requestClient.post('/admin/group/list', params);
}

/**
 * Get group detail by groupId
 * @param params Object containing groupId
 */
export function getGroupDetailApi(params: { groupId: string | number }) {
  return requestClient.post('/admin/group/detail', params);
}

/**
 * Create a new group
 * @param data Data for the new group
 */
export function createGroupApi(data: Record<string, any>) {
  return requestClient.post('/admin/group/create', data);
}

/**
 * Update an existing group
 * @param data Data for the group to update
 */
export function updateGroupApi(data: Record<string, any>) {
  return requestClient.post('/admin/group/update', data);
}

/**
 * Delete a group
 * @param params Object containing groupId
 */
export function deleteGroupApi(params: { groupId: string | number }) {
  return requestClient.post('/admin/group/delete', params);
}

/**
 * Assign permissions to a group
 * @param data Object containing groupId and permissions array
 */
export function assignPermissionsToGroupApi(data: { groupId: string | number; permissions: TreeKey[] }) {
  return requestClient.post('/admin/group/assign-permissions', data);
}

/**
 * Get permissions of a group
 * @param params Object containing groupId
 */
export function getGroupPermissionsApi(params: { groupId: string | number }) {
  return requestClient.post('/admin/group/permissions', params);
}

/**
 * Add user to a group
 * @param data Object containing groupId and userId
 */
export function addUserToGroupApi(data: { groupId: string | number; userId: TreeKey[] }) {
  return requestClient.post('/admin/group/add-user', data);
}

/**
 * Remove user from a group
 * @param data Object containing groupId and userId
 */
export function removeUserFromGroupApi(data: { groupId: string | number; userId: TreeKey[] }) {
  return requestClient.post('/admin/group/remove-user', data);
}

/**
 * Get users of a group
 * @param params Object containing groupId
 */
export function getGroupUsersApi(params: { groupId: string | number }) {
  return requestClient.post('/admin/group/users', params);
}
