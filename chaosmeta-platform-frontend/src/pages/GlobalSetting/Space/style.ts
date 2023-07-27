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
    .ant-tabs-tab {
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
  }
  .ant-pagination {
    text-align: right;
    margin-top: 16px;
  }
  .ant-alert {
    margin-bottom: 16px;
  }
`;

export const SpaceCard = styled.div`
  height: 175px;
  background-color: #fff;
  border-radius: 6px;
  padding: 16px;
  .header {
    display: flex;
    justify-content: space-between;
    .title {
      font-weight: 500;
    }
    .time {
      font-size: 12px;
      color: rgba(0, 10, 26, 0.45);
    }
    img {
      width: 32px;
      height: 32px;
      margin-right: 12px;
    }
    div:first-child {
      display: flex;
      align-items: center;
    }
    .anticon-ellipsis {
      cursor: pointer;
      height: 16px;
    }
  }
  .desc {
    margin: 16px 0;
    display: flex;
    font-size: 12px;
    div {
      flex-shrink: 0;
      color: rgba(0, 0, 0, 0.45);
    }

    span {
      color: rgba(0, 0, 0, 0.65);
      position: relative;
      overflow: hidden;
      text-overflow: ellipsis;
      word-break: keep-all;
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 2;
    }
  }
  .footer {
    display: flex;
    justify-content: space-between;
    height: 36px;
    color: rgba(0, 0, 0, 0.65);
    font-size: 12px;
    line-height: 36px;
    border-top: 1px solid rgba(0, 10, 26, 0.07);
    div:first-child {
      margin-right: 24px;
    }
    img {
      vertical-align: sub;
      margin-right: 4px;
    }
  }
`;
// https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ZsDqQJsUDgMAAAAAAAAAAAAADmKmAQ/original
