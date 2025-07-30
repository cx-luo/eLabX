<script lang="ts" setup>
import { useVbenDrawer } from '@vben/common-ui';
import { ref } from 'vue';
import type { TreeInstance } from 'element-plus';
import {
  // 新增API相关接口
  getGroupUsersApi,
  getUserListWithFilterApi,
  type UserApi,
  addUserToGroupApi,
  getGroupPermissionsApi,
  assignPermissionsToGroupApi,
} from '#/api';
import { useToast, POSITION } from 'vue-toastification';
import { $t } from '@vben/locales';
import type { TreeKey } from 'element-plus/es/components/tree/src/tree.type.mjs';

const toast = useToast();
const data = ref();
// 添加激活的标签页
const activeTab = ref('menu');

// 添加API树相关数据
const permissionTreeData = ref<any[]>([]);
const permissionTreeRef = ref<TreeInstance>();

const [Drawer, drawerApi] = useVbenDrawer({
  async onOpened() {
    data.value = drawerApi.getData<Record<string, any>>();

    // 同时获取菜单树和API树
    const [usersResult, permissionResult] = await Promise.all([
      getGroupUsersApi({ groupId: data.value?.row.groupId }),
      getGroupPermissionsApi({ groupId: data.value?.row.groupId }),
    ]);

    userTreeData.value = usersResult.items;
    permissionTreeData.value = permissionResult.items;

    if (data.value?.row?.groupId) {
      setTimeout(() => {
        if (usersTreeRef.value) {
          // Filter out null or undefined userIds to satisfy TreeKey[] type
          const checkedKeys = userTreeData.value.map((item) => item.userId);
          usersTreeRef.value.setCheckedKeys(checkedKeys as TreeKey[], false);
        }
      }, 100);
      setTimeout(() => {
        if (permissionTreeRef.value) {
          // Filter out null or undefined userIds to satisfy TreeKey[] type
          const checkedKeys = userTreeData.value.map((item) => item.userId);
          permissionTreeRef.value.setCheckedKeys(checkedKeys as TreeKey[], false);
        }
      }, 100);
    }
  },

  async onConfirm() {
    let userIds: TreeKey[] = [];
    let permissionIds: TreeKey[] = [];
    // 获取菜单权限
    if (usersTreeRef.value) {
      const checkedKeys = usersTreeRef.value.getCheckedKeys() as TreeKey[] | undefined;
      // Ensure selectedUserList is an array of TreeKey
      userIds = [...(checkedKeys ?? []), ...(selectedUserList.value ?? [])];
    }

    // 获取API权限
    if (permissionTreeRef.value) {
      const checkedKeys = permissionTreeRef.value.getCheckedKeys();
      const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys();
      permissionIds = [...checkedKeys, ...halfCheckedKeys];
    }

    if (userIds.length <= 0 && permissionIds.length <= 0) {
      toast.error($t('ui.notification.no_auth'), {
        timeout: 2000,
        position: POSITION.TOP_CENTER,
      });
      return;
    }

    setLoading(true);
    try {
      // 更新权限，同时提交菜单权限和API权限
      await addUserToGroupApi({
        groupId: data.value.row.groupId,
        userId: userIds,
      });

      await assignPermissionsToGroupApi({
        groupId: data.value.row.groupId,
        permissions: permissionIds,
      });

      toast.success($t('ui.notification.update_success'), {
        timeout: 1000,
        position: POSITION.TOP_RIGHT,
        toastClassName: 'toastification-success',
      });
      drawerApi.setData({ needRefresh: true });
    } catch {
      // toast.error($t('ui.notification.update_failed'), {
      //   timeout: 2000,
      //   position: POSITION.TOP_CENTER,
      // });
    } finally {
      drawerApi.close();
      setLoading(false);
    }
  },

  async onClosed() {
    userList.value = undefined;
  },
});

const userTreeData = ref<UserApi.ElnUser[]>([]);

const usersTreeRef = ref<TreeInstance>();

// 递归获取所有节点的 key
const getAllKeys = (data: UserApi.ElnUser[]): number[] => {
  const keys: number[] = [];
  const traverse = (nodes: UserApi.ElnUser[]) => {
    nodes.forEach((node) => {
      if (node.userId !== undefined && node.userId !== null) {
        keys.push(node.userId);
      }
    });
  };
  traverse(data);
  return keys;
};

// 展开所有节点
// const expandAll = () => {
//   if (usersTreeRef.value) {
//     const allKeys = getAllKeys(userTreeData.value);
//     allKeys.forEach((key) => {
//       usersTreeRef.value?.store?.nodesMap[key]?.expand();
//     });
//   }
// };

// 全选所有节点
const checkAll = () => {
  if (usersTreeRef.value) {
    const allKeys = getAllKeys(userTreeData.value);
    usersTreeRef.value.setCheckedKeys(allKeys);
  }
};

// 取消全选
const uncheckAll = () => {
  if (usersTreeRef.value) {
    usersTreeRef.value.setCheckedKeys([]);
  }
};

// 添加API树操作方法
// 展开所有API节点
const expandApiAll = () => {
  if (permissionTreeRef.value) {
    const allKeys = getAllApiKeys(permissionTreeData.value);
    allKeys.forEach((key) => {
      permissionTreeRef.value?.store?.nodesMap[key]?.expand();
    });
  }
};

