<template>
  <div v-for="(item, idx) in compoundsToDisplay" :key="item.reagentId || idx" class="mb-4">
    <div 
      class="bg-card text-card-foreground border-border rounded-xl border shadow-sm card-hover reagent-card-container"
      :style="{ borderLeft: `6px solid ${getRoleColor(item.reagentRole)}` }"
    >
      <!-- Card Header -->
      <div class="flex flex-col gap-y-1.5 p-5 border-b border-border/50">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <el-tag
              :style="{ backgroundColor: getRoleColor(item.reagentRole), borderColor: getRoleColor(item.reagentRole) }"
              effect="dark"
              size="small"
            >
              {{ item.reagentRole }}
            </el-tag>
            <el-popover v-if="item.reagentRole !== 'product'" placement="right" :width="100" trigger="click">
              <template #reference>
                <el-button 
                  v-if="!store.isReadonly" 
                  circle
                  size="small"
                  type="default"
                >
                  <Edit :size="16" />
                </el-button>
              </template>
              <el-select v-model="item.reagentRole" size="small" :disabled="store.isReadonly">
                <el-option v-for="role in compoundRole" :key="role.value" :label="role.value" :value="role.value" />
              </el-select>
            </el-popover>
          </div>
          <el-text size="small" type="info">#{{ (props.index !== undefined ? props.index : idx) + 1 }}</el-text>
        </div>
      </div>

      <!-- Card Content -->
      <div class="p-6 pt-0">
        <el-row :gutter="24">
          <!-- Compound Information -->
          <el-col :xs="24" :sm="24" :md="16" :lg="16" :xl="16">
            <el-form
              :model="item"
              label-width="120px"
              class="compound-info-form"
            >
              <el-row :gutter="16">
                <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
                  <el-form-item label="Compound Name" class="form-item-inline">
                    <el-text size="small" class="break-words">{{ item.reagentName }}</el-text>
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
                  <el-form-item label="SMILES" class="form-item-inline">
                    <el-text size="small" class="break-words" style="font-family: monospace;">{{ item.reagentSmiles }}</el-text>
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-row :gutter="16" class="mt-4">
                <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
                  <el-form-item label="Formula" class="form-item-inline">
                    <el-text size="small">{{ item.formula }}</el-text>
                  </el-form-item>
                </el-col>
                <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
                  <el-form-item label="Mol.Wt" class="form-item-inline">
                    <el-text size="small">{{ item.mw }}</el-text>
                  </el-form-item>
                </el-col>
                <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
                  <el-form-item label="Chiral" class="form-item-inline">
                    <el-text size="small">{{ item.isChiral === 1 ? 'Yes' : 'No' }}</el-text>
                  </el-form-item>
                </el-col>
                <el-col v-if="item.isChiral === 1" :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
                  <el-form-item label="Stereo Centers" class="form-item-inline">
                    <el-text size="small">{{ item.stereoCentersCnt }}</el-text>
                  </el-form-item>
                </el-col>
              </el-row>

              <el-form-item v-if="item.reagentRole === 'reactant'" label="Limiting" class="form-item-inline mt-4">
                <el-switch
                  v-model="item.isLimiting"
                  :active-value="1"
                  :inactive-value="0"
                  :disabled="store.isReadonly"
                  @change="handleSwitchChange(props.index !== undefined ? props.index : idx, item.isLimiting)"
                />
              </el-form-item>

              <el-form-item v-if="item.isChiral === 1" label="Chiral Descriptor" class="form-item-inline mt-4">
                <el-tag size="small" type="info">{{ item.chiralDescriptor }}</el-tag>
              </el-form-item>
            </el-form>
          </el-col>

          <!-- Compound Image -->
          <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
            <div class="flex justify-center md:justify-end">
              <div class="relative">
                <el-image
                  :src='`data:image/svg+xml;base64,${item.reagentImg}`'
                  class="w-32 h-32 rounded-lg border border-border object-cover"
                  :preview-src-list="[item.reagentImg]"
                />
              </div>
            </div>
          </el-col>
        </el-row>

        <!-- Instructions -->
        <el-alert
          :title="item.reagentRole === 'solvent' ? 'Purity - Volume - Density' : 'Limiting - Equiv - Purity - Quantity - Density/Conc. Use -1 to indicate excessive amount.'"
          type="info"
          :closable="false"
          show-icon
          class="mt-6"
        />

        <!-- Input Fields -->
        <div class="mt-6">
          <el-form
            :model="item"
            label-width="120px"
            class="reagent-form"
          >
            <el-row :gutter="16">
              <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
                <el-form-item label="Equiv." class="form-item-inline">
                  <el-input
                    v-model="item.eq"
                    size="small"
                    :disabled="item.isLimiting === 1 || store.isReadonly"
                    @input="calcMolesByEqLocal(item)"
                    class="w-full form-input"
                  />
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
                <el-form-item label="Purity (%)" class="form-item-inline">
                  <el-tooltip content="Purity of the Reactant/Product/Reagent" placement="top">
                    <el-input
                      v-model="item.purity"
                      size="small"
                      :disabled="store.isReadonly"
                      @input="calcQuantityByMoles(item)"
                      class="w-full"
                    />
                  </el-tooltip>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
                <el-form-item label="Quantity" class="form-item-inline">
                  <el-tooltip content="Mass of the Reactant/Product/Reagent" placement="top">
                    <el-input
                      v-model="item.quantity"
                      size="small"
                      :disabled="store.isReadonly"
                      @input="calcMoles(item)"
                      class="w-full"
                    >
                      <template #append>
                        <el-select v-model="item.quantityUnit" size="small" :disabled="store.isReadonly" @change="calcMoles(item)">
                          <el-option v-for="unit in unitList" :key="unit.value" :label="unit.value" :value="unit.value" />
                        </el-select>
                      </template>
                    </el-input>
                  </el-tooltip>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
                <el-form-item label="Moles" class="form-item-inline">
                  <el-input v-model="item.moles" size="small" :disabled="store.isReadonly" readonly class="w-full">
                    <template #append>
                      <span class="inline-flex items-center px-2 py-1 text-xs font-medium bg-secondary text-secondary-foreground rounded">
                        {{ item.molesUnit }}
                      </span>
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
                <el-form-item label="Density (g/mL)" class="form-item-inline">
                  <el-input
                    v-model="item.density"
                    size="small"
                    :disabled="store.isReadonly"
                    @input="calcVolume(item)"
                    class="w-full"
                  />
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
                <el-form-item label="Conc. (mol/L)" class="form-item-inline">
                  <el-input
                    v-model="item.concentration"
                    size="small"
                    :disabled="store.isReadonly"
                    @input="calcVolume(item)"
                    class="w-full"
                  />
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
                <el-form-item label="Volume" class="form-item-inline">
                  <el-input
                    v-model="item.volume"
                    size="small"
                    :disabled="store.isReadonly"
                    @input="calcMolesByVolume(item)"
                    class="w-full"
                  >
                    <template #append>
                      <el-select v-model="item.volumeUnit" size="small" :disabled="store.isReadonly" @change="calcMolesByVolume(item)">
                        <el-option v-for="unit in volumeUnitList" :key="unit.value" :label="unit.value" :value="unit.value" />
                      </el-select>
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
                <el-form-item label="CAS#" class="form-item-inline">
                  <el-input v-model="item.cas" size="small" :disabled="store.isReadonly" class="w-full" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </div>

        <!-- Action Buttons -->
        <div v-if="!store.isReadonly" class="mt-6 flex justify-end space-x-2">
          <el-button
            type="primary"
            @click="saveAdditionalInfoToDB(reagentFormRef[props.index !== undefined ? props.index : idx], item)"
            :disabled="store.isReadonly"
          >
            <Save :size="16" class="mr-1" />
            Save
          </el-button>
          <el-button
            v-if="showDeleteOption(item.reagentRole)"
            type="danger"
            @click="deleteRow(String(item.reagentId))"
            :disabled="store.isReadonly"
          >
            <Trash2 :size="16" class="mr-1" />
            Delete
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { Edit, Save, Trash2 } from 'lucide-vue-next';
import { elnStore } from '#/store';
import type { TableDataStruct } from '#/types';

