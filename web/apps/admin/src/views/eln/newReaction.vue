<script setup lang="ts">
import { ElMessage, ElNotification, ElDatePicker, type FormInstance } from 'element-plus';
import { onMounted, ref, computed, onUnmounted } from 'vue';
import { generateImgUrl, getKetcher } from '#/utils';
import type { Ketcher } from 'ketcher-core';
import ReagentCard from '#/views/eln/components/ReagentCard.vue';
import { FileText } from 'lucide-vue-next';
import { elnStore } from '#/store';

const store = elnStore();
const projectFormRef = ref<FormInstance>();
const basicInfoFormRef = ref<FormInstance>();
const reactionSmiles = ref<string | null>(null);
const isLoading = ref<boolean>(false);

// Responsive screen size detection for future use
const windowWidth = ref(window.innerWidth);

const handleResize = () => {
  windowWidth.value = window.innerWidth;
};

// Basic Information Form Data
const basicInfoForm = ref({
  projectName: '',
  batch: '',
  stepID: '',
  reactionType: '',
  comment: '',
  startDate: '',
  authorID: '',
  authorName: '',
  reference: '',
  creationDate: '',
  witnessID: '',
  witnessName: '',
  doi: '',
});

// Basic Information Form Rules
const basicInfoRules = {
  projectName: [{ required: true, message: 'Please input project name', trigger: 'blur' }],
  batch: [{ required: true, message: 'Please input batch', trigger: 'blur' }],
  stepID: [{ required: true, message: 'Please input step ID', trigger: 'blur' }],
  startDate: [{ required: true, message: 'Please select start date', trigger: 'change' }],
  authorID: [{ required: true, message: 'Please input author ID', trigger: 'blur' }],
  authorName: [{ required: true, message: 'Please input author name', trigger: 'blur' }],
  witnessID: [{ required: true, message: 'Please input witness ID', trigger: 'blur' }],
};

// Reaction Type Options from store
const reactionTypeOptions = computed(() => 
  store.reactionTypeOptions.map(type => ({ label: type, value: type }))
);

// Save Basic Info Function
async function saveBasicInfo(formEl: FormInstance | undefined) {
  if (!formEl) {
    ElMessage.error('Form is not properly initialized');
    return;
  }

  formEl.validate((valid) => {
    if (valid) {
      // TODO: Add API call to save basic info
      console.log('Basic Info:', basicInfoForm.value);
      ElNotification.success('Basic information saved successfully');
    } else {
      ElMessage.error('Please fill in all required fields');
    }
  });
}

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

  try {
    const smiles = await ketcher.getSmiles(true); 
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
    const response = await store.saveRxnToServer(reactionSmiles.value);
    if (response && response.statusCode === 200) {
      ElNotification.success('Reaction saved successfully');
    } else {
      ElNotification.error('Failed to save reaction');
    }
    
    imgUrl.value = await generateImgUrl(ketcher, reactionSmiles.value);
    ElNotification.success('Structure successfully retrieved');
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
    <ElCard class="mb-5">
      <template #header>
        <h3 class="text-lg font-semibold">Basic information</h3>
      </template>
      <el-form
        ref="basicInfoFormRef"
        :model="basicInfoForm"
        :rules="basicInfoRules"
        label-width="100px"
        label-position="top"
        class="basic-info-form"
      >
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-3 gap-4 sm:gap-5 lg:gap-6">
          <!-- Left Column -->
          <div class="space-y-4">
            <el-form-item label="ProjectName" prop="projectName" required>
              <el-input v-model="basicInfoForm.projectName" placeholder="Enter project name" />
            </el-form-item>
            <el-form-item label="Batch" prop="batch" required>
              <el-input v-model="basicInfoForm.batch" placeholder="Enter batch" />
            </el-form-item>
            <el-form-item label="StepID" prop="stepID" required>
              <el-input v-model="basicInfoForm.stepID" placeholder="Enter step ID" />
            </el-form-item>
            <el-form-item label="ReactionType" prop="reactionType">
              <el-select v-model="basicInfoForm.reactionType" placeholder="Select" style="width: 100%">
                <el-option
                  v-for="option in reactionTypeOptions"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="comment" prop="comment">
              <el-input
                v-model="basicInfoForm.comment"
                type="textarea"
                :rows="1"
                placeholder="Enter comment"
              />
            </el-form-item>
          </div>

          <!-- Middle Column -->
          <div class="space-y-4">
            <el-form-item label="StartDate" prop="startDate" required>
              <el-date-picker
                v-model="basicInfoForm.startDate"
                type="date"
                placeholder="Select date"
                style="width: 100%"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="AuthorID" prop="authorID" required>
              <el-input v-model="basicInfoForm.authorID" placeholder="Enter author ID" />
            </el-form-item>
            <el-form-item label="AuthorName" prop="authorName" required>
              <el-input v-model="basicInfoForm.authorName" placeholder="Enter author name" />
            </el-form-item>
            <el-form-item label="Reference" prop="reference">
              <el-input
                v-model="basicInfoForm.reference"
                type="textarea"
                :rows="1"
                placeholder="Enter reference"
              />
            </el-form-item>
            <div class="flex justify-end mt-6">
              <ElButton type="primary" @click="saveBasicInfo(basicInfoFormRef)">
                <FileText class="w-4 h-4 mr-2" />
                Save Basic Info
              </ElButton>
            </div>
          </div>

          <!-- Right Column -->
          <div class="space-y-4">
            <el-form-item label="CreationDate" prop="creationDate">
              <el-date-picker
                v-model="basicInfoForm.creationDate"
                type="date"
                placeholder="Select date"
                style="width: 100%"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="WitnessID" prop="witnessID" required>
              <el-input v-model="basicInfoForm.witnessID" placeholder="Enter witness ID" />
            </el-form-item>
            <el-form-item label="WitnessName" prop="witnessName">
              <el-input v-model="basicInfoForm.witnessName" placeholder="Enter witness name" />
            </el-form-item>
            <el-form-item label="DOI" prop="doi">
              <el-input v-model="basicInfoForm.doi" placeholder="Enter DOI" />
            </el-form-item>
          </div>
        </div>
      </el-form>
    </ElCard>

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

<style scoped>
/* Responsive form layout */
.basic-info-form {
  width: 100%;
}

/* Mobile devices (default) */
@media (max-width: 639px) {
  .basic-info-form :deep(.el-form-item__label) {
    width: 100% !important;
    text-align: left;
  }
  
  .basic-info-form :deep(.el-form-item) {
    margin-bottom: 18px;
  }
}

/* Tablet devices */
@media (min-width: 640px) and (max-width: 1023px) {
  .basic-info-form {
    max-width: 100%;
  }
  
  .basic-info-form :deep(.el-form-item) {
    margin-bottom: 20px;
  }
}

/* Desktop devices */
@media (min-width: 1024px) {
  .basic-info-form {
    max-width: 100%;
  }
  
  .basic-info-form :deep(.el-form-item) {
    margin-bottom: 22px;
  }
}

/* Large desktop devices */
@media (min-width: 1280px) {
  .basic-info-form {
    max-width: 100%;
  }
}

/* Ensure form items take full width on mobile */
@media (max-width: 639px) {
  .basic-info-form :deep(.el-input),
  .basic-info-form :deep(.el-select),
  .basic-info-form :deep(.el-date-editor) {
    width: 100%;
  }
}
</style>
