/**
 * 添加成员弹窗
 */

import { SearchOutlined } from '@ant-design/icons';
import { Button, Divider, Form, Input, Modal, Radio, Select } from 'antd';
import React from 'react';

interface IProps {
  open: boolean;
  setOpen: (val: boolean) => void;
}

const AddMemberModal: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;
  const [form] = Form.useForm();

  const Users = () => {
    return (
      <>
        <Select
          mode="multiple"
          showSearch={false}
          defaultActiveFirstOption={false}
          // open
          dropdownRender={(menu) => {
            return (
              <>
                <Input
                  suffix={
                    <SearchOutlined
                      onClick={() => {
                        console.log('==');
                      }}
                    />
                  }
                  onPressEnter={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    console.log('====');
                  }}
                />
                {menu}
                <div
                  style={{
                    padding: '6px 12px',
                    display: 'flex',
                    justifyContent: 'space-between',
                    borderTop: '1px solid rgba(5, 5, 5, 0.06)',
                  }}
                >
                  <a>重置</a>
                  <Button type="primary" size="small">
                    确定
                  </Button>
                </div>
              </>
            );
          }}
        >
          <Select.Option value={'ceshi'} key="1">
            diyi
          </Select.Option>
          <Select.Option value={'ceshi2'} key="2">
            diyi
          </Select.Option>
          <Select.Option value={'ceshi3'} key="3">
            diyi
          </Select.Option>
          <Select.Option value={'ceshi4'} key="4">
            diyi
          </Select.Option>
        </Select>
      </>
    );
  };
  return (
    <Modal
      title="添加成员"
      open={open}
      onCancel={() => {
        setOpen(false);
      }}
    >
      <Divider />
      <Form form={form} layout="vertical">
        <Form.Item label="用户名" name={'users'}>
          <Users />
        </Form.Item>
        <Form.Item label="用户权限" name={'role'}>
          <Radio.Group>
            <Radio value={'readonly'}>只读</Radio>
            <Radio value={'admain'}>读写</Radio>
          </Radio.Group>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default React.memo(AddMemberModal);
