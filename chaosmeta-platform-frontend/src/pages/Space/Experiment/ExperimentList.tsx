import { LightArea } from '@/components/CommonStyle';
import { QuestionCircleOutlined } from '@ant-design/icons';
import {
  Badge,
  Button,
  Col,
  Form,
  Input,
  Popover,
  Row,
  Space,
  Table,
  Tag,
  Tooltip,
} from 'antd';
import { ColumnsType } from 'antd/es/table';
import React from 'react';
/**
 * 实验列表
 * @returns
 */
const ExperimentList: React.FC<unknown> = () => {
  const [form] = Form.useForm();

  const columns: ColumnsType<any> = [
    {
      title: '名称',
      width: 160,
      dataIndex: 'name',
      render: () => {
        return (
          <>
            <div className="ellipsis row-text-gap">
              <a>
                <span>
                  我名字很长-我名字很长-我名字很长-我名字很长-我名字很长
                </span>
              </a>
            </div>
            <div className="ellipsis">
              <Popover
                title={
                  <div>
                    <Tag>基础设施</Tag>
                    <Tag>基础设施</Tag>
                    <Tag>基础设施</Tag>
                  </div>
                }
              >
                <Tag>基础设施</Tag>
                <Tag>基础设施</Tag>
                <Tag>基础设施</Tag>
              </Popover>
            </div>
          </>
        );
      },
    },
    {
      title: '创建人',
      width: 80,
      dataIndex: 'create',
    },
    {
      title: '实验次数',
      width: 120,
      dataIndex: 'count',
      sorter: true,
      render: () => {
        return <a>9</a>;
      },
    },
    {
      title: (
        <>
          <span>最近试验时间</span>
          <Tooltip title="todo">
            <QuestionCircleOutlined />
          </Tooltip>
          /<span>状态</span>
        </>
      ),
      width: 180,
      dataIndex: 'shaiy',
      sorter: true,
      render: () => {
        return (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
            }}
          >
            <div>2020-01-04 09:41:00</div>
            <div style={{ flexShrink: 0 }}>
              <Badge color="#f50" /> 成功
            </div>
          </div>
        );
      },
    },
    {
      title: '触发方式',
      width: 120,
      dataIndex: 'chuf',
      filters: [
        {
          text: '手动',
          value: 'shoudong',
        },
        {
          text: '周期性',
          value: 'zhouqi',
        },
      ],
      render: () => {
        return (
          <div>
            <div>周期性</div>
            <span className="cycle">每周二14:00</span>
          </div>
        );
      },
    },
    {
      title: '即将运行时间',
      width: 140,
      dataIndex: 'jijiang',
      sorter: true,
      render: () => {
        return (
          <div>
            <div className="run-finish">2020-01-04 09:41:00</div>
          </div>
        );
      },
    },
    {
      title: '最近编辑时间',
      width: 140,
      dataIndex: 'zuijin',
      sorter: true,
    },
    {
      title: '操作',
      width: 140,
      fixed: 'right',
      render: () => {
        return (
          <Space>
            <a onClick={() => {}}>编辑</a>
            <a onClick={() => {}}>复制</a>
            <a onClick={() => {}}>删除</a>
          </Space>
        );
      },
    },
  ];
  return (
    <div className="experiment-list ">
      <LightArea className="search">
        <Form form={form} labelCol={{ span: 6 }}>
          <Row>
            <Col span={8}>
              <Form.Item name={'name'} label="名称">
                <Input placeholder="请输入" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name={'creator'} label="创建人">
                <Input placeholder="请输入" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name={'type'} label="触发方式">
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
            <div className="title">成员列表</div>
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
  );
};

export default ExperimentList;
