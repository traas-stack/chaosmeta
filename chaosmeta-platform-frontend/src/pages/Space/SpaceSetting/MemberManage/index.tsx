/**
 * 成员管理tab
 */
import { Area } from '@/components/CommonStyle';
import { Role } from '@/pages/GlobalSetting/Account/style';
import { ExclamationCircleFilled, SearchOutlined } from '@ant-design/icons';
import { useModel } from '@umijs/max';
import {
  Alert,
  Button,
  Input,
  Modal,
  Select,
  Space,
  Table,
  message,
} from 'antd';
import { ColumnsType } from 'antd/es/table';
import React, { useState } from 'react';
import AddMemberModal from './AddMemberModal';

interface IProps {}
interface DataType {
  id: string;
  account: string;
  auth: string;
}

const MemberManage: React.FC<IProps> = () => {
  const { initialState } = useModel('@@initialState');
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [addMemberOpen, setAddMemberOpen] = useState<boolean>(false);

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

  /**
   * 删除账号
   */
  const handleDeleteAccount = () => {
    Modal.confirm({
      title: '您确认要删除当前所选成员吗？',
      icon: <ExclamationCircleFilled />,
      content: '删除空间内成员，该成员将无法进入该空间！',
      onOk() {
        return new Promise((resolve, reject) => {
          setTimeout(Math.random() > 0.5 ? resolve : reject, 1000);
          message.success('您已成功删除所选成员');
        }).catch(() => console.log('Oops errors!'));
      },
      onCancel() {},
    });
  };

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
          {/* {initialState?.userInfo?.role === 'admain' && ( */}
          <>
            <Button disabled={selectedRowKeys?.length === 0}>批量删除</Button>
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
      {addMemberOpen && (
        <AddMemberModal setOpen={setAddMemberOpen} open={addMemberOpen} />
      )}
    </Area>
  );
};

export default React.memo(MemberManage);
