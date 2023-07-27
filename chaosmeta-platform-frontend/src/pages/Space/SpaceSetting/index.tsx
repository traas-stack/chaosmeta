import ShowText from '@/components/ShowText';
import { PageContainer } from '@ant-design/pro-components';
import { Col, Form, Row, Tabs } from 'antd';
import React, { useEffect } from 'react';
import MemberManage from './MemberManage';
import { Container } from './style';
import TagManage from './TagManage';

const SpaceSetting: React.FC<unknown> = () => {
  const [form] = Form.useForm();
  const tabItems = [
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

  useEffect(() => {
    console.log(new Date(), 'date');
    form.setFieldsValue({
      spaceName:
        '空间名称空间名称空间名称空间名称空间名称空间名称空间名称空间名称',
      createTime: new Date(),
      userCount: '12',
    });
  }, []);
  return (
    <PageContainer title="空间设置">
      <Container>
        <Form form={form}>
          <Row gutter={24}>
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
          </Row>
          <Tabs items={tabItems} />
        </Form>
      </Container>
    </PageContainer>
  );
};

export default SpaceSetting;
