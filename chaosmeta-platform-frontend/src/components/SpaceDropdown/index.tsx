import { querySpaceList } from '@/services/chaosmeta/SpaceController';
import { DownOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons';
import { history, useRequest } from '@umijs/max';
import { Dropdown, Input } from 'antd';
import React, { useEffect, useState } from 'react';
import AddSpaceDrawer from '../AddSpaceDrawer';
import { SpaceContent, SpaceMenu } from './style';

export default () => {
  const [curSpace, setCurSpace] = useState<string[]>(['3st menu item']);
  const [addSpaceOpen, setAddSpaceOpen] = useState<boolean>(false);
  /**
   * 获取空间列表
   */
  const getSpaceList = useRequest(querySpaceList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      console.log(res, 'res');
    },
  });

  // side侧边菜单的展开/收起会影响这里，暂时用icon代替，todo
  const items: any = [
    {
      key: '1st menu item',
      icon: '1st menu item',
    },
    {
      key: '2st menu item',
      icon: '2st menu item',
    },
    {
      key: '3st menu item',
      icon: '3st menu item',
    },
    {
      key: '999 item',
      icon: '999 menu item',
    },
  ];

  useEffect(() => {
    console.log('=================')
    getSpaceList?.run({page: 1, pageSize: 10});
  }, []);

  return (
    <div id="spaceDropdown">
      <Dropdown
        menu={{
          items,
          selectable: true,
          onSelect: (item) => {
            console.log(item, 'item');
            setCurSpace(item?.selectedKeys);
          },
          selectedKeys: curSpace,
        }}
        overlayStyle={{ width: 'max-content' }}
        dropdownRender={(menu) => {
          console.log(
            menu,
            React.cloneElement(menu as React.ReactElement),
            'menu===',
          );
          return (
            <SpaceMenu>
              <div>
                <Input
                  placeholder="请输入关键词"
                  suffix={
                    <SearchOutlined
                      style={{ cursor: 'pointer' }}
                      onClick={() => {}}
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
              {/* <Menu items={items} /> */}
              {React.cloneElement(menu as React.ReactElement)}
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
          {curSpace[0]} <DownOutlined />
        </SpaceContent>
      </Dropdown>
      {addSpaceOpen && (
        <AddSpaceDrawer open={addSpaceOpen} setOpen={setAddSpaceOpen} />
      )}
    </div>
  );
};
