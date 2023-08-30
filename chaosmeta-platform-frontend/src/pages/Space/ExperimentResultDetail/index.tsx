import { PageContainer } from '@ant-design/pro-components';
import { history, useModel, useRequest } from '@umijs/max';
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
import { arrangeDataOriginTranstion, formatDuration } from '@/utils/format';
import { ExclamationCircleFilled } from '@ant-design/icons';
import ArrangeInfoShow from '../ExperimentDetail/ArrangeInfoShow';
import ShowLog from './ShowLog';
import { Container } from './style';

const tst = [
  {
    name: 'cpu使用率',
    uuid: 'daaef1a0431011ee94dfbd56395b052b1',
    experiment_uuid: '',
    row: 1,
    column: 0,
    duration: '90s',
    scope_id: 1,
    target_id: 1,
    exec_type: 'fault',
    exec_id: 1,
    create_time: '0001-01-01T00:00:00Z',
    update_time: '0001-01-01T00:00:00Z',
    args_value: [
      {
        id: 68,
        args_id: 1,
        workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
        value: '4',
        create_time: '2023-08-25T15:15:23+08:00',
        update_time: '2023-08-25T15:15:23+08:00',
      },
      {
        id: 69,
        args_id: 2,
        workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
        value: '',
        create_time: '2023-08-25T15:15:23+08:00',
        update_time: '2023-08-25T15:15:23+08:00',
      },
    ],
    exec_range: {
      id: 28,
      workflow_node_instance_uuid: 'daaef1a0431011ee94dfbd56395b052b',
      target_name: '',
      target_ip: '',
      target_hostname: '',
      target_label: '',
      target_app: '',
      target_namespace: '3',
      range_type: '',
      create_time: '2023-08-25T15:15:23+08:00',
      update_time: '2023-08-25T15:15:23+08:00',
    },
  },
  {
    name: 'cpu使用率',
    uuid: 'daaef1a0431011ee94dfbd56395b052b',
    experiment_uuid: '',
    row: 3,
    column: 0,
    duration: '90s',
    scope_id: 1,
    target_id: 1,
    exec_type: 'fault',
    exec_id: 1,
    create_time: '0001-01-01T00:00:00Z',
    update_time: '0001-01-01T00:00:00Z',
    args_value: [
      {
        id: 68,
        args_id: 1,
        workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
        value: '4',
        create_time: '2023-08-25T15:15:23+08:00',
        update_time: '2023-08-25T15:15:23+08:00',
      },
      {
        id: 69,
        args_id: 2,
        workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
        value: '',
        create_time: '2023-08-25T15:15:23+08:00',
        update_time: '2023-08-25T15:15:23+08:00',
      },
    ],
    exec_range: {
      id: 28,
      workflow_node_instance_uuid: 'daaef1a0431011ee94dfbd56395b052b',
      target_name: '',
      target_ip: '',
      target_hostname: '',
      target_label: '',
      target_app: '',
      target_namespace: '2023-08-25T15:15:23+08:002023-08-25T15:15:23+08:00',
      range_type: '',
      create_time: '2023-08-25T15:15:23+08:00',
      update_time: '2023-08-25T15:15:23+08:00',
    },
  },
  // {
  //   name: 'cpu使用率',
  //   uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //   experiment_uuid: '',
  //   row: 3,
  //   column: 0,
  //   duration: '90s',
  //   scope_id: 1,
  //   target_id: 1,
  //   exec_type: 'fault',
  //   exec_id: 1,
  //   create_time: '0001-01-01T00:00:00Z',
  //   update_time: '0001-01-01T00:00:00Z',
  //   args_value: [
  //     {
  //       id: 68,
  //       args_id: 1,
  //       workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //       value: '4',
  //       create_time: '2023-08-25T15:15:23+08:00',
  //       update_time: '2023-08-25T15:15:23+08:00',
  //     },
  //     {
  //       id: 69,
  //       args_id: 2,
  //       workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //       value: '',
  //       create_time: '2023-08-25T15:15:23+08:00',
  //       update_time: '2023-08-25T15:15:23+08:00',
  //     },
  //   ],
  //   exec_range: {
  //     id: 28,
  //     workflow_node_instance_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //     target_name: '',
  //     target_ip: '',
  //     target_hostname: '',
  //     target_label: '',
  //     target_app: '',
  //     target_namespace: '3',
  //     range_type: '',
  //     create_time: '2023-08-25T15:15:23+08:00',
  //     update_time: '2023-08-25T15:15:23+08:00',
  //   },
  // },
  // {
  //   name: 'cpu使用率',
  //   uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //   experiment_uuid: '',
  //   row: 4,
  //   column: 0,
  //   duration: '90s',
  //   scope_id: 1,
  //   target_id: 1,
  //   exec_type: 'fault',
  //   exec_id: 1,
  //   create_time: '0001-01-01T00:00:00Z',
  //   update_time: '0001-01-01T00:00:00Z',
  //   args_value: [
  //     {
  //       id: 68,
  //       args_id: 1,
  //       workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //       value: '4',
  //       create_time: '2023-08-25T15:15:23+08:00',
  //       update_time: '2023-08-25T15:15:23+08:00',
  //     },
  //     {
  //       id: 69,
  //       args_id: 2,
  //       workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //       value: '',
  //       create_time: '2023-08-25T15:15:23+08:00',
  //       update_time: '2023-08-25T15:15:23+08:00',
  //     },
  //   ],
  //   exec_range: {
  //     id: 28,
  //     workflow_node_instance_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //     target_name: '',
  //     target_ip: '',
  //     target_hostname: '',
  //     target_label: '',
  //     target_app: '',
  //     target_namespace: '3',
  //     range_type: '',
  //     create_time: '2023-08-25T15:15:23+08:00',
  //     update_time: '2023-08-25T15:15:23+08:00',
  //   },
  // },
  // {
  //   name: 'cpu使用率',
  //   uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //   experiment_uuid: '',
  //   row: 5,
  //   column: 0,
  //   duration: '90s',
  //   scope_id: 1,
  //   target_id: 1,
  //   exec_type: 'fault',
  //   exec_id: 1,
  //   create_time: '0001-01-01T00:00:00Z',
  //   update_time: '0001-01-01T00:00:00Z',
  //   args_value: [
  //     {
  //       id: 68,
  //       args_id: 1,
  //       workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //       value: '4',
  //       create_time: '2023-08-25T15:15:23+08:00',
  //       update_time: '2023-08-25T15:15:23+08:00',
  //     },
  //     {
  //       id: 69,
  //       args_id: 2,
  //       workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //       value: '',
  //       create_time: '2023-08-25T15:15:23+08:00',
  //       update_time: '2023-08-25T15:15:23+08:00',
  //     },
  //   ],
  //   exec_range: {
  //     id: 28,
  //     workflow_node_instance_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //     target_name: '',
  //     target_ip: '',
  //     target_hostname: '',
  //     target_label: '',
  //     target_app: '',
  //     target_namespace: '3',
  //     range_type: '',
  //     create_time: '2023-08-25T15:15:23+08:00',
  //     update_time: '2023-08-25T15:15:23+08:00',
  //   },
  // },
  // {
  //   name: 'cpu使用率',
  //   uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //   experiment_uuid: '',
  //   row: 6,
  //   column: 0,
  //   duration: '90s',
  //   scope_id: 1,
  //   target_id: 1,
  //   exec_type: 'fault',
  //   exec_id: 1,
  //   create_time: '0001-01-01T00:00:00Z',
  //   update_time: '0001-01-01T00:00:00Z',
  //   args_value: [
  //     {
  //       id: 68,
  //       args_id: 1,
  //       workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //       value: '4',
  //       create_time: '2023-08-25T15:15:23+08:00',
  //       update_time: '2023-08-25T15:15:23+08:00',
  //     },
  //     {
  //       id: 69,
  //       args_id: 2,
  //       workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //       value: '',
  //       create_time: '2023-08-25T15:15:23+08:00',
  //       update_time: '2023-08-25T15:15:23+08:00',
  //     },
  //   ],
  //   exec_range: {
  //     id: 28,
  //     workflow_node_instance_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //     target_name: '',
  //     target_ip: '',
  //     target_hostname: '',
  //     target_label: '',
  //     target_app: '',
  //     target_namespace: '3',
  //     range_type: '',
  //     create_time: '2023-08-25T15:15:23+08:00',
  //     update_time: '2023-08-25T15:15:23+08:00',
  //   },
  // },
  // {
  //   name: 'cpu使用率',
  //   uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //   experiment_uuid: '',
  //   row: 7,
  //   column: 0,
  //   duration: '90s',
  //   scope_id: 1,
  //   target_id: 1,
  //   exec_type: 'fault',
  //   exec_id: 1,
  //   create_time: '0001-01-01T00:00:00Z',
  //   update_time: '0001-01-01T00:00:00Z',
  //   args_value: [
  //     {
  //       id: 68,
  //       args_id: 1,
  //       workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //       value: '4',
  //       create_time: '2023-08-25T15:15:23+08:00',
  //       update_time: '2023-08-25T15:15:23+08:00',
  //     },
  //     {
  //       id: 69,
  //       args_id: 2,
  //       workflow_node_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //       value: '',
  //       create_time: '2023-08-25T15:15:23+08:00',
  //       update_time: '2023-08-25T15:15:23+08:00',
  //     },
  //   ],
  //   exec_range: {
  //     id: 28,
  //     workflow_node_instance_uuid: 'daaef1a0431011ee94dfbd56395b052b',
  //     target_name: '',
  //     target_ip: '',
  //     target_hostname: '',
  //     target_label: '',
  //     target_app: '',
  //     target_namespace: '3',
  //     range_type: '',
  //     create_time: '2023-08-25T15:15:23+08:00',
  //     update_time: '2023-08-25T15:15:23+08:00',
  //   },
  // },
];

