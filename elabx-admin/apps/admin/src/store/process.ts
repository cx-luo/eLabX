import { getDatabaseListApi } from '#/api';
import { computed } from 'vue';

export const databaseList = computed(async () => {
  const list = await getDatabaseListApi();
  return list.items;
});
