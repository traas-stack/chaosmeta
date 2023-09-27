import { LightArea } from '@/components/CommonStyle';
import ShowText from '@/components/ShowText';
import {
  batchDeleteUser,
  changeUserRole,
  getUserList,
} from '@/services/chaosmeta/UserController';
import { ExclamationCircleFilled, SearchOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { useIntl, useModel, useRequest } from '@umijs/max';
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
import React, { useEffect, useState } from 'react';
import { Container, Role } from './style';

interface DataType {
  id: string;
  auth?: string;
  name: string;
}

const Account: React.FC<unknown> = () => {
  const [pageData, setPageData] = useState<any>({});
  const [form] = Form.useForm();
  // 当前选中的数据
  const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([]);
  const { userInfo } = useModel('global');
  const [batchState, setBatchState] = useState<boolean>(false);
  const intl = useIntl();
  /**
   * 分页查询
   */
  const queryByPage = useRequest(getUserList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res: any) => {
      if (res?.data) {
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
    role?: any;
  }) => {
    const name = form.getFieldValue('name');
    const page = params?.page || pageData.page || 1;
    const pageSize = params?.pageSize || pageData.pageSize || 10;
    const role = params?.role?.length > 1 ? '' : params?.role?.toString();
    const queryParam = {
      name,
      ...params,
      page,
      pageSize,
      role,
    };
    queryByPage.run(queryParam);
  };

  /**
   * 修改用户角色
   */
  const handleChangeRole = useRequest(changeUserRole, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res.code === 200) {
        message.success(intl.formatMessage({ id: 'account.role.update' }));
        handleSearch();
      }
    },
  });

  /**
   * 删除账号接口
   */
  const handleBatchDelete = useRequest(batchDeleteUser, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: () => {
      message.success(intl.formatMessage({ id: 'account.delete.success' }));
      handleSearch();
    },
  });

  /**
   * 删除账号
   */
  const handleDeleteAccount = (ids: string[]) => {
    if (ids?.length > 0) {
      Modal.confirm({
        title: intl.formatMessage({ id: 'account.delete.title' }),
        icon: <ExclamationCircleFilled />,
        content: intl.formatMessage({ id: 'account.delete.content' }),
        onOk() {
          handleBatchDelete?.run({ user_ids: ids });
          handleSearch();
        },
        onCancel() {},
      });
    }
  };

  // 角色下拉选项
  const authOptions = [
    {
      label: (
        <Role>
          <span>{intl.formatMessage({ id: 'admin' })}</span>
          <span>{intl.formatMessage({ id: 'adminDescription' })}</span>
        </Role>
      ),
      value: 'admin',
      name: intl.formatMessage({ id: 'admin' }),
    },
    {
      label: (
        <Role>
          <span>{intl.formatMessage({ id: 'generalUser' })}</span>
          <span>{intl.formatMessage({ id: 'generalUserDescription' })}</span>
        </Role>
      ),
      value: 'normal',
      name: intl.formatMessage({ id: 'generalUser' }),
    },
  ];

  const columns: any = [
    {
      title: intl.formatMessage({ id: 'username' }),
      width: 80,
      dataIndex: 'name',
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
    {
      title: intl.formatMessage({ id: 'role' }),
      width: 160,
      dataIndex: 'role',
      filters: [
        {
          text: intl.formatMessage({ id: 'admin' }),
          value: 'admin',
        },
        {
          text: intl.formatMessage({ id: 'generalUser' }),
          value: 'normal',
        },
      ],
      render: (text: string, record: any) => {
        if (userInfo.role === 'admin') {
          return (
            <Select
              dropdownStyle={{ minWidth: '280px' }}
              value={text}
              style={{ minWidth: '80px', width: '100%' }}
              optionLabelProp="label"
              onChange={(value: string) => {
                handleChangeRole?.run({ user_ids: [record.id], role: value });
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
          authOptions?.filter((item) => item?.value === text)[0]?.name || text
        );
      },
    },
  ];

  const operateColumns: any = [
    {
      title: intl.formatMessage({ id: 'operate' }),
      width: 60,
      fixed: 'right',
      render: (record: DataType) => {
        return record?.name !== userInfo?.name ? (
          <a
            onClick={() => {
              handleDeleteAccount([record.id]);
            }}
          >
            {intl.formatMessage({ id: 'delete' })}
          </a>
        ) : null;
      },
    },
  ];

  useEffect(() => {
    handleSearch();
  }, []);

  return (
    <Container>
      <PageContainer title={intl.formatMessage({ id: 'account.title' })}>
        <Form form={form}>
          <LightArea>
            <div className="area-operate">
              <div className="title">
                {intl.formatMessage({ id: 'account.list' })}
              </div>
              <Space>
                <Form.Item name={'name'}>
                  <Input
                    placeholder={intl.formatMessage({
                      id: 'account.search.placeholder',
                    })}
                    onPressEnter={() => {
                      handleSearch();
                    }}
                    suffix={
                      <SearchOutlined
                        onClick={() => {
                          handleSearch();
                        }}
                      />
                    }
                  />
                </Form.Item>
                {userInfo?.role === 'admin' && (
                  <Button
                    className={batchState ? 'batch-active' : ''}
                    onClick={() => {
                      setBatchState(!batchState);
                    }}
                  >
                    {intl.formatMessage({ id: 'batchOperate' })}
                  </Button>
                )}
              </Space>
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
                        handleDeleteAccount(selectedRowKeys);
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
                  userInfo?.role === 'admin'
                    ? [...columns, ...operateColumns]
                    : columns
                }
                loading={queryByPage?.loading}
                rowKey={'id'}
                dataSource={pageData?.users || []}
                rowSelection={
                  userInfo?.role === 'admin'
                    ? {
                        selectedRowKeys,
                        getCheckboxProps: (record: DataType) => {
                          if (record?.name === userInfo?.name) {
                            return {
                              disabled: true,
                            };
                          }
                          return {};
                        },
                        onChange: (rowKeys: any[]) => {
                          setSelectedRowKeys(rowKeys);
                        },
                      }
                    : undefined
                }
                pagination={
                  pageData?.users?.length > 0
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
                  if (sorter.order) {
                    sort =
                      sorter.order === 'ascend'
                        ? 'create_time'
                        : '-create_time';
                  }
                  let role;
                  if (filters.role) {
                    role = filters.role;
                  }
                  handleSearch({
                    pageSize: pageSize,
                    page: current,
                    sort,
                    role,
                  });
                }}
              />
            </div>
          </LightArea>
        </Form>
      </PageContainer>
    </Container>
  );
};

export default Account;
