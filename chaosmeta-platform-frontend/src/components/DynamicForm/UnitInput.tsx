/**
 * 带有单位的表单特殊渲染
 */

import { getIntlText } from '@/utils/format';
import { Input, Select } from 'antd';
import { useEffect, useState } from 'react';
interface Field {
  defaultValue: string;
  unit: string;
  id: number;
}
interface IProps {
  field: Field;
  form: any;
  parentName?: string;
  onChange?: any;
  value?: string;
}

const UnitInpt = (props: IProps) => {
  const { field, form, parentName, onChange, value } = props;
  const [options, setOptions] = useState<any>([]);
  const [curUnitValue, setCurUnitValue] = useState<string>('');

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

  // 单位改变或输入框值改变时重新对输入框进行赋值
  const handleChangeUnit = ({
    value,
    isBlur,
  }: {
    value?: string;
    isBlur?: boolean;
  }) => {
    const unitValue = value || curUnitValue;
    // 如果存在值
    if (form.getFieldValue(formatName(field?.id))) {
      const inputValue = form.getFieldValue(formatName(field?.id));
      let newValue = inputValue;
      // 找出原始单位
      const temp = options.filter((item: any) => inputValue.includes(item))[0];
      // 填充值有单位且当前动作是失去焦点，不做任何操作
      if (temp && isBlur) {
        return;
      }
      if (temp) {
        // 将原始单位替换为新单位
        newValue = inputValue.replace(temp, unitValue);
      } else {
        // 没有单位直接将填充值与新单位拼接
        newValue = inputValue + unitValue;
      }
      form.setFieldValue(formatName(field?.id), newValue);
    }
  };

  // 渲染单位
  const selectAfter = () => {
    return (
      <Select
        defaultValue={curUnitValue}
        value={curUnitValue}
        onChange={(value) => {
          setCurUnitValue(value);
          handleChangeUnit({ value });
        }}
        style={{ width: 90 }}
      >
        {options?.map((item: string) => {
          return (
            <Select.Option value={item} key={item}>
              {item}
            </Select.Option>
          );
        })}
      </Select>
    );
  };

  useEffect(() => {
    if (field) {
      const { defaultValue, unit } = field || {};
      const list = unit.split(',');
      setOptions(list || []);
      let unitInit = list[0];

      const temp = list.filter((item: any) => defaultValue.includes(item))[0];
      if (temp) {
        unitInit = temp;
      }
      setCurUnitValue(unitInit);
    }
  }, [field]);

  return (
    <>
      <Input
        placeholder={getIntlText(field, 'keyCn', 'key')}
        addonAfter={selectAfter()}
        onChange={onChange}
        value={value}
        onBlur={() => {
          handleChangeUnit({ isBlur: true });
        }}
      />
    </>
  );
};

export default UnitInpt;
