import { querySpaceList } from '@/services/chaosmeta/SpaceController';
import { getSpaceUserList } from '@/services/chaosmeta/UserController';
import { DownOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons';
import { history, useModel, useRequest } from '@umijs/max';
import { Dropdown, Empty, Input } from 'antd';
import React, { useEffect, useState } from 'react';
import AddSpaceDrawer from '../AddSpaceDrawer';
import { SpaceContent, SpaceMenu } from './style';

export default () => {
  const [curSpace, setCurSpace] = useState<string[]>(['1']);
  const [addSpaceOpen, setAddSpaceOpen] = useState<boolean>(false);
  const [spaceList, setSpaceList] = useState<any>([]);
  const { userInfo, setSpacePermission } = useModel('global');

  /**
   * 根据成员名称和空间id获取成员空间内权限信息
   */
  const getUserSpaceAuth = useRequest(getSpaceUserList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res.code === 200) {
        // 存储用户空间权限
        const curUserName = userInfo?.name || localStorage.getItem('userName');
        const curUserInfo = res?.data?.users?.filter(
          (item: { name: string }) => item.name === curUserName,
        )[0];
        setSpacePermission(curUserInfo?.permission);
      }
    },
  });

  /**
   * 更新地址栏空间id，并保存
   * @param id
   */
  const handleUpdateSpaceId = (id: string, name: string) => {
    if (id) {
      history.push({
        pathname: history.location.pathname,
        query: {
          ...history.location.query,
          spaceId: id,
        },
      });
      setCurSpace([id]);
      sessionStorage.setItem('spaceId', id);
      sessionStorage.setItem('spaceName', name);
    }
  };
  /**
   * 获取空间列表
   */
  const getSpaceList = useRequest(querySpaceList, {
    manual: true,
    formatResult: (res) => res,
    debounceInterval: 300,
    onSuccess: (res) => {
      if (res?.data) {
        const namespaceList = res.data?.namespaces?.map(
          (item: { name: string; id: string }) => {
            // side侧边菜单的展开/收起会影响这里，暂时用icon代替，todo
            return {
              icon: item.name,
              key: item.id,
              id: item.id.toString(),
              name: item.name,
            };
          },
        );
        const spaceId = history.location.query.spaceId;
        // 初始化加载时，页面不存在空间id需要添加默认空间id
        if (!spaceId) {
          handleUpdateSpaceId(namespaceList[0]?.id, namespaceList[0]?.name);
        }
        // 保存当前空间id
        setSpaceList(namespaceList);
      }
    },
  });

  useEffect(() => {
    getSpaceList?.run({ page: 1, page_size: 10 });
  }, []);

  useEffect(() => {
    // 地址栏中存在空间id，需要将空间列表选项更新，并保存当前id
    if (history.location.query.spaceId) {
      setCurSpace([history.location.query.spaceId as string]);
      sessionStorage.setItem(
        'spaceId',
        history.location.query.spaceId as string,
      );
      getUserSpaceAuth?.run({
        id: history.location.query.spaceId,
        name: userInfo?.name || localStorage.getItem('userName'),
      });
    }
  }, [history.location.query.spaceId]);

  return (
    <div id="spaceDropdown">
      <Dropdown
        menu={{
          items: spaceList,
          selectable: true,
          onSelect: (item) => {
            const name = item?.item?.props?.name;
            handleUpdateSpaceId(item.key, name);
          },
          selectedKeys: curSpace,
        }}
        overlayStyle={{ width: 'max-content', maxHeight: '400px' }}
        dropdownRender={(menu) => {
          return (
            <SpaceMenu>
              <div>
                <Input
                  placeholder="请输入关键词"
                  onChange={(event) => {
                    const value = event?.target?.value;
                    getSpaceList?.run({ name: value });
                  }}
                  suffix={
                    <SearchOutlined
                      style={{ cursor: 'pointer' }}
                      onClick={() => {
                        getSpaceList?.run();
                      }}
                    />
                  }
                />
              </div>
              <div className="add-space">
                <a
                  onClick={() => {
                    setAddSpaceOpen(true);
                  }}
                >
                  <PlusOutlined /> 新建空间
                </a>
              </div>
              {menu ? (
                React.cloneElement(menu as React.ReactElement)
              ) : (
                <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
              )}
              <div className="more">
                <span>
                  没有相关空间？查看{' '}
                  <a
                    onClick={() => {
                      history.push('/setting/space');
                    }}
                  >
                    更多空间
                  </a>
                </span>
              </div>
            </SpaceMenu>
          );
        }}
      >
        <SpaceContent>
          {sessionStorage.getItem('spaceName') ||
            spaceList?.filter(
              (item: { id: string }) => item.id === curSpace[0],
            )[0]?.name}{' '}
          <DownOutlined />
        </SpaceContent>
      </Dropdown>
      {addSpaceOpen && (
        <AddSpaceDrawer open={addSpaceOpen} setOpen={setAddSpaceOpen} />
      )}
    </div>
  );
};
