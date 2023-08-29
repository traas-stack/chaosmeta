import DynamicForm from '@/components/DynamicForm';
import ShowText from '@/components/ShowText';
import { queryFaultNodeFields } from '@/services/chaosmeta/ExperimentController';
import { formatDuration } from '@/utils/format';
import {
  CheckOutlined,
  CloseOutlined,
  DeleteOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { useRequest } from '@umijs/max';
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

  console.log(activeCol, 'activeCol');

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
      message.info('请输入名称');
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
      style={{ width: 60 }}
    >
      <Select.Option value="second">秒</Select.Option>
      <Select.Option value="minute">分</Select.Option>
    </Select>
  );

  /**
   * 故障节点 - 查询节点表单配置信息
   */
  const getFaultNodeFields = useRequest(queryFaultNodeFields, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res: any) => {
      if (res?.code === 200) {
        setFieldList(res?.data?.args || []);
        // 初始化给表单赋值
        const initSecond = formatDuration(activeCol?.duration);
        form.setFieldsValue({
          ...activeCol,
          duration: initSecond,
        });
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
      // 初始化给表单赋值
      const initSecond = formatDuration(activeCol?.duration);
      form.setFieldsValue({
        ...activeCol,
        duration: initSecond,
      });
    }
  }, [activeCol?.uuid]);

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
                    rules={[{ required: true, message: '请输入' }]}
                  >
                    <Input placeholder="请输入" style={{ width: '120px' }} />
                  </Form.Item>
                ) : (
                  <Form.Item name={'name'} noStyle>
                    <ShowText value="CPU燃烧" />
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
                  title="你确定要删除吗？"
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
            <Form.Item label="节点类型" name="exec_type_name">
              <Input placeholder="节点类型" disabled />
            </Form.Item>
            <Form.Item
              label={activeCol?.exec_type === 'wait' ? '等待时长' : '持续时长'}
              name="duration"
              rules={[
                { required: true },
                {
                  validator: (_, value) => {
                    if ((value || value === 0) && value <= 0) {
                      return Promise.reject(
                        `${
                          activeCol?.exec_type === 'wait'
                            ? '等待时长'
                            : '持续时长'
                        }大于0`,
                      );
                    }
                    return Promise.resolve();
                  },
                },
              ]}
            >
              <InputNumber
                addonAfter={selectAfter}
                placeholder="请输入"
                // onChange={() => {
                //   handleTimeChange();
                // }}
                style={{ width: '100%' }}
              />
            </Form.Item>
            {/* 等待时长类型没有以下配置信息 */}
            {activeCol?.exec_type !== 'wait' && (
              <>
                {/* 动态表单部分 */}
                <DynamicForm fieldList={fieldList} parentName={'args_value'} />
                <div className="range">攻击范围</div>
                <Form.Item
                  label="Kubernetes Namespace"
                  name={['exec_range', 'target_namespace']}
                  rules={[{ required: true, message: '请输入' }]}
                >
                  <Input placeholder="请输入" />
                </Form.Item>
                <Form.Item
                  label="Kubernetes Label"
                  name={['exec_range', 'target_lebel']}
                >
                  <Input placeholder="请输入" />
                </Form.Item>
                <Form.Item label="应用" name={['exec_range', 'target_app']}>
                  <Input placeholder="请输入" />
                </Form.Item>
                <Form.Item label="name" name={['exec_range', 'target_name']}>
                  <Input placeholder="请输入" />
                </Form.Item>
                <Form.Item
                  label="Kubernetes Ip"
                  name={['exec_range', 'target_ip']}
                >
                  <Input placeholder="请输入" />
                </Form.Item>
                <Form.Item
                  label="Kubernetes Hostname"
                  name={['exec_range', 'target_hostname']}
                >
                  <Input placeholder="请输入" />
                </Form.Item>
              </>
            )}
          </div>
          <div className="config-footer">
            <Button type="primary" onClick={handleConfirm}>
              确认
            </Button>
          </div>
        </Form>
      </Spin>
    </NodeConfigContainer>
  );
};

export default NodeConfig;
