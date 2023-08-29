import EmptyCustom from '@/components/EmptyCustom';
import { triggerTypes } from '@/constants';
import {
  queryExperimentList,
  queryExperimentResultList,
  stopExperimentResult,
} from '@/services/chaosmeta/ExperimentController';
import { querySpaceOverview } from '@/services/chaosmeta/SpaceController';
import { cronTranstionCN, formatTime } from '@/utils/format';
import { useParamChange } from '@/utils/useParamChange';
import { RightOutlined } from '@ant-design/icons';
import { history, useModel, useRequest } from '@umijs/max';
import {
  Button,
  Card,
  Col,
  Popconfirm,
  Row,
  Select,
  Space,
  Spin,
  Table,
  Tabs,
  message,
} from 'antd';
import moment from 'moment';
import { useEffect, useState } from 'react';

export default () => {
  // 实验列表数据
  const [experimentData, setExperimentData] = useState<any>({});
  // 总览数量展示数据
  const [overViewInfo, setOverviewInfo] = useState<any>({});
  // 实验结果列表数据
  const [experimentResultData, setExperimentResultData] = useState<any>({});
  // 选中的tabkey
  const [tabKey, setTabkey] = useState<string>('update');
  // 当前日期范围
  const [curTime, setCurTime] = useState('7day');
  const namespaceId = history?.location?.query?.spaceId as string;
  const spaceIdChange = useParamChange('spaceId');
  const { spacePermission } = useModel('global');
  /**
   * 获取实验列表
   */
  const getExperimentList = useRequest(queryExperimentList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        setExperimentData(res?.data);
      }
    },
  });

  /**
   * 获取实验结果列表
   */
  const getExperimentResultList = useRequest(queryExperimentResultList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        setExperimentResultData(res?.data);
      }
    },
  });

  /**
   * 获取空间概览信息
   */
  const getSpaceOverview = useRequest(querySpaceOverview, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        setOverviewInfo(res?.data);
      }
    },
  });

  /**
   * 停止实验
   */
  const handleStopExperiment = useRequest(
    (params: { uuid: string; name: string }) =>
      stopExperimentResult({ uuid: params?.uuid }),
    {
      manual: true,
      formatResult: (res) => res,
      onSuccess: (res, params) => {
        if (res?.code === 200) {
          message.success(`${params[0]?.name}实验已停止`);
          getExperimentResultList?.run({
            page: experimentResultData?.page || 1,
            page_size: experimentResultData?.pageSize || 10,
            namespace_id: namespaceId,
          });
        }
      },
    },
  );

  const operations = (
    <Space>
      <span
        style={{ cursor: 'pointer' }}
        onClick={() => {
          history.push({
            pathname: '/space/experiment',
          });
        }}
      >
        查看全部实验
      </span>
      <RightOutlined />
    </Space>
  );

  // 实验列表
  const columns: any[] = [
    {
      dataIndex: 'name',
      width: 120,
      render: (text: string, record: { update_time: string }) => {
        return (
          <>
            <div className="ellipsis row-text-gap">
              <a>
                <span>{text}</span>
              </a>
            </div>
            <div className="shallow">
              最近编辑时间： {formatTime(record?.update_time)}
            </div>
          </>
        );
      },
    },
    {
      dataIndex: 'schedule_type',
      width: 110,
      render: (text: string, record: { schedule_rule: string }) => {
        const value = triggerTypes?.filter((item) => item?.value === text)[0]
          ?.label;
        if (text === 'cron') {
          return (
            <div>
              <div className="shallow row-text-gap">触发方式</div>
              <div>
                <span>{value}</span>
                <span className="shallow">
                  {`（${cronTranstionCN(record?.schedule_rule)}）`}
                </span>
              </div>
            </div>
          );
        }
        return (
          <div>
            <div className="shallow row-text-gap">触发方式</div>
            <div>
              <span>{value}</span>
            </div>
          </div>
        );
      },
    },
    {
      dataIndex: 'number',
      width: 60,
      render: (text: number, record: { uuid: string }) => {
        return (
          <div>
            <div className="shallow row-text-gap">实验次数</div>
            {text > 0 ? (
              <a
                onClick={() => {
                  history?.push({
                    pathname: '/space/space/experiment-result',
                    query: {
                      experimentId: record?.uuid,
                    },
                  });
                }}
              >
                {text}
              </a>
            ) : (
              <span>0</span>
            )}
          </div>
        );
      },
    },
  ];

  // 实验的操作列，有权限才展示
  const experimentOperate = [
    {
      width: 50,
      fixed: 'right',
      render: (record: { uuid: string }) => {
        return (
          <Space>
            <a
              onClick={() => {
                history.push({
                  pathname: '/space/experiment/add',
                  query: {
                    experimentId: record?.uuid,
                  },
                });
              }}
            >
              编辑
            </a>
          </Space>
        );
      },
    },
  ];

  // 实验结果列表
  const resultColumns: any[] = [
    {
      dataIndex: 'name',
      width: 160,
      render: (text: string) => {
        return (
          <>
            <div className="ellipsis row-text-gap">
              <a>
                <span>{text}</span>
              </a>
            </div>
          </>
        );
      },
    },
    {
      dataIndex: 'create_time',
      width: 160,
      render: (text: string) => {
        return (
          <>
            <div className="shallow row-text-gap">发起时间</div>
            {formatTime(text)}
          </>
        );
      },
    },
    {
      dataIndex: 'update_time',
      width: 160,
      render: (text: string) => {
        return (
          <>
            <div className="shallow row-text-gap">结束时间</div>
            {formatTime(text)}
          </>
        );
      },
    },
    {
      dataIndex: 'status',
      width: 100,
      render: (text: string) => {
        return (
          <>
            <div className="shallow row-text-gap">实验状态</div>
            {text || '-'}
          </>
        );
      },
    },
  ];

  // 实验结果的操作列，有权限才展示
  const operateResultColumns: any[] = [
    {
      dataIndex: 'name',
      width: 120,
      render: (record: { uuid: string; status: string; name: string }) => {
        // 运行中才可以停止
        if (record?.status === 'Running') {
          return (
            <Space>
              <Popconfirm
                title="你确定要停止吗？"
                onConfirm={() => {
                  handleStopExperiment?.run({
                    uuid: record?.uuid,
                    name: record?.name,
                  });
                }}
              >
                <a>停止</a>
              </Popconfirm>
            </Space>
          );
        }
        return null;
      },
    },
  ];

  // 最近编辑的实验
  const EditExperimental = () => {
    return (
      <Card>
        <Table
          locale={{
            emptyText: (
              <EmptyCustom
                desc={`当前页面暂无${
                  tabKey === 'running' ? '即将运行' : '最近编辑'
                }的实验`}
                title="您可以前往实验列表查看实验"
                btns={
                  <Button
                    type="primary"
                    onClick={() => {
                      history?.push('/space/experiment');
                    }}
                  >
                    前往实验列表
                  </Button>
                }
              />
            ),
          }}
          showHeader={false}
          columns={
            spacePermission === 1 ? [...columns, ...experimentOperate] : columns
          }
          rowKey={'uuid'}
          loading={getExperimentList?.loading}
          dataSource={experimentData?.experiments}
          pagination={{
            simple: true,
            total: experimentData?.total,
          }}
          onChange={(pagination: any) => {
            const { current, pageSize } = pagination;
            getExperimentList?.run({
              page: current,
              page_size: pageSize,
              namespace_id: namespaceId,
            });
          }}
          scroll={{ x: 760 }}
        />
      </Card>
    );
  };

  // 最近运行的实验结果
  const RecentlyRunExperimentalResult = () => {
    return (
      <Card>
        <Table
          locale={{
            emptyText: (
              <EmptyCustom
                desc="当前暂无最近运行的实验结果"
                title="您可以前往实验结果列表查看实验"
                btns={
                  <Button
                    type="primary"
                    onClick={() => {
                      history?.push('/space/experiment-result');
                    }}
                  >
                    前往实验结果列表
                  </Button>
                }
              />
            ),
          }}
          showHeader={false}
          columns={
            spacePermission === 1
              ? [...resultColumns, ...operateResultColumns]
              : resultColumns
          }
          rowKey={'uuid'}
          dataSource={experimentResultData?.results || []}
          loading={getExperimentResultList?.loading}
          scroll={{ x: 760 }}
          pagination={{
            simple: true,
            total: experimentResultData?.total,
          }}
          onChange={(pagination: any) => {
            const { current, pageSize } = pagination;
            getExperimentResultList?.run({
              page: current,
              page_size: pageSize,
              namespace_id: namespaceId,
            });
          }}
        />
      </Card>
    );
  };

  const items = [
    {
      label: '最近编辑的实验',
      key: 'update',
      children: <EditExperimental />,
    },
    {
      label: '即将运行的实验',
      key: 'running',
      children: <EditExperimental />,
    },
    {
      label: '最近运行的实验结果',
      key: 'runningResult',
      children: <RecentlyRunExperimentalResult />,
    },
  ];

  /**
   * 日期转换
   * @returns
   */
  const transformTime = () => {
    let start_time, end_time;
    if (curTime === '7day') {
      start_time = formatTime(moment().subtract(30, 'd').startOf('day'));
      end_time = formatTime(moment().endOf('day'));
    }
    if (curTime === '30day') {
      start_time = formatTime(moment().subtract(30, 'd').startOf('day'));
      end_time = formatTime(moment().endOf('day'));
    }
    return {
      start_time,
      end_time,
    };
  };

  // tab栏列表检索
  const handleTabSearch = (key?: string) => {
    const val = key || tabKey;
    if (val === 'runningResult') {
      getExperimentResultList?.run({
        page: 1,
        page_size: 10,
        time_search_field: 'create_time',
        start_time: transformTime()?.start_time,
        end_time: transformTime()?.end_time,
        namespace_id: namespaceId,
      });
    }
    // 即将运行
    if (val === 'running') {
      getExperimentList?.run({
        page: 1,
        page_size: 10,
        time_search_field: 'next_exec',
        start_time: transformTime()?.start_time,
        end_time: transformTime()?.end_time,
        namespace_id: namespaceId,
      });
    }
    // 最近编辑
    if (val === 'update') {
      getExperimentList?.run({
        page: 1,
        page_size: 10,
        time_search_field: 'update_time',
        start_time: transformTime()?.start_time,
        end_time: transformTime()?.end_time,
        namespace_id: namespaceId,
      });
    }
  };

  useEffect(() => {
    handleTabSearch();
    getSpaceOverview?.run({
      spaceId: namespaceId,
      recent_day: 7,
    });
  }, [spaceIdChange]);

  return (
    <>
      <div className="overview">
        <div className="top">
          <span className="title">空间总览</span>
          <Select
            onChange={(val) => {
              setCurTime(val);
              const recent_day = val === '30day' ? 30 : 7;
              getSpaceOverview?.run({
                spaceId: history?.location?.query?.spaceId as string,
                recent_day,
              });
              handleTabSearch();
            }}
            defaultValue={'7day'}
            options={[
              { label: '最近7天', value: '7day' },
              { label: '最近30天', value: '30day' },
            ]}
          />
        </div>
        <Spin spinning={getSpaceOverview?.loading}>
          <div className="result">
            <Row gutter={16} style={{ alignContent: 'center' }}>
              <Col span={8}>
                <Card style={{ display: 'flex', alignContent: 'center' }}>
                  <div>
                    <div className="shallow-65">新增实验</div>
                    <span className="count">
                      {overViewInfo?.total_experiments || 0}
                    </span>
                    <span className="unit">个</span>
                  </div>
                </Card>
              </Col>
              <Col span={16}>
                <Card>
                  <Row>
                    <Col span={12}>
                      <div style={{ position: 'relative' }}>
                        <div className="shallow-65">执行实验</div>
                        <span className="count">
                          {overViewInfo?.total_experiment_instances || 0}
                        </span>
                        <span className="unit">次</span>
                      </div>
                    </Col>
                    <Col span={12}>
                      <div>
                        <div className="shallow-65">执行失败</div>
                        <span className="count-error">
                          {overViewInfo?.failed_experiment_instances || 0}
                        </span>
                        <span className="unit">次</span>
                      </div>
                    </Col>
                  </Row>
                </Card>
              </Col>
            </Row>
          </div>
        </Spin>
      </div>

      <Tabs
        tabBarExtraContent={operations}
        items={items}
        activeKey={tabKey}
        onChange={(val) => {
          setTabkey(val);
          handleTabSearch(val);
        }}
      />
    </>
  );
};
