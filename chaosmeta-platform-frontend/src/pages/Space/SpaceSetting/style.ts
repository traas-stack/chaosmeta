import { styled } from '@umijs/max';

export const Container = styled.div`
  .ant-tabs-nav::before {
    border: none;
  }
  /* .ant-form-item {
    margin-bottom: 0;
  } */
  /* .tab-content {
    background-color: #fff;
    border-radius: 6px;
    min-height: calc(100vh - 320px);
    .tag-content {
      padding: 35px 48px;
    }
    td {
      .ant-select-selection-item,
      .ant-select-arrow {
        color: #1677ff;
      }
    }
  } */
`;

// 基本信息
export const BasicInfoContainer = styled.div`
  background-color: #fff;
  padding: 24px;
  border-radius: 6px;
  .ant-form {
    width: 50%;
    min-width: 600px;
  }
`