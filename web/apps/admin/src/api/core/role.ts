import { requestClient } from '#/api/request';

/**
 * 获取角色列表
 */
export const getRoleListApi = async (params: any) => {
  return requestClient.post('/system/role/list', params);
};

/**
 * 获取角色信息
 */
export const getRoleInfoApi = async (id: number) => {
  return requestClient.get(`/system/role/info/${id}`);
};

/**
 * 新增角色信息
 *
 * @param param 数据
 * @returns
 */
export const createRoleApi = async (param: any) => {
  return await requestClient.post('/system/role/add', param);
};

/**
 * 修改角色信息
 *
 * @param param 数据
 * @returns
 */
export const updateRoleApi = async (param: any) => {
  return await requestClient.post(`/system/role/update`, param);
};

/**
 * 删除角色信息
 *
 * @param id ID
 * @returns
 */
export const deleteRoleApi = async (id: number) => {
  return await requestClient.delete(`/system/role/delete/${id}`);
};

/**
 * 分配权限
 *
 * @param param 数据
 * @returns
 */
export const updateRoleAuthApi = async (param: any) => {
  return await requestClient.post(`/system/role/assign`, param);
};
