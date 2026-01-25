<script setup lang="ts">
import { ElMessage, ElNotification, type FormInstance } from 'element-plus';
import { onMounted, ref, onUnmounted } from 'vue';
import { getKetcher } from '#/utils';
import type { Ketcher } from 'ketcher-core';
import CompoundSubmenu from '#/views/eln/components/CompoundSubmenu.vue';
import BasicInfoForm from '#/views/eln/components/BasicInfoForm.vue';
import { elnStore } from '#/store';

const store = elnStore();
const projectFormRef = ref<FormInstance>();
const reactionSmiles = ref<string | null>(null);
const isLoading = ref<boolean>(false);

// Responsive screen size detection for future use
const windowWidth = ref(window.innerWidth);

const handleResize = () => {
  windowWidth.value = window.innerWidth;
};

onMounted(() => {
  // Setup resize listener for responsive behavior
  window.addEventListener('resize', handleResize);
  handleResize();

  // Setup Ketcher editor
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

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});

const imgUrl = ref();

async function saveNewReactionNoteToDatabase(formEl: FormInstance | undefined) {
  if (!formEl) {
    ElMessage.error('Form is not properly initialized');
    return;
  }

  const ketcher = getKetcher();
  if (!ketcher) {
    ElMessage.error('Chemical structure editor is not loaded');
    return;
  }

  if (!ketcher?.containsReaction()) {
    ElNotification.error("Don't contain reaction");
    return;
  }

  isLoading.value = true;
  const reactionId = ref<number | null>(null);
  try {
    const smiles = await ketcher.getSmiles(true); 
    reactionSmiles.value = smiles;
    console.log('获取到的 SMILES:', smiles);

    // Save reaction to server
    const response = await store.saveRxnToServer(reactionSmiles.value);
    if (response) {
      ElNotification.success('Reaction saved successfully');
      imgUrl.value = 'data:image/svg+xml;base64,' + response.imageSvg;
      reactionId.value = response.reactionId;
      store.formData = response.compounds;
    } else {
      ElNotification.error('Failed to save reaction');
    }
  } catch (error) {
    console.error('Failed to get SMILES:', error);
    ElNotification.error('Failed to get structure, please try again');
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <div class="p-5">
    <!-- Basic Information Form -->
    <BasicInfoForm />

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
    <ElCard class="mb-5">
      <div style="display: flex; justify-content: center; align-items: center; min-height: 150px">
        <ElImage :src="imgUrl" alt="rxnImg" style="max-height: 300px; display: block" />
      </div>
    </ElCard>
    
    <!-- Compounds Submenu -->
    <CompoundSubmenu />
  </div>
</template>

