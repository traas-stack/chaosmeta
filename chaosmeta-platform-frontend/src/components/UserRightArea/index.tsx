/**
 * 右上角用户信息区域
 */

import {
  getUserInfo,
  updatePassword,
} from '@/services/chaosmeta/UserController';
import cookie from '@/utils/cookie';
import { SelectLang, history, useIntl, useModel, useRequest } from '@umijs/max';
import {
  Button,
  Divider,
  Dropdown,
  Form,
  Input,
  Modal,
  Space,
  message,
} from 'antd';
import CryptoJS from 'crypto-js';
import React, { useEffect, useState } from 'react';

const UserRightArea: React.FC<any> = () => {
  const [form] = Form.useForm();
  const [passwordOpen, setPasswordOpen] = useState<boolean>(false);
  const { userInfo, setUserInfo } = useModel('global');
  const intl = useIntl();
  // 获取用户信息
  const queryUserInfo = useRequest(getUserInfo, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
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
          {intl.formatMessage({ id: 'updatePassword' })}
        </div>
      ),
      key: 'updatePassword',
    },
    {
      label: (
        <div
          style={{ color: '#FF4D4F' }}
          onClick={() => {
            cookie.clearToken('TOKEN');
            cookie.clearToken('REFRESH_TOKEN');
            localStorage.removeItem('userName');
            sessionStorage.clear();
            history.push('/login');
          }}
        >
          {intl.formatMessage({ id: 'signOut' })}
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
        message.success(intl.formatMessage({ id: 'password.success' }));

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
      const { password, oldPassword } = values;
      const params = {
        password: CryptoJS.MD5(password).toString(),
        oldPassword: CryptoJS.MD5(oldPassword).toString(),
      };
      handleUpdatePassword.run(params);
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
      {/* 语言切换 */}
      <SelectLang />
      <Dropdown menu={{ items }}>
        <Space style={{ cursor: 'pointer' }}>
          <img src={userInfo?.avatar} />
          <span>{userInfo?.name}</span>
        </Space>
      </Dropdown>
      {passwordOpen && (
        <Modal
          title={intl.formatMessage({ id: 'updatePassword' })}
          open={passwordOpen}
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
                    {intl.formatMessage({ id: 'cancel' })}
                  </Button>
                  <Button
                    onClick={() => {
                      submit();
                    }}
                    type="primary"
                    loading={handleUpdatePassword.loading}
                  >
                    {intl.formatMessage({ id: 'submit' })}
                  </Button>
                </Space>
              </div>
            </>
          }
        >
          <Divider />
          <Form form={form} layout="vertical">
            <Form.Item
              name={'oldPassword'}
              rules={[
                {
                  required: true,
                  message: intl.formatMessage({
                    id: 'password.old.placeholder',
                  }),
                },
              ]}
              label={intl.formatMessage({ id: 'password.old' })}
            >
              <Input.Password
                placeholder={intl.formatMessage({
                  id: 'password.old.placeholder',
                })}
              />
            </Form.Item>
            <Form.Item
              name={'password'}
              rules={[
                {
                  required: true,
                  message: intl.formatMessage({
                    id: 'password.new.placeholder',
                  }),
                },
              ]}
              help={intl.formatMessage({ id: 'password.rule' })}
              label={intl.formatMessage({ id: 'password.new' })}
            >
              <Input.Password
                placeholder={intl.formatMessage({
                  id: 'password.new.placeholder',
                })}
                minLength={8}
                maxLength={16}
              />
            </Form.Item>
            <Form.Item
              name={'confirmPassword'}
              rules={[
                {
                  required: true,
                  validator(rule, value) {
                    const password = form.getFieldValue('password');
                    if (!value) {
                      return Promise.reject(
                        intl.formatMessage({ id: 'password.confirm' }),
                      );
                    }
                    if (password !== value) {
                      return Promise.reject(
                        intl.formatMessage({ id: 'password.error' }),
                      );
                    }
                    return Promise.resolve();
                  },
                },
              ]}
              label={intl.formatMessage({ id: 'password.new.confirm' })}
            >
              <Input.Password
                placeholder={intl.formatMessage({
                  id: 'password.new.placeholder.again',
                })}
              />
            </Form.Item>
          </Form>
        </Modal>
      )}
    </>
  );
};

export default React.memo(UserRightArea);
