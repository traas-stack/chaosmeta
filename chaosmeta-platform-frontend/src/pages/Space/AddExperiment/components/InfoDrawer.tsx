import ShowText from '@/components/ShowText';
import { triggerTypes } from '@/constants';
import { formatTime, getIntlLabel } from '@/utils/format';
import { renderScheduleType, renderTags } from '@/utils/renderItem';
import { history, useIntl } from '@umijs/max';
import { Button, DatePicker, Drawer, Form, Input, Radio, Space } from 'antd';
import dayjs from 'dayjs';
import React, { useEffect, useState } from 'react';
import { InfoEditDrawer } from '../style';
import TagSelect from './TagSelect';

interface IProps {
  open: boolean;
  setOpen: (val: boolean) => void;
  spacePermission?: number;
  handleConfirm: any;
  baseInfo: any;
}

/**
 * 编辑和查看信息的抽屉
 * @param props
 * @returns
 */
const InfoDrawer: React.FC<IProps> = (props) => {
  const { open, setOpen, spacePermission = 1, handleConfirm, baseInfo } = props;
  const [form] = Form.useForm();
  const [addTagList, setAddTagList] = useState<any>([]);
  const intl = useIntl();

  const handleClose = () => {
    setOpen(false);
  };

  /**
   * 只读用户信息渲染
   */
  const readOnlyRender = () => {
    return (
      <div>
        <Form.Item name={'name'} label="实验名称">
          <ShowText />
        </Form.Item>
        <Form.Item name={'description'} label="实验描述">
          <ShowText />
        </Form.Item>
        <Form.Item label="标签">{renderTags(baseInfo?.labels)}</Form.Item>
        <Form.Item label="触发方式">{renderScheduleType(baseInfo)}</Form.Item>
      </div>
    );
  };

  /**
   * 确认基本信息
   */
  const handleSubmit = () => {
    form.validateFields().then((values) => {
      if (values?.schedule_type === 'once') {
        values.schedule_rule = formatTime(values?.once_time);
      }
      handleConfirm({ ...values, labels: addTagList });
      setOpen(false);
    });
  };

  useEffect(() => {
    if (open) {
      // 区分周期性和单次定时的赋值
      const { schedule_rule, schedule_type, labels, description, name } =
        baseInfo;
      form.setFieldsValue({
        schedule_type,
        labels,
        description,
        name,
      });
      if (baseInfo?.schedule_type === 'once') {
        form.setFieldValue('once_time', dayjs(schedule_rule));
      }
      if (baseInfo?.schedule_type === 'cron') {
        form.setFieldValue('schedule_rule', schedule_rule);
      }
      setAddTagList(baseInfo?.labels || []);
    }
  }, [open]);

  return (
    <Drawer
      open={open}
      onClose={handleClose}
      title="基本信息"
      width={560}
      footer={
        <div>
          <Space>
            <Button onClick={handleClose}>
              {' '}
              {intl.formatMessage({ id: 'cancel' })}
            </Button>
            <Button
              type="primary"
              onClick={() => {
                handleSubmit();
              }}
            >
              {intl.formatMessage({ id: 'confirm' })}
            </Button>
          </Space>
        </div>
      }
    >
      <InfoEditDrawer>
        <Form form={form} layout="vertical">
          {spacePermission === 1 ? (
            <>
              <Form.Item
                name={'name'}
                label={intl.formatMessage({ id: 'experimentName' })}
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'inputPlaceholder' }),
                  },
                ]}
              >
                <Input
                  placeholder={intl.formatMessage({ id: 'inputPlaceholder' })}
                />
              </Form.Item>
              <Form.Item
                name={'description'}
                label={intl.formatMessage({ id: 'experimentDescription' })}
              >
                <Input.TextArea
                  rows={3}
                  placeholder={intl.formatMessage({ id: 'inputPlaceholder' })}
                />
              </Form.Item>
              <TagSelect
                spaceId={history?.location?.query?.spaceId as string}
                setAddTagList={setAddTagList}
                addTagList={addTagList}
              />
              <Form.Item
                name={'schedule_type'}
                label={intl.formatMessage({ id: 'triggerMode' })}
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'selectPlaceholder' }),
                  },
                ]}
                initialValue={'manual'}
              >
                <Radio.Group>
                  {triggerTypes?.map((item) => {
                    return (
                      <Radio value={item?.value} key={item?.value}>
                        {getIntlLabel(item)}
                      </Radio>
                    );
                  })}
                </Radio.Group>
              </Form.Item>

              <Form.Item
                noStyle
                shouldUpdate={(pre, cur) =>
                  pre?.schedule_type !== cur?.schedule_type
                }
              >
                {({ getFieldValue }) => {
                  const triggerType = getFieldValue('schedule_type');
                  if (triggerType === 'once') {
                    return (
                      <div className="trigger-type">
                        <Form.Item
                          name={'once_time'}
                          rules={[
                            {
                              required: true,
                              message: intl.formatMessage({
                                id: 'selectPlaceholder',
                              }),
                            },
                          ]}
                        >
                          <DatePicker format="YYYY-MM-DD HH:mm:ss" showTime />
                        </Form.Item>
                      </div>
                    );
                  }
                  if (triggerType === 'cron') {
                    return (
                      <div className="trigger-type">
                        <Form.Item
                          name={'schedule_rule'}
                          label={`Cron ${intl.formatMessage({
                            id: 'expression',
                          })}`}
                          rules={[
                            {
                              required: true,
                              message: `${intl.formatMessage({
                                id: 'inputPlaceholder',
                              })} Cron ${intl.formatMessage({
                                id: 'expression',
                              })}`,
                            },
                          ]}
                        >
                          <Input
                            placeholder={`${intl.formatMessage({
                              id: 'inputPlaceholder',
                            })} ${intl.formatMessage({ id: 'expression' })}`}
                          />
                        </Form.Item>
                      </div>
                    );
                  }
                  return null;
                }}
              </Form.Item>
            </>
          ) : (
            readOnlyRender()
          )}
        </Form>
      </InfoEditDrawer>
    </Drawer>
  );
};
export default React.memo(InfoDrawer);
