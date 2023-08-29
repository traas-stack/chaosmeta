import { queryNamespaceList } from '@/services/chaosmeta/KubernetesController';
import { useRequest } from '@umijs/max';
import { Select, Spin, message } from 'antd';
import { useEffect, useState } from 'react';

interface IProps {
  value?: string;
  onChange?: (value?: string) => void;
  /**回显时可传入list */
  list?: any[];
  mode?: 'multiple' | 'tags';
  placeholder?: string;
  style?: any;
}

const KubernetesNamespaceSelect = (props: IProps) => {
  const { value, onChange, list, mode, placeholder = '请选择', style } = props;
  const [namespaceList, setNamespaceList] = useState<string[]>([]);
  const { Option } = Select;

  useEffect(() => {
    if (list && list?.length > 0) {
      setNamespaceList(list);
    }
  }, [list]);

  const getUser = useRequest(queryNamespaceList, {
    manual: true,
    formatResult: (res: any) => res,
    debounceInterval: 300,
    onSuccess: (res) => {
      if (res?.success) {
        setNamespaceList(res?.data?.list || []);
      } else {
        message.error(res?.message);
      }
    },
  });

  // 输入有值时进行搜索
  // const handleUserSearch = (val: string) => {
  //   if (val) {
  //     getUser?.run({ identifier: val });
  //   }
  // };

  useEffect(() => {
    getUser?.run({ page: 1, page_size: 500 });
  }, []);

  return (
    <Select
      mode={mode}
      value={value}
      showSearch
      // onSearch={(val) => handleUserSearch(val)}
      allowClear
      notFoundContent={getUser?.loading ? <Spin size="small" /> : null}
      filterOption={false}
      onChange={onChange}
      placeholder={placeholder}
      style={style}
    >
      {namespaceList?.map((item: any) => {
        return (
          <Option value={item?.metadata?.name} key={item?.metadata?.name}>
            {item?.metadata?.name}
          </Option>
        );
      })}
    </Select>
  );
};

export default KubernetesNamespaceSelect;
