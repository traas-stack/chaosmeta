import { queryDeploymentNameList } from '@/services/chaosmeta/KubernetesController';
import { useIntl, useRequest } from '@umijs/max';
import { Empty, Select, Spin, message } from 'antd';
import { useEffect, useState } from 'react';

interface IProps {
  value?: string;
  onChange?: (value?: string) => void;
  /**回显时可传入list */
  list?: any[];
  mode?: 'multiple' | 'tags';
  placeholder?: string;
  style?: any;
  form?: any;
  kubernetesNamespace?: any;
}

const KubernetesPodSelect = (props: IProps) => {
  const intl = useIntl();
  const {
    value,
    onChange,
    list,
    mode,
    placeholder = intl.formatMessage({ id: 'selectPlaceholder' }),
    style,
    form,
    kubernetesNamespace,
  } = props;
  const [namespaceList, setNamespaceList] = useState<string[]>([]);
  const { Option } = Select;

  useEffect(() => {
    if (list && list?.length > 0) {
      setNamespaceList(list);
    }
  }, [list]);

  const getPodList = useRequest(queryDeploymentNameList, {
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

  useEffect(() => {
    // 此项的展示是要依赖namespace的值进行检索，所以当namespace值改变时需要清空列表项和值
    setNamespaceList([]);
    form.setFieldValue(['exec_range', 'target_name'], undefined);
  }, [kubernetesNamespace]);

  return (
    <Select
      mode={mode}
      value={value}
      // showSearch
      // onSearch={(val) => handleUserSearch(val)}
      allowClear
      notFoundContent={
        getPodList?.loading ? (
          <Spin size="small" />
        ) : (
          <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
        )
      }
      filterOption={false}
      onChange={onChange}
      placeholder={placeholder}
      style={style}
      onFocus={() => {
        // 依赖的namespace有值时才进行检索
        if (kubernetesNamespace) {
          getPodList?.run({
            page: 1,
            page_size: 500,
            namespace: kubernetesNamespace,
          });
        }
      }}
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

export default KubernetesPodSelect;
