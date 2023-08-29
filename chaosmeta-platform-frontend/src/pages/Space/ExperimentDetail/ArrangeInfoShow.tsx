import DynamicForm from '@/components/DynamicForm';
import ShowText from '@/components/ShowText';
import { arrangeNodeTypeColors, scaleStepMap } from '@/constants';
import { queryFaultNodeFields } from '@/services/chaosmeta/ExperimentController';
import { formatDuration, handleTimeTransform } from '@/utils/format';
import { ZoomInOutlined, ZoomOutOutlined } from '@ant-design/icons';
import { history, useRequest } from '@umijs/max';
import { Form, Space, Spin } from 'antd';
import { useEffect, useState } from 'react';
import { ArrangeWrap, DroppableCol, DroppableRow } from './style';

interface IProps {
  arrangeList: any[];
  curExecSecond?: number; // 当前执行到的时间
  // 以下都是结果详情需要的
  isResult?: boolean;
  getExperimentArrangeNodeDetail?: any;
}
const ArrangeInfoShow: React.FC<IProps> = (props) => {
  const {
    arrangeList,
    curExecSecond,
    isResult,
    getExperimentArrangeNodeDetail,
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

  const nodeTypes = [
    {
      name: '故障节点',
      type: 'fault',
    },
    {
      name: '度量节点',
      type: 'measure',
    },
    {
      name: '压测节点',
      type: 'pressure',
    },
    {
      name: '其他节点',
      type: 'other',
    },
  ];

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
      // 获取节点动态表单部分
      getFaultNodeFields?.run({ id: el?.exec_id });
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
              // 为配置信息赋值
              configForm.setFieldsValue(curNodeDetail);
              setActiveCol({
                ...curNodeDetail,
              });
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
                      // <div style={{ display: 'flex' }}>
                        <div>
                          <div className="title ellipsis">
                            <span>{el.name}</span>
                            {/* <span>
                              测试测试测试测试测试测试测试测试测试测试测试测试
                            </span> */}
                          </div>
                          <div>{curDuration}s</div>
                        </div>
                        // <div>x</div>
                      // </div>
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

  useEffect(() => {
    handleTotalSecond();
  }, []);

  useEffect(() => {});

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
                <div className="subtitle">配置信息</div>
                <Form.Item label="节点名称" name={'name'}>
                  <ShowText ellipsis />
                </Form.Item>
                <Form.Item label="节点类型" name={'exec_type'}>
                  <ShowText />
                </Form.Item>
                <Form.Item label="持续时长" name={'duration'}>
                  <ShowText />
                </Form.Item>
                {/* 动态表单部分 */}
                <DynamicForm
                  fieldList={fieldList}
                  parentName={'args_value'}
                  readonly
                />
                <div className="subtitle range">攻击范围</div>
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
                <Form.Item label="应用" name={['exec_range', 'target_app']}>
                  <ShowText ellipsis />
                </Form.Item>
                <Form.Item label="name" name={['exec_range', 'target_name']}>
                  <ShowText ellipsis />
                </Form.Item>
                <Form.Item
                  label="Kubernetes Ip"
                  name={['exec_range', 'target_ip']}
                >
                  <ShowText ellipsis />
                </Form.Item>
                <Form.Item
                  label="Kubernetes Hostname"
                  name={['exec_range', 'target_hostname']}
                >
                  <ShowText ellipsis />
                </Form.Item>
              </Form>
            </Spin>
          </div>
        )}

        {/* 底部 */}
        <div className="footer">
          <Space style={{ alignItems: 'center' }}>
            <div>
              总时长：
              <span className="total-time">
                {handleTimeTransform(totalDuration)}
              </span>
            </div>
            <Space className="node-type">
              {nodeTypes?.map((item) => {
                return (
                  <Space key={item.name} className="node-item">
                    <div
                      style={{
                        background: arrangeNodeTypeColors[item.type],
                      }}
                    ></div>
                    {item.name}
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
