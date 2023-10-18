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
import { history, useIntl, useModel, useRequest } from '@umijs/max';
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
  const intl = useIntl();

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
        message.success(
          intl.formatMessage({ id: 'memberManageMent.delete.success.tip' }),
        );
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
        message.success(
          intl.formatMessage({ id: 'memberManageMent.permission.success.tip' }),
        );
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
      title: intl.formatMessage({ id: 'memberManageMent.delete.title' }),
      icon: <ExclamationCircleFilled />,
      content: intl.formatMessage({ id: 'memberManageMent.delete.content' }),
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
          <span>{intl.formatMessage({ id: 'write' })}</span>
          <span>
            {intl.formatMessage({ id: 'memberManageMent.write.tip' })}
          </span>
        </Role>
      ),
      value: '1',
      name: intl.formatMessage({ id: 'write' }),
    },
    {
      label: (
        <Role>
          <span>{intl.formatMessage({ id: 'readonly' })}</span>
          <span>
            {intl.formatMessage({ id: 'memberManageMent.readonly.tip' })}
          </span>
        </Role>
      ),
      value: '0',
      name: intl.formatMessage({ id: 'readonly' }),
    },
  ];

  const baseColumns: ColumnsType<DataType> = [
    {
      title: intl.formatMessage({ id: 'username' }),
      width: 80,
      dataIndex: 'name',
    },
    {
      title: intl.formatMessage({ id: 'permission' }),
      width: 160,
      dataIndex: 'permission',
      filters: [
        {
          text: intl.formatMessage({ id: 'write' }),
          value: '1',
        },
        {
          text: intl.formatMessage({ id: 'readonly' }),
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
      title: intl.formatMessage({ id: 'joinTime' }),
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
      title: intl.formatMessage({ id: 'operate' }),
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
            {intl.formatMessage({ id: 'delete' })}
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
        <div className="title">{intl.formatMessage({ id: 'memberList' })}</div>
        <Form form={form}>
          <Space>
            <Form.Item name={'name'} noStyle>
              <Input
                placeholder={intl.formatMessage({ id: 'usernamePlaceholder' })}
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
                    setAddMemberOpen(true);
                  }}
                >
                  {intl.formatMessage({ id: 'addMember' })}
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
                {intl.formatMessage({ id: 'batchDelete' })}
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
            showSizeChanger: true,
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
