/**
 * 添加成员弹窗
 */

import { spaceAddUser } from '@/services/chaosmeta/SpaceController';
import { getUserList } from '@/services/chaosmeta/UserController';
import { SearchOutlined } from '@ant-design/icons';
import { history, useRequest } from '@umijs/max';
import { Button, Divider, Form, Input, Modal, Radio, Select } from 'antd';
import React, { useEffect, useState } from 'react';

interface IProps {
  open: boolean;
  setOpen: (val: boolean) => void;
}

const AddMemberModal: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;
  const [form] = Form.useForm();
  const [userList, setUserList] = useState<any[]>([]);
  const queryUserList = useRequest(getUserList, {
    manual: true,
    onSuccess: (res) => {
      if (res?.users) {
        setUserList(res.users);
      }
      console.log(res, 'eres====');
    },
  });

  const addUser = useRequest(spaceAddUser, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      console.log(res, ' res----');
    },
  });

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

  /**
   * 添加成员
   */
  const handleAddUser = () => {
    form.validateFields().then((values) => {
      console.log(values, 'values===');
      const users = values.users?.map((item) => {
        return {
          id: item,
          permission: values.permission,
        };
      });
      const params = {
        id: history.location.query.spaceId,
        users,
      };
      console.log(params, 'params');
      addUser.run(params);
    });
  };
  useEffect(() => {
    if (open) {
      queryUserList.run({ page: 2, page_size: 10 });
    }
  }, [open]);
  return (
    <Modal
      title="添加成员"
      open={open}
      onOk={handleAddUser}
      onCancel={() => {
        setOpen(false);
      }}
    >
      <Divider />
      <Form form={form} layout="vertical">
        <Form.Item
          label="用户名"
          name={'users'}
          rules={[{ required: true, message: '请选择' }]}
        >
          {/* <Users /> */}
          <Select mode="multiple">
            {userList?.map((item) => {
              return (
                <Select.Option key={item.id} value={item.id}>
                  {item.name}
                </Select.Option>
              );
            })}
          </Select>
        </Form.Item>
        <Form.Item
          label="用户权限"
          name={'permission'}
          initialValue={0}
          rules={[{ required: true, message: '请选择' }]}
        >
          <Radio.Group>
            <Radio value={0}>只读</Radio>
            <Radio value={1}>读写</Radio>
          </Radio.Group>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default React.memo(AddMemberModal);
