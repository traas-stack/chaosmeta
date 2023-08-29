// import { useLocation } from '@umijs/max';

import { history, useModel } from '@umijs/max';

// /**
//  * 用于监听路由中指定参数的变化
//  * @param key 参数key
//  * @returns
//  */
// export const useParamChange = (key: string) => {
//   const location = useLocation();
//   // 获取当前指定字符的下标
//   const index = location.search.indexOf(key);
//   // 截取字符，使用&分隔，取数组第一个元素即是参数的值
//   if (index || index === 0) {
//     const value = location.search?.slice(index + key.length + 1)?.split('&')[0];
//     return value;
//   }
//   return undefined;
// };

export const useUpdateSpace = (spaceId: string, spaceName: string) => {
  const { setCurSpace } = useModel('global');
  if (spaceId) {
    history.push({
      pathname: history.location.pathname,
      query: {
        ...history.location.query,
        spaceId,
      },
    });
    setCurSpace([spaceId]);
    sessionStorage.setItem('spaceId', spaceId);
    sessionStorage.setItem('spaceName', spaceName);
  }
};
