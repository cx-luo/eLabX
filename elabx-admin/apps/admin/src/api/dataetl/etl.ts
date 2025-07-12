import { requestClient } from '#/api/request';

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

/**
 * Get the data rows for a specified table in a database
 * @param dbName Name of the database
 * @param tableName Name of the table
 * @param param Query parameters for filtering, pagination, etc.
 */
export const getTableDataApi = async (
  dbName: string,
  tableName: string,
  param: Record<string, any>,
) => {
  return requestClient.post(`/etl/table/data/${dbName}/${tableName}`, param);
};

/**
 * Update data for a specified table in a database
 * @param dbName Name of the database
 * @param tableName Name of the table
 * @param data Data to update (should include identifying keys and new values)
 */
export const updateTableDataApi = async (
  dbName: string,
  tableName: string,
  data: Record<string, any>,
) => {
  return requestClient.put(`/etl/table/data/${dbName}/${tableName}`, data);
};
