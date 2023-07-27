/**
 * 成员管理tab
 */
import { Area } from '@/components/CommonStyle';
import { SearchOutlined } from '@ant-design/icons';
import { Button, Input, Select, Space, Table } from 'antd';
import { ColumnsType } from 'antd/es/table';
import React from 'react';

interface IProps {}
interface DataType {
  id: string;
  account: string;
  auth: string;
}

const MemberManage: React.FC<IProps> = () => {
  const authOptions = [
    {
      label: '读写',
      value: 'writing',
    },
    {
      label: '只读',
      value: 'readonly',
    },
    {
      label: '管理员',
      value: 'admain',
    },
  ];
  const columns: ColumnsType<DataType> = [
    {
      title: '账号',
      width: 160,
      dataIndex: 'account',
    },
    {
      title: '用户名',
      width: 80,
      dataIndex: 'userName',
    },
    {
      title: '权限',
      width: 160,
      dataIndex: 'auth',
      render: (text: string) => {
        return (
          <Select
            style={{ minWidth: '80px' }}
            value={text}
            bordered={false}
            options={authOptions}
          />
        );
      },
    },
    {
      title: '加入时间',
      width: 160,
      dataIndex: 'time',
    },
    {
      title: '操作',
      width: 60,
      fixed: 'right',
      render: () => {
        return <a>删除</a>;
      },
    },
  ];
  return (
    <Area>
      <div className="area-operate">
        <div className="title">成员列表</div>
        <Space>
          <Input
            placeholder="请输入账号，用户名，权限"
            suffix={
              <SearchOutlined
                onClick={() => {
                  console.log('-==');
                }}
              />
            }
          />
          <Button>批量导入</Button>
          <Button>批量删除</Button>
          <Button type="primary">添加成员</Button>
        </Space>
      </div>
      <div className="tab-content">
        <Table
          columns={columns}
          rowKey={'id'}
          dataSource={[{ id: '1', account: 'hlt', auth: 'admain' }]}
          scroll={{ x: 'max-content' }}
          pagination={{
            showQuickJumper:true,
            total: 200,
            
          }}
        />
      </div>
    </Area>
  );
};

export default React.memo(MemberManage);