interface Props {
  compound?: TableDataStruct;
  index?: number;
}

const props = withDefaults(defineProps<Props>(), {
  compound: undefined,
  index: undefined,
});

// Store
const store = elnStore();

// Get compounds to display
const compoundsToDisplay = computed(() => {
  if (props.compound) {
    return [props.compound];
  }
  return store.formData || [];
});

// Refs
const reagentFormRef = ref();

// Constants
const compoundRole = [
  { value: 'catalyst' },
  { value: 'reagent' },
  { value: 'solvent' },
  { value: 'base' },
  { value: 'ligand' },
  { value: 'additive' },
  { value: 'other' },
];

const unitList = [
  { value: 'mg' },
  { value: 'g' },
  { value: 'kg' },
];

const volumeUnitList = [
  { value: 'mL' },
  { value: 'L' },
];

// Methods
const getRoleColor = (role: string): string => {
  const colorMap: Record<string, string> = {
    reactant: 'rgba(7,203,75,0.91)',
    product: 'rgb(11,61,232)',
    solvent: 'rgba(190,175,9,0.98)',
    base: 'rgba(5,154,119,0.93)',
    catalyst: 'rgba(241,78,7,0.98)',
    additives: 'rgba(128,7,241,0.8)',
  };
  return colorMap[role] || 'rgb(2,51,44)';
};

