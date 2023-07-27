import { Area } from '@/components/CommonStyle';
import {
  getUserInfo,
  getUserList,
  list,
} from '@/services/chaosmeta/UserController';
import { ExclamationCircleFilled, SearchOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { useModel, useRequest } from '@umijs/max';
import { Button, Input, Modal, Select, Space, Table, message } from 'antd';
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
  // 当前选中的数据
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const { initialState } = useModel('@@initialState');
  const queryUserList = useRequest(getUserList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      console.log(res, 'res=----');
    },
  });

  const test = useRequest(list, {
    manual: true,
    onSuccess: (res) => {
      console.log(res, 'res=----00000');
    },
  });

  const queryUserInfo = useRequest(getUserInfo, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      console.log(res, 'res=----00000');
    },
  });

  const authOptions = [
    {
      label: (
        <Role>
          <span>管理员</span>
          <div>拥有所有权限</div>
        </Role>
      ),
      value: 'admain',
      name: '管理员',
    },
    {
      label: (
        <Role>
          <span>普通用户</span>
          <div>可登录查看，同时叠加空间内权限</div>
        </Role>
      ),
      value: 'readonly',
      name: '普通用户',
    },
  ];
  const [baseColumns, setColumns] = useState<ColumnsType<DataType>>([]);

  const columns = [
    {
      title: '用户名',
      width: 80,
      dataIndex: 'userName',
    },
    {
      title: '加入时间',
      width: 160,
      dataIndex: 'time',
    },
    {
      title: '角色',
      width: 160,
      dataIndex: 'auth',
      filters: [
        {
          text: '管理员',
          value: 'admain',
        },
        {
          text: '普通用户',
          value: 'readonly',
        },
      ],
      render: (text: string) => {
        return (
          <Select
            dropdownStyle={{ minWidth: '200px' }}
            bordered={false}
            value={text}
            style={{ minWidth: '80px' }}
            optionLabelProp="label"
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
      },
    },
  ];
  const queryPage = (params) => {};

  /**
   * 分页查询
   */
  const queryByPage = useRequest(queryPage, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res: any) => {
      setPageData(res);
    },
  });

  const dataSource = [
    { id: '2', userName: 'ceshi' },
    {
      id: '1',
      userName: 'Serati Ma',
    },
  ];

  /**
   * 删除账号
   */
  const handleDeleteAccount = () => {
    Modal.confirm({
      title: '确认要删除当前所选账号吗？',
      icon: <ExclamationCircleFilled />,
      content: '删除账号用户将无法登录平台，要再次使用只能重新注册！',
      onOk() {
        return new Promise((resolve, reject) => {
          setTimeout(Math.random() > 0.5 ? resolve : reject, 1000);
          message.success('您已成功删除所选成员');
        }).catch(() => console.log('Oops errors!'));
      },
      onCancel() {},
    });
  };

  useEffect(() => {
    queryUserList.run();
    // queryUserInfo.run({ name: 'hlttest1' });
    // upToken.run();
    test?.run({ sceneType: 'CLASSIC_HA' });
    if (initialState?.userInfo.role === 'admain') {
      setColumns(() => {
        const newColumns: any = columns;
        newColumns.push({
          title: '操作',
          width: 60,
          fixed: 'right',
          render: (record: DataType) => {
            return record?.userName !== initialState?.userInfo?.name ? (
              <a onClick={handleDeleteAccount}>删除</a>
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
        <Area>
          <div className="area-operate">
            <div className="title">账号列表</div>
            <Space>
              <Input
                placeholder="请输入账号，用户名，权限"
                suffix={<SearchOutlined onClick={() => {}} />}
              />
              <Button
                disabled={selectedRowKeys?.length === 0}
                onClick={handleDeleteAccount}
              >
                批量删除
              </Button>
            </Space>
          </div>
          <div className="area-content">
            <Table
              columns={baseColumns}
              rowKey={'id'}
              dataSource={dataSource}
              rowSelection={{
                selectedRowKeys,
                getCheckboxProps: (record: DataType) => {
                  if (record?.userName === initialState?.userInfo?.name) {
                    return {
                      disabled: true,
                    };
                  }
                  return {};
                },
                onChange: (rowKeys: React.Key[]) => {
                  setSelectedRowKeys(rowKeys);
                },
              }}
              pagination={
                dataSource?.length > 0
                  ? {
                      showQuickJumper: true,
                      total: 200,
                    }
                  : false
              }
              onChange={(pagination: any, filters, sorter: any) => {
                const { current, pageSize } = pagination;
                queryByPage.run({ current, pageSize, auth: sorter.auth });
              }}
            />
          </div>
        </Area>
      </PageContainer>
    </Container>
  );
};

export default Account;
