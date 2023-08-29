import { arrangeNodeTypeColors } from '@/constants';
import { DroppableCol, DroppableRow, ThumbnailContainer } from './style';
const Thumbnail = () => {
  const initList = [
    {
      id: 1,
      children: [
        {
          id: '1-1',
          second: 100,
          nodeType: 'fault',
          name: 'cpu燃烧',
          nodeState: 'abnormal',
        },
        {
          id: '1-2',
          second: 30,
          nodeType: 'measure',
          name: 'cpu燃烧333',
          nodeState: 'abnormal',
        },
        {
          id: '1-3',
          second: 30,
          nodeType: 'pressure',
          nodeState: 'abnormal',
          name: 'cpu燃烧111',
        },
        {
          id: '1-4',
          second: 30,
          nodeType: 'other',
          nodeState: 'abnormal',
          name: 'cpu燃烧22222',
        },
        {
          id: '1-5',
          second: 30,
          nodeState: 'abnormal',
        },
        {
          id: '1-6',
          second: 30,
          nodeState: 'abnormal',
        },
      ],
    },
    // {
    //   id: 2,
    //   children: [
    //     {
    //       id: '2-1',
    //       second: 20,
    //     },
    //     {
    //       id: '2-2',
    //       second: 20,
    //     },
    //     {
    //       id: '2-3',
    //       second: 40,
    //     },
    //   ],
    // },
    // {
    //   id: 3,
    //   children: [
    //     {
    //       id: '3-1',
    //       second: 40,
    //     },
    //     {
    //       id: '3-2',
    //       second: 40,
    //     },
    //     {
    //       id: '3-3',
    //       second: 40,
    //     },
    //   ],
    // },
    // {
    //   id: 4,
    //   children: [
    //     {
    //       id: '4-1',
    //       second: 10,
    //     },
    //     {
    //       id: '4-2',
    //       second: 1,
    //     },
    //     {
    //       id: '4-3',
    //       second: 90,
    //     },
    //   ],
    // },

    // {
    //   id: 4,
    //   children: [
    //     {
    //       id: '4-1',
    //       second: 10,
    //     },
    //     {
    //       id: '4-2',
    //       second: 1,
    //     },
    //     {
    //       id: '4-3',
    //       second: 90,
    //     },
    //   ],
    // },
  ];

  /**
   * 每行信息的展示
   * @param props
   * @returns
   */
  const ArrangeRow = (props: any) => {
    const { item } = props;
    return (
      <div>
        <DroppableRow>
          <div className="row">
            {/* 行内子元素 */}
            {item?.children?.map((el: any, j: number) => {
              return (
                <DroppableCol
                  key={j}
                  $bg={arrangeNodeTypeColors[el?.nodeType]}
                  // 减去外边距的2px，避免子元素多时宽度偏差过大
                  style={{
                    width: `${el?.second}px`,
                  }}
                >
                  <div className="item">
                    {el.second * 3 > 30 ? (
                      <>
                        <div className="ellipsis">
                          <span>{el.name}</span>
                        </div>
                        <div>{el?.second}s</div>
                      </>
                    ) : (
                      <>
                        <div>...</div>
                        <div>...</div>
                      </>
                    )}
                  </div>
                </DroppableCol>
              );
            })}
          </div>
        </DroppableRow>
      </div>
    );
  };

  return (
    <ThumbnailContainer>
      {/* 编排元素展示 */}
      {initList?.map((item, index) => {
        return <ArrangeRow key={index} item={item} />;
      })}
    </ThumbnailContainer>
  );
};

export default Thumbnail;
