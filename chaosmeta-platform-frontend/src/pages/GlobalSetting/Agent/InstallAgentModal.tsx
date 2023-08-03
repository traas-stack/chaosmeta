/**
 * 安装Agent
 */
import { SearchOutlined } from '@ant-design/icons';
import { Button, Drawer, Form, Input, Select, Space, Tree } from 'antd';
import { DataNode } from 'antd/es/tree';
import React, { useState } from 'react';
import { InstallAgentContainer } from './style';
interface IProps {
  open: boolean;
  setOpen: (val: boolean) => void;
}
const InstallAgentModal: React.FC<IProps> = (props) => {
  const { open, setOpen } = props;
  const [form] = Form.useForm();

  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([
    '0-0-0',
    '0-0-1',
  ]);
  const [checkedKeys, setCheckedKeys] = useState<React.Key[]>(['0-0-0']);
  const [selectedKeys, setSelectedKeys] = useState<React.Key[]>([]);
  const [autoExpandParent, setAutoExpandParent] = useState<boolean>(true);

  const treeData: DataNode[] = [
    // {
    //   title: '全选',
    //   key: '0',
    //   children
    // },
    {
      title: '0-0',
      key: '0-0',
      children: [
        {
          title: '0-0-0',
          key: '0-0-0',
          children: [
            { title: '0-0-0-0', key: '0-0-0-0' },
            { title: '0-0-0-1', key: '0-0-0-1' },
            { title: '0-0-0-2', key: '0-0-0-2' },
          ],
        },
        {
          title: '0-0-1',
          key: '0-0-1',
          children: [
            { title: '0-0-1-0', key: '0-0-1-0' },
            { title: '0-0-1-1', key: '0-0-1-1' },
            { title: '0-0-1-2', key: '0-0-1-2' },
          ],
        },
        {
          title: '0-0-2',
          key: '0-0-2',
        },
      ],
    },
    {
      title: '0-1',
      key: '0-1',
      children: [
        { title: '0-1-0-0', key: '0-1-0-0' },
        { title: '0-1-0-1', key: '0-1-0-1' },
        { title: '0-1-0-2', key: '0-1-0-2' },
      ],
    },
    {
      title: '0-2',
      key: '0-2',
    },
  ];

  const onExpand = (expandedKeysValue: React.Key[]) => {
    console.log('onExpand', expandedKeysValue);
    // if not set autoExpandParent to false, if children expanded, parent can not collapse.
    // or, you can remove all expanded children keys.
    setExpandedKeys(expandedKeysValue);
    setAutoExpandParent(false);
  };

  const onCheck = (checkedKeysValue: React.Key[]) => {
    console.log('onCheck', checkedKeysValue);
    setCheckedKeys(checkedKeysValue);
  };

  const onSelect = (selectedKeysValue: React.Key[], info: any) => {
    console.log('onSelect', info);
    setSelectedKeys(selectedKeysValue);
  };
  const handleClose = () => {
    setOpen(false);
  };
  return (
    <Drawer
      open={open}
      onClose={handleClose}
      title="安装Agent"
      width={960}
      bodyStyle={{ paddingTop: 0 }}
      footer={
        <div>
          <Space>
            <Button onClick={handleClose}>取消</Button>
            <Button type="primary">确定</Button>
          </Space>
        </div>
      }
    >
      <InstallAgentContainer>
        <Form form={form}>
          <div className="content">
            <div className="left">
              <Input.Group compact>
                <Select defaultValue="app" style={{ width: '120px' }}>
                  <Select.Option value="app">应用名称</Select.Option>
                  <Select.Option value="k8s">K8s Label</Select.Option>
                  <Select.Option value="hostName">HostName</Select.Option>
                  <Select.Option value="ip">IP</Select.Option>
                </Select>
                <Input
                  style={{ width: '320px' }}
                  onPressEnter={() => {}}
                  suffix={
                    <SearchOutlined
                      style={{ cursor: 'pointer' }}
                      onClick={() => {}}
                    />
                  }
                />
              </Input.Group>

              <Tree
                checkable
                onExpand={onExpand}
                expandedKeys={expandedKeys}
                autoExpandParent={autoExpandParent}
                onCheck={onCheck}
                checkedKeys={checkedKeys}
                onSelect={onSelect}
                selectedKeys={selectedKeys}
                treeData={treeData}
                titleRender={(nodeData: any) => {
                  console.log(nodeData, 'nodeData====');
                  const selectCount = checkedKeys.filter((item: string) => {
                    return item.includes(nodeData.key)
                  })?.length
                  console.log(selectCount, 'selectCount')
                  return (
                    <div>
                      {nodeData.title}
                      {nodeData.children?.length > 0 && <span>({9}/{nodeData.children})</span>}
                    </div>
                  );
                }}
              />
            </div>
            <div></div>
          </div>
          {/* <Form.Item name={'name'} label="名称">
          <Input placeholder="请输入" />
        </Form.Item> */}
        </Form>
      </InstallAgentContainer>
    </Drawer>
  );
};
export default React.memo(InstallAgentModal);
