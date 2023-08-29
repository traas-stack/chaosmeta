import { login, register } from '@/services/chaosmeta/UserController';
import { history, useRequest } from '@umijs/max';
import { Button, Form, Input, message } from 'antd';
import CryptoJS from 'crypto-js';
import { useState } from 'react';
import { Container } from './style';

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
              {operateType === 'login' ? '登录' : '注册'}
            </div>
            <div className="tip">欢迎使用ChaosMeta</div>
            <Form.Item
              name={'name'}
              rules={[{ required: true, message: '请输入用户名' }]}
              help="用户名可以使用中英文，长度不超过64个字符"
            >
              <Input placeholder="用户名" maxLength={64} />
            </Form.Item>
            <Form.Item
              name={'password'}
              rules={[{ required: true, message: '请输入密码' }]}
              help={'密码8-16位中英文大小写及下划线等特殊字符'}
            >
              <Input.Password
                placeholder="密码"
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
                        return Promise.reject('请确认密码');
                      }
                      if (password !== value) {
                        return Promise.reject('两次密码不一致');
                      }
                      return Promise.resolve();
                    },
                  },
                ]}
              >
                <Input.Password
                  placeholder="确认密码"
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
              {operateType === 'login' ? '登录' : '注册'}
            </Button>
            <div>
              {operateType === 'login' ? (
                <div>
                  还没有账号？{' '}
                  <a
                    onClick={() => {
                      setOperateType('register');
                      form.resetFields();
                    }}
                  >
                    注册
                  </a>
                </div>
              ) : (
                <div
                  onClick={() => {
                    setOperateType('login');
                  }}
                >
                  已经有账号？ <a>登录</a>
                </div>
              )}
            </div>
          </div>
        </div>
      </Form>
    </Container>
  );
};
