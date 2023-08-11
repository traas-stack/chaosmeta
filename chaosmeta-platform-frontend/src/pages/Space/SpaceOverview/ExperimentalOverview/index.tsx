import EmptyCustom from '@/components/EmptyCustom';
import { RightOutlined } from '@ant-design/icons';
import { history } from '@umijs/max';
import { Button, Card, Space, Table, Tabs } from 'antd';

export default () => {
  const operations = (
    <Space>
      <span>查看全部实验</span>
      <RightOutlined />
    </Space>
  );

  const columns: any[] = [
    {
      dataIndex: 'name',
      key: 'name',
      width: 120,
      render: () => {
        return (
          <>
            <div className="ellipsis row-text-gap">
              <a>
                <span>
                  我名字很长-我名字很长-我名字很长-我名字很长-我名字很长
                </span>
              </a>
            </div>
            <div className="shallow">最近编辑时间： 2021-07-11 11:20:33</div>
          </>
        );
      },
    },
    {
      dataIndex: 'type',
      key: 'type',
      width: 110,
      render: () => {
        return (
          <div>
            <div className="shallow row-text-gap">触发方式</div>
            <div>
              <span>周期性</span>
              <span className="shallow">（2023-07-11 14:00:00）</span>
            </div>
          </div>
        );
      },
    },
    {
      dataIndex: 'count',
      key: 'count',
      width: 60,
      render: () => {
        return (
          <div>
            <div className="shallow row-text-gap">实验次数</div>
            <div>22</div>
          </div>
        );
      },
    },
    {
      dataIndex: 'action',
      key: 'action',
      width: 50,
      fixed: 'right',
      render: () => {
        return (
          <Space>
            <a>运行</a>
            <a>编辑</a>
          </Space>
        );
      },
    },
  ];

  // 最近编辑的实验
  const EditExperimental = () => {
    return (
      <Card>
        <Table
          locale={{
            emptyText: (
              <EmptyCustom
                desc="当前页面暂无最近编辑的实验"
                title="您可以前往实验列表编辑实验"
                // btnText="前往实验列表"
                // btnOperate={() => {
                //   history?.push('/space/experiment');
                // }}
                btns={
                  <Button
                    type="primary"
                    onClick={() => {
                      history?.push('/space/experiment');
                    }}
                  >
                    前往实验列表
                  </Button>
                }
              />
            ),
          }}
          showHeader={false}
          columns={columns}
          rowKey={'id'}
          // dataSource={[
          //   { name: '999', id: '1' },
          //   { name: '999', id: '5' },
          // ]}
          dataSource={[]}
          pagination={false}
          scroll={{ x: 760 }}
        />
      </Card>
    );
  };

  // 即将运行的实验
  const SoonRunExperimental = () => {
    return (
      <Card>
        <Table
          showHeader={false}
          columns={columns}
          rowKey={'id'}
          dataSource={[]}
          pagination={false}
          scroll={{ x: 760 }}
        />
      </Card>
    );
  };

  // 最近运行的实验结果
  const RecentlyRunExperimentalResult = () => {
    return (
      <Card>
        <Table
          showHeader={false}
          columns={columns}
          rowKey={'id'}
          dataSource={[
            { name: '999', id: '1' },
            { name: '999', id: '2' },
            { name: '999', id: '3' },
            { name: '999', id: '4' },
            { name: '999', id: '5' },
          ]}
          pagination={false}
          scroll={{ x: 760 }}
        />
      </Card>
    );
  };

  const items = [
    {
      label: '最近编辑的实验',
      key: '1',
      children: <EditExperimental />,
    },
    {
      label: '即将运行的实验',
      key: '2',
      children: <SoonRunExperimental />,
    },
    {
      label: '最近运行的实验结果',
      key: '3',
      children: <RecentlyRunExperimentalResult />,
    },
  ];

  return <Tabs tabBarExtraContent={operations} items={items} />;
};
