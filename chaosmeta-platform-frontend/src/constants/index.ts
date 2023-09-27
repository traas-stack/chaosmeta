export const DEFAULT_NAME = 'Umi Max';

// KubernetesController使用
// -1: 开发环境，0：生产环境
export const envType = 0;

export const tagColors = [
  {
    color: '#EDEEEF',
    type: 'default',
    borderColor: '#DADADA',
  },
  {
    color: '#FFD7D7',
    type: 'red',
    borderColor: '#F8B4B4',
  },
  {
    color: '#FFF2B5',
    type: 'yellow',
    borderColor: '#FFE361',
  },
  {
    color: '#CDCCFF',
    type: 'purple',
    borderColor: '#B2B1FF',
  },
  {
    color: '#FFE0CB',
    type: 'orange',
    borderColor: '#FFCDAA',
  },
  {
    color: '#DAFFA7',
    type: 'green',
    borderColor: '#C5FF71',
  },
];

// 编排节点类型对应颜色
export const arrangeNodeTypeColors: any = {
  fault: '#F5E2CC',
  measure: '#C6F8E0',
  flow: '#FFD5D5',
  other: '#D4E3F1',
};

export const nodeMode = {
  fault: '故障节点',
};

// 计算会有小数问题，直接在这里列举处理了
// secondStep 每段时间轴的间隔时间
// 宽度对应时间，默认1s为3px
// widthSecond 1s对应的宽度
export const scaleStepMap: any = {
  33: {
    secondStep: 90,
    widthSecond: 1,
  },
  66: {
    secondStep: 45,
    widthSecond: 2,
  },
  100: {
    secondStep: 30,
    widthSecond: 3,
  },
  150: {
    secondStep: 20,
    widthSecond: 4.5,
  },
  200: {
    secondStep: 15,
    widthSecond: 6,
  },
  300: {
    secondStep: 10,
    widthSecond: 9,
  },
};

// 触发方式选项
export const triggerTypes = [
  { label: '手动触发', value: 'manual', labelUS: 'manual trigger' },
  { label: '单次定时', value: 'once', labelUS: 'single timing' },
  { label: '周期性', value: 'cron', labelUS: 'periodicity' },
];

// 实验结果状态
export const experimentResultStatus = [
  {
    value: 'Pending',
    label: '等待中',
    labelUS: 'pending',
    color: 'blue',
    type: 'info',
  },
  {
    value: 'Running',
    labelUS: 'running',
    label: '运行中',
    color: 'blue',
    type: 'info',
  },
  {
    value: 'Succeeded',
    labelUS: 'succeeded',
    label: '成功',
    color: 'green',
    type: 'success',
  },

  {
    value: 'Failed',
    labelUS: 'failed',
    label: '失败',
    color: 'red',
    type: 'error',
  },
  {
    value: 'error',
    labelUS: 'error',
    label: '错误',
    color: 'red',
    type: 'error',
  },
];

// 实验状态
export const experimentStatus = [
  {
    value: 0,
    label: '待执行',
    labelUS: 'to be executed',
    color: 'blue',
  },
  {
    value: 1,
    label: '执行成功',
    labelUS: 'execution succeed',
    color: 'green',
  },
  {
    value: 2,
    label: '执行失败',
    labelUS: 'execution failed',
    color: 'red',
  },
  {
    value: 3,
    label: '执行中',
    labelUS: 'executing',
    color: 'blue',
  },
];

// 节点类型
export const nodeTypeMap: any = {
  fault: '故障节点',
  measure: '度量引擎',
  flow: '流量注入',
  wait: '等待时长',
};

// 节点类型
export const nodeTypeMapUS: any = {
  fault: 'faulty node',
  measure: 'measurement engine',
  flow: 'flow injection',
  wait: 'waiting time',
};


export const nodeTypes = [
  {
    label: '故障节点',
    labelUS: 'faulty node',
    type: 'fault',
  },
  {
    label: '度量引擎',
    labelUS: 'measurement engine',
    type: 'measure',
  },
  {
    label: '流量注入',
    labelUS: 'flow injection',
    type: 'flow',
  },
  {
    label: '其他节点',
    labelUS: 'other nodes',
    type: 'other',
  },
];
