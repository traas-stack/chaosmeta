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
import { history, useModel, useRequest } from '@umijs/max';
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
        message.success('你已成功删除标签');
        handlePageSearch();
      }
    },
  });

  const baseColumns: ColumnsType<any> = [
    {
      title: '名称',
      width: 120,
      dataIndex: 'name',
    },
    {
      title: '标签颜色',
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
      title: '标签样式',
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
      title: '创建人',
      width: 160,
      dataIndex: 'creator',
    },
    {
      title: '创建时间',
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
      title: '操作',
      width: 60,
      fixed: 'right',
      dataIndex: 'id',
      render: (text: number) => {
        return (
          <Popconfirm
            placement="bottom"
            title={'你确定要删除吗？'}
            onConfirm={() => {
              return handleDelete?.run({
                id: text,
                ns_id: history?.location?.query?.spaceId as string,
              });
            }}
            okText="确定"
            cancelText="取消"
          >
            <a>删除</a>
          </Popconfirm>
        );
      },
    },
  ];

  useEffect(() => {
    handlePageSearch();
  }, []);
  return (
    <LightArea>
      <div className="area-operate">
        <div className="title">标签列表</div>
        <Form form={form}>
          <Space>
            <Form.Item name={'name'} noStyle>
              <Input
                placeholder="请输入标签名称"
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
                  批量操作
                </Button>
                <Button
                  type="primary"
                  onClick={() => {
                    setAddTagDrawerOpen(true);
                  }}
                >
                  新建标签
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
                已选择 <a>{selectedRowKeys.length}</a> 项
                <a
                  style={{ paddingLeft: '16px' }}
                  onClick={() => {
                    setSelectedRowKeys([]);
                  }}
                >
                  清空
                </a>
              </>
            }
            style={{ marginBottom: '16px' }}
            type="info"
            action={
              <Popconfirm
                placement="bottom"
                title={'你确定要删除吗？'}
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
                    message.success('你已成功删除标签');
                    setSelectedRowKeys([]);
                    handlePageSearch();
                  }
                  return result;
                }}
                okText="确定"
                cancelText="取消"
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
                  批量删除
                </Button>
              </Popconfirm>
            }
            showIcon
          />
        )}
        <Table
          locale={{
            emptyText: () => {
              if (spacePermission === 0) {
                return (
                  <EmptyCustom
                    desc="请新建标签，提前创建好标签在创建实验的时候可以直接选用快速为实验打上标签。"
                    topTitle="您还没有标签数据"
                    btns={
                      <Button
                        type="primary"
                        onClick={() => {
                          setAddTagDrawerOpen(true);
                        }}
                      >
                        新建标签
                      </Button>
                    }
                    imgSrc="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*q5JMSKh4woYAAAAAAAAAAAAADmKmAQ/original"
                  />
                );
              }
              return (
                <EmptyCustom
                  desc="您在该空间只是只读权限，暂不支持添加标签。若想添加标签请去成员管理中找空间内有读写权限的成员修改权限"
                  topTitle="当前还没有标签数据"
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
                      前往成员管理
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
