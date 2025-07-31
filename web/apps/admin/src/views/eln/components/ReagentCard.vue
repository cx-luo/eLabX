<template>
  <div v-for="(item, index) in store.formData">
    <el-card style="margin-top: 5px" :style="{ border: '6px solid ' + getRoleColor(item.reagentRole) }">
      <el-form ref="reagentFormRef" :model="item" label-width="auto" style="width: 100%" :rules="store.formRules">
        <el-row>
          <el-col :md="16" :sm="16" :xs="16">
            <el-row>
              <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
                <el-form-item class="label-right-align" prop="reagentRole" label="Role">
                  <el-tag
                    effect="dark"
                    size="small"
                    :style="{
                      fontWeight: 'bolder',
                      backgroundColor: getRoleColor(item.reagentRole),
                    }"
                  >
                    {{ item.reagentRole }}
                  </el-tag>
                  <el-popover v-if="item.reagentRole != 'product'" placement="right" :width="100" trigger="click">
                    <template #reference>
                      <el-button v-if="!store.isReadonly" link style="color: blue" size="small">
                        <el-icon :size="20" style="vertical-align: middle">
                          <Edit />
                        </el-icon>
                        <!--                      <span style="vertical-align: middle"> Change </span>-->
                      </el-button>
                    </template>
                    <el-select v-model="item.reagentRole" size="small" :disabled="store.isReadonly">
                      <el-option
                        v-for="role in compoundRole"
                        :key="role.value"
                        :label="role.value"
                        :value="role.value"
                      />
                    </el-select>
                  </el-popover>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
                <el-form-item class="label-right-align" fixed="left" prop="id" label="№" type="index">
                  {{ item.reagentRole }}
                  {{ index + 1 }}
                </el-form-item>
              </el-col>
            </el-row>

            <el-row>
              <el-form-item class="label-right-align" label="CmpdName" prop="reagentName" show-overflow-tooltip>
                {{ item.reagentName }}
              </el-form-item>
            </el-row>

            <el-row>
              <el-form-item class="label-right-align" label="Smiles" prop="reagentSmiles" show-overflow-tooltip>
                {{ item.reagentSmiles }}
              </el-form-item>
            </el-row>

            <el-row>
              <el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="8">
                <el-form-item class="label-right-align" label="Formula" prop="formula" show-overflow-tooltip>
                  {{ item.formula }}
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="12" :md="12" :lg="6" :xl="4">
                <el-form-item class="label-right-align" prop="mw" label="Mol.Wt">{{ item.mw }} </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="12" :md="12" :lg="6" :xl="4">
                <el-form-item class="label-right-align" label="isChiral" prop="isChiral" show-overflow-tooltip>
                  {{ item.isChiral == 1 ? item.isChiral : 0 }}
                </el-form-item>
              </el-col>

              <el-col v-if="item.isChiral == 1" :xs="24" :sm="12" :md="12" :lg="6" :xl="4">
                <el-form-item
                  class="label-right-align"
                  label="stereoCentersCnt"
                  prop="stereoCentersCnt"
                  show-overflow-tooltip
                  >{{ item.stereoCentersCnt }}
                </el-form-item>
              </el-col>

              <el-col v-if="item.reagentRole == 'reactant'" :md="6" :sm="6" :xs="24">
                <el-form-item class="label-right-align" label="Limiting">
                  <el-switch
                    v-model="item.isLimiting"
                    :active-value="1"
                    :inactive-value="0"
                    style="--el-switch-on-color: #9507dc"
                    :disabled="store.isReadonly"
                    @change="handleSwitchChange(index, item.isLimiting)"
                  />
                </el-form-item>
              </el-col>

              <el-col v-if="item.isChiral == 1" :md="18" :sm="18" :xs="24">
                <el-form-item
                  class="label-right-align"
                  label="Chiral Descriptor"
                  prop="chiralDescriptor"
                  :required="item.isChiral == 1"
                >
                  <el-tag>{{ item.chiralDescriptor }}</el-tag>
                </el-form-item>
              </el-col>
            </el-row>

            <el-row>
              <el-form-item
                class="label-right-align"
                label="Fill in the order"
                prop="reagentName"
                show-overflow-tooltip
              >
                <el-text class="notice-text">
                  <span v-if="item.reagentRole === 'solvent'"> Purity - Volume - density </span>
                  <span v-else>
                    Limiting - Equiv - Purity - Quantity - Density/Conc. Use
                    <span style="color: red"> -1 </span> to indicate an excessive amount.
                  </span>
                </el-text>
              </el-form-item>
            </el-row>
          </el-col>
          <el-col :md="8" :sm="8" :xs="8" class="myImage">
            <el-form-item class="label-right-align" prop="reagentImg">
              <el-image
                :src="item.reagentImg"
                shape="square"
                fit="scale-down"
                size="small"
                :preview-src-list="[item.reagentImg]"
                style="height: 128px"
              />

              <!--              <div v-html="item.reagentImg"></div>-->
            </el-form-item>
          </el-col>
        </el-row>

        <!--details-->
        <el-row>
          <!--    eq -->
          <el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="4">
            <el-form-item class="label-right-align" label="Equiv." prop="eq">
              <el-input
                v-model="item.eq"
                size="small"
                class="workbook-input"
                @input="calcMolesByEqLocal(item)"
                :disabled="item.isLimiting == 1 || store.isReadonly"
              />
            </el-form-item>
          </el-col>

          <!--purity-->
          <el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="4">
            <el-form-item class="label-right-align" label="Purity(%)" prop="purity">
              <el-tooltip
                class="my-tooltip"
                effect="dark"
                content="Purity of the Reactant/Product/Reagent"
                placement="top-start"
              >
                <el-input
                  v-model="item.purity"
                  size="small"
                  class="workbook-input"
                  :disabled="store.isReadonly"
                  @input="calcQuantityByMoles(item)"
                />
              </el-tooltip>
            </el-form-item>
          </el-col>

          <!--quantity-->
          <el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="4">
            <el-form-item class="label-right-align" label="Quantity" prop="quantity">
              <el-tooltip
                class="my-tooltip"
                effect="dark"
                content="Mass of the Reactant/Product/Reagent. 实际参加反应的质量。"
                placement="top-start"
              >
                <el-input
                  v-model="item.quantity"
                  size="small"
                  class="workbook-input-no-padding"
                  :disabled="store.isReadonly"
                  @input="calcMoles(item)"
                >
                  <template #append>
                    <el-form-item class="label-right-align" prop="quantityUnit">
                      <el-select
                        v-model="item.quantityUnit"
                        size="small"
                        :disabled="store.isReadonly"
                        @change="calcMoles(item)"
                      >
                        <el-option v-for="item in unitList" :key="item.value" :label="item.value" :value="item.value" />
                      </el-select>
                    </el-form-item>
                  </template>
                </el-input>
              </el-tooltip>
            </el-form-item>
          </el-col>

          <!--moles-->
          <el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="4">
            <el-form-item class="label-right-align" label="Moles" prop="moles">
              <el-input
                v-model="item.moles"
                size="small"
                class="workbook-input-no-padding"
                :disabled="store.isReadonly"
                readonly
              >
                <template #append style="width: 20px">
                  <el-tag class="molesTag" size="small">{{ item.molesUnit }}</el-tag>
                </template>
              </el-input>
            </el-form-item>
          </el-col>

          <!--density-->
          <el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="4">
            <el-form-item class="label-right-align" prop="density" label="Density(g/mL)">
              <template #default="scope">
                <el-input
                  v-model="item.density"
                  size="small"
                  class="workbook-input"
                  :disabled="store.isReadonly"
                  @input="calcVolume(item)"
                >
                </el-input>
              </template>
            </el-form-item>
          </el-col>

          <!--concentration-->
          <el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="4">
            <el-form-item class="label-right-align" label="Conc.(mol/L)" prop="concentration">
              <el-input
                v-model="item.concentration"
                size="small"
                class="workbook-input-no-padding"
                :disabled="store.isReadonly"
                @input="calcVolume(item)"
              >
              </el-input>
            </el-form-item>
          </el-col>

          <!--volume-->
          <el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="4">
            <el-form-item class="label-right-align" prop="volume" label="Volume">
              <el-input
                v-model="item.volume"
                size="small"
                class="workbook-input-no-padding"
                @input="calcMolesByVolume(item)"
                :disabled="store.isReadonly"
              >
                <template #append>
                  <el-form-item class="label-right-align" prop="volumeUnit">
                    <el-select
                      v-model="item.volumeUnit"
                      size="small"
                      :disabled="store.isReadonly"
                      @change="calcMolesByVolume(item)"
                    >
                      <el-option
                        v-for="item in volumeUnitList"
                        :key="item.value"
                        :label="item.value"
                        :value="item.value"
                      />
                    </el-select>
                  </el-form-item>
                </template>
              </el-input>
            </el-form-item>
          </el-col>

          <!--cas-->
          <el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="4">
            <el-form-item class="label-right-align" label="CAS#">
              <el-input v-model="item.cas" size="small" class="workbook-input" :disabled="store.isReadonly"> </el-input>
            </el-form-item>
          </el-col>

          <!--button-->
          <el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="4">
            <el-form-item class="label-right-align" v-if="!store.isReadonly" fixed="right" label="operation">
              <template #default="scope">
                <el-button-group>
                  <el-button
                    type="primary"
                    @click="saveAdditionalInfoToDB(reagentFormRef[index], item)"
                    size="small"
                    style="background-color: #13b29e"
                    :disabled="store.isReadonly"
                  >
                    <el-icon>
                      <Save />
                    </el-icon>
                  </el-button>
                  <el-button
                    type="primary"
                    size="small"
                    @click.prevent="deleteRow(item.reagentId)"
                    style="margin-left: 2px; background-color: rgba(236, 21, 21, 0.91)"
                    v-if="showDeleteOption(item.reagentRole)"
                    :disabled="store.isReadonly"
                  >
                    <el-icon>
                      <Delete />
                    </el-icon>
                  </el-button>
                </el-button-group>
              </template>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { elnStore } from '#/store';
