import { LightArea } from '@/components/CommonStyle';
import { DownOutlined, RightOutlined, UpOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { history, useIntl, useModel } from '@umijs/max';
import { Button, Card, Col, Row, Space, Tag } from 'antd';
import React, { useState } from 'react';
import ExperimentalOverview from './ExperimentalOverview';
import { Container, SpaceContent, TopStep } from './style';

const MySpace: React.FC<unknown> = () => {
  const [panelState, setPanelState] = useState<boolean>(true);
  const { spacePermission } = useModel('global');
  const intl = useIntl();
  return (
    <Container>
      <PageContainer title={intl.formatMessage({ id: 'overview.workbench' })}>
        <LightArea>
          <TopStep>
            <div className="panel">
              <div className="title">
                {intl.formatMessage({ id: 'overview.tip' })}
              </div>
              <div className="panel-state">
                {panelState ? (
                  <Space>
                    <span>
                      {intl.formatMessage({ id: 'overview.panel.close' })}
                    </span>
                    <div
                      className="icon"
                      onClick={() => {
                        setPanelState(!panelState);
                      }}
                    >
                      <UpOutlined />
                    </div>
                  </Space>
                ) : (
                  <Space>
                    <span>
                      {intl.formatMessage({ id: 'overview.panel.expand' })}
                    </span>
                    <div
                      className="icon"
                      onClick={() => {
                        setPanelState(!panelState);
                      }}
                    >
                      <DownOutlined style={{ marginTop: '3px' }} />
                    </div>
                  </Space>
                )}
              </div>
            </div>
            {/* {panelState && ( */}
            <Row gutter={16} className={panelState ? 'card' : 'card-hidden'}>
              <Col span={8}>
                <Card>
                  <Space>
                    <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*h_acR7jTCrgAAAAAAAAAAAAADmKmAQ/original" />
                    <div>
                      <div className="title">
                        {intl.formatMessage({ id: 'overview.step1.title' })}
                      </div>
                      <div className="desc">
                        {intl.formatMessage({
                          id: 'overview.step1.description',
                        })}
                      </div>
                      {spacePermission === 1 && (
                        <Space className="buttons">
                          <Button
                            type="primary"
                            onClick={() => {
                              history.push({
                                pathname: '/space/experiment/add',
                                query: {
                                  spaceId: history?.location?.query?.spaceId,
                                },
                              });
                            }}
                          >
                            {intl.formatMessage({ id: 'overview.step1.title' })}
                          </Button>
                          {/* <Button>实验模版</Button> */}
                        </Space>
                      )}
                    </div>
                  </Space>
                </Card>
              </Col>
              <Col span={8}>
                <Card>
                  <Space>
                    <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*MelqSodcfO8AAAAAAAAAAAAADmKmAQ/original" />
                    <div>
                      <div className="title">
                        {intl.formatMessage({ id: 'overview.step2.title' })}
                      </div>
                      <div className="desc">
                        {intl.formatMessage({
                          id: 'overview.step2.description',
                        })}
                      </div>
                    </div>
                  </Space>
                </Card>
              </Col>
              <Col span={8}>
                <Card>
                  <Space>
                    <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*in2BQ4sjkicAAAAAAAAAAAAADmKmAQ/original" />
                    <div>
                      <div className="title">
                        {intl.formatMessage({ id: 'overview.step3.title' })}
                      </div>
                      <div className="desc">
                        {intl.formatMessage({
                          id: 'overview.step3.description',
                        })}
                      </div>
                      {spacePermission === 1 && (
                        <Space className="buttons">
                          <Button
                            type="primary"
                            onClick={() => {
                              history.push('/space/experiment-result');
                            }}
                          >
                            {intl.formatMessage({ id: 'overview.step3.title' })}
                          </Button>
                        </Space>
                      )}
                    </div>
                  </Space>
                </Card>
              </Col>
            </Row>
            {/* )} */}
          </TopStep>
        </LightArea>
        <SpaceContent>
          <Row gutter={16}>
            <Col span={24} className="left">
              <LightArea>
                <ExperimentalOverview />
              </LightArea>
            </Col>
            {/* 推荐实验暂时没有 */}
            <Col span={8} className="right" style={{ display: 'none' }}>
              <LightArea className="recommend">
                <div className="top">
                  <span className="title">推荐实验</span>
                  <Space className="shallow-65">
                    <span>查看更多</span>
                    <RightOutlined />
                  </Space>
                </div>
                <Card>
                  <div className="items">
                    <div className="item">
                      <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*h_acR7jTCrgAAAAAAAAAAAAADmKmAQ/original" />
                      <div>
                        <div>K8s-docker-service-kill</div>
                        <div>
                          <Tag>标签</Tag>
                          <Tag>标签</Tag>
                        </div>
                      </div>
                    </div>
                    <div className="item">
                      <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*h_acR7jTCrgAAAAAAAAAAAAADmKmAQ/original" />
                      <div>
                        <div>K8s-docker-service-kill</div>
                        <div>
                          <Tag>标签</Tag>
                          <Tag>标签</Tag>
                        </div>
                      </div>
                    </div>
                    {/* <div className="item">
                      <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*h_acR7jTCrgAAAAAAAAAAAAADmKmAQ/original" />
                      <div>
                        <div>K8s-docker-service-kill</div>
                        <div>
                          <Tag>标签</Tag>
                          <Tag>标签</Tag>
                        </div>
                      </div>
                    </div> */}
                  </div>
                </Card>
              </LightArea>
            </Col>
          </Row>
        </SpaceContent>
      </PageContainer>
    </Container>
  );
};

export default MySpace;
