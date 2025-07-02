<script setup lang="ts">
import { ElMessage, ElNotification, type FormInstance } from 'element-plus';
import { onMounted, ref } from 'vue';
import { getKetcher } from './utils';
import type { Ketcher } from 'ketcher-core';

const projectFormRef = ref<FormInstance>();
const reactionSmiles = ref<string | null>(null);
const isLoading = ref<boolean>(false);

onMounted(() => {
  document.addEventListener('DOMContentLoaded', () => {
    const ketcherFrame = document.getElementById(
      'ketcher-js-editor',
    ) as HTMLIFrameElement | null;

    if (ketcherFrame) {
      ketcherFrame.addEventListener('load', () => {
        const ketcher = <Ketcher>getKetcher();
        if (ketcher) {
          // 正常使用
          console.log('ketcher is ready');
        }
      });
    }
  });
});

async function saveNewReactionNoteToDatabase(formEl: FormInstance | undefined) {
  if (!formEl) {
    ElMessage.error('表单未正确初始化');
    return;
  }

  const ketcher = getKetcher();

  if (!ketcher) {
    ElMessage.error('化学结构编辑器未加载完成');
    return;
  }

  isLoading.value = true;

  try {
    const smiles = await ketcher.getSmiles(true); // 假设返回 Promise<string>
    reactionSmiles.value = smiles;
    console.log('获取到的 SMILES:', smiles);

    // 此处可添加保存到后端的逻辑
    // await saveReactionToServer(smiles);

    ElNotification.success('结构已成功获取');
  } catch (error) {
    console.error('获取 SMILES 失败:', error);
    ElNotification.error('获取结构失败，请重试');
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <div class="p-5">
    <ElCard title="Create a reaction">
      <el-form ref="projectFormRef" label-width="120px">
        <div id="marvin-js" style="margin-top: 10px">
          <iframe
            id="ketcher-js-editor"
            src="/static/ketcher/index.html"
            width="100%"
            height="450px"
          ></iframe>
        </div>
        <div style="margin-top: 10px">
          <ElButton
            type="primary"
            @click="saveNewReactionNoteToDatabase(projectFormRef)"
          >
            Save Reaction
          </ElButton>
        </div>
      </el-form>
    </ElCard>
  </div>
</template>

<style scoped></style>
