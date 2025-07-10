import {requestClient} from '#/api/request';

/**
 * Fetch database names and table names from the database.
 */
export const getDatabaseListApi = async () => {
  return requestClient.get('/etl/database/list');
};

/**
 * Get the list of table names for a specified database
 * @param dbName Name of the database
 */
export const getTableListApi = async (dbName: string) => {
  return requestClient.get(`/etl/table/list/${dbName}`);
};

/**
 * Get the list of columns for a specified table in a database
 * @param dbName Name of the database
 * @param tableName Name of the table
 */
export const getTableColumnsApi = async (dbName: string, tableName: string) => {
  return requestClient.get(`/etl/table/columns/${dbName}/${tableName}`);
};
