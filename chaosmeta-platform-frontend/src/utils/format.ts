import { getLocale } from '@umijs/max';
import cronstrue from 'cronstrue';
import 'cronstrue/locales/zh_CN';
import dayjs from 'dayjs';
import { v1 } from 'uuid';

export function trim(str: string) {
  return str.trim();
}

export const formatTime = (time?: dayjs.Dayjs | string | number) => {
  if (time) {
    return dayjs(time)?.format('YYYY-MM-DD HH:mm:ss');
  }
  return '--';
};

/**
 * 将传入的参数转换为两位数格式，不足两位时在前面补零。
 * @param {any} val - 需要转换的参数。
 * @returns {string|number} 返回转换后的字符串或数字。
 */
export const montageTime = (val: any) => {
  // 如果参数长度为一，则在前面补零。
  if (val?.toString()?.length === 1) {
    return `0${val}`;
  }
  // 否则直接返回原值。
  return val;
};

/**
 * 传入秒数，返回时分秒的拼接值
 * @param second
 * @returns
 */
export const handleTimeTransform = (second: number) => {
  // 计算剩余的秒数、分数和小时数
  const renderSecond = second % 60; // 剩余秒数
  const minute = Math.floor(second / 60); // 分钟数
  const residueMinute = minute % 60; // 剩余分钟数
  const hour = Math.floor(second / 60 / 60); // 小时数
  // 格式化时间字符串
  const text = `${montageTime(hour)}:${montageTime(
    residueMinute,
  )}:${montageTime(renderSecond)}`;
  return text;
};

/**
 * 判断持续时长的类型是分或秒，分就转换为对应秒数返回，秒返回对应值
 * @param value
 * @returns
 */
export const formatDuration = (value: string): number => {
  if (value) {
    if (value?.includes('m')) {
      const duration = Number(value?.split('m')[0]);
      return duration * 60;
    }
    return Number(value?.split('s')[0]);
  }
  return 0;
};

/**
 * 传入list，将其转换为编排可用的二维数组
 * 将原始数据进行转换和排序，返回目标数组
 * @param list
 * @param isDetail // 详情用于展示时不需要拼接空数据
 * @returns
 */
export const arrangeDataOriginTranstion = (list: any[], isDetail?: boolean) => {
  if (list?.length > 0) {
    // 如果列表不为空
    const originList = JSON.parse(JSON.stringify(list)); // 克隆一份原始列表
    originList?.sort(
      // 根据行号和列号排序
      (
        a: { row: number; column: number },
        b: { row: number; column: number },
      ) => {
        if (a.row === b.row) {
          return a.column - b.column;
        }
        return a.row - b.row;
      },
    );
    // row和column的起始值都是0，这里统一加1，为了后面id绑定拖拽组件时不出问题（id不能为0）
    // 这里只是作为展示编排组件用，提交时会再处理一遍
    originList?.forEach((item: { row: number; column: number }) => {
      item.row = item.row + 1;
      item.column = item.column + 1;
    });
    const targetList = originList?.reduce((groups: any, item: any) => {
      const rowIndex = groups?.findIndex(
        // 查找当前项所在的组
        (el: { row: number }) => el?.row === item.row,
      );
      if (rowIndex !== -1 && groups[rowIndex]?.children) {
        groups[rowIndex].children.push(item); // 将该项添加到该组中
      } else {
        // 否则新建一个组
        groups[groups?.length] = {
          // row初始可能为0，但这里id不能为0，默认加一
          row: item.row,
          id: item.row,
          // 手动添加一个id，拖拽组件库绑定数据时需要 todo -- 后期优化看是否可以去除
          children: [{ ...item }],
        };
      }
      return groups;
    }, []);

    // 详情不需要首位拼接空数据占位
    if (isDetail) {
      return targetList;
    }
    // 在首尾各插入一行，用于占位，targetList?.length + 2，上面逻辑为了避免id为0需要加1，这里就加2
    const popItem = targetList[targetList?.length - 1];
    targetList?.unshift({
      row: popItem?.row + 1,
      id: popItem?.row + 1,
      children: [],
    }); // 在第一个位置插入一个空白项
    targetList?.push({
      row: popItem?.row + 2,
      id: popItem?.row + 2,
      children: [],
    }); // 在最后插入一个空白项
    return targetList; // 返回目标数组
  }
  return [
    // 如果列表为空，则返回一个包含三个空白项的数组
    { row: 1, id: 1, children: [] },
    { row: 2, id: 2, children: [] },
    { row: 3, id: 3, children: [] },
  ];
};

/**
 * 将编排数据转换为后端需要的坐标类型数据
 * @param result
 */
export const arrangeDataResultTranstion = (result: any) => {
  // 过滤出有子项的项，并转换为新数组
  const newResult = result?.filter(
    (item: { children: any[] }) => item?.children?.length > 0,
  );
  // 新建一个空数组作为存储结果的容器
  const newList: any[] = [];
  // 遍历原始数组中的每一项
  newResult.forEach((item: { children: any[] }, index: number) => {
    // 遍历该项的每一项
    item.children?.forEach((el, j) => {
      const newArgs: any[] = [];
      // 处理args参数（动态表单部分）
      if (el?.args_value) {
        Object.keys(el?.args_value)?.forEach((key) => {
          if (el?.args_value[key] !== undefined) {
            newArgs?.push({
              args_id: Number(key),
              value: el?.args_value[key]?.toString(),
            });
          }
        });
      }
      newList.push({ ...el, row: index, column: j, args_value: newArgs });
    });
  });
  return newList;
};

