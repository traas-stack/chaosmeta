/**
 * 成员管理tab
 */
import { LightArea } from '@/components/CommonStyle';
import ShowText from '@/components/ShowText';
import { Role } from '@/pages/GlobalSetting/Account/style';
import {
  querySpaceUserList,
  spaceBatchDeleteUser,
  spaceModifyUserPermission,
} from '@/services/chaosmeta/SpaceController';
import { useParamChange } from '@/utils/useParamChange';
import { ExclamationCircleFilled, SearchOutlined } from '@ant-design/icons';
import { history, useModel, useRequest } from '@umijs/max';
import {
  Alert,
  Button,
  Form,
  Input,
  Modal,
  Select,
  Space,
  Table,
  message,
} from 'antd';
import { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';
import AddMemberModal from './AddMemberModal';

interface DataType {
  id: string;
  account: string;
  auth: string;
}

interface PageData {
  page: number;
  pageSize: number;
  total: number;
  users: any[];
}
const MemberManage: React.FC<unknown> = () => {
  const [form] = Form.useForm();
  const [selectedRowKeys, setSelectedRowKeys] = useState<number[]>([]);
  const [addMemberOpen, setAddMemberOpen] = useState<boolean>(false);
  const spaceIdChange = useParamChange('spaceId');
  const [batchState, setBatchState] = useState<boolean>(false);
  const [pageData, setPageData] = useState<PageData>({
    page: 1,
    pageSize: 10,
    total: 0,
    users: [],
  });
  const { spacePermission } = useModel('global');

  /**
   * 分页接口
   */
  const getUserList = useRequest(querySpaceUserList, {
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
  const handlePageSearch = (values?: any) => {
    const { page, pageSize, sort, permission } = values || {};
    const name = form.getFieldValue('name');
    const params = {
      id: Number(spaceIdChange),
      page: page || pageData.page || 1,
      page_size: pageSize || pageData.pageSize || 10,
      sort,
      username: name,
      permission,
    };
    getUserList?.run(params);
  };

  /**
   * 删除成员
   */
  const deleteUser = useRequest(spaceBatchDeleteUser, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res.code === 200) {
        message.success('您已成功删除所选成员');
        handlePageSearch();
      }
    },
  });

  /**
   * 修改空间内成员权限
   */
  const editUserRole = useRequest(spaceModifyUserPermission, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res.code === 200) {
        message.success('权限修改成功');
        handlePageSearch();
      }
    },
  });

  /**
   * 删除成员
   */
  const handleDeleteUser = (id?: number) => {
    const user_ids = id || id === 0 ? [id] : selectedRowKeys;
    Modal.confirm({
      title: '您确认要删除当前所选成员吗？',
      icon: <ExclamationCircleFilled />,
      content: '删除空间内成员，该成员将无法进入该空间！',
      onOk() {
        return deleteUser?.run({
          user_ids,
          id: Number(history?.location.query?.spaceId),
        });
      },
    });
  };

  const authOptions = [
    {
      label: (
        <Role>
          <span>读写</span>
          <span>在当前空间有编辑权限</span>
        </Role>
      ),
      value: '1',
      name: '读写',
    },
    {
      label: (
        <Role>
          <span>只读</span>
          <span>在当前空间只能查看</span>
        </Role>
      ),
      value: '0',
      name: '只读',
    },
  ];

  const baseColumns: ColumnsType<DataType> = [
    {
      title: '用户名',
      width: 80,
      dataIndex: 'name',
    },
    {
      title: '权限',
      width: 160,
      dataIndex: 'permission',
      filters: [
        {
          text: '读写',
          value: '1',
        },
        {
          text: '只读',
          value: '0',
        },
      ],
      render: (text: string, record: any) => {
        if (spacePermission === 1) {
          return (
            <Select
              dropdownStyle={{ minWidth: '200px' }}
              value={text}
              style={{ minWidth: '80px', width: '100%' }}
              optionLabelProp="label"
              onChange={(val) => {
                editUserRole?.run({
                  user_ids: [record?.id],
                  permission: Number(val),
                  id: history?.location?.query?.spaceId as string,
                });
              }}
            >
              {authOptions?.map((item) => {
                return (
                  <Select.Option
                    value={item?.value}
                    label={item.name}
                    key={item.value}
                  >
                    {item.label}
                  </Select.Option>
                );
              })}
            </Select>
          );
        }
        return (
          authOptions?.filter((item) => item.value === text)[0]?.name || text
        );
      },
    },
    {
      title: '加入时间',
      width: 160,
      dataIndex: 'create_time',
      sorter: true,
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
  ];

  const operateColumns: ColumnsType<DataType> = [
    {
      title: '操作',
      width: 60,
      fixed: 'right',
      dataIndex: 'id',
      render: (text: number) => {
        return (
          <a
            onClick={() => {
              handleDeleteUser(text);
            }}
          >
            删除
          </a>
        );
      },
    },
  ];

  useEffect(() => {
    getUserList?.run({ id: Number(spaceIdChange), page: 1, page_size: 10 });
  }, [spaceIdChange]);

  return (
    <LightArea>
      <div className="area-operate">
        <div className="title">成员列表</div>
        <Form form={form}>
          <Space>
            <Form.Item name={'name'} noStyle>
              <Input
                placeholder="请输入用户名"
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
                    setAddMemberOpen(true);
                  }}
                >
                  添加成员
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
              <Button
                type="link"
                disabled={selectedRowKeys?.length === 0}
                style={{
                  color:
                    selectedRowKeys?.length > 0
                      ? 'rgba(255,77,79,1)'
                      : 'rgba(0,10,26,0.26)',
                }}
                onClick={() => {
                  setSelectedRowKeys([]);
                  handleDeleteUser();
                }}
              >
                批量删除
              </Button>
            }
            showIcon
          />
        )}
        <Table
          columns={
            spacePermission === 1
              ? [...baseColumns, ...operateColumns]
              : baseColumns
          }
          rowKey={'id'}
          loading={getUserList.loading}
          dataSource={pageData.users}
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
          }}
          onChange={(pagination: any, filters, sorter: any) => {
            const { current, pageSize } = pagination;
            let sort;
            if (sorter.order) {
              sort = sorter.order === 'ascend' ? 'create_time' : '-create_time';
            }
            let permission;
            if (filters?.permission?.length === 1) {
              permission = filters.permission[0];
            }
            handlePageSearch({
              pageSize: pageSize,
              page: current,
              sort,
              permission,
            });
          }}
        />
      </div>
      {addMemberOpen && (
        <AddMemberModal
          setOpen={setAddMemberOpen}
          open={addMemberOpen}
          handlePageSearch={handlePageSearch}
          spaceId={history?.location?.query?.spaceId as string}
        />
      )}
    </LightArea>
  );
};

export default React.memo(MemberManage);
