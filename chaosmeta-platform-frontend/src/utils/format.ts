// 示例方法，没有实际意义
import moment from 'moment';
export function trim(str: string) {
  return str.trim();
}

// 转换为日期时间格式
export const formatTime = (time: moment.MomentInput) => {
  if (time) {
    return moment(time)?.format('YYYY-MM-DD HH:mm:ss');
  }
  return '--';
};
