import { LightArea } from '@/components/CommonStyle';
import ShowText from '@/components/ShowText';
import {
  batchDeleteUser,
  changeUserRole,
  getUserList,
} from '@/services/chaosmeta/UserController';
import { ExclamationCircleFilled, SearchOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { useModel, useRequest } from '@umijs/max';
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
  userName: string;
}

const Account: React.FC<unknown> = () => {
  const [pageData, setPageData] = useState<any>({});
  const [form] = Form.useForm();
  // 当前选中的数据
  const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([]);
  const { userInfo } = useModel('global');
  const [batchState, setBatchState] = useState<boolean>(false);
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
        message.success('用户角色修改成功');
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
      message.success('您已成功删除所选成员');
      handleSearch();
    },
  });

  /**
   * 删除账号
   */
  const handleDeleteAccount = (ids: string[]) => {
    if (ids?.length > 0) {
      Modal.confirm({
        title: '确认要删除当前所选账号吗？',
        icon: <ExclamationCircleFilled />,
        content: '删除账号用户将无法登录平台，要再次使用只能重新注册！',
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
          <span>管理员</span>
          <span>拥有所有权限</span>
        </Role>
      ),
      value: 'admin',
      name: '管理员',
    },
    {
      label: (
        <Role>
          <span>普通用户</span>
          <span>可登录查看，同时叠加空间内权限</span>
        </Role>
      ),
      value: 'normal',
      name: '普通用户',
    },
  ];

  const columns: any = [
    {
      title: '用户名',
      width: 80,
      dataIndex: 'name',
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
    {
      title: '角色',
      width: 160,
      dataIndex: 'role',
      filters: [
        {
          text: '管理员',
          value: 'admin',
        },
        {
          text: '普通用户',
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
      title: '操作',
      width: 60,
      fixed: 'right',
      render: (record: DataType) => {
        return record?.userName !== userInfo?.name ? (
          <a
            onClick={() => {
              handleDeleteAccount([record.id]);
            }}
          >
            删除
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
      <PageContainer title="账号管理">
        <Form form={form}>
          <LightArea>
            <div className="area-operate">
              <div className="title">账号列表</div>
              <Space>
                <Form.Item name={'name'}>
                  <Input
                    placeholder="请输入用户名进行搜索"
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
                    批量操作
                  </Button>
                )}
              </Space>
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
                        handleDeleteAccount(selectedRowKeys);
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
                          if (record?.userName === userInfo?.name) {
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
