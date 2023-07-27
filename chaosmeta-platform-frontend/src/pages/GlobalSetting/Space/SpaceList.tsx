import {
  EllipsisOutlined,
  ExclamationCircleFilled,
  SearchOutlined,
} from '@ant-design/icons';
import {
  Alert,
  Col,
  Dropdown,
  Input,
  Modal,
  Pagination,
  Row,
  Space,
  Tooltip,
  message,
} from 'antd';
import React from 'react';
import { SpaceCard } from './style';

interface IProps {
  listType: 'my' | 'all';
  title: string;
}

const SpaceList: React.FC<IProps> = (props) => {
  const { title, listType } = props;

  /**
   * 删除空间
   */
  const handleDeleteAccount = () => {
    console.log('====');
    Modal.confirm({
      title: '确认要删除当前所选账号吗？',
      icon: <ExclamationCircleFilled />,
      content: '删除账号用户将无法登录平台，要再次使用只能重新注册！',
      onOk() {
        return new Promise((resolve, reject) => {
          setTimeout(Math.random() > 0.5 ? resolve : reject, 1000);
          message.success('您已成功删除所选成员');
        }).catch(() => console.log('Oops errors!'));
      },
      onCancel() {},
    });
  };

  const items = [
    {
      label: <div>空间设置</div>,
      key: 'spaceSetting',
    },
    {
      label: (
        <div
          onClick={() => {
            console.log('click');
            handleDeleteAccount();
          }}
        >
          删除1
        </div>
      ),
      key: 'delete',
    },
  ];

  return (
    <div className="space-list">
      <div className="area-operate">
        <div className="title">{title}</div>
        <Space>
          <Input
            placeholder="请输入空间名称、空间成员"
            suffix={<SearchOutlined onClick={() => {}} />}
          />
        </Space>
      </div>
      {listType === 'all' && (
        <Alert
          message="可联系空间内具有读写权限的成员添加为空间成员"
          type="info"
          showIcon
          closable
        />
      )}
      <div>
        <Row gutter={[16, 16]}>
          {[1, 2, 3, 4, 5].map((item, index) => {
            return (
              <Col span={6} key={index}>
                <SpaceCard>
                  <div className="header">
                    <div>
                      {listType === 'all' ? (
                        <Tooltip title="未加入">
                          <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*5OzsS5d_il8AAAAAAAAAAAAADmKmAQ/original" />
                        </Tooltip>
                      ) : (
                        <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*3rVvS7yMa38AAAAAAAAAAAAADmKmAQ/original" />
                      )}
                      <div>
                        <div className="title">工作空间</div>
                        <span className="time">2022-01-01 10:11:22</span>
                      </div>
                    </div>
                    <Dropdown
                      // open
                      disabled
                      placement="bottomRight"
                      menu={{
                        items,
                      }}
                    >
                      <Tooltip title="您不具有该空间的读写权限，暂无法使用此功能">
                        <EllipsisOutlined />
                      </Tooltip>
                    </Dropdown>
                  </div>
                  <div className="desc">
                    <div>描述：</div>
                    <span>
                      对这个空间进行简短的描述，Lorem ipsum dolor sit, amet
                      consectetur adipisicing elit. Iste consequuntur non,
                      eligendi exercitationem, nihil ex architecto ea, explicabo
                      tenetur a corporis ipsum cupiditate. Vero modi quae saepe
                      iure dolor recusandae.
                    </span>
                  </div>
                  <div className="footer">
                    <div>
                      <Tooltip title="具有读写权限的成员">
                        <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*_TiCQ6O9B_oAAAAAAAAAAAAADmKmAQ/original" />
                        <span>张三、李四</span>
                      </Tooltip>
                    </div>
                    <div>
                      <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*lps4TYQ9p4MAAAAAAAAAAAAADmKmAQ/original" />
                      <span>12</span>
                    </div>
                    <div>
                      <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*GLyEQrfTN68AAAAAAAAAAAAADmKmAQ/original" />
                      <span>13</span>
                    </div>
                  </div>
                </SpaceCard>
              </Col>
            );
          })}
        </Row>
      </div>
      <Pagination
        showQuickJumper
        defaultCurrent={2}
        total={500}
        onChange={(page, pageSize) => {
          console.log(page, pageSize, 'page, pageSize');
        }}
      />
    </div>
  );
};

export default SpaceList;
