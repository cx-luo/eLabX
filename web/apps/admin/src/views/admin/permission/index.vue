<script lang="ts" setup>
import { useVbenVxeGrid, type VxeGridProps } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { Page, useVbenDrawer, type VbenFormProps } from '@vben/common-ui';
import { ElButton } from 'element-plus';
import GroupDrawer from './drawer.vue';
import { deletePermissionApi, getGroupListApi, getPermissionListApi, updateGroupApi, updatePermissionApi } from '#/api';
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
        return await getPermissionListApi({
          page: page.currentPage,
          pageSize: page.pageSize,
          query: formValues.name,
          status: formValues.status,
        });
      },
    },
  },

  columns: [
    {
      title: $t('admin.permission.permissionId'),
      field: 'permissionId',
      fixed: 'left',
      width: 'auto',
    },
    {
      title: $t('admin.permission.permissionName'),
      field: 'permissionName',
      width: 'auto',
    },
    {
      title: $t('admin.permission.createdBy'),
      field: 'createdBy',
      minWidth: 120,
    },
    {
      title: $t('admin.permission.description'),
      field: 'description',
      minWidth: 160,
    },
    {
      title: $t('ui.table.status'),
      field: 'status',
      slots: { default: 'status' },
      width: 80,
    },
    {
      title: $t('admin.permission.permissionType'),
      field: 'permissionType',
      minWidth: 120,
      showOverflow: true,
    },
    {
      title: $t('ui.table.createTime'),
      field: 'createAt',
      formatter: 'formatDateTime',
      width: 'auto',
    },
    {
      title: $t('ui.table.updateTime'),
      field: 'updateAt',
      formatter: 'formatDateTime',
      width: 'auto',
    },
    {
      title: $t('ui.table.action'),
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      width: '120',
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
    await updatePermissionApi({ ...row });

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
  connectedComponent: GroupDrawer,
  onClosed() {
    const data = drawerApi.getData();
    if (data && data.needRefresh) {
      gridApi.query();
    }
  },
});

function openDrawer(create: boolean, row?: any) {
  drawerApi.setData({ create, row });
  drawerApi.open();
}

/* 创建 */
function handleCreate() {
  openDrawer(true);
}

function handleEdit(row: any) {
  openDrawer(false, row);
}

function handleDelete(row: any) {
  deletePermissionApi({ permissionId: row.permissionId }).then(() => {
    toast.success($t('ui.notification.delete_success'), {
      timeout: 1000,
      position: POSITION.TOP_RIGHT,
      toastClassName: 'toastification-success',
    });
    gridApi.query();
  });
}
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('admin.permission.title')">
      <template #toolbar-tools>
        <el-button class="mr-2" type="primary" v-permission="['admin:permission:create']" @click="handleCreate">
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
        <span :style="{ marginRight: '15px' }">{{ row.permissionName }}</span>
      </template>

      <template #icon="{ row }">
        <div class="flex h-full items-center justify-center">
          <Icon v-if="row.permissionType === 'menu'" icon="mdi:menu" class="size-4" />
          <Icon v-if="row.permissionType === 'api'" icon="mdi:api" class="size-4" />
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
          :disabled="!hasPermission(['admin:permission:update'])"
        />
      </template>

      <template #action="{ row }">
        <el-button type="primary" link v-permission="['admin:permission:update']" @click="handleEdit(row)">
          <Icon icon="mdi:pencil" class="size-4" />
        </el-button>
        <el-button type="danger" link v-permission="['admin:permission:delete']" @click="handleDelete(row)">
          <Icon icon="mdi:trash-can" class="size-4" />
        </el-button>
      </template>
    </Grid>
    <Drawer />
  </Page>
</template>
