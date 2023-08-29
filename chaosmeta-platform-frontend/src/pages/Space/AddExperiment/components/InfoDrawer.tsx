import ShowText from '@/components/ShowText';
import { triggerTypes } from '@/constants';
import { formatTime } from '@/utils/format';
import { history } from '@umijs/max';
import { Button, DatePicker, Drawer, Form, Input, Radio, Space } from 'antd';
import moment from 'moment';
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
        <Form.Item name={'desc'} label="实验描述">
          <ShowText />
        </Form.Item>
        <Form.Item name={'labels'} label="标签">
          <ShowText />
        </Form.Item>
        <Form.Item name={'triggerType'} label="触发方式">
          <ShowText />
        </Form.Item>
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
        form.setFieldValue('once_time', moment(schedule_rule));
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
            <Button onClick={handleClose}>取消</Button>
            <Button
              type="primary"
              onClick={() => {
                handleSubmit();
              }}
            >
              确定
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
                label="实验名称"
                rules={[{ required: true, message: '请输入' }]}
              >
                <Input placeholder="请输入" />
              </Form.Item>
              <Form.Item name={'description'} label="实验描述">
                <Input.TextArea rows={3} placeholder="请输入" />
              </Form.Item>
              <TagSelect
                spaceId={history?.location?.query?.spaceId as string}
                setAddTagList={setAddTagList}
                addTagList={addTagList}
              />
              <Form.Item
                name={'schedule_type'}
                label="触发方式"
                rules={[{ required: true, message: '请选择' }]}
              >
                <Radio.Group>
                  {triggerTypes?.map((item) => {
                    return (
                      <Radio value={item?.value} key={item?.value}>
                        {item?.label}
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
                          rules={[{ required: true, message: '请选择' }]}
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
                          label="Cron表达式"
                          rules={[
                            { required: true, message: '请输入Cron表达式' },
                          ]}
                        >
                          <Input placeholder="请输入表达式" />
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
