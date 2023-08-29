import { tagColors, triggerTypes } from '@/constants';
import { Popover, Tag } from 'antd';
import { cronTranstionCN, formatTime } from './format';

/**
 * 用于各重复组件方法的渲染
 */

/**
 * 渲染触发方式
 * @returns
 */
export const renderScheduleType = (baseInfo: any) => {
  const { schedule_type, schedule_rule } = baseInfo;
  const text = triggerTypes?.filter((item) => item?.value === schedule_type)[0]
    ?.label;
  if (schedule_type === 'cron') {
    return (
      <div>
        {text}
        <span>{`（${cronTranstionCN(schedule_rule)}）`}</span>
      </div>
    );
  }
  if (schedule_type === 'once') {
    return (
      <div>
        {text}
        <span>{`（${formatTime(schedule_rule)}）`}</span>
      </div>
    );
  }
  return <div>{text}</div>;
};

/**
 * 渲染标签
 * @param labels
 * @returns
 */
export const renderTags = (labels: any) => {
  const tagItems = labels?.map((item: any) => {
    const temp = tagColors?.filter((el) => el.type === item?.color)[0];
    return (
      <Tag key={item?.name} color={temp?.type}>
        {item?.name}
      </Tag>
    );
  });

  return <Popover title={<div>{tagItems}</div>}>{tagItems}</Popover>;
};
