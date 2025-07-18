import { defineStore } from 'pinia';
import { computed, reactive, ref } from 'vue';
import type { FormInstance } from 'element-plus';
import type { TableDataStruct, ReactionBasicInfo, ReactionSample } from '#/types';

// You need to import or define loadReagents for this to work.
// For now, let's assume it's imported from an API module:
import { loadReagents } from '#/api/enote'; // <-- Make sure this path is correct for your project

export const elnStore = defineStore('eln', () => {
  const reactionId = ref(0);
  const formData = ref<TableDataStruct[]>();
  const formRules = ref();
  const samplesTable = ref<ReactionSample[]>([]);
  const reactionSmiles = ref('');
  const cdStructure = ref('');
  const otherReportTable = ref([]);
  const projectBasicInfo = ref<ReactionBasicInfo>();

  const nmrTable = ref([]);
  const procedureText = ref('');
  const procedureHtml = ref('');
  const isGetTargetProduct = ref(false);
  const procedureComments = ref('');

  const procedureTableData = ref([]);
  const tableData = ref<TableDataStruct[]>([]);

  const tableDataMap = new Map();
  const samplesTableMap = new Map();
  const lcmsTableMap = new Map();
  const nmrTableMap = new Map();
  const otherReportTableMap = new Map();
  const workbookMap = new Map();
  const projectBasicInfoMap = new Map();
  const procedureTextMap = new Map();
  const procedureHtmlMap = new Map();
  const isReadOnlyMap = new Map();
  const rxnImgMap = new Map();
  const reactionConditionsDataMap = new Map();

  const lcmsTable = ref([]);

  const formRefData = ref<FormInstance[]>([]);
  const reactantFormValidate = ref(false);
  const productFormValidate = ref(false);
  const solventFormValidate = ref(false);
  const otherFormValidate = ref(false);
  const reagentFormValidate = ref(false);
  const samplesFormValidate = ref(false);
  const basicInfoFormValidate = ref(false);
  const reactionConditionFormValidate = ref(false);
  const seenProductIds = new Set();

  const projectSearchInfo = reactive({
    total: 0,
    projectInfo: [],
    conditionForm: {
      userId: null,
      projectName: '',
      pageName: '',
    },
  });

  const productList = computed(() => tableData.value.filter((item) => item.reagentRole === 'product'));

  async function loadReagent(rxnId: number) {
    try {
      const response = await loadReagents(rxnId);
      // If your API returns a statusCode, check it. Otherwise, just check for data.
      if ('statusCode' in response && response.statusCode === 404) {
        tableData.value = [];
        return;
      }
      if (response && response.data && Array.isArray(response.data.reagents)) {
        tableData.value = response.data.reagents;
        tableDataMap.set(rxnId, response.data.reagents);
      } else {
        tableData.value = [];
      }
    } catch (error) {
      // Handle error, e.g., network error
      tableData.value = [];
    }
  }

  const isReadonly = computed(() => {
    if (projectBasicInfo.value) {
      return projectBasicInfo.value.rxnStatus !== 'open';
    } else {
      return false;
    }
  });

  const colorList = [
    'beige',
    'black',
    'Brick red',
    'brown',
    'Brown yellow',
    'colorless',
    'dark blue',
    'dark green',
    'dark grey',
    'dark yellow',
    'green',
    'grey',
    'Grey white',
    'light blue',
    'light brown',
    'light green',
    'light grey',
    'light orange',
    'light pink',
    'light purple',
    'Light reddish brown',
    'light yellow',
    'off-white',
    'orange',
    'oyster white',
    'pink',
    'purple',
    'red',
    'reddish brown',
    'Reddish brown yellow',
    'tan',
    'white',
    'yellow',
    'yellow green',
  ];
  const statusList = [
    'amorphous solid',
    'crystalline',
    'film',
    'liquid',
    'oil',
    'powder',
    'semi-solid',
    'solid',
    'viscous oil',
  ];

  const rxnList = [
    // "Enzyme-catalyzed aminolysis of ester",
    // "Enzyme-catalyzed hydrolysis of ester",
    // "Enzyme-catalyzed nitrile hydrolysis",
    // "Enzyme-catalyzed nitroreduction",
    // "Enzyme-catalyzed oxidation of alcohol",
    // "Enzyme-catalyzed reduction of imine",
    // "Enzyme-catalyzed reduction of ketone",
    // "Enzyme-catalyzed transamination",
    // "Enzyme-catalyzed transesterification",
    'Epoxidation',
    'Fluorination',
    'Heck coupling',
    'Heterogenous hydrogenation',
    'Hiyama coupling',
    'Hydroamination',
    'Kumada coupling',
    'Liebeskind-Srogl coupling',
    'Mitsunobu reaction',
    'Miyaura boration',
    'Negishi coupling',
    'Ni-catalyzed reductive coupling',
    'Photoredox coupling',
    'Salt formation with chiral acid/base',
    'Sonogashira coupling',
    'Stannylation',
    'Stille coupling',
    'Suzuki-Miyaura coupling (sp2-sp2)',
    'Suzuki-Miyaura coupling (sp2-sp3)',
    'Trifluoromethylation',
    'Tsuji-Trost reaction',
  ];

  const qualifierList = [
    'Crude Product',
    'Experiment Discarded',
    'Product NOt isolated',
    'Purified Product',
    'Yield form Conversion Rate',
    'Yield from HPLC',
  ];

  return {
    reactionId,
    reactionSmiles,
    cdStructure,
    procedureTableData,
    formRefData,
    reactantFormValidate,
    productFormValidate,
    solventFormValidate,
    otherFormValidate,
    reagentFormValidate,
    samplesFormValidate,
    basicInfoFormValidate,
    reactionConditionFormValidate,
    // 用于初始化空列表
    tableData,
    rxnImgMap,
    otherReportTable,
    productList,
    statusList,
    qualifierList,
    rxnList,
    // 用于尚未加载的数据
    // user: null as UserInfo | null,
    lcmsTable,
    nmrTable,
    samplesTable,
    isReadonly,
    reactionConditionsDataMap,
    colorList,
    loadReagent,
    otherReportTableMap,
    procedureHtml,
    procedureHtmlMap,
    procedureText,
    isGetTargetProduct,
    procedureComments,
    tableDataMap,
    samplesTableMap,
    lcmsTableMap,
    nmrTableMap,
    projectBasicInfo,
    projectBasicInfoMap,
    procedureTextMap,
    workbookMap,
    isReadOnlyMap,
    seenProductIds,
    projectSearchInfo,
    // loadWorkbookAndReagents
    formData,
    formRules,
  };
});
