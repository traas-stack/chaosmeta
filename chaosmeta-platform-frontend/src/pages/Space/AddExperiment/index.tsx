import ShowText from '@/components/ShowText';
import { CheckOutlined, CloseOutlined, EditOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { Button, Form, Input, Space, Tag } from 'antd';
import { useEffect, useState } from 'react';
import ArrangeContent from './ArrangeContent';
import { Container } from './style';

const initList = [
  { id: '999', children: [] },
  {
    id: 1,
    children: [
      {
        id: '1-1',
        second: 30,
      },
      {
        id: '1-2',
        second: 30,
      },
      {
        id: '1-3',
        second: 30,
      },
      {
        id: '1-4',
        second: 30,
      },
      {
        id: '1-5',
        second: 30,
      },
      {
        id: '1-6',
        second: 30,
      },
    ],
  },
  {
    id: 2,
    children: [
      {
        id: '2-1',
        second: 20,
      },
      {
        id: '2-2',
        second: 20,
      },
      {
        id: '2-3',
        second: 40,
      },
    ],
  },
  {
    id: 3,
    children: [
      {
        id: '3-1',
        second: 40,
      },
      {
        id: '3-2',
        second: 40,
      },
      {
        id: '3-3',
        second: 40,
      },
    ],
  },
  {
    id: 4,
    children: [
      {
        id: '4-1',
        second: 10,
      },
      {
        id: '4-2',
        second: 5,
      },
      {
        id: '4-3',
        second: 90,
      },
    ],
  },
  { id: '9990', children: [] },
];

const leftList = [
  {
    id: 'a',
    name: '故障节点',
    children: [
      {
        id: 'a-1',
        name: '容器',
        children: [
          {
            id: 'a-1-1',
            name: 'CPU燃烧',
            nodeType: 'fault',
            second: 40,
          },
          {
            id: 'a-1-2',
            name: 'CPU负载高',
            nodeType: 'fault',
            second: 40,
          },
          {
            id: 'a-1-3',
            name: 'CPU缓存',
            nodeType: 'fault',
            second: 40,
          },
          {
            id: 'a-1-4',
            name: '内存燃烧',
            nodeType: 'fault',
            second: 40,
          },
        ],
      },
      {
        id: 'a-2',
        name: '物理机',
        children: [
          {
            id: 'a-2-1',
            name: 'CPU燃烧',
            nodeType: 'fault',
          },
          {
            id: 'a-2-2',
            name: 'CPU负载高',
          },
          {
            id: 'a-2-3',
            name: 'CPU缓存',
          },
          {
            id: 'a-2-4',
            name: '内存燃烧',
          },
        ],
      },
    ],
  },
  {
    id: 'b',
    name: '校验节点',
    children: [
      {
        id: 'b-1',
        name: 'HTTP',
        nodeType: 'other',
        second: 112,
      },
      {
        id: 'b-2',
        name: 'Postman',
        nodeType: 'other',
        second: 40,
      },
    ],
  },
];

const AddExperiment = () => {
  const [editTitleState, setEditTitleState] = useState<boolean>(false);
  const [form] = Form.useForm();
  // 编排的数据
  const [arrangeList, setArrangeList] = useState(initList);
  // 左侧节点数据
  const [leftNodeList, setLeftNodeList] = useState(leftList);
  // 标题渲染

  const renderTitle = () => {
    return (
      <Form form={form}>
        <Space>
          {editTitleState ? (
            <Form.Item name={'title'} noStyle>
              <Input placeholder="请输入" />
            </Form.Item>
          ) : (
            <Form.Item name={'title'} label="触发方式">
              <ShowText ellipsis />
            </Form.Item>
          )}
          {editTitleState ? (
            <Space>
              <CloseOutlined
                className="cancel"
                onClick={() => {
                  setEditTitleState(false);
                }}
              />
              <CheckOutlined
                className="confirm"
                onClick={() => {
                  setEditTitleState(false);
                }}
              />
            </Space>
          ) : (
            <EditOutlined
              className="edit"
              style={{ color: '#1890FF' }}
              onClick={() => {
                setEditTitleState(true);
              }}
            />
          )}
        </Space>
        <div className="ellipsis tags">
          <Tag>标签</Tag>
          <Tag>标签</Tag>
          <Tag>标签</Tag>
          <Tag>标签</Tag>
          <Tag>标签</Tag>
          <Tag>标签</Tag>
          <Tag>标签</Tag>
          <Tag>标签</Tag>
        </div>
      </Form>
    );
  };

  const headerExtra = () => {
    return (
      <Form form={form}>
        <div className="header-extra">
          <div>
            <Form.Item name={'triggerMode'} label="触发方式">
              <ShowText />
            </Form.Item>
            <Form.Item name={'desc'} label="描述">
              <ShowText />
            </Form.Item>
          </div>
          <Space>
            <Button ghost danger>
              删除
            </Button>
            <Button ghost type="primary" onClick={() => {}}>
              完成
            </Button>
          </Space>
        </div>
      </Form>
    );
  };

  useEffect(() => {
    form.setFieldsValue({
      title:
        '实验名称实验名称实验名称实验名称实验名称实验名称实验名称实验名称实验名称实验名称实验名称',
      triggerMode: '手动触发',
      desc: '我是描述',
    });
    // document.body.style.overflow = 'hidden'
  }, []);

  return (
    <Container>
      <>
        <PageContainer
          header={{
            title: renderTitle(),
            onBack: () => {
              history.push('/space/experiment');
            },
            extra: headerExtra(),
          }}
        >
          <ArrangeContent
            arrangeList={arrangeList}
            leftNodeList={leftNodeList}
            setArrangeList={setArrangeList}
          />
        </PageContainer>
      </>
    </Container>
  );
};

export default AddExperiment;
