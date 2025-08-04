<script lang="ts" setup>
import { useVbenDrawer } from '@vben/common-ui';
import { ref } from 'vue';
import type { TreeInstance } from 'element-plus';
import {
  getGroupUsersApi,
  getUserListWithFilterApi,
  type UserApi,
  addUserToGroupApi,
  getGroupPermissionsApi,
  assignPermissionsToGroupApi,
  getPermsListWithFilterApi,
  removeUserFromGroupApi,
} from '#/api';
import { useToast, POSITION } from 'vue-toastification';
import { $t } from '@vben/locales';
import type { TreeKey } from 'element-plus/es/components/tree/src/tree.type.mjs';

const toast = useToast();
const data = ref();
// 添加激活的标签页
const activeTab = ref('user');

const userTreeData = ref<UserApi.ElnUser[]>([]);
const userTreeRef = ref<TreeInstance>();
// 添加API树相关数据
const permissionTreeData = ref<UserApi.ElnPermission[]>([]);
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
        if (userTreeRef.value) {
          // Filter out null or undefined userIds to satisfy TreeKey[] type
          const checkedKeys = userTreeData.value.map((item) => item.userId);
          userTreeRef.value.setCheckedKeys(checkedKeys as TreeKey[], false);
        }
      }, 100);
      setTimeout(() => {
        if (permissionTreeRef.value) {
          // Filter out null or undefined userIds to satisfy TreeKey[] type
          const checkedKeys = permissionTreeData.value.map((item) => item.permissionId);
          permissionTreeRef.value.setCheckedKeys(checkedKeys as TreeKey[], false);
        }
      }, 100);
    }
  },

  async onConfirm() {
    let userIds: TreeKey[] = [];
    let permissionIds: TreeKey[] = [];

    // 获取用户权限
    if (userTreeRef.value) {
      const checkedKeys = userTreeRef.value.getCheckedKeys() as TreeKey[] | undefined;
      // 合并已选用户列表，去重
      userIds = Array.from(new Set([...(checkedKeys ?? []), ...((selectedUserList.value ?? []) as TreeKey[])]));
    }

    // 获取权限
    if (permissionTreeRef.value) {
      const checkedKeys = permissionTreeRef.value.getCheckedKeys() as TreeKey[] | undefined;
      // 合并已选权限列表，去重
      permissionIds = Array.from(
        new Set([...(checkedKeys ?? []), ...((selectedPermissionList.value ?? []) as TreeKey[])]),
      );
    }

    if (userIds.length === 0 && permissionIds.length === 0) {
      toast.error($t('ui.notification.no_auth'), {
        timeout: 2000,
        position: POSITION.TOP_CENTER,
      });
      return;
    }

    setLoading(true);
    try {
      // 更新权限，同时提交用户和权限
      // 先获取当前分组下的所有用户
      const currentUserIds = userTreeData.value.map((item) => item.userId);

      // 需要添加的用户
      const usersToAdd = userIds.filter((id) => !currentUserIds.includes(id));
      // 需要移除的用户
      const usersToRemove = currentUserIds.filter((id) => !userIds.includes(id));

      if (usersToAdd.length > 0) {
        await addUserToGroupApi({
          groupId: data.value.row.groupId,
          userId: usersToAdd,
        });
      }

      if (usersToRemove.length > 0) {
        await removeUserFromGroupApi({
          groupId: data.value.row.groupId,
          userId: usersToRemove,
        });
      }

      if (permissionIds.length > 0) {
        await assignPermissionsToGroupApi({
          groupId: data.value.row.groupId,
          permissions: permissionIds,
        });
      }

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
    permissionList.value = undefined;
  },
});

