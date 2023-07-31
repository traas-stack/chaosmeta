import { Area } from '@/components/CommonStyle';
import { DownOutlined, RightOutlined, UpOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { Button, Card, Col, Divider, Row, Select, Space, Tag } from 'antd';
import React, { useState } from 'react';
import ExperimentalOverview from './ExperimentalOverview';
import { Container, SpaceContent, TopStep } from './style';

const MySpace: React.FC<unknown> = () => {
  const [panelState, setPanelState] = useState<boolean>(true);
  return (
    <Container>
      <PageContainer title="工作台">
        <Area>
          <TopStep>
            <div className="panel">
              <div className="title">开始您的实验，只需要3步！</div>
              <div className="panel-state">
                {panelState ? (
                  <Space>
                    <span>收起</span>
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
                    <span>展开</span>
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
            {panelState && (
              <Row gutter={16} className="card">
                <Col span={8}>
                  <Card>
                    <Space>
                      <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*h_acR7jTCrgAAAAAAAAAAAAADmKmAQ/original" />
                      <div>
                        <div className="title">创建实验</div>
                        <div className="desc">
                          可选择实验模版快速构建实验场景，进行基础资源，如cpu燃烧等实验来验证应用系统的可靠性
                        </div>
                        <Space className="buttons">
                          <Button type="primary">创建实验</Button>
                          <Button>实验模版</Button>
                        </Space>
                      </div>
                    </Space>
                  </Card>
                </Col>
                <Col span={8}>
                  <Card>
                    <Space>
                      <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*MelqSodcfO8AAAAAAAAAAAAADmKmAQ/original" />
                      <div>
                        <div className="title">执行实验</div>
                        <div className="desc">针对配置好的实验可发起攻击</div>
                      </div>
                    </Space>
                  </Card>
                </Col>
                <Col span={8}>
                  <Card>
                    <Space>
                      <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*in2BQ4sjkicAAAAAAAAAAAAADmKmAQ/original" />
                      <div>
                        <div className="title">查看实验结果</div>
                        <div className="desc">
                          实验过程中可观测系统指标，实验完成后可查看实验结果，系统会自动度量
                        </div>
                        <Space className="buttons">
                          <Button type="primary">查看实验结果</Button>
                        </Space>
                      </div>
                    </Space>
                  </Card>
                </Col>
              </Row>
            )}
          </TopStep>
        </Area>
        <SpaceContent>
          <Row gutter={16}>
            <Col span={16} className="left">
              <Area>
                <ExperimentalOverview />
              </Area>
            </Col>
            <Col span={8} className="right">
              <Area className="overview">
                <div className="top">
                  <span className="title">空间总览</span>
                  <Select
                    defaultValue={'7day'}
                    options={[{ label: '最近7天', value: '7day' }]}
                  />
                </div>
                <Card>
                  <div className="result">
                    <div>
                      <span className="count">19</span>
                      <div className="shallow-65">新增实验</div>
                    </div>
                    <div style={{ position: 'relative' }}>
                      <Divider type="vertical" />
                      <span className="count">36</span>
                      <div className="shallow-65">执行实验</div>
                    </div>
                    <div>
                      <span className="count-error">9</span>
                      <div className="shallow-65">执行失败</div>
                    </div>
                  </div>
                </Card>
              </Area>
              <Area className="recommend">
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
              </Area>
            </Col>
          </Row>
        </SpaceContent>
      </PageContainer>
    </Container>
  );
};

export default MySpace;