const AddExperiment = () => {
  // 编排的数据
  const [arrangeList, setArrangeList] = useState([]);
  // 用户权限
  const { spacePermission } = useModel('global');
  const [tabKey, setTabKey] = useState<'log' | 'index' | string>('log');
  const curExecSecond = '180s';
  // 单个节点详情
  const [curNodeDetail, setCurNodeDetail] = useState<any>({});
  // 结果详情
  const [resultDetail, setResultDetail] = useState<any>({});

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
        message.success(`${resultDetail?.name}实验已停止`);
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
      title: '确认要停止实验吗？',
      icon: <ExclamationCircleFilled />,
      onOk() {
        return stopExperiment?.run({ uuid: resultDetail?.uuid });
      },
    });
  };

  const headerExtra = () => {
    return (
      <Space>
        <Button>查看实验配置</Button>
        {resultDetail?.status === 'Running' && (
          <Button
            type="primary"
            onClick={() => {
              handleDeleteAccount();
            }}
          >
            停止
          </Button>
        )}
      </Space>
    );
  };

  const items: TabsProps['items'] = [
    {
      key: 'log',
      label: `实验日志`,
      children: <ShowLog message={curNodeDetail?.message || ''} />,
    },
    // todot -- 后端暂不支持
    // {
    //   key: 'index',
    //   label: `实验观测指标`,
    //   children: <ObservationCharts />,
    // },
  ];

  /**
   * 当前状态匹配
   */
  const handleMateStatus = () => {
    const temp = experimentResultStatus?.filter(
      (item) => item?.value === resultDetail?.status,
    )[0];
    return temp;
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
        <Tag color={handleMateStatus()?.color}>{handleMateStatus()?.label}</Tag>
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
            <div>实验进度</div>
            {/* 后端不支持展示进度，只有成功展示进度条，其他情况展示当前状态 */}
            {resultDetail?.status === 'Succeeded' ? (
              <Progress percent={100} size="small" />
            ) : (
              <span>
                <Badge color={handleMateStatus()?.color} />{' '}
                {handleMateStatus()?.label}
              </span>
            )}
          </div>
          {resultDetail?.status &&
            resultDetail?.status !== 'Running' &&
            resultDetail?.status !== 'Pending' && (
              <Alert
                message={
                  <>{`${statusText[resultDetail?.status]}${
                    resultDetail?.message || ''
                  }`}</>
                }
                style={{ marginBottom: '16px' }}
                type={handleMateStatus()?.type}
                action={
                  (resultDetail?.status === 'error' ||
                    resultDetail?.status === 'Failed') && (
                    <Button type="link" onClick={() => {}}>
                      查看详情
                    </Button>
                  )
                }
                showIcon
              />
            )}

          {/* 编排信息的展示 */}
          <ArrangeInfoShow
            arrangeList={arrangeList}
            curExecSecond={formatDuration(curExecSecond)}
            isResult
            getExperimentArrangeNodeDetail={getExperimentArrangeNodeDetail}
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
