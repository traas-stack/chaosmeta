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
import { formatTime } from '@/utils/format';
import { useParamChange } from '@/utils/useParamChange';
import { PageContainer } from '@ant-design/pro-components';
import { history, useModel, useRequest } from '@umijs/max';
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

  // 时间类型
  const timeTypes = [
    {
      value: 'create_time',
      label: '实验开始时间',
    },
    {
      value: 'update_time',
      label: '实验结束时间',
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
    const sort = params?.sort || '-update_time';
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
          message.success(`${params[0]?.name}实验已停止`);
          handleSearch();
        }
      },
    },
  );

  const columns: ColumnsType<any> = [
    {
      title: '名称',
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
      title: '执行人',
      width: 120,
      dataIndex: 'creator_name',
    },
    {
      title: '实验开始时间',
      width: 180,
      dataIndex: 'create_time',
      sorter: true,
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
    {
      title: '实验结束时间',
      width: 180,
      dataIndex: 'update_time',
      sorter: true,
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
    {
      title: '实验状态',
      width: 100,
      dataIndex: 'status',
      render: (text: string) => {
        const temp = experimentResultStatus?.filter(
          (item) => item?.value === text,
        )[0];
        if (temp) {
          return (
            <div>
              <Badge color={temp?.color} /> {temp?.label}
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
      title: '操作',
      width: 60,
      fixed: 'right',
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

  useEffect(() => {
    handleSearch();
  }, [spaceIdChange]);

  return (
    <>
      <PageContainer title="实验结果">
        <Container>
          <div className="result-list ">
            <LightArea className="search">
              <Form form={form} labelCol={{ span: 6 }}>
                <Row gutter={16}>
                  <Col span={8}>
                    <Form.Item name={'name'} label="名称">
                      <Input placeholder="请输入" />
                    </Form.Item>
                  </Col>
                  <Col span={8}>
                    <Form.Item name={'creator_name'} label="执行人">
                      <Input placeholder="请输入" />
                    </Form.Item>
                  </Col>
                  <Col span={8}>
                    <Form.Item name={'status'} label="实验状态">
                      <Select placeholder="请选择">
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
                              spaceId: history?.location?.query?.spaceId as string,
                            }
                          })
                          handleSearch({ page: 1, pageSize: 10 });
                        }}
                      >
                        重置
                      </Button>
                      <Button
                        type="primary"
                        onClick={() => {
                          handleSearch();
                        }}
                      >
                        查询
                      </Button>
                    </Space>
                  </Col>
                </Row>
              </Form>
            </LightArea>
            <LightArea>
              <div className="table">
                <div className="area-operate">
                  <div className="title">实验结果列表</div>
                </div>
                <Table
                  locale={{
                    emptyText: (
                      <EmptyCustom
                        desc="当前暂无实验结果数据"
                        title={
                          spacePermission === 1
                            ? '您可以前往实验详情页面运行实验'
                            : '您可以前往实验列表页面查看实验'
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
                              前往实验列表
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
