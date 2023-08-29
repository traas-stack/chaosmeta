import { LightArea } from '@/components/CommonStyle';
import ShowText from '@/components/ShowText';
import { ExclamationCircleFilled } from '@ant-design/icons';
import {
  Alert,
  Button,
  Col,
  Form,
  Input,
  Modal,
  Row,
  Space,
  Table,
} from 'antd';
import React from 'react';

interface IProps {
  title: string;
  selectedRows: any[];
  setSelectedRows: any;
  setUpgradationOpen: any;
}

const SpaceList: React.FC<IProps> = (props) => {
  const { title, selectedRows, setSelectedRows, setUpgradationOpen } = props;
  const [form] = Form.useForm();

  /**
   * 卸载
   */
  const handleUnload = () => {
    Modal.confirm({
      title: '确认要卸载已安装的Agent吗？',
      icon: <ExclamationCircleFilled />,
      onOk() {
        // handleBatchDelete?.run({ user_ids: ids });
        // handleSearch();
        //   return new Promise((resolve, reject) => {
        // }).catch(() => console.log('Oops errors!'));
      },
      onCancel() {},
    });
  };

  const columns: any = [
    {
      title: 'Hostname',
      width: 80,
      dataIndex: 'Hostname',
    },
    {
      title: 'IP',
      width: 160,
      dataIndex: 'ip',
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
    {
      title: '应用',
      width: 160,
      dataIndex: 'app',
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
    {
      title: 'Agent版本',
      width: 160,
      dataIndex: 'version',
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
    {
      title: 'Agent状态',
      width: 160,
      dataIndex: 'state',
      render: (text: string) => {
        return <ShowText isTime value={text} />;
      },
    },
    {
      title: '操作',
      width: 160,
      dataIndex: 'id',
      render: (text: string) => {
        return (
          <Space>
            <a
              onClick={() => {
                setUpgradationOpen({
                  open: true,
                  title: '升级',
                  id: text,
                });
              }}
            >
              升级Agent
            </a>
            <a onClick={handleUnload}>卸载</a>
          </Space>
        );
      },
    },
  ];

  const dataSource = [
    {
      id: 1,
      ip: '99',
    },
  ];

  return (
    <LightArea>
      <div className="search">
        <Form form={form} labelCol={{ span: 6 }}>
          <Row gutter={24}>
            <Col span={8}>
              <Form.Item name={'hostName'} label="Hostname">
                <Input placeholder="请输入" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name={'IP'} label="IP">
                <Input placeholder="请输入" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name={'app'} label="应用">
                <Input placeholder="请选择" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name={'agentVersion'} label="Agent版本">
                <Input placeholder="请选择" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name={'agentState'} label="Agent状态">
                <Input placeholder="请选择" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Space>
                <Button>重置</Button>
                <Button type="primary">查询</Button>
              </Space>
            </Col>
          </Row>
        </Form>
      </div>
      <div className="title">{title}</div>
      <div className="area-content">
        <div>
          {selectedRows?.length > 0 && (
            <Alert
              message={
                <>
                  已选择 <a>{selectedRows.length}</a> 项
                  <a
                    style={{ paddingLeft: '24px' }}
                    onClick={() => {
                      setSelectedRows([]);
                    }}
                  >
                    清空
                  </a>
                </>
              }
              style={{ marginBottom: '16px' }}
              type="info"
              action={
                <Space>
                  <a style={{ color: '#FF4D4F' }} onClick={handleUnload}>
                    批量卸载
                  </a>
                  <a
                    onClick={() => {
                      setUpgradationOpen({
                        open: true,
                        title: '批量升级',
                      });
                    }}
                  >
                    批量升级
                  </a>
                </Space>
              }
              showIcon
            />
          )}
        </div>
        <Table
          columns={columns}
          dataSource={dataSource}
          rowKey={'id'}
          rowSelection={{
            selectedRowKeys: selectedRows?.map((item) => {
              return item.id;
            }),
            onChange: (rowKeys: any[], rows) => {
              setSelectedRows(rows);
            },
          }}
          pagination={
            dataSource?.length > 0
              ? {
                  showQuickJumper: true,
                  total: 100,
                  // total: pageData?.total,
                  // current: pageData?.page,
                  // pageSize: pageData?.pageSize,
                }
              : false
          }
          // onChange={(pagination: any, filters) => {
          //   const { current, pageSize } = pagination;
          //   let role;
          //   if (filters.role) {
          //     role = filters.role;
          //   }
          //   handleSearch({
          //     pageSize: pageSize,
          //     page: current,
          //     sort,
          //     role,
          //   });
          // }}
        />
      </div>
    </LightArea>
  );
};

export default SpaceList;
