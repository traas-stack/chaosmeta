/**
 * 封装一层request，用于处理刷新token
 */
import {
  AxiosResponse,
  getLocale,
  history,
  request as requestUmi,
} from '@umijs/max';
import { message } from 'antd';
import cookie from './cookie';

// 由于有可能同时存在多个接口为401，会有调取多次刷新token接口的情况，需定义一个存储当前刷新token接口的变量
let reqFlag: any = null;
const loginTextMap: any = {
  'zh-CN': {
    notLogin: '未登录，请登录后查看',
    loginTimeout: '登录超时，需重新登录',
  },
  'en-US': {
    notLogin: 'Not logged in, please log in to view',
    loginTimeout: 'Login timeout, need to log in again',
  },
};

/**
 * 更新token
 * @param body
 * @param options
 * @returns
 */
export async function updateToken(
  body?: any,
  options?: { [key: string]: any },
) {
  // 如果不存在更新token请求时，刷新token接口的request赋值给reqFlag

  if (!reqFlag) {
    reqFlag = requestUmi<any>(`/users/token/refresh`, {
      method: 'POST',
      data: body,
      ...(options || {}),
    });
  }
  const res = await reqFlag;
  // 请求完成之后重新赋值为null
  reqFlag = null;
  // 当刷新token 的接口返回的也是401时，需要重新登录
  if (res.code === 401) {
    message.destroy();
    message.info(loginTextMap[getLocale()]?.loginTimeout);
    history.push('/login');
  }

  return res;
}

/**
 * 基于umi-request 封装一层，目的为了处理401刷新token问题
 * @param url
 * @param options
 * @returns
 */
const request = async <T>(
  url: string,
  options: any,
): Promise<AxiosResponse<T, any>> => {
  // 基于umi/max的request
  const res: any = await requestUmi(url, options);
  // 如果当前返回值为401时，则表示需要更新token或者token不存在
  if (res.code === 401) {
    let token = cookie.getToken('TOKEN');
    // token不存在时跳转到登录页面
    if (!token) {
      message.destroy();
      message.info(loginTextMap[getLocale()]?.notLogin);
      history.push('/login');
    } else {
      // 调取更新token的接口
      const result = await updateToken();
      // 如果更新token的接口返回最新的token之后，重新发起请求
      if (result?.data?.token) {
        return await requestUmi(url, options);
      }
    }
  }
  return res;
};

export default request;
