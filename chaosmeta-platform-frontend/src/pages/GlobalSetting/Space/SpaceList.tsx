import { formatTime } from '@/utils/format';
import { EllipsisOutlined } from '@ant-design/icons';
import { history, useModel } from '@umijs/max';
import { Col, Dropdown, Row, Tooltip } from 'antd';
import React from 'react';
import { SpaceCard } from './style';

interface IProps {
  pageData: any;
  handleDelete: (id: number) => void;
}

const SpaceList: React.FC<IProps> = (props) => {
  const { pageData, handleDelete } = props;
  const { userInfo, setCurSpace } = useModel('global');

  const items = (spaceId: number) => [
    {
      label: <div>空间设置</div>,
      key: 'spaceSetting',
    },
    {
      label: (
        <div
          onClick={() => {
            handleDelete(spaceId);
          }}
        >
          删除
        </div>
      ),
      key: 'delete',
    },
  ];

  /**
   * 渲染读写权限成员
   */
  // const renderAdmin = (userData: any) => {
  //   const users = userData?.user?.filter((item) => item?.permission === 1);
  //   const text = users?.map((item) => item);
  // };

  // 当前用户相对应某个空间是否有权限
  const getPermission = (userData: any) => {
    // permission 0 = 只读， 1 = 读写， 2 = 未加入
    let permission = 2;
    // 全局管理员角色，默认拥有所有空间权限
    if (userInfo?.role === 'admin') {
      permission = 1;
    } else {
      permission = userData?.users?.filter(
        (item) => item?.user_id === userInfo?.id,
      )[0]?.permission;
    }
    return permission || permission === 0 ? permission : 2;
  };

  // const handleUpdateSpaceId = (id: any) => {
  //   if (id) {
  //     const name = form.getFieldValue('name');
  //     history.push({
  //       pathname: history.location.pathname,
  //       query: {
  //         ...history.location.query,
  //         spaceId: id,
  //       },
  //     });
  //     setCurSpace([id]);
  //     sessionStorage.setItem('spaceId', id);
  //     sessionStorage.setItem('spaceName', name);
  //   }
  // };

  /**
   * 点击卡片跳转到对应空间
   * @param record
   */
  const handleClickSpace = (record: any) => {
    const { name, id } = record;
    if (id) {
      history.push({
        pathname: '/space/overview',
        query: {
          spaceId: id,
        },
      });
      setCurSpace([id]);
      sessionStorage.setItem('spaceId', id);
      sessionStorage.setItem('spaceName', name);
    }
  };

  return (
    <Row gutter={[16, 16]}>
      {pageData?.namespaces?.map((item: any, index: number) => {
        const permission = getPermission(item?.userData);
        return (
          <Col span={6} key={index}>
            <SpaceCard $permission={permission}>
              <Tooltip
                // 未加入的提示
                title={permission === 2 && '你没有该空间的权限，请联系读写成员'}
              >
                <div className="header">
                  <div>
                    {permission !== 1 ? (
                      <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*5OzsS5d_il8AAAAAAAAAAAAADmKmAQ/original" />
                    ) : (
                      <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*3rVvS7yMa38AAAAAAAAAAAAADmKmAQ/original" />
                    )}

                    <div
                      onClick={() => {
                        // 未加入不允许进入
                        if (permission !== 2) {
                          handleClickSpace(item?.namespaceInfo);
                        }
                      }}
                      style={{ cursor: 'pointer' }}
                    >
                      <div className="title">{item.namespaceInfo?.name}</div>
                      <span className="time">
                        {formatTime(item?.namespaceInfo?.create_time)}
                      </span>
                    </div>
                  </div>
                  <Dropdown
                    disabled={permission !== 1}
                    placement="bottomRight"
                    menu={{
                      items: items(item.namespaceInfo?.namespaceInfo?.id),
                    }}
                  >
                    <Tooltip
                      title={permission === 0 ? '只读用户占无法使用此功能' : ''}
                    >
                      <EllipsisOutlined />
                    </Tooltip>
                  </Dropdown>
                </div>
              </Tooltip>

              <div className="desc">
                <div>描述：</div>
                <Tooltip title="haoduohaoduo">
                  <span>{item?.namespaceInfo?.description}</span>
                </Tooltip>
              </div>
              <div className="footer">
                <div>
                  <Tooltip title="读写成员">
                    <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*_TiCQ6O9B_oAAAAAAAAAAAAADmKmAQ/original" />
                    {/* <span>{renderAdmin(item?.userData)}</span> */}
                  </Tooltip>
                </div>

                <div>
                  <Tooltip title="实验数量">
                    <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*lps4TYQ9p4MAAAAAAAAAAAAADmKmAQ/original" />
                    <span>{item?.experimentrData?.toTal}</span>
                  </Tooltip>
                </div>

                <div>
                  <Tooltip title="空间成员">
                    <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*GLyEQrfTN68AAAAAAAAAAAAADmKmAQ/original" />
                    <span>{item?.userData?.toTal}</span>
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
