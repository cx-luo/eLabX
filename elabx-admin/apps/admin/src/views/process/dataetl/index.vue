<script lang="ts" setup>
import {computed, h, watch} from 'vue';
import {useVbenVxeGrid, type VxeGridProps} from '#/adapter/vxe-table';
import {$t} from '#/locales';
import {Page, useVbenDrawer, type VbenFormProps} from '@vben/common-ui';
import {LucideFilePenLine} from '@vben/icons';
import {ElButton} from 'element-plus';
import ApiDrawer from './drawer.vue';
// import {useToast, POSITION} from 'vue-toastification';
import {formatDateTime} from '@vben/utils';
import {
  getTableColumnsApi,
  getTableDataApi,
  getTableListApi,
  getDatabaseListApi,
} from '#/api';
import {ref, onMounted} from 'vue';

// const toast = useToast();

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
const columnsOptions = ref<{ label: string; value: string }[]>([]);

// 当前选中的数据库
const selectedDatabase = ref<string | undefined>(undefined);
const selectedTable = ref<string | undefined>(undefined);
const selectedColumns = ref<{ title: string }[] | undefined>(undefined);

// 监听 columnsOptions 变化，自动选中主键列
watch(
  columnsOptions,
  (cols) => {
    if (!cols || cols.length === 0) {
      selectedColumns.value = undefined;
      return;
    }

    const pkCols = (cols as any[]).filter((col) => col.isPrimaryKey);
    // 如果没有主键，默认选第一个
    const defaultCols = pkCols.length > 0 ? pkCols : [cols[0]];
    // 只在未手动选择时设置
    if (!selectedColumns.value || selectedColumns.value.length === 0) {
      selectedColumns.value = defaultCols.map((col) => ({
        title: col.label,
      }));
    }
  },
  {immediate: true}
);

// 禁止清空主键列
function handleColumnsChange(newVal: { title: string }[] | undefined) {
  if (!newVal || newVal.length === 0) {
    return;
  }
  // columnsOptions 里找主键
  const pkCols = (columnsOptions.value as any[]).filter((col) => col.isPrimaryKey);
  const pkTitles = pkCols.map((col) => col.label);
  // 检查主键是否被选中
  const selectedTitles = newVal.map((c) => c.title);
  const missingPK = pkTitles.filter((pk) => !selectedTitles.includes(pk));
  if (missingPK.length > 0) {
    // 自动补回主键
    selectedColumns.value = [
      ...newVal,
      ...pkCols
        .filter((col) => !selectedTitles.includes(col.label))
        .map((col) => ({title: col.label})),
    ];
  } else {
    selectedColumns.value = newVal;
  }
}

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

// 获取字段（列）列表
const fetchColumnsList = async (dbName: string, tableName: string) => {
  if (!dbName || !tableName) {
    columnsOptions.value = [];
    return;
  }
  try {
    const res = await getTableColumnsApi(dbName, tableName);
    // 适配下拉格式
    columnsOptions.value = (res.items || []).map((col: any) => ({
      label:
        typeof col === 'object' && col.columnName
          ? col.columnName
          : String(col),
      value:
        typeof col === 'object' && col.columnName
          ? col.columnName
          : String(col),
    }));
  } catch (e) {
    columnsOptions.value = [];
  }
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
      placeholder: $t('ui.placeholder.select'),
      options: tableOptions,
      onChange: (val: string | undefined) => {
        selectedTable.value = val;
        if (selectedDatabase.value && val) {
          fetchColumnsList(selectedDatabase.value, val);
        } else {
          columnsOptions.value = [];
        }
      },
      // 选中后才可选
      disabled: computed(() => !selectedDatabase.value),
    },
    rules: 'required',
  },
  {
    component: 'Select',
    fieldName: 'columns',
    label: $t('process.etl.columnName'),
    componentProps: {
      placeholder: $t('ui.placeholder.select'),
      options: columnsOptions,
      multiple: true,
      clearable: true,
      filterable: true,
      collapseTags: true,
      collapseTagsTooltip: true,
      onChange: () => handleColumnsChange(selectedColumns.value),
      disabled: computed(() => !selectedTable.value),
    },
  },
];

