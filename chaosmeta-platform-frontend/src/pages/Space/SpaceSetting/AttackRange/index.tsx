/**
 * 实验结果页面
 * @returns
 */
import { LightArea } from '@/components/CommonStyle';
import { Button, Col, Form, Input, Row, Space } from 'antd';
import Table, { ColumnsType } from 'antd/es/table';
import React, { useState } from 'react';
import { AttackRangeContainer } from '../style';

const ExperimentResult: React.FC<unknown> = () => {
  const [form] = Form.useForm();
  const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([]);
  const columns: ColumnsType<any> = [
    {
      title: 'Hostname',
      width: 120,
      dataIndex: 'name',
      render: () => {
        return <a>xxxx</a>;
      },
    },
    {
      title: 'IP',
      width: 120,
      dataIndex: 'ip',
    },
    {
      title: '操作',
      width: 80,
      fixed: 'right',
      render: () => {
        return (
          <Space>
            <a onClick={() => {}}>移除</a>
          </Space>
        );
      },
    },
  ];
  return (
    <AttackRangeContainer>
      <LightArea>
        <div className="search">
          <Form form={form} labelCol={{ span: 6 }}>
            <Row>
              <Col span={8}>
                <Form.Item name={'name'} label="集群">
                  <Input placeholder="请输入" />
                </Form.Item>
              </Col>
              {/* <Col style={{ textAlign: 'right' }} span={16}>
                <Space>
                  <Button>重置</Button>
                  <Button type="primary">查询</Button>
                </Space>
              </Col> */}
            </Row>
          </Form>
        </div>
        <div className="table">
          <div className="operate">
            <div>集群外host</div>
            <Space>
              <Button>批量移除</Button>
              <Button type="primary">添加Host</Button>
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
            rowSelection={{
              selectedRowKeys,
              onChange: (rowKeys: any[]) => {
                setSelectedRowKeys(rowKeys);
              },
            }}
          />
        </div>
      </LightArea>
    </AttackRangeContainer>
  );
};

export default ExperimentResult;
