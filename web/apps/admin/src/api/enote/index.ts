import { requestClient } from '#/api/request';

/**
 * Load reagents for a specified reaction
 * @param reactionId The ID of the reaction
 */
export const loadReagents = async (reactionId: number) => {
  return requestClient.get(`/eln/loadReagents/${reactionId}`);
};
