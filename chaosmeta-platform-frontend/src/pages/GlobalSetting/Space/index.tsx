import AddSpaceDrawer from '@/components/AddSpaceDrawer';
import {
  deleteSpace,
  queryClassSpaceList,
} from '@/services/chaosmeta/SpaceController';
import { ExclamationCircleFilled, SearchOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { useIntl, useRequest } from '@umijs/max';
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
    pageSize: 12,
    total: 0,
    namespaces: [],
  });
  const intl = useIntl();

  /**
   * 获取空间列表接口
   */
  const getSpaceList = useRequest(queryClassSpaceList, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
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
    let userName, spaceName;
    if (searchType === 'spaceMember') {
      userName = member;
      spaceName = undefined;
    }
    if (searchType === 'spaceName') {
      userName = undefined;
      spaceName = name;
    }
    const params = {
      page: page || pageData.page || 1,
      page_size: pageSize || pageData.pageSize || 12,
      sort,
      name: spaceName,
      namespaceClass: namespaceClass || spaceType,
      userName: userName,
    };
    getSpaceList?.run(params);
  };

  const handleDeleteSpace = useRequest(deleteSpace, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success(
          intl.formatMessage({ id: 'spaceManagement.delete.success' }),
        );
        handlePageSearch();
      }
    },
  });

  const tabItems = [
    {
      label: intl.formatMessage({ id: 'spaceManagement.tab.all' }),
      key: 'all',
    },
    {
      label: intl.formatMessage({ id: 'spaceManagement.tab.related' }),
      key: 'relevant',
    },
  ];

  /**
   * 删除账号
   */
  const handleDelete = (id: number) => {
    // return
    Modal.confirm({
      title: intl.formatMessage({ id: 'spaceManagement.delete.title' }),
      icon: <ExclamationCircleFilled />,
      onOk() {
        return handleDeleteSpace?.run({ id });
      },
    });
  };

  // 空间权限类型
  const spaceTypes = [
    {
      label: intl.formatMessage({ id: 'all' }),
      value: 'all',
    },
    // {
    //   label: '未加入',
    //   value: 'not',
    // },
    {
      label: intl.formatMessage({ id: 'readonly' }),
      value: 'read',
    },
    {
      label: intl.formatMessage({ id: 'write' }),
      value: 'write',
    },
  ];

  // 检索类型
  const searchOptions = [
    {
      label: intl.formatMessage({ id: 'spaceName' }),
      value: 'spaceName',
    },
    {
      label: intl.formatMessage({ id: 'spaceManagement.member' }),
      value: 'spaceMember',
    },
  ];

  useEffect(() => {
    handlePageSearch();
  }, []);

  return (
    <PageContainer title={intl.formatMessage({ id: 'spaceManagement.title' })}>
      <Container>
        <Form form={form}>
          <Tabs
            items={tabItems}
            activeKey={tabKey}
            onChange={(val) => {
              setTabKey(val);
              const namespaceClass =
                val === 'relevant' ? 'relevant' : undefined;
              handlePageSearch({ namespaceClass });
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
                              placeholder={intl.formatMessage({
                                id: 'spaceManagement.spaceName.placeholder',
                              })}
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
                            placeholder={intl.formatMessage({
                              id: 'spaceManagement.spaceMember.placeholder',
                            })}
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
                  {intl.formatMessage({ id: 'createSpace' })}
                </Button>
              </Space>
            }
          />
        </Form>

        <Alert
          message={intl.formatMessage({ id: 'spaceManagement.alert' })}
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
                  showSizeChanger
                  pageSize={pageData.pageSize}
                  current={pageData.page}
                  total={pageData.total}
                  pageSizeOptions={[12, 20, 50, 100]}
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
