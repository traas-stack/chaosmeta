import ShowText from '@/components/ShowText';
import { CheckOutlined, CloseOutlined, EditOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { Button, Form, Input, Space, Tag } from 'antd';
import { useEffect, useState } from 'react';
import { Container } from './style';

const AddExperiment = () => {
  const [editTitleState, setEditTitleState] = useState<boolean>(false);
  const [form] = Form.useForm();
  // 标题渲染
  const renderTitle = () => {
    return (
      <div>

      <Space>
        {editTitleState ? (
          <Form.Item name={'title'} noStyle>
            <Input placeholder="请输入" />
          </Form.Item>
        ) : (
          <div className="ellipsis">实验名称实验名称实验名称实验名称</div>
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
      <div className='ellipsis tags'>
        <Tag>标签</Tag>
        <Tag>标签</Tag>
        <Tag>标签</Tag>
        <Tag>标签</Tag>
        <Tag>标签</Tag>
        <Tag>标签</Tag>
        <Tag>标签</Tag>
        <Tag>标签</Tag>
      </div>
      </div>

    );
  };

  const headerExtra = () => {
    return (
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
          <Button ghost type="primary">
            完成
          </Button>
        </Space>
      </div>
    );
  };

  useEffect(() => {
    form.setFieldsValue({
      title:
        '实验名称实验名称实验名称实验名称实验名称实验名称实验名称实验名称实验名称实验名称实验名称',
      triggerMode: '手动触发',
      desc: '我是描述',
    });
  }, []);

  return (
    <Container>
      <Form form={form}>
        <PageContainer
          header={{
            title: renderTitle(),
            onBack: () => {
              history.push('/space/experiment');
            },
            extra: headerExtra(),
          }}
        >

          <div>
            <div className=''></div>
          </div>
        </PageContainer>
      </Form>
    </Container>
  );
};

export default AddExperiment;
