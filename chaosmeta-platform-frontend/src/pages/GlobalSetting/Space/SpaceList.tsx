import { formatTime } from '@/utils/format';
import { EllipsisOutlined } from '@ant-design/icons';
import { Col, Dropdown, Row, Tooltip } from 'antd';
import React from 'react';
import { SpaceCard } from './style';

interface IProps {
  pageData: any;
  handleDelete: (id: number) => void;
}

const SpaceList: React.FC<IProps> = (props) => {
  const { pageData, handleDelete } = props;

  const items = (spaceId: number) => [
    {
      label: <div>空间设置</div>,
      key: 'spaceSetting',
    },
    {
      label: (
        <div
          onClick={() => {
            console.log(spaceId, 'spaceId');
            handleDelete(spaceId);
          }}
        >
          删除
        </div>
      ),
      key: 'delete',
    },
  ];

  return (
    <Row gutter={[16, 16]}>
      {pageData?.namespaces?.map((item: any, index: number) => {
        return (
          <Col span={6} key={index}>
            <SpaceCard style={{ background: 'rgba(0,0,0,0.05)' }}>
              <div className="header">
                <div>
                  {/* {listType === 'all' ? (
                <Tooltip title="未加入">
                  <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*5OzsS5d_il8AAAAAAAAAAAAADmKmAQ/original" />
                </Tooltip>
              ) : (
                <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*3rVvS7yMa38AAAAAAAAAAAAADmKmAQ/original" />
              )} */}
                  <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*3rVvS7yMa38AAAAAAAAAAAAADmKmAQ/original" />
                  <Tooltip title="你没有该空间的权限，请联系读写成员">
                    <div>
                      <div className="title">{item.namespaceInfo?.name}</div>
                      <span className="time">
                        {formatTime(item?.namespaceInfo?.create_time)}
                      </span>
                    </div>
                  </Tooltip>
                </div>
                <Dropdown
                  // open
                  // disabled
                  placement="bottomRight"
                  menu={{
                    items: items(item.namespaceInfo?.namespaceInfo?.id),
                  }}
                >
                  <Tooltip title="您不具有该空间的读写权限，暂无法使用此功能">
                    <EllipsisOutlined />
                  </Tooltip>
                </Dropdown>
              </div>
              <div className="desc">
                <div>描述：</div>
                <Tooltip title="haoduohaoduo">
                  <span>{item?.namespaceInfo?.description}</span>
                </Tooltip>
              </div>
              <div className="footer">
                <div>
                  <Tooltip title="具有读写权限的成员">
                    <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*_TiCQ6O9B_oAAAAAAAAAAAAADmKmAQ/original" />
                    <span>张三、李四</span>
                  </Tooltip>
                </div>

                <div>
                  <Tooltip title="实验数量">
                    <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*lps4TYQ9p4MAAAAAAAAAAAAADmKmAQ/original" />
                    <span>12</span>
                  </Tooltip>
                </div>

                <div>
                  <Tooltip title="空间成员">
                    <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*GLyEQrfTN68AAAAAAAAAAAAADmKmAQ/original" />
                    <span>13</span>
                  </Tooltip>
                </div>
              </div>
            </SpaceCard>
          </Col>
        );
      })}
    </Row>
  );
};

export default SpaceList;
