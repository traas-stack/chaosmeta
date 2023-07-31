import { createSpace } from '@/services/chaosmeta/SpaceController';
import { history, useRequest } from '@umijs/max';
import { Button, Drawer, Form, Input, Space, message } from 'antd';
import React from 'react';

interface IProps {
  open: boolean;
  setOpen: (open: boolean) => void;
}

const AddSpaceDrawer: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;
  const [form] = Form.useForm();

  const handleCancel = () => {
    setOpen(false);
  };

  /**
   * 创建空间接口
   */
  const create = useRequest(createSpace, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      console.log(res, '----');
      if (res.code === 200) {
        message.success('创建成功');
        setOpen(false)
        history.push('/space/setting');
      }
    },
  });

  /**
   * 创建空间
   */
  const handleCreate = () => {
    form.validateFields().then((values) => {
      console.log(values, 'values');
      create?.run(values);
    });
  };

  return (
    <Drawer
      title="新建空间"
      open={open}
      onClose={handleCancel}
      width={480}
      footer={
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button onClick={handleCancel}>取消</Button>
            <Button
              type="primary"
              onClick={handleCreate}
              loading={create.loading}
            >
              创建完成并去配置
            </Button>
          </Space>
        </div>
      }
    >
      <Form layout="vertical" form={form}>
        <Form.Item
          name={'name'}
          label="空间名称"
          rules={[{ required: true, message: '请输入' }]}
          help="请尽量保持项目名称的简洁，不超过 64 个字符"
        >
          <Input placeholder="请输入空间名称" maxLength={64} />
        </Form.Item>
        <Form.Item
          name={'description'}
          label="空间描述"
          style={{ marginTop: '36px' }}
        >
          <Input.TextArea placeholder="请输入空间描述" rows={5} />
        </Form.Item>
      </Form>
    </Drawer>
  );
};

export default React.memo(AddSpaceDrawer);