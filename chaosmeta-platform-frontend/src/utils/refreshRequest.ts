import { request } from '@umijs/max';
import { message } from 'antd';

export const refreshRequest = (response: { config: any; headers: any; }) => {
  const { config, headers } = response;
  const {method, data, params, url} = config || {};

  // get请求
  const againRequest = () => {
    return request<any>(url, {
      method,
      data: config.data,
      params: config.params,
      headers: {...headers},
      // ...(options || {}),
    });
  };

  // refreshtoken过期需要重新登录
  if(url === '/users/token/refresh') {
    message.error('登录过期，请重新登录')
  } else {
    againRequest().then()
  }
};
