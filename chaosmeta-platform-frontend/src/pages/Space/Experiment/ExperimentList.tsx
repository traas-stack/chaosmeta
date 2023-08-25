import { LightArea } from '@/components/CommonStyle';
import EmptyCustom from '@/components/EmptyCustom';
import TimeTypeRangeSelect from '@/components/Select/TimeTypeRangeSelect';
import { QuestionCircleOutlined } from '@ant-design/icons';
import { history } from '@umijs/max';
import {
  Badge,
  Button,
  Col,
  Form,
  Input,
  Popover,
  Row,
  Select,
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

  const triggerMode = [
    {
      text: '手动',
      value: 'shoudong',
    },
    {
      text: '单次定时',
      value: 'danci',
    },
    {
      text: '周期性',
      value: 'zhouqi',
    },
  ];

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
      filters: triggerMode,
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
          <Row gutter={16} style={{ whiteSpace: 'nowrap', overflow: 'hidden' }}>
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
                <Select placeholder="请选择">
                  {triggerMode?.map((item) => {
                    return (
                      <Select.Option key={item.value} value={item.value}>
                        {item.text}
                      </Select.Option>
                    );
                  })}
                </Select>
              </Form.Item>
            </Col>
            <TimeTypeRangeSelect />
            <Col style={{ textAlign: 'right', flex: 1 }}>
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
              <Button
                type="primary"
                onClick={() => {
                  history.push('/space/experiment/choose');
                }}
              >
                创建实验
              </Button>
            </Space>
          </div>
          <Table
            columns={columns}
            locale={{
              emptyText: (
                <EmptyCustom
                  desc="请先创建实验，您可选择自己创建实验也可以通过推荐实验来快速构建实验场景，来验证应用系统的可靠性。"
                  topTitle="当前空间还没有实验数据"
                  btns={
                    <Space>
                      <Button
                        type="primary"
                        onClick={() => {
                          history.push({
                            pathname: '/space/setting',
                            query: {
                              tabKey: 'user',
                              spaceId: history?.location?.query?.spaceId,
                            },
                          });
                        }}
                      >
                        创建实验
                      </Button>
                      <Button
                        onClick={() => {
                          history.push({
                            pathname: '/space/setting',
                            query: {
                              tabKey: 'user',
                              spaceId: history?.location?.query?.spaceId,
                            },
                          });
                        }}
                      >
                        推荐实验
                      </Button>
                    </Space>
                  }
                />
              ),
            }}
            rowKey={'id'}
            // dataSource={[{ id: '1', account: 'hlt', auth: 'admain' }]}
            dataSource={[]}
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
