import DynamicForm from '@/components/DynamicForm';
import KubernetesNamespaceSelect from '@/components/Select/KubernetesNamespaceSelect';
import KubernetesPodNodeSelect from '@/components/Select/KubernetesPodNodeSelect';
import KubernetesPodSelect from '@/components/Select/KubernetesPodSelect';
import ShowText from '@/components/ShowText';
import { nodeTypeMap, nodeTypeMapUS } from '@/constants';
import { queryFaultNodeFields } from '@/services/chaosmeta/ExperimentController';
import { formatDuration } from '@/utils/format';
import {
  CheckOutlined,
  CloseOutlined,
  DeleteOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { getLocale, useIntl, useRequest } from '@umijs/max';
import {
  Button,
  Form,
  Input,
  InputNumber,
  Popconfirm,
  Select,
  Space,
  Spin,
  message,
} from 'antd';
import React, { useEffect, useState } from 'react';
import { NodeConfigContainer } from '../style';

interface IProps {
  form: any;
  activeCol: any;
  arrangeList: any[];
  setArrangeList: any;
  setActiveCol: any;
  disabled?: boolean;
}

/**
 * 节点信息配置
 * @param props
 * @returns
 */
const NodeConfig: React.FC<IProps> = (props) => {
  const {
    form,
    activeCol,
    setArrangeList,
    setActiveCol,
    disabled = false,
  } = props;
  const [editTitleState, setEditTitleState] = useState<boolean>(false);
  const [fieldList, setFieldList] = useState<any[]>([]);
  const [durationType, setDurationType] = useState<string>('second');
  const [kubernetesNamespace, setKubernetesNamespace] = useState<string>('');
  const intl = useIntl();

  /**
   * 更新节点属性的方法
   * @param key 属性名
   * @param value 属性值，可以是字符串或数字
   */
  const handleEditNode = (key: string, value: any) => {
    setArrangeList((result: any) => {
      const values = JSON.parse(JSON.stringify(result)); // 将 result 对象深拷贝一份
      const parentIndex = values?.findIndex(
        (item: { row: any }) => item?.row === activeCol?.parentId,
      );
      if (parentIndex !== -1 && activeCol?.index >= 0) {
        values[parentIndex].children[activeCol.index][key] = value;
        // 同时更新保存当前节点信息
        setActiveCol((origin: any) => {
          return { ...origin, [key]: value, state: true };
        });
      }
      return values; // 返回更新后的 values 数组
    });
  };

  /**
   * 修改标题需同步修改编排节点title
   */
  const handleEditTitle = () => {
    const curTitle = form.getFieldValue('name');
    if (!curTitle) {
      message.info(
        `${intl.formatMessage({ id: 'inputPlaceholder' })} ${intl.formatMessage(
          { id: 'name' },
        )}`,
      );
      return;
    }
    handleEditNode('name', curTitle);
    setEditTitleState(false);
  };

  /**
   * 确定更新节点信息
   */
  const hanldeUpdateNode = (params: any) => {
    let curTime = form.getFieldValue('duration');
    if (durationType === 'minute') {
      curTime = `${curTime}m`;
    } else {
      curTime = `${curTime}s`;
    }
    setArrangeList((result: any) => {
      const values = JSON.parse(JSON.stringify(result)); // 将 result 对象深拷贝一份
      const parentIndex = values?.findIndex(
        (item: { row: any }) => item?.row === activeCol?.parentId,
      );
      if (parentIndex !== -1 && activeCol?.index >= 0) {
        const oldValue = values[parentIndex].children[activeCol.index];
        values[parentIndex].children[activeCol.index] = {
          ...oldValue,
          ...params,
          // 节点信息配置完成标识
          nodeInfoState: true,
          duration: curTime,
        };
      }
      return values; // 返回更新后的 values 数组
    });
    // 关闭配置栏
    setActiveCol({});
  };

  const selectAfter = (
    <Select
      defaultValue={durationType}
      onChange={(value) => {
        setDurationType(value);
      }}
      style={{ width: 90 }}
    >
      <Select.Option value="second">
        {intl.formatMessage({ id: 'second' })}
      </Select.Option>
      <Select.Option value="minute">
        {intl.formatMessage({ id: 'minute' })}
      </Select.Option>
    </Select>
  );

  // 表单赋值
  const handleFormAssignment = () => {
    // 初始化给表单赋值
    const initSecond = formatDuration(activeCol?.duration);
    let target_name = activeCol?.exec_range?.target_name;
    if (!Array.isArray(target_name)) {
      target_name = activeCol?.exec_range?.target_name
        ? activeCol?.exec_range?.target_name?.split(',')
        : undefined;
    }
    setKubernetesNamespace(activeCol?.exec_range?.target_namespace);
    form.setFieldsValue({
      ...activeCol,
      duration: initSecond,
      exec_type_name: (getLocale() === 'en-US' ? nodeTypeMapUS : nodeTypeMap)[
        activeCol?.exec_type
      ],
      exec_range: {
        ...activeCol?.exec_range,
        target_name,
      },
    });
  };

  /**
   * 故障节点 - 查询节点表单配置信息
   */
  const getFaultNodeFields = useRequest(queryFaultNodeFields, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res: any) => {
      if (res?.code === 200) {
        setFieldList(res?.data?.args || []);
        handleFormAssignment();
      }
    },
  });

  /**
   * @description: 处理删除节点的方法
   */
  const handleDeleteNode = () => {
    setArrangeList((result: any) => {
      const values = JSON.parse(JSON.stringify(result));
      const parentIndex = values?.findIndex(
        (item: { row: any }) => item?.row === activeCol?.parentId,
      );
      if (parentIndex !== -1 && activeCol?.index >= 0) {
        values[parentIndex]?.children?.splice(activeCol?.index, 1);
      }
      return values;
    });
    // 设置选中状态为false
    setActiveCol({ state: false });
  };

  /**
   * 确认节点配置
   */
  const handleConfirm = () => {
    form.validateFields().then((values: any) => {
      hanldeUpdateNode(values);
    });
  };

  useEffect(() => {
    form.resetFields();
    if (activeCol?.uuid) {
      // 请求动态表单渲染
      if (activeCol?.exec_id) {
        getFaultNodeFields?.run({ id: activeCol?.exec_id });
      }
      handleFormAssignment();
    }
  }, [activeCol?.uuid]);

  useEffect(() => {
    const initSecond = formatDuration(activeCol?.duration);
    form.setFieldsValue({
      duration: initSecond,
    });
  }, [activeCol?.duration]);

  return (
    <NodeConfigContainer>
      <Spin spinning={getFaultNodeFields?.loading}>
        <Form form={form} layout="vertical" disabled={disabled}>
          <div className="header">
            {!disabled ? (
              <Space>
                {editTitleState ? (
                  <Form.Item
                    name={'name'}
                    noStyle
                    rules={[
                      {
                        required: true,
                        message: intl.formatMessage({ id: 'inputPlaceholder' }),
                      },
                    ]}
                  >
                    <Input
                      placeholder={intl.formatMessage({
                        id: 'inputPlaceholder',
                      })}
                      style={{ width: '120px' }}
                    />
                  </Form.Item>
                ) : (
                  <Form.Item name={'name'} noStyle>
                    <ShowText />
                  </Form.Item>
                )}
                {editTitleState ? (
                  <Space>
                    <CloseOutlined
                      className="cancel"
                      style={{ color: '#FF4D4F' }}
                      onClick={() => {
                        setEditTitleState(false);
                      }}
                    />
                    <CheckOutlined
                      style={{ color: '#1890FF' }}
                      className="confirm"
                      onClick={() => {
                        handleEditTitle();
                      }}
                    />
                  </Space>
                ) : (
                  <EditOutlined
                    className="edit"
                    style={{ color: '#1890FF' }}
                    onClick={() => {
                      setEditTitleState(true);
                    }}
                  />
                )}
              </Space>
            ) : (
              <Form.Item name={'name'} noStyle>
                <ShowText />
              </Form.Item>
            )}
            <Space size={24}>
              {!disabled && (
                <Popconfirm
                  title={intl.formatMessage({ id: 'deleteConfirmText' })}
                  onConfirm={handleDeleteNode}
                >
                  <DeleteOutlined style={{ color: '#FF4D4F' }} />
                </Popconfirm>
              )}

              <CloseOutlined
                onClick={() => {
                  setActiveCol({ state: false });
                }}
              />
            </Space>
          </div>
          <div className="form">
            <Form.Item
              label={intl.formatMessage({ id: 'nodeType' })}
              name="exec_type_name"
            >
              <Input disabled />
            </Form.Item>
            <Form.Item
              label={
                activeCol?.exec_type === 'wait'
                  ? intl.formatMessage({ id: 'waitTime' })
                  : intl.formatMessage({ id: 'duration' })
              }
              name="duration"
              rules={[
                { required: true },
                {
                  validator: (_, value) => {
                    if ((value || value === 0) && value <= 0) {
                      return Promise.reject(
                        `${
                          activeCol?.exec_type === 'wait'
                            ? intl.formatMessage({ id: 'waitTime' })
                            : intl.formatMessage({ id: 'duration' })
                        } ${intl.formatMessage({ id: 'limit' })}`,
                      );
                    }
                    return Promise.resolve();
                  },
                },
              ]}
            >
              <InputNumber
                addonAfter={selectAfter}
                placeholder={intl.formatMessage({ id: 'inputPlaceholder' })}
                style={{ width: '100%' }}
              />
            </Form.Item>
            {/* 等待时长类型没有以下配置信息 */}
            {activeCol?.exec_type !== 'wait' && (
              <>
                {/* 动态表单部分 */}
                <DynamicForm
                  fieldList={fieldList}
                  parentName={'args_value'}
                  form={form}
                />
                <div className="range">
                  {intl.formatMessage({ id: 'attackRange' })}
                </div>
                {/* 一级节点为node时不展示Namespace */}
                {activeCol?.scope_id !== 2 && (
                  <Form.Item
                    label="Kubernetes Namespace"
                    name={['exec_range', 'target_namespace']}
                    rules={[
                      {
                        required: true,
                        message: intl.formatMessage({ id: 'inputPlaceholder' }),
                      },
                    ]}
                  >
                    <KubernetesNamespaceSelect
                      onChange={(val: any) => {
                        setKubernetesNamespace(val);
                      }}
                    />
                  </Form.Item>
                )}

                <Form.Item
                  label="Kubernetes Label"
                  name={['exec_range', 'target_label']}
                >
                  <Input
                    placeholder={intl.formatMessage({ id: 'inputPlaceholder' })}
                  />
                </Form.Item>
                {/* <Form.Item label="应用" name={['exec_range', 'target_app']}>
                  <Input placeholder="请输入" />
                </Form.Item> */}
                <Form.Item label="PodName" name={['exec_range', 'target_name']}>
                  {activeCol?.scope_id === 2 ? (
                    <KubernetesPodNodeSelect mode="multiple" />
                  ) : (
                    <KubernetesPodSelect
                      mode="multiple"
                      form={form}
                      kubernetesNamespace={kubernetesNamespace}
                    />
                  )}
                </Form.Item>
                {/* 一级节点为node时展示，node的id为2 */}
                {activeCol?.scope_id === 2 && (
                  <>
                    <Form.Item label="Ip" name={['exec_range', 'target_ip']}>
                      <Input
                        placeholder={intl.formatMessage({
                          id: 'inputPlaceholder',
                        })}
                      />
                    </Form.Item>
                  </>
                )}
                {/* <Form.Item
                  label="Hostname"
                  name={['exec_range', 'target_hostname']}
                >
                  <Input placeholder="请输入" />
                </Form.Item> */}
              </>
            )}
          </div>
          <div className="config-footer">
            <Button type="primary" onClick={handleConfirm}>
              {intl.formatMessage({ id: 'confirm' })}
            </Button>
          </div>
        </Form>
      </Spin>
    </NodeConfigContainer>
  );
};

export default NodeConfig;