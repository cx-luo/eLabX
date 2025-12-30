<script setup lang="ts">
import {computed, onMounted, ref} from 'vue';
import {useVbenForm, type VbenFormProps} from '#/adapter/form';
import {preferences, updatePreferences} from '@vben/preferences';
import {useToast} from 'vue-toastification';
import {getProjectList} from '#/api';
import {ElButton, ElCard, ElDivider, ElSwitch, ElTag} from 'element-plus';
import {Save, Settings} from 'lucide-vue-next';

const toast = useToast();
const columnsOptions = ref([])
// Reactive data
const projectList = ref<any[]>([]);
const loading = ref(false);
const selectedProjects = ref<string[]>([]);
const sidebarVisibility = ref<Record<string, boolean>>({});

// Form configuration for project selection
const projectFormOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: false,
  submitOnEnter: false,
  wrapperClass: 'grid-cols-1',
  schema: [
    {
      component: 'Select',
      fieldName: 'search',
      label: 'Search Projects',
      componentProps: {
        placeholder: 'Type to search projects...',
        allowClear: true,
        prefixIcon: 'lucide:search',
        filterable: true,
        multiple: true,
        options: columnsOptions,
      },
    },
  ],
};


// Computed properties
const filteredProjects = computed(() => {
  const searchTerm = projectFormValues.value?.search?.toLowerCase() || '';
  return projectList.value.filter(project =>
      project.projectName?.toLowerCase().includes(searchTerm) ||
      project.description?.toLowerCase().includes(searchTerm)
  );
});

const selectedProjectCount = computed(() => selectedProjects.value.length);

// Form instances
const [ProjectForm] = useVbenForm(projectFormOptions);

// Form values
const projectFormValues = ref<{ search?: string }>({});

// Methods
const loadProjects = async () => {
  try {
    loading.value = true;
    const response = await getProjectList({
      page: 1,
      pageSize: 100,
    });
    projectList.value = response.data || [];
  } catch (error) {
    console.error('Failed to load projects:', error);
    toast.error('Failed to load projects');
  } finally {
    loading.value = false;
  }
};


const saveUserPreferences = async () => {
  try {
    // Update sidebar preferences
    const updatedPreferences = {
      ...preferences,
      sidebar: {
        ...preferences.sidebar,
        selectedProjects: selectedProjects.value,
        sectionVisibility: sidebarVisibility.value,
      } as any,
    };

    await updatePreferences(updatedPreferences);
    toast.success('Preferences saved successfully');
  } catch (error) {
    console.error('Failed to save preferences:', error);
    toast.error('Failed to save preferences');
  }
};

const toggleProjectSelection = (projectId: string) => {
  const index = selectedProjects.value.indexOf(projectId);
  if (index > -1) {
    selectedProjects.value.splice(index, 1);
  } else {
    selectedProjects.value.push(projectId);
  }
};

const selectAllProjects = () => {
  selectedProjects.value = projectList.value.map(project => project.projectId);
};

const clearAllProjects = () => {
  selectedProjects.value = [];
};

const resetToDefaults = () => {
  selectedProjects.value = [];
};

// Lifecycle
onMounted(async () => {
  await loadProjects();
});

// Watch for form changes
const handleProjectFormChange = (values: { search?: string }) => {
  projectFormValues.value = values;
  console.log(values.search)
};
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-foreground">Project Selection & Sidebar Configuration</h1>
        <p class="text-muted-foreground mt-1">Choose which projects to display and configure sidebar visibility</p>
      </div>
      <div class="flex items-center gap-2">
        <ElButton @click="resetToDefaults" type="default" size="small">
          <Settings class="w-4 h-4 mr-2"/>
          Reset to Defaults
        </ElButton>
        <ElButton @click="saveUserPreferences" type="primary" size="small">
          <Save class="w-4 h-4 mr-2"/>
          Save Preferences
        </ElButton>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Project Selection -->
      <ElCard class="h-fit">
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Project Selection</h3>
            <div class="flex items-center gap-2">
              <ElButton @click="selectAllProjects" size="small" type="primary" plain>
                Select All
              </ElButton>
              <ElButton @click="clearAllProjects" size="small" type="default" plain>
                Clear All
              </ElButton>
            </div>
          </div>
        </template>

        <!-- Project Search Form -->
        <div class="mb-4">
          <ProjectForm @input="handleProjectFormChange"/>
        </div>

        <!-- Project List -->
        <div class="space-y-3 max-h-96 overflow-y-auto">
          <div v-if="loading" class="text-center py-8 text-muted-foreground">
            Loading projects...
          </div>
          <div v-else-if="filteredProjects.length === 0" class="text-center py-8 text-muted-foreground">
            No projects found
          </div>
          <div
              v-else
              v-for="project in filteredProjects"
              :key="project.projectId"
              class="flex items-center justify-between p-3 border border-border rounded-lg hover:bg-muted/50 transition-colors"
              :class="{ 'bg-primary/10 border-primary': selectedProjects.includes(project.projectId) }"
          >
            <div class="flex-1 min-w-0">
              <h4 class="font-medium text-foreground truncate">{{ project.projectName }}</h4>
              <p class="text-sm text-muted-foreground truncate">{{ project.description }}</p>
              <div class="flex items-center gap-2 mt-1">
                <ElTag size="small" :type="project.status === 1 ? 'success' : 'info'">
                  {{ project.status === 1 ? 'Active' : 'Inactive' }}
                </ElTag>
                <span class="text-xs text-muted-foreground">
                  Created by: {{ project.createdBy }}
                </span>
              </div>
            </div>
            <ElSwitch
                :model-value="selectedProjects.includes(project.projectId)"
                @change="toggleProjectSelection(project.projectId)"
                size="small"
            />
          </div>
        </div>

        <!-- Selection Summary -->
        <ElDivider/>
        <div class="flex items-center justify-between text-sm">
          <span class="text-muted-foreground">
            {{ selectedProjectCount }} of {{ projectList.length }} projects selected
          </span>
        </div>
      </ElCard>
    </div>
  </div>
</template>

<style scoped lang="scss">
// vben-admin optimized styles
.project-card {
  transition: all 200ms ease-in-out;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
}

.section-item {
  transition: all 200ms ease-in-out;

  &:hover {
    background-color: rgba(0, 0, 0, 0.02);
  }
}

// Custom scrollbar
.overflow-y-auto {
  scrollbar-width: thin;
  scrollbar-color: rgba(0, 0, 0, 0.2) transparent;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 3px;
  }

  &::-webkit-scrollbar-thumb:hover {
    background-color: rgba(0, 0, 0, 0.3);
  }
}
</style>