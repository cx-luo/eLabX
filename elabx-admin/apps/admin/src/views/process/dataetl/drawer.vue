<script lang="ts" setup>
import { ref } from 'vue';
import { useVbenDrawer, z } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { useVbenForm } from '#/adapter/form';
import { updateTableDataRowApi } from '#/api';
import { useToast, POSITION } from 'vue-toastification';

const toast = useToast();
const data = ref();
const newSchema = ref<any[]>([]);

const primaryKey = ref<string[]>([]);

const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,
  // 所有表单项共用，可单独在表单内覆盖
  commonConfig: {
    // 所有表单项
    componentProps: {
      class: 'w-full',
    },
  },
  schema: newSchema.value,
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
      await updateTableDataRowApi(
        data.value.dbName,
        data.value.tableName,
        primaryKey.value,
        { ...values },
      );

      toast.success(
        data.value?.create
          ? $t('ui.notification.create_success')
          : $t('ui.notification.update_success'),
        {
          timeout: 1000,
          position: POSITION.TOP_RIGHT,
          toastClassName: 'toastification-success',
        },
      );

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

      (data.value?.columns || []).map((col: any) => {
        newSchema.value.push({
          component: 'Input',
          fieldName: col.value || col.columnName || String(col),
          label: col.label || col.columnName || String(col),
          componentProps: {
            placeholder: $t('ui.placeholder.input'),
            allowClear: true,
          },
          // rules: z.string().min(1, {message: $t('ui.formRules.required')}),
        });
      });

      // 根据 columns 得到 primaryKey
      if (data.value?.columns) {
        primaryKey.value = data.value.columns
          .filter(
            (col: any) => col.primaryKey === true || col.isPrimaryKey === true,
          )
          .map((col: any) => col.value || col.columnName || String(col));
      }

      baseFormApi.updateSchema(newSchema.value);

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
  <Drawer :title="$t('ui.modal.update')">
    <BaseForm />
  </Drawer>
</template>
