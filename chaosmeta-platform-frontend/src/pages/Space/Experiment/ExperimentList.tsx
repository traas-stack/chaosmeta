import { LightArea } from '@/components/CommonStyle';
import EmptyCustom from '@/components/EmptyCustom';
import TimeTypeRangeSelect from '@/components/Select/TimeTypeRangeSelect';
import { experimentStatus, triggerTypes } from '@/constants';
import {
  createExperiment,
  deleteExperiment,
  queryExperimentList,
} from '@/services/chaosmeta/ExperimentController';
import {
  copyExperimentFormatData,
  cronTranstionCN,
  formatTime,
  getIntlLabel,
} from '@/utils/format';
import { renderTags } from '@/utils/renderItem';
import { useParamChange } from '@/utils/useParamChange';
import {
  EllipsisOutlined,
  ExclamationCircleFilled,
  QuestionCircleOutlined,
} from '@ant-design/icons';
import { history, useIntl, useModel, useRequest } from '@umijs/max';
import {
  Badge,
  Button,
  Col,
  Dropdown,
  Form,
  Input,
  Modal,
  Popover,
  Row,
  Select,
  Space,
  Table,
  Tooltip,
  message,
} from 'antd';
import { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';

interface PageData {
  experiments: any[];
  page: number;
  pageSize: number;
  total: number;
}
/**
 * 实验列表
 * @returns
 */
const ExperimentList: React.FC<unknown> = () => {
  const [form] = Form.useForm();
  const [pageData, setPageData] = useState<PageData>({
    page: 1,
    pageSize: 10,
    total: 0,
    experiments: [],
  });
  const spaceIdChange = useParamChange('spaceId');
  const { spacePermission } = useModel('global');
  const intl = useIntl();
  // 请输入文案
  const pleaseInput = intl.formatMessage({
    id: 'pleaseInput',
  });

  // 请选择文案
  const pleaseSelect = intl.formatMessage({
    id: 'pleaseSelect',
  });

  // 时间类型
  const timeTypes = [
    // 最近实验时间后端暂没有提供该字段，暂时使用update_time -- todo
    {
      value: 'update_time',
      label: intl.formatMessage({
        id: 'latestExperimentalTime',
      }),
    },
    {
      value: 'update_time',
      label: intl.formatMessage({
        id: 'lastEditTime',
      }),
    },
    {
      value: 'next_exec',
      label: intl.formatMessage({
        id: 'upcomingRunningTime',
      }),
    },
  ];

  /**
   * 分页接口
   */
  const queryByPage = useRequest(queryExperimentList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        setPageData(res?.data);
      }
    },
  });

  /**
   * 分页查询
   * @param params
   */
  const handleSearch = (params?: {
    sort?: string;
    page?: number;
    pageSize?: number;
  }) => {
    const { name, creator, schedule_type, time, timeType } =
      form.getFieldsValue();
    const page = params?.page || pageData.page || 1;
    const pageSize = params?.pageSize || pageData.pageSize || 10;
    let start_time, end_time;
    if (time?.length > 0) {
      start_time = formatTime(time[0]?.format());
      end_time = formatTime(time[1]?.format());
    }
    // 默认使用更新时间倒序
    const sort = params?.sort || '-update_time';
    const queryParam = {
      sort,
      name,
      page,
      page_size: pageSize,
      creator,
      schedule_type,
      start_time,
      end_time,
      time_search_field: timeType,
      namespace_id: history?.location?.query?.spaceId as string,
    };
    queryByPage.run(queryParam);
  };

  /**
   * 删除实验接口
   */
  const handleDeleteExperiment = useRequest(deleteExperiment, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success(
          intl.formatMessage({
            id: 'deleteText',
          }),
        );
        handleSearch();
      }
    },
  });

  /**
   * 创建实验 -- 复制实验使用
   */
  const handleCreateExperiment = useRequest(createExperiment, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success(
          intl.formatMessage({
            id: 'copyText',
          }),
        );
        handleSearch();
      }
    },
  });

  /**
   * 确认删除实验
   */
  const handleDeleteConfirm = (uuid: string) => {
    if (uuid) {
      Modal.confirm({
        title: intl.formatMessage({
          id: 'experiment.delete.title',
        }),
        icon: <ExclamationCircleFilled />,
        content: intl.formatMessage({
          id: 'experiment.delete.content',
        }),
        onOk() {
          handleDeleteExperiment?.run({ uuid });
          handleSearch();
        },
        onCancel() {},
      });
    }
  };

  /**
   * 复制实验
   */
  const handleCopyExperiment = (record: any) => {
    const params = copyExperimentFormatData(record);
    handleCreateExperiment?.run(params);
  };

  /**
   * 操作下拉列表
   * @param record
   * @returns
   */
  const operateItems = (record: any) => [
    {
      key: '1',
      label: (
        <div
          onClick={() => {
            handleCopyExperiment(record);
          }}
        >
          {intl.formatMessage({
            id: 'copy',
          })}
        </div>
      ),
    },
    {
      key: '2',
      label: (
        <div
          onClick={() => {
            handleDeleteConfirm(record?.uuid);
          }}
        >
          {intl.formatMessage({
            id: 'delete',
          })}
        </div>
      ),
    },
  ];

  const columns: ColumnsType<any> = [
    {
      title: intl.formatMessage({
        id: 'name',
      }),
      width: 160,
      dataIndex: 'name',
      render: (text: string, record: { uuid: string; labels: any[] }) => {
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
            <div className="ellipsis">{renderTags(record?.labels)}</div>
          </>
        );
      },
    },
    {
      title: intl.formatMessage({
        id: 'creator',
      }),
      width: 80,
      dataIndex: 'creator_name',
    },
    {
      title: intl.formatMessage({
        id: 'numberOfExperiments',
      }),
      width: 120,
      dataIndex: 'number',
      // sorter: true,
      render: (text: number, record: { uuid: string }) => {
        if (text > 0) {
          return (
            <a
              onClick={() => {
                history?.push({
                  pathname: '/space/experiment-result',
                  query: {
                    experimentId: record?.uuid,
                  },
                });
              }}
            >
              {text}
            </a>
          );
        }
        return 0;
      },
    },
    {
      title: (
        <>
          <span>
            {intl.formatMessage({
              id: 'latestExperimentalTime',
            })}
          </span>
          <Tooltip
            title={intl.formatMessage({
              id: 'lastStartTime',
            })}
          >
            <QuestionCircleOutlined />
          </Tooltip>
          /
          <span>
            {intl.formatMessage({
              id: 'status',
            })}
          </span>
        </>
      ),
      width: 180,
      sorter: true,
      render: (record: any) => {
        const statusTemp: any = experimentStatus.filter(
          (item) => item.value === record?.status,
        )[0];
        return (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
            }}
          >
            <div>
              <Popover
                title={getIntlLabel(statusTemp)}
                overlayStyle={{ textAlign: 'center' }}
              >
                <Badge color={statusTemp?.color} />{' '}
              </Popover>
              {/* 这里展示时间 -- 后端没有相应字段 --todo */}
              <span>{formatTime(record?.last_instance)}</span>
            </div>
          </div>
        );
      },
    },
    {
      title: intl.formatMessage({
        id: 'triggerMode',
      }),
      width: 120,
      dataIndex: 'schedule_type',
      render: (text: string, record: { schedule_rule: string }) => {
        const value: any = triggerTypes?.filter(
          (item) => item?.value === text,
        )[0];
        if (text === 'cron') {
          return (
            <div>
              <div>{getIntlLabel(value)}</div>
              <span className="cycle">
                {cronTranstionCN(record?.schedule_rule)}
              </span>
            </div>
          );
        }
        return (
          <div>
            <div>{getIntlLabel(value)}</div>
          </div>
        );
      },
    },
    {
      title: intl.formatMessage({
        id: 'upcomingRunningTime',
      }),
      width: 140,
      dataIndex: 'next_exec',
      sorter: true,
      render: (text: string) => {
        return formatTime(text) || '-';
      },
    },
    {
      title: intl.formatMessage({
        id: 'lastEditTime',
      }),
      width: 180,
      dataIndex: 'update_time',
      sorter: true,
      render: (text: string) => {
        return formatTime(text) || '-';
      },
    },
  ];

  // 操作项columns，管理员才展示
  const operateColumn: any[] = [
    {
      title: intl.formatMessage({
        id: 'operate',
      }),
      width: 90,
      fixed: 'right',
      render: (record: any) => {
        return (
          <Space>
            <a
              onClick={() => {
                window.open(
                  `${window.origin}/space/experiment/add?experimentId=${record?.uuid}&spaceId=${history?.location?.query?.spaceId}`,
                );
              }}
            >
              {intl.formatMessage({
                id: 'edit',
              })}
            </a>
            <Dropdown
              menu={{ items: operateItems(record) }}
              placement="bottom"
              arrow
            >
              <EllipsisOutlined className="operate-icon" />
            </Dropdown>
          </Space>
        );
      },
    },
  ];

  useEffect(() => {
    handleSearch();
  }, [spaceIdChange]);

  return (
    <div className="experiment-list ">
      <LightArea className="search">
        <Form form={form} labelCol={{ span: 6 }}>
          <Row gutter={16} style={{ whiteSpace: 'nowrap', overflow: 'hidden' }}>
            <Col span={8}>
              <Form.Item
                name={'name'}
                label={intl.formatMessage({
                  id: 'name',
                })}
              >
                <Input placeholder={pleaseInput} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name={'creator'}
                label={intl.formatMessage({
                  id: 'creator',
                })}
              >
                <Input placeholder={pleaseInput} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name={'schedule_type'}
                label={intl.formatMessage({
                  id: 'triggerMode',
                })}
              >
                <Select placeholder={pleaseSelect}>
                  {triggerTypes?.map((item: any) => {
                    return (
                      <Select.Option key={item.value} value={item.value}>
                        {getIntlLabel(item)}
                      </Select.Option>
                    );
                  })}
                </Select>
              </Form.Item>
            </Col>
            <TimeTypeRangeSelect form={form} timeTypes={timeTypes} />
            <Col style={{ textAlign: 'right', flex: 1 }}>
              <Space>
                <Button
                  onClick={() => {
                    form.resetFields();
                    handleSearch({ page: 1, pageSize: 10 });
                  }}
                >
                  {intl.formatMessage({
                    id: 'reset',
                  })}
                </Button>
                <Button
                  type="primary"
                  onClick={() => {
                    handleSearch();
                  }}
                >
                  {intl.formatMessage({
                    id: 'query',
                  })}
                </Button>
              </Space>
            </Col>
          </Row>
        </Form>
      </LightArea>
      <LightArea>
        <div className="table">
          <div className="area-operate">
            <div className="title">
              {intl.formatMessage({
                id: 'experimentList',
              })}
            </div>
            {spacePermission === 1 && (
              <Space>
                <Button
                  type="primary"
                  onClick={() => {
                    history.push({
                      pathname: '/space/experiment/add',
                      query: {
                        spaceId: history?.location?.query?.spaceId,
                      },
                    });
                  }}
                >
                  {intl.formatMessage({
                    id: 'createExperiment',
                  })}
                </Button>
              </Space>
            )}
          </div>
          <Table
            locale={{
              emptyText: (
                <EmptyCustom
                  desc={
                    spacePermission === 1
                      ? intl.formatMessage({
                          id: 'experiment.table.noAuth.description',
                        })
                      : intl.formatMessage({
                          id: 'experiment.table.description',
                        })
                  }
                  topTitle={intl.formatMessage({
                    id: 'experiment.table.title',
                  })}
                  btns={
                    <Space>
                      {spacePermission === 1 ? (
                        <Button
                          type="primary"
                          onClick={() => {
                            history.push({
                              pathname: '/space/experiment/add',
                              query: {
                                spaceId: history?.location?.query?.spaceId,
                              },
                            });
                          }}
                        >
                          {intl.formatMessage({
                            id: 'createExperiment',
                          })}
                        </Button>
                      ) : (
                        <Button
                          type="primary"
                          onClick={() => {
                            history.push({
                              pathname: '/space/setting',
                              query: {
                                spaceId: history?.location?.query?.spaceId,
                                tabKey: 'user',
                              },
                            });
                          }}
                        >
                          {intl.formatMessage({
                            id: 'goToMemberManagement',
                          })}
                        </Button>
                      )}
                      {/* <Button
                        onClick={() => {
                          history.push({
                            pathname: '/space/setting',
                            query: {
                              tabKey: 'user',
                              spaceId: history?.location?.query?.spaceId,
                            },
                          });
                        }}
                      >
                        推荐实验
                      </Button> */}
                    </Space>
                  }
                />
              ),
            }}
            columns={
              spacePermission === 1 ? [...columns, ...operateColumn] : columns
            }
            loading={queryByPage?.loading}
            rowKey={'uuid'}
            scroll={{ x: 1000 }}
            dataSource={pageData?.experiments || []}
            pagination={
              pageData?.experiments?.length > 0
                ? {
                    showQuickJumper: true,
                    total: pageData?.total,
                    current: pageData?.page,
                    pageSize: pageData?.pageSize,
                    showSizeChanger: true,
                  }
                : false
            }
            onChange={(pagination: any, filters, sorter: any) => {
              const { current, pageSize } = pagination;
              let sort;
              const sortKey = sorter?.field;
              if (sorter.order) {
                sort = sorter.order === 'ascend' ? sortKey : `-${sortKey}`;
              }
              handleSearch({
                pageSize: pageSize,
                page: current,
                sort,
              });
            }}
          />
        </div>
      </LightArea>
    </div>
  );
};

export default ExperimentList;
