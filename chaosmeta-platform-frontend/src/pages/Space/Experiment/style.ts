import { styled } from '@umijs/max';

export const Container = styled.div`
  .ant-tabs-nav::before {
    border: none;
  }
  /* 实验列表 */
  .experiment-list {
    .search {
      padding-bottom: 0;
      margin-bottom: 24px;
    }
    .table {
      background-color: #fff;
      padding: 16px;
      border-radius: 6px;
    }
    .ellipsis {
      position: relative;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      word-break: keep-all;
    }
    .tag-tip {
      .ant-tooltip-inner {
        background-color: red !important;
      }
    }
    .cycle {
      font-size: 12px;
      color: rgba(0, 0, 0, 0.45);
    }
    .run-finish {
      color: rgba(0, 0, 0, 0.26);
    }
  }
  /* 推荐实验 */
  .recommend {
    background-color: #fff;
    border-radius: 6px;
    /* padding: 0 16px; */
    display: flex;
    .left {
      width: 210px;
      border-right: 1px solid #e3e4e6;
      /* padding-right: 16px; */

      /* 推荐实验下分类tab */
      .tab {
        .tab-item {
          position: relative;
          background-color: #fff;
          text-align: center;
          border-bottom: 1px solid #e3e4e6;
          padding: 8px 0;
          border-radius: 6px 6px 0 0;
          cursor: pointer;
          /* transition: 0.2s; */
        }
        .tab-item-active {
          border-bottom: none;
          color: #1677ff;
        }
        .tab-item-active-border-left {
          border-left: 1px solid #e3e4e6;
        }
        .tab-item-active-border-right {
          border-right: 1px solid #e3e4e6;
        }
        /* 底部外圆角实现 */
        .tab-item-active-before::before,
        .tab-item-active-after::after {
          position: absolute;
          z-index: 999;
          bottom: 0;
          content: '';
          width: 20px;
          height: 20px;
          border-radius: 100%;
          box-shadow: 0 0 0 40px #fff;
          transition: 0.2s;
          border: 1px solid #e3e4e6;
        }
        .tab-item-active-before::before {
          left: -20px;
          clip-path: inset(50% -10px 0 50%);
        }
        .tab-item-active-after::after {
          right: -20px;
          clip-path: inset(50% 50% 0 -10px);
        }
        .tab-item-content::before {
          content: '';
          position: absolute;
          left: 0;
          height: 24px;
          width: 1px;
          background-color: #e3e4e6;
        }
      }
      /* tab下分类 */
      .group {
        padding: 16px;
        .group-item {
          height: 40px;
          line-height: 40px;
          padding-left: 16px;
          border-radius: 6px;
          cursor: pointer;
        }
        .group-item:hover {
          background-color: rgba(0, 0, 0, 0.04);
        }
        .group-item-active {
          background-color: rgba(0, 0, 0, 0.06);
        }
      }
    }
    .right {
      flex: 1;
      padding: 16px 24px;
      height: calc(100vh - 240px);
      overflow-y: auto;
      .tip {
        margin-bottom: 16px;
      }
      .title {
        font-size: 16px;
        font-weight: 500;
        color: rgba(0, 0, 0, 0.85);
        margin-bottom: 16px;
      }
      .result-item {
        border-radius: 6px;
        border: 1px solid #edeeef;
        max-width: 288px;
        padding: 16px;
      }
      .introduce {
        display: flex;
        align-items: flex-start;
        margin-top: 16px;
        img {
          flex-shrink: 0;
          margin-right: 8px;
        }
      }
    }
  }
`;
