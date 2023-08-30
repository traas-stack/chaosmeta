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
} from '@/utils/format';
import { renderTags } from '@/utils/renderItem';
import { useParamChange } from '@/utils/useParamChange';
import {
  EllipsisOutlined,
  ExclamationCircleFilled,
  QuestionCircleOutlined,
} from '@ant-design/icons';
import { history, useModel, useRequest } from '@umijs/max';
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

  // 时间类型
  const timeTypes = [
    // 最近实验时间后端暂没有提供该字段，暂时使用update_time -- todo
    {
      value: 'update_time',
      label: '最近实验时间',
    },
    {
      value: 'update_time',
      label: '最近编辑时间',
    },
    {
      value: 'next_exec',
      label: '即将运行时间',
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
        message.success('删除成功！');
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
        message.success('复制成功');
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
        title: '确认要删除这个实验吗？',
        icon: <ExclamationCircleFilled />,
        content: '删除实验将会删除该实验的配置，但不会删除历史实验结果！',
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

  // 触发方式
  const triggerMode = [
    {
      text: '手动',
      value: 'manual',
    },
    {
      text: '单次定时',
      value: 'once',
    },
    {
      text: '周期性',
      value: 'cron',
    },
  ];

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
          复制
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
          删除
        </div>
      ),
    },
  ];

  const columns: ColumnsType<any> = [
    {
      title: '名称',
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
      title: '创建人',
      width: 80,
      dataIndex: 'creator_name',
    },
    {
      title: '实验次数',
      width: 120,
      dataIndex: 'number',
      sorter: true,
      render: (text: number, record: { uuid: string }) => {
        if (text > 0) {
          return (
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
          );
        }
        return 0;
      },
    },
    {
      title: (
        <>
          <span>最近实验时间</span>
          <Tooltip title="最近一次实验开始的时间">
            <QuestionCircleOutlined />
          </Tooltip>
          /<span>状态</span>
        </>
      ),
      width: 180,
      sorter: true,
      render: (record: any) => {
        const statusTemp = experimentStatus.filter(
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
                title={statusTemp?.label}
                overlayStyle={{ maxWidth: '80px' }}
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
      title: '触发方式',
      width: 120,
      dataIndex: 'schedule_type',
      // filters: triggerMode,
      render: (text: string, record: { schedule_rule: string }) => {
        const value = triggerTypes?.filter((item) => item?.value === text)[0]
          ?.label;

        if (text === 'cron') {
          return (
            <div>
              <div>{value}</div>
              <span className="cycle">
                {cronTranstionCN(record?.schedule_rule)}
              </span>
            </div>
          );
        }
        return (
          <div>
            <div>{value}</div>
          </div>
        );
      },
    },
    {
      title: '即将运行时间',
      width: 140,
      dataIndex: 'next_exec',
      sorter: true,
      render: (text: string) => {
        return formatTime(text) || '-';
      },
    },
    {
      title: '最近编辑时间',
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
      title: '操作',
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
              编辑
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
              <Form.Item name={'name'} label="名称">
                <Input placeholder="请输入" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name={'creator'} label="创建人">
                <Input placeholder="请输入" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name={'schedule_type'} label="触发方式">
                <Select placeholder="请选择">
                  {triggerMode?.map((item) => {
                    return (
                      <Select.Option key={item.value} value={item.value}>
                        {item.text}
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
            <div className="title">实验列表</div>
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
                  创建实验
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
                      ? '请先创建实验，您可选择自己创建实验也可以通过推荐实验来快速构建实验场景，来验证应用系统的可靠性。'
                      : '您在该空间是只读权限，暂不支持创建实验。若想创建实验请去成员管理中找空间内有读写权限的成员修改权限'
                  }
                  topTitle="当前空间还没有实验数据"
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
                          创建实验
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
                          前往成员管理
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
            // dataSource={ [{uuid: 1}]}
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
