import { DownOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons';
import { Dropdown, Input, MenuProps } from 'antd';
import React, { useState } from 'react';
import AddSpaceDrawer from '../AddSpaceDrawer';
import { SpaceContent, SpaceMenu } from './style';

export default () => {
  const [curSpace, setCurSpace] = useState<string[]>(['3st menu item']);
  const [addSpaceOpen, setAddSpaceOpen] = useState<boolean>(false);
  const items: MenuProps['items'] = [
    {
      key: '1st menu item',
      label: '1st menu item',
    },
    {
      key: '2st menu item',
      label: '2st menu item',
    },
    {
      key: '3st menu item',
      label: '3st menu item',
    },
  ];

  return (
    <>
      <Dropdown
        // open
        menu={{
          items,
          selectable: true,
          onSelect: (item) => {
            setCurSpace(item?.selectedKeys);
          },
          selectedKeys: curSpace,
        }}
        dropdownRender={(menu) => (
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

            {React.cloneElement(menu as React.ReactElement)}
            <a
              onClick={() => {
                setAddSpaceOpen(true);
              }}
            >
              <PlusOutlined /> 新建空间
            </a>
          </SpaceMenu>
        )}
      >
        <SpaceContent>
          {curSpace[0]} <DownOutlined />
        </SpaceContent>
      </Dropdown>
      {addSpaceOpen && (
        <AddSpaceDrawer open={addSpaceOpen} setOpen={setAddSpaceOpen} />
      )}
    </>
  );
};
