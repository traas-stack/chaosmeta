import { PageContainer } from '@ant-design/pro-components';
import { Tabs } from 'antd';
import React, { useEffect } from 'react';
import BasicInfo from './BasicInfo';
import MemberManage from './MemberManage';
import TagManage from './TagManage';
import { Container } from './style';

const SpaceSetting: React.FC<unknown> = () => {
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
    },
  ];

  useEffect(() => {}, []);
  return (
    <PageContainer title="空间设置">
      <Container>
        {/* <Row gutter={24}>
            <Col span={8}>
              <Form.Item name={'spaceName'} label="空间名称">
                <ShowText ellipsis isEdit />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name={'createTime'} label="创建时间">
                <ShowText isTime />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name={'userCount'} label="成员数量">
                <ShowText />
              </Form.Item>
            </Col>
          </Row> */}
        <Tabs items={tabItems} />
      </Container>
    </PageContainer>
  );
};

export default SpaceSetting;
