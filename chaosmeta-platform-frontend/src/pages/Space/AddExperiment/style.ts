import { styled } from '@umijs/max';

export const Container = styled.div`
  font-size: 14px;
  color: rgba(0, 0, 0, 0.65);
  .ellipsis {
    position: relative;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    word-break: keep-all;
  }
  /* 头部样式 */
  .ant-page-header {
    height: 72px;
    padding: 0 24px;
    border-bottom: 1px solid #dadde7;
    box-shadow: inset 0 -1px 0 0 rgba(6, 17, 120, 0.1);
    background-color: #fff;
    z-index: 9999;
    /* title部分 */
    .ant-page-header-heading {
      .ant-form-item {
        margin: 0;
        color: rgba(0, 0, 0, 0.65);
      }
      .ant-form-item-label {
        label {
          color: rgba(0, 0, 0, 0.65);
        }
      }
      .ant-page-header-heading-title {
        width: 300px;
        color: rgba(4, 24, 94, 0.85);
        font-size: 16px;
        font-weight: 500;
        .ellipsis {
          width: 260px;
        }
        .cancel {
          color: #ff4d4f;
        }
        .edit,
        .confirm {
          color: #1890ff;
        }
        .tags {
          font-weight: 400;
          span {
            color: rgba(0, 10, 26, 0.47);
          }
        }
      }
      /* 右侧展示部分 */
      .ant-page-header-heading-extra {
        flex: 1;
      }
      .ant-page-header-heading-extra > .ant-space,
      .ant-page-header-heading-extra > .ant-space > .ant-space-item {
        width: 100%;
      }
      .header-extra {
        flex: 1;
        display: flex;
        justify-content: space-between;
      }
    }
  }
  .ant-pro-page-container-children-container {
    padding: 0;
  }
  /* 内容区域 */
  .content {
    display: flex;
    position: relative;
  }
`;

/**
 * 节点库样式
 */
export const NodeLibraryContainer = styled.div`
  .wrap {
    min-width: 220px;
    position: relative;
    background-color: #f7f9fb;
    min-height: calc(100vh - 72px);
    height: 100%;
    box-shadow: 0 2px 4px 0 rgba(4, 24, 94, 0.04),
      0 4px 8px 0 rgba(4, 24, 94, 0.06);
    .title {
      height: 40px;
      background-color: #fff;
      font-size: 14px;
      color: #293a76;
      line-height: 40px;
      padding-left: 12px;
      border-bottom: 1px solid #dadde7;
    }
    .node {
      background-color: #f7f9fb;
      padding: 12px;
      .ant-input-affix-wrapper {
        background-color: #eef0f5;
        margin-bottom: 12px;
      }
      /* 外层折叠面板header */
      .ant-collapse-header {
        padding: 4px;
        color: #293a76;
      }
      .collapse-second {
        /* 内层层折叠面板header */
        .ant-collapse-header {
          padding: 8px 16px;
        }
      }
      .ant-collapse-content-box {
        padding: 0;
      }
      /* 节点库下拖拽节点 */
      .temp-item {
        padding: 6px 16px;
        background-color: #fff;
        border-radius: 6px;
        border: 1px solid #dadde7;
        box-shadow: 0 2px 4px 0 rgba(4, 24, 94, 0.04);
        margin-bottom: 8px;
        font-weight: 400;
        color: #293a76;
        cursor: pointer;
        img {
          vertical-align: middle;
          margin-right: 6px;
          margin-top: -4px;
        }
      }
      .temp-item:hover {
        border: 1px solid #1890ff;
      }
    }
    .fold-icon {
      position: absolute;
      right: -12px;
      top: calc(50% - 40px);
      cursor: pointer;
      span {
        transform: translatex(12px);
      }
    }
  }
`;

/**
 * 编排区域样式
 */
export const ArrangeContainer = styled.div<{ $activeColState?: boolean }>`
  flex: 1;
  overflow: hidden;
  /* position: absolute; */
  .flow {
    /* height: calc(100vh - 115px); */
    /* width: 100%; */
    height: calc(100% - 40px);
    overflow: auto;
    border-left: 1px solid rgba(0, 10, 26, 0.26);
    margin-left: 52px;
    /* 时间轴 */
    .time-axis {
      height: 40px;
      display: flex;
      background-color: #fafafc;
      line-height: 40px;
      font-size: 14px;
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
    /* 编排内容区域 */
    .center-content {
      min-height: calc(100vh - 176px);
      height: 100%;
      /* 背景行渐变，方格 */
      background-repeat: repeat-y;
      background-image: linear-gradient(
        0deg,
        transparent 63px,
        rgb(232, 237, 244) 63px,
        rgb(232, 237, 244) 64px,
        transparent 64px
      );
      background-size: 100% 64px;
      .row {
        height: 64px;
      }
    }
  }
  .footer {
    width: calc(
      100vw - 220px -
        (
          ${(props) => {
            console.log(props?.$activeColState, 'props?.$activeColState---');
            return props?.$activeColState ? '280px' : '0px';
          }}
        )
    );
    height: 40px;
    padding: 0 16px;
    position: fixed;
    bottom: 0;
    display: flex;
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
`;

/**
 * 节点信息配置样式
 */
export const NodeConfigContainer = styled.div`
  width: 280px;
  flex-shrink: 0;
  border-left: 1px solid #e2e3ee;
  padding: 0 12px;
  background-color: #fff;
  .header {
    height: 40px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin: 0 -12px;
    padding: 0 12px;
    border-bottom: 1px solid #e2e3ee;
    margin-bottom: 16px;
  }
`;

