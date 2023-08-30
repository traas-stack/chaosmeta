import {
  queryClassSpaceList,
  querySpaceUserPermission,
} from '@/services/chaosmeta/SpaceController';
import { DownOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons';
import { history, useModel, useRequest } from '@umijs/max';
import { Dropdown, Empty, Input, Spin } from 'antd';
import React, { useEffect, useState } from 'react';
import AddSpaceDrawer from '../AddSpaceDrawer';
import { SpaceContent, SpaceMenu } from './style';

export default () => {
  const [addSpaceOpen, setAddSpaceOpen] = useState<boolean>(false);
  const [spaceList, setSpaceList] = useState<any>([]);
  const { setSpacePermission, curSpace, setCurSpace } = useModel('global');

  // 一级路由，切换空间时页面在以下路由时，不需要跳转，刷新页面接口即可，其他页面需要跳转回空间概览
  const parentRoute = [
    '/space/overview',
    '/space/experiment',
    '/space/experiment-result',
    '/space/setting',
  ];

  /**
   * 根据成员名称和空间id获取成员空间内权限信息
   */
  const getUserSpaceAuth = useRequest(querySpaceUserPermission, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res.code === 200) {
        setSpacePermission(res?.data);
      }
    },
  });

  /**
   * 更新地址栏空间id，并保存
   * @param id
   */
  const handleUpdateSpaceId = (id: string, name: string) => {
    if (id) {
      if (parentRoute.includes(history.location.pathname)) {
        history.push({
          pathname: history.location.pathname,
          query: {
            ...history.location.query,
            spaceId: id,
          },
        });
      } else {
        history.push({
          pathname: '/space/overview',
          query: {
            spaceId: id,
          },
        });
      }
      setCurSpace([id]);
      sessionStorage.setItem('spaceId', id);
      sessionStorage.setItem('spaceName', name);
    }
  };

  /**
   * 获取空间列表 -- 当前用户有查看权限的空间只读和读写
   */
  const getSpaceList = useRequest(queryClassSpaceList, {
    manual: true,
    formatResult: (res) => res,
    debounceInterval: 300,
    onSuccess: (res) => {
      if (res?.code === 200) {
        const namespaceList = res.data?.namespaces?.map(
          (item: { namespaceInfo: any }) => {
            // side侧边菜单的展开/收起会影响这里，暂时用icon代替，todo
            return {
              icon: item?.namespaceInfo?.name,
              key: item?.namespaceInfo?.id,
              id: item?.namespaceInfo?.id?.toString(),
              name: item?.namespaceInfo?.name,
            };
          },
        );
        const spaceId = history.location.query.spaceId;
        console.log(spaceId, 'spaceId---')
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
    // 获取空间列表
    getSpaceList?.run({ page: 1, page_size: 10, namespaceClass: 'relevant' });
  }, []);

  useEffect(() => {
    // getSpaceList?.run({ page: 1, page_size: 10, namespaceClass: 'relevant' });
    // 地址栏中存在空间id，需要将空间列表选项更新，并保存当前id
    if (history.location.query.spaceId) {
      setCurSpace([history.location.query.spaceId as string]);
      sessionStorage.setItem(
        'spaceId',
        history.location.query.spaceId as string,
      );
      getUserSpaceAuth?.run({
        id: history.location.query.spaceId as string,
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
        overlayStyle={{ width: 'max-content', maxHeight: '300px' }}
        dropdownRender={(menu) => {
          return (
            <SpaceMenu>
              <div className="search">
                <Input
                  placeholder="请输入关键词"
                  onChange={(event) => {
                    const value = event?.target?.value;
                    getSpaceList?.run({
                      name: value,
                      namespaceClass: 'relevant',
                    });
                  }}
                  suffix={
                    <SearchOutlined
                      style={{ cursor: 'pointer' }}
                      onClick={() => {
                        getSpaceList?.run({ namespaceClass: 'relevant' });
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
              <Spin spinning={getSpaceList?.loading}>
                <div style={{ maxHeight: '200px', overflowY: 'auto' }}>
                  {menu ? (
                    React.cloneElement(menu as React.ReactElement)
                  ) : (
                    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
                  )}
                </div>
              </Spin>

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
