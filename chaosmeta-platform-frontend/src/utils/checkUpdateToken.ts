import { updateToken } from '@/services/chaosmeta/UserController';

/**
 * 检查当前token是否过期
 */
export const checkUpdateToken = async () => {
  // 过期时间5分钟, 提前10秒刷新token，防止token过期导致
  const time = 5 * 60 * 60 - 10000;
  //  token存入时的时间戳
  const createTime = localStorage.getItem('tokenCreateTime');
  // 当前时间戳
  if (createTime) {
    const curTime = Date.now();
    const diffTime = curTime - Number(createTime);
    // 当前时间与存入时间差大于5分钟时，更新token
    if (diffTime >= time) {
      const result = await updateToken();
      console.log(result, 'result===');
    }
  }
};
