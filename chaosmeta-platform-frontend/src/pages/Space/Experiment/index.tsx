import { PageContainer } from '@ant-design/pro-components';
import { Tabs } from 'antd';
import React from 'react';
import ExperimentList from './ExperimentList';
// import RecommendExperiment from './RecommendExperiment';
import { Container } from './style';
/**
 * 实验列表页面
 * @returns
 */
const Experiment: React.FC<unknown> = () => {
  const tabItems = [
    {
      label: '实验列表',
      key: 'list',
      children: <ExperimentList />,
    },
    // 一期暂时隐藏
    // {
    //   label: '推荐实验',
    //   key: '',
    //   children: <RecommendExperiment />,
    // },
  ];
  
  return (
    <>
      <PageContainer title="实验">
        <Container>
          <Tabs items={tabItems} />
        </Container>
      </PageContainer>
    </>
  );
};

export default Experiment;
