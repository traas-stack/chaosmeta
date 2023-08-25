/**
 * 添加集群
 */
import { Divider, Form, Input, Modal } from 'antd';
import React from 'react';
interface IProps {
  open: boolean;
  setOpen: (val: boolean) => void;
}
const AddColonyModal: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;
  const [form] = Form.useForm();
  const handleClose = () => {
    setOpen(false);
  };
  return (
    <Modal
      open={open}
      onCancel={handleClose}
      onOk={() => {}}
      title="添加集群"
      width={480}
    >
      <Divider style={{marginTop: '16px'}}/>
      <Form form={form}>
        <Form.Item name={'name'}>
          <Input.TextArea
            placeholder="请把目标K8s集群的配置文件拷贝进来一般为“/home/.kube/config”"
            rows={4}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};
export default React.memo(AddColonyModal);
