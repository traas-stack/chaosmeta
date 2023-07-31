import { Form, Modal, Select } from "antd";
import React from "react";

interface IProps {
  open: boolean;
  setOpen: (val: boolean) => void;
}
const AddTagModal: React.FC<IProps> = (props) => {
  const {open, setOpen}  = props
  const [form] = Form.useForm();

  return <Modal title="新建标签" open={open} onCancel={() => {
    setOpen(false)
  }}>
    <Form form={form} layout="vertical">
      <Form.Item name={'tags'} label="标签">
        <Select>
          <Select.Option>
            
          </Select.Option>
        </Select>
      </Form.Item>
    </Form>

  </Modal>

}

export default React.memo(AddTagModal) 