// 收起所有API节点
const collapseApiAll = () => {
  if (permissionTreeRef.value) {
    const allKeys = getAllApiKeys(permissionTreeData.value);
    allKeys.forEach((key) => {
      permissionTreeRef.value?.store?.nodesMap[key]?.collapse();
    });
  }
};

// 全选所有API节点
const checkApiAll = () => {
  if (permissionTreeRef.value) {
    const allKeys = getAllApiKeys(permissionTreeData.value);
    permissionTreeRef.value.setCheckedKeys(allKeys);
  }
};

// 取消全选API节点
const uncheckApiAll = () => {
  if (permissionTreeRef.value) {
    permissionTreeRef.value.setCheckedKeys([]);
  }
};

// 递归获取所有API节点的key
const getAllApiKeys = (data: any[]): number[] => {
  const keys: number[] = [];
  const traverse = (nodes: any[]) => {
    nodes.forEach((node) => {
      if (node.id !== undefined && node.id !== null) {
        keys.push(node.id);
      }
      if (node.children?.length) {
        traverse(node.children);
      }
    });
  };
  traverse(data);
  return keys;
};

function setLoading(loading: boolean) {
  drawerApi.setState({ loading });
}

const userList = ref<UserApi.ElnUser[]>();
const selectedUserList = ref<TreeKey[]>();

const getAllUsersByFilter = (query: string) => {
  if (query && query.length >= 3) {
    setLoading(true);
    setTimeout(() => {
      getUserListWithFilterApi({ query: query })
        .then((res) => {
          // Assign userList.value to res.items if present, otherwise empty array
          userList.value = res?.items;
          setLoading(false);
        })
        .catch(() => {
          userList.value = [];
          setLoading(false);
        });
    }, 200);
  } else {
    userList.value = [];
  }
};
</script>

<template>
  <Drawer :title="$t('page.system.role.button.auth')">
    <el-tabs v-model="activeTab" class="mb-4" type="border-card">
      <!-- 用户列表标签页 -->
      <el-tab-pane :label="$t('admin.group.drawer.userList')" name="menu">
        <div class="flex flex-col gap-4">
          <!--          <div class="flex flex-col gap-4">-->
          <!--            <ElInput v-model="userList" @input="getAllUsersByFilter" />-->
          <!--          </div>-->
          <el-select
            v-model="selectedUserList"
            class="w-50 m-2"
            placeholder="Enter at least 3 characters"
            multiple
            clearable
            filterable
            remote
            :remote-method="getAllUsersByFilter"
          >
            <el-option
              v-for="item in userList"
              :key="item.userId"
              :label="item.username"
              :value="item.userId"
              style="width: 400px"
            >
              <div class="option-row">
                <span class="flex items-center gap-2">
                  <el-avatar :size="14" :src="data.avatar" />
                </span>
                <span>{{ data.username }}</span>
                <span class="ml-2 text-xs text-gray-500" v-if="data.email">{{ data.email }} </span>
              </div>
            </el-option>
          </el-select>
          <div class="flex gap-2">
            <el-button @click="checkAll">{{ $t('ui.tree.select_all') }}</el-button>
            <el-button @click="uncheckAll">{{ $t('ui.tree.unselect_all') }}</el-button>
          </div>

          <el-tree
            ref="usersTreeRef"
            :data="userTreeData"
            show-checkbox
            node-key="userId"
            :check-strictly="false"
            class="w-full"
          >
            <template #default="{ data }">
              <div class="flex items-center">
                <span class="flex items-center gap-2">
                  <el-avatar :size="14" :src="data.avatar" />
                  <span>{{ data.username }}</span>
                </span>
                <span class="ml-2 text-xs text-gray-500" v-if="data.email">{{ data.email }} </span>
              </div>
            </template>
          </el-tree>
        </div>
      </el-tab-pane>

      <!-- 用户组权限管理标签页 -->
      <el-tab-pane :label="$t('admin.group.drawer.addPermissions')" name="api">
        <div class="flex flex-col gap-4">
          <div class="flex gap-2">
            <el-button @click="expandApiAll">{{ $t('ui.tree.expand_all') }}</el-button>
            <el-button @click="collapseApiAll">{{ $t('ui.tree.collapse_all') }}</el-button>
            <el-button @click="checkApiAll">{{ $t('ui.tree.select_all') }}</el-button>
            <el-button @click="uncheckApiAll">{{ $t('ui.tree.unselect_all') }}</el-button>
          </div>

          <el-tree
            ref="permissionTreeRef"
            :data="permissionTreeData"
            show-checkbox
            node-key="id"
            :check-strictly="false"
            class="w-full"
          >
            <template #default="{ data }">
              <div class="flex items-center">
                <span>{{ data.permissionName }}</span>
                <span v-if="data.description" class="ml-2 text-xs text-gray-400">{{ data.description }}</span>
              </div>
            </template>
          </el-tree>
        </div>
      </el-tab-pane>
    </el-tabs>
  </Drawer>
</template>

<style scoped lang="scss">
.option-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  font-size: 14px;
}

.col-name {
  flex: 0 0 15%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-right: 8px;
}

.col-cas {
  flex: 0 0 30%;
  color: blue;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.col-role {
  flex: 0 0 auto;
  text-align: right;
  color: var(--el-text-color-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
