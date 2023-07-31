import ShowText from '@/components/ShowText';
import {
  editSpaceBasic,
  querySpaceDetail,
} from '@/services/chaosmeta/SpaceController';
import { useRequest } from '@umijs/max';
import { Button, Form, Input, Space, Spin } from 'antd';
import React, { useEffect, useState } from 'react';
import { BasicInfoContainer } from './style';

interface IProps {}

const BasicInfo: React.FC<IProps> = () => {
  const [form] = Form.useForm();
  const [saveDisabled, setSaveDisabled] = useState<boolean>(true);
  const [spaceInfo, setSpaceInfo] = useState<any>({});

  /**
   * 修改接口
   */
  const editInfo = useRequest(editSpaceBasic, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      console.log(res, 'res');
      setSaveDisabled(true);
    },
  });

  /**
   * 获取空间信息
   */
  const getSpaceInfo = useRequest(querySpaceDetail, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      console.log(res, 'res===');
      setSpaceInfo(res?.data?.namespace);
      form.setFieldsValue(res?.data?.namespace);
    },
  });

  /**
   * 修改
   */
  const handleEdit = () => {
    form.validateFields().then((values) => {
      console.log(values, 'values----');
      editInfo.run(values);
    });
  };

  useEffect(() => {
    form.setFieldsValue({
      createTime: '2023-07-12 14:00:00',
      count: '33',
    });
    getSpaceInfo?.run({ id: 1 });
  }, []);
  return (
    <Spin spinning={getSpaceInfo.loading}>
      <BasicInfoContainer>
        <Form
          form={form}
          layout="vertical"
          onValuesChange={() => {
            setSaveDisabled(false);
          }}
        >
          <Form.Item
            name={'name'}
            label="空间名称"
            rules={[{ required: true, message: '请输入空间名称' }]}
            help="请尽量保持项目名称的简洁，不超过 64 个字符"
          >
            <Input placeholder="请输入空间名称" maxLength={64} />
          </Form.Item>
          <Form.Item name={'description'} label="空间描述">
            <Input.TextArea
              placeholder="请输入空间描述"
              style={{ resize: 'none' }}
              rows={4}
              maxLength={200}
              showCount
            />
          </Form.Item>
          <Form.Item name={'create_time'} label="创建时间">
            <ShowText />
          </Form.Item>
          <Form.Item name={'count'} label="成员数量">
            <ShowText />
          </Form.Item>
          <Space>
            <Button
              type="primary"
              disabled={saveDisabled}
              onClick={handleEdit}
              loading={editInfo.loading}
            >
              保存
            </Button>
            <Button
              onClick={() => {
                form.setFieldsValue(spaceInfo);
                setSaveDisabled(true);
              }}
            >
              取消
            </Button>
          </Space>
        </Form>
      </BasicInfoContainer>
    </Spin>
  );
};

export default React.memo(BasicInfo);