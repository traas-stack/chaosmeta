import { Area } from '@/components/CommonStyle';
import ShowText from '@/components/ShowText';
import {
  batchDeleteUser,
  changeUserRole,
  getUserList,
  list,
} from '@/services/chaosmeta/UserController';
import { ExclamationCircleFilled, SearchOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { useModel, useRequest } from '@umijs/max';
import {
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
  const [baseColumns, setColumns] = useState<ColumnsType<DataType>>([]);

  const test = useRequest(list, {
    manual: true,
    formatResult: (res) => {
      console.log(res, 'res-test===');
      return res;
    },
    onSuccess: (res) => {
      console.log(res, 'res=----00000');
    },
  });

  /**
   * 分页查询
   */
  const queryByPage = useRequest(getUserList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res: any) => {
      console.log(res, 'res=====999999');
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
    console.log(queryParam, 'value===');
    // return;
    queryByPage.run(queryParam);
  };

  const dataSource = [
    { id: '2', userName: 'ceshi' },
    {
      id: '1',
      userName: 'Serati Ma',
    },
  ];

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

  // 角色下拉选项
  const authOptions = [
    {
      label: (
        <Role>
          <span>管理员</span>
          <div>拥有所有权限</div>
        </Role>
      ),
      value: 'admin',
      name: '管理员',
    },
    {
      label: (
        <Role>
          <span>普通用户</span>
          <div>可登录查看，同时叠加空间内权限</div>
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
      render: (text: string, record: { id: number }) => {
        if (userInfo.role === 'admin') {
          return (
            <Select
              dropdownStyle={{ minWidth: '200px' }}
              bordered={false}
              value={text}
              style={{ minWidth: '80px' }}
              optionLabelProp="label"
              onChange={(value: string) => {
                console.log(value, 'value');
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
        return authOptions?.filter((item) => item?.value === text)[0]?.name;
      },
    },
  ];

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
          //   return new Promise((resolve, reject) => {
          // }).catch(() => console.log('Oops errors!'));
        },
        onCancel() {},
      });
    }
  };

  useEffect(() => {
    handleSearch();
    // test.run({sceneType: 'CLASSIC_HA'})
    if (userInfo.role === 'admin') {
      setColumns(() => {
        const newColumns: any = columns;
        newColumns.push({
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
        });
        return newColumns;
      });
    } else {
      setColumns(columns);
    }
  }, []);

  return (
    <Container>
      <PageContainer title="账号管理">
        <Form form={form}>
          <Area>
            <div className="area-operate">
              <div
                className="title"
                onClick={() => {
                  test.run({ sceneType: 'CLASSIC_HA' });
                }}
              >
                账号列表
              </div>
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
                    disabled={selectedRowKeys?.length === 0}
                    onClick={() => {
                      handleDeleteAccount(selectedRowKeys);
                    }}
                  >
                    批量删除
                  </Button>
                )}
              </Space>
            </div>
            <div className="area-content">
              <Table
                columns={baseColumns}
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
                  dataSource?.length > 0
                    ? {
                        showQuickJumper: true,
                        total: pageData?.total,
                        current: pageData?.page,
                        pageSize: pageData?.pageSize,
                      }
                    : false
                }
                onChange={(pagination: any, filters, sorter: any) => {
                  console.log(pagination, 'filters');
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
          </Area>
        </Form>
      </PageContainer>
    </Container>
  );
};

export default Account;
