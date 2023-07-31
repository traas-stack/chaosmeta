import { styled } from '@umijs/max';

export const Container = styled.div`
  padding-bottom: 24px;

  .shallow {
    color: rgba(0, 0, 0, 0.45);
  }
  .shallow-65 {
    color: rgba(0, 0, 0, 0.65);
  }
`;

export const TopStep = styled.div`
  .panel {
    display: flex;
    justify-content: space-between;
    .title {
      font-size: 20px;
    }
    .panel-state {
      font-size: 12px;
      .icon {
        width: 16px;
        height: 16px;
        border-radius: 50%;
        text-align: center;
        background-color: rgba(0, 0, 0, 0.1);
        cursor: pointer;
      }
    }
  }
  .card {
    margin-top: 16px;
    .ant-col {
      display: flex;
      .ant-card {
        flex: 1;
      }
    }
    .ant-card-body {
      .ant-space {
        align-items: flex-start;
      }
      .title {
        font-size: 16px;
        margin-bottom: 8px;
      }
      .desc {
        opacity: 0.45;
      }
      .buttons {
        margin-top: 16px;
      }
    }
  }
`;

export const SpaceContent = styled.div`
  margin-top: 16px;
  .ellipsis {
    position: relative;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    word-break: keep-all;
  }
  .left {
    .ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn {
      color: rgba(0, 0, 0, 0.85);
    }
    .ant-tabs-nav::before {
      border: none;
    }
    .ant-tabs-ink-bar {
      background-color: transparent;
      height: 0;
    }
    .row-text-gap {
      padding-bottom: 5px;
    }
  }
  .right {
    .top {
      display: flex;
      justify-content: space-between;
      align-items: center;
      .title {
        font-size: 16px;
        font-weight: 500;
        color: rgba(0, 0, 0, 0.85);
      }
      .ant-select-selector {
        width: 100px;
        .ant-select-selection-item {
          font-size: 12px;
          color: rgba(0, 0, 0, 0.65);
        }
      }
    }
    .overview {
      margin-bottom: 16px;
      .ant-select-selector {
        background-color: rgba(255, 255, 255, 0.5);
        margin-bottom: 8px;
      }
      .ant-card {
        .ant-card-body {
          height: 110px;
          padding: 0 24px;
        }
      }
      .result {
        height: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
        .count {
          font-size: 24px;
          font-weight: 500;
        }
        .count-error {
          font-size: 24px;
          font-weight: 500;
          color: #ff4d4f;
        }
        .ant-divider {
          position: absolute;
          left: -25px;
          height: 100%;
          /* background-color: rgba(0, 0, 0, 0.1); */
        }
      }
    }
    .recommend {
      .top {
        margin-bottom: 8px;
      }
      .item {
        display: flex;
        align-items: center;
        border-radius: 6px;
        border: 1px solid #e3e4e6;
        margin-top: 8px;
        padding: 12px;
        img {
          margin-right: 12px;
        }
      }
    }
  }
`;
