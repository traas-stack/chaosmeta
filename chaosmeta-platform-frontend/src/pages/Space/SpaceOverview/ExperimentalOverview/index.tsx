import EmptyCustom from '@/components/EmptyCustom';
import { experimentResultStatus, triggerTypes } from '@/constants';
import {
  queryExperimentList,
  queryExperimentResultList,
  stopExperimentResult,
} from '@/services/chaosmeta/ExperimentController';
import { querySpaceOverview } from '@/services/chaosmeta/SpaceController';
import { cronTranstionCN, formatTime, getIntlLabel } from '@/utils/format';
import { useParamChange } from '@/utils/useParamChange';
import { RightOutlined } from '@ant-design/icons';
import { history, useIntl, useModel, useRequest } from '@umijs/max';
import {
  Badge,
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
import dayjs from 'dayjs';
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
  const [timeType, setCurTime] = useState('7day');
  const namespaceId = history?.location?.query?.spaceId as string;
  const spaceIdChange = useParamChange('spaceId');
  const { spacePermission } = useModel('global');
  const intl = useIntl();
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
   * 日期转换
   * @returns
   */
  const transformTime = (value?: string) => {
    let start_time, end_time;
    const curTime = value || timeType;
    if (curTime === '7day') {
      start_time = formatTime(dayjs().subtract(6, 'd').startOf('day'));
      end_time = formatTime(dayjs().endOf('day'));
    }
    if (curTime === '30day') {
      start_time = formatTime(dayjs().subtract(30, 'd').startOf('day'));
      end_time = formatTime(dayjs().endOf('day'));
    }
    return {
      start_time,
      end_time,
    };
  };

  // 最近运行实验结果接口
  const handleExperimentResultList = (params: {
    curTime?: string;
    pageSize?: number;
    pageIndex?: number;
  }) => {
    const { curTime, pageIndex, pageSize } = params;
    getExperimentResultList?.run({
      page: pageIndex || 1,
      page_size: pageSize || 10,
      time_search_field: 'create_time',
      start_time: transformTime(curTime)?.start_time,
      end_time: transformTime(curTime)?.end_time,
      namespace_id: namespaceId,
      time_type: 'range',
    });
  };

  // 实验列表
  const handleExperimentList = (params: {
    curTime?: string;
    pageSize?: number;
    pageIndex?: number;
    timeField?: string;
  }) => {
    const { curTime, pageIndex, pageSize } = params;
    getExperimentList?.run({
      page: pageIndex || 1,
      page_size: pageSize || 10,
      time_search_field: 'create_time',
      start_time: transformTime(curTime)?.start_time,
      end_time: transformTime(curTime)?.end_time,
      namespace_id: namespaceId,
      time_type: 'range',
    });
  };

  // tab栏列表检索
  const handleTabSearch = (params?: { key?: string; curTime?: string }) => {
    const { key, curTime } = params || {};
    const val = key || tabKey;
    if (val === 'runningResult') {
      handleExperimentResultList({ curTime });
    }
    // 即将运行
    if (val === 'running') {
      handleExperimentList({ timeField: 'next_exec', curTime });
    }
    // 最近编辑
    if (val === 'update') {
      handleExperimentList({ timeField: 'update_time', curTime });
    }
  };

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
          handleExperimentResultList({
            pageIndex: experimentResultData?.page || 1,
            pageSize: experimentResultData?.pageSize || 10,
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
        {intl.formatMessage({
          id: 'overview.spaceOverview.tab.more',
        })}
      </span>
      <RightOutlined />
    </Space>
  );

  // 实验列表
  const columns: any[] = [
    {
      dataIndex: 'name',
      width: 120,
      render: (text: string, record: { update_time: string; uuid: string }) => {
        return (
          <>
            <div className="ellipsis row-text-gap">
              <a
                onClick={() => {
                  history.push({
                    pathname: '/space/experiment/detail',
                    query: {
                      experimentId: record.uuid,
                      spaceId: history?.location?.query?.spaceId,
                    },
                  });
                }}
              >
                <span>{text}</span>
              </a>
            </div>
            <div className="shallow">
              {intl.formatMessage({ id: 'lastEditTime' })}：{' '}
              {formatTime(record?.update_time)}
            </div>
          </>
        );
      },
    },
    {
      dataIndex: 'schedule_type',
      width: 110,
      render: (text: string, record: { schedule_rule: string }) => {
        const temp = triggerTypes?.filter((item) => item?.value === text)[0];
        if (text === 'cron') {
          return (
            <div>
              <div className="shallow row-text-gap">
                {intl.formatMessage({ id: 'triggerMode' })}
              </div>
              <div>
                <span>{getIntlLabel(temp)}</span>
                <span className="shallow">
                  {`（${cronTranstionCN(record?.schedule_rule)}）`}
                </span>
              </div>
            </div>
          );
        }
        return (
          <div>
            <div className="shallow row-text-gap">
              {intl.formatMessage({ id: 'triggerMode' })}
            </div>
            <div>
              <span>{getIntlLabel(temp)}</span>
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
            <div className="shallow row-text-gap">
              {intl.formatMessage({ id: 'numberOfExperiments' })}
            </div>
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
              {intl.formatMessage({ id: 'edit' })}
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
      render: (text: string, record: { uuid: string }) => {
        return (
          <>
            <div className="ellipsis row-text-gap">
              <a>
                <span
                  onClick={() => {
                    history.push({
                      pathname: '/space/experiment-result/detail',
                      query: {
                        resultId: record?.uuid,
                        spaceId: history?.location?.query?.spaceId as string,
                      },
                    });
                  }}
                >
                  {text}
                </span>
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
            <div className="shallow row-text-gap">
              {intl.formatMessage({ id: 'startTime' })}
            </div>
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
            <div className="shallow row-text-gap">
              {intl.formatMessage({ id: 'endTime' })}
            </div>
            {formatTime(text)}
          </>
        );
      },
    },
    {
      dataIndex: 'status',
      width: 100,
      render: (text: string) => {
        const temp: any = experimentResultStatus?.filter(
          (item) => item?.value === text,
        )[0];
        return (
          <>
            <div className="shallow row-text-gap">
              {intl.formatMessage({ id: 'experimentStatus' })}
            </div>
            {temp ? (
              <div>
                <Badge color={temp?.color} /> {getIntlLabel(temp)}
              </div>
            ) : (
              '-'
            )}
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
                desc={
                  tabKey === 'running'
                    ? intl.formatMessage({
                        id: 'overview.spaceOverview.tab2.noAuth.empty.description',
                      })
                    : intl.formatMessage({
                        id: 'overview.spaceOverview.tab1.noAuth.empty.description',
                      })
                }
                title={intl.formatMessage({
                  id: 'overview.spaceOverview.tab1.noAuth.empty.title',
                })}
                btns={
                  <Button
                    type="primary"
                    onClick={() => {
                      history?.push('/space/experiment');
                    }}
                  >
                    {intl.formatMessage({
                      id: 'overview.spaceOverview.tab1.noAuth.empty.btn',
                    })}
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
            pageSize: experimentData?.page_size,
            current: experimentData?.page,
          }}
          onChange={(pagination: any) => {
            const { current, pageSize } = pagination;
            handleExperimentList({ pageIndex: current, pageSize });
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
                desc={intl.formatMessage({
                  id: 'overview.spaceOverview.tab3.noAuth.empty.description',
                })}
                title={intl.formatMessage({
                  id: 'overview.spaceOverview.tab3.noAuth.empty.title',
                })}
                btns={
                  <Button
                    type="primary"
                    onClick={() => {
                      history?.push('/space/experiment-result');
                    }}
                  >
                    {intl.formatMessage({
                      id: 'overview.spaceOverview.tab3.noAuth.empty.btn',
                    })}
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
            pageSize: experimentResultData?.page_size,
            current: experimentResultData?.page,
          }}
          onChange={(pagination: any) => {
            const { current, pageSize } = pagination;
            handleExperimentResultList({
              pageSize: pageSize,
              pageIndex: current,
            });
          }}
        />
      </Card>
    );
  };

  const items = [
    {
      label: intl.formatMessage({
        id: 'overview.spaceOverview.tab1.title',
      }),
      key: 'update',
      children: <EditExperimental />,
    },
    {
      label: intl.formatMessage({
        id: 'overview.spaceOverview.tab2.title',
      }),
      key: 'running',
      children: <EditExperimental />,
    },
    {
      label: intl.formatMessage({
        id: 'overview.spaceOverview.tab3.title',
      }),
      key: 'runningResult',
      children: <RecentlyRunExperimentalResult />,
    },
  ];

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
          <span className="title">
            {intl.formatMessage({ id: 'overview.spaceOverview' })}
          </span>
          <Select
            onChange={(val) => {
              setCurTime(timeType);
              const recent_day = val === '30day' ? 30 : 7;
              getSpaceOverview?.run({
                spaceId: history?.location?.query?.spaceId as string,
                recent_day,
              });
              handleTabSearch({ curTime: val });
            }}
            defaultValue={'7day'}
            options={[
              {
                label: intl.formatMessage({
                  id: 'overview.statistics.option.7',
                }),
                value: '7day',
              },
              {
                label: intl.formatMessage({
                  id: 'overview.statistics.option.30',
                }),
                value: '30day',
              },
            ]}
          />
        </div>
        <Spin spinning={getSpaceOverview?.loading}>
          <div className="result">
            <Row gutter={16} style={{ alignContent: 'center' }}>
              <Col span={8}>
                <Card style={{ display: 'flex', alignContent: 'center' }}>
                  <div>
                    <div className="shallow-65">
                      {intl.formatMessage({
                        id: 'overview.statistics.newExperiment',
                      })}
                    </div>
                    <span className="count">
                      {overViewInfo?.total_experiments || 0}
                    </span>
                    <span className="unit">
                      {intl.formatMessage({
                        id: 'overview.statistics.count',
                      })}
                    </span>
                  </div>
                </Card>
              </Col>
              <Col span={16}>
                <Card>
                  <Row>
                    <Col span={12}>
                      <div style={{ position: 'relative' }}>
                        <div className="shallow-65">
                          {intl.formatMessage({
                            id: 'overview.statistics.performingExperiments',
                          })}
                        </div>
                        <span className="count">
                          {overViewInfo?.total_experiment_instances || 0}
                        </span>
                        <span className="unit">
                          {intl.formatMessage({
                            id: 'overview.statistics.times',
                          })}
                        </span>
                      </div>
                    </Col>
                    <Col span={12}>
                      <div>
                        <div className="shallow-65">
                          {intl.formatMessage({
                            id: 'overview.statistics.executionFailed',
                          })}
                        </div>
                        <span className="count-error">
                          {overViewInfo?.failed_experiment_instances || 0}
                        </span>
                        <span className="unit">
                          {intl.formatMessage({
                            id: 'overview.statistics.times',
                          })}
                        </span>
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
          handleTabSearch({ key: val });
        }}
      />
    </>
  );
};
