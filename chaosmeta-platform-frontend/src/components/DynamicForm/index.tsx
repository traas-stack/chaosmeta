import { Form, Input, InputNumber, Radio, Select } from 'antd';
import ShowText from '../ShowText';

interface Field {
  id: number;
  execType: string;
  injectId: number;
  key: string;
  keyCn: string;
  valueType: string;
  valueRule: string;
  description: string;
  descriptionCn: string;
  unit: string;
  unitCn: string;
  defaultValue: string;
  required: boolean;
  create_time: string;
  update_time: string;
}
interface Props {
  fieldList: Field[];
  parentName?: string;
  readonly?: boolean;
}

/**
 * 动态表单渲染
 * @param props
 * @returns
 */
const DynamicForm = (props: Props) => {
  const { fieldList, parentName, readonly } = props;

  /**
   * 解析动态表单配置内容，根据不同条件渲染为不同表单元素
   * @param field
   * @returns
   */
  const renderItem = (field: Field) => {
    const { valueRule, valueType, keyCn } = field;
    if (valueType === 'int') {
      return <InputNumber placeholder={keyCn} style={{ width: '100%' }} />;
    }
    if (valueType === 'bool') {
      return (
        <Radio.Group>
          <Radio value={true}>是</Radio>
          <Radio value={false}>否</Radio>
        </Radio.Group>
      );
    }
    if (valueType === 'stringlist') {
      let options: string[] = [];
      if (valueRule) {
        options = valueRule?.split(',');
      }
      return (
        <Select mode={'tags'} placeholder={keyCn}>
          {options.map((option) => {
            return (
              <Select.Option key={option} value={option}>
                {option}
              </Select.Option>
            );
          })}
        </Select>
      );
    }
    if (valueRule) {
      let options: string[] = [];
      if (valueRule) {
        options = valueRule?.split(',');
      }
      return (
        <Select placeholder={keyCn}>
          {options.map((option) => {
            return (
              <Select.Option key={option} value={option}>
                {option}
              </Select.Option>
            );
          })}
        </Select>
      );
    }
    return <Input placeholder={keyCn} />;
  };

  const initValue = (defaultValue: any) => {
    if (defaultValue || defaultValue === 0) {
      return defaultValue;
    }
    return undefined;
  };

  /**
   * 如果存在父级名字，则将其与当前名字合并成一个字符串数组返回；否则直接返回当前名字。
   * @param name
   * @returns
   */
  const formatName = (name: number) => {
    if (parentName) {
      return [parentName, name.toString()];
    }
    return name;
  };

  /**
   * 判断规则函数，根据传入的参数进行判断并返回相应的结果
   * @param {Object} item - 包含规则信息的对象
   * @param {*} value - 待判断的值
   * @returns {Promise<string|undefined>} - 返回一个 Promise 对象，包含判断结果或 undefined
   */
  const rule = (item: any, value: number) => {
    // 获取规则信息中的 valueType、valueRule 和 keyCn 属性
    const { valueType, valueRule, keyCn } = item;
    if (!value && value !== 0) {
      // 其他情况都返回通过, 为空时让form自动去判断
      return Promise.resolve();
    }
    // 如果 value 的类型为 int
    if (valueType === 'int') {
      // 根据 valueRule 中是否包含 '-' 分别处理小于和大于情况
      if (valueRule?.includes('-')) {
        const valueRuleList = valueRule.split('-');
        if (
          value < Number(valueRuleList[0]) ||
          value > Number(valueRuleList[1])
        ) {
          // 不满足大于规定范围条件，返回错误信息
          return Promise.reject(`${keyCn}的取值为 ${valueRule}`);
        }
      }
      if (valueRule?.includes('>=')) {
        const valueRuleList = valueRule.split('>=');
        if (value < Number(valueRuleList[1])) {
          // 不满足大于规定范围条件，返回错误信息
          return Promise.reject(`${keyCn}的取值为 ${valueRule}`);
        }
      }
      if (valueRule?.includes('>')) {
        const valueRuleList = valueRule.split('>');
        if (value <= Number(valueRuleList[1])) {
          // 不满足大于规定范围条件，返回错误信息
          return Promise.reject(`${keyCn}的取值为 ${valueRule}`);
        }
      }
    }
    // 其他情况都返回通过
    return Promise.resolve();
  };

  return (
    <>
      {fieldList?.map((item: Field) => {
        const { keyCn, id, required, defaultValue } = item;
        if (readonly) {
          return (
            <Form.Item
              name={formatName(id)}
              label={keyCn}
              key={id}
              // rules={[{ required, message: keyCn }]}
              initialValue={initValue(defaultValue)}
            >
              <ShowText />
            </Form.Item>
          );
        }
        return (
          <Form.Item
            name={formatName(id)}
            label={keyCn}
            key={id}
            rules={[
              { required },
              {
                validator: (_, value) => {
                  return rule(item, value);
                },
              },
            ]}
            initialValue={initValue(defaultValue)}
          >
            {renderItem(item)}
          </Form.Item>
        );
      })}
    </>
  );
};

export default DynamicForm;
