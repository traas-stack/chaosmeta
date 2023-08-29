import { createSpace } from '@/services/chaosmeta/SpaceController';
import { history, useModel, useRequest } from '@umijs/max';
import { Button, Drawer, Form, Input, Space, message } from 'antd';
import React from 'react';

interface IProps {
  open: boolean;
  setOpen: (open: boolean) => void;
}

const AddSpaceDrawer: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;
  const [form] = Form.useForm();
  const { setCurSpace } = useModel('global');

  const handleCancel = () => {
    setOpen(false);
  };

  /**
   * 更新地址栏空间id，并保存
   * @param id
   */
  const handleUpdateSpaceId = (id: any) => {
    if (id) {
      const name = form.getFieldValue('name');
      history.push({
        pathname: history.location.pathname,
        query: {
          ...history.location.query,
          spaceId: id,
        },
      });
      setCurSpace([id]);
      sessionStorage.setItem('spaceId', id);
      sessionStorage.setItem('spaceName', name);
    }
  };

  /**
   * 创建空间接口
   */
  const create = useRequest(createSpace, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: async (res) => {
      if (res.code === 200) {
        message.success('创建成功');
        // 更新空间信息
        handleUpdateSpaceId(res?.data?.id);
        setOpen(false);
        history.push('/space/setting');
      }
    },
  });

  /**
   * 创建空间
   */
  const handleCreate = () => {
    // setCurSpace()
    // handleUpdateSpaceId('20', 'ceshi')
    // return
    form.validateFields().then((values) => {
      create?.run(values);
    });
  };

  return (
    <Drawer
      title="新建空间"
      open={open}
      onClose={handleCancel}
      width={480}
      footer={
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button onClick={handleCancel}>取消</Button>
            <Button
              type="primary"
              onClick={handleCreate}
              loading={create.loading}
            >
              创建完成并去配置
            </Button>
          </Space>
        </div>
      }
    >
      <Form layout="vertical" form={form}>
        <Form.Item
          name={'name'}
          label="空间名称"
          rules={[{ required: true, message: '请输入' }]}
          help="请尽量保持空间名称的简洁，不超过64个字符"
        >
          <Input placeholder="请输入空间名称" maxLength={64} />
        </Form.Item>
        <Form.Item
          name={'description'}
          label="空间描述"
          style={{ marginTop: '36px' }}
        >
          <Input.TextArea placeholder="请输入空间描述" rows={5} />
        </Form.Item>
      </Form>
    </Drawer>
  );
};

export default React.memo(AddSpaceDrawer);
