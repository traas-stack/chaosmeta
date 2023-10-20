/**
 * 实验结果页面
 * @returns
 */
import { LightArea } from '@/components/CommonStyle';
import EmptyCustom from '@/components/EmptyCustom';
import TimeTypeRangeSelect from '@/components/Select/TimeTypeRangeSelect';
import ShowText from '@/components/ShowText';
import { experimentResultStatus } from '@/constants';
import {
  queryExperimentResultList,
  stopExperimentResult,
} from '@/services/chaosmeta/ExperimentController';
import { formatTime, getIntlLabel } from '@/utils/format';
import { useParamChange } from '@/utils/useParamChange';
import { PageContainer } from '@ant-design/pro-components';
import { history, useIntl, useModel, useRequest } from '@umijs/max';
import {
  Badge,
  Button,
  Col,
  Form,
  Input,
  Popconfirm,
  Row,
  Select,
  Space,
  Table,
  message,
} from 'antd';
import { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';
import { Container } from './style';

interface PageData {
  results: any[];
  page: number;
  pageSize: number;
  total: number;
}
const ExperimentResult: React.FC<unknown> = () => {
  const [form] = Form.useForm();
  const [pageData, setPageData] = useState<PageData>({
    page: 1,
    pageSize: 10,
    total: 0,
    results: [],
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
    {
      value: 'create_time',
      label: intl.formatMessage({
        id: 'experimentStartTime',
      }),
    },
    {
      value: 'update_time',
      label: intl.formatMessage({
        id: 'experimentEndTime',
      }),
    },
  ];

  /**
   * 分页接口
   */
  const queryByPage = useRequest(queryExperimentResultList, {
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
    const { name, status, creator_name, timeType, time } =
      form.getFieldsValue();
    const page = params?.page || pageData.page || 1;
    const pageSize = params?.pageSize || pageData.pageSize || 10;
    let start_time, end_time;
    if (time?.length > 0) {
      start_time = formatTime(time[0]?.format());
      end_time = formatTime(time[1]?.format());
    }
    // 默认使用更新时间倒序
    const sort = params?.sort || '-create_time';
    const queryParam = {
      sort,
      name,
      page,
      page_size: pageSize,
      experiment_uuid: history?.location?.query?.experimentId as string,
      creator_name,
      status,
      time_search_field: timeType,
      start_time,
      end_time,
      namespace_id: history?.location?.query?.spaceId as string,
      time_type: end_time ? 'range' : undefined,
    };
    queryByPage.run(queryParam);
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
          message.success(
            `${params[0]?.name} ${intl.formatMessage({
              id: 'experimentResult.stop.text',
            })}`,
          );
          handleSearch();
        }
      },
    },
  );

  const columns: ColumnsType<any> = [
    {
      title: intl.formatMessage({
        id: 'name',
      }),
      width: 160,
      dataIndex: 'name',
      render: (text: string, record: { uuid: string }) => {
        return (
          <a
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
          </a>
        );
      },
    },
    {
      title: intl.formatMessage({
        id: 'creator',
      }),
      width: 120,
      dataIndex: 'creator_name',
    },
    {
      title: intl.formatMessage({
        id: 'experimentStartTime',
      }),
      width: 180,
      dataIndex: 'create_time',
      sorter: true,
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
    {
      title: intl.formatMessage({
        id: 'experimentEndTime',
      }),
      width: 180,
      dataIndex: 'update_time',
      sorter: true,
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
    {
      title: intl.formatMessage({
        id: 'experimentStatus',
      }),
      width: 100,
      dataIndex: 'status',
      render: (text: string) => {
        const temp: any = experimentResultStatus?.filter(
          (item) => item?.value === text,
        )[0];
        if (temp) {
          return (
            <div>
              <Badge color={temp?.color} /> {getIntlLabel(temp)}
            </div>
          );
        }
        return '-';
      },
    },
  ];

  // 操作column，有权限才可查看
  const operateColumn: any[] = [
    {
      title: intl.formatMessage({
        id: 'operate',
      }),
      width: 60,
      fixed: 'right',
      render: (record: { uuid: string; status: string; name: string }) => {
        // 运行中才可以停止
        if (record?.status === 'Running') {
          return (
            <Space>
              <Popconfirm
                title={intl.formatMessage({
                  id: 'stopConfirmText',
                })}
                onConfirm={() => {
                  handleStopExperiment?.run({
                    uuid: record?.uuid,
                    name: record?.name,
                  });
                }}
              >
                <a>
                  {intl.formatMessage({
                    id: 'stop',
                  })}
                </a>
              </Popconfirm>
            </Space>
          );
        }
        return null;
      },
    },
  ];

  useEffect(() => {
    handleSearch();
  }, [spaceIdChange]);

  return (
    <>
      <PageContainer
        title={intl.formatMessage({
          id: 'experimentResult',
        })}
      >
        <Container>
          <div className="result-list ">
            <LightArea className="search">
              <Form form={form} labelCol={{ span: 6 }}>
                <Row gutter={16}>
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
                      name={'creator_name'}
                      label={intl.formatMessage({
                        id: 'creator',
                      })}
                    >
                      <Input placeholder={pleaseInput} />
                    </Form.Item>
                  </Col>
                  <Col span={8}>
                    <Form.Item
                      name={'status'}
                      label={intl.formatMessage({
                        id: 'experimentStatus',
                      })}
                    >
                      <Select placeholder={pleaseSelect}>
                        {experimentResultStatus.map((item) => {
                          return (
                            <Select.Option key={item.value} value={item.value}>
                              {item.label}
                            </Select.Option>
                          );
                        })}
                      </Select>
                    </Form.Item>
                  </Col>
                  <TimeTypeRangeSelect form={form} timeTypes={timeTypes} />
                  <Col style={{ textAlign: 'right', flex: 1 }} span={16}>
                    <Space>
                      <Button
                        onClick={() => {
                          form.resetFields();
                          history.push({
                            pathname: '/space/experiment-result',
                            query: {
                              spaceId: history?.location?.query
                                ?.spaceId as string,
                            },
                          });
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
                    {' '}
                    {intl.formatMessage({
                      id: 'experimentResultList',
                    })}
                  </div>
                </div>
                <Table
                  locale={{
                    emptyText: (
                      <EmptyCustom
                        desc={intl.formatMessage({
                          id: 'experimentResult.table.description',
                        })}
                        title={
                          spacePermission === 1
                            ? intl.formatMessage({
                                id: 'experimentResult.table.title',
                              })
                            : intl.formatMessage({
                                id: 'experimentResult.table.noAuth.title',
                              })
                        }
                        btns={
                          <Space>
                            <Button
                              type="primary"
                              onClick={() => {
                                history.push({
                                  pathname: '/space/experiment',
                                  query: {
                                    spaceId: history?.location?.query?.spaceId,
                                  },
                                });
                              }}
                            >
                              {intl.formatMessage({
                                id: 'experimentResult.table.btn',
                              })}
                            </Button>
                          </Space>
                        }
                      />
                    ),
                  }}
                  columns={
                    spacePermission === 1
                      ? [...columns, ...operateColumn]
                      : columns
                  }
                  loading={queryByPage?.loading}
                  rowKey={'uuid'}
                  scroll={{ x: 1000 }}
                  dataSource={pageData?.results || []}
                  pagination={
                    pageData?.results?.length > 0
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
                      sort =
                        sorter.order === 'ascend' ? sortKey : `-${sortKey}`;
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
        </Container>
      </PageContainer>
    </>
  );
};

export default ExperimentResult;
