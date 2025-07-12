<script lang="ts" setup>
import { computed } from 'vue';
import { useVbenVxeGrid, type VxeGridProps } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { Page, type VbenFormProps } from '@vben/common-ui';
import { formatDateTime } from '@vben/utils';
import { getTableColumnsApi } from '#/api/core/etl';
import { ref, onMounted } from 'vue';
import { getDatabaseListApi, getTableListApi } from '#/api/core/etl';

const formOptions: VbenFormProps = {
  // 默认展开
  collapsed: false,
  // 控制表单是否显示折叠按钮
  showCollapseButton: false,
  // 按下回车时是否提交表单
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      fieldName: 'database',
      label: $t('process.etl.database'),
      componentProps: {
        allowClear: true,
        placeholder: $t('ui.placeholder.input'),
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'table',
      label: $t('process.etl.table'),
      componentProps: {
        allowClear: true,
        placeholder: $t('ui.placeholder.input'),
      },
      rules: 'required',
    },
  ],
};

// 数据库和数据表下拉选项
const databaseOptions = ref<{ label: string; value: string }[]>([]);
const tableOptions = ref<{ label: string; value: string }[]>([]);

// 当前选中的数据库
const selectedDatabase = ref<string | undefined>(undefined);

// 获取数据库列表
const fetchDatabaseList = async () => {
  const res = await getDatabaseListApi();
  databaseOptions.value = res.items || [];
};

// 获取表列表
const fetchTableList = async (dbName: string) => {
  if (!dbName) {
    tableOptions.value = [];
    return;
  }
  const res = await getTableListApi(dbName);
  // 适配下拉格式
  tableOptions.value = (res.items || []).map((t: string) => ({
    label: t,
    value: t,
  }));
};

// 页面加载时获取数据库列表
onMounted(() => {
  fetchDatabaseList();
});

// 表单 schema 优化为下拉选择
formOptions.schema = [
  {
    component: 'Select',
    fieldName: 'database',
    label: $t('process.etl.database'),
    componentProps: {
      allowClear: true,
      placeholder: $t('ui.placeholder.select'),
      options: databaseOptions,
      onChange: (val: string) => {
        selectedDatabase.value = val;
        // 清空表选择
        (formOptions.schema?.[1]?.componentProps as any).value = undefined;
        fetchTableList(val);
      },
    },
    rules: 'required',
  },
  {
    component: 'Select',
    fieldName: 'table',
    label: $t('process.etl.table'),
    componentProps: {
      allowClear: true,
      placeholder: $t('ui.placeholder.select'),
      options: tableOptions,
      // 选中数据库后才可选
      disabled: computed(() => !selectedDatabase.value),
    },
    rules: 'required',
  },
];

const gridOptions: VxeGridProps = {
  toolbarConfig: {
    custom: true,
    export: true,
    // import: true,
    refresh: true,
    zoom: true,
  },
  height: 'auto',
  exportConfig: {},
  pagerConfig: {
    enabled: false,
  },
  rowConfig: {
    isHover: true,
  },
  stripe: true,
  // proxyConfig 用于配置表格的数据代理，主要用于处理数据的查询、分页、排序等操作，将这些操作与后端 API 解耦，简化表格与后端的数据交互。
  // 优化建议：将 query 的数据处理与异常处理做更清晰的分离，返回结构与 columns 结构保持一致，避免直接返回 API 原始数据。
  proxyConfig: {
    ajax: {
      query: async ({}, formValues) => {
        const dbName = formValues?.database;
        const tableName = formValues?.table;
        if (!dbName || !tableName) {
          return { result: [], total: 0 };
        }
        try {
          return await getTableColumnsApi(dbName, tableName);
        } catch (e) {
          // 可以根据需要弹出错误提示
          return { result: [], total: 0 };
        }
      },
    },
    // 关闭分页（由 pagerConfig 控制），如需分页可在此扩展
    autoLoad: false, // 避免初始自动加载，等表单选择后再加载
  },

  columns: [
    {
      title: $t('ui.table.seq'),
      type: 'seq',
      width: 70,
    },
    {
      title: $t('process.etl.columnName'),
      field: 'columnName',
    },
    {
      title: $t('process.etl.dataType'),
      field: 'dataType',
    },
    {
      title: $t('process.etl.isNullable'),
      field: 'isNullable',
    },
    {
      title: $t('process.etl.columnDefault'),
      field: 'columnDefault',
    },
    {
      title: $t('process.etl.columnComment'),
      field: 'columnComment',
    },
  ],
};

const [Grid] = useVbenVxeGrid({ gridOptions, formOptions });
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('page.system.api.title')">
      <template #createdAt="{ row }">
        {{ formatDateTime(row?.createdAt) }}
      </template>

      <template #parentId="{ row }">
        <span :style="{ marginRight: '15px' }">
          <template v-if="row.parentId === 0"> 根API </template></span
        >
      </template>
    </Grid>
  </Page>
</template>
