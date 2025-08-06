<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { useVbenForm } from '#/adapter/form';
import { createPermissionApi, updatePermissionApi } from '#/api';
import { statusList } from '#/store';
import { useToast, POSITION } from 'vue-toastification';

const toast = useToast();

const data = ref();

const getTitle = computed(() =>
  data.value?.create
    ? $t('ui.modal.create', { moduleName: $t('admin.permission.title') })
    : $t('ui.modal.update', { moduleName: $t('admin.permission.title') }),
);
const newSchema = ref([
  {
    component: 'Input',
    fieldName: 'permissionId',
    label: $t('admin.permission.permissionId'),
    disabled: true,
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'permissionName',
    label: $t('admin.permission.permissionName'),
    rules: 'required',
    componentProps: {
      placeholder: $t('ui.placeholder.input'),
    },
  },
  {
    component: 'Input',
    fieldName: 'permissionType',
    label: $t('admin.permission.permissionType'),
    rules: 'required',
    componentProps: {
      placeholder: $t('ui.placeholder.input'),
      allowClear: true,
    },
  },
  {
    component: 'Input',
    fieldName: 'description',
    label: $t('admin.permission.description'),
    componentProps: {
      placeholder: $t('ui.placeholder.input'),
      allowClear: true,
    },
  },
  {
    component: 'InputNumber',
    fieldName: 'createdBy',
    label: $t('admin.permission.createdBy'),
    componentProps: {
      placeholder: $t('ui.placeholder.input'),
      allowClear: true,
    },
  },
  // Permissions are usually set separately, so omit from main group form
  {
    component: 'RadioGroup',
    fieldName: 'status',
    defaultValue: 2,
    label: $t('ui.table.status'),
    rules: 'selectRequired',
    componentProps: {
      optionType: 'button',
      class: 'flex flex-wrap',
      options: statusList,
    },
  },
]);
const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,
  // 所有表单项共用，可单独在表单内覆盖
  commonConfig: {
    // 所有表单项
    componentProps: {
      class: 'w-full',
    },
  },
  // Add new group or modify group info
  schema: computed(() => {
    if (data.value?.create) {
      return newSchema.value.filter((item) => item.fieldName !== 'permissionId');
    }
    return newSchema.value;
  }),
});

const [Drawer, drawerApi] = useVbenDrawer({
  onCancel() {
    drawerApi.close();
  },

  async onConfirm() {
    // 校验输入的数据
    const validate = await baseFormApi.validate();
    if (!validate.valid) {
      return;
    }

    setLoading(true);

    // 获取表单数据
    const values = await baseFormApi.getValues();

    try {
      await (data.value?.create ? createPermissionApi(values) : updatePermissionApi(values));

      toast.success(data.value?.create ? $t('ui.notification.create_success') : $t('ui.notification.update_success'), {
        timeout: 1000,
        position: POSITION.TOP_RIGHT,
        toastClassName: 'toastification-success',
      });
      drawerApi.setData({ needRefresh: true });
    } catch {
      // toast.error(
      //   data.value?.create
      //     ? $t('ui.notification.create_failed')
      //     : $t('ui.notification.update_failed'),
      //   {
      //     timeout: 2000,
      //     position: POSITION.TOP_CENTER,
      //   },
      // );
    } finally {
      drawerApi.close();
      setLoading(false);
    }
  },

  onOpenChange(isOpen) {
    if (isOpen) {
      // 获取传入的数据
      data.value = drawerApi.getData<Record<string, any>>();

      // 为表单赋值
      baseFormApi.setValues(data.value?.row);

      setLoading(false);
    }
  },
});

function setLoading(loading: boolean) {
  drawerApi.setState({ loading });
}
</script>

<template>
  <Drawer :title="getTitle">
    <BaseForm />
  </Drawer>
</template>