/**
 * 将cron表达式转换为中文描述
 * @param cronRule
 * @returns
 */
export const cronTranstionCN = (cronRule: string) => {
  // 获取当前语言，用于cron表达式描述展示
  let curIntl = getLocale();
  curIntl = curIntl.replace('-', '_');
  let cronCN = '';
  if (cronRule) {
    try {
      cronCN = cronstrue.toString(cronRule, { locale: curIntl });
    } catch (error) {
      cronCN = 'error';
    }
  }
  return cronCN;
};

/**
 * 获取国际化文本
 * @param temp 包含中文文案和对应英文文案的对象
 * @returns 返回当前环境下的文案
 */
export const getIntlLabel = (temp: { label: string; labelUS: string }) => {
  // 获取当前环境
  const curIntl = getLocale();
  // 如果当前环境为英文环境
  if (curIntl === 'en-US') {
    // 返回对应的英文文案
    return temp?.labelUS;
  }
  // 返回对应的中文文案
  return temp?.label;
};

/**
 * 获取国际化文本
 * @param temp 包含中文文案和对应英文文案的对象
 * @returns 返回当前环境下的文案
 */
export const getIntlName = (temp: { name: string; nameCn: string }) => {
  // 获取当前环境
  const curIntl = getLocale();
  if (temp?.name && temp?.nameCn) {
    // 如果当前环境为英文环境
    if (curIntl === 'en-US') {
      // 返回对应的英文文案
      return temp?.name;
    }
    // 返回对应的中文文案
    return temp?.nameCn;
  }
  return temp?.name;
};

/**
 * 根据传入的对象和key返回对应的中文/英文文案
 * @param temp 包含中文文案和对应英文文案的对象
 * @param cnKey 中文key
 * @param usKey 英文key
 * @returns
 */
export const getIntlText = (temp: any, cnKey: string, usKey: string) => {
  // 获取当前环境
  const curIntl = getLocale();
  // 如果当前环境为英文环境
  if (curIntl === 'en-US') {
    // 返回对应的英文文案
    return temp?.[usKey];
  }
  // 返回对应的中文文案
  return temp?.[cnKey];
};

/**
 * 用于生成uuid
 * v1基于时间戳生成
 */
export const generateUuid = () => {
  // 将生成的uuid中的连接线去除
  return v1()?.replaceAll('-', '');
};

/**
 * 复制实验时将数据中不需要的字段去除
 * @param data
 * @returns
 */
export const copyExperimentFormatData = (data: any) => {
  const result = JSON.parse(JSON.stringify(data));
  const newLabels = result?.labels?.map((item: { id: number }) => item?.id);
  result?.workflow_nodes?.forEach(
    (item: {
      uuid: string;
      experiment_uuid: undefined;
      create_time: undefined;
      update_time: undefined;
      exec_range: any;
      args_value: { args_id: any; value: any }[];
    }) => {
      item.uuid = generateUuid();
      item.experiment_uuid = undefined;
      item.create_time = undefined;
      item.update_time = undefined;
      item.exec_range = {
        // 使用的是复制来的数据，需要去除创建时不需要的字段
        ...item.exec_range,
        create_time: undefined,
        id: undefined,
        workflow_node_instance_uuid: undefined,
        update_time: undefined,
      };
      if (item.args_value) {
        item.args_value = item?.args_value?.map(
          (el: { args_id: any; value: any }) => {
            return {
              args_id: el?.args_id,
              value: el?.value,
            };
          },
        );
      }
    },
  );
  const params = {
    ...result,
    name: `${result?.name}副本`,
    uuid: undefined,
    labels: newLabels,
  };
  return params;
};

/**
 * 格式化表单项名称
 * @param item 表单项对象，包含id、key和execType三个属性
 * @param parentName 父级表单项名称，可选参数
 * @returns 返回格式化后的表单项名称数组
 */
export const formatFormName = (
  item: { id: number; key: string; execType: string },
  parentName?: string,
) => {
  const { id, key, execType } = item;
  let name: any = '';
  if (parentName) {
    // 拼接父name
    name = [parentName, id.toString()];
  }
  if (execType === 'flow_common') {
    // 如果执行类型为'flow_common'，则和key进行拼接
    name = ['flow_range', key];
  }
  if (execType === 'measure_common') {
    // 如果执行类型为'measure_common'，则和key进行拼接
    name = ['measure_range', key];
  }
  if (key === 'duration') {
    // 如果键名为'duration'，则直接返回该键名
    name = key;
  }
  return name;
};

/**
 * 时间戳字符串转时间戳
 * @param time 时间字符串，"1706150902050"转为1706150902050，"2024-01-12 20:02:16"则不变
 * @returns
 */
export const timesStampString = (time: string) => {
  return !time || isNaN(Number(time)) ? time : Number(time);
};
