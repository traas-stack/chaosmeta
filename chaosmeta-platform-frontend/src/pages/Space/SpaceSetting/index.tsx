import { useParamChange } from '@/utils/useParamChange';
import { PageContainer } from '@ant-design/pro-components';
import { history, useIntl } from '@umijs/max';
import { Tabs } from 'antd';
import React, { useEffect, useState } from 'react';
import BasicInfo from './BasicInfo';
import MemberManage from './MemberManage';
import TagManage from './TagManage';
import { Container } from './style';
// import { createHistory } from '@umijs/max';

const SpaceSetting: React.FC<unknown> = () => {
  const [activeKey, setActiveKey] = useState<string>('basic');
  const tabKeyChange = useParamChange('tabKey');
  const intl = useIntl();
  const tabItems = [
    {
      label: intl.formatMessage({ id: 'basicInfo' }),
      key: 'basic',
      children: <BasicInfo />,
    },
    {
      label: intl.formatMessage({ id: 'spaceSetting.tab.member' }),
      key: 'user',
      children: <MemberManage />,
    },
    {
      label: intl.formatMessage({ id: 'spaceSetting.tab.tag' }),
      key: 'tag',
      children: <TagManage />,
    },
    // {
    //   label: '实验攻击范围配置',
    //   key: 'range',
    //   children: <AttackRange />,
    // },
  ];

  useEffect(() => {
    if (history.location.query.tabKey) {
      setActiveKey(history.location.query.tabKey as string);
    }
  }, [tabKeyChange]);

  return (
    <PageContainer title={intl.formatMessage({ id: 'spaceSetting' })}>
      <Container>
        <Tabs
          items={tabItems}
          activeKey={activeKey}
          onChange={(val) => {
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
