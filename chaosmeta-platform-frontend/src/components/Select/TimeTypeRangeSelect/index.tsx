import { Col, DatePicker, Form, Radio, Select } from 'antd';
import moment from 'moment';
import React from 'react';

interface IProps {
  form: any;
  timeTypes: any[];
}
/**
 * 时间类型下拉框及时间选择器
 * @returns
 */
const TimeTypeRangeSelect: React.FC<IProps> = (props) => {
  const { form, timeTypes } = props;
  const { RangePicker } = DatePicker;

  const timeFastType = [
    {
      value: 'all',
      label: '全部',
    },
    {
      value: '7day',
      label: '近7天',
    },
    {
      value: '30day',
      label: '近30天',
    },
  ];
  return (
    <>
      <Col span={6}>
        <Form.Item
          name={'timeType'}
          label="时间类型"
          labelCol={{ span: 8 }}
          // style={{minWidth: '260px'}}
        >
          <Select placeholder="请选择" allowClear>
            {timeTypes?.map((item) => {
              return (
                <Select.Option key={item.value} value={item.value}>
                  {item.label}
                </Select.Option>
              );
            })}
          </Select>
        </Form.Item>
      </Col>
      <Col>
        <Form.Item name={'timePresets'}>
          <Radio.Group
            onChange={(event: any) => {
              const value = event.target?.value;
              if (value === 'all') {
                form.setFieldValue('time', undefined);
              }
              if (value === '7day') {
                const timeRange = [
                  moment().subtract(6, 'd').startOf('day'),
                  moment().endOf('day'),
                ];
                form.setFieldValue('time', timeRange);
              }
              if (value === '30day') {
                const timeRange = [
                  moment().subtract(30, 'd').startOf('day'),
                  moment().endOf('day'),
                ];
                form.setFieldValue('time', timeRange);
              }
            }}
          >
            {timeFastType?.map((item) => {
              return (
                <Radio.Button value={item.value} key={item.value}>
                  {item.label}
                </Radio.Button>
              );
            })}
          </Radio.Group>
        </Form.Item>
      </Col>
      <Col>
        <Form.Item name={'time'}>
          <RangePicker showTime />
        </Form.Item>
      </Col>
    </>
  );
};

export default React.memo(TimeTypeRangeSelect);
