<script setup lang="ts">
import { ElMessage, ElNotification, type FormInstance } from 'element-plus';
import { onMounted, ref } from 'vue';
import { generateImgUrl, getKetcher } from '#/utils';
import type { Ketcher } from 'ketcher-core';
import ReagentCard from '#/views/eln/components/ReagentCard.vue';

const projectFormRef = ref<FormInstance>();
const reactionSmiles = ref<string | null>(null);
const isLoading = ref<boolean>(false);

onMounted(() => {
  document.addEventListener('DOMContentLoaded', () => {
    const ketcherFrame = document.getElementById('ketcher-js-editor') as HTMLIFrameElement | null;

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

const imgUrl = ref();

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

  if (!ketcher?.containsReaction()) {
    ElNotification.error("Don't contain reaction");
    return;
  }

  isLoading.value = true;

  try {
    const smiles = await ketcher.getSmiles(true); // 假设返回 Promise<string>
    reactionSmiles.value = smiles;
    console.log('获取到的 SMILES:', smiles);
    ketcher
      .calculate({
        properties: ['molecular-weight', 'gross'],
        struct: reactionSmiles.value,
      })
      .then((res) => {
        console.log(res.gross);
      });

    // 此处可添加保存到后端的逻辑
    // await saveReactionToServer(smiles);
    imgUrl.value = await generateImgUrl(ketcher, reactionSmiles.value);
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
          <iframe id="ketcher-js-editor" src="/static/ketcher/index.html" width="100%" height="450px"></iframe>
        </div>
        <div style="display: flex; justify-content: flex-end; margin-top: 10px">
          <ElButton type="primary" @click="saveNewReactionNoteToDatabase(projectFormRef)"> Save Reaction </ElButton>
        </div>
      </el-form>
    </ElCard>
    <ElCard>
      <div style="display: flex; justify-content: center; align-items: center; min-height: 150px">
        <ElImage :src="imgUrl" alt="rxnImg" style="max-height: 300px; display: block" />
      </div>
    </ElCard>
    <reagent-card :form-data="{ reagentName: 'lcx' }" :reactant-table-data="{ reagentName: 'lcx' }" />
  </div>
</template>

<style scoped></style>