// export interface VxeColumnProps<D = any> {
//   colId?: VxeColumnPropTypes.ColId
//   /**
//    * 渲染类型
//    */
//   type?: VxeColumnPropTypes.Type
//   /**
//    * 列字段名
//    */
//   field?: VxeColumnPropTypes.Field
//   /**
//    * 列标题
//    */
//   title?: VxeColumnPropTypes.Title
//   /**
//    * 列宽度
//    */
//   width?: VxeColumnPropTypes.Width
//   /**
//    * 列最小宽度，把剩余宽度按比例分配
//    */
//   minWidth?: VxeColumnPropTypes.MinWidth
//   /**
//    * 列最大宽度
//    */
//   maxWidth?: VxeColumnPropTypes.MaxWidth
//   /**
//    * 是否允许拖动列宽调整大小
//    */
//   resizable?: VxeColumnPropTypes.Resizable
//   /**
//    * 将列固定在左侧或者右侧
//    */
//   fixed?: VxeColumnPropTypes.Fixed
//   /**
//    * 列对其方式
//    */
//   align?: VxeColumnPropTypes.Align
//   /**
//    * 表头对齐方式
//    */
//   headerAlign?: VxeColumnPropTypes.HeaderAlign
//   /**
//    * 表尾列的对齐方式
//    */
//   footerAlign?: VxeColumnPropTypes.FooterAlign
//   /**
//    * 当内容过长时显示为省略号
//    */
//   showOverflow?: VxeColumnPropTypes.ShowOverflow
//   /**
//    * 当表头内容过长时显示为省略号
//    */
//   showHeaderOverflow?: VxeColumnPropTypes.ShowHeaderOverflow
//   /**
//    * 当表尾内容过长时显示为省略号
//    */
//   showFooterOverflow?: VxeColumnPropTypes.ShowFooterOverflow
//   /**
//    * 给单元格附加 className
//    */
//   className?: VxeColumnPropTypes.ClassName<D>
//   /**
//    * 给表头单元格附加 className
//    */
//   headerClassName?: VxeColumnPropTypes.HeaderClassName<D>
//   /**
//    * 给表尾单元格附加 className
//    */
//   footerClassName?: VxeColumnPropTypes.FooterClassName<D>
//   /**
//    * 格式化显示内容
//    */
//   formatter?: VxeColumnPropTypes.Formatter<D>
//   /**
//    * 格式化表尾显示内容
//    */
//   footerFormatter?: VxeColumnPropTypes.FooterFormatter<D>
//   /**
//    * 单元格默认高度
//    */
//   padding?: VxeColumnPropTypes.Padding
//   /**
//    * 垂直对齐方式
//    */
//   verticalAlign?: VxeColumnPropTypes.VerticalAlign
//   /**
//    * 是否允许排序
//    */
//   sortable?: VxeColumnPropTypes.Sortable
//   /**
//    * 自定义排序的属性
//    */
//   sortBy?: VxeColumnPropTypes.SortBy<D>
//   /**
//    * 排序的字段类型，比如字符串转数值等
//    */
//   sortType?: VxeColumnPropTypes.SortType
//   /**
//    * 配置筛选条件数组
//    */
//   filters?: VxeColumnPropTypes.Filters
//   /**
//    * 筛选是否允许多选
//    */
//   filterMultiple?: VxeColumnPropTypes.FilterMultiple
//   /**
//    * 自定义筛选方法
//    */
//   filterMethod?: VxeColumnPropTypes.FilterMethod<D>
//   /**
//    * 筛选模板配置项
//    */
//   filterRender?: VxeColumnPropTypes.FilterRender
//   /**
//    * 设置为分组节点
//    */
//   rowGroupNode?: VxeColumnPropTypes.RowGroupNode
//   /**
//    * 设置为树节点
//    */
//   treeNode?: VxeColumnPropTypes.TreeNode
//   /**
//    * 设置为拖拽排序
//    */
//   dragSort?: VxeColumnPropTypes.DragSort
//   /**
//    * 设置为行高拖拽
//    */
//   rowResize?: VxeColumnPropTypes.RowResize
//   /**
//    * 是否可视
//    */
//   visible?: VxeColumnPropTypes.Visible
//   /**
//    * 指定聚合函数
//    */
//   aggFunc?: VxeColumnPropTypes.AggFunc
//   /**
//    * 自定义表尾单元格数据导出方法
//    */
//   headerExportMethod?: VxeColumnPropTypes.HeaderExportMethod<D>
//   /**
//    * 自定义单元格数据导出方法
//    */
//   exportMethod?: VxeColumnPropTypes.ExportMethod<D>
//   /**
//    * 自定义表尾单元格数据导出方法
//    */
//   footerExportMethod?: VxeColumnPropTypes.FooterExportMethod<D>
//   /**
//    * 已废弃，被 titlePrefix 替换
//    * @deprecated
//    */
//   titleHelp?: VxeColumnPropTypes.TitleHelp
//   /**
//    * 标题前缀图标配置项
//    */
//   titlePrefix?: VxeColumnPropTypes.TitlePrefix
//   /**
//    * 标题后缀图标配置项
//    */
//   titleSuffix?: VxeColumnPropTypes.TitleSuffix
//   /**
//    * 单元格值类型
//    */
//   cellType?: VxeColumnPropTypes.CellType
//   /**
//    * 单元格渲染配置项
//    */
//   cellRender?: VxeColumnPropTypes.CellRender<D>
//   /**
//    * 单元格编辑渲染配置项
//    */
//   editRender?: VxeColumnPropTypes.EditRender<D>
//   /**
//    * 内容渲染配置项
//    */
//   contentRender?: VxeColumnPropTypes.ContentRender
//   /**
//    * 额外的参数
//    */
//   params?: VxeColumnPropTypes.Params
// }

