import { styled } from '@umijs/max';

export const Container = styled.div`
  background-color: #fff;
  border-radius: 6px;
  padding: 16px 24px;
  .ant-tabs-nav::before {
    border: none;
  }
  .ant-tabs-nav {
    .ant-form-item {
      margin-bottom: 0;
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
  /* height: 175px; */
  background-color: #fff;
  border-radius: 6px;
  padding: 16px;
  padding-bottom: 0;
  border: 1px solid #d6d8da;
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
    margin: 24px 0;
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
    height: 38px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: rgba(0, 0, 0, 0.65);
    font-size: 12px;
    border-top: 1px solid rgba(0, 10, 26, 0.07);
    margin: 0 -16px;
    padding: 0 16px;
    img {
      vertical-align: sub;
      margin-right: 4px;
    }
  }
`;