const showDeleteOption = (role: string): boolean => {
  return !(role === 'reactant' || role === 'product');
};

// Placeholder methods - these should be implemented based on your store methods
const handleSwitchChange = (_index: number, _value: number) => {
  // Implementation needed
};

const calcMolesByEqLocal = (_item: any) => {
  // Implementation needed
};

const calcQuantityByMoles = (_item: any) => {
  // Implementation needed
};

const calcMoles = (_item: any) => {
  // Implementation needed
};

const calcVolume = (_item: any) => {
  // Implementation needed
};

const calcMolesByVolume = (_item: any) => {
  // Implementation needed
};

const saveAdditionalInfoToDB = (_formRef: any, _item: any) => {
  // Implementation needed
};

const deleteRow = (_reagentId: string) => {
  // Implementation needed
};
</script>
<style scoped lang="scss">
// vben-admin optimized styles - using Tailwind CSS classes in template
// Custom styles only for specific vben-admin theming if needed

// Ensure proper spacing and transitions
.reagent-card-container {
  transition: all 200ms ease-in-out;
}

// Custom role color styling
.role-badge {
  display: inline-flex;
  align-items: center;
  padding: 0.125rem 0.625rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
  color: white;
}

// Button hover effects
.action-button {
  transition: all 200ms ease-in-out;
  
  &:hover {
    transform: scale(1.05);
  }
  
  &:active {
    transform: scale(0.95);
  }
}

// Form input focus states
.form-input {
  transition: all 200ms ease-in-out;
  
  &:focus {
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
    border-color: #3b82f6;
  }
}

// Card hover effects
.card-hover {
  transition: all 300ms ease-in-out;
  
  &:hover {
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
    transform: translateY(-0.25rem);
  }
}

// Form item inline layout - label and input on same line
.reagent-form,
.compound-info-form {
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

// Responsive adjustments for form items
@media (max-width: 768px) {
  .reagent-form,
  .compound-info-form {
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

// Responsive adjustments
@media (max-width: 640px) {
  .mobile-stack {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .mobile-full {
    width: 100%;
  }
}
</style>
