import { PlusOutlined, SearchOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { Col, Form, Input, Row, Tabs, Tag } from 'antd';
import { useState } from 'react';
import { Container } from './style';

const AddExperiment = () => {
  const [tabKey, setTabKey] = useState<string>('rong');
  const [tempKey, setTempKey] = useState<string>('node');

  const tabItems = [
    {
      label: '容器',
      key: 'rong',
    },
    {
      label: '物理机',
      key: 'wuli',
    },
    {
      label: 'K8s',
      key: 'K8s',
    },
  ];

  const tempList = [
    // {
    //   label: '所有模版',
    //   key: 'all',
    // },
    {
      label: 'Node',
      key: 'node',
    },
    {
      label: 'Pod',
      key: 'pod',
    },
    {
      label: 'Containe',
      key: 'containe',
    },
    {
      label: 'Ngnix',
      key: 'ngnix',
    },
  ];
  const result = [1, 2, 3, 4, 5, 6, 7, 8, 9];

  return (
    <PageContainer
      header={{
        title: '创建实验',
        onBack: () => {
          history.push('/space/experiment');
        },
      }}
    >
      <Container>
        <div>自定义创建</div>
        <div
          className="custom"
          onClick={() => {
            history?.push('/space/experiment/add');
          }}
        >
          <div>
            <PlusOutlined />
            <div>自定义创建</div>
          </div>
        </div>
        <div>使用推荐实验创建</div>
        <Form>
          <Tabs
            items={tabItems}
            activeKey={tabKey}
            onChange={(val) => {
              setTabKey(val);
            }}
            tabBarExtraContent={
              <>
                <Form.Item name={'member'}>
                  <Input
                    style={{ width: '220px' }}
                    placeholder="请输入搜索关键字"
                    onPressEnter={() => {
                      // handlePageSearch();
                    }}
                    suffix={
                      <SearchOutlined
                        onClick={() => {
                          // handlePageSearch();
                        }}
                      />
                    }
                  />
                </Form.Item>
              </>
            }
          />
          <div className="content">
            <div className="left">
              {tempList?.map((item) => {
                return (
                  <div
                    key={item.key}
                    onClick={() => {
                      setTempKey(item.key);
                    }}
                    className={
                      tempKey === item.key
                        ? 'group-item-active group-item'
                        : 'group-item'
                    }
                  >
                    {item.label}
                  </div>
                );
              })}
            </div>
            <div className="right">
              <Row gutter={[24, 24]}>
                {result?.map((item, index) => {
                  return (
                    <Col key={index} span={8}>
                      <div className="result-item">
                        <img
                          style={{ width: '100%' }}
                          src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*fYqDRaaaj7gAAAAAAAAAAAAADmKmAQ/original"
                          alt=""
                        />
                        <div className="introduce">
                          <img
                            src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*8jDfQoFH9NcAAAAAAAAAAAAADmKmAQ/original"
                            alt=""
                          />
                          <div>
                            <div>容器-docker-service-kill</div>
                            <div>
                              <Tag>我是标签</Tag>
                              <Tag>标签</Tag>
                            </div>
                          </div>
                        </div>
                      </div>
                    </Col>
                  );
                })}
              </Row>
            </div>
          </div>
        </Form>
      </Container>
    </PageContainer>
  );
};

export default AddExperiment;
