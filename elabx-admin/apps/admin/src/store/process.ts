import { getDatabaseListApi } from '#/api/core/etl';
import { computed } from 'vue';

export const databaseList = computed(async () => {
  const list = await getDatabaseListApi();
  return list.items;
});
