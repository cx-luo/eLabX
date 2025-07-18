import type { Recordable, UserInfo } from '@vben/types';

import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { DEFAULT_HOME_PATH, LOGIN_PATH } from '@vben/constants';
import { resetAllStores, useAccessStore, useUserStore } from '@vben/stores';

import { ElNotification } from 'element-plus';
import { defineStore } from 'pinia';
import { md5 } from 'js-md5';
import { getAccessCodesApi, getUserInfoApi, loginApi, logoutApi } from '#/api';
import { $t } from '#/locales';

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore();
  const userStore = useUserStore();
  const router = useRouter();

  const loginLoading = ref(false);

  /**
   * 异步处理登录操作
   * Asynchronously handle the login process
   * @param params 登录表单数据
   * @param onSuccess
   */
  async function authLogin(
    params: Recordable<any>,
    onSuccess?: () => Promise<void> | void,
  ) {
    // 异步处理用户登录操作并获取 accessToken
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;
      const user = {
        username: params.username,
        password: md5(params.password),
      };

      await loginApi(user)
        .then(async (response) => {
          // 如果成功获取到 accessToken
          if (response) {
            // 将 accessToken 存储到 accessStore 中
            accessStore.setAccessToken(response.accessToken);
            // 获取用户信息并存储到 accessStore 中
            const [fetchUserInfoResult, accessCodes] = await Promise.all([
              fetchUserInfo(),
              getAccessCodesApi(),
            ]);

            userInfo = fetchUserInfoResult;

            userStore.setUserInfo(userInfo);
            accessStore.setAccessCodes(accessCodes);

            if (accessStore.loginExpired) {
              accessStore.setLoginExpired(false);
            } else {
              onSuccess
                ? await onSuccess?.()
                : await router.push(userInfo.homePath || DEFAULT_HOME_PATH);
            }

            if (userInfo?.realName) {
              ElNotification({
                message: `${$t('authentication.loginSuccessDesc')}:${userInfo?.realName}`,
                title: $t('authentication.loginSuccess'),
                type: 'success',
              });
            }
          }
        })
        .catch((e) => {
          console.log('login error', e);
        });
    } finally {
      loginLoading.value = false;
    }

    return {
      userInfo,
    };
  }

  async function logout(redirect: boolean = true) {
    try {
      await logoutApi();
    } catch {
      // 不做任何处理
    }
    resetAllStores();
    accessStore.setLoginExpired(false);

    // 回登录页带上当前路由地址
    await router.replace({
      path: LOGIN_PATH,
      query: redirect
        ? {
            redirect: encodeURIComponent(router.currentRoute.value.fullPath),
          }
        : {},
    });

    router.go(0);
  }

  async function fetchUserInfo() {
    let userInfo: null | UserInfo = null;
    userInfo = await getUserInfoApi();
    userStore.setUserInfo(userInfo);
    return userInfo;
  }

  function $reset() {
    loginLoading.value = false;
  }

  return {
    $reset,
    authLogin,
    fetchUserInfo,
    loginLoading,
    logout,
  };
});
