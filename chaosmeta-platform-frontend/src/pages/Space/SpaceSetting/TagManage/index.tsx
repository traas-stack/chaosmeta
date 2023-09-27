/**
 * 成员管理tab
 */
import { LightArea } from '@/components/CommonStyle';
import EmptyCustom from '@/components/EmptyCustom';
import ShowText from '@/components/ShowText';
import { tagColors } from '@/constants';
import {
  querySpaceTagList,
  spaceDeleteTag,
} from '@/services/chaosmeta/SpaceController';
import { useParamChange } from '@/utils/useParamChange';
import { SearchOutlined } from '@ant-design/icons';
import { history, useIntl, useModel, useRequest } from '@umijs/max';
import {
  Alert,
  Button,
  Form,
  Input,
  Popconfirm,
  Space,
  Table,
  Tag,
  message,
} from 'antd';
import { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';
import AddTagDrawer from './AddTagDrawer';

interface PageData {
  page: number;
  pageSize: number;
  total: number;
  labels: any[];
}

const TagManage: React.FC<any> = () => {
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const { spacePermission } = useModel('global');
  const [pageData, setPageData] = useState<PageData>({
    page: 1,
    pageSize: 10,
    total: 0,
    labels: [],
  });
  const spaceIdChange = useParamChange('spaceId');
  const [batchState, setBatchState] = useState<boolean>(false);
  const [form] = Form.useForm();
  const [addTagDrawerOpen, setAddTagDrawerOpen] = useState<boolean>(false);
  const intl = useIntl();
  /**
   * 分页接口
   */
  const getTagList = useRequest(querySpaceTagList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        setPageData(res.data);
      }
    },
  });

  /**
   * 分页数据获取
   * @param values
   */
  const handlePageSearch = (values?: {
    page?: number;
    pageSize?: number;
    sort?: string;
  }) => {
    const { page, pageSize, sort } = values || {};
    const name = form.getFieldValue('name');
    const params = {
      id: Number(spaceIdChange),
      page: page || pageData.page || 1,
      page_size: pageSize || pageData.pageSize || 10,
      sort,
      name,
    };
    getTagList?.run(params);
  };

  /**
   * 删除标签接口
   */
  const handleDelete = useRequest(spaceDeleteTag, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success(
          intl.formatMessage({ id: 'tagManageMent.delete.success.tip' }),
        );
        handlePageSearch();
      }
    },
  });

  const baseColumns: ColumnsType<any> = [
    {
      title: intl.formatMessage({ id: 'name' }),
      width: 120,
      dataIndex: 'name',
    },
    {
      title: intl.formatMessage({ id: 'tagManageMent.column.tagColor' }),
      width: 100,
      dataIndex: 'color',
      render: (text: string) => {
        const temp = tagColors?.filter((item) => item.type === text)[0];
        return (
          <Tag style={{ width: '20px', height: '20px' }} color={temp.type} />
        );
      },
    },
    {
      title: intl.formatMessage({ id: 'tagManageMent.column.tagStyle' }),
      width: 100,
      dataIndex: 'name',
      render: (text: string, record: { color: string }) => {
        const temp = tagColors?.filter(
          (item) => item.type === record?.color,
        )[0];
        return <Tag color={temp?.type}>{text}</Tag>;
      },
    },
    {
      title: intl.formatMessage({ id: 'creator' }),
      width: 160,
      dataIndex: 'creator',
    },
    {
      title: intl.formatMessage({ id: 'createTime' }),
      width: 160,
      dataIndex: 'create_time',
      sorter: true,
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
  ];

  const operateColumns: ColumnsType<any> = [
    {
      title: intl.formatMessage({ id: 'operate' }),
      width: 60,
      fixed: 'right',
      dataIndex: 'id',
      render: (text: number) => {
        return (
          <Popconfirm
            placement="bottom"
            title={intl.formatMessage({ id: 'deleteConfirmText' })}
            onConfirm={() => {
              return handleDelete?.run({
                id: text,
                ns_id: history?.location?.query?.spaceId as string,
              });
            }}
            okText={intl.formatMessage({ id: 'confirm' })}
            cancelText={intl.formatMessage({ id: 'cancel' })}
          >
            <a>{intl.formatMessage({ id: 'delete' })}</a>
          </Popconfirm>
        );
      },
    },
  ];

  useEffect(() => {
    handlePageSearch();
  }, [spaceIdChange]);

  return (
    <LightArea>
      <div className="area-operate">
        <div className="title">
          {intl.formatMessage({ id: 'tagManageMent.title' })}
        </div>
        <Form form={form}>
          <Space>
            <Form.Item name={'name'} noStyle>
              <Input
                placeholder={intl.formatMessage({
                  id: 'tagManageMent.search.placeholder',
                })}
                onPressEnter={() => {
                  handlePageSearch();
                }}
                suffix={
                  <SearchOutlined
                    onClick={() => {
                      handlePageSearch();
                    }}
                  />
                }
              />
            </Form.Item>
            {spacePermission === 1 && (
              <>
                <Button
                  className={batchState ? 'batch-active' : ''}
                  onClick={() => {
                    setBatchState(!batchState);
                  }}
                >
                  {intl.formatMessage({ id: 'batchOperate' })}
                </Button>
                <Button
                  type="primary"
                  onClick={() => {
                    setAddTagDrawerOpen(true);
                  }}
                >
                  {intl.formatMessage({ id: 'tagManageMent.add.title' })}
                </Button>
              </>
            )}
          </Space>
        </Form>
      </div>
      <div className="area-content">
        {(selectedRowKeys?.length > 0 || batchState) && (
          <Alert
            message={
              <>
                {intl.formatMessage({ id: 'selected' })}{' '}
                <a>{selectedRowKeys.length}</a>{' '}
                {intl.formatMessage({ id: 'item' })}
                <a
                  style={{ paddingLeft: '16px' }}
                  onClick={() => {
                    setSelectedRowKeys([]);
                  }}
                >
                  {intl.formatMessage({ id: 'clear' })}
                </a>
              </>
            }
            style={{ marginBottom: '16px' }}
            type="info"
            action={
              <Popconfirm
                placement="bottom"
                title={intl.formatMessage({ id: 'deleteConfirmText' })}
                disabled={selectedRowKeys?.length === 0}
                onConfirm={async () => {
                  const queryList: any[] = [];
                  // 批量删除需循环调用
                  selectedRowKeys?.forEach((item: any) => {
                    const curPromise = spaceDeleteTag({
                      id: item,
                      ns_id: history?.location?.query?.spaceId as string,
                    });
                    queryList.push(curPromise);
                  });
                  const result = await Promise.all(queryList);
                  if (result.every((item) => item.code === 200)) {
                    message.success(
                      intl.formatMessage({
                        id: 'tagManageMent.delete.success.tip',
                      }),
                    );
                    setSelectedRowKeys([]);
                    handlePageSearch();
                  }
                  return result;
                }}
                okText={intl.formatMessage({ id: 'confirm' })}
                cancelText={intl.formatMessage({ id: 'cancel' })}
              >
                <Button
                  type="link"
                  disabled={selectedRowKeys?.length === 0}
                  style={{
                    color:
                      selectedRowKeys?.length > 0
                        ? 'rgba(255,77,79,1)'
                        : 'rgba(0,10,26,0.26)',
                  }}
                >
                  {intl.formatMessage({ id: 'batchDelete' })}
                </Button>
              </Popconfirm>
            }
            showIcon
          />
        )}
        <Table
          locale={{
            emptyText: () => {
              if (spacePermission === 1) {
                return (
                  <EmptyCustom
                    desc={intl.formatMessage({
                      id: 'tagManageMent.table.empty.description',
                    })}
                    topTitle={intl.formatMessage({
                      id: 'tagManageMent.table.empty.title',
                    })}
                    btns={
                      <Button
                        type="primary"
                        onClick={() => {
                          setAddTagDrawerOpen(true);
                        }}
                      >
                        {intl.formatMessage({ id: 'tagManageMent.add.title' })}
                      </Button>
                    }
                    imgSrc="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*q5JMSKh4woYAAAAAAAAAAAAADmKmAQ/original"
                  />
                );
              }
              return (
                <EmptyCustom
                  desc={intl.formatMessage({
                    id: 'tagManageMent.table.empty.noAuth.description',
                  })}
                  topTitle={intl.formatMessage({
                    id: 'tagManageMent.table.empty.noAuth.title',
                  })}
                  btns={
                    <Button
                      type="primary"
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
                      {intl.formatMessage({ id: 'goToMemberManagement' })}
                    </Button>
                  }
                  imgSrc="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*q5JMSKh4woYAAAAAAAAAAAAADmKmAQ/original"
                />
              );
            },
          }}
          columns={
            spacePermission === 1
              ? [...baseColumns, ...operateColumns]
              : baseColumns
          }
          rowKey={'id'}
          loading={getTagList.loading}
          // dataSource={[]}
          dataSource={pageData.labels}
          scroll={{ x: 'max-content' }}
          rowSelection={
            spacePermission === 1
              ? {
                  selectedRowKeys,
                  onChange: (rowKeys: any[]) => {
                    setSelectedRowKeys(rowKeys);
                  },
                }
              : undefined
          }
          pagination={{
            showQuickJumper: true,
            total: pageData.total,
            showSizeChanger: true,
          }}
          onChange={(pagination: any, filters, sorter: any) => {
            const { current, pageSize } = pagination;
            let sort;
            if (sorter.order) {
              sort = sorter.order === 'ascend' ? 'create_time' : '-create_time';
            }
            handlePageSearch({
              pageSize: pageSize,
              page: current,
              sort,
            });
          }}
        />
      </div>
      {addTagDrawerOpen && (
        <AddTagDrawer
          open={addTagDrawerOpen}
          setOpen={setAddTagDrawerOpen}
          spaceId={history.location.query?.spaceId as string}
          handlePageSearch={handlePageSearch}
        />
      )}
    </LightArea>
  );
};

export default React.memo(TagManage);
