/**
 * Agent升级抽屉
 */
import { Button, Drawer, Input, Space, Table } from 'antd';
import React, { useEffect, useState } from 'react';
import { UpgradationContainer } from './style';
interface IProps {
  drawerData: {
    id: string;
    open: boolean;
    title: string;
  };
  setDrawerData: (val: any) => void;
  selectedRows: any[];
}
const UpgradationDrawer: React.FC<IProps> = (props) => {
  const { drawerData, setDrawerData, selectedRows } = props;
  const [newSelectedRows, setNewSelectedRows] = useState<any[]>([]);
  const handleClose = () => {
    setDrawerData({
      id: '',
      title: '',
      open: false,
    });
    setNewSelectedRows([]);
  };

  console.log(newSelectedRows, 'newSelectedRows===');

  const columns: any = [
    {
      title: 'Nodename',
      width: 80,
      dataIndex: 'Nodename',
    },
    {
      title: 'Hostname',
      width: 80,
      dataIndex: 'Hostname',
    },
    {
      title: 'IP',
      width: 160,
      dataIndex: 'ip',
    },
    {
      title: 'Agent版本',
      width: 160,
      dataIndex: 'version',
    },
  ];

  useEffect(() => {
    if (drawerData.open) {
      setNewSelectedRows(selectedRows);
    }
  }, [drawerData]);

  return (
    <Drawer
      open={drawerData.open}
      onClose={handleClose}
      title={drawerData.title}
      width={800}
      footer={
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button onClick={handleClose}>取消</Button>
            <Button type="primary">确定</Button>
          </Space>
        </div>
      }
    >
      <UpgradationContainer>
        <div className="version">
          <div>Agent版本升级至：</div>{' '}
          <Input placeholder="请输入" disabled value={'2.0.1 (稳定版本)'} />
        </div>
      </UpgradationContainer>
      {!drawerData?.id && (
        <Table
          columns={columns}
          dataSource={selectedRows}
          rowKey={'id'}
          pagination={false}
          rowSelection={{
            selectedRowKeys: newSelectedRows?.map((item) => {
              return item.id;
            }),
            onChange: (rowKeys: any[], rows) => {
              setNewSelectedRows(rows);
            },
          }}
        />
      )}
    </Drawer>
  );
};
export default React.memo(UpgradationDrawer);
