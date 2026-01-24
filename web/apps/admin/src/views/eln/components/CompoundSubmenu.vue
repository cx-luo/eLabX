<template>
  <ElCard class="compound-submenu-card">
    <template #header>
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold">Compounds</h3>
        <div class="flex items-center gap-2">
          <el-badge :value="totalCompounds" :max="99" class="item">
            <span class="text-sm text-muted-foreground">Total: {{ totalCompounds }}</span>
          </el-badge>
        </div>
      </div>
    </template>

    <div class="space-y-4">
      <!-- Reactants Section -->
      <el-collapse v-model="activeMenus" class="compound-menu-collapse">
        <el-collapse-item name="reactants" :disabled="reactants.length === 0">
          <template #title>
            <div class="flex items-center gap-3 w-full">
              <div 
                class="w-4 h-4 rounded-full flex-shrink-0"
                :style="{ backgroundColor: getRoleColor('reactant') }"
              />
              <span class="font-medium">Reactants</span>
              <el-badge :value="reactants.length" :max="99" class="ml-auto" />
            </div>
          </template>
          <div class="pt-4 space-y-4">
            <ReagentCard 
              v-for="(item, index) in reactants" 
              :key="item.reagentId || `reactant-${index}`"
              :compound="item"
              :index="index"
            />
            <div v-if="reactants.length === 0" class="text-center py-8 text-muted-foreground text-sm">
              No reactants added yet
            </div>
          </div>
        </el-collapse-item>

        <!-- Products Section -->
        <el-collapse-item name="products" :disabled="products.length === 0">
          <template #title>
            <div class="flex items-center gap-3 w-full">
              <div 
                class="w-4 h-4 rounded-full flex-shrink-0"
                :style="{ backgroundColor: getRoleColor('product') }"
              />
              <span class="font-medium">Products</span>
              <el-badge :value="products.length" :max="99" class="ml-auto" />
            </div>
          </template>
          <div class="pt-4 space-y-4">
            <ReagentCard 
              v-for="(item, index) in products" 
              :key="item.reagentId || `product-${index}`"
              :compound="item"
              :index="index"
            />
            <div v-if="products.length === 0" class="text-center py-8 text-muted-foreground text-sm">
              No products added yet
            </div>
          </div>
        </el-collapse-item>

        <!-- Reagents Section -->
        <el-collapse-item name="reagents" :disabled="reagents.length === 0">
          <template #title>
            <div class="flex items-center gap-3 w-full">
              <div 
                class="w-4 h-4 rounded-full flex-shrink-0"
                :style="{ backgroundColor: getRoleColor('reagent') }"
              />
              <span class="font-medium">Reagents</span>
              <el-badge :value="reagents.length" :max="99" class="ml-auto" />
            </div>
          </template>
          <div class="pt-4 space-y-4">
            <ReagentCard 
              v-for="(item, index) in reagents" 
              :key="item.reagentId || `reagent-${index}`"
              :compound="item"
              :index="index"
            />
            <div v-if="reagents.length === 0" class="text-center py-8 text-muted-foreground text-sm">
              No reagents added yet
            </div>
          </div>
        </el-collapse-item>

        <!-- Catalysts Section -->
        <el-collapse-item name="catalysts" :disabled="catalysts.length === 0">
          <template #title>
            <div class="flex items-center gap-3 w-full">
              <div 
                class="w-4 h-4 rounded-full flex-shrink-0"
                :style="{ backgroundColor: getRoleColor('catalyst') }"
              />
              <span class="font-medium">Catalysts</span>
              <el-badge :value="catalysts.length" :max="99" class="ml-auto" />
            </div>
          </template>
          <div class="pt-4 space-y-4">
            <ReagentCard 
              v-for="(item, index) in catalysts" 
              :key="item.reagentId || `catalyst-${index}`"
              :compound="item"
              :index="index"
            />
            <div v-if="catalysts.length === 0" class="text-center py-8 text-muted-foreground text-sm">
              No catalysts added yet
            </div>
          </div>
        </el-collapse-item>

        <!-- Solvents Section -->
        <el-collapse-item name="solvents" :disabled="solvents.length === 0">
          <template #title>
            <div class="flex items-center gap-3 w-full">
              <div 
                class="w-4 h-4 rounded-full flex-shrink-0"
                :style="{ backgroundColor: getRoleColor('solvent') }"
              />
              <span class="font-medium">Solvents</span>
              <el-badge :value="solvents.length" :max="99" class="ml-auto" />
            </div>
          </template>
          <div class="pt-4 space-y-4">
            <ReagentCard 
              v-for="(item, index) in solvents" 
              :key="item.reagentId || `solvent-${index}`"
              :compound="item"
              :index="index"
            />
            <div v-if="solvents.length === 0" class="text-center py-8 text-muted-foreground text-sm">
              No solvents added yet
            </div>
          </div>
        </el-collapse-item>

        <!-- Other Compounds Section -->
        <el-collapse-item name="others" :disabled="others.length === 0">
          <template #title>
            <div class="flex items-center gap-3 w-full">
              <div 
                class="w-4 h-4 rounded-full flex-shrink-0"
                :style="{ backgroundColor: getRoleColor('other') }"
              />
              <span class="font-medium">Other Compounds</span>
              <el-badge :value="others.length" :max="99" class="ml-auto" />
            </div>
          </template>
          <div class="pt-4 space-y-4">
            <ReagentCard 
              v-for="(item, index) in others" 
              :key="item.reagentId || `other-${index}`"
              :compound="item"
              :index="index"
            />
            <div v-if="others.length === 0" class="text-center py-8 text-muted-foreground text-sm">
              No other compounds added yet
            </div>
          </div>
        </el-collapse-item>
      </el-collapse>

      <!-- Empty State -->
      <div v-if="totalCompounds === 0" class="text-center py-12">
        <div class="text-muted-foreground">
          <p class="text-lg font-medium mb-2">No compounds added</p>
          <p class="text-sm">Save a reaction to see compounds here</p>
        </div>
      </div>
    </div>
  </ElCard>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { ElCard, ElCollapse, ElCollapseItem, ElBadge } from 'element-plus';
