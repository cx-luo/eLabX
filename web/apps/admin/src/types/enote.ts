export interface TableDataStruct {
  reagentId: number;
  reagentName: string;
  reagentSmiles: string;
  mw: number;
  reagentRole: string;
  formula: string;
  cas: string;
  eq: number;
  purity: number;
  quantity: number;
  quantityUnit: string;
  concentration: number;
  density: number;
  densityUnit: string;
  compoundId: number;
  yield: number;
  isLimiting: number;
  isChiral: number;
  stereoCentersCnt: number;
  chiralDescriptor: string;
  reagentImg: string;
  productAlias: string;
  moles: number;
  molesUnit: string;
  volume: number;
  volumeUnit: string;
  cdStructure?: string;
}

export interface ReactionBasicInfo {
  projectName: string;
  startDate: string;
  creationDate: string;

  batch: null | number;
  authorId: null | number;
  witnessId: null | number;

  stepId: string;
  authorName: string;
  witnessName: string;

  rxnType: string;
  reference: string;
  pageName: string;

  doi: string;
  rxnStatus: string;
  comment: string;
  rxnTypeFromAI: string;
}

export interface ReactionSample {
  reactionId: number; // 反应ID
  reagentId: number; // 试剂ID
  sampleStatus: string; // 样品状态
  isTestSamples: number; // 是否为测试样品
  mf: string; // 分子式 (Molecular Formula)
  mass: number; // 质量
  purity: number; // 纯度百分比
  synthesised: number; // 合成量
  synthesisedUnit: string; // 合成单位
  sampleYield: number; // 样品产率
  moles: number; // 摩尔数
  molesUnit: string; // 摩尔单位
  amountSubmitted: number; // 提交的数量
  amountSubmittedUnit: string; // 提交数量的单位
  productAlias: string; // 产品别名
  color: string; // 颜色
  qualifier: string; // 定性描述
}

const checksumCas = (part1: any, part2: any) => {
  // 将part1和part2拼接起来并转换为数组，然后反转数组
  const items = Array.from(part1 + part2).reverse();
  let sum = 0;

  // 由于数组是从尾部开始的，因此我们从1开始计数，而不是0
  items.forEach((digit, index) => {
    if (typeof digit === 'string') {
      sum += parseInt(digit, 10) * (index + 1);
    }
  });
  // 返回计算出的校验码
  return sum % 10;
};

export const validateCAS = (value: any, callback: any) => {
  // 正则表达式用于验证CAS号的格式
  // 2到6位数字，接着是2位数字，最后是1位数字，中间用连字符隔开
  const casRegex = /^(\d{2,7})-(\d{2})-(\d)$/;
  const match = value.match(casRegex);

  // 如果格式不匹配，返回false
  if (!match) {
    callback(new Error('The format is incorrect'));
  } else {
    // 提取CAS号的前两部分
    const [, part1, part2, part3] = match;

    // 检查计算出的校验码与提供的CAS号的第三部分是否一致
    if (checksumCas(part1, part2) === parseInt(part3, 10)) {
      callback();
    } else {
      callback(new Error('CAS error, please re-enter'));
    }
  }
};
