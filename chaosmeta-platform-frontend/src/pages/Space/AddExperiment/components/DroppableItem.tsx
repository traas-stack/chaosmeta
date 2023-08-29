import { arrangeNodeTypeColors, scaleStepMap } from '@/constants';
import { formatDuration } from '@/utils/format';
import { InfoCircleOutlined } from '@ant-design/icons';
import { useSortable } from '@dnd-kit/sortable';
import { Descriptions, Popover } from 'antd';
import React, { useEffect, useRef } from 'react';
import { DroppableCol, NodeItemHover } from '../style';

interface IProps {
  item: any;
  index: number;
  parentId: string;
  activeCol?: any;
  setActiveCol?: any;
  disabled?: boolean;
  curProportion?: number;
  setArrangeList?: any;
}

/**
 * 编排区域行内块节点
 */
const DroppableItem: React.FC<IProps> = (props) => {
  const {
    item,
    index,
    parentId,
    setActiveCol,
    activeCol,
    disabled,
    curProportion = 100,
    setArrangeList,
  } = props;
  const { setNodeRef, transform, listeners, isDragging } = useSortable({
    id: item?.uuid,
    disabled,
    data: {
      parentId,
      index,
      dragtype: 'item',
      ...item,
    },
  });
  const curDuration = formatDuration(item?.duration);
  // 节点拖动宽度的手柄ref
  const scaleRef = useRef<any>(null);
  // nodeInfoState当前节点信息是否配置完成
  const { nodeInfoState } = item || {};

  /**
   * 悬浮到节点展示信息
   */
  const hoverRenderInfo = (item: any) => {
    const { name } = item || {};
    const renterTitle = () => {
      return (
        <div className="title">
          <div>
            {!nodeInfoState && <InfoCircleOutlined className="icon" />}
            <span className="name">{name}</span>
          </div>
          {!nodeInfoState && <div style={{ color: '#FF4D4F' }}>未完成</div>}
        </div>
      );
    };
    return (
      <NodeItemHover>
        <Descriptions title={renterTitle()} column={1}>
          <Descriptions.Item label="持续时长">
            {formatDuration(item?.duration)}s
          </Descriptions.Item>
        </Descriptions>
      </NodeItemHover>
    );
  };

  useEffect(() => {
    if (scaleRef.current) {
      // 拖动修改节点宽度
      scaleRef.current.onmousedown = (event: any) => {
        const duration = formatDuration(item?.duration);
        // 初始宽度
        const oldWidth =
          duration * (scaleStepMap[curProportion]?.widthSecond || 3) - 2;
        event.preventDefault();
        const startX = event.clientX;
        document.onmousemove = (moveEvent: any) => {
          moveEvent.preventDefault();
          // 拖动后计算新的节点宽度
          const newWidth = moveEvent.clientX - startX + oldWidth;
          // 将宽度计算为对应秒数
          const second = newWidth / scaleStepMap[curProportion]?.widthSecond;
          // 更新数据
          setArrangeList((result: any) => {
            const values = JSON.parse(JSON.stringify(result)); // 将 result 对象深拷贝一份
            const parentIndex = values?.findIndex(
              (item: { row: any }) => item?.row === parentId,
            );
            if (parentIndex !== -1 && index >= 0) {
              values[parentIndex].children[index]['duration'] = `${Math.round(
                second,
              )}s`; // 更新子节点对应属性的值
            }
            // 同时更新选中项数据
            setActiveCol((value: any) => {
              return { ...value, duration: `${Math.round(second)}s` };
            });
            return values; // 返回更新后的 values 数组
          });
        };
        // 释放鼠标时，清除移动事件监听器
        document.onmouseup = () => {
          document.onmousemove = document.onmouseup = null;
        };
      };
    }
  }, []);

  return (
    <DroppableCol
      ref={setNodeRef}
      $isDragging={isDragging}
      $bg={arrangeNodeTypeColors[item?.exec_type]}
      // 减去外边距的2px，避免子元素多时宽度偏差过大
      style={{
        width: `${
          curDuration * (scaleStepMap[curProportion]?.widthSecond || 3) - 2
        }px`,
        // 最小宽度为1s对应的px
        minWidth: `${scaleStepMap[curProportion]?.widthSecond}px`,
      }}
      $transform={transform}
      $nodeState={nodeInfoState}
      $activeState={activeCol?.uuid === item?.uuid}
      onClick={() => {
        if (activeCol?.uuid === item?.uuid) {
          setActiveCol({});
        } else {
          setActiveCol({ ...item, parentId, index });
        }
      }}
    >
      <Popover title={hoverRenderInfo(item)} placement="topLeft">
        <div className="item" {...listeners}>
          {curDuration * (scaleStepMap[curProportion]?.widthSecond || 3) >
          30 ? (
            <>
              <div className="title ellipsis">
                {!nodeInfoState && (
                  <InfoCircleOutlined
                    style={{ color: '#FF4D4F', marginRight: '4px' }}
                  />
                )}
                <span>{item.name}</span>
              </div>
              <div>{curDuration}s</div>
            </>
          ) : (
            <>
              <div>...</div>
              <div>...</div>
            </>
          )}
        </div>
      </Popover>
      {/* 拖动节点宽度的手柄 */}
      {activeCol?.uuid === item?.uuid && !isDragging && (
        <Popover overlayStyle={{ width: '56px' }} title={activeCol?.duration}>
          <div className="scale" ref={scaleRef}></div>
        </Popover>
      )}
    </DroppableCol>
  );
};

export default DroppableItem;
