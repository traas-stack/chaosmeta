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
import Arrange from './components/Arrange';
import DroppableItem from './components/DroppableItem';
import NodeConfig from './components/NodeConfig';
import NodeLibrary from './components/NodeLibrary';
import { DroppableCol, DroppableRow, NodeItem } from './style';

interface IProps {
  arrangeList: any[];
  leftNodeList: any[];
  setArrangeList: any;
}

const ArrangeContent: React.FC<IProps> = (props) => {
  const { arrangeList, leftNodeList, setArrangeList } = props;
  const [form] = Form.useForm();
  // 当前正在拖动元素的id
  const [activeId, setActiveId] = useState<UniqueIdentifier | null>(null);
  // 时间轴个数
  const [timeCount, setTimeCount] = useState<number>(16);
  // 当前拖动元素的数据，左侧节点/右侧row/右侧row-item
  const [curDragData, setCurDragData] = useState<any>(null);
  const recentlyMovedToNewContainer = useRef(false);
  const lastOverId = useRef<UniqueIdentifier | null>(null);
  // 当前选中的行内子元素
  const [activeCol, setActiveCol] = useState<any>({ state: false });
  // 判断数组中是否包含该项
  const isExist = (id: any, arr: any[]) => {
    return arr?.some((item) => item?.id === id);
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
          droppableContainers: args.droppableContainers.filter((container) =>
            isExist(container?.id, arrangeList),
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
            (item) => item?.id === overId,
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
    const { dragtype, name, index, id, second } = curDragData || {};
    let renderItem = null;
    // 左侧节点拖动
    const leftNodeRender = () => {
      return (
        <NodeItem>
          <div className="temp-item">
            {/* <div> */}
            <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*rOAzRrDGQoAAAAAAAAAAAAAADmKmAQ/original" />
            {name}
            {/* </div> */}
          </div>
        </NodeItem>
      );
    };
    // 右侧行拖动
    const rowRender = () => {
      return (
        <DroppableRow $isMoving={true}>
          <div className="row">
            {/* {itemData?.id} */}
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
                      parentId={curDragData?.id}
                      // curProportion={curProportion}
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
          // $isMoving={isMoving}
          $bg={'#C6F8E0'}
          style={{ width: `${second * 3}px` }}
        >
          <div className="item">{id}</div>
        </DroppableCol>
        // <DroppableItem
        //   index={index}
        //   item={index}
        //   parentId={curDragData?.id}
        // />
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
   * 判断是否添加行
   */
  const handleAddRow = () => {
    // 第一行有子元素/最后一行有子元素，添加行
    setArrangeList((params: any) => {
      const values = JSON.parse(JSON.stringify(params));
      if (values[0]?.children?.length > 0) {
        values.unshift({ id: values?.length, children: [] });
      }
      if (values[values?.length - 1]?.children?.length > 0) {
        values.push({ id: values?.length, children: [] });
      }
      // 遍历编排数组，将其中行中最长秒数对比当前默认时间轴，若时间轴宽度不够就再加长
      const maxSecondList: any = [];
      values?.forEach((item: { children: any[] }) => {
        let secondSum = 0;
        item?.children?.forEach((el) => {
          secondSum += el?.second;
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
      return values;
    });
  };
  // 150

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
    if (activeType === 'item' && activeParentId) {
      setArrangeList((values: any[]) => {
        // 1. 直接拖动到其他容器的子元素上，可利用子元素及父容器的id  && activeParentId !== overParentId
        if (overType === 'item') {
          const activeParentIndex = values?.findIndex(
            (item) => item?.id === activeParentId,
          );
          const overParentIndex = values?.findIndex(
            (item) => item?.id === overParentId,
          );
          const activeItem = values[activeParentIndex]?.children[activeIndex];
          values[activeParentIndex]?.children?.splice(activeIndex, 1);
          values[overParentIndex]?.children?.splice(overIndex, 0, activeItem);
          recentlyMovedToNewContainer.current = true;
        }
        // 2. 拖动到其他容器但不在子元素上，在该容器末尾加入元素
        if (overType === 'row') {
          const activeParentIndex = values?.findIndex(
            (item) => item?.id === activeParentId,
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
            (item) => item?.id === overParentId,
          );
          const activeItem = active?.data?.current;
          values[overParentIndex]?.children?.splice(overIndex, 0, {
            ...activeItem,
            dragtype: 'item',
          });
          recentlyMovedToNewContainer.current = true;
        }
        // 2. 拖动到其他容器但不在子元素上，在该容器末尾加入元素
        if (overType === 'row') {
          const activeItem = active?.data?.current;
          values[overIndex]?.children?.push({
            ...activeItem,
            dragtype: 'item',
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
        // const curList = JSON.parse(JSON.stringify(values));
        // 拖动元素后修改数据
        const parentTemp = values?.filter(
          (item: { id: any }) => item.id === activeParentId,
        )[0];
        const parentIndex = values.findIndex(
          (item: { id: any }) => item.id === activeParentId,
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
          values[overIndex].children[newOverIndex].id = `${
            newOverIndex + 1
          }-move`;
          values[overIndex].children[newOverIndex].isNode = false;
        }
        // 拖动至行内的子元素上
        if (overType === 'item') {
          const overParentIndex = values?.findIndex(
            (item) => item.id === overParentId,
          );
          const newOverIndex = values[overParentIndex]?.children?.length;
          values[overParentIndex].children[
            overIndex
          ].id = `${newOverIndex}-move`;
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
  }, [arrangeList]);

  return (
    <>
      <DndContext
        collisionDetection={collisionDetectionStrategy}
        sensors={sensors}
        onDragStart={handleDragStart}
        onDragOver={handleDragOver}
        onDragEnd={handleDragEnd}
        onDragCancel={onDragCancel}
      >
        <div className="content">
          {/* 节点库 */}
          <NodeLibrary leftNodeList={leftNodeList} />
          {/* 编排区域 */}
          <Arrange
            arrangeList={arrangeList}
            setArrangeList={setArrangeList}
            timeCount={timeCount}
            setTimeCount={setTimeCount}
            activeCol={activeCol}
            setActiveCol={setActiveCol}
          />
          {/* 右侧编辑信息区域 */}
          {activeCol?.state && (
            <NodeConfig
              form={form}
              activeCol={activeCol}
              setActiveCol={setActiveCol}
              arrangeList={arrangeList}
              setArrangeList={setArrangeList}
            />
          )}
        </div>
        <MoveingRender />
      </DndContext>
    </>
  );
};

export default ArrangeContent;
