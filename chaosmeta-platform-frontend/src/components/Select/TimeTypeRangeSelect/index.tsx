import { Col, DatePicker, Form, Radio, Row, Select } from 'antd';
import React from 'react';

interface IProps {
  colSpan?: number;
}
/**
 * 时间类型下拉框及时间选择器
 * @returns
 */
const TimeTypeRangeSelect: React.FC<IProps> = (props) => {
  const { colSpan = 20 } = props;
  const { RangePicker } = DatePicker;
  const timeTypes = [
    {
      value: 'recentExperiments',
      label: '最近实验时间',
    },
    {
      value: 'recentlyEdited',
      label: '最近编辑时间',
    },
    {
      value: 'run',
      label: '即将运行时间',
    },
  ];

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
          <Select placeholder="请选择">
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
        <Radio.Group>
          {timeFastType?.map((item) => {
            return (
              <Radio.Button value={item.value} key={item.value}>
                {item.label}
              </Radio.Button>
            );
          })}
        </Radio.Group>
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
