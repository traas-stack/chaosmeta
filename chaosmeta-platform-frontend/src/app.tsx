// 运行时配置
// 全局初始化数据配置，用于 Layout 用户信息和权限初始化

import { RequestConfig, history } from '@umijs/max';
import SpaceDropdown from './components/SpaceDropdown';
import UserRightArea from './components/UserRightArea';
import cookie from './utils/cookie';
import errorHandler from './utils/errorHandler';

// 更多信息见文档：https://umijs.org/docs/api/runtime-config#getinitialstate
// export async function getInitialState(): Promise<{
//   userInfo: {
//     name?: string;
//     avatar: string;
//     role?: string;
//   };
// }> {
//   let userInfo = {
//     avatar:
//       'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*RG7jSIPO-pQAAAAAAAAAAAAADmKmAQ/original',
//   };
//   // 获取当前用户信息
//   const userName = localStorage.getItem('userName') || undefined;
//   if (userName && history.location.pathname !== '/login') {
//     // const userResult = await getUserInfo({ name: userName });
//     userInfo = {
//       ...userInfo,
//       // ...userResult?.data,
//     };
//   }
//   return { userInfo };
// }

/**
 * 请求配置
 */
export const request: RequestConfig = {
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
  errorConfig: {
    errorThrower: (res: any) => {
      const { success, message } = res;
      if (!success) {
        const error: any = new Error(message);
        error.name = 'BizError';
        error.info = { ...res };
        throw error; // 抛出自制的错误进入errorHandler统一处理
      }
    },
    errorHandler,
  },
  requestInterceptors: [
    (url: string, options: any) => {
      // 获取token后拼接Bearer
      let token = cookie.getToken('TOKEN')
        ? `Bearer ${cookie.getToken('TOKEN')}`
        : '';
      // 刷新token的接口需要用REFRESH_TOKEN作为token
      if (url === '/users/token/refresh') {
        token = cookie.getToken('REFRESH_TOKEN')
          ? `Bearer ${cookie.getToken('REFRESH_TOKEN')}`
          : '';
      }
      return {
        url,
        options: {
          ...options,
          headers: {
            ...options.headers,
            Authorization: token,
          },
        },
      };
    },
  ],
  // 响应拦截器
  responseInterceptors: [
    (response: any) => {
      // code为200或401时，设置为true，不进行错误处理，401时涉及token刷新及重新发起请求问题，在封装层的requst中处理
      response.data.success =
        response?.data?.code === 200 || response?.data?.code === 401;
      return response;
    },
  ],
};

export const layout = () => {
  return {
    logo: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*lMXkRKmd8WcAAAAAAAAAAAAADmKmAQ/original',
    title: '',
    menu: {
      locale: false,
    },
    siderWidth: 208,
    layout: 'mix',
    splitMenus: true,
    rightContentRender: () => {
      return <UserRightArea />;
    },
    menuExtraRender: (props: any) => {
      if (props?.matchMenuKeys[0] === '/setting') {
        return null;
      }
      return <SpaceDropdown />;
    },
    onPageChange: () => {
      // 页面切换时需要加上空间id参数, 全局设置下不需要
      const spaceId = sessionStorage.getItem('spaceId');
      const { pathname, query } = history?.location || {};
      if (
        !history.location.query.spaceId &&
        pathname.split('/')[1] === 'space'
      ) {
        history.push({
          pathname,
          query: { ...query, spaceId },
        });
      }
    },
  };
};
