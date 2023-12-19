import ShowText from '@/components/ShowText';
import {
  arrangeNodeTypeColors,
  nodeTypeMap,
  nodeTypeMapUS,
  nodeTypes,
  scaleStepMap,
} from '@/constants';
import {
  queryFaultNodeFields,
  queryFlowNodeFields,
  queryMeasureNodeFields,
} from '@/services/chaosmeta/ExperimentController';
import { queryFaultNodeDetail } from '@/services/chaosmeta/KubernetesController';
import {
  formatDuration,
  getIntlLabel,
  handleTimeTransform,
} from '@/utils/format';
import {
  CheckCircleFilled,
  ZoomInOutlined,
  ZoomOutOutlined,
} from '@ant-design/icons';
import { getLocale, history, useIntl, useRequest } from '@umijs/max';
import { Form, Space, Spin } from 'antd';
import { useEffect, useState } from 'react';
import DynamicFormRender from '../AddExperiment/components/DynamicFormRender';
import { ArrangeWrap, DroppableCol, DroppableRow } from './style';

interface IProps {
  arrangeList: any[];
  curExecSecond?: number; // 当前执行到的时间
  // 以下都是结果详情需要的
  isResult?: boolean;
  getExperimentArrangeNodeDetail?: any;
  setCurNodeDetail?: any;
}
const ArrangeInfoShow: React.FC<IProps> = (props) => {
  const {
    arrangeList,
    curExecSecond,
    isResult,
    getExperimentArrangeNodeDetail,
    setCurNodeDetail,
  } = props;
  // 当前占比
  const [curProportion, setCurProportion] = useState<number>(100);
  const [timeCount, setTimeCount] = useState<number>(16);
  const scaleStep = [33, 66, 100, 150, 200, 300];
  // 统计总时长
  const [totalDuration, setTotalDuration] = useState(0);
  const listMin = [...Array(timeCount)].map((x, i) => i);
  const [activeCol, setActiveCol] = useState<any>({ state: false });
  const [fieldList, setFieldList] = useState<any[]>([]);
  const [configForm] = Form.useForm();
  const intl = useIntl();
  // 用于判断当前节点是否为 kubernetes node或 kubernetes pod节点下
  const [targetName, setTargetName] = useState<string>('');
  /**
   * 故障节点 - 查询节点表单配置信息
   */
  const getFaultNodeFields = useRequest(queryFaultNodeFields, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res: any) => {
      if (res?.code === 200) {
        setFieldList(res?.data?.args || []);
      }
    },
  });

  /**
   * 流量注入 - 查询节点表单配置信息
   */
  const getFlowNodeFields = useRequest(queryFlowNodeFields, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res: any) => {
      if (res?.code === 200) {
        setFieldList(res?.data?.args || []);
      }
    },
  });

  /**
   * 度量引擎 - 查询节点表单配置信息
   */
  const getMeasureNodeFields = useRequest(queryMeasureNodeFields, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res: any) => {
      if (res?.code === 200) {
        setFieldList(res?.data?.args || []);
      }
    },
  });

  /**
   * 根据targetid获取该节点信息，用于判断该节点是否位于node或pod下
   */
  const getFaultNodeDetail = useRequest(queryFaultNodeDetail, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res: any) => {
      if (res?.code === 200) {
        setTargetName(res?.data?.name);
      }
    },
  });

  /**
   * @description: 处理函数，计算二级列表中所有子项的总时长并更新到总时长上
   */
  const handleTotalSecond = () => {
    // 初始化总时长为0
    let totalSecond: number = 0;
    // 遍历二级列表中的每一个元素
    arrangeList?.forEach((item) => {
      // 遍历该元素的每一个子项
      item?.children?.forEach((el: { duration: string }) => {
        const second = formatDuration(el?.duration);
        // 将子项的秒数累加到总时长上
        totalSecond += second;
      });
    });
    // 更新总时长
    setTotalDuration(totalSecond);
  };

  /**
   * 时间轴渲染
   * @param index
   * @returns
   */
  const renderTimeItem = (index: number) => {
    const secondStep = scaleStepMap[curProportion]?.secondStep;
    const second = index * secondStep;
    const text = handleTimeTransform(second);
    // 时间轴距离间隔固定90px
    return (
      <div key={index} className="time-item" style={{ width: `90px` }}>
        {text}
      </div>
    );
  };

  // 点击节点后的操作
  const handleNodeClick = async (el: any) => {
    if (activeCol?.uuid === el?.uuid) {
      setActiveCol({});
    } else {
      // wait类型不需要请求接口
      if (el?.exec_type === 'wait') {
        // 为配置信息赋值，实验详情的操作
        configForm.setFieldsValue(el);
        setActiveCol({
          ...el,
        });
        if (setCurNodeDetail) {
          setCurNodeDetail(el);
        }
        return;
      }
      // 获取节点动态表单部分
      if (el?.exec_type === 'flow') {
        getFlowNodeFields?.run({ id: el?.exec_id });
      }
      if (el?.exec_type === 'measure') {
        getMeasureNodeFields?.run({ id: el?.exec_id });
      }
      if (el?.exec_type === 'fault') {
        getFaultNodeFields?.run({ id: el?.exec_id });
        getFaultNodeDetail?.run({ targetId: el?.target_id });
      }
      // 结果详情中需要通过接口获取节点信息，实验详情中则在详情中直接返回了
      if (isResult) {
        // 实验结果时点击节点，获取单个节点详情后赋值
        getExperimentArrangeNodeDetail
          ?.run({
            uuid: history?.location?.query?.resultId,
            node_id: el?.uuid,
          })
          .then((res: { code: number; data: { workflow_node: any } }) => {
            if (res?.code === 200) {
              // 将args_value即动态表单部分数据转化为form需要的
              const curNodeDetail = res?.data?.workflow_node;
              if (curNodeDetail?.args_value) {
                const newArgs: any = {};
                curNodeDetail?.args_value?.forEach((arg: any) => {
                  newArgs[arg?.args_id] = arg?.value;
                });
                curNodeDetail.args_value = newArgs;
              }
              if (el?.exec_type === 'flow') {
                // 为配置信息赋值
                configForm.setFieldsValue({
                  ...curNodeDetail,
                  flow_range: {
                    ...curNodeDetail?.flow_subtasks,
                  },
                });
                setActiveCol({
                  ...curNodeDetail,
                  flow_range: {
                    ...curNodeDetail?.flow_subtasks,
                  },
                });
              }
              if (el?.exec_type === 'measure') {
                // 为配置信息赋值
                configForm.setFieldsValue({
                  ...curNodeDetail,
                  measure_range: {
                    ...curNodeDetail?.measure_subtasks,
                  },
                });
                setActiveCol({
                  ...curNodeDetail,
                  measure_range: {
                    ...curNodeDetail?.measure_subtasks,
                  },
                });
              }
              if (el?.exec_type === 'fault') {
                // 为配置信息赋值
                configForm.setFieldsValue({
                  ...curNodeDetail,
                  exec_range: {
                    ...curNodeDetail?.subtasks,
                  },
                });
                setActiveCol({
                  ...curNodeDetail,
                  exec_range: {
                    ...curNodeDetail?.subtasks,
                  },
                });
              }
            }
          });
        return;
      }
      // 为配置信息赋值，实验详情的操作
      configForm.setFieldsValue(el);
      setActiveCol({
        ...el,
      });
    }
  };

  /**
   * 每行信息的展示
   * @param props
   * @returns
   */
  const ArrangeRow = (props: any) => {
    const { item, index } = props;
    return (
      <div>
        <DroppableRow>
          <div className="row">
            {/* 行内子元素 */}
            {item?.children?.map((el: any, j: number) => {
              const curDuration = formatDuration(el?.duration);
              // 失败或者错误状态
              const isError = el?.status === 'Failed' || el?.status === 'error';
              return (
                <DroppableCol
                  key={j}
                  $bg={arrangeNodeTypeColors[el?.exec_type]}
                  $nodeStutas={isResult && el?.status}
                  // 减去外边距的2px，避免子元素多时宽度偏差过大
                  style={{
                    width: `${
                      curDuration *
                        (scaleStepMap[curProportion]?.widthSecond || 3) -
                      2
                    }px`,
                    // 最小宽度为1s对应的px
                    minWidth: `${scaleStepMap[curProportion]?.widthSecond}px`,
                    flexShrink: 0,
                  }}
                  $activeState={activeCol?.uuid === el?.uuid}
                  onClick={() => {
                    handleNodeClick(el);
                  }}
                >
                  <div className="item">
                    {curDuration *
                      (scaleStepMap[curProportion]?.widthSecond || 3) >
                    30 ? (
                      <div>
                        <div
                          className="title ellipsis"
                          style={{ paddingRight: isError ? '12px' : '4px' }}
                        >
                          <span>{el.name}</span>
                        </div>
                        <div className="duration">
                          <span>{curDuration}s</span>
                        </div>
                        {isError && (
                          <span className="tip-icon">
                            <img
                              src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*Qp1MT7UkGCQAAAAAAAAAAAAADmKmAQ/original"
                              alt=""
                            />
                          </span>
                        )}
                        {el?.status === 'Succeeded' && (
                          <span
                            className="tip-icon"
                            style={{ color: '#52c41a' }}
                          >
                            <CheckCircleFilled />
                          </span>
                        )}
                      </div>
                    ) : (
                      <>
                        <div>...</div>
                        <div>...</div>
                      </>
                    )}
                  </div>
                </DroppableCol>
              );
            })}
          </div>
          <div className="handle">{index + 1}</div>
        </DroppableRow>
      </div>
    );
  };

  // 是否为node下的节点
  const isNode = () => {
    // 只有故障节点下的才有可能有node节点
    if (activeCol?.exec_type !== 'fault') {
      return false;
    }
    // 父节点有两个node，一种是scope_id为2，另一种是通过接口查询name为node
    return activeCol?.scope_id === 2 || targetName === 'node';
  };

  // 是否为pod下的节点
  const isPod = () => {
    // 只有故障节点下的才有可能有node节点
    if (activeCol?.exec_type !== 'fault') {
      return false;
    }
    // 父节点有两个pod，一种是scope_id为1，另一种是通过接口查询name为pod
    return activeCol?.scope_id === 1 || targetName === 'pod';
  };

  /**
   * 遍历编排数组，将其中行中最长秒数对比当前默认时间轴，若时间轴宽度不够就再加长
   * @param values
   */
  const handleAddTimeAxis = (values: any) => {
    const maxSecondList: any = [];
    values?.forEach((item: { children: any[] }) => {
      let secondSum = 0;
      item?.children?.forEach((el) => {
        // 将持续时长转化为数字形式进行计算
        const duration = formatDuration(el?.duration);
        secondSum += duration;
      });
      maxSecondList.push(secondSum);
    });
    const maxSecond = maxSecondList.sort((a: number, b: number) => b - a)[0];
    // 提前一段加长
    const curSecond = timeCount * 30 - 60;
    if (maxSecond > curSecond) {
      setTimeCount(() => {
        // 多加入90s
        const newCount = Math.round(maxSecond / 30) + 3;
        return newCount;
      });
    }
  };

  useEffect(() => {
    // 比例变化时，修改时间间隔的数量，不低于屏幕宽度的秒数，避免出现空白区域，默认1000s
    const doc = document.body;
    const secondStep = scaleStepMap[curProportion]?.secondStep;
    const widthSecond = scaleStepMap[curProportion]?.widthSecond;
    const second = doc.clientWidth / widthSecond;
    setTimeCount(() => {
      const minCount = Math.round(second / secondStep);
      const curCount = Math.round(1000 / secondStep);
      return curCount > minCount ? curCount : minCount;
    });
  }, [curProportion]);

  // 攻击范围下不同节点渲染不同
  const attackRangeRender = () => {
    // 父节点为node时
    if (isNode()) {
      return (
        <>
          <Form.Item
            label="Kubernetes Label"
            name={['exec_range', 'target_label']}
          >
            <ShowText ellipsis />
          </Form.Item>
          <Form.Item label={'NodeName'} name={['exec_range', 'target_name']}>
            <ShowText ellipsis />
          </Form.Item>
          <Form.Item label="Ip" name={['exec_range', 'target_ip']}>
            <ShowText ellipsis />
          </Form.Item>
        </>
      );
    }

    // 父节点为pod时
    if (isPod()) {
      return (
        <>
          <Form.Item
            label="Kubernetes Namespace"
            name={['exec_range', 'target_namespace']}
          >
            <ShowText ellipsis />
          </Form.Item>
          <Form.Item
            label="Kubernetes Label"
            name={['exec_range', 'target_label']}
          >
            <ShowText ellipsis />
          </Form.Item>
          <Form.Item label={'PodName'} name={['exec_range', 'target_name']}>
            <ShowText ellipsis />
          </Form.Item>
          <Form.Item label={'ContainersName'} name={['exec_range', 'target_sub_name']}>
            <ShowText ellipsis />
          </Form.Item>
        </>
      );
    }
    // 父节点为deployment时
    if (targetName === 'deployment') {
      return (
        <>
          <Form.Item
            label="Kubernetes Namespace"
            name={['exec_range', 'target_namespace']}
          >
            <ShowText ellipsis />
          </Form.Item>
          <Form.Item
            label={'DeploymentName'}
            name={['exec_range', 'target_name']}
          >
            <ShowText ellipsis />
          </Form.Item>
        </>
      );
    }
  };

  useEffect(() => {
    handleAddTimeAxis(arrangeList);
  }, [arrangeList]);

  useEffect(() => {
    handleTotalSecond();
  }, []);

  return (
    <ArrangeWrap
      $activeItem={activeCol?.state}
      $curExecSecond={curExecSecond}
      $widthSecond={scaleStepMap[curProportion]?.widthSecond}
    >
      <div className="wrap">
        <div className="arrange">
          {/* 当前执行到的标识线 -- todo后端暂不支持 */}
          {/* {curExecSecond && (
            <div className="time-stage">
              <img
                src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*wTyzQ4qeH2kAAAAAAAAAAAAADmKmAQ/original"
                alt=""
              />
              <div className="line"></div>
            </div>
          )} */}
          {/* 顶部时间轴 */}
          <div
            className="time-axis"
            style={{
              minWidth: `${timeCount * 90}px`,
            }}
          >
            {listMin?.map((item, index) => {
              return renderTimeItem(index);
            })}
          </div>
          {/* 编排元素展示 */}
          {arrangeList?.map((item, index) => {
            return <ArrangeRow key={index} item={item} index={index} />;
          })}
        </div>
        {activeCol?.uuid && (
          <div className="info">
            <Spin spinning={getFaultNodeFields?.loading}>
              <Form form={configForm}>
                <div className="subtitle">
                  {intl.formatMessage({ id: 'configInfo' })}
                </div>
                <Form.Item
                  label={intl.formatMessage({ id: 'nodeName' })}
                  name={'name'}
                >
                  <ShowText ellipsis />
                </Form.Item>
                <Form.Item
                  label={intl.formatMessage({ id: 'nodeType' })}
                  // name={'exec_type'}
                >
                  <ShowText
                    value={
                      (getLocale() === 'en-US' ? nodeTypeMapUS : nodeTypeMap)[
                        activeCol?.exec_type
                      ] || activeCol?.exec_type
                    }
                  />
                </Form.Item>
                {activeCol?.exec_type !== 'flow' &&
                  activeCol?.exec_type !== 'measure' && (
                    <>
                      <Form.Item
                        label={intl.formatMessage({ id: 'atomicCapabilities' })}
                        name="exec_name"
                      >
                        <ShowText />
                      </Form.Item>
                      <div className="subtitle range">
                        {intl.formatMessage({ id: 'commonParameters' })}
                      </div>
                      <Form.Item
                        label={`${
                          activeCol?.exec_type === 'wait'
                            ? intl.formatMessage({ id: 'waitTime' })
                            : intl.formatMessage({ id: 'duration' })
                        }`}
                        name={'duration'}
                      >
                        <ShowText />
                      </Form.Item>
                    </>
                  )}

                {activeCol?.exec_type !== 'wait' && (
                  <>
                    {/* 动态表单部分 */}
                    <DynamicFormRender
                      fieldList={fieldList}
                      nodeType={activeCol?.exec_type}
                      readonly
                    />
                    {/* 节点父类型为node或pod 或deployment时才展示 攻击范围 */}
                    {(isNode() || isPod() || targetName === 'deployment') && (
                      <>
                        <div className="subtitle range">
                          {intl.formatMessage({ id: 'attackRange' })}
                        </div>
                        {/* node下的节点时不展示 */}
                        {attackRangeRender()}
                      </>
                    )}
                  </>
                )}
              </Form>
            </Spin>
          </div>
        )}

        {/* 底部 */}
        <div className="footer">
          <Space style={{ alignItems: 'center' }}>
            <div>
              {intl.formatMessage({ id: 'totalDuration' })}：
              <span className="total-time">
                {handleTimeTransform(totalDuration)}
              </span>
            </div>
            <Space className="node-type">
              {nodeTypes?.map((item: any) => {
                return (
                  <Space key={item.label} className="node-item">
                    <div
                      style={{
                        background: arrangeNodeTypeColors[item.type],
                      }}
                    ></div>
                    {getIntlLabel(item)}
                  </Space>
                );
              })}
            </Space>
          </Space>
          <Space>
            <ZoomOutOutlined
              style={{
                color: curProportion === 33 ? 'rgba(0,0,0,0.16)' : '',
              }}
              onClick={() => {
                if (curProportion > 33) {
                  setCurProportion(() => {
                    const curIndex = scaleStep?.findIndex(
                      (item) => item === curProportion,
                    );
                    return scaleStep[curIndex - 1];
                  });
                }
              }}
            />
            <span>{curProportion}%</span>
            <ZoomInOutlined
              style={{
                color: curProportion === 300 ? 'rgba(0,0,0,0.16)' : '',
              }}
              onClick={() => {
                if (curProportion < 300) {
                  setCurProportion(() => {
                    const curIndex = scaleStep?.findIndex(
                      (item) => item === curProportion,
                    );
                    return scaleStep[curIndex + 1];
                  });
                }
              }}
            />
          </Space>
        </div>
      </div>
    </ArrangeWrap>
  );
};

export default ArrangeInfoShow;
