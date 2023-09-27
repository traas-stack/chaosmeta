import {
  createSpace,
  queryClassSpaceList,
} from '@/services/chaosmeta/SpaceController';
import { history, useIntl, useModel, useRequest } from '@umijs/max';
import { Button, Drawer, Form, Input, Space, message } from 'antd';
import React from 'react';

interface IProps {
  open: boolean;
  setOpen: (open: boolean) => void;
}

const AddSpaceDrawer: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;
  const [form] = Form.useForm();
  const { setCurSpace, setSpaceList } = useModel('global');
  const intl = useIntl();

  const handleCancel = () => {
    setOpen(false);
  };

  /**
   * 获取空间列表 -- 当前用户有查看权限的空间只读和读写
   */
  const getSpaceList = useRequest(queryClassSpaceList, {
    manual: true,
    formatResult: (res) => res,
    debounceInterval: 300,
    onSuccess: (res) => {
      if (res?.code === 200) {
        const namespaceList = res.data?.namespaces?.map(
          (item: { namespaceInfo: any }) => {
            // side侧边菜单的展开/收起会影响这里，暂时用icon代替，todo
            return {
              icon: item?.namespaceInfo?.name,
              key: item?.namespaceInfo?.id,
              id: item?.namespaceInfo?.id?.toString(),
              name: item?.namespaceInfo?.name,
            };
          },
        );
        // 更新空间列表
        setSpaceList(namespaceList);
      }
    },
  });

  /**
   * 更新地址栏空间id，并保存
   * @param id
   */
  const handleUpdateSpaceId = (id: any) => {
    if (id) {
      const name = form.getFieldValue('name');
      getSpaceList?.run({ page: 1, page_size: 10, namespaceClass: 'relevant' });
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
        message.success(intl.formatMessage({ id: 'createText' }));
        // 更新空间信息
        handleUpdateSpaceId(res?.data?.id);
        setOpen(false);
        history.push({
          pathname: '/space/setting',
          query: {
            spaceId: res?.data?.id,
          },
        });
      }
    },
  });

  /**
   * 创建空间
   */
  const handleCreate = () => {
    form.validateFields().then((values) => {
      create?.run(values);
    });
  };

  return (
    <Drawer
      title={intl.formatMessage({ id: 'createSpace' })}
      open={open}
      onClose={handleCancel}
      width={480}
      footer={
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button onClick={handleCancel}>
              {intl.formatMessage({ id: 'cancel' })}
            </Button>
            <Button
              type="primary"
              onClick={handleCreate}
              loading={create.loading}
            >
              {intl.formatMessage({ id: 'createSpace.confirm' })}
            </Button>
          </Space>
        </div>
      }
    >
      <Form layout="vertical" form={form}>
        <Form.Item
          name={'name'}
          label={intl.formatMessage({ id: 'spaceName' })}
          rules={[
            {
              required: true,
              message: intl.formatMessage({ id: 'inputPlaceholder' }),
            },
          ]}
          help={intl.formatMessage({ id: 'spaceDescriptionTip' })}
        >
          <Input
            placeholder={`${intl.formatMessage({
              id: 'inputPlaceholder',
            })} ${intl.formatMessage({ id: 'spaceName' })}`}
            maxLength={64}
          />
        </Form.Item>
        <Form.Item
          name={'description'}
          label={intl.formatMessage({ id: 'spaceDescription' })}
          style={{ marginTop: '36px' }}
        >
          <Input.TextArea
            placeholder={`${intl.formatMessage({
              id: 'inputPlaceholder',
            })} ${intl.formatMessage({ id: 'spaceDescription' })}`}
            rows={5}
          />
        </Form.Item>
      </Form>
    </Drawer>
  );
};

export default React.memo(AddSpaceDrawer);
