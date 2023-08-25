import ShowText from '@/components/ShowText';
import { scaleStepMap } from '@/constants';
import { useSortable } from '@dnd-kit/sortable';
import React from 'react';
import { DroppableCol } from '../style';

interface IProps {
  item: any;
  index: number;
  parentId: string;
  activeCol?: any;
  setActiveCol?: any;
  disabled?: boolean;
  curProportion?: number;
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
  } = props;
  const { setNodeRef, transform, listeners, isDragging } = useSortable({
    id: item?.id,
    disabled,
    data: {
      parentId,
      index,
      dragtype: 'item',
      ...item,
    },
  });

  return (
    <DroppableCol
      ref={setNodeRef}
      $isDragging={isDragging}
      $bg={'#C6F8E0'}
      // 减去外边距的2px，避免子元素多时宽度偏差过大
      style={{
        width: `${
          item?.second * (scaleStepMap[curProportion]?.widthSecond || 3) - 2
        }px`,
      }}
      $transform={transform}
      $activeState={activeCol?.id === item?.id && activeCol?.state}
      onClick={() => {
        if (disabled) {
          return;
        }
        if (activeCol?.id === item?.id) {
          setActiveCol({ state: false });
        } else {
          setActiveCol({ ...item, parentId, index, state: true });
        }
      }}
    >
      <div className="item ellipsis" {...listeners}>
        {item.second > 10 ? (
          <>
            <ShowText
              value={'测试超出测试超出测试超出测试超出测试超出测试超出'}
              ellipsis
            />
            {/* <div>{item.id}</div> */}
            <div>{item?.second}s</div>
          </>
        ) : (
          <>
            <div>...</div>
            <div>...</div>
          </>
        )}
      </div>
      {activeCol?.id === item?.id && activeCol?.state && (
        <div className="scale"></div>
      )}
    </DroppableCol>
  );
};

export default DroppableItem;