const columnsList = ref([
  {
    title: $t('ui.table.seq'),
    type: 'seq',
    width: 70,
    field: null,
    showOverflow: null,
  },
]);

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
  rowConfig: {
    isHover: true,
  },
  verticalAlign: true,
  stripe: true,
  pagerConfig: {
    enabled: true,
    pageSizes: [10, 20, 50, 100],
    layouts: ['PrevPage', 'JumpNumber', 'NextPage', 'Sizes', 'Total'],
  },
  proxyConfig: {
    // Enable server-side paging, sorting, and filtering
    sortConfig: {
      remote: true, // Enable remote sorting
    },
    filterConfig: {
      remote: true, // Enable remote filtering if needed in the future
    },

    ajax: {
      query: async ({page}, formValues) => {
        const dbName = formValues?.database;
        const tableName = formValues?.table;
        const columns = formValues?.columns;
        selectedColumns.value = columns;
        if (!dbName || !tableName) {
          return {result: [], total: 0};
        }
        try {
          // Fetch columns before fetching table data
          if (columns && Array.isArray(columns)) {
            // Dynamically set columns for the grid
            columnsList.value.length = 1; // Keep only the first column (seq)
            columns.forEach((col: any) => {
              columnsList.value.push({
                title:
                  typeof col === 'object' && col.columnName
                    ? col.columnName
                    : String(col),
                type: 'text',
                width: 130,
                showOverflow: true,
                field:
                  typeof col === 'object' && col.columnName
                    ? col.columnName
                    : String(col),
              });
            });

            columnsList.value.push({
              title: $t('ui.table.action'),
              field: 'action',
              fixed: 'right',
              slots: {default: 'action'},
              width: 120,
            });
          }
          return await getTableDataApi(dbName, tableName, {
            page: page.currentPage,
            pageSize: page.pageSize,
            columns: columns,
          });
        } catch (e) {
          // You can display an error message here if needed
          return {result: [], total: 0};
        }
      },
    },

    autoLoad: false, // Prevent initial auto loading; wait for form selection before loading
  },

  columns: columnsList.value,
};

const [Grid, gridApi] = useVbenVxeGrid({gridOptions, formOptions});

const [Drawer, drawerApi] = useVbenDrawer({
  connectedComponent: ApiDrawer,
  onClosed() {
    const data = drawerApi.getData();
    if (data && data.needRefresh) {
      gridApi.query();
    }
  },
});

function openDrawer(create: boolean, row?: any) {
  drawerApi.setData({
    create,
    row,
    columns: selectedColumns.value,
  });
  drawerApi.open();
}

/* 编辑 */
function handleEdit(row: any) {
  openDrawer(false, row);
}

/* 删除 */
// async function handleDelete(row: any) {
//   row.pending = true;
//   try {
//     await deleteApiApi({id: row.id});

//     toast.success($t('ui.notification.delete_success'), {
//       timeout: 1000,
//       position: POSITION.TOP_RIGHT,
//       toastClassName: 'toastification-success',
//     });
//   } catch {
//     toast.error($t('ui.notification.delete_failed'), {
//       timeout: 2000,
//       position: POSITION.TOP_CENTER,
//     });
//   } finally {
//     row.pending = false;
//     await gridApi.query();
//   }
// }
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

      <template #action="{ row }">
        <ElButton
          type="primary"
          link
          v-permission="['system:api:update']"
          :icon="h(LucideFilePenLine)"
          @click="() => handleEdit(row)"
        />

        <!-- <el-popconfirm
          :title="
            $t('ui.text.do_you_want_delete', {
              moduleName: $t('page.system.api.module'),
            })
          "
          :confirm-button-text="$t('ui.button.ok')"
          :cancElButton-text="$t('ui.button.cancel')"
          @confirm="() => handleDelete(row)"
        >
          <template #reference>
            <ElButton
              type="danger"
              v-permission="['system:api:delete']"
              link
              :icon="LucideTrash2"
            />
          </template>
        </el-popconfirm> -->
      </template>
    </Grid>
    <Drawer/>
  </Page>
</template>
