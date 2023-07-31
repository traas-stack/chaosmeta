/**
 * 右上角用户信息区域
 */

import {
  getUserInfo,
  updatePassword,
} from '@/services/chaosmeta/UserController';
import cookie from '@/utils/cookie';
import { history, useModel, useRequest } from '@umijs/max';
import { Button, Dropdown, Form, Input, Modal, Space, message } from 'antd';
import CryptoJS from 'crypto-js';
import React, { useEffect, useState } from 'react';

const UserRightArea: React.FC<any> = () => {
  const [form] = Form.useForm();
  const [passwordOpen, setPasswordOpen] = useState<boolean>(false);
  const { userInfo, setUserInfo } = useModel('global');

  // 获取用户信息
  const queryUserInfo = useRequest(getUserInfo, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      console.log(res);
      setUserInfo({ ...userInfo, ...res?.data });
    },
  });
  // 下拉菜单选项
  const items = [
    {
      label: (
        <div
          onClick={() => {
            setPasswordOpen(true);
          }}
        >
          修改密码
        </div>
      ),
      key: 'updatePassword',
    },
    {
      label: (
        <div
          onClick={() => {
            cookie.clearToken('TOKEN');
            cookie.clearToken('REFRESH_TOKEN');
            localStorage.removeItem('userName');
            history.push('/login');
          }}
        >
          退出登录
        </div>
      ),
      key: 'logout',
    },
  ];

  /**
   * 修改密码接口
   */
  const handleUpdatePassword = useRequest(updatePassword, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success('密码修改成功，即将跳转到登录页面重新登录');
        setTimeout(() => {
          cookie.clearToken('TOKEN');
          cookie.clearToken('REFRESH_TOKEN');
          localStorage.removeItem('userName');
          history.push('/login');
        }, 1500);
      }
    },
  });

  /**
   * 提交修改
   */
  const submit = () => {
    form.validateFields().then((values) => {
      const { password } = values;
      console.log(values, 'values===');
      handleUpdatePassword.run({ password: CryptoJS.MD5(password).toString() });
    });
  };

  useEffect(() => {
    const name = localStorage.getItem('userName');
    if (name) {
      queryUserInfo.run({ name });
    }
  }, [localStorage.getItem('userName')]);

  return (
    <>
      <Dropdown menu={{ items }}>
        <Space style={{ cursor: 'pointer' }}>
          <img src={userInfo?.avatar} />
          <span>{userInfo?.name}</span>
        </Space>
      </Dropdown>
      {passwordOpen && (
        <Modal
          open={passwordOpen}
          closeIcon={false}
          onCancel={() => {
            setPasswordOpen(false);
          }}
          footer={
            <>
              <div>
                <Space>
                  <Button
                    onClick={() => {
                      setPasswordOpen(false);
                    }}
                  >
                    取消
                  </Button>
                  <Button
                    onClick={() => {
                      submit();
                    }}
                    type="primary"
                    loading={handleUpdatePassword.loading}
                  >
                    提交
                  </Button>
                </Space>
              </div>
            </>
          }
        >
          <Form form={form}>
            <Form.Item
              name={'password'}
              rules={[{ required: true, message: '请输入密码' }]}
            >
              <Input.Password placeholder="密码" />
            </Form.Item>
            <Form.Item
              name={'confirmPassword'}
              rules={[
                {
                  validator(rule, value) {
                    const password = form.getFieldValue('password');
                    if (!value) {
                      return Promise.reject('请确认密码');
                    }
                    if (password !== value) {
                      return Promise.reject('密码不正确');
                    }
                    return Promise.resolve();
                  },
                },
              ]}
            >
              <Input.Password placeholder="确认密码" />
            </Form.Item>
          </Form>
        </Modal>
      )}
    </>
  );
};

export default React.memo(UserRightArea);