const store = elnStore();
const compoundRole = [
  {
    value: 'catalyst',
  },
  {
    value: 'reagent',
  },
  {
    value: 'solvent',
  },
  {
    value: 'base',
  },
  {
    value: 'ligand',
  },
  {
    value: 'additive',
  },
  {
    value: 'other',
  },
];

const getRoleColor = (role: string) => {
  switch (role) {
    case 'reactant':
      return 'rgba(7,203,75,0.91)';
    case 'product':
      return 'rgb(11,61,232)';
    case 'solvent':
      return 'rgba(190,175,9,0.98)';
    case 'base':
      return 'rgba(5,154,119,0.93)';
    case 'catalyst':
      return 'rgba(241,78,7,0.98)';
    case 'additives':
      return 'rgba(128,7,241,0.8)';
    // ... 其他角色颜色
    default:
      return 'rgb(2,51,44)'; // 默认颜色
  }
};

const activeRow = ref(-1);

const unitList = [
  {
    value: 'mg',
  },
  {
    value: 'g',
  },
  {
    value: 'kg',
  },
];

const volumeUnitList = [
  // {
  //   value: 'uL',
  // },
  {
    value: 'mL',
  },
  {
    value: 'L',
  },
];

const showDeleteOption = (role) => {
  return !(role == 'reactant' || role == 'product');
};
</script>
<style scoped lang="scss">
.label-right-align :deep(.el-form-item__label) {
  text-align: right;
  font-weight: bolder;
  font-size: 12px;
  color: black;
  font-family: 'Microsoft YaHei', Arial, Helvetica, sans-serif;
}

.label-right-align :deep(.el-form-item__content) {
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.my-tooltip .el-tooltip__content {
  color: #409eff; /* 设置字体颜色 */
}

.notice-text {
  color: #0c5bec;
  font-weight: bold;
}
</style>
