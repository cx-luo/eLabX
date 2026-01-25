<template>
  <ElCard class="mb-5">
    <template #header>
      <h3 class="text-lg font-semibold">Basic information</h3>
    </template>
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
      class="basic-info-form"
    >
      <el-row :gutter="20">
        <!-- Left Column -->
        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="ProjectName" prop="projectName" required>
            <el-input v-model="formData.projectName" placeholder="Enter project name" />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="Batch" prop="batch" required>
            <el-input v-model="formData.batch" placeholder="Enter batch" />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="StepID" prop="stepID" required>
            <el-input v-model="formData.stepID" placeholder="Enter step ID" />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="ReactionType" prop="reactionType">
            <el-select v-model="formData.reactionType" placeholder="Select" style="width: 100%">
              <el-option
                v-for="option in reactionTypeOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="StartDate" prop="startDate" required>
            <el-date-picker
              v-model="formData.startDate"
              type="date"
              placeholder="Select date"
              style="width: 100%"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="CreationDate" prop="creationDate">
            <el-date-picker
              v-model="formData.creationDate"
              type="date"
              placeholder="Select date"
              style="width: 100%"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="AuthorID" prop="authorID" required>
            <el-input v-model="formData.authorID" placeholder="Enter author ID" />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="AuthorName" prop="authorName" required>
            <el-input v-model="formData.authorName" placeholder="Enter author name" />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="WitnessID" prop="witnessID" required>
            <el-input v-model="formData.witnessID" placeholder="Enter witness ID" />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="WitnessName" prop="witnessName">
            <el-input v-model="formData.witnessName" placeholder="Enter witness name" />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-form-item label="DOI" prop="doi">
            <el-input v-model="formData.doi" placeholder="Enter DOI" />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
          <el-form-item label="comment" prop="comment">
            <el-input
              v-model="formData.comment"
              type="textarea"
              :rows="1"
              placeholder="Enter comment"
            />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
          <el-form-item label="Reference" prop="reference">
            <el-input
              v-model="formData.reference"
              type="textarea"
              :rows="1"
              placeholder="Enter reference"
            />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
          <el-form-item>
            <el-button type="primary" @click="handleSave">
              <FileText class="w-4 h-4 mr-2" />
              Save Basic Info
            </el-button>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </ElCard>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { ElMessage, ElNotification, type FormInstance } from 'element-plus';
import { FileText } from 'lucide-vue-next';
import { elnStore } from '#/store';

const store = elnStore();
const formRef = ref<FormInstance>();

// Form Data
const formData = ref({
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

// Form Rules
const rules = {
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
async function handleSave() {
  if (!formRef.value) {
    ElMessage.error('Form is not properly initialized');
    return;
  }

  formRef.value.validate((valid) => {
    if (valid) {
      // TODO: Add API call to save basic info
      console.log('Basic Info:', formData.value);
      ElNotification.success('Basic information saved successfully');
    } else {
      ElMessage.error('Please fill in all required fields');
    }
  });
}

// Expose form data and ref for parent component
defineExpose({
  formData,
  formRef,
  handleSave,
});
</script>

<style scoped lang="scss">
.basic-info-form {
  :deep(.el-form-item) {
    margin-bottom: 18px;
    display: flex;
    flex-direction: row;
    align-items: center;
    
    .el-form-item__label {
      font-weight: bold;
      line-height: 32px;
      padding-bottom: 0;
      padding-right: 12px;
      text-align: left;
      flex-shrink: 0;
      white-space: nowrap;
    }
    
    .el-form-item__content {
      line-height: 32px;
      flex: 1;
      margin-left: 0 !important;
    }
  }
}

// Responsive adjustments
@media (max-width: 768px) {
  .basic-info-form {
    :deep(.el-form-item) {
      margin-bottom: 16px;
      flex-direction: column;
      align-items: flex-start;
      
      .el-form-item__label {
        width: 100%;
        padding-bottom: 8px;
        padding-right: 0;
      }
      
      .el-form-item__content {
        width: 100%;
      }
    }
  }
}
</style>