/**
 * 节点库单个节点样式
 */
export const NodeItem = styled.div<{ $isDragging?: boolean }>`
  /* 节点库下拖拽节点 */
  background-color: ${(props) => {
    return props?.$isDragging ? '#edeff2' : '#fff';
  }};
  border-radius: 6px;
  .temp-item {
    opacity: ${(props: any) => {
      return props?.$isDragging ? 0 : 1;
    }};
    font-size: 14px;
    padding: 6px 16px;
    border: 1px solid #dadde7;
    border-radius: 6px;
    box-shadow: 0 2px 4px 0 rgba(4, 24, 94, 0.04);
    margin-bottom: 8px;
    font-weight: 400;
    color: #293a76;
    cursor: pointer;
    img {
      vertical-align: middle;
      margin-right: 6px;
      margin-top: -4px;
    }
  }
  .temp-item:hover {
    border: 1px solid #1890ff;
  }
`;

/**
 * 行的拖动手柄，默认状态，row下还有一个手柄样式，是拖动起的样式，两个需要分离开，避免拖动时样式错乱
 */
export const HandleMove = styled.div<{
  $index?: number;
  $activeState?: boolean;
  $hoverState?: boolean;
  $scrollTop?: number;
}>`
  display: flex;
  position: absolute;
  z-index: 999;
  left: 240px;
  cursor: pointer;
  transition: all 0.3s;
  top: ${(props) => {
    return `${((props?.$index || 0) + 1) * 64 - 16 - props?.$scrollTop}px`;
  }};
  .delete {
    position: absolute;
    top: 18px;
    left: -16px;
    color: #ff4d4f;
  }
  .handle {
    height: 52px;
    width: 24px;
    font-size: 14px;
    margin-left: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(0, 10, 26, 0.26);
    border-radius: 6px;
    background-color: ${(props) => {
      let color = '#fff';
      if (props?.$activeState) {
        color = '#1890ff';
        if (props?.$hoverState) {
          color = '#ef7b77';
        }
      }
      return color;
    }};
    color: ${(props) => {
      return props?.$activeState ? '#fff' : '#000';
    }};
  }
`;

/**
 * 编排区域行样式
 */
export const DroppableRow = styled.div<{
  $isDragging?: boolean;
  $isMoving?: boolean;
  $transform?: any;
  $offsetTop?: number;
  $index?: number;
  $activeState?: boolean;
  $hoverState?: boolean;
}>`
  height: 64px;
  box-sizing: border-box;
  border-top: 1px solid rgb(232, 237, 244);
  cursor: pointer;
  border: ${(props) => {
    const { $isMoving, $activeState, $hoverState } = props;
    let border = '';
    if ($isMoving || $activeState) {
      border = '1px solid rgba(6,17,120,0.40)';
    }
    if ($hoverState && $activeState) {
      border = '1px solid #f4a19d';
    }
    return border;
  }};
  background-color: ${(props) => {
    const { $isDragging, $isMoving, $activeState, $hoverState } = props;
    let color = '#f5f6f9';
    if ($isMoving || $activeState) {
      color = '#E6E8EF';
      if ($hoverState) {
        color = '#fdeeee';
      }
    }
    if ($isDragging) {
      color = '#e5e6e8';
    }
    return color;
  }};
  transition: 0.5s;
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
  /* 拖动手柄，拖动起的状态 */
  .handle {
    height: 52px;
    width: 24px;
    font-size: 14px;
    position: absolute;
    left: -30px;
    top: 6px;
    z-index: 9999;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(0, 10, 26, 0.26);
    border-radius: 6px;
    background-color: #1890ff;
    color: #fff;
  }
  /* 拖动时行下方的占位蓝条 */
  .moveing {
    height: 2px;
    width: 100%;
    background-color: #1890ff;
    position: absolute;
    z-index: 9;
    display: ${(props) => {
      return props?.$isDragging ? 'block' : 'none';
    }};
    transition: 0.5s;
    transform: ${(props) => {
      const { $transform } = props;
      return `translate3d(${Math.round($transform?.x)}px, ${Math.round(
        $transform?.y,
      )}px, 0)`;
    }};
  }
`;

/**
 * 编排区域行内子节点样式
 */
export const DroppableCol = styled.div<{
  $isDragging?: boolean;
  $bg?: string;
  $transform?: any;
  $activeState?: boolean;
}>`
  height: 100%;
  position: relative;
  box-sizing: border-box;
  font-size: 14px;
  margin-right: 2px;
  padding: 6px;
  border: ${(props) => {
    let border = 'none';
    if (props?.$activeState) {
      border = '2px solid #597EF7';
    }
    if (props?.$isDragging) {
      border = '2px dashed #597EF7';
    }
    return border;
  }};
  background-color: ${(props) => {
    let color = '#fff';
    color = props?.$bg ?? color;
    if (props?.$isDragging) {
      color = '#d2e1f8';
    }
    return color;
  }};
  transition: 0.5s;
  transform: ${(props) => {
    const { $transform } = props;
    return `translate3d(${Math.round($transform?.x)}px, ${Math.round(
      $transform?.y,
    )}px, 0)`;
  }};
  border-radius: 6px;
  .item {
    width: 100%;
    height: 100%;
    opacity: ${(props) => {
      return props?.$isDragging ? 0 : 1;
    }};
  }
  /* 拖拽宽度 */
  .scale {
    position: absolute;
    right: -4px;
    top: 16px;
    height: 24px;
    width: 6px;
    background-color: #597ef7;
    border-radius: 3px;
    cursor: col-resize;
    
  }
`;
