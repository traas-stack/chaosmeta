/**
 * 成员管理tab
 */
import { LightArea } from '@/components/CommonStyle';
import { ExclamationCircleFilled, SearchOutlined } from '@ant-design/icons';
import { Button, Input, Modal, Space, Table, message } from 'antd';
import { ColumnsType } from 'antd/es/table';
import React, { useState } from 'react';

interface IProps {}

const TagManage: React.FC<IProps> = () => {
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  /**
   * 删除标签
   */
  const handleDeleteAccount = () => {
    Modal.confirm({
      title: '确认要删除所选标签吗？',
      icon: <ExclamationCircleFilled />,
      content: '删除将会对应删除当前空间里关联实验的标签',
      onOk() {
        return new Promise((resolve, reject) => {
          setTimeout(Math.random() > 0.5 ? resolve : reject, 1000);
          message.success('删除成功');
        }).catch(() => console.log('Oops errors!'));
      },
      onCancel() {},
    });
  };

  const columns: ColumnsType<any> = [
    {
      title: '名称',
      width: 160,
      dataIndex: 'name',
    },
    {
      title: '标签颜色',
      width: 80,
      dataIndex: 'color',
    },
    {
      title: '标签样式',
      width: 160,
      dataIndex: 'time',
      sorter: true,
    },
    {
      title: '创建人',
      width: 160,
      dataIndex: 'user',
      sorter: true,
    },
    {
      title: '创建时间',
      width: 160,
      dataIndex: 'time',
      sorter: true,
    },
    {
      title: '操作',
      width: 60,
      fixed: 'right',
      render: () => {
        return <a onClick={handleDeleteAccount}>删除</a>;
      },
    },
  ];
  return (
    <LightArea>
      <div className="area-operate">
        <div className="title">标签列表</div>
        <Space>
          <Input
            placeholder="请输入标签名称"
            suffix={
              <SearchOutlined
                onClick={() => {
                  console.log('-==');
                }}
              />
            }
          />
          <Button>批量删除</Button>
          <Button type="primary">新建标签</Button>
        </Space>
      </div>
      <div className="area-content">
        <Table
          columns={columns}
          rowKey={'id'}
          dataSource={[{ id: '1', account: 'hlt', auth: 'admain' }]}
          scroll={{ x: 'max-content' }}
          rowSelection={{
            selectedRowKeys,
            onChange: (rowKeys: React.Key[]) => {
              setSelectedRowKeys(rowKeys);
            },
          }}
          pagination={{
            showQuickJumper: true,
            total: 200,
          }}
        />
      </div>
    </LightArea>
  );
};

export default React.memo(TagManage);
