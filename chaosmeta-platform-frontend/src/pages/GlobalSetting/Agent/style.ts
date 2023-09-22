import { styled } from '@umijs/max';

export const Container = styled.div`
  padding-bottom: 24px;
  .ant-tabs-content {
    margin-top: -4px;
  }
  .ant-tabs-nav {
    background-color: #fff;
    margin: 0;
    padding: 0 16px;
    border-radius: 6px 6px 0 0;
    .ant-tabs-tab,
    .ant-tabs-nav-add {
      background-color: #fff;
      border: none;
      margin-top: 4px;
      .ant-tabs-tab-btn {
        margin: 0 4px;
      }
    }
    .ant-tabs-tab-active {
      background-image: url('https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ZsDqQJsUDgMAAAAAAAAAAAAADmKmAQ/original');
      background-size: 100% 38px;
      background-repeat: no-repeat;
    }
    .ant-tabs-tab-remove {
      margin: 0;
    }
  }
  .ant-alert {
    margin-bottom: 16px;
  }
  .search {
    padding: 16px;
    padding-bottom: 0;
    background-color: #fff;
    border-radius: 6px;
    margin-bottom: 16px;
    .ant-form-item {
      margin-bottom: 16px;
    }
    .ant-col:last-child {
      text-align: right;
    }
  }
  .title {
    font-weight: 500;
    font-size: 16px;
    color: rgba(0, 0, 0, 0.85);
    margin-bottom: 16px;
  }
`;

/**
 * 升级弹窗
 */
export const UpgradationContainer = styled.div`
  .version {
    display: flex;
    align-items: center;
    margin-bottom: 24px;
    font-weight: 500;
    color: rgba(0, 0, 0, 0.85);
    div {
      flex-shrink: 0;
    }
  }
`;
/**
 * 应用配置
 */
export const AppConfigContainer = styled.div`
  font-weight: 500;
  color: rgba(0, 0, 0, 0.85);
  .desc {
    margin-bottom: 24px;
  }
`;

export const InstallAgentContainer = styled.div`
  .title {
    margin: 24px 0;
  }
  .content {
    display: flex;
    height: calc(100vh - 200px);
    .transfer {
      button {
        width: 24px;
        height: 24px;
        display: flex;
        justify-content: center;
        align-items: center;
        margin: 0 8px;
        margin-top: 180px;
        padding: 0;
      }
    }
    .left,
    .right {
      width: 420px;
      border: 1px solid rgba(0, 10, 26, 0.16);
      border-radius: 6px;
      overflow-y: hidden;
      .header {
        /* display: flex; */
        height: 46px;
        line-height: 46px;
        text-align: right;
        padding: 0 12px;
        margin: 0 4px;
        border-bottom: 1px solid rgba(0, 10, 26, 0.16);
      }
      .search {
        .ant-space-compact {
          display: flex;
          .ant-form-item:last-child {
            flex: 1;
          }
        }
        padding: 16px;
        padding-bottom: 0;
        .ant-form-item {
          margin: 0;
        }
        .ant-select-selector {
          width: 110px;
        }
      }
      .table {
        padding: 16px 16px 0 16px;
        margin-bottom: 24px;
        .ant-table-body {
          padding-bottom: 48px;
        }
      }
    }
  }
`;
