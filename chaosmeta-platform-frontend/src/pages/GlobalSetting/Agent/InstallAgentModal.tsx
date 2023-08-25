/**
 * 安装Agent
 */
import { RightOutlined, SearchOutlined } from '@ant-design/icons';
import { Button, Drawer, Form, Input, Select, Space, Table } from 'antd';
import React, { useState } from 'react';
import { InstallAgentContainer } from './style';
interface IProps {
  open: boolean;
  setOpen: (val: boolean) => void;
}

const InstallAgentModal: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;
  const [form] = Form.useForm();
  const [rightData, setRightData] = useState<any>([]);
  const [leftSelectedRows, setLeftSelectedRows] = useState<any[]>([]);
  const [rightSelectedRows, setRightSelectedRows] = useState<any[]>([]);

  const handleClose = () => {
    setOpen(false);
  };

  const leftData = [
    { hostname: '11', ip: '22', id: 1 },
    { hostname: '112', ip: '22', id: 2 },
    { hostname: '113', ip: '22', id: 3 },
    { hostname: '114', ip: '22', id: 4 },
    { hostname: '11', ip: '22', id: 11 },
    { hostname: '112', ip: '22', id: 22 },
    { hostname: '113', ip: '22', id: 33 },
    { hostname: '114', ip: '22', id: 44 },
    { hostname: '11', ip: '22', id: 15 },
    { hostname: '112', ip: '22', id: 26 },
    { hostname: '113', ip: '22', id: 37 },
    { hostname: '114', ip: '22', id: 48 },
  ];

  const columns = [
    {
      title: 'hostname',
      dataIndex: 'hostname',
    },
    {
      title: 'IP',
      dataIndex: 'ip',
    },
  ];

  const options = [
    {
      value: 'appName',
      label: '应用名称',
    },
    {
      value: 'k8s',
      label: 'K8s Label',
    },
    {
      value: 'hostname',
      label: 'Hostname',
    },
    {
      value: 'ip',
      label: 'IP',
    },
  ];

  return (
    <Drawer
      open={open}
      onClose={handleClose}
      title="安装Agent"
      width={960}
      bodyStyle={{ paddingTop: 0 }}
      footer={
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button onClick={handleClose}>取消</Button>
            <Button type="primary">确定</Button>
          </Space>
        </div>
      }
    >
      <InstallAgentContainer>
        <Form form={form}>
          <div className="title">选择机器</div>
          <div className="content">
            <div className="left">
              <div className="header">机器列表</div>
              <div className="search">
                <Space.Compact>
                  <Form.Item name={'type'}>
                    <Select defaultValue="appName" options={options} />
                  </Form.Item>
                  <Form.Item name={'name'}>
                    <Input
                      placeholder="请输入"
                      suffix={
                        <SearchOutlined
                          onClick={() => {
                            // handleSearch();
                          }}
                        />
                      }
                    />
                  </Form.Item>
                </Space.Compact>
              </div>
              <div className="table">
                <Table
                  columns={columns}
                  dataSource={leftData}
                  rowKey={'id'}
                  pagination={false}
                  rowSelection={{
                    selectedRowKeys: leftSelectedRows.map((item) => item.id),
                    onChange: (rowKeys: any[], rows: any[]) => {
                      setLeftSelectedRows(rows);
                    },
                  }}
                  scroll={{ y: 600 }}
                />
              </div>
            </div>
            <div className="transfer">
              <Button
                disabled={leftSelectedRows?.length === 0}
                type={leftSelectedRows?.length > 0 ? 'primary' : 'default'}
                onClick={() => {
                  if (leftSelectedRows?.length > 0) {
                    const rightKeys = rightData?.map((item) => item.id);
                    const newList = leftSelectedRows?.filter(
                      (item) => !rightKeys?.includes(item.id),
                    );
                    setRightData([...rightData, ...newList]);
                    setRightSelectedRows([...rightSelectedRows, ...newList]);
                  }
                }}
              >
                <RightOutlined />
              </Button>
            </div>
            <div className="right">
              <div className="header">
                已选机器（{rightSelectedRows.length}）
              </div>
              <div className="table">
                <Table
                  columns={columns}
                  dataSource={rightData}
                  rowKey={'id'}
                  pagination={false}
                  rowSelection={{
                    selectedRowKeys: rightSelectedRows.map((item) => item.id),
                    onChange: (rowKeys: any[], rows: any[]) => {
                      setRightSelectedRows(rows);
                    },
                  }}
                  scroll={{ y: 400 }}
                />
              </div>
              <div style={{ height: '32px' }}></div>
            </div>
          </div>
        </Form>
      </InstallAgentContainer>
    </Drawer>
  );
};
export default React.memo(InstallAgentModal);
