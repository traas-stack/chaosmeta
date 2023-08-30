import { arrangeNodeTypeColors } from '@/constants';
import { formatDuration } from '@/utils/format';
import { InfoCircleOutlined } from '@ant-design/icons';
import {
  CollisionDetection,
  DndContext,
  DragOverlay,
  DropAnimation,
  MouseSensor,
  TouchSensor,
  UniqueIdentifier,
  closestCenter,
  defaultDropAnimationSideEffects,
  getFirstCollision,
  pointerWithin,
  rectIntersection,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import {
  SortableContext,
  arrayMove,
  horizontalListSortingStrategy,
} from '@dnd-kit/sortable';
import { Form } from 'antd';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { v1 } from 'uuid';
import Arrange from './components/Arrange';
import DroppableItem from './components/DroppableItem';
import NodeConfig from './components/NodeConfig';
import NodeLibrary from './components/NodeLibrary';
import { DroppableCol, DroppableRow, NodeItem } from './style';

interface IProps {
  arrangeList: any[];
  setArrangeList: any;
  disabled?: boolean;
}

const ArrangeContent: React.FC<IProps> = (props) => {
  const { arrangeList, setArrangeList, disabled = false } = props;
  const [form] = Form.useForm();
  // 当前正在拖动元素的id
  const [activeId, setActiveId] = useState<UniqueIdentifier | null>(null);
  // 时间轴个数
  const [timeCount, setTimeCount] = useState<number>(35);
  // 当前拖动元素的数据，左侧节点/右侧row/右侧row-item
  const [curDragData, setCurDragData] = useState<any>(null);
  const recentlyMovedToNewContainer = useRef(false);
  const lastOverId = useRef<UniqueIdentifier | null>(null);
  // 当前选中的行内子元素
  const [activeCol, setActiveCol] = useState<any>({ state: false });
  // 判断数组中是否包含该项
  const isExist = (dragId: any, arr: any[]) => {
    return arr?.some((item) => item?.row === dragId);
  };

  /**
   * 碰撞算法 -- 参考dnd-kit多容器
   * Custom collision detection strategy optimized for multiple containers
   * - First, find any droppable containers intersecting with the pointer.
   * - If there are none, find intersecting containers with the active draggable.
   * - If there are no intersecting containers, return the last matched intersection
   *
   */
  const collisionDetectionStrategy: CollisionDetection = useCallback(
    (args) => {
      if (activeId && isExist(activeId, arrangeList)) {
        return closestCenter({
          ...args,
          droppableContainers: args.droppableContainers.filter(
            (container: any) => {
              return isExist(container?.id, arrangeList);
            },
          ),
        });
      }
      // Start by finding any intersecting droppable
      const pointerIntersections = pointerWithin(args);
      const intersections =
        pointerIntersections.length > 0
          ? // If there are droppables intersecting with the pointer, return those
            pointerIntersections
          : rectIntersection(args);
      let overId: any = getFirstCollision(intersections, 'id');
      if (overId !== null) {
        if (overId === 'void') {
          // If the intersecting droppable is the trash, return early
          // Remove this if you're not using trashable functionality in your app
          return intersections;
        }
        if (isExist(activeId, arrangeList)) {
          // const containerItems = arrangeList[overId];
          const containerItems = arrangeList?.filter(
            (item) => item?.row === overId,
          )[0]?.children;
          // If a container is matched and it contains items (columns 'A', 'B', 'C')
          if (containerItems.length > 0) {
            // Return the closest droppable within that container
            overId = closestCenter({
              ...args,
              droppableContainers: args.droppableContainers.filter(
                (container) =>
                  container.id !== overId &&
                  isExist(container.id, containerItems),
              ),
            })[0]?.id;
          }
        }

        lastOverId.current = overId;

        return [{ id: overId }];
      }
      // When a draggable item moves to a new container, the layout may shift
      // and the `overId` may become `null`. We manually set the cached `lastOverId`
      // to the id of the draggable item that was moved to the new container, otherwise
      // the previous `overId` will be returned which can cause items to incorrectly shift positions
      if (recentlyMovedToNewContainer.current) {
        lastOverId.current = activeId;
      }
      // If no droppable is matched, return the last match
      return lastOverId.current ? [{ id: lastOverId.current }] : [];
    },
    [activeId, arrangeList],
  );

  // 传感器，可配置拖动激活的条件
  const sensors = useSensors(
    useSensor(MouseSensor, {
      // 激活前要求鼠标移动5px
      activationConstraint: {
        distance: 2,
      },
    }),
    useSensor(TouchSensor, {
      // 激活延迟100ms，移动2px
      activationConstraint: {
        delay: 100,
        tolerance: 2,
      },
    }),
  );

  // 放置动画
  const dropAnimation: DropAnimation = {
    sideEffects: defaultDropAnimationSideEffects({
      styles: {
        active: {
          opacity: '0.5',
        },
      },
    }),
  };

  /**
   * 拖动中节点的渲染
   */
  const MoveingRender = () => {
    const { dragtype, name, index, duration, nameCn, exec_type } =
      curDragData || {};
    const second: number = formatDuration(duration);
    let renderItem = null;
    // 左侧节点拖动
    const leftNodeRender = () => {
      return (
        <NodeItem>
          <div className="temp-item">
            <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*rOAzRrDGQoAAAAAAAAAAAAAADmKmAQ/original" />
            {nameCn}
          </div>
        </NodeItem>
      );
    };
    // 右侧行拖动
    const rowRender = () => {
      return (
        <DroppableRow $isMoving={true}>
          <div className="row">
            {curDragData?.children && (
              <SortableContext
                items={curDragData?.children}
                strategy={horizontalListSortingStrategy}
              >
                {curDragData?.children?.map((el: any, j: number) => {
                  return (
                    <DroppableItem
                      key={j}
                      index={j}
                      item={el}
                      parentId={curDragData?.row}
                    />
                  );
                })}
              </SortableContext>
            )}
          </div>
          <div className="handle">{index}</div>
        </DroppableRow>
      );
    };
    // 右侧行内子元素拖动
    const rowItemRender = () => {
      return (
        <DroppableCol
          $bg={arrangeNodeTypeColors[exec_type]}
          style={{ width: `${second * 3}px` }}
        >
          <div className="item">
            <>
              <div className="title ellipsis">
                {!curDragData?.nodeInfoState && (
                  <InfoCircleOutlined
                    style={{ color: '#FF4D4F', marginRight: '4px' }}
                  />
                )}
                <span>{name}</span>
              </div>
              <div>{second}s</div>
            </>
          </div>
        </DroppableCol>
      );
    };
    if (dragtype === 'node') {
      renderItem = leftNodeRender();
    }
    if (dragtype === 'row') {
      renderItem = rowRender();
    }
    if (dragtype === 'item') {
      renderItem = rowItemRender();
    }
    return (
      <>
        {createPortal(
          <DragOverlay adjustScale={false} dropAnimation={dropAnimation}>
            {activeId ? renderItem : null}
          </DragOverlay>,
          document.body,
        )}
      </>
    );
  };

  /**
   * 拖动开始
   */
  const handleDragStart = (params: any) => {
    const { active } = params;
    setCurDragData(active?.data?.current);
    setActiveId(active.id);
  };

  /**
   * 遍历编排数组，将其中行中最长秒数对比当前默认时间轴，若时间轴宽度不够就再加长
   * @param values
   */
  const handleAddTimeAxis = (values: any) => {
    const maxSecondList: any = [];
    values?.forEach((item: { children: any[] }) => {
      let secondSum = 0;
      item?.children?.forEach((el) => {
        // 将持续时长转化为数字形式进行计算
        const duration = formatDuration(el?.duration);
        secondSum += duration;
      });
      maxSecondList.push(secondSum);
    });
    const maxSecond = maxSecondList.sort((a: number, b: number) => b - a)[0];
    // 提前一段加长
    const curSecond = timeCount * 30 - 60;
    if (maxSecond > curSecond) {
      setTimeCount(() => {
        // 多加入90s
        const newCount = Math.round(maxSecond / 30) + 3;
        return newCount;
      });
    }
  };

  /**
   * 判断是否添加行
   */
  const handleAddRow = () => {
    // 第一行有子元素/最后一行有子元素，添加行
    setArrangeList((params: any) => {
      const values = JSON.parse(JSON.stringify(params));
      let maxRow = 0;
      // 找到当前数据中的最大行号
      values?.forEach((el: { row: number }) => {
        if (el?.row > maxRow) {
          maxRow = el?.row;
        }
      });
      if (values[0]?.children?.length > 0) {
        values.unshift({
          id: maxRow + 1,
          row: maxRow + 1,
          children: [],
        });
      }
      if (values[values?.length - 1]?.children?.length > 0) {
        values.push({
          id: maxRow + 1,
          row: maxRow + 1,
          children: [],
        });
      }
      // 遍历编排数组，将其中行中最长秒数对比当前默认时间轴，若时间轴宽度不够就再加长
      handleAddTimeAxis(values);
      return values;
    });
  };

  /**
   * 拖动至某个元素
   */
  const handleDragOver = (params: any) => {
    const { active, over } = params;
    if (!active.id || !over?.id || active?.id === over?.id) {
      return;
    }

    // 拖动的元素
    const {
      index: activeIndex,
      dragtype: activeType,
      parentId: activeParentId,
    } = active?.data?.current || {};
    // 拖放至的元素
    const {
      index: overIndex,
      dragtype: overType,
      parentId: overParentId,
    } = over?.data?.current || {};
    // 行内子元素的跨容器拖动，两种情况，
    // 1. 直接拖动到其他容器的子元素上，可利用子元素及父容器的id
    // 2. 拖动到其他容器但不在子元素上，在该容器末尾加入元素
    if (activeType === 'item' && (activeParentId || activeParentId === 0)) {
      setArrangeList((values: any[]) => {
        // 1. 直接拖动到其他容器的子元素上，可利用子元素及父容器的id  && activeParentId !== overParentId
        if (overType === 'item') {
          const activeParentIndex = values?.findIndex(
            (item) => item?.row === activeParentId,
          );
          const overParentIndex = values?.findIndex(
            (item) => item?.row === overParentId,
          );
          const activeItem = values[activeParentIndex]?.children[activeIndex];
          values[activeParentIndex]?.children?.splice(activeIndex, 1);
          values[overParentIndex]?.children?.splice(overIndex, 0, activeItem);
          recentlyMovedToNewContainer.current = true;
        }
        // 2. 拖动到其他容器但不在子元素上，在该容器末尾加入元素
        if (overType === 'row') {
          const activeParentIndex = values?.findIndex(
            (item) => item?.row === activeParentId,
          );
          const activeItem = values[activeParentIndex]?.children[activeIndex];

          values[activeParentIndex]?.children?.splice(activeIndex, 1);
          values[overIndex]?.children?.push(activeItem);
          recentlyMovedToNewContainer.current = true;
        }
        return values;
      });
      return;
    }
    // 从节点库拖入
    if (activeType === 'node' && over?.id) {
      setArrangeList((values: any[]) => {
        // 拖动到容器的子元素上
        if (overType === 'item') {
          const overParentIndex = values?.findIndex(
            (item) => item?.row === overParentId,
          );
          const activeItem = active?.data?.current;
          values[overParentIndex]?.children?.splice(overIndex, 0, {
            ...activeItem,
            // 暂时将uuid设置为当前节点的id用于拖拽，拖拽结束后这里的uuid需要重新生成，避免拖拽绑定重复id
            uuid: active.id,
            duration: '60s',
            dragtype: 'item',
            exec_id: activeItem?.id,
            name: activeItem?.nameCn,
            // 当前节点信息是否进行配置完成, wait类型只需要配置时长，所以这里默认为true
            nodeInfoState: activeItem?.exec_type === 'wait',
          });
          recentlyMovedToNewContainer.current = true;
        }
        // 2. 拖动到其他容器但不在子元素上，在该容器末尾加入元素
        if (overType === 'row') {
          const activeItem = active?.data?.current;
          values[overIndex]?.children?.push({
            ...activeItem,
            uuid: active.id,
            // 拖入时长默认15m
            duration: '60s',
            dragtype: 'item',
            // 将节点库id保存
            exec_id: activeItem?.id,
            name: activeItem?.nameCn,
            // 当前节点信息是否进行配置完成
            nodeInfoState: activeItem?.exec_type === 'wait',
          });
          recentlyMovedToNewContainer.current = true;
        }
        return values;
      });
    }
  };

  /**
   * 拖动结束后
   * 不跨容器时拖动结束后再进行数据的修改，拖动效果可用dnd-kit提供的transform实现，跨容器则需要再handleDragOver中实现
   */
  const handleDragEnd = (params: any) => {
    const { active, over } = params;
    if (!active.id || !over?.id) {
      setActiveId(null);
      return;
    }
    // 拖动的元素
    const {
      index: activeIndex,
      dragtype: activeType,
      parentId: activeParentId,
      isNode,
    } = active?.data?.current || {};
    // 拖放至的元素
    const {
      index: overIndex,
      dragtype: overType,
      parentId: overParentId,
    } = over?.data?.current || {};
    // 行之间的拖动
    if (activeType === 'row' && overType === 'row' && !isNode) {
      if (activeIndex !== overIndex) {
        setArrangeList((values: any[]) => {
          return arrayMove(values, activeIndex, overIndex);
        });
        handleAddRow();
      }
      setActiveId(null);
      return;
    }
    // 行内元素的拖动，不跨容器
    if (
      activeType === 'item' &&
      overType === 'item' &&
      activeParentId === overParentId &&
      !isNode
    ) {
      setArrangeList((values: any) => {
        // 拖动元素后修改数据
        const parentTemp = values?.filter(
          (item: { row: any }) => item.row === activeParentId,
        )[0];
        const parentIndex = values.findIndex(
          (item: { row: any }) => item.row === activeParentId,
        );
        const newParentList = arrayMove(
          parentTemp?.children,
          activeIndex,
          overIndex,
        );
        if (parentIndex !== -1) {
          values[parentIndex] = { ...parentTemp, children: newParentList };
        }
        return values;
      });
      setActiveId(null);
      handleAddRow();
      return;
    }
    // 行内元素的拖动，跨容器
    if (activeType === 'item' && overType === 'row' && !isNode) {
      handleAddRow();
      return;
    }
    // 从节点处拖入的元素需手动修改id，拖拽绑定了ID，避免id重复
    if (isNode) {
      setArrangeList((values: any[]) => {
        // 拖动至行内，为最后一个元素
        if (overType === 'row') {
          const newOverIndex = values[overIndex]?.children?.length - 1;
          values[overIndex].children[newOverIndex].uuid = v1().replaceAll(
            '-',
            '',
          );
          values[overIndex].children[newOverIndex].isNode = false;
        }
        // 拖动至行内的子元素上
        if (overType === 'item') {
          const overParentIndex = values?.findIndex(
            (item) => item.row === overParentId,
          );
          values[overParentIndex].children[overIndex].uuid = v1().replaceAll(
            '-',
            '',
          );
          values[overParentIndex].children[overIndex].isNode = false;
        }
        return values;
      });
      handleAddRow();
      setActiveId(null);
      return;
    }
  };

  /**
   * 取消拖动
   */
  const onDragCancel = () => {
    setActiveId(null);
  };

  useEffect(() => {
    requestAnimationFrame(() => {
      recentlyMovedToNewContainer.current = false;
    });
    handleAddTimeAxis(arrangeList);
  }, [arrangeList]);

  return (
    <>
      <DndContext
        collisionDetection={collisionDetectionStrategy}
        sensors={sensors}
        onDragStart={handleDragStart}
        onDragOver={(params: any) => {
          try {
            handleDragOver(params);
          } catch (error) {
            console.log(error, 'error');
          }
        }}
        onDragEnd={(params: any) => {
          try {
            handleDragEnd(params);
          } catch (error) {
            console.log(error, 'error');
          }
        }}
        onDragCancel={onDragCancel}
      >
        <div className="content">
          {/* 节点库 */}
          <NodeLibrary disabled={disabled} />
          {/* 编排区域 */}
          <Arrange
            arrangeList={arrangeList}
            setArrangeList={setArrangeList}
            timeCount={timeCount}
            setTimeCount={setTimeCount}
            activeCol={activeCol}
            setActiveCol={setActiveCol}
            disabled={disabled}
          />
          {/* 右侧编辑信息区域 */}
          {activeCol?.uuid && (
            <NodeConfig
              form={form}
              activeCol={activeCol}
              setActiveCol={setActiveCol}
              arrangeList={arrangeList}
              setArrangeList={setArrangeList}
              disabled={disabled}
            />
          )}
        </div>
        <MoveingRender />
      </DndContext>
    </>
  );
};

export default ArrangeContent;
