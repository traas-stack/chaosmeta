import { formatFormName, getIntlText } from '@/utils/format';
import { useIntl } from '@umijs/max';
import { Form, Input, Radio, Select } from 'antd';
import ShowText from '../ShowText';
import UnitInput from './UnitInput';
import { DividerLine } from './style';

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
  form?: any;
}

/**
 * 动态表单渲染
 * @param props
 * @returns
 */
const DynamicForm = (props: Props) => {
  const { fieldList, parentName, readonly, form } = props;
  const intl = useIntl();

  /**
   * 解析动态表单配置内容，根据不同条件渲染为不同表单元素
   * @param field
   * @returns
   */
  const renderItem = (field: any) => {
    const { valueRule, valueType } = field;
    // 无论valueType类型是什么，后端统一接收string
    // if (valueType === 'int') {
    //   return (
    //     <InputNumber
    //       placeholder={getIntlText(field, 'descriptionCn', 'description')}
    //       style={{ width: '100%' }}
    //     />
    //   );
    // }
    if (valueType === 'bool') {
      return (
        <Radio.Group>
          <Radio value={'true'}>{intl.formatMessage({ id: 'yes' })}</Radio>
          <Radio value={'false'}>{intl.formatMessage({ id: 'no' })}</Radio>
        </Radio.Group>
      );
    }
    if (valueType === 'stringlist') {
      let options: string[] = [];
      if (valueRule) {
        options = valueRule?.split(',');
      }
      return (
        <Select mode={'tags'} placeholder={getIntlText(field, 'keyCn', 'key')}>
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
    if (
      valueRule &&
      !valueRule.includes('-') &&
      !valueRule.includes('>') &&
      !valueRule.includes('<') &&
      !valueRule.includes('=')
    ) {
      let options: string[] = [];
      if (valueRule) {
        options = valueRule?.split(',');
      }
      return (
        <Select placeholder={getIntlText(field, 'keyCn', 'key')}>
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

    return <Input placeholder={getIntlText(field, 'keyCn', 'key')} />;
  };

  const initValue = (defaultValue: any) => {
    if (defaultValue || defaultValue === 0) {
      return defaultValue;
    }
    return undefined;
  };

  /**
   * 判断规则函数，根据传入的参数进行判断并返回相应的结果
   * @param {Object} item - 包含规则信息的对象
   * @param {*} value - 待判断的值
   * @returns {Promise<string|undefined>} - 返回一个 Promise 对象，包含判断结果或 undefined
   */
  const rule = (item: any, value: number) => {
    // 获取规则信息中的 valueType、valueRule 和 keyCn 属性
    const { valueType, valueRule } = item;
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
          return Promise.reject(
            `${getIntlText(item, 'keyCn', 'key')} ${intl.formatMessage({
              id: 'ruleText',
            })} ${valueRule}`,
          );
        }
      }
      if (valueRule?.includes('>=')) {
        const valueRuleList = valueRule.split('>=');
        if (value < Number(valueRuleList[1])) {
          // 不满足大于规定范围条件，返回错误信息
          return Promise.reject(
            `${getIntlText(item, 'keyCn', 'key')} ${intl.formatMessage({
              id: 'ruleText',
            })} ${valueRule}`,
          );
        }
      }
      if (valueRule?.includes('>')) {
        const valueRuleList = valueRule.split('>');
        if (value <= Number(valueRuleList[1])) {
          // 不满足大于规定范围条件，返回错误信息
          return Promise.reject(
            `${getIntlText(item, 'keyCn', 'key')} ${intl.formatMessage({
              id: 'ruleText',
            })} ${valueRule}`,
          );
        }
      }
    }
    // 其他情况都返回通过
    return Promise.resolve();
  };

  return (
    <>
      {fieldList?.map((item: Field, index: number) => {
        const { id, required, defaultValue } = item;
        if (readonly) {
          return (
            <>
              {index === 0 && item?.execType?.includes('common') && (
                <div className="subtitle range">
                  {intl.formatMessage({ id: 'generalParameters' })}
                </div>
              )}
              {!item?.execType?.includes('common') &&
                fieldList[index - 1]?.execType?.includes('common') && (
                  <DividerLine />
                )}
              <Form.Item
                name={formatFormName(item, parentName)}
                label={getIntlText(item, 'keyCn', 'key')}
                key={id}
                initialValue={initValue(defaultValue)}
              >
                <ShowText ellipsis />
              </Form.Item>
            </>
          );
        }
        return (
          <>
            {index === 0 && item?.execType?.includes('common') && (
              <div className="range">
                {intl.formatMessage({ id: 'generalParameters' })}
              </div>
            )}
            {!item?.execType?.includes('common') &&
              fieldList[index - 1]?.execType?.includes('common') && (
                <DividerLine />
              )}
            <Form.Item
              tooltip={getIntlText(item, 'descriptionCn', 'description')}
              name={formatFormName(item, parentName)}
              label={getIntlText(item, 'keyCn', 'key')}
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
              {/* 带有单位的特殊处理 */}
              {item.unit ? (
                <UnitInput field={item} form={form} parentName={parentName} />
              ) : (
                renderItem(item)
              )}
            </Form.Item>
          </>
        );
      })}
    </>
  );
};

export default DynamicForm;
