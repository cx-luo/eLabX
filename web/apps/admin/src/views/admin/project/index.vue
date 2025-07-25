<script lang="ts" setup>
import { h } from 'vue';
import { useVbenVxeGrid, type VxeGridProps } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { Page, useVbenDrawer, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';
import { ElButton } from 'element-plus';
import ProjectDrawer from './drawer.vue';
import { deleteProjectApi, getProjectList, updateProjectApi } from '#/api';
import { statusList } from '#/store';
import { Icon } from '@iconify/vue';
import { useToast, POSITION } from 'vue-toastification';
import { hasPermission } from '#/directives/permissions';

const toast = useToast();

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
      fieldName: 'name',
      label: $t('page.system.menu.name'),
      defaultValue: '',
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('ui.table.status'),
      componentProps: {
        options: statusList,
        placeholder: $t('ui.placeholder.select'),
      },
    },
  ],
};

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
    keyField: 'projectId',
  },
  cellConfig: {
    height: 56,
  },
  stripe: true,
  pagerConfig: {
    enabled: true,
    pageSizes: [10, 20, 50, 100],
    layouts: ['PrevPage', 'JumpNumber', 'NextPage', 'Sizes', 'Total'],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        return await getProjectList({
          page: page.currentPage,
          pageSize: page.pageSize,
          name: formValues.name,
          status: formValues.status,
        });
      },
    },
  },

  columns: [
    {
      title: $t('admin.project.projectId'),
      field: 'projectId',
      width: 100,
      fixed: 'left',
    },
    {
      title: $t('admin.project.projectName'),
      field: 'projectName',
      minWidth: 120,
    },
    {
      title: $t('admin.project.createdBy'),
      field: 'createdBy',
      minWidth: 120,
    },
    {
      title: $t('admin.project.description'),
      field: 'description',
      minWidth: 200,
    },
    {
      title: $t('ui.table.status'),
      field: 'status',
      slots: { default: 'status' },
      width: 100,
    },
    {
      title: $t('admin.project.permissions'),
      field: 'permissions',
      minWidth: 150,
      showOverflow: true,
    },
    {
      title: $t('ui.table.createTime'),
      field: 'createAt',
      formatter: 'formatDateTime',
      width: 160,
    },
    {
      title: $t('ui.table.updateTime'),
      field: 'updateAt',
      formatter: 'formatDateTime',
      width: 160,
    },
    {
      title: $t('ui.table.action'),
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      width: 80,
    },
  ],
};

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions,
  formOptions,
});

async function handleStatusChanged(row: any, checked: boolean) {
  row.pending = true;
  row.status = checked ? 1 : 2;
  try {
    await updateProjectApi({ ...row });

    toast.success($t('ui.notification.update_success'), {
      timeout: 1000,
      position: POSITION.TOP_RIGHT,
      toastClassName: 'toastification-success',
    });
  } catch {
    // toast.error($t('ui.notification.update_failed'), {
    //   timeout: 2000,
    //   position: POSITION.TOP_CENTER,
    // });
  } finally {
    row.pending = false;
    await gridApi.query();
  }
}

const [Drawer, drawerApi] = useVbenDrawer({
  connectedComponent: ProjectDrawer,
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
  });
  drawerApi.open();
}

/* 创建 */
function handleCreate() {
  openDrawer(true);
}

/* 编辑 */
function handleEdit(row: any) {
  openDrawer(false, row);
}

/* 删除 */
async function handleDelete(row: any) {
  row.pending = true;
  try {
    await deleteProjectApi({ projectId: row.projectId });

    toast.success($t('ui.notification.delete_success'), {
      timeout: 1000,
      position: POSITION.TOP_RIGHT,
      toastClassName: 'toastification-success',
    });
  } catch {
    // toast.error($t('ui.notification.delete_failed'), {
    //   timeout: 2000,
    //   position: POSITION.TOP_CENTER,
    // });
  } finally {
    row.pending = false;
    await gridApi.query();
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('admin.title')">
      <template #toolbar-tools>
        <el-button class="mr-2" type="primary" v-permission="['admin:project:create']" @click="handleCreate">
          {{ $t('admin.button.create') }}
        </el-button>
        <!--        <el-button class="mr-2" @click="expandAll">-->
        <!--          {{ $t('page.admin.project.button.create') }}-->
        <!--        </el-button>-->
        <!--        <el-button class="mr-2" @click="collapseAll">-->
        <!--          {{ $t('page.admin.project.button.create') }}-->
        <!--        </el-button>-->
      </template>

      <template #title="{ row }">
        <span :style="{ marginRight: '15px' }">{{ $t(row.meta.name) }}</span>
      </template>

      <template #icon="{ row }">
        <div class="flex h-full items-center justify-center">
          <Icon v-if="row.meta.icon !== undefined" :icon="row.meta.icon" class="size-4" />
        </div>
      </template>

      <template #status="{ row }">
        <el-switch
          :model-value="row.status === 1"
          :loading="row.pending"
          :inline-prompt="true"
          :active-text="$t('ui.switch.active')"
          :inactive-text="$t('ui.switch.inactive')"
          @change="(checked: boolean) => handleStatusChanged(row, checked)"
          :disabled="!hasPermission(['admin:project:update'])"
        />
      </template>

      <template #action="{ row }">
        <ElButton
          type="primary"
          link
          v-permission="['admin:project:update']"
          :icon="h(LucideFilePenLine)"
          @click="() => handleEdit(row)"
        />

        <el-popconfirm
          :title="
            $t('ui.text.do_you_want_delete', {
              moduleName: $t('page.system.menu.module'),
            })
          "
          :confirm-button-text="$t('ui.button.ok')"
          :cancElButton-text="$t('ui.button.cancel')"
          @confirm="() => handleDelete(row)"
        >
          <template #reference>
            <ElButton type="danger" v-permission="['admin:project:delete']" link :icon="LucideTrash2" />
          </template>
        </el-popconfirm>
      </template>
    </Grid>
    <Drawer />
  </Page>
</template>
