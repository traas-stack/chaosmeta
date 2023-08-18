import { CaretRightOutlined, SearchOutlined } from '@ant-design/icons';
import { useSortable } from '@dnd-kit/sortable';
import { Collapse, Input } from 'antd';
import React from 'react';
import { NodeItem, NodeLibraryContainer } from '../style';

interface IProps {
  leftNodeList: any[];
}
/**
 * 节点库
 */
const NodeLibrary: React.FC<IProps> = (props) => {
  const { leftNodeList } = props;

  /**
   * 左侧节点渲染
   */
  const LeftNodeItem = (props: any) => {
    const { itemData } = props;
    // 用于绑定拖拽，listeners需配置到拖动元素上，拖动就靠它
    const { setNodeRef, listeners, isDragging } = useSortable({
      id: itemData?.id,
      // 额外数据，用于悬浮态数据的渲染和判断
      data: {
        dragtype: 'node',
        isNode: true,
        ...itemData,
      },
    });
    return (
      <NodeItem ref={setNodeRef} $isDragging={isDragging}>
        <div
          {...listeners}
          className="temp-item"
          // {...attributes}
        >
          <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*rOAzRrDGQoAAAAAAAAAAAAAADmKmAQ/original" />
          {itemData?.name}
        </div>
      </NodeItem>
    );
  };

  return (
    <NodeLibraryContainer>
      <div className="wrap">
        <div className="title">节点库</div>
        <div className="node">
          <div>
            <Input
              placeholder="搜索节点名称"
              onPressEnter={() => {
                // handlePageSearch();
              }}
              bordered={false}
              prefix={
                <SearchOutlined
                  onClick={() => {
                    // handlePageSearch();
                  }}
                />
              }
            />
          </div>
          {/* 左侧节点折叠面板 */}
          <Collapse
            defaultActiveKey={['1']}
            className="collapse"
            expandIcon={({ isActive }) => (
              <CaretRightOutlined rotate={isActive ? 90 : 0} />
            )}
            ghost
          >
            {leftNodeList?.map((item) => {
              return (
                <Collapse.Panel
                  key={item.id}
                  header={item.name}
                  collapsible="icon"
                >
                  <Collapse
                    className="collapse-second"
                    ghost
                    expandIcon={({ isActive }) => (
                      <CaretRightOutlined rotate={isActive ? 90 : 0} />
                    )}
                  >
                    {item?.children?.map((el: any) => {
                      // 判断当前节点下是否还有子集，有就继续遍历向下找，没有则直接渲染
                      if (el?.children?.length > 0) {
                        return (
                          <Collapse.Panel key={el.id} header={el.name}>
                            {el?.children?.map((temp: { id: any }) => {
                              return (
                                <LeftNodeItem itemData={temp} key={temp?.id} />
                              );
                            })}
                          </Collapse.Panel>
                        );
                      }
                      // 节点
                      return <LeftNodeItem itemData={el} key={el?.id} />;
                    })}
                  </Collapse>
                </Collapse.Panel>
              );
            })}
          </Collapse>
        </div>
        {/* <div className="fold-icon">
    <LeftOutlined />
    <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*j-zJRJOhf7YAAAAAAAAAAAAADmKmAQ/original" />
  </div> */}
      </div>
    </NodeLibraryContainer>
  );
};

export default NodeLibrary;
