import { PageContainer } from '@ant-design/pro-components';
import { getLocale, history, useIntl, useRequest } from '@umijs/max';
import {
  Alert,
  Badge,
  Button,
  Modal,
  Progress,
  Space,
  Tabs,
  TabsProps,
  Tag,
  message,
} from 'antd';
import { useEffect, useState } from 'react';
// import ArrangeContent from './ArrangeContent';
// import InfoDrawer from './components/InfoDrawer';
// import ArrangeInfoShow from './ArrangeInfoShow';
import { experimentResultStatus } from '@/constants';
import {
  queryExperimentResultArrangeNodeDetail,
  queryExperimentResultArrangeNodeList,
  queryExperimentResultDetail,
  stopExperimentResult,
} from '@/services/chaosmeta/ExperimentController';
import {
  arrangeDataOriginTranstion,
  formatDuration,
  getIntlLabel,
} from '@/utils/format';
import { ExclamationCircleFilled } from '@ant-design/icons';
import ArrangeInfoShow from '../ExperimentDetail/ArrangeInfoShow';
import ShowLog from './ShowLog';
import { Container } from './style';

const AddExperiment = () => {
  // 编排的数据
  const [arrangeList, setArrangeList] = useState([]);
  // 用户权限
  const [tabKey, setTabKey] = useState<'log' | 'index' | string>('log');
  const curExecSecond = '180s';
  // 单个节点详情
  const [curNodeDetail, setCurNodeDetail] = useState<any>({});
  // 结果详情
  const [resultDetail, setResultDetail] = useState<any>({});
  const intl = useIntl();

  /**
   * 获取实验结果详情
   */
  const getResultDetail = useRequest(queryExperimentResultDetail, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        setResultDetail(res?.data);
      }
    },
  });

  /**
   * 获取实验结果编排节点list
   */
  const getExperimentArrangeDetail = useRequest(
    queryExperimentResultArrangeNodeList,
    {
      manual: true,
      formatResult: (res) => res,
      onSuccess: (res) => {
        if (res?.code === 200) {
          setArrangeList(
            arrangeDataOriginTranstion(res?.data?.workflow_nodes || [], true),
          );
        }
      },
    },
  );

  /**
   * 获取实验结果单个编排节点
   */
  const getExperimentArrangeNodeDetail = useRequest(
    queryExperimentResultArrangeNodeDetail,
    {
      manual: true,
      formatResult: (res) => res,
      onSuccess: (res) => {
        if (res?.code === 200) {
          const data = res?.data?.workflow_node;
          setCurNodeDetail(data);
        }
      },
    },
  );

  /**
   * 停止实验
   */
  const stopExperiment = useRequest(stopExperimentResult, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success(
          `${resultDetail?.name}${intl.formatMessage({
            id: 'experimentResult.stop.text',
          })}`,
        );
        getResultDetail?.run({
          uuid: history?.location?.query?.resultId as string,
        });
      }
    },
  });

  /**
   * 停止实验
   */
  const handleDeleteAccount = () => {
    Modal.confirm({
      title: intl.formatMessage({ id: 'stopConfirmText' }),
      icon: <ExclamationCircleFilled />,
      onOk() {
        return stopExperiment?.run({ uuid: resultDetail?.uuid });
      },
    });
  };

  const headerExtra = () => {
    return (
      <Space>
        {/* <Button>查看实验配置</Button> */}
        {resultDetail?.status === 'Running' && (
          <Button
            type="primary"
            onClick={() => {
              handleDeleteAccount();
            }}
          >
            {intl.formatMessage({ id: 'stop' })}
          </Button>
        )}
      </Space>
    );
  };

  const items: TabsProps['items'] = [
    {
      key: 'log',
      label: intl.formatMessage({ id: 'experimentLog' }),
      children: <ShowLog message={curNodeDetail?.message || ''} />,
    },
    // todo -- 后端暂不支持
    // {
    //   key: 'index',
    //   label: `实验观测指标`,
    //   children: <ObservationCharts />,
    // },
  ];

  /**
   * 当前状态匹配
   */
  const handleMateStatus: any = () => {
    const temp = experimentResultStatus?.filter(
      (item) => item?.value === resultDetail?.status,
    )[0];
    return temp;
  };

  // 不同状态展示不同文案
  const statusTextUS: any = {
    Succeeded: 'The run is over and the experiment is successful.',
    Failed: 'The run ends and the experiment fails. Reason for failure: ',
    error: 'End of run, experiment error. wrong reason: ',
  };
  // 不同状态展示不同文案
  const statusText: any = {
    Succeeded: '运行结束，实验成功。',
    Failed: '运行结束，实验失败。失败原因：',
    error: '运行结束，实验错误。错误原因：',
  };

  const renderTitle = () => {
    return (
      <div>
        {resultDetail?.name}{' '}
        <Tag color={handleMateStatus()?.color}>
          {getIntlLabel(handleMateStatus())}
        </Tag>
      </div>
    );
  };

  useEffect(() => {
    const { resultId } = history?.location?.query || {};
    if (resultId) {
      getResultDetail?.run({ uuid: resultId as string });
      getExperimentArrangeDetail?.run({ uuid: resultId as string });
    } else {
      setArrangeList(arrangeDataOriginTranstion([], true));
    }
  }, []);

  return (
    <Container>
      <PageContainer
        header={{
          title: renderTitle(),
          onBack: () => {
            history.back();
          },
          extra: headerExtra(),
        }}
      >
        <div className="content">
          <div className="content-title">
            <div>{intl.formatMessage({ id: 'experimentProgress' })}</div>
            {/* 后端不支持展示进度，只有成功展示进度条，其他情况展示当前状态 */}
            {resultDetail?.status === 'Succeeded' ? (
              <Progress percent={100} size="small" />
            ) : (
              <span>
                <Badge color={handleMateStatus()?.color} />{' '}
                {getIntlLabel(handleMateStatus())}
              </span>
            )}
          </div>
          {resultDetail?.status &&
            resultDetail?.status !== 'Running' &&
            resultDetail?.status !== 'Pending' && (
              <Alert
                message={
                  <>{`${
                    (getLocale() === 'en-US' ? statusTextUS : statusText)[
                      resultDetail?.status
                    ]
                  }${resultDetail?.message || ''}`}</>
                }
                style={{ marginBottom: '16px' }}
                type={handleMateStatus()?.type}
                showIcon
              />
            )}

          {/* 编排信息的展示 */}
          <ArrangeInfoShow
            arrangeList={arrangeList}
            curExecSecond={formatDuration(curExecSecond)}
            isResult
            getExperimentArrangeNodeDetail={getExperimentArrangeNodeDetail}
            setCurNodeDetail={setCurNodeDetail}
          />
          {/* 日志信息 */}
          <div className="log">
            <Tabs
              defaultActiveKey="log"
              activeKey={tabKey}
              items={items}
              onChange={(key: string) => {
                setTabKey(key);
              }}
            />
          </div>
        </div>
      </PageContainer>
    </Container>
  );
};

export default AddExperiment;
