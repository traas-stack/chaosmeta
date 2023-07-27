// 运行时配置
// 全局初始化数据配置，用于 Layout 用户信息和权限初始化

import { RequestConfig, history, request as requ } from '@umijs/max';
import { Dropdown, Space } from 'antd';
import cookie from 'react-cookies';
import SpaceDropdown from './components/SpaceDropdown';
import { updateToken } from './services/chaosmeta/UserController';
import errorHandler from './utils/errorHandler';

let count = 1;

// 更多信息见文档：https://umijs.org/docs/api/runtime-config#getinitialstate
export async function getInitialState(): Promise<{
  userInfo: {
    name?: string;
    avatar: string;
    role?: string;
  };
}> {
  let userInfo = {
    avatar:
      'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*RG7jSIPO-pQAAAAAAAAAAAAADmKmAQ/original',
  };
  // 获取当前用户信息
  const userName = localStorage.getItem('userName') || undefined;
  if (userName && history.location.pathname !== '/login') {
    // const userResult = await getUserInfo({ name: userName });
    userInfo = {
      ...userInfo,
      // ...userResult?.data,
    };
  }
  return { userInfo };
}

/**
 * 相应拦截
 * @param response
 * @returns
 */
const responseInterceptors = async (response: any): Promise<any> => {
  // code401时token失效，需要刷新token或者重新登录
  if (response.data.code === 401 && count < 3) {
    response.data.success = true;
    // 刷新token接口仍然为401，则重新登录
    if (response.config.url === '/users/token/refresh') {
      history.push('/login');
      return;
    }
    // 获取新的token存到cookie中，重新调取接口
    const { config } = response;
    const { method, data, params, url } = config || {};
    const newToken: any = await updateToken();
    cookie.save('TOKEN', newToken.data?.token, {
      domain: document.domain,
    });
    // 重新调取接口
    const againRequest = () => {
      return requ<any>(url, {
        method,
        data,
        params,
        // headers: { ...headers },
        // ...(options || {}),
      });
    };
    await againRequest();
  }
  console.log(response, 'response===');
  // 401时另做处理
  response.data.success =
    response?.data?.code === 200 || response?.data?.code === 401;
  return response;
};

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
      // const test = await updateToken();
      // 获取token后拼接Bearer
      let token = cookie.load('TOKEN') ? `Bearer ${cookie.load('TOKEN')}` : '';
      if (url === '/users/token/refresh') {
        token = cookie.load('REFRESH_TOKEN')
          ? `Bearer ${cookie.load('REFRESH_TOKEN')}`
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
  // responseInterceptors: [
  //   (response: any) => {
  //     console.log(response, 'response===');
  //     // 拦截响应数据，进行个性化处理 200成功其他失败，设置success为false进入到 errorThrower进行自定义错误处理
  //     // token过期
  //     if(response.data.code === 401) {
  //       // if(response.config.url )
  //     }
  //     response.data.success = false;
  //     // response.data.success = response?.data?.code === 200;
  //     requ(response.config)
  //     return response;
  //   },
  // ],
  responseInterceptors: [responseInterceptors],
};

export const layout = (props: { initialState: { userInfo: any } }) => {
  const { initialState } = props;
  return {
    logo: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*AH-NT5UMv3UAAAAAAAAAAAAADmKmAQ/original',
    title: 'Chaosmeta',
    menu: {
      locale: false,
    },
    siderWidth: 208,
    layout: 'mix',
    // breakpoint: false,
    splitMenus: true,
    // collapsedButtonRender: () => {
    //   return null;
    // },
    // collapsed: false,
    rightContentRender: (param) => {
      const items = [
        {
          label: <span>修改密码</span>,
          key: 'updatePassword',
        },
        {
          label: (
            <span
              onClick={() => {
                cookie.remove('TOKEN');
                cookie.remove('REFRESH_TOKEN');
                localStorage.removeItem('userName');
                history.push('/login');
              }}
            >
              退出登录
            </span>
          ),
          key: 'logout',
        },
      ];
      return (
        <Dropdown menu={{ items }}>
          <Space style={{ cursor: 'pointer' }}>
            <img src={initialState?.userInfo?.avatar} />
            <span>{initialState?.userInfo?.name}</span>
          </Space>
        </Dropdown>
      );
    },
    // headerRender: (layoutProps: any) => {
    //   const { title, logo } = layoutProps;
    //   return (
    //     <div
    //       style={{
    //         display: 'flex',
    //         justifyContent: 'space-between',
    //         padding: '0 16px',
    //       }}
    //     >
    //       <div>
    //         {title}
    //         <img src={logo} />
    //       </div>
    //       <div>
    //         <img src={initialState?.avatar} />
    //         {initialState?.name}
    //       </div>
    //     </div>
    //   );
    // },
    menuExtraRender: (props: any) => {
      if (props?.matchMenuKeys[0] === '/setting') {
        return null;
      }
      return <SpaceDropdown />;
    },
  };
};
