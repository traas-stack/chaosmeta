import AddSpaceDrawer from '@/components/AddSpaceDrawer';
import { Area } from '@/components/CommonStyle';
import { ExclamationCircleFilled } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { Button, Modal, Tabs, message } from 'antd';
import React, { useState } from 'react';
import SpaceList from './SpaceList';
import { Container } from './style';

interface DataType {
  id: string;
  auth?: string;
  userName: string;
}

const SpaceManage: React.FC<unknown> = () => {
  const [pageData, setPageData] = useState<any>({});
  const [addSpaceOpen, setAddSpaceOpen] = useState<boolean>(false);

  const tabItems = [
    {
      label: '全部空间',
      key: 'all',
      children: (
        <Area>
          <SpaceList listType="all" title="空间列表" />
        </Area>
      ),
    },
    {
      label: '我管理的空间',
      key: 'my',
      children: (
        <Area>
          <SpaceList listType="my" title="我的空间列表" />
        </Area>
      ),
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

  return (
    <Container>
      <PageContainer title="空间管理">
        <Tabs
          items={tabItems}
          type="card"
          tabPosition="top"
          tabBarExtraContent={
            <>
              <Button
                type="primary"
                onClick={() => {
                  setAddSpaceOpen(true);
                }}
              >
                新建空间
              </Button>
            </>
          }
        />
      </PageContainer>
      {addSpaceOpen && (
        <AddSpaceDrawer open={addSpaceOpen} setOpen={setAddSpaceOpen} />
      )}
    </Container>
  );
};

export default SpaceManage;
