<template>
  <div v-for="(item, index) in store.formData" :key="item.reagentId || index" class="mb-4">
    <div 
      class="bg-card text-card-foreground border-border rounded-xl border shadow-sm card-hover reagent-card-container"
      :style="{ borderLeft: `6px solid ${getRoleColor(item.reagentRole)}` }"
    >
      <!-- Card Header -->
      <div class="flex flex-col gap-y-1.5 p-5 border-b border-border/50">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span 
              class="role-badge"
              :style="{ backgroundColor: getRoleColor(item.reagentRole) }"
            >
              {{ item.reagentRole }}
            </span>
            <el-popover v-if="item.reagentRole !== 'product'" placement="right" :width="100" trigger="click">
              <template #reference>
                <button 
                  v-if="!store.isReadonly" 
                  class="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-8 w-8"
                >
                  <el-icon class="h-4 w-4"><Edit /></el-icon>
                </button>
              </template>
              <el-select v-model="item.reagentRole" size="small" :disabled="store.isReadonly">
                <el-option v-for="role in compoundRole" :key="role.value" :label="role.value" :value="role.value" />
              </el-select>
            </el-popover>
          </div>
          <span class="text-sm font-medium text-muted-foreground">#{{ index + 1 }}</span>
        </div>
      </div>

      <!-- Card Content -->
      <div class="p-6 pt-0">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- Compound Information -->
          <div class="lg:col-span-2 space-y-4">
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                  Compound Name
                </label>
                <p class="text-sm text-muted-foreground break-words">{{ item.reagentName }}</p>
              </div>
              <div class="space-y-2">
                <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                  SMILES
                </label>
                <p class="text-sm text-muted-foreground break-words font-mono">{{ item.reagentSmiles }}</p>
              </div>
            </div>
            
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
              <div class="space-y-2">
                <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                  Formula
                </label>
                <p class="text-sm text-muted-foreground">{{ item.formula }}</p>
              </div>
              <div class="space-y-2">
                <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                  Mol.Wt
                </label>
                <p class="text-sm text-muted-foreground">{{ item.mw }}</p>
              </div>
              <div class="space-y-2">
                <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                  Chiral
                </label>
                <p class="text-sm text-muted-foreground">{{ item.isChiral === 1 ? 'Yes' : 'No' }}</p>
              </div>
              <div v-if="item.isChiral === 1" class="space-y-2">
                <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                  Stereo Centers
                </label>
                <p class="text-sm text-muted-foreground">{{ item.stereoCentersCnt }}</p>
              </div>
            </div>

            <div v-if="item.reagentRole === 'reactant'" class="flex items-center space-x-2">
              <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Limiting
              </label>
              <el-switch
                v-model="item.isLimiting"
                :active-value="1"
                :inactive-value="0"
                :disabled="store.isReadonly"
                @change="handleSwitchChange(index, item.isLimiting)"
              />
            </div>

            <div v-if="item.isChiral === 1" class="flex items-center space-x-2">
              <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Chiral Descriptor
              </label>
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-secondary text-secondary-foreground">
                {{ item.chiralDescriptor }}
              </span>
            </div>
          </div>

          <!-- Compound Image -->
          <div class="flex justify-center lg:justify-end">
            <div class="relative">
              <el-image
                :src="item.reagentImg"
                class="w-32 h-32 rounded-lg border border-border object-cover"
                :preview-src-list="[item.reagentImg]"
              />
            </div>
          </div>
        </div>

        <!-- Instructions -->
        <div class="mt-6 p-4 bg-muted/50 rounded-lg border border-border/50">
          <p class="text-sm font-medium text-foreground">
            <span v-if="item.reagentRole === 'solvent'">Purity - Volume - Density</span>
            <span v-else>
              Limiting - Equiv - Purity - Quantity - Density/Conc. Use
              <span class="text-destructive font-semibold">-1</span> to indicate excessive amount.
            </span>
          </p>
        </div>

        <!-- Input Fields -->
        <div class="mt-6">
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <div class="space-y-2">
              <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Equiv.
              </label>
              <el-input
                v-model="item.eq"
                size="small"
                :disabled="item.isLimiting === 1 || store.isReadonly"
                @input="calcMolesByEqLocal(item)"
                class="w-full form-input"
              />
            </div>

            <div class="space-y-2">
              <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Purity (%)
              </label>
              <el-tooltip content="Purity of the Reactant/Product/Reagent" placement="top">
                <el-input
                  v-model="item.purity"
                  size="small"
                  :disabled="store.isReadonly"
                  @input="calcQuantityByMoles(item)"
                  class="w-full"
                />
              </el-tooltip>
            </div>

            <div class="space-y-2">
              <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Quantity
              </label>
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
            </div>

            <div class="space-y-2">
              <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Moles
              </label>
              <el-input v-model="item.moles" size="small" :disabled="store.isReadonly" readonly class="w-full">
                <template #append>
                  <span class="inline-flex items-center px-2 py-1 text-xs font-medium bg-secondary text-secondary-foreground rounded">
                    {{ item.molesUnit }}
                  </span>
                </template>
              </el-input>
            </div>

            <div class="space-y-2">
              <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Density (g/mL)
              </label>
              <el-input
                v-model="item.density"
                size="small"
                :disabled="store.isReadonly"
                @input="calcVolume(item)"
                class="w-full"
              />
            </div>

            <div class="space-y-2">
              <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Conc. (mol/L)
              </label>
              <el-input
                v-model="item.concentration"
                size="small"
                :disabled="store.isReadonly"
                @input="calcVolume(item)"
                class="w-full"
              />
            </div>

            <div class="space-y-2">
              <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Volume
              </label>
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
            </div>

            <div class="space-y-2">
              <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                CAS#
              </label>
              <el-input v-model="item.cas" size="small" :disabled="store.isReadonly" class="w-full" />
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div v-if="!store.isReadonly" class="mt-6 flex justify-end space-x-2">
          <button
            @click="saveAdditionalInfoToDB(reagentFormRef[index], item)"
            :disabled="store.isReadonly"
            class="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-3 action-button"
          >
            <el-icon class="mr-2 h-4 w-4"><Save /></el-icon>
            Save
          </button>
          <button
            v-if="showDeleteOption(item.reagentRole)"
            @click="deleteRow(String(item.reagentId))"
            :disabled="store.isReadonly"
            class="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-destructive text-destructive-foreground hover:bg-destructive/90 h-9 px-3 action-button"
          >
            <el-icon class="mr-2 h-4 w-4"><Delete /></el-icon>
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { elnStore } from '#/store';

// Store
const store = elnStore();

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
