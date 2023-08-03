/**
 * 成员管理tab
 */
import { LightArea } from '@/components/CommonStyle';
import ShowText from '@/components/ShowText';
import { Role } from '@/pages/GlobalSetting/Account/style';
import {
  querySpaceUserList,
  spaceBatchDeleteUser,
} from '@/services/chaosmeta/SpaceController';
import { useParamChange } from '@/utils/useParamChange';
import { ExclamationCircleFilled, SearchOutlined } from '@ant-design/icons';
import { history, useRequest } from '@umijs/max';
import { Alert, Button, Form, Input, Modal, Select, Space, Table } from 'antd';
import { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';
import AddMemberModal from './AddMemberModal';

interface IProps {}
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
const MemberManage: React.FC<IProps> = () => {
  const [form] = Form.useForm();
  const [selectedRowKeys, setSelectedRowKeys] = useState<number[]>([]);
  const [addMemberOpen, setAddMemberOpen] = useState<boolean>(false);
  const spaceIdChange = useParamChange('spaceId');
  const [pageData, setPageData] = useState<PageData>({
    page: 1,
    pageSize: 10,
    total: 0,
    users: [],
  });

  /**
   * 分页接口
   */
  const getUserList = useRequest(querySpaceUserList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      console.log(res, 'res');
      if (res?.data?.users) {
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
      // permission,
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
      console.log(res, '-------------');
    },
  });

  /**
   * 删除标签
   */
  // const handleDeleteUser = (id?: number) => {
  //   const user_ids = id ? [id] : selectedRowKeys;
  //   Modal.confirm({
  //     title: '确认要删除所选标签吗？',
  //     icon: <ExclamationCircleFilled />,
  //     content: '删除将会对应删除当前空间里关联实验的标签',
  //     onOk() {
  //       deleteUser?.run({
  //         user_ids,
  //         id: Number(history?.location.query?.spaceId),
  //       });
  //       // return new Promise((resolve, reject) => {
  //       //   setTimeout(Math.random() > 0.5 ? resolve : reject, 1000);
  //       //   message.success('删除成功');
  //       // }).catch(() => console.log('Oops errors!'));
  //     },
  //     onCancel() {},
  //   });
  // };

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
        console.log(user_ids, 'user_ids')
        // return
        deleteUser?.run({
          user_ids,
          id: Number(history?.location.query?.spaceId),
        });
        // return new Promise((resolve, reject) => {
        //   setTimeout(Math.random() > 0.5 ? resolve : reject, 1000);
        //   message.success('您已成功删除所选成员');
        // }).catch(() => console.log('Oops errors!'));
      },
      onCancel() {},
    });
  };

  const authOptions = [
    {
      label: (
        <Role>
          <span>读写</span>
          <div>拥有所有权限</div>
        </Role>
      ),
      value: 'admain',
      name: '读写',
    },
    {
      label: (
        <Role>
          <span>只读</span>
          <div>只能查看</div>
        </Role>
      ),
      value: 'readonly',
      name: '只读',
    },
  ];

  const columns: ColumnsType<DataType> = [
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
          value: 'admain',
        },
        {
          text: '只读',
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
                placeholder="请输入用户名进行搜索"
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
            <>
              <Button
                disabled={selectedRowKeys?.length === 0}
                onClick={() => {
                  handleDeleteUser();
                }}
              >
                批量删除
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
            {/* )} */}
          </Space>
        </Form>
      </div>
      <div className="area-content">
        {selectedRowKeys?.length > 0 && (
          <Alert
            message={
              <>
                已选择 <a>{selectedRowKeys.length}</a> 项
              </>
            }
            style={{ marginBottom: '16px' }}
            type="info"
            action={
              <a
                onClick={() => {
                  setSelectedRowKeys([]);
                }}
              >
                清空
              </a>
            }
            showIcon
          />
        )}

        <Table
          columns={columns}
          rowKey={'id'}
          loading={getUserList.loading}
          dataSource={pageData.users}
          scroll={{ x: 'max-content' }}
          rowSelection={{
            selectedRowKeys,
            onChange: (rowKeys: any[]) => {
              setSelectedRowKeys(rowKeys);
            },
          }}
          pagination={{
            showQuickJumper: true,
            total: 200,
          }}
          onChange={(pagination: any, filters, sorter: any) => {
            console.log(pagination, 'filters');
            const { current, pageSize } = pagination;
            let sort;
            if (sorter.order) {
              sort = sorter.order === 'ascend' ? 'create_time' : '-create_time';
            }
            let permission;
            if (filters.permission) {
              permission = filters.permission;
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
        <AddMemberModal setOpen={setAddMemberOpen} open={addMemberOpen} />
      )}
    </LightArea>
  );
};

export default React.memo(MemberManage);
