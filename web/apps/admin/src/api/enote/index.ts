import { requestClient } from '#/api/request';

/**
 * Load reagents for a specified reaction
 * @param reactionId The ID of the reaction
 */
export const loadReagents = async (reactionId: number) => {
  return requestClient.get(`/eln/loadReagents/${reactionId}`);
};


/**
 * Update a row in a specified table in a database
 * @param dbName Name of the database
 * @param tableName Name of the table
 * @param primaryKey Name of the primary key column
 * @param data Data to update (must include the primary key value and fields to update)
 */
export const updateTableDataRowApi = async (
    dbName: string,
    tableName: string,
    primaryKey: string[],
    data: Record<string, any>,
) => {
  return requestClient.post('/etl/table/data/update', {
    dbName,
    tableName,
    primaryKey,
    data,
  });
};