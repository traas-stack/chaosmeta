import ShowText from '@/components/ShowText';
import {
  createExperiment,
  deleteExperiment,
  queryExperimentDetail,
  updateExperiment,
} from '@/services/chaosmeta/ExperimentController';
import { querySpaceUserPermission } from '@/services/chaosmeta/SpaceController';
import {
  arrangeDataOriginTranstion,
  arrangeDataResultTranstion,
} from '@/utils/format';
import { renderScheduleType, renderTags } from '@/utils/renderItem';
import { EditOutlined, ExclamationCircleFilled } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { history, useModel, useRequest } from '@umijs/max';
import { Button, Form, Modal, Space, Spin, message } from 'antd';
import { useEffect, useState } from 'react';
import ArrangeContent from './ArrangeContent';
import InfoDrawer from './components/InfoDrawer';
import { Container } from './style';

const AddExperiment = () => {
  const [form] = Form.useForm();
  // 用户权限
  const { setSpacePermission, spacePermission } = useModel('global');
  // 编排的数据
  const [arrangeList, setArrangeList] = useState<any>([]);
  // 编辑基本信息抽屉
  const [infoDrawerOpen, setInfoDrawerOpen] = useState(false);
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
        // 已经保存过的信息，完善度设为true，已完善
        const newList = experiments?.workflow_nodes?.map((item: any) => {
          // 将动态表单args_value的值处理为form可以使用的
          const newArgs: any = {};
          item?.args_value?.forEach((arg: any) => {
            newArgs[arg?.args_id] = arg?.value;
          });
          return { ...item, nodeInfoState: true, args_value: newArgs };
        });
        form.setFieldsValue(experiments);
        setBaseInfo(experiments);
        setArrangeList(arrangeDataOriginTranstion(newList || []));
      }
    },
  });

  /**
   * 编辑更新实验
   */
  const editExperiment = useRequest(updateExperiment, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success('更新成功');
        history?.push('/space/experiment');
      }
    },
  });

  /**
   * 创建实验
   */
  const handleCreateExperiment = useRequest(createExperiment, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success('创建成功');
        history?.push('/space/experiment');
      }
    },
  });

  /**
   * 根据成员名称和空间id获取成员空间内权限信息
   */
  const getUserSpaceAuth = useRequest(querySpaceUserPermission, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res.code === 200) {
        // 存储用户空间权限
        setSpacePermission(res?.data);
      }
    },
  });

  // 标题渲染
  const renderTitle = () => {
    return (
      <Form form={form}>
        <Space>
          <Form.Item name={'name'}>
            <ShowText ellipsis />
          </Form.Item>
          <>
            {spacePermission === 1 ? (
              <EditOutlined
                className="edit"
                style={{ color: '#1890FF' }}
                onClick={() => {
                  setInfoDrawerOpen(true);
                }}
              />
            ) : (
              <a
                onClick={() => {
                  setInfoDrawerOpen(true);
                }}
              >
                查看
              </a>
            )}
          </>
        </Space>
        <Form.Item>{renderTags(baseInfo?.labels)}</Form.Item>
      </Form>
    );
  };

  /**
   * 提交实验信息
   */
  const handleSubmit = () => {
    form.validateFields().then((values) => {
      const arrangeResult = arrangeDataResultTranstion(arrangeList);
      if (!baseInfo?.name || !baseInfo?.schedule_type) {
        message.info('请完善基本信息');
        return;
      }
      if (
        !arrangeResult?.length ||
        arrangeResult?.some((item) => !item?.nodeInfoState)
      ) {
        message.info('请完善节点信息');
        return;
      }
      const newLabels = baseInfo?.labels?.map(
        (item: { id: number }) => item?.id,
      );
      const newList = arrangeResult?.map((item) => {
        const {
          args_value,
          exec_range,
          exec_id,
          row,
          column,
          uuid,
          duration,
          scope_id,
          target_id,
          exec_type,
          name,
        } = item;
        // return;
        let target_name = exec_range?.target_name;
        if (Array.isArray(target_name)) {
          target_name = exec_range?.target_name?.join(',');
        }
        const newExecRange = {
          ...exec_range,
          target_name: target_name || undefined,
        };

        return {
          name,
          args_value,
          exec_range: newExecRange,
          exec_id,
          row,
          column,
          uuid,
          duration,
          scope_id,
          target_id,
          exec_type,
        };
      });
      const params = {
        ...values,
        labels: newLabels,
        schedule_rule: baseInfo?.schedule_rule,
        namespace_id: Number(history?.location?.query?.spaceId),
        workflow_nodes: newList,
      };
      const experimentId = history?.location?.query?.experimentId;
      if (experimentId) {
        editExperiment?.run({ ...params, uuid: experimentId });
      } else {
        handleCreateExperiment?.run(params);
      }
    });
  };

  /**
   * 删除实验接口
   */
  const handleDeleteExperiment = useRequest(deleteExperiment, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success('删除成功！');
        history.push('/space/experiment');
      }
    },
  });

  /**
   * 确认删除实验
   */
  const handleDeleteConfirm = () => {
    const uuid = history?.location?.query?.experimentId as string;
    if (uuid) {
      Modal.confirm({
        title: '确认要删除这个实验吗？',
        icon: <ExclamationCircleFilled />,
        content: '删除实验将会删除该实验的配置，但不会删除历史实验结果！',
        onOk() {
          handleDeleteExperiment?.run({ uuid });
        },
      });
    }
  };

  const headerExtra = () => {
    return (
      <Form form={form}>
        <div className="header-extra">
          <div>
            <Form.Item name={'schedule_type'} label="触发方式">
              {renderScheduleType(baseInfo)}
            </Form.Item>
            <Form.Item name={'description'} label="描述">
              <ShowText />
            </Form.Item>
          </div>
          {spacePermission === 1 && (
            <Space>
              <Button
                ghost
                danger
                onClick={() => {
                  handleDeleteConfirm();
                }}
              >
                删除
              </Button>
              <Button
                ghost
                type="primary"
                loading={handleCreateExperiment?.loading}
                onClick={() => {
                  handleSubmit();
                }}
              >
                完成
              </Button>
            </Space>
          )}
        </div>
      </Form>
    );
  };

  /**
   * 更新基础信息
   * @param values
   */
  const handleConfirm = (values: any) => {
    form.setFieldsValue(values);
    setBaseInfo({ ...baseInfo, ...values });
  };

  useEffect(() => {
    const { experimentId, spaceId } = history?.location?.query || {};
    // 地址栏中存在空间id，需要将空间列表选项更新，并保存当前id
    if (spaceId) {
      getUserSpaceAuth?.run({
        id: spaceId as string,
      });
    }
    if (experimentId) {
      getExperimentDetail?.run({ uuid: experimentId as string });
    } else {
      setArrangeList(arrangeDataOriginTranstion([]));
      form.setFieldValue('name', '实验名称');
    }
  }, [history.location.query]);

  return (
    <Container>
      <Spin spinning={getExperimentDetail.loading}>
        <PageContainer
          header={{
            title: renderTitle(),
            onBack: () => {
              history.push('/space/experiment');
            },
            extra: headerExtra(),
          }}
        >
          <ArrangeContent
            arrangeList={arrangeList}
            setArrangeList={setArrangeList}
            disabled={spacePermission !== 1}
          />
          {infoDrawerOpen && (
            <InfoDrawer
              open={infoDrawerOpen}
              setOpen={setInfoDrawerOpen}
              spacePermission={spacePermission}
              handleConfirm={handleConfirm}
              baseInfo={baseInfo}
            />
          )}
        </PageContainer>
      </Spin>
    </Container>
  );
};

export default AddExperiment;
