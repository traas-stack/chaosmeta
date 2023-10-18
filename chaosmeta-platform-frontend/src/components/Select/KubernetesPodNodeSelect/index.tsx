import { queryPodNodeList } from '@/services/chaosmeta/KubernetesController';
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
  } = props;
  const [namespaceList, setNamespaceList] = useState<string[]>([]);
  const { Option } = Select;

  useEffect(() => {
    if (list && list?.length > 0) {
      setNamespaceList(list);
    }
  }, [list]);

  const getPodList = useRequest(queryPodNodeList, {
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
    getPodList?.run({
      page: 1,
      page_size: 500,
    });
  }, []);

  return (
    <Select
      mode={mode}
      value={value}
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
