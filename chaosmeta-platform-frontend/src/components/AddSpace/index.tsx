import { Divider, Form, Input, Modal } from 'antd';
import React from 'react';

interface IProps {
  open: boolean;
  setOpen: (open: boolean) => void;
}

const AddSpace: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;

  const handleCancel = () => {
    setOpen(false);
  };

  return (
    <Modal title="新建空间" open={open} onCancel={handleCancel}>
      <Divider style={{ margin: '16px 0' }} />
      <Form layout="vertical">
        <Form.Item
          name={'spaceName'}
          label="工作空间名称"
          rules={[{ required: true, message: '请输入空间名称' }]}
        >
          <Input placeholder="请输入空间名称" />
        </Form.Item>
        <Form.Item name={'spaceDesc'} label="空间描述">
          <Input.TextArea placeholder="空间描述" rows={5} />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default React.memo(AddSpace);
