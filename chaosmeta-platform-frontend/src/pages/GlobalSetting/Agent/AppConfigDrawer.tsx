/**
 * 应用配置抽屉
 */
import { Button, Drawer, Form, Input, Space } from 'antd';
import React from 'react';
import { AppConfigContainer } from './style';
interface IProps {
  open: boolean;
  setOpen: (val: boolean) => void;
}
const AppConfigDrawer: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;
  const [form] = Form.useForm();
  const handleClose = () => {
    setOpen(false);
  };
  return (
    <Drawer
      open={open}
      onClose={handleClose}
      width={560}
      title="应用配置"
      footer={
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button onClick={handleClose}>取消</Button>
            <Button type="primary">提交</Button>
          </Space>
        </div>
      }
    >
      <AppConfigContainer>
        <div className="desc">
          根据K8s标签识别应用，需维护好应用映射关系：
          <br />
          如Labels：ChaosMeta/app:应用名
        </div>
        <Form form={form} layout="vertical">
          <Form.Item name={'appName'} label="应用名Key值：">
            <Input placeholder="app的名称" />
          </Form.Item>
        </Form>
      </AppConfigContainer>
    </Drawer>
  );
};
export default React.memo(AppConfigDrawer);
