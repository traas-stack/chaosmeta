// 运行时配置
// 全局初始化数据配置，用于 Layout 用户信息和权限初始化

import { RequestConfig } from '@umijs/max';
import SpaceDropdown from './components/SpaceDropdown';

// 更多信息见文档：https://umijs.org/docs/api/runtime-config#getinitialstate
export async function getInitialState(): Promise<{
  name: string;
  avatar: string;
}> {
  return {
    name: 'Serati Ma',
    avatar: 'https://img.alicdn.com/tfs/TB1YHEpwUT1gK0jSZFhXXaAtVXa-28-27.svg',
  };
}

export const request: RequestConfig = {
  timeout: 1000,
  // other axios options you want
  errorConfig: {
    errorHandler() {},
    errorThrower() {},
  },
  requestInterceptors: [
    (url: string, options: any) => {
      return {
        url,
        options: {
          ...options,
          headers: {
            ...options.headers,
          },
        },
      };
    },
  ],
  responseInterceptors: [],
};

export const layout = (props: {
  initialState: { name: string; avatar: string };
}) => {
  const { initialState } = props;
  console.log(props, 'props');
  return {
    logo: 'https://img.alicdn.com/tfs/TB1YHEpwUT1gK0jSZFhXXaAtVXa-28-27.svg',
    title: 'ceshi',
    menu: {
      locale: false,
    },
    layout: 'mix',
    breakpoint: false,
    // collapsedButtonRender: () => {
    //   return null;
    // },
    // collapsed: false,
    // rightContentRender: () => {
    //   return <>99888</>
    // },
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
    menuHeaderRender: () => {
      return <SpaceDropdown />;
    },
  };
};
