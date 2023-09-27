import ShowText from '@/components/ShowText';
import {
  editSpaceBasic,
  querySpaceDetail,
} from '@/services/chaosmeta/SpaceController';
import { formatTime } from '@/utils/format';
import { useParamChange } from '@/utils/useParamChange';
import { history, useIntl, useModel, useRequest } from '@umijs/max';
import { Button, Form, Input, Space, Spin } from 'antd';
import React, { useEffect, useState } from 'react';
import { BasicInfoContainer } from './style';

const BasicInfo: React.FC<any> = () => {
  const [form] = Form.useForm();
  const [saveDisabled, setSaveDisabled] = useState<boolean>(true);
  const [spaceInfo, setSpaceInfo] = useState<any>({});
  const spaceIdChange = useParamChange('spaceId');
  const { spacePermission } = useModel('global');
  const intl = useIntl();

  /**
   * 修改接口
   */
  const editInfo = useRequest(editSpaceBasic, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        setSaveDisabled(true);
      }
    },
  });

  /**
   * 获取空间信息
   */
  const getSpaceInfo = useRequest(querySpaceDetail, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      setSpaceInfo(res?.data?.namespace);
      form.setFieldsValue({
        ...res?.data?.namespace,
        create_time: formatTime(res?.data?.namespace?.create_time),
      });
    },
  });

  /**
   * 修改
   */
  const handleEdit = () => {
    form.validateFields().then((values) => {
      const params = {
        id: Number(history.location.query.spaceId),
        description: values.description,
        name: values.name,
      };
      editInfo.run(params);
    });
  };
  useEffect(() => {
    if (history.location.query.spaceId) {
      getSpaceInfo?.run({ id: history.location.query.spaceId as string });
    }
  }, [spaceIdChange]);
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
          {/* 读写权限时 */}
          {spacePermission === 1 ? (
            <div>
              <Form.Item
                name={'name'}
                label={intl.formatMessage({ id: 'spaceName' })}
                rules={[
                  {
                    required: true,
                    message: `${intl.formatMessage({
                      id: 'inputPlaceholder',
                    })} ${intl.formatMessage({ id: 'spaceName' })}`,
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
              >
                <Input.TextArea
                  placeholder={`${intl.formatMessage({
                    id: 'inputPlaceholder',
                  })} ${intl.formatMessage({ id: 'spaceDescription' })}`}
                  style={{ resize: 'none' }}
                  rows={4}
                  maxLength={200}
                  showCount
                />
              </Form.Item>
              <Form.Item
                name={'create_time'}
                label={intl.formatMessage({ id: 'createTime' })}
              >
                <ShowText />
              </Form.Item>
              <Form.Item
                name={'count'}
                label={intl.formatMessage({ id: 'memberNumber' })}
              >
                <ShowText />
              </Form.Item>
              <Space>
                <Button
                  type="primary"
                  disabled={saveDisabled}
                  onClick={handleEdit}
                  loading={editInfo.loading}
                >
                  {intl.formatMessage({ id: 'save' })}
                </Button>
                <Button
                  onClick={() => {
                    form.setFieldsValue(spaceInfo);
                    setSaveDisabled(true);
                  }}
                >
                  {intl.formatMessage({ id: 'cancel' })}
                </Button>
              </Space>
            </div>
          ) : (
            // 只读权限时
            <div>
              <Form.Item
                name={'name'}
                label={intl.formatMessage({ id: 'spaceName' })}
              >
                <ShowText />
              </Form.Item>
              <Form.Item
                name={'description'}
                label={intl.formatMessage({ id: 'spaceDescription' })}
              >
                <ShowText />
              </Form.Item>
              <Form.Item
                name={'create_time'}
                label={intl.formatMessage({ id: 'createTime' })}
              >
                <ShowText />
              </Form.Item>
              <Form.Item
                name={'count'}
                label={intl.formatMessage({ id: 'memberNumber' })}
              >
                <ShowText />
              </Form.Item>
            </div>
          )}
        </Form>
      </BasicInfoContainer>
    </Spin>
  );
};

export default BasicInfo;
