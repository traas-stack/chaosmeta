import { LightArea } from '@/components/CommonStyle';
import { InfoCircleOutlined } from '@ant-design/icons';
import { Col, Row, Tag } from 'antd';
import React, { useState } from 'react';
import { RecommendContainer } from './style';
/**
 * 推荐实验
 * @returns
 */
const RecommendExperiment: React.FC<unknown> = () => {
  const [tempKey, setTempKey] = useState<string>('all');
  const [groupKey, setGroupKey] = useState('rong');
  const tempList = [
    {
      label: '所有模版',
      key: 'all',
    },
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

  const result = [1, 2, 3, 4];
  // const result = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15];

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

  /**
   * tab栏类
   * @param index
   * @param key
   * @returns
   */
  const groupItemClass = (index: number, key: string) => {
    let name = 'tab-item';
    if (groupKey === key) {
      // 第一个选中只会有右侧圆角及边框
      if (index === 0) {
        name = `${name} tab-item-active tab-item-active-after tab-item-active-border-right`;
        // 最后一个选中只会有左侧圆角及边框
      } else if (index === tabItems.length - 1) {
        name = `${name} tab-item-active tab-item-active-before tab-item-active-border-left`;
        // 其他情况选中左右两侧都需要圆角及边框
      } else {
        name = `${name} tab-item-active tab-item-active-before tab-item-active-after tab-item-active-border-left tab-item-active-border-right`;
      }
    }
    return name;
  };

  /**
   * 竖线的展示时机
   * @param index
   * @param key
   * @returns
   */
  const groupContentClass = (index: number) => {
    let name = 'tab-item-content';
    let groupKeyIndex = -1;
    tabItems.forEach((item, i) => {
      if (item.key === groupKey) {
        groupKeyIndex = i;
      }
    });
    // 第一个tab不会有竖线，选中的和选中的下一个tab也不需要竖线
    if (index === 0 || groupKeyIndex === index || groupKeyIndex + 1 === index) {
      name = '';
    }
    return name;
  };
  return (
    <RecommendContainer>
      <LightArea>
        <div className="recommend">
          <div className="left">
            <div>
              <Row className="tab">
                {tabItems?.map((item, index) => {
                  return (
                    <>
                      <Col
                        key={item.key}
                        className={groupItemClass(index, item.key)}
                        span={8}
                        onClick={() => {
                          setGroupKey(item.key);
                        }}
                      >
                        <span className={groupContentClass(index)}>
                          {item.label}
                        </span>
                      </Col>
                    </>
                  );
                })}
              </Row>
            </div>
            <div className="group">
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
          </div>
          <div className="right">
            <div className="tip">
              <InfoCircleOutlined
                style={{ color: '#1677ff', paddingRight: '8px' }}
              />
              搜索“{'xxx'}”，
              {/* <span>未找到相关结果！</span> */}
              <span>共为您找到以下相关结果 {'xx'} 条</span>
            </div>
            {/* <div style={{marginTop: '160px'}}>
            <Empty
              description="暂无相关实验模版，试试换一下关键词哦！"
              image="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*IVXqQJVKr04AAAAAAAAAAAAADmKmAQ/original"
            />
          </div> */}
            <div>
              <div className="title">Node</div>
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
        </div>
      </LightArea>
    </RecommendContainer>
  );
};

export default RecommendExperiment;
