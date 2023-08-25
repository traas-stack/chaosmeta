export const DEFAULT_NAME = 'Umi Max';

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
  pressure: '#FFD5D5',
  other: '#D4E3F1',
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
