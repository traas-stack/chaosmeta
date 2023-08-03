// 全局共享数据示例
import { useState } from 'react';
interface userInfo {
  name: string;
  role: string;
  avatar: string;
}

const useUser = () => {
  // 登录人信息
  const [userInfo, setUserInfo] = useState<userInfo>({
    avatar:
      'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*RG7jSIPO-pQAAAAAAAAAAAAADmKmAQ/original',
    name: '',
    role: 'normal',
  });
  // 空间id
  const [spaceId, setSpaceId] = useState<string>('');
  return {
    userInfo,
    setUserInfo,
    spaceId,
    setSpaceId,
  };
};

export default useUser;
