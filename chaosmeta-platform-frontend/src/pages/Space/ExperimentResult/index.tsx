/**
 * 实验结果页面
 * @returns
 */
import { LightArea } from '@/components/CommonStyle';
import ShowText from '@/components/ShowText';
import { PageContainer } from '@ant-design/pro-components';
import { Badge, Button, Col, Form, Input, Row, Space, Table } from 'antd';
import { ColumnsType } from 'antd/es/table';
import React from 'react';
import { Container } from './style';

const ExperimentResult: React.FC<unknown> = () => {
  const [form] = Form.useForm();
  const columns: ColumnsType<any> = [
    {
      title: '名称',
      width: 80,
      dataIndex: 'name',
      render: () => {
        return <a>xxxx</a>;
      },
    },
    {
      title: '执行人',
      width: 80,
      dataIndex: 'name',
    },
    {
      title: '实验开始时间',
      width: 140,
      dataIndex: 'jijiang',
      sorter: true,
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
    {
      title: '实验结束时间',
      width: 140,
      dataIndex: 'jijiang',
      sorter: true,
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
    {
      title: '实验状态',
      width: 180,
      dataIndex: 'shaiy',
      sorter: true,
      render: () => {
        return (
          <div>
            <Badge color="#f50" /> 成功
          </div>
        );
      },
    },
    {
      title: '操作',
      width: 140,
      fixed: 'right',
      render: () => {
        return (
          <Space>
            <a onClick={() => {}}>停止</a>
          </Space>
        );
      },
    },
  ];
  return (
    <>
      <PageContainer title="实验结果">
        <Container>
          <div className="result-list ">
            <LightArea className="search">
              <Form form={form} labelCol={{ span: 6 }}>
                <Row>
                  <Col span={8}>
                    <Form.Item name={'name'} label="名称">
                      <Input placeholder="请输入" />
                    </Form.Item>
                  </Col>
                  <Col span={8}>
                    <Form.Item name={'creator'} label="执行人">
                      <Input placeholder="请输入" />
                    </Form.Item>
                  </Col>
                  <Col span={8}>
                    <Form.Item name={'type'} label="实验状态">
                      <Input placeholder="请输入" />
                    </Form.Item>
                  </Col>
                  <Col span={8}>
                    <Form.Item name={'time'} label="时间类型">
                      <Input placeholder="请输入" />
                    </Form.Item>
                  </Col>
                  <Col style={{ textAlign: 'right' }} span={16}>
                    <Space>
                      <Button>重置</Button>
                      <Button type="primary">查询</Button>
                    </Space>
                  </Col>
                </Row>
              </Form>
            </LightArea>
            <LightArea>
              <div className="table">
                <div className="area-operate">
                  <div className="title">实验结果列表</div>
                  <Space>
                    <Button type="primary" onClick={() => {}}>
                      创建实验
                    </Button>
                  </Space>
                </div>
                <Table
                  columns={columns}
                  rowKey={'id'}
                  dataSource={[{ id: '1', account: 'hlt', auth: 'admain' }]}
                  pagination={{
                    showQuickJumper: true,
                    total: 200,
                  }}
                  scroll={{ x: 760 }}
                />
              </div>
            </LightArea>
          </div>
        </Container>
      </PageContainer>
    </>
  );
};

export default ExperimentResult;
