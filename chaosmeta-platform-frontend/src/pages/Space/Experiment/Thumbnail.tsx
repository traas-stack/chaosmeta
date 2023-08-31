import { arrangeNodeTypeColors } from '@/constants';
import { DroppableCol, DroppableRow, ThumbnailContainer } from './style';
const Thumbnail = () => {
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
      {[]?.map((item, index) => {
        return <ArrangeRow key={index} item={item} />;
      })}
    </ThumbnailContainer>
  );
};

export default Thumbnail;
