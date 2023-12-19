import { queryContainersNameList } from '@/services/chaosmeta/KubernetesController';
import { history, useIntl, useRequest } from '@umijs/max';
import { Empty, Select, Spin, Tooltip, message } from 'antd';
import { useEffect, useState } from 'react';
import { GrayText, GroupLabel, OptionRow } from './style';

interface IProps {
  value?: string;
  onChange?: (value?: string) => void;
  /**回显时可传入list */
  list?: any[];
  mode?: 'multiple' | 'tags';
  placeholder?: string;
  style?: any;
  form?: any;
  kubernetesNamespace?: string;
  popupMatchSelectWidth: number;
}

interface IContainerName {
  container: string;
  pods: string;
}
const { Option, OptGroup } = Select;

const KubernetesContainersNameSelect = (props: IProps) => {
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
    popupMatchSelectWidth,
  } = props;

  const [containerList, setContainerList] = useState<IContainerName[]>([]);
  // 是否上第一次请求container列表,编辑配置时第一次请求不需要设置
  const [firstInitFlag, setFirstInitFlag] = useState<boolean>(true);

  useEffect(() => {
    if (list && list?.length > 0) {
      setContainerList(list);
    }
  }, [list]);

  const { run, loading } = useRequest(queryContainersNameList, {
    manual: true,
    formatResult: (res: any) => res,
    debounceInterval: 300,
    onSuccess: (res) => {
      const { success, data } = res;
      if (success) {
        if (!data?.containers) {
          return;
        }
        setContainerList(data.containers);

        // 编辑实验配置
        const { experimentId } = history?.location?.query || {};
        if (experimentId && firstInitFlag) {
          setFirstInitFlag(false);
          return;
        }
        // 创建实验默认是firstcontainer
        form.setFieldValue(['exec_range', 'target_sub_name'], 'firstcontainer');
      } else {
        message.error(res?.message);
      }
    },
  });

  useEffect(() => {
    // 此项的展示是要依赖namespace的值进行检索，所以当namespace值改变时需要重新请求
    if (kubernetesNamespace) {
      run(kubernetesNamespace);
    }
  }, [kubernetesNamespace]);

  return (
    <Select
      mode={mode}
      value={value}
      defaultValue={'firstcontainer'}
      allowClear
      popupMatchSelectWidth={popupMatchSelectWidth}
      notFoundContent={
        loading ? (
          <Spin size="small" />
        ) : (
          <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
        )
      }
      filterOption={false}
      onChange={onChange}
      placeholder={placeholder}
      style={style}
      optionLabelProp={'value'}
    >
      {containerList.length > 0 && (
        <OptGroup
          label={
            <GroupLabel>
              <div className="container">container</div>
              <div className="pods">related pods</div>
            </GroupLabel>
          }
        >
          {containerList.map(({ pods, container }) => {
            return (
              <Option value={container} key={container}>
                <OptionRow className={'ellipsis'}>
                  <GrayText>{container}</GrayText>
                  &nbsp;&nbsp;
                  <Tooltip title={pods}>
                    <span>{pods}</span>
                  </Tooltip>
                </OptionRow>
              </Option>
            );
          })}
        </OptGroup>
      )}
    </Select>
  );
};

export default KubernetesContainersNameSelect;