// 递归获取所有节点的 key
const getAllKeys = (data: UserApi.ElnUser[]): number[] => {
  const keys: number[] = [];
  const traverse = (nodes: UserApi.ElnUser[]) => {
    nodes.forEach((node) => {
      keys.push(node.userId as number);
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
  if (userTreeRef.value) {
    const allKeys = getAllKeys(userTreeData.value);
    userTreeRef.value.setCheckedKeys(allKeys);
  }
};

// 取消全选
const uncheckAll = () => {
  if (userTreeRef.value) {
    userTreeRef.value.setCheckedKeys([]);
  }
};

// 添加API树操作方法
// 展开所有API节点
// const expandApiAll = () => {
//   if (permissionTreeRef.value) {
//     const allKeys = getAllApiKeys(permissionTreeData.value);
//     allKeys.forEach((key) => {
//       permissionTreeRef.value?.store?.nodesMap[key]?.expand();
//     });
//   }
// };
//
// // 收起所有API节点
// const collapseApiAll = () => {
//   if (permissionTreeRef.value) {
//     const allKeys = getAllApiKeys(permissionTreeData.value);
//     allKeys.forEach((key) => {
//       permissionTreeRef.value?.store?.nodesMap[key]?.collapse();
//     });
//   }
// };

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
const getAllApiKeys = (data: UserApi.ElnPermission[]): number[] => {
  const keys: number[] = [];
  const traverse = (nodes: UserApi.ElnPermission[]) => {
    nodes.forEach((node) => {
      if (node.permissionId !== undefined && node.permissionId !== null) {
        keys.push(node.permissionId);
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

const permissionList = ref<UserApi.ElnPermission[]>();
const selectedPermissionList = ref<TreeKey[]>();

const getAllPermissionsByFilter = (query: string) => {
  if (query && query.length >= 3) {
    setLoading(true);
    setTimeout(() => {
      getPermsListWithFilterApi({ query: query })
        .then((res) => {
          permissionList.value = res?.items;
          setLoading(false);
        })
        .catch(() => {
          permissionList.value = [];
          setLoading(false);
        });
    }, 200);
  } else {
    permissionList.value = [];
  }
};
</script>

<template>
  <Drawer :title="$t('page.system.role.button.auth')">
    <el-tabs v-model="activeTab" class="mb-4" type="border-card">
      <!-- 用户列表标签页 -->
      <el-tab-pane :label="$t('admin.group.drawer.userList')" name="user">
        <div class="flex flex-col gap-4">
          <el-select
            v-model="selectedUserList"
            placeholder="Enter at least 3 characters"
            multiple
            clearable
            filterable
            remote
            :remote-method="getAllUsersByFilter"
          >
            <el-option v-for="item in userList" :key="item.userId" :label="item.username" :value="item.userId">
              <div class="option-row flex items-center">
                <span class="ml-2 text-xs text-gray-500" style="width: 60px; min-width: 60px; display: inline-block">{{
                  item.userId
                }}</span>
                <span class="flex items-center justify-center" style="width: 24px; min-width: 24px">
                  <el-avatar :size="14" :src="item.avatar" />
                </span>
                <span class="flex-1 pl-2 text-left" style="min-width: 80px">{{ item.username }}</span>
                <span class="ml-2 text-xs text-gray-500">{{ item.email }} </span>
              </div>
            </el-option>
          </el-select>
          <div class="flex gap-2">
            <el-button @click="checkAll">{{ $t('ui.tree.select_all') }}</el-button>
            <el-button @click="uncheckAll">{{ $t('ui.tree.unselect_all') }}</el-button>
          </div>

          <el-tree
            ref="userTreeRef"
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
      <el-tab-pane :label="$t('admin.group.drawer.addPermissions')" name="permission">
        <div class="flex flex-col gap-4">
          <el-select
            v-model="selectedPermissionList"
            placeholder="Enter at least 3 characters"
            multiple
            clearable
            filterable
            remote
            :remote-method="getAllPermissionsByFilter"
          >
            <el-option
              v-for="item in permissionList"
              :key="item.permissionId"
              :label="item.permissionName"
              :value="item.permissionId"
            >
              <div class="option-row flex items-center">
                <span class="ml-2 text-xs text-gray-500" style="width: 60px; min-width: 60px; display: inline-block">{{
                  item.permissionId
                }}</span>
                <span class="flex-1 pl-2 text-left" style="min-width: 80px">{{ item.permissionName }}</span>
                <span class="ml-2 text-xs text-gray-500">{{ item.description }} </span>
              </div>
            </el-option>
          </el-select>

          <div class="flex gap-2">
            <!--            <el-button @click="expandApiAll">{{ $t('ui.tree.expand_all') }}</el-button>-->
            <!--            <el-button @click="collapseApiAll">{{ $t('ui.tree.collapse_all') }}</el-button>-->
            <el-button @click="checkApiAll">{{ $t('ui.tree.select_all') }}</el-button>
            <el-button @click="uncheckApiAll">{{ $t('ui.tree.unselect_all') }}</el-button>
          </div>

          <el-tree
            ref="permissionTreeRef"
            :data="permissionTreeData"
            show-checkbox
            node-key="permissionId"
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
</style>
