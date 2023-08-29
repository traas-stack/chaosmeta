import { styled } from '@umijs/max';

export const Container = styled.div`
  height: calc(100vh - 130px);
  background-color: #fff;
  border-radius: 8px;
  padding: 24px;
  .ant-tabs-nav::before {
    border: none;
  }
  .ant-tabs-tab {
    width: 58px;
    justify-content: center;
  }
  .custom {
    height: 107px;
    width: 217px;
    display: flex;
    justify-content: center;
    align-items: center;
    text-align: center;
    background-color: #ffffff;
    border: 2px dashed rgba(0, 10, 26, 0.16);
    border-radius: 6px;
    margin: 16px 0 24px 0;
    color: rgba(0, 10, 26, 0.68);
    cursor: pointer;
    svg {
      width: 20px;
      height: 32px;
    }
  }
  .content {
    display: flex;
    border: 1px solid #edeeef;
    border-radius: 6px;
    .left {
      width: 210px;
      border-right: 1px solid #edeeef;
      padding: 16px 0;
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
        background-color: rgba(0, 145, 255, 0.1);
      }
    }
    .right {
      flex: 1;
      padding: 16px 24px;
      height: calc(100vh - 430px);
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
