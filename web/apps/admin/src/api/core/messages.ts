import { requestClient } from '#/api/request';

export interface UserMessage {
  id: number;
  userId: number;
  title: string;
  content: string;
  isRead: boolean;
  type: string;
  createAt: string;
}

export interface GetUserMessagesParams {
  userId: number;
  page?: number;
  pageSize?: number;
  type?: string;
  isRead?: boolean;
}

export interface MarkMessageAsReadParams {
  userId: number;
  msgId: number;
}

export interface MarkAllMessagesAsReadParams {
  userId: number;
}

/**
 * 获取用户消息列表
 */
export function getUserMessages(params: GetUserMessagesParams) {
  return requestClient.post('/api/system/msg/list', params);
}

/**
 * 标记单条消息为已读
 */
export function markMessageAsRead(params: MarkMessageAsReadParams) {
  return requestClient.post('/api/system/msg/read', params);
}

/**
 * 标记所有消息为已读
 */
export function markAllMessagesAsRead(params: MarkAllMessagesAsReadParams) {
  return requestClient.post('/api/system/msg/read/all', params);
}