import ReagentCard from './ReagentCard.vue';
import { elnStore } from '#/store';
import type { TableDataStruct } from '#/types';

const store = elnStore();

// Active menu items (expanded by default)
const activeMenus = ref<string[]>(['reactants', 'products', 'reagents']);

// Get compounds from store
const compounds = computed<TableDataStruct[]>(() => store.formData || store.tableData || []);

// Group compounds by role
const reactants = computed(() => 
  compounds.value.filter(item => item.reagentRole === 'reactant')
);

const products = computed(() => 
  compounds.value.filter(item => item.reagentRole === 'product')
);

const reagents = computed(() => 
  compounds.value.filter(item => item.reagentRole === 'reagent')
);

const catalysts = computed(() => 
  compounds.value.filter(item => item.reagentRole === 'catalyst')
);

const solvents = computed(() => 
  compounds.value.filter(item => item.reagentRole === 'solvent')
);

const others = computed(() => 
  compounds.value.filter(item => {
    const role = item.reagentRole?.toLowerCase();
    return role && !['reactant', 'product', 'reagent', 'catalyst', 'solvent'].includes(role);
  })
);

const totalCompounds = computed(() => compounds.value.length);

// Get role color
const getRoleColor = (role: string): string => {
  const colorMap: Record<string, string> = {
    reactant: 'rgba(7,203,75,0.91)',
    product: 'rgb(11,61,232)',
    solvent: 'rgba(190,175,9,0.98)',
    base: 'rgba(5,154,119,0.93)',
    catalyst: 'rgba(241,78,7,0.98)',
    reagent: 'rgba(128,7,241,0.8)',
    additive: 'rgba(128,7,241,0.8)',
    ligand: 'rgba(128,7,241,0.8)',
    other: 'rgb(2,51,44)',
  };
  return colorMap[role.toLowerCase()] || 'rgb(2,51,44)';
};
</script>

<style scoped lang="scss">
.compound-submenu-card {
  width: 100%;
}

.compound-menu-collapse {
  :deep(.el-collapse-item__header) {
    padding: 12px 16px;
    font-size: 14px;
    border-radius: 6px;
    transition: all 200ms ease-in-out;
    
    &:hover {
      background-color: rgba(0, 0, 0, 0.02);
    }
    
    &.is-disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }
  
  :deep(.el-collapse-item__content) {
    padding: 0 16px 16px;
  }
  
  :deep(.el-collapse-item__wrap) {
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
    margin-bottom: 8px;
    
    &:last-child {
      border-bottom: none;
      margin-bottom: 0;
    }
  }
}

// Badge styling
:deep(.el-badge__content) {
  border: none;
  font-weight: 500;
}
</style>

