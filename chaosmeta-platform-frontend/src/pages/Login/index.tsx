import { login, register } from '@/services/chaosmeta/UserController';
import { history, useIntl, useRequest } from '@umijs/max';
import { Button, Form, Input, message } from 'antd';
import CryptoJS from 'crypto-js';
import { useState } from 'react';
import { Container } from './style';

export default () => {
  const [form] = Form.useForm();
  const [operateType, setOperateType] = useState<'login' | 'register'>('login');
  const intl = useIntl();

  /**
   * 注册
   */
  const handleRegister = useRequest(register, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res.code === 200) {
        message.success(intl.formatMessage({ id: 'reister.success' }));
        form.resetFields();
        setOperateType('login');
      }
    },
  });

  /**
   * 登录
   */
  const handleLogin = useRequest(login, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res: any, params: any) => {
      if (res.code === 200) {
        localStorage.setItem('userName', params[0]?.name);
        history.push('/space/overview');
        form.resetFields();
      }
    },
  });

  /**
   * 提交
   */
  const submit = () => {
    form.validateFields().then((values) => {
      const { name, password } = values;
      const params = {
        name,
        password: CryptoJS.MD5(password).toString(),
      };
      if (operateType === 'login') {
        handleLogin?.run(params);
      } else {
        handleRegister?.run(params);
      }
    });
  };

  return (
    <Container>
      <Form form={form}>
        {/* 占位用 */}
        <div className="seize"></div>
        {/* 卡片部分 */}
        <div className="card">
          <div className="bg"></div>
          <div className="content">
            <div className="img">
              <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*lMXkRKmd8WcAAAAAAAAAAAAADmKmAQ/original" />
            </div>
            <div className="title">
              {operateType === 'login'
                ? intl.formatMessage({ id: 'login' })
                : intl.formatMessage({ id: 'register' })}
            </div>
            <div className="tip">
              {intl.formatMessage({ id: 'welcome' })} ChaosMeta
            </div>
            <Form.Item
              name={'name'}
              rules={[
                {
                  required: true,
                  message: intl.formatMessage({ id: 'usernamePlaceholder' }),
                },
              ]}
              help={intl.formatMessage({ id: 'username.rule' })}
            >
              <Input
                placeholder={intl.formatMessage({ id: 'username' })}
                maxLength={64}
              />
            </Form.Item>
            <Form.Item
              name={'password'}
              rules={[
                {
                  required: true,
                  message: intl.formatMessage({ id: 'password.placeholder' }),
                },
              ]}
              help={intl.formatMessage({ id: 'password.rule' })}
            >
              <Input.Password
                placeholder={intl.formatMessage({ id: 'password' })}
                minLength={8}
                maxLength={16}
                onPressEnter={() => {
                  submit();
                }}
              />
            </Form.Item>
            {operateType === 'register' && (
              <Form.Item
                name={'confirmPassword'}
                rules={[
                  {
                    validator(rule, value) {
                      const password = form.getFieldValue('password');
                      if (!value) {
                        return Promise.reject(
                          intl.formatMessage({ id: 'password.confirm' }),
                        );
                      }
                      if (password !== value) {
                        return Promise.reject(
                          intl.formatMessage({ id: 'password.inconsistent' }),
                        );
                      }
                      return Promise.resolve();
                    },
                  },
                ]}
              >
                <Input.Password
                  placeholder={intl.formatMessage({
                    id: 'password.confirm.again',
                  })}
                  minLength={8}
                  maxLength={16}
                />
              </Form.Item>
            )}

            <Button
              type="primary"
              loading={handleLogin?.loading || handleRegister?.loading}
              onClick={submit}
            >
              {operateType === 'login'
                ? intl.formatMessage({ id: 'login' })
                : intl.formatMessage({ id: 'register' })}
            </Button>
            <div>
              {operateType === 'login' ? (
                <div>
                  {intl.formatMessage({ id: 'notAccount' })}{' '}
                  <a
                    onClick={() => {
                      setOperateType('register');
                      form.resetFields();
                    }}
                  >
                    {intl.formatMessage({ id: 'register' })}
                  </a>
                </div>
              ) : (
                <div
                  onClick={() => {
                    setOperateType('login');
                  }}
                >
                  {intl.formatMessage({ id: 'haveAccount' })}{' '}
                  <a>{intl.formatMessage({ id: 'login' })}</a>
                </div>
              )}
            </div>
          </div>
        </div>
      </Form>
    </Container>
  );
};
