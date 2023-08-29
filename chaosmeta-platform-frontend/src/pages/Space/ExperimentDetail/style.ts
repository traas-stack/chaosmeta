import { styled } from '@umijs/max';

export const Container = styled.div`
  .content {
    background-color: #fff;
    padding: 16px 24px;
    margin-bottom: 24px;
    border-radius: 8px;
    font-size: 14px;
    color: rgba(0, 10, 26, 0.89);
    .ant-descriptions-title {
      font-weight: 500;
    }
    .experiment {
      .experiment-title {
        font-size: 16px;
        font-weight: 500;
        color: rgba(0, 0, 0, 0.8);
        margin: 16px 0;
      }
    }
  }
`;

/**
 * 编排展示区域
 */
export const ArrangeWrap = styled.div<{
  $activeItem?: boolean;
  $curExecSecond?: number;
  $widthSecond: number;
}>`
  .wrap {
    display: flex;
    background-color: #f7f7fa;
    border-radius: 6px;
    border: 1px solid rgba(0, 10, 26, 0.1);
    position: relative;
    .arrange {
      min-height: calc(100vh - 380px);
      padding-left: 38px;
      margin-bottom: 40px;
      flex: 1;
      overflow: auto;
      position: relative;
      .time-stage {
        height: 100%;
        position: absolute;
        z-index: 20;
        top: 2px;
        left: ${(props) => {
          return (
            props?.$curExecSecond &&
            `${props?.$curExecSecond * props?.$widthSecond + 30}px`
          );
        }};
        .line {
          height: 100%;
          width: 1px;
          background-color: #1890ff;
          transform: translateX(8px);
        }
      }
      /* 时间轴 */
      .time-axis {
        height: 40px;
        display: flex;
        background-color: #fafafc;
        line-height: 40px;
        font-size: 14px;
        border-bottom: 1px solid rgb(232, 237, 244);
        .time-item {
          position: relative;
          color: rgba(0, 10, 26, 0.47);

          font-weight: 400;
          flex-shrink: 0;
          background-color: #fafafc;
        }
        .time-item::before {
          content: '';
          position: absolute;
          display: block;
          bottom: 0;
          left: 0;
          width: 1px;
          height: 6px;
          background-color: rgb(216, 225, 235);
        }
      }
    }
    .info {
      width: 230px;
      padding: 16px;
      right: 0;
      top: 0;
      border: 1px solid rgba(0, 10, 26, 0.1);
      .subtitle {
        font-weight: 500;
        margin-bottom: 8px;
      }
      .range {
        margin-top: 6px;
      }
      .ant-form-item {
        margin-bottom: 0;
        font-weight: 400;
      }
    }
    .footer {
      height: 40px;
      padding: 0 16px;
      position: absolute;
      bottom: 0;
      width: ${(props) => {
        return props?.$activeItem ? `calc(100% - 206px)` : '100%';
      }};
      display: flex;
      color: #878c93;
      justify-content: space-between;
      align-items: center;
      border: 1px solid #dadde7;
      background-color: #fff;
      .total-time {
        display: inline-block;
        height: 22px;
        background-color: #f4f4f6;
        padding: 4px 12px;
        border-radius: 6px;
      }
      .node-type {
        margin-left: 46px;
        align-items: center;
      }
      .node-item {
        display: flex;
        div:first-child {
          width: 16px;
          height: 16px;
          display: inline-block;
          border-radius: 4px;
        }
      }
    }
  }
`;

/**
 * 编排区域行样式
 */
export const DroppableRow = styled.div<{
  $activeState?: boolean;
}>`
  height: 64px;
  width: 100%;
  box-sizing: border-box;
  border-bottom: 1px solid rgb(232, 237, 244);
  position: relative;
  cursor: pointer;
  border: ${(props) => {
    const { $activeState } = props;
    let border = '';
    if ($activeState) {
      border = '1px solid rgba(6,17,120,0.40)';
    }
    return border;
  }};
  background-color: #f5f6f9;
  .row {
    height: 100%;
    width: 100%;
    padding: 2px;
    display: flex;
    align-items: center;
    opacity: ${(props: any) => {
      return props?.$isDragging ? 0 : 1;
    }};
  }
  .handle {
    height: 52px;
    width: 24px;
    font-size: 14px;
    color: rgba(0, 10, 26, 0.47);
    position: absolute;
    left: -30px;
    top: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(0, 10, 26, 0.26);
    border-radius: 6px;
    background-color: #fff;
  }
`;

/**
 * 编排区域行内子节点样式
 */
export const DroppableCol = styled.div<{
  $bg?: string;
  $activeState?: boolean;
  $nodeStutas?: string;
}>`
  height: 100%;
  position: relative;
  box-sizing: border-box;
  overflow: hidden;
  font-size: 14px;
  margin-right: 2px;
  color: rgba(0, 10, 26, 0.89);
  border: ${(props) => {
    let border = 'none';
    if (props?.$activeState) {
      border = '2px solid #597EF7';
    }
    if (props?.$nodeStutas === 'Failed' || props?.$nodeStutas === 'error') {
      border = '2px solid #FF4D4F';
    }
    return border;
  }};
  background-color: ${(props) => {
    let color = '#D4E3F1';
    color = props?.$bg ?? color;
    return color;
  }};
  border-radius: 6px;
  .item {
    width: 100%;
    height: 100%;
    margin: 6px;
    .title {
      margin: 8px 0 6px 0;
    }
  }
`;
