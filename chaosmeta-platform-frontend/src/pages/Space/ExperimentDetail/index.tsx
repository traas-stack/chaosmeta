import { PageContainer } from '@ant-design/pro-components';
import { history, useModel } from '@umijs/max';
import { Button, Descriptions, Space, Spin, message } from 'antd';
import { useEffect, useState } from 'react';
// import ArrangeContent from './ArrangeContent';
// import InfoDrawer from './components/InfoDrawer';
import {
  createExperiment,
  queryExperimentDetail,
  runExperiment,
} from '@/services/chaosmeta/ExperimentController';
import {
  arrangeDataOriginTranstion,
  copyExperimentFormatData,
  formatTime,
} from '@/utils/format';
import { renderScheduleType, renderTags } from '@/utils/renderItem';
import { useRequest } from '@umijs/max';
import ArrangeInfoShow from './ArrangeInfoShow';
import { Container } from './style';

const AddExperiment = () => {
  // 编排的数据
  const [arrangeList, setArrangeList] = useState([]);
  // 用户权限
  const { spacePermission } = useModel('global');
  const [baseInfo, setBaseInfo] = useState<any>({});

  /**
   * 获取实验详情
   */
  const getExperimentDetail = useRequest(queryExperimentDetail, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        const experiments = res?.data?.experiments;
        // 将args_value转换为form可以使用的
        const newList = experiments?.workflow_nodes?.map((item: any) => {
          const newArgs: any = {};
          item?.args_value?.forEach((arg: any) => {
            newArgs[arg?.args_id] = arg?.value;
          });
          return { ...item, args_value: newArgs };
        });
        setBaseInfo(experiments);
        setArrangeList(arrangeDataOriginTranstion(newList || [], true));
      }
    },
  });

  /**
   * 运行试验
   */
  const handleRunExperiment = useRequest(runExperiment, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success('开始运行实验');
        getExperimentDetail?.run({
          uuid: history?.location?.query?.experimentId as string,
        });
      }
    },
  });

  /**
   * 创建实验 -- 复制实验使用
   */
  const handleCreateExperiment = useRequest(createExperiment, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success('复制成功');
        history.push({
          pathname: '/space/experiment/detail',
          query: {
            ...history?.location?.query,
            experimentId: res?.data?.uuid,
          },
        });
      }
    },
  });

  /**
   * 复制实验
   */
  const handleCopyExperiment = () => {
    const params = copyExperimentFormatData(baseInfo);
    handleCreateExperiment?.run(params);
  };

  const headerExtra = () => {
    if (spacePermission === 1) {
      return (
        <Space>
          <Button
            loading={handleCreateExperiment?.loading}
            onClick={() => {
              handleCopyExperiment();
            }}
          >
            复制
          </Button>
          <Button
            onClick={() => {
              history.push({
                pathname: '/space/experiment-result',
                query: {
                  experimentId: history?.location?.query?.experimentId,
                  spaceId: history?.location?.query?.spaceId as string,
                },
              });
            }}
          >
            实验结果
          </Button>
          <Button
            onClick={() => {
              history.push({
                pathname: '/space/experiment/add',
                query: history?.location?.query,
              });
            }}
          >
            编辑
          </Button>
          {/* 手动时才展示 */}
          {baseInfo?.schedule_type === 'manual' && (
            <Button
              type="primary"
              // status === 3 为运行中
              loading={handleRunExperiment?.loading || baseInfo?.status === 3}
              onClick={() => {
                handleRunExperiment?.run({
                  uuid: history?.location?.query?.experimentId as string,
                });
              }}
            >
              运行
            </Button>
          )}
        </Space>
      );
    }
    return (
      <Button
        onClick={() => {
          history.push({
            pathname: '/space/experiment-result',
            query: {
              experimentId: history?.location?.query?.experimentId,
              spaceId: history?.location?.query?.spaceId as string,
            },
          });
        }}
      >
        实验结果
      </Button>
    );
  };

  useEffect(() => {
    const { experimentId } = history?.location?.query || {};
    if (experimentId) {
      getExperimentDetail?.run({ uuid: experimentId as string });
    } else {
      setArrangeList(arrangeDataOriginTranstion([]));
    }
  }, [history?.location?.query?.experimentId]);

  return (
    <Container>
      <PageContainer
        header={{
          title: baseInfo?.name || '',
          onBack: () => {
            history.back();
          },
          extra: headerExtra(),
        }}
      >
        <Spin spinning={getExperimentDetail?.loading}>
          <div className="content">
            <Descriptions title="基本信息">
              <Descriptions.Item label="创建人">
                {baseInfo?.creator_name}
              </Descriptions.Item>
              <Descriptions.Item label="最近操作时间">
                {formatTime(baseInfo?.update_time)}
              </Descriptions.Item>
              <Descriptions.Item label="标签">
                {renderTags(baseInfo?.labels) || '--'}
              </Descriptions.Item>
              <Descriptions.Item label="触发方式">
                {renderScheduleType(baseInfo)}
              </Descriptions.Item>
              <Descriptions.Item label="描述">
                {baseInfo?.description}
              </Descriptions.Item>
            </Descriptions>
            <div className="experiment">
              <div className="experiment-title">实验配置</div>
              <ArrangeInfoShow arrangeList={arrangeList} />
            </div>
          </div>
        </Spin>
      </PageContainer>
    </Container>
  );
};

export default AddExperiment;
