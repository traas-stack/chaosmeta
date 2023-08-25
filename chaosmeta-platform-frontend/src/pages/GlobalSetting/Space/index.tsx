import AddSpaceDrawer from '@/components/AddSpaceDrawer';
import {
  deleteSpace,
  queryClassSpaceList,
} from '@/services/chaosmeta/SpaceController';
import { ExclamationCircleFilled, SearchOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { useRequest } from '@umijs/max';
import {
  Alert,
  Button,
  Empty,
  Form,
  Input,
  Modal,
  Pagination,
  Radio,
  Select,
  Space,
  Spin,
  Tabs,
  message,
} from 'antd';
import React, { useEffect, useState } from 'react';
import SpaceList from './SpaceList';
import { Container } from './style';

interface DataType {
  id: string;
  auth?: string;
  userName: string;
}

interface PageData {
  page: number;
  pageSize: number;
  total: number;
  namespaces: any[];
}
const SpaceManage: React.FC<unknown> = () => {
  const [addSpaceOpen, setAddSpaceOpen] = useState<boolean>(false);
  const [spaceType, setSpaceType] = useState<string>('all');
  const [tabKey, setTabKey] = useState<string>('all');
  const [form] = Form.useForm();
  const [pageData, setPageData] = useState<PageData>({
    page: 1,
    pageSize: 10,
    total: 0,
    namespaces: [],
  });

  /**
   * 获取空间列表接口
   */
  const getSpaceList = useRequest(queryClassSpaceList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      console.log(res, 'res----');
      if (res?.code === 200) {
        setPageData(res?.data || {});
      }
    },
  });

  /**
   * 分页数据获取
   * @param values
   */
  const handlePageSearch = (values?: {
    page?: number;
    pageSize?: number;
    sort?: string;
    namespaceClass?: string;
  }) => {
    const { page, pageSize, sort, namespaceClass } = values || {};
    const { searchType, name, member } = form.getFieldsValue();
    const params = {
      page: page || pageData.page || 1,
      page_size: pageSize || pageData.pageSize || 10,
      sort,
      name,
      namespaceClass: namespaceClass || spaceType,
    };
    getSpaceList?.run(params);
  };

  const handleDeleteSpace = useRequest(deleteSpace, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success('您已成功删除所选空间');
        handlePageSearch();
      }
    },
  });

  const tabItems = [
    {
      label: '全部空间',
      key: 'all',
    },
    {
      label: '我管理的空间',
      key: 'myAdmin',
    },
  ];

  /**
   * 删除账号
   */
  const handleDelete = (id: number) => {
    // return
    Modal.confirm({
      title: '确认要删除当前所选空间吗？',
      icon: <ExclamationCircleFilled />,
      onOk() {
        return handleDeleteSpace?.run({ id });
      },
    });
  };

  // 空间权限类型
  const spaceTypes = [
    {
      label: '全部',
      value: 'all',
    },
    {
      label: '未加入',
      value: 'not',
    },
    {
      label: '只读',
      value: 'read',
    },
    {
      label: '读写',
      value: 'write',
    },
  ];

  // 检索类型
  const searchOptions = [
    {
      label: '空间名称',
      value: 'spaceName',
    },
    {
      label: '空间成员',
      value: 'spaceMember',
    },
  ];

  useEffect(() => {
    handlePageSearch();
  }, []);
  return (
    <PageContainer title="空间管理">
      <Container>
        <Form form={form}>
          <Tabs
            items={tabItems}
            activeKey={tabKey}
            onChange={(val) => {
              setTabKey(val);
            }}
            tabBarExtraContent={
              <Space>
                {tabKey === 'all' && (
                  <Radio.Group defaultValue={spaceType} value={spaceType}>
                    {spaceTypes?.map((item) => {
                      return (
                        <Radio.Button
                          value={item.value}
                          key={item.value}
                          onChange={(event) => {
                            setSpaceType(event?.target?.value);
                            handlePageSearch({
                              namespaceClass: event?.target?.value,
                            });
                          }}
                        >
                          {item.label}
                        </Radio.Button>
                      );
                    })}
                  </Radio.Group>
                )}
                <Space.Compact>
                  <Form.Item name={'searchType'} initialValue={'spaceName'}>
                    <Select
                      options={searchOptions}
                      style={{ width: '120px' }}
                    />
                  </Form.Item>
                  <Form.Item
                    noStyle
                    shouldUpdate={(pre, cur) =>
                      pre.searchType !== cur.searchType
                    }
                  >
                    {({ getFieldValue }) => {
                      const searchType = getFieldValue('searchType');
                      if (searchType === 'spaceName') {
                        return (
                          <Form.Item name={'name'}>
                            <Input
                              style={{ width: '220px' }}
                              placeholder="请输入空间名称"
                              onPressEnter={() => {
                                handlePageSearch();
                              }}
                              suffix={
                                <SearchOutlined
                                  onClick={() => {
                                    handlePageSearch();
                                  }}
                                />
                              }
                            />
                          </Form.Item>
                        );
                      }
                      return (
                        <Form.Item name={'member'}>
                          <Input
                            style={{ width: '220px' }}
                            placeholder="请输入空间成员"
                            onPressEnter={() => {
                              handlePageSearch();
                            }}
                            suffix={
                              <SearchOutlined
                                onClick={() => {
                                  handlePageSearch();
                                }}
                              />
                            }
                          />
                        </Form.Item>
                      );
                    }}
                  </Form.Item>
                </Space.Compact>

                <Button
                  type="primary"
                  onClick={() => {
                    setAddSpaceOpen(true);
                  }}
                >
                  新建空间
                </Button>
              </Space>
            }
          />
        </Form>

        <Alert
          message="可联系空间内具有读写权限的成员添加为空间成员"
          type="info"
          showIcon
          closable
        />
        <div>
          <Spin spinning={getSpaceList?.loading}>
            {pageData?.namespaces?.length > 0 ? (
              <>
                <SpaceList pageData={pageData} handleDelete={handleDelete} />
                <Pagination
                  showQuickJumper
                  pageSize={pageData.pageSize}
                  current={pageData.page}
                  total={pageData.total}
                  onChange={(page, pageSize) => {
                    handlePageSearch({ page, pageSize });
                  }}
                />
              </>
            ) : (
              <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
            )}
          </Spin>
        </div>
        {addSpaceOpen && (
          <AddSpaceDrawer open={addSpaceOpen} setOpen={setAddSpaceOpen} />
        )}
      </Container>
    </PageContainer>
  );
};

export default SpaceManage;
