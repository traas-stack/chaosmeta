import { Button, Drawer, Form, Input, Space } from 'antd';
import React from 'react';

interface IProps {
  open: boolean;
  setOpen: (open: boolean) => void;
}

const AddSpaceDrawer: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;

  const handleCancel = () => {
    setOpen(false);
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
            <Button>取消</Button>
            <Button type="primary">创建完成并去配置</Button>
          </Space>
        </div>
      }
    >
      <Form layout="vertical">
        <Form.Item
          name={'spaceName'}
          label="空间名称"
          rules={[{ required: true, message: '请输入空间名称' }]}
          help="请尽量保持项目名称的简洁，不超过 64 个字符"
        >
          <Input placeholder="请输入空间名称" />
        </Form.Item>
        <Form.Item
          name={'spaceDesc'}
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
