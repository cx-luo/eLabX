import type { UserInfo } from '@vben/types';

import { requestClient } from '#/api/request';

export namespace UserApi {
  export interface ElnUser {
    userId: number;
    username: string;
    realName: string;
    email: string;
    avatar: string;
    roles: string;
    permissions: string;
    status: number;
    groupId: string;
    createAt: string;
    updateAt: string;
  }
}

/**
 * 获取用户信息
 */
export const getUserInfoApi = async () => {
  return requestClient.get<UserInfo>('/system/user/info');
};

/**
 * 获取用户列表
 */
export const getUserListApi = async (param: any) => {
  return requestClient.post<UserApi.ElnUser[]>('/system/user/list', param);
};

/**
 * 新增用户信息
 *
 * @param param 数据
 * @returns
 */
export const createUserApi = async (param: any) => {
  return await requestClient.post('/system/user/add', param);
};

/**
 * 修改用户信息
 *
 * @param id ID
 * @param param 数据
 * @returns
 */
export const updateUserApi = async (id: number, param: any) => {
  return await requestClient.put(`/system/user/update/${id}`, param);
};

/**
 * 删除用户信息
 *
 * @param id ID
 * @returns
 */
export const deleteUserApi = async (id: number) => {
  return await requestClient.delete(`/system/user/delete/${id}`);
};

/**
 * 获取筛选后的用户列表
 *
 * @param param 筛选参数
 * @returns 用户列表
 */
export const getUserListWithFilterApi = async (param: any) => {
  return await requestClient.post<UserApi.ElnUser[]>('/system/user/list/filter', param);
};
