import { PageContainer } from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { Tabs } from 'antd';
import React, { useEffect, useState } from 'react';
import AttackRange from './AttackRange';
import BasicInfo from './BasicInfo';
import MemberManage from './MemberManage';
import TagManage from './TagManage';
import { Container } from './style';
// import { createHistory } from '@umijs/max';

const SpaceSetting: React.FC<unknown> = () => {
  const [activeKey, setActiveKey] = useState<string>('basic');
  const tabItems = [
    {
      label: '基本信息',
      key: 'basic',
      children: <BasicInfo />,
    },
    {
      label: '成员管理',
      key: 'user',
      children: <MemberManage />,
    },
    {
      label: '标签管理',
      key: 'tag',
      children: <TagManage />,
    },
    {
      label: '实验攻击范围配置',
      key: 'range',
      children: <AttackRange />,
    },
  ];

  useEffect(() => {
    if (history.location.query.tabKey) {
      setActiveKey(history.location.query.tabKey as string);
    }
  }, []);
  return (
    <PageContainer title="空间设置">
      <Container>
        <Tabs
          items={tabItems}
          activeKey={activeKey}
          onChange={(val) => {
            console.log(val, 'val===');
            setActiveKey(val);
            history.push({
              pathname: history.location.pathname,
              query: { ...history.location.query, tabKey: val },
            });
          }}
        />
      </Container>
    </PageContainer>
  );
};

export default SpaceSetting;
