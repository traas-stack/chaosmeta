import { ExclamationCircleFilled, PlusOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { Button, Modal, Space, Tabs, Tooltip } from 'antd';
import React, { useState } from 'react';
import AddColonyModal from './AddColonyModal';
import AppConfigDrawer from './AppConfigDrawer';
import InstallAgentModal from './InstallAgentModal';
import TableList from './TableList';
import UpgradationDrawer from './UpgradationDrawer';
import { Container } from './style';

interface DataType {
  id: string;
  auth?: string;
  userName: string;
}

const Agent: React.FC<unknown> = () => {
  const [pageData, setPageData] = useState<any>({});
  const [selectedRows, setSelectedRows] = useState<any[]>([]);
  // 升级弹窗
  const [upgradationOpen, setUpgradationOpen] = useState<any>({
    id: '',
    title: '',
    open: false,
  });

  // 应用配置
  const [appConfigOpen, setAppConfigOpen] = useState<boolean>(false);
  // 添加集群
  const [addColonyOpen, setAddColonyOpen] = useState<boolean>(false);
  // 安装Agent
  const [installAgentOpen, setInstallAgentOpen] = useState<boolean>(false);
  const tabItems = [
    {
      label: '集群外Host',
      key: 'else',
      children: (
        <TableList
          selectedRows={selectedRows}
          setSelectedRows={setSelectedRows}
          setUpgradationOpen={setUpgradationOpen}
          title="集群外Host"
        />
      ),
    },
    {
      label: '默认集群',
      key: 'normal',
      children: (
        <TableList
          selectedRows={selectedRows}
          setSelectedRows={setSelectedRows}
          setUpgradationOpen={setUpgradationOpen}
          title="默认集群"
        />
      ),
      closable: false,
    },
  ];

  /**
   * 操作tab时
   */
  const handleEditTab = (action, key) => {
    if (key === 'add') {
      setAddColonyOpen(true);
    } else {
      Modal.confirm({
        title: '确认要删除当前集群吗？',
        icon: <ExclamationCircleFilled />,
        onOk() {
          // handleBatchDelete?.run({ user_ids: ids });
          // handleSearch();
          //   return new Promise((resolve, reject) => {
          // }).catch(() => console.log('Oops errors!'));
        },
        onCancel() {},
      });
    }
  };

  return (
    <Container>
      <PageContainer title="空间管理">
        <Tabs
          addIcon={
            <>
              <Tooltip title="添加集群">
                <PlusOutlined />
              </Tooltip>
            </>
          }
          items={tabItems}
          onEdit={handleEditTab}
          type="editable-card"
          animated
          tabPosition="top"
          tabBarExtraContent={
            <Space>
              <Button
                disabled={selectedRows.length === 0}
                onClick={() => {
                  setUpgradationOpen({
                    open: true,
                    title: '批量升级',
                  });
                }}
              >
                批量升级
              </Button>
              <Button
                onClick={() => {
                  setAppConfigOpen(true);
                }}
              >
                应用配置
              </Button>
              <Button
                type="primary"
                onClick={() => {
                  setInstallAgentOpen(true);
                }}
              >
                安装Agent
              </Button>
            </Space>
          }
        />
      </PageContainer>
      {/* 升级Agent */}
      {upgradationOpen?.open && (
        <UpgradationDrawer
          drawerData={upgradationOpen}
          setDrawerData={setUpgradationOpen}
          selectedRows={selectedRows}
        />
      )}
      {/* 应用配置 */}
      {appConfigOpen && (
        <AppConfigDrawer open={appConfigOpen} setOpen={setAppConfigOpen} />
      )}
      {/* 应用配置 */}
      {addColonyOpen && (
        <AddColonyModal open={addColonyOpen} setOpen={setAddColonyOpen} />
      )}
      {/* 安装Agent */}
      {installAgentOpen && (
        <InstallAgentModal
          open={installAgentOpen}
          setOpen={setInstallAgentOpen}
        />
      )}
    </Container>
  );
};

export default Agent;
