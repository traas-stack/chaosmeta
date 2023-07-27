import { login, register } from '@/services/chaosmeta/UserController';
import { history, useRequest } from '@umijs/max';
import { Button, Form, Input, message } from 'antd';
import CryptoJS from 'crypto-js';
import { useState } from 'react';
import cookie from 'react-cookies';
import { Container, OperateArea } from './style';

export default () => {
  const [form] = Form.useForm();
  const [operateType, setOperateType] = useState<'login' | 'register'>('login');

  /**
   * 注册
   */
  const handleRegister = useRequest(register, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res.code === 200) {
        message.success('注册成功，请登录');
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
        // cookie.save('TOKEN', res.data?.refreshToken, {
        //   domain: document.domain,
        // });
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
      <div className="login">
        <Form form={form}>
          <Form.Item
            name={'name'}
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="用户名" />
          </Form.Item>
          <Form.Item
            name={'password'}
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="密码" />
          </Form.Item>
          {operateType === 'register' && (
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
          )}
        </Form>
        <OperateArea operatetype={operateType}>
          {operateType === 'login' ? (
            <>
              <Button
                onClick={() => {
                  setOperateType('register');
                }}
              >
                注册
              </Button>
              <Button
                type="primary"
                loading={handleLogin?.loading}
                onClick={submit}
              >
                登录
              </Button>
            </>
          ) : (
            <>
              <div>
                已有账号？
                <a
                  onClick={() => {
                    setOperateType('login');
                  }}
                >
                  登录
                </a>
              </div>
              <Button
                type="primary"
                loading={handleRegister?.loading}
                onClick={submit}
              >
                注册
              </Button>
            </>
          )}
        </OperateArea>
      </div>
    </Container>
  );
};
